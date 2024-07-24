package keeper_test

import "github.com/bandprotocol/chain/v2/x/feeds/types"

func (suite *KeeperTestSuite) TestMsgSubmitSignals() {
	testCases := []struct {
		name      string
		input     *types.MsgSubmitSignals
		expErr    bool
		expErrMsg string
		postCheck func()
	}{
		{
			name: "no delegation",
			input: &types.MsgSubmitSignals{
				Delegator: InvalidDelegator.String(),
				Signals: []types.Signal{
					{
						ID:    "CS:BAND-USD",
						Power: 10,
					},
				},
			},
			expErr:    true,
			expErrMsg: "delegation not enough",
			postCheck: func() {},
		},
		{
			name: "1 signal more than delegations",
			input: &types.MsgSubmitSignals{
				Delegator: ValidDelegator.String(),
				Signals: []types.Signal{
					{
						ID:    "CS:BAND-USD",
						Power: 1e10 + 1,
					},
				},
			},
			expErr:    true,
			expErrMsg: "delegation not enough",
			postCheck: func() {},
		},
		{
			name: "2 signals more than delegations",
			input: &types.MsgSubmitSignals{
				Delegator: ValidDelegator.String(),
				Signals: []types.Signal{
					{
						ID:    "CS:BAND-USD",
						Power: 1e10,
					},
					{
						ID:    "CS:ATOM-USD",
						Power: 1,
					},
				},
			},
			expErr:    true,
			expErrMsg: "delegation not enough",
			postCheck: func() {},
		},
		{
			name: "valid request",
			input: &types.MsgSubmitSignals{
				Delegator: ValidDelegator.String(),
				Signals: []types.Signal{
					{
						ID:    "CS:BAND-USD",
						Power: 1e10,
					},
				},
			},
			expErr:    false,
			expErrMsg: "",
			postCheck: func() {
				suite.Require().Equal(
					[]types.Signal{
						{
							ID:    "CS:BAND-USD",
							Power: 1e10,
						},
					},
					suite.feedsKeeper.GetDelegatorSignals(suite.ctx, ValidDelegator),
				)
				suite.Require().Equal(
					[]types.Signal{
						{
							ID:    "CS:BAND-USD",
							Power: 1e10,
						},
					},
					suite.feedsKeeper.GetSignalTotalPowersByPower(suite.ctx, 300),
				)
			},
		},
		{
			name: "valid request (replace)",
			input: &types.MsgSubmitSignals{
				Delegator: ValidDelegator.String(),
				Signals: []types.Signal{
					{
						ID:    "CS:BAND-USD",
						Power: 1e9,
					},
					{
						ID:    "CS:BTC-USD",
						Power: 1e9,
					},
				},
			},
			expErr:    false,
			expErrMsg: "",
			postCheck: func() {
				suite.Require().Equal(
					[]types.Signal{
						{
							ID:    "CS:BAND-USD",
							Power: 1e9,
						},
						{
							ID:    "CS:BTC-USD",
							Power: 1e9,
						},
					},
					suite.feedsKeeper.GetDelegatorSignals(suite.ctx, ValidDelegator),
				)
				suite.Require().Equal(
					[]types.Signal{
						{
							ID:    "CS:BAND-USD",
							Power: 1e9,
						},
						{
							ID:    "CS:BTC-USD",
							Power: 1e9,
						},
					},
					suite.feedsKeeper.GetSignalTotalPowersByPower(suite.ctx, 300),
				)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.SubmitSignals(suite.ctx, tc.input)

			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}

			tc.postCheck()
		})
	}
}

func (suite *KeeperTestSuite) TestMsgSubmitSignalPrices() {
	suite.feedsKeeper.SetCurrentFeeds(suite.ctx, []types.Feed{{
		SignalID: "CS:BAND-USD",
		Interval: 100,
	}})

	testCases := []struct {
		name      string
		input     *types.MsgSubmitSignalPrices
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid validator",
			input: &types.MsgSubmitSignalPrices{
				Validator: InvalidValidator.String(),
				Timestamp: suite.ctx.BlockTime().Unix(),
				Prices: []types.SignalPrice{
					{
						PriceStatus: types.PriceStatusAvailable,
						SignalID:    "CS:BAND-USD",
						Price:       10e12,
					},
				},
			},
			expErr:    true,
			expErrMsg: "not bonded validator",
		},
		{
			name: "invalid symbol",
			input: &types.MsgSubmitSignalPrices{
				Validator: ValidValidator.String(),
				Timestamp: suite.ctx.BlockTime().Unix(),
				Prices: []types.SignalPrice{
					{
						PriceStatus: types.PriceStatusAvailable,
						SignalID:    "CS:BTC-USD",
						Price:       10e12,
					},
				},
			},
			expErr:    true,
			expErrMsg: "signal id is not supported",
		},
		{
			name: "invalid timestamp",
			input: &types.MsgSubmitSignalPrices{
				Validator: ValidValidator.String(),
				Timestamp: suite.ctx.BlockTime().Unix() - 200,
				Prices: []types.SignalPrice{
					{
						PriceStatus: types.PriceStatusAvailable,
						SignalID:    "CS:BAND-USD",
						Price:       10e12,
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid timestamp",
		},
		{
			name: "valid message",
			input: &types.MsgSubmitSignalPrices{
				Validator: ValidValidator.String(),
				Timestamp: suite.ctx.BlockTime().Unix(),
				Prices: []types.SignalPrice{
					{
						PriceStatus: types.PriceStatusAvailable,
						SignalID:    "CS:BAND-USD",
						Price:       10e12,
					},
				},
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name: "price too fast",
			input: &types.MsgSubmitSignalPrices{
				Validator: ValidValidator.String(),
				Timestamp: suite.ctx.BlockTime().Unix(),
				Prices: []types.SignalPrice{
					{
						PriceStatus: types.PriceStatusAvailable,
						SignalID:    "CS:BAND-USD",
						Price:       10e12,
					},
				},
			},
			expErr:    true,
			expErrMsg: "price is submitted too early",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.SubmitSignalPrices(suite.ctx, tc.input)

			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgUpdateReferenceSourceConfig() {
	params := suite.feedsKeeper.GetParams(suite.ctx)
	referenceSourceConfig := types.DefaultReferenceSourceConfig()

	testCases := []struct {
		name      string
		input     *types.MsgUpdateReferenceSourceConfig
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid admin",
			input: &types.MsgUpdateReferenceSourceConfig{
				Admin:                 "invalid",
				ReferenceSourceConfig: referenceSourceConfig,
			},
			expErr:    true,
			expErrMsg: "invalid admin",
		},
		{
			name: "all good",
			input: &types.MsgUpdateReferenceSourceConfig{
				Admin:                 params.Admin,
				ReferenceSourceConfig: referenceSourceConfig,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.UpdateReferenceSourceConfig(suite.ctx, tc.input)

			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgUpdateParams() {
	params := types.DefaultParams()

	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "all good",
			input: &types.MsgUpdateParams{
				Authority: suite.feedsKeeper.GetAuthority(),
				Params:    params,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.UpdateParams(suite.ctx, tc.input)

			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
