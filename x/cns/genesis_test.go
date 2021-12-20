package cns_test

import (
	"testing"

	keepertest "github.com/strangelove-ventures/cns/testutil/keeper"
	"github.com/strangelove-ventures/cns/testutil/nullify"
	"github.com/strangelove-ventures/cns/x/cns"
	"github.com/strangelove-ventures/cns/x/cns/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CnsKeeper(t)
	cns.InitGenesis(ctx, *k, genesisState)
	got := cns.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
