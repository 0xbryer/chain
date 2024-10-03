package keeper_test

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/bandprotocol/chain/v2/x/feeds/types"
)

func (suite *KeeperTestSuite) TestQueryAllCurrentPrices() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup
	prices := []types.Price{
		{
			SignalID:  "CS:ATOM-USD",
			Price:     100000000,
			Timestamp: 1234567890,
		},
		{
			SignalID:  "CS:BAND-USD",
			Price:     200000000,
			Timestamp: 1234567890,
		},
	}

	for _, price := range prices {
		suite.feedsKeeper.SetPrice(ctx, price)
	}

	// query and check
	res, err := queryClient.AllCurrentPrices(context.Background(), &types.QueryAllCurrentPricesRequest{})
	suite.Require().NoError(err)
	// signal ids are not in the current feeds
	suite.Require().Equal(&types.QueryAllCurrentPricesResponse{
		Prices: []types.Price(nil),
	}, res)

	// set current feeds
	feeds := []types.Feed{
		{
			SignalID: "CS:ATOM-USD",
			Interval: 100,
		},
		{
			SignalID: "CS:BAND-USD",
			Interval: 100,
		},
	}

	suite.feedsKeeper.SetCurrentFeeds(ctx, feeds)

	// query and check
	res, err = queryClient.AllCurrentPrices(context.Background(), &types.QueryAllCurrentPricesRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllCurrentPricesResponse{
		Prices: prices,
	}, res)
}

func (suite *KeeperTestSuite) TestQueryCurrentPrices() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup
	prices := []types.Price{
		{
			SignalID:  "CS:ATOM-USD",
			Price:     100000000,
			Timestamp: 1234567890,
		},
		{
			SignalID:  "CS:BAND-USD",
			Price:     200000000,
			Timestamp: 1234567890,
		},
	}

	for _, price := range prices {
		suite.feedsKeeper.SetPrice(ctx, price)
	}

	// set current feeds with only BAND
	feeds := []types.Feed{
		{
			SignalID: "CS:BAND-USD",
			Interval: 100,
		},
	}

	suite.feedsKeeper.SetCurrentFeeds(ctx, feeds)

	// query and check
	// ATOM is not in the current feeds so it should return unavailable
	expectedCurrentPrices := []types.Price{
		types.NewPrice(types.PriceStatusUnavailable, "CS:ATOM-USD", 0, 0),
		prices[1],
	}
	res, err := queryClient.CurrentPrices(context.Background(), &types.QueryCurrentPricesRequest{
		SignalIds: []string{"CS:ATOM-USD", "CS:BAND-USD"},
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryCurrentPricesResponse{
		Prices: expectedCurrentPrices,
	}, res)
}

func (suite *KeeperTestSuite) TestQueryDelegatorSignals() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup
	signals := []types.Signal{
		{
			ID:    "CS:BAND-USD",
			Power: 1e9,
		},
		{
			ID:    "CS:BTC-USD",
			Power: 1e9,
		},
	}
	_, err := suite.msgServer.SubmitSignals(ctx, &types.MsgSubmitSignals{
		Delegator: ValidDelegator.String(),
		Signals:   signals,
	})
	suite.Require().NoError(err)

	// query and check
	res, err := queryClient.DelegatorSignals(context.Background(), &types.QueryDelegatorSignalsRequest{
		DelegatorAddress: ValidDelegator.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryDelegatorSignalsResponse{
		Signals: signals,
	}, res)
}

func (suite *KeeperTestSuite) TestQueryPrices() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup
	prices := []*types.Price{
		{
			SignalID:  "CS:ATOM-USD",
			Price:     100000000,
			Timestamp: 1234567890,
		},
		{
			SignalID:  "CS:BAND-USD",
			Price:     200000000,
			Timestamp: 1234567890,
		},
	}

	for _, price := range prices {
		suite.feedsKeeper.SetPrice(ctx, *price)
	}

	// query and check
	var (
		req    *types.QueryPricesRequest
		expRes *types.QueryPricesResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"all prices",
			func() {
				req = &types.QueryPricesRequest{}
				expRes = &types.QueryPricesResponse{
					Prices: prices,
				}
			},
			true,
		},
		{
			"limit 1",
			func() {
				req = &types.QueryPricesRequest{
					Pagination: &query.PageRequest{Limit: 1},
				}
				expRes = &types.QueryPricesResponse{
					Prices: prices[:1],
				}
			},
			true,
		},
		{
			"filter",
			func() {
				req = &types.QueryPricesRequest{
					SignalIds: []string{"CS:BAND-USD"},
				}
				expRes = &types.QueryPricesResponse{
					Prices: prices[1:],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			res, err := queryClient.Prices(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetPrices(), res.GetPrices())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(expRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryPrice() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup

	price := types.Price{
		SignalID:  "CS:BAND-USD",
		Price:     100000000,
		Timestamp: 1234567890,
	}
	suite.feedsKeeper.SetPrice(ctx, price)

	// query and check
	res, err := queryClient.Price(context.Background(), &types.QueryPriceRequest{
		SignalId: "CS:BAND-USD",
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryPriceResponse{
		Price: price,
	}, res)

	res, err = queryClient.Price(context.Background(), &types.QueryPriceRequest{
		SignalId: "CS:ATOM-USD",
	})
	suite.Require().ErrorContains(err, "price not found")
	suite.Require().Nil(res)
}

func (suite *KeeperTestSuite) TestQueryValidatorPrices() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup
	feeds := []types.Feed{
		{
			SignalID: "CS:ATOM-USD",
			Interval: 100,
		},
		{
			SignalID: "CS:BAND-USD",
			Interval: 100,
		},
	}

	suite.feedsKeeper.SetCurrentFeeds(ctx, feeds)

	valPrices := []types.ValidatorPrice{
		{
			Validator: ValidValidator.String(),
			SignalID:  "CS:ATOM-USD",
			Price:     1e9,
			Timestamp: ctx.BlockTime().Unix(),
		},
		{
			Validator: ValidValidator.String(),
			SignalID:  "CS:BAND-USD",
			Price:     1e9,
			Timestamp: ctx.BlockTime().Unix(),
		},
	}

	err := suite.feedsKeeper.SetValidatorPriceList(ctx, ValidValidator, valPrices)
	suite.Require().NoError(err)

	// query all prices
	res, err := queryClient.ValidatorPrices(context.Background(), &types.QueryValidatorPricesRequest{
		ValidatorAddress: ValidValidator.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryValidatorPricesResponse{
		ValidatorPrices: valPrices,
	}, res)

	// query with specific SignalIds
	res, err = queryClient.ValidatorPrices(context.Background(), &types.QueryValidatorPricesRequest{
		ValidatorAddress: ValidValidator.String(),
		SignalIds:        []string{"CS:ATOM-USD"},
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryValidatorPricesResponse{
		ValidatorPrices: []types.ValidatorPrice(nil),
	}, res)

	// query with invalid validator
	res, err = queryClient.ValidatorPrices(context.Background(), &types.QueryValidatorPricesRequest{
		ValidatorAddress: InvalidValidator.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryValidatorPricesResponse{
		ValidatorPrices: []types.ValidatorPrice(nil),
	}, res)

	// query with specific SignalIds for invalid validator
	res, err = queryClient.ValidatorPrices(context.Background(), &types.QueryValidatorPricesRequest{
		ValidatorAddress: InvalidValidator.String(),
		SignalIds:        []string{"CS:ATOM-USD"},
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryValidatorPricesResponse{
		ValidatorPrices: []types.ValidatorPrice(nil),
	}, res)
}

func (suite *KeeperTestSuite) TestQueryValidValidator() {
	queryClient := suite.queryClient

	// query and check
	res, err := queryClient.ValidValidator(context.Background(), &types.QueryValidValidatorRequest{
		ValidatorAddress: ValidValidator.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryValidValidatorResponse{
		Valid: true,
	}, res)

	res, err = queryClient.ValidValidator(context.Background(), &types.QueryValidValidatorRequest{
		ValidatorAddress: InvalidValidator.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryValidValidatorResponse{
		Valid: false,
	}, res)
}

func (suite *KeeperTestSuite) TestQuerySignalTotalPowers() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// setup
	signals := []*types.Signal{
		{
			ID:    "CS:ATOM-USD",
			Power: 100000000,
		},
		{
			ID:    "CS:BAND-USD",
			Power: 100000000,
		},
	}

	for _, signal := range signals {
		suite.feedsKeeper.SetSignalTotalPower(ctx, *signal)
	}

	// query and check
	var (
		req    *types.QuerySignalTotalPowersRequest
		expRes *types.QuerySignalTotalPowersResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"all feeds",
			func() {
				req = &types.QuerySignalTotalPowersRequest{}
				expRes = &types.QuerySignalTotalPowersResponse{
					SignalTotalPowers: signals,
				}
			},
			true,
		},
		{
			"limit 1",
			func() {
				req = &types.QuerySignalTotalPowersRequest{
					Pagination: &query.PageRequest{Limit: 1},
				}
				expRes = &types.QuerySignalTotalPowersResponse{
					SignalTotalPowers: signals[:1],
				}
			},
			true,
		},
		{
			"filter",
			func() {
				req = &types.QuerySignalTotalPowersRequest{
					SignalIds: []string{"CS:BAND-USD"},
				}
				expRes = &types.QuerySignalTotalPowersResponse{
					SignalTotalPowers: signals[1:],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			res, err := queryClient.SignalTotalPowers(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.SignalTotalPowers, res.SignalTotalPowers)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryCurrentFeeds() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// query and check
	var (
		req    *types.QueryCurrentFeedsRequest
		expRes *types.QueryCurrentFeedsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"no current feeds",
			func() {
				req = &types.QueryCurrentFeedsRequest{}
				expRes = &types.QueryCurrentFeedsResponse{
					CurrentFeeds: types.CurrentFeedWithDeviations{
						Feeds:               nil,
						LastUpdateTimestamp: ctx.BlockTime().Unix(),
						LastUpdateBlock:     ctx.BlockHeight(),
					},
				}
			},
			true,
		},
		{
			"1 current symbol",
			func() {
				feeds := []types.Feed{
					{
						SignalID: "CS:BAND-USD",
						Power:    36000000000,
						Interval: 100,
					},
				}

				suite.feedsKeeper.SetCurrentFeeds(ctx, feeds)

				feedWithDeviations := []types.FeedWithDeviation{
					{
						SignalID:            "CS:BAND-USD",
						Power:               36000000000,
						Interval:            100,
						DeviationBasisPoint: 83,
					},
				}

				req = &types.QueryCurrentFeedsRequest{}
				expRes = &types.QueryCurrentFeedsResponse{
					CurrentFeeds: types.CurrentFeedWithDeviations{
						Feeds:               feedWithDeviations,
						LastUpdateTimestamp: ctx.BlockTime().Unix(),
						LastUpdateBlock:     ctx.BlockHeight(),
					},
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			res, err := queryClient.CurrentFeeds(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(expRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryReferenceSourceConfig() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// query and check
	res, err := queryClient.ReferenceSourceConfig(context.Background(), &types.QueryReferenceSourceConfigRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryReferenceSourceConfigResponse{
		ReferenceSourceConfig: suite.feedsKeeper.GetReferenceSourceConfig(ctx),
	}, res)
}

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx, queryClient := suite.ctx, suite.queryClient

	// query and check
	res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{
		Params: suite.feedsKeeper.GetParams(ctx),
	}, res)
}
