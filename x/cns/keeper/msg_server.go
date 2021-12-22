package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/cns/x/cns/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterChain(goCtx context.Context, msg *types.MsgRegisterChain) (*types.MsgRegisterChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO check that the signer of the msg is authorized to write to the KVStore and handle accordingly
	_ = ctx

	return &types.MsgRegisterChainResponse{}, nil
}
