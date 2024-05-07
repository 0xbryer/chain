package keeper_test

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/bandprotocol/chain/v2/pkg/tss/testutil"
	bandtesting "github.com/bandprotocol/chain/v2/testing"
	"github.com/bandprotocol/chain/v2/x/bandtss/keeper"
	"github.com/bandprotocol/chain/v2/x/bandtss/types"
	tsstypes "github.com/bandprotocol/chain/v2/x/tss/types"
)

var Coins1000000uband = sdk.NewCoins(sdk.NewInt64Coin("uband", 1000000))

func defaultVotes() []abci.VoteInfo {
	return []abci.VoteInfo{{
		Validator: abci.Validator{
			Address: bandtesting.Validators[0].PubKey.Address(),
			Power:   70,
		},
		SignedLastBlock: true,
	}, {
		Validator: abci.Validator{
			Address: bandtesting.Validators[1].PubKey.Address(),
			Power:   20,
		},
		SignedLastBlock: true,
	}, {
		Validator: abci.Validator{
			Address: bandtesting.Validators[2].PubKey.Address(),
			Power:   10,
		},
		SignedLastBlock: true,
	}}
}

func SetupFeeCollector(app *bandtesting.TestingApp, ctx sdk.Context, k keeper.Keeper) (authtypes.ModuleAccountI, error) {
	// Set collected fee to 1000000uband and 50% tss reward proportion.
	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	if err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, Coins1000000uband); err != nil {
		return nil, err
	}

	if err := app.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		Coins1000000uband,
	); err != nil {
		return nil, err
	}
	app.AccountKeeper.SetAccount(ctx, feeCollector)

	params := k.GetParams(ctx)
	params.RewardPercentage = 50
	if err := k.SetParams(ctx, params); err != nil {
		return nil, err
	}

	return feeCollector, nil
}

func (s *KeeperTestSuite) TestAllocateTokenNoActiveValidators() {
	app, ctx := bandtesting.CreateTestApp(s.T(), false)
	feeCollector, err := SetupFeeCollector(app, ctx, *app.BandtssKeeper)
	s.Require().NoError(err)

	s.Require().Equal(Coins1000000uband, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	// No active tss validators so nothing should happen.
	app.OracleKeeper.AllocateTokens(ctx, defaultVotes())

	distAccount := app.AccountKeeper.GetModuleAccount(ctx, disttypes.ModuleName)
	s.Require().Equal(Coins1000000uband, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	s.Require().Empty(app.BankKeeper.GetAllBalances(ctx, distAccount.GetAddress()))
}

func (s *KeeperTestSuite) TestAllocateTokensOneActive() {
	app, ctx := bandtesting.CreateTestApp(s.T(), false)
	tssKeeper, k := app.TSSKeeper, app.BandtssKeeper
	feeCollector, err := SetupFeeCollector(app, ctx, *k)
	s.Require().NoError(err)

	s.Require().Equal(Coins1000000uband, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	// From 50% of fee, 1% should go to community pool, the rest goes to the only active validator.
	err = tssKeeper.HandleSetDEs(ctx, bandtesting.Validators[1].Address, []tsstypes.DE{
		{
			PubD: testutil.HexDecode("dddd"),
			PubE: testutil.HexDecode("eeee"),
		},
	})
	s.Require().NoError(err)

	for _, validator := range bandtesting.Validators {
		err := k.AddNewMember(ctx, validator.Address)
		s.Require().NoError(err)
	}

	k.AllocateTokens(ctx, defaultVotes())

	distAccount := app.AccountKeeper.GetModuleAccount(ctx, disttypes.ModuleName)
	s.Require().Equal(
		sdk.NewCoins(sdk.NewInt64Coin("uband", 500000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	s.Require().Equal(
		sdk.NewCoins(sdk.NewInt64Coin("uband", 500000)),
		app.BankKeeper.GetAllBalances(ctx, distAccount.GetAddress()),
	)
	s.Require().Equal(
		sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(10000)}},
		app.DistrKeeper.GetFeePool(ctx).CommunityPool,
	)
	s.Require().Empty(app.DistrKeeper.GetValidatorOutstandingRewards(ctx, bandtesting.Validators[0].ValAddress))
	s.Require().Equal(
		sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(490000)}},
		app.DistrKeeper.GetValidatorOutstandingRewards(ctx, bandtesting.Validators[1].ValAddress).Rewards,
	)
	s.Require().Empty(app.DistrKeeper.GetValidatorOutstandingRewards(ctx, bandtesting.Validators[2].ValAddress))
}

func (s *KeeperTestSuite) TestAllocateTokensAllActive() {
	ctx, app, k := s.ctx, s.app, s.app.BandtssKeeper

	feeCollector, err := SetupFeeCollector(app, ctx, *k)
	s.Require().NoError(err)
	s.Require().Equal(Coins1000000uband, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))

	for _, validator := range bandtesting.Validators {
		err := k.AddNewMember(ctx, validator.Address)
		s.Require().NoError(err)
		deCount := s.app.TSSKeeper.GetDECount(ctx, validator.Address)
		s.Require().Greater(deCount, uint64(0))
	}

	// From 50% of fee, 1% should go to community pool, the rest get split to validators.
	k.AllocateTokens(ctx, defaultVotes())

	distAccount := app.AccountKeeper.GetModuleAccount(ctx, disttypes.ModuleName)
	s.Require().Equal(
		sdk.NewCoins(sdk.NewInt64Coin("uband", 500000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	s.Require().Equal(
		sdk.NewCoins(sdk.NewInt64Coin("uband", 500000)),
		app.BankKeeper.GetAllBalances(ctx, distAccount.GetAddress()),
	)
	s.Require().Equal(
		sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(10000)}},
		app.DistrKeeper.GetFeePool(ctx).CommunityPool,
	)
	s.Require().Equal(
		sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(343000)}},
		app.DistrKeeper.GetValidatorOutstandingRewards(ctx, bandtesting.Validators[0].ValAddress).Rewards,
	)
	s.Require().Equal(
		sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(98000)}},
		app.DistrKeeper.GetValidatorOutstandingRewards(ctx, bandtesting.Validators[1].ValAddress).Rewards,
	)
	s.Require().Equal(
		sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(49000)}},
		app.DistrKeeper.GetValidatorOutstandingRewards(ctx, bandtesting.Validators[2].ValAddress).Rewards,
	)
}

func (s *KeeperTestSuite) TestHandleInactiveValidators() {
	ctx, k := s.ctx, s.app.BandtssKeeper
	s.SetupGroup(tsstypes.GROUP_STATUS_ACTIVE)
	address := bandtesting.Validators[0].Address

	member := types.Member{
		Address:    address.String(),
		IsActive:   true,
		Since:      time.Time{},
		LastActive: time.Time{},
	}
	k.SetMember(ctx, member)
	s.app.TSSKeeper.SetMember(ctx, tsstypes.Member{
		ID:       1,
		GroupID:  1,
		Address:  address.String(),
		IsActive: true,
	})
	ctx = ctx.WithBlockTime(time.Now())

	k.HandleInactiveValidators(ctx)

	member, err := k.GetMember(ctx, address)
	s.Require().NoError(err)
	s.Require().False(member.IsActive)
}
