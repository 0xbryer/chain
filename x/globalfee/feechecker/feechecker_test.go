package feechecker_test

import (
	"math"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/stretchr/testify/suite"

	"github.com/bandprotocol/chain/v2/pkg/tss"
	bandtesting "github.com/bandprotocol/chain/v2/testing"
	"github.com/bandprotocol/chain/v2/x/globalfee/feechecker"
	oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"
	tsstypes "github.com/bandprotocol/chain/v2/x/tss/types"
)

var (
	BasicCalldata = []byte("BASIC_CALLDATA")
	BasicClientID = "BASIC_CLIENT_ID"
)

type StubTx struct {
	sdk.Tx
	sdk.FeeTx
	Msgs      []sdk.Msg
	GasPrices sdk.DecCoins
}

func (st *StubTx) GetMsgs() []sdk.Msg {
	return st.Msgs
}

func (st *StubTx) ValidateBasic() error {
	return nil
}

func (st *StubTx) GetGas() uint64 {
	return 1000000
}

func (st *StubTx) GetFee() sdk.Coins {
	fees := make(sdk.Coins, len(st.GasPrices))

	// Determine the fees by multiplying each gas prices
	glDec := sdk.NewDec(int64(st.GetGas()))
	for i, gp := range st.GasPrices {
		fee := gp.Amount.Mul(glDec)
		fees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
	}

	return fees
}

type FeeCheckerTestSuite struct {
	suite.Suite
	FeeChecker feechecker.FeeChecker
	ctx        sdk.Context
	requestID  oracletypes.RequestID
}

func (suite *FeeCheckerTestSuite) SetupTest() {
	app, ctx := bandtesting.CreateTestApp(suite.T(), true)

	suite.ctx = ctx.WithBlockHeight(999).
		WithIsCheckTx(true).
		WithMinGasPrices(sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(1, 4)}})

	err := app.OracleKeeper.GrantReporter(suite.ctx, bandtesting.Validators[0].ValAddress, bandtesting.Alice.Address)
	suite.Require().NoError(err)

	expiration := ctx.BlockTime().Add(1000 * time.Hour)
	err = app.AuthzKeeper.SaveGrant(
		ctx,
		bandtesting.Alice.Address,
		bandtesting.Validators[0].Address,
		authz.NewGenericAuthorization(
			sdk.MsgTypeURL(&tsstypes.MsgSubmitDEs{}),
		),
		&expiration,
	)
	suite.Require().NoError(err)

	req := oracletypes.NewRequest(
		1,
		BasicCalldata,
		[]sdk.ValAddress{bandtesting.Validators[0].ValAddress},
		1,
		1,
		bandtesting.ParseTime(0),
		"",
		nil,
		nil,
		0,
		0,
		0,
		bandtesting.FeePayer.Address.String(),
		bandtesting.Coins100000000uband,
	)
	suite.requestID = app.OracleKeeper.AddRequest(suite.ctx, req)

	suite.FeeChecker = feechecker.NewFeeChecker(
		&app.AuthzKeeper,
		&app.OracleKeeper,
		&app.GlobalfeeKeeper,
		app.StakingKeeper,
		app.TSSKeeper,
		app.BandtssKeeper,
	)
}

func (suite *FeeCheckerTestSuite) TestValidRawReport() {
	msgs := []sdk.Msg{
		oracletypes.NewMsgReportData(suite.requestID, []oracletypes.RawReport{}, bandtesting.Validators[0].ValAddress),
	}
	stubTx := &StubTx{Msgs: msgs}

	// test - check report tx
	isReportTx := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, msgs[0])
	suite.Require().True(isReportTx)

	// test - check tx fee with min gas prices
	fee, priority, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.Coins{}, fee)
	suite.Require().Equal(int64(math.MaxInt64), priority)
}

func (suite *FeeCheckerTestSuite) TestNotValidRawReport() {
	msgs := []sdk.Msg{oracletypes.NewMsgReportData(1, []oracletypes.RawReport{}, bandtesting.Alice.ValAddress)}
	stubTx := &StubTx{Msgs: msgs}

	// test - check report tx
	isReportTx := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, msgs[0])
	suite.Require().False(isReportTx)

	// test - check tx fee with min gas prices
	_, _, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().Error(err)
}

func (suite *FeeCheckerTestSuite) TestValidReport() {
	reportMsgs := []sdk.Msg{
		oracletypes.NewMsgReportData(suite.requestID, []oracletypes.RawReport{}, bandtesting.Validators[0].ValAddress),
	}
	authzMsg := authz.NewMsgExec(bandtesting.Alice.Address, reportMsgs)
	stubTx := &StubTx{Msgs: []sdk.Msg{&authzMsg}}

	// test - check bypass min fee
	isBypassMinFeeMsg := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, reportMsgs[0])
	suite.Require().True(isBypassMinFeeMsg)

	// test - check tx fee with min gas prices
	fee, priority, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.Coins{}, fee)
	suite.Require().Equal(int64(math.MaxInt64), priority)
}

func (suite *FeeCheckerTestSuite) TestNoAuthzReport() {
	reportMsgs := []sdk.Msg{
		oracletypes.NewMsgReportData(suite.requestID, []oracletypes.RawReport{}, bandtesting.Validators[0].ValAddress),
	}
	authzMsg := authz.NewMsgExec(bandtesting.Bob.Address, reportMsgs)
	stubTx := &StubTx{Msgs: []sdk.Msg{&authzMsg}, GasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uband", sdk.NewInt(1)))}

	// test - check bypass min fee
	isBypassMinFeeMsg := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, &authzMsg)
	suite.Require().False(isBypassMinFeeMsg)

	// test - check tx fee with min gas prices
	_, _, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().NoError(err)
}

func (suite *FeeCheckerTestSuite) TestNotValidReport() {
	reportMsgs := []sdk.Msg{
		oracletypes.NewMsgReportData(suite.requestID+1, []oracletypes.RawReport{}, bandtesting.Validators[0].ValAddress),
	}
	authzMsg := authz.NewMsgExec(bandtesting.Alice.Address, reportMsgs)
	stubTx := &StubTx{Msgs: []sdk.Msg{&authzMsg}}

	// test - check bypass min fee
	isBypassMinFeeMsg := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, &authzMsg)
	suite.Require().False(isBypassMinFeeMsg)

	// test - check tx fee with min gas prices
	_, _, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().Error(err)
}

func (suite *FeeCheckerTestSuite) TestNotReportMsg() {
	requestMsg := oracletypes.NewMsgRequestData(
		1,
		BasicCalldata,
		1,
		1,
		BasicClientID,
		bandtesting.Coins100000000uband,
		bandtesting.TestDefaultPrepareGas,
		bandtesting.TestDefaultExecuteGas,
		bandtesting.FeePayer.Address,
		0,
		0,
	)
	stubTx := &StubTx{
		Msgs: []sdk.Msg{requestMsg},
		GasPrices: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec("uaaaa", sdk.NewDecWithPrec(100, 3)),
			sdk.NewDecCoinFromDec("uaaab", sdk.NewDecWithPrec(1, 3)),
			sdk.NewDecCoinFromDec("uaaac", sdk.NewDecWithPrec(0, 3)),
			sdk.NewDecCoinFromDec("uband", sdk.NewDecWithPrec(3, 3)),
			sdk.NewDecCoinFromDec("uccca", sdk.NewDecWithPrec(0, 3)),
			sdk.NewDecCoinFromDec("ucccb", sdk.NewDecWithPrec(1, 3)),
			sdk.NewDecCoinFromDec("ucccc", sdk.NewDecWithPrec(100, 3)),
		),
	}

	// test - check bypass min fee
	isBypassMinFeeMsg := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, requestMsg)
	suite.Require().False(isBypassMinFeeMsg)

	// test - check tx fee with min gas prices
	fee, priority, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().NoError(err)
	suite.Require().Equal(stubTx.GetFee(), fee)
	suite.Require().Equal(int64(30), priority)
}

func (suite *FeeCheckerTestSuite) TestReportMsgAndOthersTypeMsgInTheSameAuthzMsgs() {
	reportMsg := oracletypes.NewMsgReportData(suite.requestID, []oracletypes.RawReport{}, bandtesting.Validators[0].ValAddress)
	requestMsg := oracletypes.NewMsgRequestData(
		1,
		BasicCalldata,
		1,
		1,
		BasicClientID,
		bandtesting.Coins100000000uband,
		bandtesting.TestDefaultPrepareGas,
		bandtesting.TestDefaultExecuteGas,
		bandtesting.FeePayer.Address,
		0,
		0,
	)
	msgs := []sdk.Msg{reportMsg, requestMsg}
	authzMsg := authz.NewMsgExec(bandtesting.Alice.Address, msgs)
	stubTx := &StubTx{Msgs: []sdk.Msg{&authzMsg}, GasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uband", sdk.NewInt(1)))}

	// test - check bypass min fee
	isBypassMinFeeMsg := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, &authzMsg)
	suite.Require().False(isBypassMinFeeMsg)

	// test - check tx fee with min gas prices
	fee, priority, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().NoError(err)
	suite.Require().Equal(stubTx.GetFee(), fee)
	suite.Require().Equal(int64(10000), priority)
}

func (suite *FeeCheckerTestSuite) TestReportMsgAndOthersTypeMsgInTheSameTx() {
	reportMsg := oracletypes.NewMsgReportData(suite.requestID, []oracletypes.RawReport{}, bandtesting.Validators[0].ValAddress)
	requestMsg := oracletypes.NewMsgRequestData(
		1,
		BasicCalldata,
		1,
		1,
		BasicClientID,
		bandtesting.Coins100000000uband,
		bandtesting.TestDefaultPrepareGas,
		bandtesting.TestDefaultExecuteGas,
		bandtesting.FeePayer.Address,
		0,
		0,
	)
	stubTx := &StubTx{
		Msgs:      []sdk.Msg{reportMsg, requestMsg},
		GasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uband", sdk.NewInt(1))),
	}

	// test - check bypass min fee
	isBypassMinFeeMsg := suite.FeeChecker.IsBypassMinFeeMsg(suite.ctx, requestMsg)
	suite.Require().False(isBypassMinFeeMsg)

	// test - check tx fee with min gas prices
	fee, priority, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
	suite.Require().NoError(err)
	suite.Require().Equal(stubTx.GetFee(), fee)
	suite.Require().Equal(int64(10000), priority)
}

func (suite *FeeCheckerTestSuite) TestIsBypassMinFeeTxAndCheckTxFeeWithMinGasPrices() {
	testCases := []struct {
		name                string
		stubTx              func() *StubTx
		expIsBypassMinFeeTx bool
		expErr              error
		expFee              sdk.Coins
		expPriority         int64
	}{
		{
			name: "valid MsgReportData",
			stubTx: func() *StubTx {
				return &StubTx{
					Msgs: []sdk.Msg{
						oracletypes.NewMsgReportData(
							suite.requestID,
							[]oracletypes.RawReport{},
							bandtesting.Validators[0].ValAddress,
						),
					},
				}
			},
			expIsBypassMinFeeTx: true,
			expErr:              nil,
			expFee:              sdk.Coins{},
			expPriority:         math.MaxInt64,
		},
		{
			name: "valid MsgSubmitDEs",
			stubTx: func() *StubTx {
				privD, _ := tss.GenerateSigningNonce([]byte{})
				privE, _ := tss.GenerateSigningNonce([]byte{})

				return &StubTx{
					Msgs: []sdk.Msg{
						&tsstypes.MsgSubmitDEs{
							DEs: []tsstypes.DE{
								{
									PubD: privD.Point(),
									PubE: privE.Point(),
								},
							},
							Address: bandtesting.Validators[0].Address.String(),
						},
					},
				}
			},
			expIsBypassMinFeeTx: true,
			expErr:              nil,
			expFee:              sdk.Coins{},
			expPriority:         math.MaxInt64,
		},
		{
			name: "invalid MsgSubmitDEs",
			stubTx: func() *StubTx {
				return &StubTx{
					Msgs: []sdk.Msg{
						&tsstypes.MsgSubmitDEs{
							DEs: []tsstypes.DE{
								{
									PubD: nil,
									PubE: nil,
								},
							},
							Address: "wrong address",
						},
					},
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              sdkerrors.ErrInsufficientFee,
			expFee:              nil,
			expPriority:         0,
		},
		{
			name: "valid MsgSubmitDEs in valid MsgExec",
			stubTx: func() *StubTx {
				privD, _ := tss.GenerateSigningNonce([]byte{})
				privE, _ := tss.GenerateSigningNonce([]byte{})

				msgExec := authz.NewMsgExec(bandtesting.Alice.Address, []sdk.Msg{
					&tsstypes.MsgSubmitDEs{
						DEs: []tsstypes.DE{
							{
								PubD: privD.Point(),
								PubE: privE.Point(),
							},
						},
						Address: bandtesting.Validators[0].Address.String(),
					},
				})

				return &StubTx{
					Msgs: []sdk.Msg{
						&msgExec,
					},
				}
			},
			expIsBypassMinFeeTx: true,
			expErr:              nil,
			expFee:              sdk.Coins{},
			expPriority:         math.MaxInt64,
		},
		{
			name: "valid MsgSubmitDEs in invalid MsgExec",
			stubTx: func() *StubTx {
				privD, _ := tss.GenerateSigningNonce([]byte{})
				privE, _ := tss.GenerateSigningNonce([]byte{})

				msgExec := authz.NewMsgExec(bandtesting.Bob.Address, []sdk.Msg{
					&tsstypes.MsgSubmitDEs{
						DEs: []tsstypes.DE{
							{
								PubD: privD.Point(),
								PubE: privE.Point(),
							},
						},
						Address: bandtesting.Validators[0].Address.String(),
					},
				})

				return &StubTx{
					Msgs: []sdk.Msg{
						&msgExec,
					},
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              sdkerrors.ErrInsufficientFee,
			expFee:              nil,
			expPriority:         0,
		},
		{
			name: "valid MsgReportData in valid MsgExec",
			stubTx: func() *StubTx {
				msgExec := authz.NewMsgExec(bandtesting.Alice.Address, []sdk.Msg{
					oracletypes.NewMsgReportData(
						suite.requestID,
						[]oracletypes.RawReport{},
						bandtesting.Validators[0].ValAddress,
					),
				})

				return &StubTx{
					Msgs: []sdk.Msg{
						&msgExec,
					},
				}
			},
			expIsBypassMinFeeTx: true,
			expErr:              nil,
			expFee:              sdk.Coins{},
			expPriority:         math.MaxInt64,
		},
		{
			name: "invalid MsgReportData with not enough fee",
			stubTx: func() *StubTx {
				return &StubTx{
					Msgs: []sdk.Msg{
						oracletypes.NewMsgReportData(1, []oracletypes.RawReport{}, bandtesting.Alice.ValAddress),
					},
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              sdkerrors.ErrInsufficientFee,
			expFee:              nil,
			expPriority:         0,
		},
		{
			name: "invalid MsgReportData in valid MsgExec with not enough fee",
			stubTx: func() *StubTx {
				msgExec := authz.NewMsgExec(bandtesting.Alice.Address, []sdk.Msg{
					oracletypes.NewMsgReportData(
						suite.requestID+1,
						[]oracletypes.RawReport{},
						bandtesting.Validators[0].ValAddress,
					),
				})

				return &StubTx{
					Msgs: []sdk.Msg{
						&msgExec,
					},
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              sdkerrors.ErrInsufficientFee,
			expFee:              nil,
			expPriority:         0,
		},
		{
			name: "valid MsgReportData in invalid MsgExec with enough fee",
			stubTx: func() *StubTx {
				msgExec := authz.NewMsgExec(bandtesting.Bob.Address, []sdk.Msg{
					oracletypes.NewMsgReportData(
						suite.requestID,
						[]oracletypes.RawReport{},
						bandtesting.Validators[0].ValAddress,
					),
				})

				return &StubTx{
					Msgs: []sdk.Msg{
						&msgExec,
					},
					GasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uband", sdk.NewInt(1))),
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              nil,
			expFee:              sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(1000000))),
			expPriority:         10000,
		},
		{
			name: "valid MsgRequestData",
			stubTx: func() *StubTx {
				msgRequestData := oracletypes.NewMsgRequestData(
					1,
					BasicCalldata,
					1,
					1,
					BasicClientID,
					bandtesting.Coins100000000uband,
					bandtesting.TestDefaultPrepareGas,
					bandtesting.TestDefaultExecuteGas,
					bandtesting.FeePayer.Address,
					0,
					0,
				)

				return &StubTx{
					Msgs: []sdk.Msg{msgRequestData},
					GasPrices: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("uaaaa", sdk.NewDecWithPrec(100, 3)),
						sdk.NewDecCoinFromDec("uaaab", sdk.NewDecWithPrec(1, 3)),
						sdk.NewDecCoinFromDec("uaaac", sdk.NewDecWithPrec(0, 3)),
						sdk.NewDecCoinFromDec("uband", sdk.NewDecWithPrec(3, 3)),
						sdk.NewDecCoinFromDec("uccca", sdk.NewDecWithPrec(0, 3)),
						sdk.NewDecCoinFromDec("ucccb", sdk.NewDecWithPrec(1, 3)),
						sdk.NewDecCoinFromDec("ucccc", sdk.NewDecWithPrec(100, 3)),
					),
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              nil,
			expFee: sdk.NewCoins(
				sdk.NewCoin("uaaaa", sdk.NewInt(100000)),
				sdk.NewCoin("uaaab", sdk.NewInt(1000)),
				sdk.NewCoin("uband", sdk.NewInt(3000)),
				sdk.NewCoin("ucccb", sdk.NewInt(1000)),
				sdk.NewCoin("ucccc", sdk.NewInt(100000)),
			),
			expPriority: 30,
		},
		{
			name: "valid MsgRequestData and valid MsgReport in valid MsgExec with enough fee",
			stubTx: func() *StubTx {
				msgReportData := oracletypes.NewMsgReportData(
					suite.requestID,
					[]oracletypes.RawReport{},
					bandtesting.Validators[0].ValAddress,
				)
				msgRequestData := oracletypes.NewMsgRequestData(
					1,
					BasicCalldata,
					1,
					1,
					BasicClientID,
					bandtesting.Coins100000000uband,
					bandtesting.TestDefaultPrepareGas,
					bandtesting.TestDefaultExecuteGas,
					bandtesting.FeePayer.Address,
					0,
					0,
				)
				msgs := []sdk.Msg{msgReportData, msgRequestData}
				authzMsg := authz.NewMsgExec(bandtesting.Alice.Address, msgs)

				return &StubTx{
					Msgs:      []sdk.Msg{&authzMsg},
					GasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uband", sdk.NewInt(1))),
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              nil,
			expFee: sdk.NewCoins(
				sdk.NewCoin("uband", sdk.NewInt(1000000)),
			),
			expPriority: 10000,
		},
		{
			name: "valid MsgRequestData and valid MsgReport with enough fee",
			stubTx: func() *StubTx {
				msgReportData := oracletypes.NewMsgReportData(
					suite.requestID,
					[]oracletypes.RawReport{},
					bandtesting.Validators[0].ValAddress,
				)
				msgRequestData := oracletypes.NewMsgRequestData(
					1,
					BasicCalldata,
					1,
					1,
					BasicClientID,
					bandtesting.Coins100000000uband,
					bandtesting.TestDefaultPrepareGas,
					bandtesting.TestDefaultExecuteGas,
					bandtesting.FeePayer.Address,
					0,
					0,
				)

				return &StubTx{
					Msgs:      []sdk.Msg{msgReportData, msgRequestData},
					GasPrices: sdk.NewDecCoins(sdk.NewDecCoin("uband", sdk.NewInt(1))),
				}
			},
			expIsBypassMinFeeTx: false,
			expErr:              nil,
			expFee: sdk.NewCoins(
				sdk.NewCoin("uband", sdk.NewInt(1000000)),
			),
			expPriority: 10000,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			stubTx := tc.stubTx()

			// test - IsByPassMinFeeTx
			isByPassMinFeeTx := suite.FeeChecker.IsBypassMinFeeTx(suite.ctx, stubTx)
			suite.Require().Equal(tc.expIsBypassMinFeeTx, isByPassMinFeeTx)

			// test - CheckTxFeeWithMinGasPrices
			fee, priority, err := suite.FeeChecker.CheckTxFeeWithMinGasPrices(suite.ctx, stubTx)
			suite.Require().ErrorIs(err, tc.expErr)
			suite.Require().Equal(fee, tc.expFee)
			suite.Require().Equal(tc.expPriority, priority)
		})
	}
}

func (suite *FeeCheckerTestSuite) TestGetBondDenom() {
	denom := suite.FeeChecker.GetBondDenom(suite.ctx)
	suite.Require().Equal("uband", denom)
}

func (suite *FeeCheckerTestSuite) TestDefaultZeroGlobalFee() {
	coins, err := suite.FeeChecker.DefaultZeroGlobalFee(suite.ctx)

	suite.Require().Equal(1, len(coins))
	suite.Require().Equal("uband", coins[0].Denom)
	suite.Require().Equal(sdk.NewDec(0), coins[0].Amount)
	suite.Require().NoError(err)
}

func TestFeeCheckerTestSuite(t *testing.T) {
	suite.Run(t, new(FeeCheckerTestSuite))
}
