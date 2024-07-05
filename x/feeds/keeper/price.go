package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/chain/v2/x/feeds/types"
)

// ValidateSubmitSignalPricesRequest validates price submission.
func (k Keeper) ValidateSubmitSignalPricesRequest(
	ctx sdk.Context,
	blockTime int64,
	req *types.MsgSubmitSignalPrices,
	val sdk.ValAddress,
) error {
	isValid := k.IsBondedValidator(ctx, val)
	if !isValid {
		return types.ErrNotBondedValidator
	}

	status := k.oracleKeeper.GetValidatorStatus(ctx, val)
	if !status.IsActive {
		return types.ErrOracleStatusNotActive.Wrapf("val: %s", val.String())
	}

	if types.AbsInt64(req.Timestamp-blockTime) > k.GetParams(ctx).AllowableBlockTimeDiscrepancy {
		return types.ErrInvalidTimestamp.Wrapf(
			"block_time: %d, timestamp: %d",
			blockTime,
			req.Timestamp,
		)
	}

	return nil
}

// NewValidatorPrice creates new ValidatorPrice.
func (k Keeper) NewValidatorPrice(
	val sdk.ValAddress,
	price types.SignalPrice,
	blockTime int64,
	blockHeight int64,
) types.ValidatorPrice {
	return types.ValidatorPrice{
		PriceStatus: price.PriceStatus,
		Validator:   val.String(),
		SignalID:    price.SignalID,
		Price:       price.Price,
		Timestamp:   blockTime,
		BlockHeight: blockHeight,
	}
}
