package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgRegisterChain = "register_chain"

var _ sdk.Msg = &MsgRegisterChain{}

func NewMsgRegisterChain(chainID string, chainInfo string) *MsgRegisterChain {
	return &MsgRegisterChain{
		ChainID:   chainID,
		ChainInfo: chainInfo,
	}
}

func (msg *MsgRegisterChain) Route() string {
	return RouterKey
}

func (msg *MsgRegisterChain) Type() string {
	return TypeMsgRegisterChain
}

func (msg *MsgRegisterChain) GetSigners() []sdk.AccAddress {
	return msg.GetSigners()
}

func (msg *MsgRegisterChain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterChain) ValidateBasic() error {
	// TODO ensure signer is authorized to send MsgRegisterChain and manipulate the KVStore
	return nil
}
