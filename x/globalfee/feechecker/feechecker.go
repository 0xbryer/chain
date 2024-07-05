package feechecker

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	feedskeeper "github.com/bandprotocol/chain/v2/x/feeds/keeper"
	feedstypes "github.com/bandprotocol/chain/v2/x/feeds/types"
	"github.com/bandprotocol/chain/v2/x/globalfee/keeper"
	oraclekeeper "github.com/bandprotocol/chain/v2/x/oracle/keeper"
	oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"
)

type FeeChecker struct {
	AuthzKeeper     *authzkeeper.Keeper
	OracleKeeper    *oraclekeeper.Keeper
	GlobalfeeKeeper *keeper.Keeper
	StakingKeeper   *stakingkeeper.Keeper
	FeedsKeeper     *feedskeeper.Keeper

	FeedsMsgServer feedstypes.MsgServer
}

func NewFeeChecker(
	authzKeeper *authzkeeper.Keeper,
	oracleKeeper *oraclekeeper.Keeper,
	globalfeeKeeper *keeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	feedsKeeper *feedskeeper.Keeper,
) FeeChecker {
	feedsMsgServer := feedskeeper.NewMsgServerImpl(*feedsKeeper)

	return FeeChecker{
		AuthzKeeper:     authzKeeper,
		OracleKeeper:    oracleKeeper,
		GlobalfeeKeeper: globalfeeKeeper,
		StakingKeeper:   stakingKeeper,
		FeedsKeeper:     feedsKeeper,
		FeedsMsgServer:  feedsMsgServer,
	}
}

// CheckTxFee is responsible for verifying whether a transaction contains the necessary fee.
func (fc FeeChecker) CheckTxFee(
	ctx sdk.Context,
	tx sdk.Tx,
) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, sdkerrors.ErrTxDecode.Wrap("Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()
	priority := getTxPriority(feeCoins, int64(gas), fc.GetBondDenom(ctx))

	// Ensure that the provided fees meet minimum-gas-prices and globalFees,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if !ctx.IsCheckTx() {
		return feeCoins, priority, nil
	}

	// Check if this tx should be free or not
	if fc.IsBypassMinFeeTx(ctx, tx) {
		return sdk.Coins{}, int64(math.MaxInt64), nil
	}

	minGasPrices := getMinGasPrices(ctx)
	globalMinGasPrices, err := fc.GetGlobalMinGasPrices(ctx)
	if err != nil {
		return nil, 0, err
	}

	allGasPrices := CombinedGasPricesRequirement(minGasPrices, globalMinGasPrices)

	// Calculate all fees from all gas prices
	var allFees sdk.Coins
	if !allGasPrices.IsZero() {
		glDec := sdk.NewDec(int64(gas))
		for _, gp := range allGasPrices {
			if !gp.IsZero() {
				fee := gp.Amount.Mul(glDec)
				allFees = append(allFees, sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt()))
			}
		}
	}

	if !allFees.IsZero() && !feeCoins.IsAnyGTE(allFees) {
		return nil, 0, sdkerrors.ErrInsufficientFee.Wrapf(
			"insufficient fees; got: %s required: %s",
			feeCoins,
			allFees,
		)
	}

	return feeCoins, priority, nil
}

// IsBypassMinFeeTx checks whether tx is min fee bypassable.
func (fc FeeChecker) IsBypassMinFeeTx(ctx sdk.Context, tx sdk.Tx) bool {
	newCtx, _ := ctx.CacheContext()

	// Check if all messages are free
	for _, msg := range tx.GetMsgs() {
		if !fc.IsBypassMinFeeMsg(newCtx, msg) {
			return false
		}
	}

	return true
}

// IsBypassMinFeeMsg checks whether msg is min fee bypassable.
func (fc FeeChecker) IsBypassMinFeeMsg(ctx sdk.Context, msg sdk.Msg) bool {
	switch msg := msg.(type) {
	case *oracletypes.MsgReportData:
		if err := checkValidMsgReport(ctx, fc.OracleKeeper, msg); err != nil {
			return false
		}
	case *feedstypes.MsgSubmitSignalPrices:
		if _, err := fc.FeedsMsgServer.SubmitSignalPrices(ctx, msg); err != nil {
			return false
		}
	case *authz.MsgExec:
		msgs, err := msg.GetMessages()
		if err != nil {
			return false
		}

		grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
		if err != nil {
			return false
		}

		for _, m := range msgs {
			// Check if this grantee have authorization for the message.
			cap, _ := fc.AuthzKeeper.GetAuthorization(
				ctx,
				grantee,
				m.GetSigners()[0],
				sdk.MsgTypeURL(m),
			)
			if cap == nil {
				return false
			}

			// Check if this message should be free or not.
			if !fc.IsBypassMinFeeMsg(ctx, m) {
				return false
			}
		}
	default:
		return false
	}

	return true
}

// GetGlobalMinGasPrices returns global min gas prices
func (fc FeeChecker) GetGlobalMinGasPrices(ctx sdk.Context) (sdk.DecCoins, error) {
	globalMinGasPrices := fc.GlobalfeeKeeper.GetParams(ctx).MinimumGasPrices
	if len(globalMinGasPrices) != 0 {
		return globalMinGasPrices.Sort(), nil
	}
	// global fee is empty set, set global fee to 0uband (bondDenom)
	globalMinGasPrices, err := fc.DefaultZeroGlobalFee(ctx)
	if err != nil {
		return globalMinGasPrices, err
	}

	return globalMinGasPrices.Sort(), nil
}

// DefaultZeroGlobalFee returns a zero coin with the staking module bond denom
func (fc FeeChecker) DefaultZeroGlobalFee(ctx sdk.Context) ([]sdk.DecCoin, error) {
	bondDenom := fc.GetBondDenom(ctx)
	if bondDenom == "" {
		return nil, sdkerrors.ErrNotFound.Wrap("empty staking bond denomination")
	}

	return []sdk.DecCoin{sdk.NewDecCoinFromDec(bondDenom, sdk.NewDec(0))}, nil
}

// GetBondDenom returns Bondable coin denomination
func (fc FeeChecker) GetBondDenom(ctx sdk.Context) string {
	return fc.StakingKeeper.BondDenom(ctx)
}
