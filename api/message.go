package api

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/vapor/blockchain/signers"
	"github.com/vapor/common"
	"github.com/vapor/consensus"
	"github.com/vapor/crypto"
	vcrypto "github.com/vapor/crypto"
	"github.com/vapor/crypto/csp"
	"github.com/vapor/crypto/ed25519"
	chainjson "github.com/vapor/encoding/json"
)

// SignMsgResp is response for sign message
type SignMsgResp struct {
	Signature   string            `json:"signature"`
	DerivedXPub vcrypto.XPubKeyer `json:"derived_xpub"`
}

func (a *API) signMessage(ctx context.Context, ins struct {
	Address  string             `json:"address"`
	Message  chainjson.HexBytes `json:"message"`
	Password string             `json:"password"`
}) Response {
	cp, err := a.wallet.AccountMgr.GetLocalCtrlProgramByAddress(ins.Address)
	if err != nil {
		return NewErrorResponse(err)
	}

	account, err := a.wallet.AccountMgr.GetAccountByProgram(cp)
	if err != nil {
		return NewErrorResponse(err)
	}

	path, err := signers.Path(account.Signer, signers.AccountKeySpace, cp.Change, cp.KeyIndex)
	if err != nil {
		return NewErrorResponse(err)
	}
	derivedXPubs := csp.DeriveXPubs(account.XPubs, path)

	sig, err := a.wallet.Hsm.XSign(account.XPubs[0], path, ins.Message, ins.Password)
	if err != nil {
		return NewErrorResponse(err)
	}
	return NewSuccessResponse(SignMsgResp{
		Signature:   hex.EncodeToString(sig),
		DerivedXPub: derivedXPubs[0],
	})
}

// VerifyMsgResp is response for verify message
type VerifyMsgResp struct {
	VerifyResult bool `json:"result"`
}

func (a *API) verifyMessage(ctx context.Context, ins struct {
	Address     string             `json:"address"`
	DerivedXPub vcrypto.XPubKeyer  `json:"derived_xpub"`
	Message     chainjson.HexBytes `json:"message"`
	Signature   string             `json:"signature"`
}) Response {
	sig, err := hex.DecodeString(ins.Signature)
	if err != nil {
		return NewErrorResponse(err)
	}

	derivedPK := ins.DerivedXPub.PublicKey()
	pubHash := crypto.Ripemd160(derivedPK)
	addressPubHash, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.ActiveNetParams)
	if err != nil {
		return NewErrorResponse(err)
	}

	address := addressPubHash.EncodeAddress()
	if address != strings.TrimSpace(ins.Address) {
		return NewSuccessResponse(VerifyMsgResp{VerifyResult: false})
	}

	if ed25519.Verify(ins.DerivedXPub.PublicKey(), ins.Message, sig) {
		return NewSuccessResponse(VerifyMsgResp{VerifyResult: true})
	}
	return NewSuccessResponse(VerifyMsgResp{VerifyResult: false})
}
