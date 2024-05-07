package keeper_test

import (
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/bandprotocol/chain/v2/pkg/tss"
	"github.com/bandprotocol/chain/v2/pkg/tss/testutil"
	bandtesting "github.com/bandprotocol/chain/v2/testing"
	bandtsskeeper "github.com/bandprotocol/chain/v2/x/bandtss/keeper"
	bandtsstypes "github.com/bandprotocol/chain/v2/x/bandtss/types"
	"github.com/bandprotocol/chain/v2/x/tss/keeper"
	"github.com/bandprotocol/chain/v2/x/tss/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *bandtesting.TestingApp
	ctx         sdk.Context
	queryClient types.QueryClient
	msgSrvr     types.MsgServer
	authority   sdk.AccAddress
}

var (
	PrivD = testutil.HexDecode("de6aedbe8ba688dd6d342881eb1e67c3476e825106477360148e2858a5eb565c")
	PrivE = testutil.HexDecode("3ff4fb2beac0cee0ab230829a5ae0881310046282e79c978ca22f44897ea434a")
	PubD  = tss.Scalar(PrivD).Point()
	PubE  = tss.Scalar(PrivE).Point()
)

func (s *KeeperTestSuite) SetupTest() {
	app, ctx := bandtesting.CreateTestApp(s.T(), true)
	s.app = app
	s.ctx = ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(app.TSSKeeper))
	queryClient := types.NewQueryClient(queryHelper)

	s.queryClient = queryClient
	s.msgSrvr = keeper.NewMsgServerImpl(app.TSSKeeper)
	s.authority = authtypes.NewModuleAddress(govtypes.ModuleName)
}

func (s *KeeperTestSuite) setupCreateGroup() {
	ctx, bandtssKeeper := s.ctx, s.app.BandtssKeeper
	bandtssMsgSrvr := bandtsskeeper.NewMsgServerImpl(bandtssKeeper)

	// Create group from testutil
	for _, tc := range testutil.TestCases {
		// Initialize members
		var members []string
		for _, m := range tc.Group.Members {
			members = append(members, sdk.AccAddress(m.PubKey()).String())
		}

		// Create group
		_, err := bandtssMsgSrvr.CreateGroup(ctx, &bandtsstypes.MsgCreateGroup{
			Members:   members,
			Threshold: tc.Group.Threshold,
			Authority: s.authority.String(),
		})
		s.Require().NoError(err)

		// Set DKG context
		s.app.TSSKeeper.SetDKGContext(ctx, tc.Group.ID, tc.Group.DKGContext)
	}
}

func (s *KeeperTestSuite) setupRound1() {
	s.setupCreateGroup()

	ctx, app, msgSrvr := s.ctx, s.app, s.msgSrvr
	for _, tc := range testutil.TestCases {
		for _, m := range tc.Group.Members {
			// Submit Round 1 information for each member
			_, err := msgSrvr.SubmitDKGRound1(ctx, &types.MsgSubmitDKGRound1{
				GroupID: tc.Group.ID,
				Round1Info: types.Round1Info{
					MemberID:           m.ID,
					CoefficientCommits: m.CoefficientCommits,
					OneTimePubKey:      m.OneTimePubKey(),
					A0Signature:        m.A0Signature,
					OneTimeSignature:   m.OneTimeSignature,
				},
				Address: sdk.AccAddress(m.PubKey()).String(),
			})
			s.Require().NoError(err)
		}
	}

	// Execute the EndBlocker to process groups
	app.EndBlocker(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight() + 1})
}

func (s *KeeperTestSuite) setupRound2() {
	s.setupRound1()

	ctx, app, msgSrvr := s.ctx, s.app, s.msgSrvr
	for _, tc := range testutil.TestCases {
		for _, m := range tc.Group.Members {
			// Submit Round 2 information for each member
			_, err := msgSrvr.SubmitDKGRound2(ctx, &types.MsgSubmitDKGRound2{
				GroupID: tc.Group.ID,
				Round2Info: types.Round2Info{
					MemberID:              m.ID,
					EncryptedSecretShares: m.EncSecretShares,
				},
				Address: sdk.AccAddress(m.PubKey()).String(),
			})
			s.Require().NoError(err)
		}
	}

	// Execute the EndBlocker to process groups
	app.EndBlocker(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight() + 1})
}

func (s *KeeperTestSuite) setupConfirm() {
	s.setupRound2()

	ctx, app, msgSrvr := s.ctx, s.app, s.msgSrvr
	for _, tc := range testutil.TestCases {
		for _, m := range tc.Group.Members {
			// Confirm the group participation for each member
			_, err := msgSrvr.Confirm(ctx, &types.MsgConfirm{
				GroupID:      tc.Group.ID,
				MemberID:     m.ID,
				OwnPubKeySig: m.PubKeySignature,
				Address:      sdk.AccAddress(m.PubKey()).String(),
			})
			s.Require().NoError(err)
		}
	}

	// Execute the EndBlocker to process groups
	app.EndBlocker(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight() + 1})
}

func (s *KeeperTestSuite) setupDE() {
	ctx, msgSrvr := s.ctx, s.msgSrvr

	for _, tc := range testutil.TestCases {
		for _, m := range tc.Group.Members {
			// Submit DEs for each member
			_, err := msgSrvr.SubmitDEs(ctx, &types.MsgSubmitDEs{
				DEs: []types.DE{
					{PubD: PubD, PubE: PubE},
					{PubD: PubD, PubE: PubE},
					{PubD: PubD, PubE: PubE},
					{PubD: PubD, PubE: PubE},
					{PubD: PubD, PubE: PubE},
				},
				Address: sdk.AccAddress(m.PubKey()).String(),
			})
			s.Require().NoError(err)
		}
	}
}

func (s *KeeperTestSuite) SetupGroup(groupStatus types.GroupStatus) {
	switch groupStatus {
	case types.GROUP_STATUS_ROUND_1:
		s.setupCreateGroup()
	case types.GROUP_STATUS_ROUND_2:
		s.setupRound1()
	case types.GROUP_STATUS_ROUND_3:
		s.setupRound2()
	case types.GROUP_STATUS_ACTIVE:
		s.setupConfirm()
		s.setupDE()
	}
}

func (s *KeeperTestSuite) TestGetSetGroupCount() {
	ctx, k := s.ctx, s.app.TSSKeeper
	k.SetGroupCount(ctx, 1)

	groupCount := k.GetGroupCount(ctx)
	s.Require().Equal(uint64(1), groupCount)
}

func (s *KeeperTestSuite) TestGetNextGroupID() {
	ctx, k := s.ctx, s.app.TSSKeeper

	// Initial group count
	k.SetGroupCount(ctx, 0)

	groupID1 := k.GetNextGroupID(ctx)
	s.Require().Equal(tss.GroupID(1), groupID1)
	groupID2 := k.GetNextGroupID(ctx)
	s.Require().Equal(tss.GroupID(2), groupID2)
}

func (s *KeeperTestSuite) TestIsGrantee() {
	ctx, k := s.ctx, s.app.TSSKeeper
	expTime := time.Unix(0, 0)

	// Init grantee address
	grantee, _ := sdk.AccAddressFromBech32("band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs")

	// Init granter address
	granter, _ := sdk.AccAddressFromBech32("band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun")

	// Save grant msgs to grantee
	for _, m := range types.GetTSSGrantMsgTypes() {
		err := s.app.AuthzKeeper.SaveGrant(ctx, grantee, granter, authz.NewGenericAuthorization(m), &expTime)
		s.Require().NoError(err)
	}

	isGrantee := k.CheckIsGrantee(ctx, granter, grantee)
	s.Require().True(isGrantee)
}

func (s *KeeperTestSuite) TestCreateNewGroup() {
	ctx, k := s.ctx, s.app.TSSKeeper

	group := types.Group{
		Size_:     5,
		Threshold: 3,
		PubKey:    nil,
		Status:    types.GROUP_STATUS_ROUND_1,
	}

	// Create new group
	groupID := k.CreateNewGroup(ctx, group)

	// init group ID
	group.ID = groupID

	// Get group by id
	got, err := k.GetGroup(ctx, groupID)
	s.Require().NoError(err)
	s.Require().Equal(group, got)
}

func (s *KeeperTestSuite) TestSetGroup() {
	ctx, k := s.ctx, s.app.TSSKeeper
	group := types.Group{
		Size_:     5,
		Threshold: 3,
		PubKey:    nil,
		Status:    types.GROUP_STATUS_ROUND_1,
	}

	// Set new group
	groupID := k.CreateNewGroup(ctx, group)

	// Update group size value
	group.Size_ = 6

	// Add group ID
	group.ID = groupID

	k.SetGroup(ctx, group)

	// Get group from chain state
	got, err := k.GetGroup(ctx, groupID)

	// Validate group size value
	s.Require().NoError(err)
	s.Require().Equal(group.Size_, got.Size_)
}

func (s *KeeperTestSuite) TestGetGroups() {
	ctx, k := s.ctx, s.app.TSSKeeper
	group := types.Group{
		ID:        1,
		Size_:     5,
		Threshold: 3,
		PubKey:    nil,
		Status:    types.GROUP_STATUS_ROUND_1,
	}

	// Set new group
	k.SetGroup(ctx, group)

	// Get group from chain state
	got := k.GetGroups(ctx)
	s.Require().Equal([]types.Group{group}, got)
}

func (s *KeeperTestSuite) TestDeleteGroup() {
	ctx, k := s.ctx, s.app.TSSKeeper

	// Create a sample group ID
	groupID := tss.GroupID(123)

	// Set up a sample group in the store
	group := types.Group{
		ID: groupID,
		// Set other fields as needed
	}
	k.SetGroup(ctx, group)

	// Delete the group
	k.DeleteGroup(ctx, groupID)

	// Verify that the group is deleted
	_, err := k.GetGroup(ctx, groupID)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestGetSetDKGContext() {
	ctx, k := s.ctx, s.app.TSSKeeper

	dkgContext := []byte("dkg-context sample")
	k.SetDKGContext(ctx, 1, dkgContext)

	got, err := k.GetDKGContext(ctx, 1)
	s.Require().NoError(err)
	s.Require().Equal(dkgContext, got)
}

func (s *KeeperTestSuite) TestGetSetMember() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID, memberID := tss.GroupID(1), tss.MemberID(1)
	member := types.Member{
		ID:          1,
		GroupID:     groupID,
		Address:     "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
		PubKey:      nil,
		IsMalicious: false,
	}
	k.SetMember(ctx, member)

	got, err := k.GetMember(ctx, groupID, memberID)
	s.Require().NoError(err)
	s.Require().Equal(member, got)
}

func (s *KeeperTestSuite) TestGetMembers() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	members := []types.Member{
		{
			ID:          1,
			GroupID:     groupID,
			Address:     "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
			PubKey:      nil,
			IsMalicious: false,
		},
		{
			ID:          2,
			GroupID:     groupID,
			Address:     "band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun",
			PubKey:      nil,
			IsMalicious: false,
		},
	}

	// Set members
	for _, m := range members {
		k.SetMember(ctx, m)
	}

	got, err := k.GetGroupMembers(ctx, groupID)
	s.Require().NoError(err)
	s.Require().Equal(members, got)
}

func (s *KeeperTestSuite) TestGetSetMemberIsActive() {
	ctx, k := s.ctx, s.app.TSSKeeper

	groupID := tss.GroupID(10)
	address := sdk.MustAccAddressFromBech32("band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs")
	k.SetMember(ctx, types.Member{
		ID:       tss.MemberID(1),
		GroupID:  groupID,
		Address:  address.String(),
		PubKey:   nil,
		IsActive: true,
	})

	// check when being set to active
	members, err := k.GetGroupMembers(ctx, groupID)
	s.Require().NoError(err)
	s.Require().Len(members, 1)

	for _, member := range members {
		s.Require().True(member.IsActive)
	}

	err = k.SetMemberIsActive(ctx, groupID, address, false)
	s.Require().NoError(err)

	members, err = k.GetGroupMembers(ctx, groupID)
	s.Require().NoError(err)
	s.Require().Len(members, 1)

	for _, member := range members {
		s.Require().False(member.IsActive)
	}
}

func (s *KeeperTestSuite) TestSetLastExpiredGroupID() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	k.SetLastExpiredGroupID(ctx, groupID)

	got := k.GetLastExpiredGroupID(ctx)
	s.Require().Equal(groupID, got)
}

func (s *KeeperTestSuite) TestGetSetLastExpiredGroupID() {
	ctx, k := s.ctx, s.app.TSSKeeper

	// Set the last expired group ID
	groupID := tss.GroupID(98765)
	k.SetLastExpiredGroupID(ctx, groupID)

	// Get the last expired group ID
	got := k.GetLastExpiredGroupID(ctx)

	// Assert equality
	s.Require().Equal(groupID, got)
}

func (s *KeeperTestSuite) TestProcessExpiredGroups() {
	ctx, k := s.ctx, s.app.TSSKeeper

	// Create group
	groupID := k.CreateNewGroup(ctx, types.Group{})
	k.SetMember(ctx, types.Member{
		ID:          1,
		GroupID:     groupID,
		Address:     "band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun",
		PubKey:      nil,
		IsMalicious: false,
	})

	// Set the current block height
	blockHeight := int64(30001)
	ctx = ctx.WithBlockHeight(blockHeight)

	// Handle expired groups
	k.HandleExpiredGroups(ctx)

	// Assert that the last expired group ID is updated correctly
	lastExpiredGroupID := k.GetLastExpiredGroupID(ctx)
	s.Require().Equal(groupID, lastExpiredGroupID)
}

func (s *KeeperTestSuite) TestGetSetPendingProcessGroups() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)

	// Set the pending process group in the store
	k.SetPendingProcessGroups(ctx, types.PendingProcessGroups{
		GroupIDs: []tss.GroupID{groupID},
	})

	got := k.GetPendingProcessGroups(ctx)

	// Check if the retrieved pending process groups match the original sample
	s.Require().Len(got, 1)
	s.Require().Equal(groupID, got[0])
}

func (s *KeeperTestSuite) TestHandleProcessGroup() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID, memberID := tss.GroupID(1), tss.MemberID(1)
	member := types.Member{
		ID:          memberID,
		GroupID:     groupID,
		IsMalicious: false,
	}

	k.SetMember(ctx, member)

	k.SetGroup(ctx, types.Group{
		ID:     groupID,
		Status: types.GROUP_STATUS_ROUND_1,
	})
	k.HandleProcessGroup(ctx, groupID)
	group := k.MustGetGroup(ctx, groupID)
	s.Require().Equal(types.GROUP_STATUS_ROUND_2, group.Status)

	k.SetGroup(ctx, types.Group{
		ID:     groupID,
		Status: types.GROUP_STATUS_ROUND_2,
	})
	k.HandleProcessGroup(ctx, groupID)
	group = k.MustGetGroup(ctx, groupID)
	s.Require().Equal(types.GROUP_STATUS_ROUND_3, group.Status)

	k.SetGroup(ctx, types.Group{
		ID:     groupID,
		Status: types.GROUP_STATUS_FALLEN,
	})
	k.HandleProcessGroup(ctx, groupID)
	group = k.MustGetGroup(ctx, groupID)
	s.Require().Equal(types.GROUP_STATUS_FALLEN, group.Status)

	k.SetGroup(ctx, types.Group{
		ID:     groupID,
		Status: types.GROUP_STATUS_ROUND_3,
	})
	k.HandleProcessGroup(ctx, groupID)
	group = k.MustGetGroup(ctx, groupID)
	s.Require().Equal(types.GROUP_STATUS_ACTIVE, group.Status)

	// if member is malicious
	k.SetGroup(ctx, types.Group{
		ID:     groupID,
		Status: types.GROUP_STATUS_ROUND_3,
	})
	member.IsMalicious = true
	k.SetMember(ctx, member)
	k.HandleProcessGroup(ctx, groupID)
	group = k.MustGetGroup(ctx, groupID)
	s.Require().Equal(types.GROUP_STATUS_FALLEN, group.Status)
}

func (s *KeeperTestSuite) TestParams() {
	k := s.app.TSSKeeper

	testCases := []struct {
		name         string
		input        types.Params
		expectErr    bool
		expectErrStr string
	}{
		{
			name: "set invalid params",
			input: types.Params{
				MaxGroupSize:   0,
				MaxDESize:      0,
				CreatingPeriod: 1,
				SigningPeriod:  1,
			},
			expectErr:    true,
			expectErrStr: "must be positive:",
		},
		{
			name: "set full valid params",
			input: types.Params{
				MaxGroupSize:   types.DefaultMaxGroupSize,
				MaxDESize:      types.DefaultMaxDESize,
				CreatingPeriod: types.DefaultCreatingPeriod,
				SigningPeriod:  types.DefaultSigningPeriod,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			expected := k.GetParams(s.ctx)
			err := k.SetParams(s.ctx, tc.input)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expectErrStr)
			} else {
				expected = tc.input
				s.Require().NoError(err)
			}

			p := k.GetParams(s.ctx)
			s.Require().Equal(expected, p)
		})
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
