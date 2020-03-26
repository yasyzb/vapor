package config

import (
	"encoding/hex"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/bytom/vapor/crypto/ed25519/chainkd"
)

var (
	// CommonConfig means config object
	CommonConfig *Config
)

type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`
	// Options for services
	P2P        *P2PConfig        `mapstructure:"p2p"`
	Wallet     *WalletConfig     `mapstructure:"wallet"`
	Auth       *RPCAuthConfig    `mapstructure:"auth"`
	Web        *WebConfig        `mapstructure:"web"`
	Websocket  *WebsocketConfig  `mapstructure:"ws"`
	Federation *FederationConfig `mapstructure:"federation"`
	CrossChain *CrossChainConfig `mapstructure:"cross_chain"`
}

// Default configurable parameters.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		P2P:        DefaultP2PConfig(),
		Wallet:     DefaultWalletConfig(),
		Auth:       DefaultRPCAuthConfig(),
		Web:        DefaultWebConfig(),
		Websocket:  DefaultWebsocketConfig(),
		Federation: DefaultFederationConfig(),
		CrossChain: DefaultCrossChainConfig(),
	}
}

// Set the RootDir for all Config structs
func (cfg *Config) SetRoot(root string) *Config {
	cfg.BaseConfig.RootDir = root
	return cfg
}

// NodeKey retrieves the currently configured private key of the node, checking
// first any manually set key, falling back to the one found in the configured
// data folder. If no key can be found, a new one is generated.
func (cfg *Config) PrivateKey() *chainkd.XPrv {
	if cfg.XPrv != nil {
		return cfg.XPrv
	}

	filePath := rootify(cfg.PrivateKeyFile, cfg.BaseConfig.RootDir)
	fildReader, err := os.Open(filePath)
	if err != nil {
		log.WithField("err", err).Panic("fail on open private key file")
	}

	defer fildReader.Close()
	buf := make([]byte, 128)
	if _, err = io.ReadFull(fildReader, buf); err != nil {
		log.WithField("err", err).Panic("fail on read private key file")
	}

	var xprv chainkd.XPrv
	if _, err := hex.Decode(xprv[:], buf); err != nil {
		log.WithField("err", err).Panic("fail on decode private key")
	}

	cfg.XPrv = &xprv
	xpub := cfg.XPrv.XPub()
	cfg.XPub = &xpub
	return cfg.XPrv
}

//-----------------------------------------------------------------------------
// BaseConfig
type BaseConfig struct {
	// The root directory for all data.
	// This should be set in viper so it can unmarshal into this struct
	RootDir string `mapstructure:"home"`

	//The ID of the network to json
	ChainID string `mapstructure:"chain_id"`

	//log level to set
	LogLevel string `mapstructure:"log_level"`

	// A custom human readable name for this node
	Moniker string `mapstructure:"moniker"`

	// TCP or UNIX socket address for the profiling server to listen on
	ProfListenAddress string `mapstructure:"prof_laddr"`

	Mining bool `mapstructure:"mining"`

	// Database backend: leveldb | memdb
	DBBackend string `mapstructure:"db_backend"`

	// Database directory
	DBPath string `mapstructure:"db_dir"`

	// Keystore directory
	KeysPath string `mapstructure:"keys_dir"`

	ApiAddress string `mapstructure:"api_addr"`

	VaultMode bool `mapstructure:"vault_mode"`

	// log file name
	LogFile string `mapstructure:"log_file"`

	PrivateKeyFile string `mapstructure:"private_key_file"`
	XPrv           *chainkd.XPrv
	XPub           *chainkd.XPub

	// Federation file name
	FederationFileName string `mapstructure:"federation_file"`
}

// Default configurable base parameters.
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Moniker:            "anonymous",
		ProfListenAddress:  "",
		Mining:             false,
		DBBackend:          "leveldb",
		DBPath:             "data",
		KeysPath:           "keystore",
		LogFile:            "log",
		PrivateKeyFile:     "node_key.txt",
		FederationFileName: "federation.json",
	}
}

func (b BaseConfig) DBDir() string {
	return rootify(b.DBPath, b.RootDir)
}

func (b BaseConfig) LogDir() string {
	return rootify(b.LogFile, b.RootDir)
}

func (b BaseConfig) KeysDir() string {
	return rootify(b.KeysPath, b.RootDir)
}

func (b BaseConfig) FederationFile() string {
	return rootify(b.FederationFileName, b.RootDir)
}

// P2PConfig
type P2PConfig struct {
	ListenAddress    string `mapstructure:"laddr"`
	Seeds            string `mapstructure:"seeds"`
	SkipUPNP         bool   `mapstructure:"skip_upnp"`
	LANDiscover      bool   `mapstructure:"lan_discoverable"`
	MaxNumPeers      int    `mapstructure:"max_num_peers"`
	HandshakeTimeout int    `mapstructure:"handshake_timeout"`
	DialTimeout      int    `mapstructure:"dial_timeout"`
	ProxyAddress     string `mapstructure:"proxy_address"`
	ProxyUsername    string `mapstructure:"proxy_username"`
	ProxyPassword    string `mapstructure:"proxy_password"`
	KeepDial         string `mapstructure:"keep_dial"`
	Compression      string `mapstructure:"compression_backend"`
}

// Default configurable p2p parameters.
func DefaultP2PConfig() *P2PConfig {
	return &P2PConfig{
		ListenAddress:    "tcp://0.0.0.0:56656",
		SkipUPNP:         false,
		LANDiscover:      true,
		MaxNumPeers:      20,
		HandshakeTimeout: 30,
		DialTimeout:      3,
		ProxyAddress:     "",
		ProxyUsername:    "",
		ProxyPassword:    "",
		Compression:      "snappy",
	}
}

//-----------------------------------------------------------------------------
type WalletConfig struct {
	Disable  bool   `mapstructure:"disable"`
	Rescan   bool   `mapstructure:"rescan"`
	TxIndex  bool   `mapstructure:"txindex"`
	MaxTxFee uint64 `mapstructure:"max_tx_fee"`
}

type RPCAuthConfig struct {
	Disable bool `mapstructure:"disable"`
}

type WebConfig struct {
	Closed bool `mapstructure:"closed"`
}

type WebsocketConfig struct {
	MaxNumWebsockets     int `mapstructure:"max_num_websockets"`
	MaxNumConcurrentReqs int `mapstructure:"max_num_concurrent_reqs"`
}

type FederationConfig struct {
	Xpubs  []chainkd.XPub `json:"xpubs"`
	Quorum int            `json:"quorum"`
}

type CrossChainConfig struct {
	AssetWhitelist string `mapstructure:"asset_whitelist"`
}

// Default configurable rpc's auth parameters.
func DefaultRPCAuthConfig() *RPCAuthConfig {
	return &RPCAuthConfig{
		Disable: false,
	}
}

// Default configurable web parameters.
func DefaultWebConfig() *WebConfig {
	return &WebConfig{
		Closed: false,
	}
}

// Default configurable wallet parameters.
func DefaultWalletConfig() *WalletConfig {
	return &WalletConfig{
		Disable:  false,
		Rescan:   false,
		TxIndex:  false,
		MaxTxFee: uint64(1000000000),
	}
}

// Default configurable websocket parameters.
func DefaultWebsocketConfig() *WebsocketConfig {
	return &WebsocketConfig{
		MaxNumWebsockets:     25,
		MaxNumConcurrentReqs: 20,
	}
}

// Default configurable federation parameters.
func DefaultFederationConfig() *FederationConfig {
	return &FederationConfig{
		Xpubs: []chainkd.XPub{
			xpub("580daf48fa8962100047cb1391da890bb7f2c849fdbc9b368cb4394a4c7cbb0977e2e7ebbf055dc0ef90af6a0d2af01ce7ec56b735d016aab597815ec48552e5"),
			xpub("f3f6bcf61b65fa9d1566455a5688ca8b395efdc22e654963134b5e5cb0a45d8be522d21abc384a73177a7b9d64eba915fcfe2862d86a508a3c46dc410bdd72ad"),
			xpub("53559612f2b7bcada18948b7de39d63947a0e2bd7336d07db1350c54ba5743996b84bf9d18ff7a2457e1a5c70ce5013e4a3b62666ddb03294c53051d5f5c70c0"),
			xpub("7c88cc58adfc71818b08308d43c29de22460b0ea6895449cbec6e458d7dc09e0aea243fa5075ee6621da0d805bd047f6bb207329c5bd2ca3253b172fb323b512"),
		},
		Quorum: 2,
	}
}

// Default configurable crosschain parameters.
func DefaultCrossChainConfig() *CrossChainConfig {
	return &CrossChainConfig{}
}

func xpub(str string) (xpub chainkd.XPub) {
	if err := xpub.UnmarshalText([]byte(str)); err != nil {
		log.Panicf("Fail converts a string to xpub")
	}
	return xpub
}

//-----------------------------------------------------------------------------
// Utils

// helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home == "" {
		return "./.vapor"
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Vapor")
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "Vapor")
	default:
		return filepath.Join(home, ".vapor")
	}
}

func isFolderNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
