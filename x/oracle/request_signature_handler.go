package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bandprotocol/chain/v2/x/oracle/keeper"
	"github.com/bandprotocol/chain/v2/x/oracle/types"
	tsstypes "github.com/bandprotocol/chain/v2/x/tss/types"
)

// NewRequestingSignatureHandler creates a new TSS Handler for requesting the signature
func NewRequestingSignatureHandler(k keeper.Keeper) tsstypes.Handler {
	return func(ctx sdk.Context, content tsstypes.Content) ([]byte, error) {
		switch c := content.(type) {
		case *types.OracleResultRequestingSignature:
			return k.GetByteResult(ctx, c.RequestID)

		default:
			return nil, sdkerrors.Wrapf(
				sdkerrors.ErrUnknownRequest,
				"unrecognized tss request signature type: %s",
				c.RequestingSignatureType(),
			)
		}
	}
}
