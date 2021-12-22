package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/strangelove-ventures/cns/x/cns/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) RegisteredChain(goCtx context.Context, req *types.QueryRegisteredChainRequest) (*types.QueryRegisteredChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO check that the requested chainID exists in the CNS registrar and respond accordingly
	_ = ctx

	return &types.QueryRegisteredChainResponse{}, nil
}

func (k Keeper) RegisteredChainsAll(c context.Context, req *types.QueryAllRegisteredChainsRequest) (*types.QueryAllRegisteredChainsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var chains = make(map[string]string)
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	registeredChainsStore := prefix.NewStore(store, types.KeyPrefix(types.RegisteredChainsKeyPrefix))

	// TODO this may be wrong? tried to work backwards from another example but with a map perhaps this isn't what we want
	pageRes, err := query.Paginate(registeredChainsStore, req.Pagination, func(key []byte, value []byte) error {
		var registeredChains types.RegisteredChains
		if err := k.cdc.Unmarshal(value, &registeredChains); err != nil {
			return err
		}

		for key, val := range registeredChains.Chains {
			chains[key] = val
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRegisteredChainsResponse{RegisteredChains: &types.RegisteredChains{Chains: chains}, Pagination: pageRes}, nil
}
