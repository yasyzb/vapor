package dht

import (
	"bytes"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/bytom/vapor/common"
	"github.com/bytom/vapor/errors"
	"github.com/bytom/vapor/p2p/signlib"
)

func TestPacketCodec(t *testing.T) {
	var testPackets = []struct {
		ptype      byte
		wantErr    error
		wantPacket interface{}
	}{
		{
			ptype:   byte(pingPacket),
			wantErr: nil,
			wantPacket: &ping{
				Version:    4,
				From:       rpcEndpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
				To:         rpcEndpoint{net.ParseIP("::1"), 2222, 3333},
				Expiration: 1136239445,
				Topics:     []Topic{"test topic"},
				Rest:       []byte{},
			},
		},
		{
			ptype:   byte(pingPacket),
			wantErr: nil,
			wantPacket: &ping{
				Version:    4,
				From:       rpcEndpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
				To:         rpcEndpoint{net.ParseIP("::1"), 2222, 3333},
				Expiration: 1136239445,
				Topics:     []Topic{"test topic"},
				Rest:       []byte{0x01, 0x02},
			},
		},
		{
			ptype:   byte(pingPacket),
			wantErr: nil,
			wantPacket: &ping{
				Version:    555,
				From:       rpcEndpoint{net.ParseIP("2001:db8:3c4d:15::abcd:ef12"), 3322, 5544},
				To:         rpcEndpoint{net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"), 2222, 33338},
				Expiration: 1136239445,
				Topics:     []Topic{"test topic"},
				Rest:       []byte{0xC5, 0x01, 0x02, 0x03, 0x04, 0x05},
			},
		},
		{
			ptype:   byte(pongPacket),
			wantErr: nil,
			wantPacket: &pong{
				To:          rpcEndpoint{net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"), 2222, 33338},
				ReplyTok:    []byte("fbc914b16819237dcd8801d7e53f69e9719adecb3cc0e790c57e91ca4461c954"),
				Expiration:  1136239445,
				WaitPeriods: []uint32{},
				Rest:        []byte{0xC6, 0x01, 0x02, 0x03, 0xC2, 0x04, 0x05, 0x06},
			},
		},
		{
			ptype:   byte(findnodePacket),
			wantErr: nil,
			wantPacket: &findnode{
				Target:     MustHexID("a2cb4c36765430f2e72564138c36f30fbc8af5a8bb91649822cd937dedbb8748"),
				Expiration: 1136239445,
				Rest:       []byte{0x82, 0x99, 0x99, 0x83, 0x99, 0x99, 0x99},
			},
		},
		{
			ptype:   byte(neighborsPacket),
			wantErr: nil,
			wantPacket: &neighbors{
				Nodes: []rpcNode{
					{
						ID:  MustHexID("a2cb4c36765430f2e72564138c36f30fbc8af5a8bb91649822cd937dedbb8748"),
						IP:  net.ParseIP("99.33.22.55").To4(),
						UDP: 4444,
						TCP: 4445,
					},
					{
						ID:  MustHexID("312c55512422cf9b8a4097e9a6ad79402e87a15ae909a4bfefa22398f03d2095"),
						IP:  net.ParseIP("1.2.3.4").To4(),
						UDP: 1,
						TCP: 1,
					},
					{
						ID:  MustHexID("38643200b172dcfef857492156971f0e6aa2c538d8b74010f8e140811d53b98c"),
						IP:  net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
						UDP: 3333,
						TCP: 3333,
					},
					{
						ID:  MustHexID("8dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2"),
						IP:  net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"),
						UDP: 999,
						TCP: 1000,
					},
				},
				Expiration: 1136239445,
				Rest:       []byte{0x01, 0x02, 0x03},
			},
		},
		{
			ptype:   byte(findnodeHashPacket),
			wantErr: nil,
			wantPacket: &findnodeHash{
				Target:     common.Hash{0x0, 0x1, 0x2, 0x3},
				Expiration: 1136239445,
				Rest:       []byte{0x01, 0x02, 0x03},
			},
		},
		{
			ptype:   byte(topicRegisterPacket),
			wantErr: nil,
			wantPacket: &topicRegister{
				Topics: []Topic{"test topic"},
				Idx:    uint(0x01),
				Pong:   []byte{0x01, 0x02, 0x03},
			},
		},
		{
			ptype:   byte(topicQueryPacket),
			wantErr: nil,
			wantPacket: &topicQuery{
				Topic:      "test topic",
				Expiration: 1136239445,
			},
		},
		{
			ptype:   byte(topicNodesPacket),
			wantErr: nil,
			wantPacket: &topicNodes{
				Echo: common.Hash{0x00, 0x01, 0x02},
				Nodes: []rpcNode{
					{
						ID:  MustHexID("a2cb4c36765430f2e72564138c36f30fbc8af5a8bb91649822cd937dedbb8748"),
						IP:  net.ParseIP("99.33.22.55").To4(),
						UDP: 4444,
						TCP: 4445,
					},
					{
						ID:  MustHexID("312c55512422cf9b8a4097e9a6ad79402e87a15ae909a4bfefa22398f03d2095"),
						IP:  net.ParseIP("1.2.3.4").To4(),
						UDP: 1,
						TCP: 1,
					},
					{
						ID:  MustHexID("38643200b172dcfef857492156971f0e6aa2c538d8b74010f8e140811d53b98c"),
						IP:  net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
						UDP: 3333,
						TCP: 3333,
					},
					{
						ID:  MustHexID("8dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2"),
						IP:  net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"),
						UDP: 999,
						TCP: 1000,
					},
				},
			},
		},
		{
			ptype:      byte(topicNodesPacket + 1),
			wantErr:    errPacketType,
			wantPacket: &topicNodes{},
		},
	}

	privateKey, _ := signlib.NewPrivKey()
	netID := uint64(0x12345)
	for i, test := range testPackets {
		packet, h, err := encodePacket(privateKey, test.ptype, test.wantPacket, netID)
		if err != nil {
			t.Fatal(err)
		}

		var pkt ingressPacket
		if err := decodePacket(packet, &pkt, netID); err != nil {
			if errors.Root(err) != test.wantErr {
				t.Errorf("index %d did not accept packet %s\n%v", i, packet, err)
			}
			continue
		}

		if !reflect.DeepEqual(pkt.hash, h) {
			t.Fatalf("packet hash err. got %x, want %x", pkt.hash, h)
		}

		if !reflect.DeepEqual(pkt.data, test.wantPacket) {
			t.Errorf("got %s\nwant %s", spew.Sdump(pkt.data), spew.Sdump(test.wantPacket))
		}
	}
}

type testConn struct {
	conn net.Conn
}

func (tc *testConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	n, err = tc.conn.Read(b)
	return n, nil, err
}

func (tc *testConn) WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error) {
	return tc.conn.Write(b)
}

func (tc *testConn) Close() error {
	return tc.conn.Close()
}

func (tc *testConn) LocalAddr() net.Addr {
	return tc.conn.LocalAddr()
}

type testNetWork struct {
	read chan ingressPacket // ingress packets arrive here
	IP   net.IP
}

func (tw *testNetWork) reqReadPacket(pkt ingressPacket) {
	tw.read <- pkt
}

func (tw *testNetWork) selfIP() net.IP {
	return tw.IP
}

func TestPacketTransport(t *testing.T) {
	c1, c2 := net.Pipe()
	inConn := &testConn{conn: c1}
	realaddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 40000}
	toAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 40000}
	inPrivKey, _ := signlib.NewPrivKey()
	outPrivKey, _ := signlib.NewPrivKey()
	netID := uint64(0x12345)

	udpInput, err := listenUDP(inPrivKey, inConn, realaddr, netID)
	if err != nil {
		t.Fatal(err)
	}
	node := &Node{ID: MustHexID("8dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2"),
		IP:  net.ParseIP("99.33.22.55").To4(),
		UDP: 4444,
		TCP: 4445,
	}

	udpInput.net = &testNetWork{read: make(chan ingressPacket, 100)}
	go udpInput.readLoop()

	outConn := &testConn{conn: c2}
	udp, err := listenUDP(outPrivKey, outConn, realaddr, netID)
	if err != nil {
		t.Fatal(err)
	}
	udp.net = &testNetWork{IP: node.IP}
	var hash []byte

	//test sendPing
	hash = udp.sendPing(node, toAddr, nil)
	pkts := receivePacket(udpInput)
	if !bytes.Equal(pkts[0].hash, hash) {
		t.Fatal("pingPacket transport err")
	}

	//test sendFindnodeHash
	target := common.Hash{0x01, 0x02}
	udp.sendFindnodeHash(node, target)
	pkts = receivePacket(udpInput)
	if !bytes.Equal(pkts[0].data.(*findnodeHash).Target.Bytes(), target.Bytes()) {
		t.Fatal("findnodeHashPacket transport err")
	}

	//test sendNeighbours
	nodes := []*Node{
		{
			ID:  MustHexID("a2cb4c36765430f2e72564138c36f30fbc8af5a8bb91649822cd937dedbb8748"),
			IP:  net.ParseIP("99.33.22.55").To4(),
			UDP: 4444,
			TCP: 4445,
		},
		{
			ID:  MustHexID("312c55512422cf9b8a4097e9a6ad79402e87a15ae909a4bfefa22398f03d2095"),
			IP:  net.ParseIP("1.2.3.4").To4(),
			UDP: 1,
			TCP: 1,
		},
		{
			ID:  MustHexID("38643200b172dcfef857492156971f0e6aa2c538d8b74010f8e140811d53b98c"),
			IP:  net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
			UDP: 3333,
			TCP: 3333,
		},
		{
			ID:  MustHexID("8dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2"),
			IP:  net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"),
			UDP: 999,
			TCP: 1000,
		},
	}

	udp.sendNeighbours(node, nodes)
	pkts = receivePacket(udpInput)
	var gotNodes []rpcNode
	for _, pkt := range pkts {
		gotNodes = append(gotNodes, pkt.data.(*neighbors).Nodes[:]...)
	}
	for i := 0; i < len(nodes); i++ {
		if !reflect.DeepEqual(nodeToRPC(nodes[i]), gotNodes[i]) {
			t.Fatal("sendNeighboursPacket transport err")
		}
	}

	//test sendFindnode
	targetNode := NodeID{0x01, 0x02, 0x03}
	udp.sendFindnode(node, targetNode)
	pkts = receivePacket(udpInput)
	if pkts[0].data.(*findnode).Target != targetNode {
		t.Fatal("sendFindnode transport err")
	}

	//test sendTopicRegister
	topics := []Topic{"topic1", "topic2", "topic3"}
	idx := 0xff
	pong := []byte{0x01, 0x02, 0x03}
	udp.sendTopicRegister(node, topics, idx, pong)
	pkts = receivePacket(udpInput)
	if !bytes.Equal(pkts[0].data.(*topicRegister).Pong, pong) {
		t.Fatal("sendTopicRegister pong field err")
	}
	if pkts[0].data.(*topicRegister).Idx != uint(idx) {
		t.Fatal("sendTopicRegister idx field err")
	}
	if !reflect.DeepEqual(pkts[0].data.(*topicRegister).Topics, topics) {
		t.Fatal("sendTopicRegister topic field err")
	}

	//test sendTopicNodes
	queryHash := common.Hash{0x01, 0x02, 0x03}
	udp.sendTopicNodes(node, queryHash, nodes)
	pkts = receivePacket(udpInput)
	gotNodes = gotNodes[:0]
	for _, pkt := range pkts {
		gotNodes = append(gotNodes, pkt.data.(*topicNodes).Nodes[:]...)
	}

	for i := 0; i < 2; i++ {
		if !reflect.DeepEqual(nodeToRPC(nodes[i]), gotNodes[i]) {
			t.Fatal("sendTopicNodes node field err")
		}
	}

	if pkts[0].data.(*topicNodes).Echo != queryHash {
		t.Fatal("sendTopicNodes echo field err")
	}
}

func TestSendTopicNodes(t *testing.T) {
	c1, c2 := net.Pipe()
	inConn := &testConn{conn: c1}
	realaddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 40000}
	inPrivKey, _ := signlib.NewPrivKey()
	outPrivKey, _ := signlib.NewPrivKey()
	netID := uint64(0x12345)

	udpInput, err := listenUDP(inPrivKey, inConn, realaddr, netID)
	if err != nil {
		t.Fatal(err)
	}
	node := &Node{ID: MustHexID("8dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2"),
		IP:  net.ParseIP("99.33.22.55").To4(),
		UDP: 4444,
		TCP: 4445,
	}

	udpInput.net = &testNetWork{read: make(chan ingressPacket, 100)}
	go udpInput.readLoop()

	outConn := &testConn{conn: c2}
	udp, err := listenUDP(outPrivKey, outConn, realaddr, netID)
	if err != nil {
		t.Fatal(err)
	}
	udp.net = &testNetWork{IP: node.IP}

	//test sendTopicNodes
	queryHash := common.Hash{0x01, 0x02, 0x03}
	var nodes []*Node
	var gotNodes []rpcNode
	for i := 0; i < 100; i++ {
		node := &Node{
			ID:  MustHexID("a2cb4c36765430f2e72564138c36f30fbc8af5a8bb91649822cd937dedbb8748"),
			IP:  net.ParseIP("1.2.3.4").To4(),
			UDP: uint16(i),
			TCP: uint16(i),
		}
		nodes = append(nodes, node)
	}
	udp.sendTopicNodes(node, queryHash, nodes)
	pkts := receivePacket(udpInput)
	for _, pkt := range pkts {
		gotNodes = append(gotNodes, pkt.data.(*topicNodes).Nodes[:]...)
	}
	for i := 0; i < len(gotNodes); i++ {
		if !reflect.DeepEqual(nodeToRPC(nodes[i]), gotNodes[i]) {
			t.Fatal("sendTopicNodes node field err")
		}
	}

	nodes = nodes[:0]
	gotNodes = gotNodes[:0]
	udp.sendTopicNodes(node, queryHash, nodes)
	pkts = receivePacket(udpInput)
	for _, pkt := range pkts {
		gotNodes = append(gotNodes, pkt.data.(*topicNodes).Nodes[:]...)
	}
	for i := 0; i < len(gotNodes); i++ {
		if !reflect.DeepEqual(nodeToRPC(nodes[i]), gotNodes[i]) {
			t.Fatal("sendTopicNodes node field err")
		}
	}
}

func receivePacket(udpInput *udp) []ingressPacket {
	waitTicker := time.NewTimer(10 * time.Millisecond)
	defer waitTicker.Stop()
	var msgs []ingressPacket
	for {
		select {
		case msg := <-udpInput.net.(*testNetWork).read:
			msgs = append(msgs, msg)
		case <-waitTicker.C:
			return msgs
		}
	}
	return msgs
}
