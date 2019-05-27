package txbuilder

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	edchainkd "github.com/vapor/crypto/ed25519/chainkd"
	chainjson "github.com/vapor/encoding/json"
	"github.com/vapor/testutil"
)

func TestWitnessJSON(t *testing.T) {
	si := &SigningInstruction{
		Position: 17,
		WitnessComponents: []witnessComponent{
			DataWitness{1, 2, 3},
			&SignatureWitness{
				Quorum: 4,
				Keys: []keyID{{
					XPub:           testutil.TestXPub,
					DerivationPath: []chainjson.HexBytes{{5, 6, 7}},
				}},
				Sigs: []chainjson.HexBytes{{8, 9, 10}},
			},
			&RawTxSigWitness{
				Quorum: 20,
				Keys: []keyID{{
					XPub:           testutil.TestXPub,
					DerivationPath: []chainjson.HexBytes{{21, 22}},
				}},
				Sigs: []chainjson.HexBytes{{23, 24, 25}},
			},
		},
	}

	b, err := json.MarshalIndent(si, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var got SigningInstruction
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("error on input %s: %s", b, err)
	}
	for _, wc := range got.WitnessComponents {
		switch w := wc.(type) {
		case *SignatureWitness:
			for i, key := range w.Keys {
				if reflect.TypeOf(key.XPub).String() == "string" {
					if xpub, err := edchainkd.NewXPub(reflect.ValueOf(key.XPub).String()); err != nil {
						panic(err)
					} else {
						w.Keys[i].XPub = *xpub
					}
				}
			}
		case *RawTxSigWitness:
			for i, key := range w.Keys {
				if reflect.TypeOf(key.XPub).String() == "string" {
					if xpub, err := edchainkd.NewXPub(reflect.ValueOf(key.XPub).String()); err != nil {
						panic(err)
					} else {
						w.Keys[i].XPub = *xpub
					}
				}
			}
		}
	}

	if !testutil.DeepEqual(si, &got) {
		t.Errorf("got:\n%s\nwant:\n%s\nJSON was: %s", spew.Sdump(&got), spew.Sdump(si), string(b))
	}
}
