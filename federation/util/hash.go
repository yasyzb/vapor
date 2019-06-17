package util

import (
	// "encoding/hex"

	// "github.com/vapor/errors"
	"github.com/vapor/protocol/bc"
)

func StringToAssetID(s string) (*bc.Hash, error) {
	// h, err := hex.DecodeString(s)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "decode asset string")
	// }

	// panic(len(h))

	// var b [32]byte
	// copy(b[:], h)
	// assetID := bc.NewAssetID(b)
	assetID := &bc.Hash{}
	if err := assetID.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}

	return assetID, nil
}
