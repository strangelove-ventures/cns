package keeper

import (
	"github.com/strangelove-ventures/cns/x/cns/types"
)

var _ types.QueryServer = Keeper{}
