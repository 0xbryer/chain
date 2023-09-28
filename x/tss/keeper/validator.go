package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/bandprotocol/chain/v2/x/tss/types"
)

// valWithPower is an internal type to track validator with voting power inside of AllocateTokens.
type valWithPower struct {
	val   stakingtypes.ValidatorI
	power int64
}

// AllocateTokens allocates a portion of fee collected in the previous blocks to validators that
// that are actively performing tss tasks. Note that this reward is also subjected to comm tax
// and this reward is calculate after allocate to active tss validators
func (k Keeper) AllocateTokens(ctx sdk.Context, previousVotes []abci.VoteInfo) {
	toReward := []valWithPower{}
	totalPower := int64(0)
	for _, vote := range previousVotes {
		val := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		if k.GetStatus(ctx, vote.Validator.Address).Status == types.MEMBER_STATUS_ACTIVE {
			toReward = append(toReward, valWithPower{val: val, power: vote.Validator.Power})
			totalPower += vote.Validator.Power
		}
	}
	if totalPower == 0 {
		// No active validators performing tss tasks, nothing needs to be done here.
		return
	}

	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	totalFee := sdk.NewDecCoinsFromCoins(k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())...)
	// Compute the fee allocated for tss module to distribute to active validators.
	tssRewardRatio := sdk.NewDecWithPrec(int64(k.GetParams(ctx).RewardPercentage), 2)
	tssRewardInt, _ := totalFee.MulDecTruncate(tssRewardRatio).TruncateDecimal()
	// Transfer the tss reward portion from fee collector to distr module.
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distrtypes.ModuleName, tssRewardInt)
	if err != nil {
		panic(err)
	}
	// Convert the transferred tokens back to DecCoins for internal distr allocations.
	tssReward := sdk.NewDecCoinsFromCoins(tssRewardInt...)
	remaining := tssReward
	rewardMultiplier := sdk.OneDec().Sub(k.distrKeeper.GetCommunityTax(ctx))
	// Allocate non-community pool tokens to active validators weighted by voting power.
	for _, each := range toReward {
		powerFraction := sdk.NewDec(each.power).QuoTruncate(sdk.NewDec(totalPower))
		reward := tssReward.MulDecTruncate(rewardMultiplier).MulDecTruncate(powerFraction)
		k.distrKeeper.AllocateTokensToValidator(ctx, each.val, reward)
		remaining = remaining.Sub(reward)
	}
	// Allocate the remaining coins to the community pool.
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.distrKeeper.SetFeePool(ctx, feePool)
}
