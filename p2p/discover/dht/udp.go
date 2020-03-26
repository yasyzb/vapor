package dht

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"path"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tendermint/go-wire"

	"github.com/bytom/vapor/common"
	cfg "github.com/bytom/vapor/config"
	"github.com/bytom/vapor/crypto"
	"github.com/bytom/vapor/p2p/netutil"
	"github.com/bytom/vapor/p2p/signlib"
	"github.com/bytom/vapor/version"
)

const (
	//Version dht discover protocol version
	Version   = 5
	logModule = "discover"
)

// Errors
var (
	errPacketTooSmall = errors.New("too small")
	errPrefixMismatch = errors.New("prefix mismatch")
	errNetIDMismatch  = errors.New("network id mismatch")
	errPacketType     = errors.New("unknown packet type")
)

// Timeouts
const (
	respTimeout = 1 * time.Second
	expiration  = 20 * time.Second
)

// ReadPacket is sent to the unhandled channel when it could not be processed
type ReadPacket struct {
	Data []byte
	Addr *net.UDPAddr
}

// Config holds Table-related settings.
type Config struct {
	// These settings are required and configure the UDP listener:
	PrivateKey *ecdsa.PrivateKey

	// These settings are optional:
	AnnounceAddr *net.UDPAddr // local address announced in the DHT
	NodeDBPath   string       // if set, the node database is stored at this filesystem location
	//NetRestrict  *netutil.Netlist  // network whitelist
	Bootnodes []*Node           // list of bootstrap nodes
	Unhandled chan<- ReadPacket // unhandled packets are sent on this channel
}

// RPC request structures
type (
	ping struct {
		Version    uint
		From, To   rpcEndpoint
		Expiration uint64

		// v5
		Topics []Topic

		// Ignore additional fields (for forward compatibility).
		Rest []byte
	}

	// pong is the reply to ping.
	pong struct {
		// This field should mirror the UDP envelope address
		// of the ping packet, which provides a way to discover the
		// the external address (after NAT).
		To rpcEndpoint

		ReplyTok   []byte // This contains the hash of the ping packet.
		Expiration uint64 // Absolute timestamp at which the packet becomes invalid.

		// v5
		TopicHash    common.Hash
		TicketSerial uint32
		WaitPeriods  []uint32

		// Ignore additional fields (for forward compatibility).
		Rest []byte
	}

	// findnode is a query for nodes close to the given target.
	findnode struct {
		Target     NodeID // doesn't need to be an actual public key
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte
	}

	// findnode is a query for nodes close to the given target.
	findnodeHash struct {
		Target     common.Hash
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte
	}

	// reply to findnode
	neighbors struct {
		Nodes      []rpcNode
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte
	}

	topicRegister struct {
		Topics []Topic
		Idx    uint
		Pong   []byte
	}

	topicQuery struct {
		Topic      Topic
		Expiration uint64
	}

	// reply to topicQuery
	topicNodes struct {
		Echo  common.Hash
		Nodes []rpcNode
	}

	rpcNode struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
		ID  NodeID
	}

	rpcEndpoint struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
	}
)

var (
	netIDSize  = 8
	nodeIDSize = 32
	sigSize    = 520 / 8
	headSize   = netIDSize + nodeIDSize + sigSize // space of packet frame data
)

// Neighbors replies are sent across multiple packets to
// stay below the 1280 byte limit. We compute the maximum number
// of entries by stuffing a packet until it grows too large.
var maxNeighbors = func() int {
	p := neighbors{Expiration: ^uint64(0)}
	maxSizeNode := rpcNode{IP: make(net.IP, 16), UDP: ^uint16(0), TCP: ^uint16(0)}
	for n := 0; ; n++ {
		p.Nodes = append(p.Nodes, maxSizeNode)
		var size int
		var err error
		b := new(bytes.Buffer)
		wire.WriteJSON(p, b, &size, &err)
		if err != nil {
			// If this ever happens, it will be caught by the unit tests.
			panic("cannot encode: " + err.Error())
		}
		if headSize+size+1 >= 1280 {
			return n
		}
	}
}()

var maxTopicNodes = func() int {
	p := topicNodes{}
	maxSizeNode := rpcNode{IP: make(net.IP, 16), UDP: ^uint16(0), TCP: ^uint16(0)}
	for n := 0; ; n++ {
		p.Nodes = append(p.Nodes, maxSizeNode)
		var size int
		var err error
		b := new(bytes.Buffer)
		wire.WriteJSON(p, b, &size, &err)
		if err != nil {
			// If this ever happens, it will be caught by the unit tests.
			panic("cannot encode: " + err.Error())
		}
		if headSize+size+1 >= 1280 {
			return n
		}
	}
}()

func makeEndpoint(addr *net.UDPAddr, tcpPort uint16) rpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return rpcEndpoint{IP: ip, UDP: uint16(addr.Port), TCP: tcpPort}
}

func (e1 rpcEndpoint) equal(e2 rpcEndpoint) bool {
	return e1.UDP == e2.UDP && e1.TCP == e2.TCP && e1.IP.Equal(e2.IP)
}

func nodeFromRPC(sender *net.UDPAddr, rn rpcNode) (*Node, error) {
	if err := netutil.CheckRelayIP(sender.IP, rn.IP); err != nil {
		return nil, err
	}
	n := NewNode(rn.ID, rn.IP, rn.UDP, rn.TCP)
	err := n.validateComplete()
	return n, err
}

func nodeToRPC(n *Node) rpcNode {
	return rpcNode{ID: n.ID, IP: n.IP, UDP: n.UDP, TCP: n.TCP}
}

type ingressPacket struct {
	remoteID   NodeID
	remoteAddr *net.UDPAddr
	ev         nodeEvent
	hash       []byte
	data       interface{} // one of the RPC structs
	rawData    []byte
}

type conn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

type netWork interface {
	reqReadPacket(pkt ingressPacket)
	selfIP() net.IP
}

// udp implements the RPC protocol.
type udp struct {
	conn conn
	priv signlib.PrivKey
	//netID used to isolate subnets
	netID       uint64
	ourEndpoint rpcEndpoint
	net         netWork
}

//NewDiscover create new dht discover
func NewDiscover(config *cfg.Config, privKey signlib.PrivKey, port uint16, netID uint64) (*Network, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("0.0.0.0", strconv.FormatUint(uint64(port), 10)))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)
	ntab, err := ListenUDP(privKey, conn, realaddr, path.Join(config.DBDir(), "discover"), nil, netID)
	if err != nil {
		return nil, err
	}
	seeds, err := QueryDNSSeeds(net.LookupHost)
	if err != nil {
		log.WithFields(log.Fields{"module": logModule, "err": err}).Error("fail on query dns seeds")
	}

	codedSeeds := netutil.CheckAndSplitAddresses(config.P2P.Seeds)
	seeds = append(seeds, codedSeeds...)
	if len(seeds) == 0 {
		return ntab, nil
	}

	var nodes []*Node
	for _, seed := range seeds {
		version.Status.AddSeed(seed)
		url := "enode://" + hex.EncodeToString(crypto.Sha256([]byte(seed))) + "@" + seed
		nodes = append(nodes, MustParseNode(url))
	}

	if err = ntab.SetFallbackNodes(nodes); err != nil {
		return nil, err
	}
	return ntab, nil
}

// ListenUDP returns a new table that listens for UDP packets on laddr.
func ListenUDP(privKey signlib.PrivKey, conn conn, realaddr *net.UDPAddr, nodeDBPath string, netrestrict *netutil.Netlist, netID uint64) (*Network, error) {
	transport, err := listenUDP(privKey, conn, realaddr, netID)
	if err != nil {
		return nil, err
	}

	net, err := newNetwork(transport, privKey.XPub(), nodeDBPath, netrestrict)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"module": logModule, "net": net.tab.self}).Info("UDP listener up v5")
	transport.net = net
	go transport.readLoop()
	return net, nil
}

func listenUDP(priv signlib.PrivKey, conn conn, realaddr *net.UDPAddr, netID uint64) (*udp, error) {
	return &udp{conn: conn, priv: priv, netID: netID, ourEndpoint: makeEndpoint(realaddr, uint16(realaddr.Port))}, nil
}

func (t *udp) localAddr() *net.UDPAddr {
	return t.conn.LocalAddr().(*net.UDPAddr)
}

func (t *udp) Close() {
	t.conn.Close()
}

func (t *udp) send(remote *Node, ptype nodeEvent, data interface{}) (hash []byte) {
	hash, _ = t.sendPacket(remote.ID, remote.addr(), byte(ptype), data)
	return hash
}

func (t *udp) sendPing(remote *Node, toaddr *net.UDPAddr, topics []Topic) (hash []byte) {
	hash, _ = t.sendPacket(remote.ID, toaddr, byte(pingPacket), ping{
		Version:    Version,
		From:       t.ourEndpoint,
		To:         makeEndpoint(toaddr, uint16(toaddr.Port)), // TODO: maybe use known TCP port from DB
		Expiration: uint64(time.Now().Add(expiration).Unix()),
		Topics:     topics,
	})
	return hash
}

func (t *udp) sendFindnode(remote *Node, target NodeID) {
	t.sendPacket(remote.ID, remote.addr(), byte(findnodePacket), findnode{
		Target:     target,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
}

func (t *udp) sendNeighbours(remote *Node, results []*Node) {
	// Send neighbors in chunks with at most maxNeighbors per packet
	// to stay below the 1280 byte limit.
	p := neighbors{Expiration: uint64(time.Now().Add(expiration).Unix())}
	for i, result := range results {
		p.Nodes = append(p.Nodes, nodeToRPC(result))
		if len(p.Nodes) == maxNeighbors || i == len(results)-1 {
			t.sendPacket(remote.ID, remote.addr(), byte(neighborsPacket), p)
			p.Nodes = p.Nodes[:0]
		}
	}
}

func (t *udp) sendFindnodeHash(remote *Node, target common.Hash) {
	t.sendPacket(remote.ID, remote.addr(), byte(findnodeHashPacket), findnodeHash{
		Target:     target,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
}

func (t *udp) sendTopicRegister(remote *Node, topics []Topic, idx int, pong []byte) {
	t.sendPacket(remote.ID, remote.addr(), byte(topicRegisterPacket), topicRegister{
		Topics: topics,
		Idx:    uint(idx),
		Pong:   pong,
	})
}

func (t *udp) sendTopicNodes(remote *Node, queryHash common.Hash, nodes []*Node) {
	p := topicNodes{Echo: queryHash}
	var sent bool
	for _, result := range nodes {
		if result.IP.Equal(t.net.selfIP()) || netutil.CheckRelayIP(remote.IP, result.IP) == nil {
			p.Nodes = append(p.Nodes, nodeToRPC(result))
		}
		if len(p.Nodes) == maxTopicNodes {
			t.sendPacket(remote.ID, remote.addr(), byte(topicNodesPacket), p)
			p.Nodes = p.Nodes[:0]
			sent = true
		}
	}
	if !sent || len(p.Nodes) > 0 {
		t.sendPacket(remote.ID, remote.addr(), byte(topicNodesPacket), p)
	}
}

func (t *udp) sendPacket(toid NodeID, toaddr *net.UDPAddr, ptype byte, req interface{}) (hash []byte, err error) {
	packet, hash, err := encodePacket(t.priv, ptype, req, t.netID)
	if err != nil {
		return hash, err
	}
	log.WithFields(log.Fields{"module": logModule, "event": nodeEvent(ptype), "to id": hex.EncodeToString(toid[:8]), "to addr": toaddr}).Debug("send packet")
	if _, err = t.conn.WriteToUDP(packet, toaddr); err != nil {
		log.WithFields(log.Fields{"module": logModule, "error": err}).Info(fmt.Sprint("UDP send failed"))
	}
	return hash, err
}

// zeroed padding space for encodePacket.
var headSpace = make([]byte, headSize)

func encodePacket(privKey signlib.PrivKey, ptype byte, req interface{}, netID uint64) (p, hash []byte, err error) {
	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(ptype)
	var size int
	wire.WriteJSON(req, b, &size, &err)
	if err != nil {
		log.WithFields(log.Fields{"module": logModule, "error": err}).Error("error encoding packet")
		return nil, nil, err
	}
	packet := b.Bytes()
	nodeID := privKey.XPub()
	sig := privKey.Sign(common.BytesToHash(packet[headSize:]).Bytes())
	id := []byte(strconv.FormatUint(netID, 16))
	copy(packet[:], id[:])
	copy(packet[netIDSize:], nodeID[:nodeIDSize])
	copy(packet[netIDSize+nodeIDSize:], sig)

	hash = common.BytesToHash(packet[:]).Bytes()
	return packet, hash, nil
}

// readLoop runs in its own goroutine. it injects ingress UDP packets
// into the network loop.
func (t *udp) readLoop() {
	defer t.conn.Close()
	// Discovery packets are defined to be no larger than 1280 bytes.
	// Packets larger than this size will be cut at the end and treated
	// as invalid because their hash won't match.
	buf := make([]byte, 1280)
	for {
		nbytes, from, err := t.conn.ReadFromUDP(buf)
		if netutil.IsTemporaryError(err) {
			// Ignore temporary read errors.
			log.WithFields(log.Fields{"module": logModule, "error": err}).Debug("Temporary read error")
			continue
		} else if err != nil {
			// Shut down the loop for permament errors.
			log.WithFields(log.Fields{"module": logModule, "error": err}).Debug("Read error")
			return
		}
		if err := t.handlePacket(from, buf[:nbytes]); err != nil {
			log.WithFields(log.Fields{"module": logModule, "from": from, "error": err}).Error("handle packet err")
		}
	}
}

func (t *udp) handlePacket(from *net.UDPAddr, buf []byte) error {
	pkt := ingressPacket{remoteAddr: from}
	if err := decodePacket(buf, &pkt, t.netID); err != nil {
		return err
	}
	t.net.reqReadPacket(pkt)
	return nil
}

func (t *udp) getNetID() uint64 {
	return t.netID
}

func decodePacket(buffer []byte, pkt *ingressPacket, netID uint64) error {
	if len(buffer) < headSize+1 {
		return errPacketTooSmall
	}
	buf := make([]byte, len(buffer))
	copy(buf, buffer)
	fromID, sigdata := buf[netIDSize:netIDSize+nodeIDSize], buf[headSize:]

	if !bytes.Equal(buf[:netIDSize], []byte(strconv.FormatUint(netID, 16))[:netIDSize]) {
		return errNetIDMismatch
	}

	pkt.rawData = buf
	pkt.hash = common.BytesToHash(buf[:]).Bytes()
	pkt.remoteID = ByteID(fromID)
	switch pkt.ev = nodeEvent(sigdata[0]); pkt.ev {
	case pingPacket:
		pkt.data = new(ping)
	case pongPacket:
		pkt.data = new(pong)
	case findnodePacket:
		pkt.data = new(findnode)
	case neighborsPacket:
		pkt.data = new(neighbors)
	case findnodeHashPacket:
		pkt.data = new(findnodeHash)
	case topicRegisterPacket:
		pkt.data = new(topicRegister)
	case topicQueryPacket:
		pkt.data = new(topicQuery)
	case topicNodesPacket:
		pkt.data = new(topicNodes)
	default:
		return errPacketType
	}
	var err error
	wire.ReadJSON(pkt.data, sigdata[1:], &err)
	if err != nil {
		log.WithFields(log.Fields{"module": logModule, "error": err}).Error("wire readjson err")
	}

	return err
}
