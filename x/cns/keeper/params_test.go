package keeper_test

import (
	"testing"

	testkeeper "github.com/strangelove-ventures/cns/testutil/keeper"
	"github.com/strangelove-ventures/cns/x/cns/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CnsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
