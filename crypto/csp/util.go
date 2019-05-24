// csp is a package of cipher service provider
package csp

import (
	"crypto"
	"fmt"
	"io"
	"reflect"

	vcrypto "github.com/vapor/crypto"
	edchainkd "github.com/vapor/crypto/ed25519/chainkd"
)

// Utility functions

func NewXKeys(r io.Reader) (xprv vcrypto.XPrvKeyer, xpub vcrypto.XPubKeyer, err error) {
	// TODO: if ... create sm2 xprv and xpub
	// return .....

	// if ... create ed25519 xprv and xpub
	return edchainkd.NewXKeys(r)
}

func XPubKeys(xpubs []vcrypto.XPubKeyer) []crypto.PublicKey {
	res := make([]crypto.PublicKey, 0, len(xpubs))
	for _, xpub := range xpubs {
		switch xpb := xpub.(type) {
		case edchainkd.XPub:
			res = append(res, xpb.PublicKey())
		}
	}
	return res
}

// func DeriveXPubs(xpubs []vcrypto.XPubKeyer, path [][]byte) []vcrypto.XPubKeyer {
// 	res := make([]vcrypto.XPubKeyer, 0, len(xpubs))
// 	for _, xpub := range xpubs {
// 		switch xpb := xpub.(type) {
// 		case edchainkd.XPub:
// 			d := xpb.Derive(path)
// 			res = append(res, d)
// 			fmt.Println("DeriveXPubs d is:", d)
// 		}
// 	}
// 	fmt.Println("DeriveXPubs len(res) is:", len(res))
// 	return res
// }

func DeriveXPubs(xpubs []vcrypto.XPubKeyer, path [][]byte) []vcrypto.XPubKeyer {
	res := make([]vcrypto.XPubKeyer, 0, len(xpubs))
	for _, xpub := range xpubs {
		fmt.Println("DeriveXPubs type xpubs:", reflect.TypeOf(xpubs), "type xpubs[0]:", reflect.TypeOf(xpubs[0]))
		fmt.Println("xpubs[0] is:", xpubs[0])
		newxpub, err := edchainkd.NewXPub(reflect.ValueOf(xpub).String())
		if err != nil {
			fmt.Println("csp DeriveXPubs err:", err)
		} else {
			d := newxpub.Derive(path)
			res = append(res, d)
			fmt.Println("DeriveXPubs d is:", d)
		}

		// switch xpb := xpub.(type) {
		// case edchainkd.XPub:
		// 	d := xpb.Derive(path)
		// 	res = append(res, d)
		// 	fmt.Println("DeriveXPubs d is:", d)
		// }
	}
	fmt.Println("DeriveXPubs len(res) is:", len(res))
	return res
}
