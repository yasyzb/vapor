package config

// public key => address and private key
var keyMap = map[string]*addrAndPriKey{
	"public_key_1": {
		"address_1",
		"private_key_1",
	},
	"public_key_2": {
		"address_2",
		"private_key_2",
	},
}

type addrAndPriKey struct {
	Address string
	PriKey string
}

// GetAddrAndPriKey gets address and private key according to public key
func GetAddrAndPriKey(pubKey string) *addrAndPriKey {
	return keyMap[pubKey]
}
