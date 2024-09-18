package oracle_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	bandtest "github.com/bandprotocol/chain/v3/app"
	oracletypes "github.com/bandprotocol/chain/v3/x/oracle/types"
)

func init() {
	bandtest.SetBech32AddressPrefixesAndBip44CoinTypeAndSeal(sdk.GetConfig())
	sdk.DefaultBondDenom = "uband"
}

type IBCTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path

	// shortcut to chainB (bandchain)
	bandApp *bandtest.BandApp
}

func (suite *IBCTestSuite) SetupTest() {
	ibctesting.DefaultTestingAppInit = bandtest.CreateTestingAppFn(suite.T())

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.path.EndpointA.ChannelConfig.PortID = oracletypes.ModuleName
	suite.path.EndpointA.ChannelConfig.Version = oracletypes.Version
	suite.path.EndpointB.ChannelConfig.PortID = oracletypes.ModuleName
	suite.path.EndpointB.ChannelConfig.Version = oracletypes.Version

	suite.bandApp = suite.chainB.App.(*bandtest.BandApp)

	suite.coordinator.Setup(suite.path)

	// Activate oracle validator on chain B (bandchain)
	for _, v := range suite.chainB.Vals.Validators {
		err := suite.bandApp.OracleKeeper.Activate(
			suite.chainB.GetContext(),
			sdk.ValAddress(v.Address),
		)
		suite.Require().NoError(err)
	}

	suite.coordinator.CommitBlock(suite.chainB)
}

func (suite *IBCTestSuite) sendReport(requestID oracletypes.RequestID, report oracletypes.Report, needToResolve bool) {
	suite.bandApp.OracleKeeper.SetReport(suite.chainB.GetContext(), requestID, report)
	if needToResolve {
		suite.bandApp.OracleKeeper.AddPendingRequest(suite.chainB.GetContext(), requestID)
	}

	suite.coordinator.CommitBlock(suite.chainB)
}

func (suite *IBCTestSuite) sendOracleRequestPacket(
	path *ibctesting.Path,
	seq uint64,
	oracleRequestPacket oracletypes.OracleRequestPacketData,
	timeoutHeight clienttypes.Height,
) channeltypes.Packet {
	packet := channeltypes.NewPacket(
		oracleRequestPacket.GetBytes(),
		seq,
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		timeoutHeight,
		0,
	)
	_, err := path.EndpointA.SendPacket(timeoutHeight, 0, oracleRequestPacket.GetBytes())
	suite.Require().NoError(err)
	return packet
}

func (suite *IBCTestSuite) checkChainBTreasuryBalances(expect sdk.Coins) {
	treasuryBalances := suite.bandApp.BankKeeper.GetAllBalances(
		suite.chainB.GetContext(),
		bandtest.Treasury.Address,
	)
	suite.Require().Equal(expect, treasuryBalances)
}

func (suite *IBCTestSuite) checkChainBSenderBalances(expect sdk.Coins) {
	b := suite.bandApp.BankKeeper.GetAllBalances(
		suite.chainB.GetContext(),
		suite.chainB.SenderAccount.GetAddress(),
	)
	suite.Require().Equal(expect, b)
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *IBCTestSuite) TestHandleIBCRequestSuccess() {
	path := suite.path
	// send request from A to B
	timeoutHeight := clienttypes.NewHeight(10, 110)
	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
		path.EndpointA.ClientID,
		1,
		[]byte("test"),
		4,
		2,
		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(12000000))),
		bandtest.TestDefaultPrepareGas,
		bandtest.TestDefaultExecuteGas,
	)
	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

	err := path.RelayPacket(packet)
	suite.Require().NoError(err) // relay committed

	// Treasury get fees from relayer
	suite.checkChainBTreasuryBalances(sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(12000000))))

	raws1 := []oracletypes.RawReport{
		oracletypes.NewRawReport(1, 0, []byte("data1")),
		oracletypes.NewRawReport(2, 0, []byte("data2")),
		oracletypes.NewRawReport(3, 0, []byte("data3")),
	}
	suite.sendReport(
		oracletypes.RequestID(1),
		oracletypes.NewReport(sdk.ValAddress(suite.chainB.Vals.Validators[0].Address), true, raws1),
		false,
	)

	raws2 := []oracletypes.RawReport{
		oracletypes.NewRawReport(1, 0, []byte("data1")),
		oracletypes.NewRawReport(2, 0, []byte("data2")),
		oracletypes.NewRawReport(3, 0, []byte("data3")),
	}
	suite.sendReport(
		oracletypes.RequestID(1),
		oracletypes.NewReport(sdk.ValAddress(suite.chainB.Vals.Validators[2].Address), true, raws2),
		true,
	)

	oracleResponsePacket := oracletypes.NewOracleResponsePacketData(
		path.EndpointA.ClientID,
		1,
		2,
		1577923360,
		1577923385,
		oracletypes.RESOLVE_STATUS_SUCCESS,
		[]byte("test"),
	)
	responsePacket := channeltypes.NewPacket(
		oracleResponsePacket.GetBytes(),
		1,
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		clienttypes.ZeroHeight(),
		uint64(time.Unix(1577923385, 0).Add(10*time.Minute).UnixNano()),
	)
	expectCommitment := channeltypes.CommitPacket(suite.chainB.Codec, responsePacket)
	commitment := suite.chainB.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(
		suite.chainB.GetContext(),
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		1,
	)
	suite.Equal(expectCommitment, commitment)
}

// func (suite *OracleTestSuite) TestIBCPrepareValidateBasicFail() {
// 	path := suite.path

// 	clientID := path.EndpointA.ClientID
// 	coins := sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(6000000)))

// 	oracleRequestPackets := []oracletypes.OracleRequestPacketData{
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte(strings.Repeat("beeb", 130)),
// 			1,
// 			1,
// 			coins,
// 			bandtesting.TestDefaultPrepareGas,
// 			bandtesting.TestDefaultExecuteGas,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			0,
// 			coins,
// 			bandtesting.TestDefaultPrepareGas,
// 			bandtesting.TestDefaultExecuteGas,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			2,
// 			coins,
// 			bandtesting.TestDefaultPrepareGas,
// 			bandtesting.TestDefaultExecuteGas,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			strings.Repeat(clientID, 9),
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			1,
// 			coins,
// 			bandtesting.TestDefaultPrepareGas,
// 			bandtesting.TestDefaultExecuteGas,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			1,
// 			coins,
// 			0,
// 			bandtesting.TestDefaultExecuteGas,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			1,
// 			coins,
// 			bandtesting.TestDefaultPrepareGas,
// 			0,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			1,
// 			coins,
// 			oracletypes.MaximumOwasmGas,
// 			oracletypes.MaximumOwasmGas,
// 		),
// 		oracletypes.NewOracleRequestPacketData(
// 			clientID,
// 			1,
// 			[]byte("beeb"),
// 			1,
// 			1,
// 			bandtesting.BadCoins,
// 			bandtesting.TestDefaultPrepareGas,
// 			bandtesting.TestDefaultExecuteGas,
// 		),
// 	}

// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	for i, requestPacket := range oracleRequestPackets {
// 		packet := suite.sendOracleRequestPacket(path, uint64(i)+1, requestPacket, timeoutHeight)

// 		err := path.RelayPacket(packet)
// 		suite.Require().NoError(err) // relay committed
// 	}
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestNotEnoughFund() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)

// 	// Use Carol as a relayer
// 	carol := bandtesting.Carol
// 	carolExpectedBalance := sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(2492500)))
// 	_, err := suite.chainB.SendMsgs(banktypes.NewMsgSend(
// 		suite.chainB.SenderAccount.GetAddress(),
// 		carol.Address,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(2500000))),
// 	))
// 	suite.Require().NoError(err)

// 	suite.chainB.SenderPrivKey = carol.PrivKey
// 	suite.chainB.SenderAccount = suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), carol.Address)

// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err = path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed

// 	carolBalance := suite.chainB.App.BankKeeper.GetAllBalances(suite.chainB.GetContext(), carol.Address)
// 	suite.Require().Equal(carolExpectedBalance, carolBalance)
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestNotEnoughFeeLimit() {
// 	path := suite.path
// 	expectedBalance := suite.chainB.App.BankKeeper.GetAllBalances(
// 		suite.chainB.GetContext(),
// 		suite.chainB.SenderAccount.GetAddress(),
// 	).Sub(sdk.NewCoin("uband", math.NewInt(7500)))

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(2000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed

// 	suite.checkChainBSenderBalances(expectedBalance)
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestInvalidCalldataSize() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte(strings.Repeat("beeb", 2000)),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestNotEnoughPrepareGas() {
// 	path := suite.path
// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		1,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestInvalidAskCountFail() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		17,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed

// 	oracleRequestPacket = oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		3,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet = suite.sendOracleRequestPacket(path, 2, oracleRequestPacket, timeoutHeight)

// 	err = path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestBaseOwasmFeePanic() {
// 	path := suite.path

// 	params := suite.chainB.App.OracleKeeper.GetParams(suite.chainB.GetContext())
// 	params.BaseOwasmGas = 100000000
// 	params.PerValidatorRequestGas = 0
// 	err := suite.chainB.App.OracleKeeper.SetParams(suite.chainB.GetContext(), params)
// 	suite.Require().NoError(err)

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	// ConsumeGas panics due to insufficient gas, so ErrAcknowledgement is not created.
// 	err = path.RelayPacket(packet)
// 	suite.Require().Contains(err.Error(), "BASE_OWASM_FEE; gasWanted: 1000000")
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestPerValidatorRequestFeePanic() {
// 	path := suite.path

// 	params := suite.chainB.App.OracleKeeper.GetParams(suite.chainB.GetContext())
// 	params.PerValidatorRequestGas = 100000000
// 	err := suite.chainB.App.OracleKeeper.SetParams(suite.chainB.GetContext(), params)
// 	suite.Require().NoError(err)

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	// ConsumeGas panics due to insufficient gas, so ErrAcknowledgement is not created.
// 	err = path.RelayPacket(packet)
// 	suite.Require().Contains(err.Error(), "PER_VALIDATOR_REQUEST_FEE; gasWanted: 1000000")
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestOracleScriptNotFound() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		100,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestBadWasmExecutionFail() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		2,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestWithEmptyRawRequest() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		3,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestUnknownDataSource() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		4,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestInvalidDataSourceCount() {
// 	path := suite.path

// 	params := suite.chainB.App.OracleKeeper.GetParams(suite.chainB.GetContext())
// 	params.MaxRawRequestCount = 3
// 	err := suite.chainB.App.OracleKeeper.SetParams(suite.chainB.GetContext(), params)
// 	suite.Require().NoError(err)

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		4,
// 		obi.MustEncode(testdata.Wasm4Input{
// 			IDs:      []int64{1, 2, 3, 4},
// 			Calldata: "beeb",
// 		}),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(4000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err = path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestTooMuchWasmGas() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		6,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCPrepareRequestTooLargeCalldata() {
// 	path := suite.path
// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		8,
// 		[]byte("beeb"),
// 		1,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed
// }

// func (suite *OracleTestSuite) TestIBCResolveRequestOutOfGas() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		[]byte("beeb"),
// 		2,
// 		1,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(6000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		1,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed

// 	suite.checkChainBTreasuryBalances(sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(6000000))))
// 	suite.checkChainBSenderBalances(sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(3970000))))

// 	raws := []oracletypes.RawReport{
// 		oracletypes.NewRawReport(1, 0, []byte("data1")),
// 		oracletypes.NewRawReport(2, 0, []byte("data2")),
// 		oracletypes.NewRawReport(3, 0, []byte("data3")),
// 	}
// 	_, err = suite.chainB.SendReport(1, raws, bandtesting.Validators[0])
// 	suite.Require().NoError(err)

// 	commitment := suite.chainB.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(
// 		suite.chainB.GetContext(),
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		1,
// 	)

// 	oracleResponsePacket := oracletypes.NewOracleResponsePacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		1,
// 		1577923380,
// 		1577923400,
// 		oracletypes.RESOLVE_STATUS_FAILURE,
// 		[]byte{},
// 	)
// 	responsePacket := channeltypes.NewPacket(
// 		oracleResponsePacket.GetBytes(),
// 		1,
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		path.EndpointA.ChannelConfig.PortID,
// 		path.EndpointA.ChannelID,
// 		clienttypes.ZeroHeight(),
// 		1577924000000000000,
// 	)
// 	expectCommitment := channeltypes.CommitPacket(suite.chainB.Codec, responsePacket)
// 	suite.Equal(expectCommitment, commitment)
// }

// func (suite *OracleTestSuite) TestIBCResolveReadNilExternalData() {
// 	path := suite.path

// 	// send request from A to B
// 	timeoutHeight := clienttypes.NewHeight(0, 110)
// 	oracleRequestPacket := oracletypes.NewOracleRequestPacketData(
// 		path.EndpointA.ClientID,
// 		4,
// 		obi.MustEncode(testdata.Wasm4Input{IDs: []int64{1, 2}, Calldata: string("beeb")}),
// 		2,
// 		2,
// 		sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(4000000))),
// 		bandtesting.TestDefaultPrepareGas,
// 		bandtesting.TestDefaultExecuteGas,
// 	)
// 	packet := suite.sendOracleRequestPacket(path, 1, oracleRequestPacket, timeoutHeight)

// 	err := path.RelayPacket(packet)
// 	suite.Require().NoError(err) // relay committed

// 	suite.checkChainBTreasuryBalances(sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(4000000))))
// 	suite.checkChainBSenderBalances(sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(5970000))))

// 	raws1 := []oracletypes.RawReport{oracletypes.NewRawReport(0, 0, nil), oracletypes.NewRawReport(1, 0, []byte("beebd2v1"))}
// 	_, err = suite.chainB.SendReport(1, raws1, bandtesting.Validators[0])
// 	suite.Require().NoError(err)

// 	raws2 := []oracletypes.RawReport{oracletypes.NewRawReport(0, 0, []byte("beebd1v2")), oracletypes.NewRawReport(1, 0, nil)}
// 	_, err = suite.chainB.SendReport(1, raws2, bandtesting.Validators[1])
// 	suite.Require().NoError(err)

// 	commitment := suite.chainB.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(
// 		suite.chainB.GetContext(),
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		1,
// 	)

// 	oracleResponsePacket := oracletypes.NewOracleResponsePacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		2,
// 		1577923380,
// 		1577923405,
// 		oracletypes.RESOLVE_STATUS_SUCCESS,
// 		obi.MustEncode(testdata.Wasm4Output{Ret: "beebd1v2beebd2v1"}),
// 	)
// 	responsePacket := channeltypes.NewPacket(
// 		oracleResponsePacket.GetBytes(),
// 		1,
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		path.EndpointA.ChannelConfig.PortID,
// 		path.EndpointA.ChannelID,
// 		clienttypes.ZeroHeight(),
// 		1577924005000000000,
// 	)
// 	expectCommitment := channeltypes.CommitPacket(suite.chainB.Codec, responsePacket)
// 	suite.Equal(expectCommitment, commitment)
// }

// func (suite *OracleTestSuite) TestIBCResolveRequestNoReturnData() {
// 	path := suite.path

// 	suite.chainB.App.OracleKeeper.SetRequest(suite.chainB.GetContext(), 1, oracletypes.NewRequest(
// 		// 3rd Wasm - do nothing
// 		3,
// 		[]byte("beeb"),
// 		[]sdk.ValAddress{bandtesting.Validators[0].ValAddress, bandtesting.Validators[1].ValAddress},
// 		1,
// 		suite.chainB.GetContext().
// 			BlockHeight()-
// 			1,
// 		bandtesting.ParseTime(1577923380),
// 		path.EndpointA.ClientID,
// 		[]oracletypes.RawRequest{
// 			oracletypes.NewRawRequest(1, 1, []byte("beeb")),
// 		},
// 		&oracletypes.IBCChannel{PortId: path.EndpointB.ChannelConfig.PortID, ChannelId: path.EndpointB.ChannelID},
// 		0,
// 	))

// 	raws := []oracletypes.RawReport{oracletypes.NewRawReport(1, 0, []byte("beeb"))}
// 	_, err := suite.chainB.SendReport(1, raws, bandtesting.Validators[0])
// 	suite.Require().NoError(err)

// 	commitment := suite.chainB.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(
// 		suite.chainB.GetContext(),
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		1,
// 	)

// 	oracleResponsePacket := oracletypes.NewOracleResponsePacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		1,
// 		1577923380,
// 		1577923355,
// 		oracletypes.RESOLVE_STATUS_FAILURE,
// 		[]byte{},
// 	)
// 	responsePacket := channeltypes.NewPacket(
// 		oracleResponsePacket.GetBytes(),
// 		1,
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		path.EndpointA.ChannelConfig.PortID,
// 		path.EndpointA.ChannelID,
// 		clienttypes.ZeroHeight(),
// 		1577923955000000000,
// 	)
// 	expectCommitment := channeltypes.CommitPacket(suite.chainB.Codec, responsePacket)
// 	suite.Equal(expectCommitment, commitment)
// }

// func (suite *OracleTestSuite) TestIBCResolveRequestWasmFailure() {
// 	path := suite.path

// 	suite.chainB.App.OracleKeeper.SetRequest(suite.chainB.GetContext(), 1, oracletypes.NewRequest(
// 		// 6th Wasm - out-of-gas
// 		6,
// 		[]byte("beeb"),
// 		[]sdk.ValAddress{bandtesting.Validators[0].ValAddress, bandtesting.Validators[1].ValAddress},
// 		1,
// 		suite.chainB.GetContext().
// 			BlockHeight()-
// 			1,
// 		bandtesting.ParseTime(1577923380),
// 		path.EndpointA.ClientID,
// 		[]oracletypes.RawRequest{
// 			oracletypes.NewRawRequest(1, 1, []byte("beeb")),
// 		},
// 		&oracletypes.IBCChannel{PortId: path.EndpointB.ChannelConfig.PortID, ChannelId: path.EndpointB.ChannelID},
// 		bandtesting.TestDefaultExecuteGas,
// 	))

// 	raws := []oracletypes.RawReport{oracletypes.NewRawReport(1, 0, []byte("beeb"))}
// 	_, err := suite.chainB.SendReport(1, raws, bandtesting.Validators[0])
// 	suite.Require().NoError(err)

// 	commitment := suite.chainB.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(
// 		suite.chainB.GetContext(),
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		1,
// 	)

// 	oracleResponsePacket := oracletypes.NewOracleResponsePacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		1,
// 		1577923380,
// 		1577923355,
// 		oracletypes.RESOLVE_STATUS_FAILURE,
// 		[]byte{},
// 	)
// 	responsePacket := channeltypes.NewPacket(
// 		oracleResponsePacket.GetBytes(),
// 		1,
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		path.EndpointA.ChannelConfig.PortID,
// 		path.EndpointA.ChannelID,
// 		clienttypes.ZeroHeight(),
// 		1577923955000000000,
// 	)
// 	expectCommitment := channeltypes.CommitPacket(suite.chainB.Codec, responsePacket)
// 	suite.Equal(expectCommitment, commitment)
// }

// func (suite *OracleTestSuite) TestIBCResolveRequestCallReturnDataSeveralTimes() {
// 	path := suite.path

// 	suite.chainB.App.OracleKeeper.SetRequest(suite.chainB.GetContext(), 1, oracletypes.NewRequest(
// 		// 9th Wasm - set return data several times
// 		9,
// 		[]byte("beeb"),
// 		[]sdk.ValAddress{bandtesting.Validators[0].ValAddress, bandtesting.Validators[1].ValAddress},
// 		1,
// 		suite.chainB.GetContext().
// 			BlockHeight()-
// 			1,
// 		bandtesting.ParseTime(1577923380),
// 		path.EndpointA.ClientID,
// 		[]oracletypes.RawRequest{
// 			oracletypes.NewRawRequest(1, 1, []byte("beeb")),
// 		},
// 		&oracletypes.IBCChannel{PortId: path.EndpointB.ChannelConfig.PortID, ChannelId: path.EndpointB.ChannelID},
// 		bandtesting.TestDefaultExecuteGas,
// 	))

// 	raws := []oracletypes.RawReport{oracletypes.NewRawReport(1, 0, []byte("beeb"))}
// 	_, err := suite.chainB.SendReport(1, raws, bandtesting.Validators[0])
// 	suite.Require().NoError(err)

// 	commitment := suite.chainB.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(
// 		suite.chainB.GetContext(),
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		1,
// 	)

// 	oracleResponsePacket := oracletypes.NewOracleResponsePacketData(
// 		path.EndpointA.ClientID,
// 		1,
// 		1,
// 		1577923380,
// 		1577923355,
// 		oracletypes.RESOLVE_STATUS_FAILURE,
// 		[]byte{},
// 	)
// 	responsePacket := channeltypes.NewPacket(
// 		oracleResponsePacket.GetBytes(),
// 		1,
// 		path.EndpointB.ChannelConfig.PortID,
// 		path.EndpointB.ChannelID,
// 		path.EndpointA.ChannelConfig.PortID,
// 		path.EndpointA.ChannelID,
// 		clienttypes.ZeroHeight(),
// 		1577923955000000000,
// 	)
// 	expectCommitment := channeltypes.CommitPacket(suite.chainB.Codec, responsePacket)
// 	suite.Equal(expectCommitment, commitment)
// }

func TestIBCTestSuite(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}
