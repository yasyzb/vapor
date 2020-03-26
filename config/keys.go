package config

import (
	"errors"
	"path"

	cmn "github.com/tendermint/tmlibs/common"
	"github.com/spf13/viper"

	"github.com/vapor/common"
	"github.com/vapor/protocol/vm/vmutil"
	"github.com/vapor/consensus"
)

// public key => address and private key
var keyMap = map[string]*addrAndPriKey{}

type addrAndPriKey struct {
	address []byte
	priKey  string
}

// read from file
type keyConfig struct {
	Addresses []string `mapstructure:"addresses"`
	PubKeys   []string `mapstructure:"public_keys"`
	PriKyes   []string `mapstructure:"private_keys"`
}

func InitKeyConfig(root string) {
	keyPath := path.Join(root, "keys.yml")
	if err := cmn.EnsureDir(keyPath, 0700); err != nil {
		cmn.Exit("Error: " + err.Error())
	}

	v := viper.New()
	v.SetConfigFile(keyPath)
	if err := v.ReadInConfig(); err != nil {
		cmn.Exit("Error: " + err.Error())
	}
	var cfg keyConfig
	v.Unmarshal(&cfg)

	if len(cfg.Addresses) == 0 || len(cfg.Addresses) != len(cfg.PriKyes) || len(cfg.Addresses) != len(cfg.PubKeys) {
		cmn.Exit("invalid keys config")
	}

	for i := 0; i < len(cfg.Addresses); i++ {
		program, err := getProgramByAddress(cfg.Addresses[i])
		if err != nil {
			cmn.Exit("invalid address " + cfg.Addresses[i])
		}

		keyMap[cfg.PubKeys[i]] = &addrAndPriKey{
			address: program,
			priKey:  cfg.PriKyes[i],
		}
	}
}

// GetAddrAndPriKey gets address and private key according to public key
func GetAddrAndPriKey(pubKey string) ([]byte, string) {
	key, ok := keyMap[pubKey]
	if !ok {
		return nil, ""
	}
	return key.address, key.priKey
}

func getProgramByAddress(address string) ([]byte, error) {
	addr, err := common.DecodeAddress(address, &consensus.ActiveNetParams)
	if err != nil {
		return nil, err
	}

	redeemContract := addr.ScriptAddress()
	program := []byte{}
	switch addr.(type) {
	case *common.AddressWitnessPubKeyHash:
		program, err = vmutil.P2WPKHProgram(redeemContract)
	case *common.AddressWitnessScriptHash:
		program, err = vmutil.P2WSHProgram(redeemContract)
	default:
		return nil, errors.New("invalid address")
	}
	if err != nil {
		return nil, err
	}

	return program, nil
}
