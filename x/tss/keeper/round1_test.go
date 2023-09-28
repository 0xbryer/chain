package keeper_test

import (
	"github.com/bandprotocol/chain/v2/pkg/tss"
	"github.com/bandprotocol/chain/v2/x/tss/types"
)

func (s *KeeperTestSuite) TestGetSetRound1Info() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	memberID := tss.MemberID(1)
	round1Info := types.Round1Info{
		MemberID: memberID,
		CoefficientCommits: tss.Points{
			[]byte("point1"),
			[]byte("point2"),
		},
		OneTimePubKey:    []byte("OneTimePubKeySimple"),
		A0Signature:      []byte("A0SignatureSimple"),
		OneTimeSignature: []byte("OneTimeSignatureSimple"),
	}

	// Set round 1 info
	k.SetRound1Info(ctx, groupID, round1Info)

	got, err := k.GetRound1Info(ctx, groupID, memberID)
	s.Require().NoError(err)
	s.Require().Equal(round1Info, got)
}

func (s *KeeperTestSuite) TestAddRound1Info() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	memberID := tss.MemberID(1)
	round1Info := types.Round1Info{
		MemberID: memberID,
		CoefficientCommits: tss.Points{
			[]byte("point1"),
			[]byte("point2"),
		},
		OneTimePubKey:    []byte("OneTimePubKeySimple"),
		A0Signature:      []byte("A0SignatureSimple"),
		OneTimeSignature: []byte("OneTimeSignatureSimple"),
	}

	// Add round 1 info
	k.AddRound1Info(ctx, groupID, round1Info)

	gotR1, err := k.GetRound1Info(ctx, groupID, memberID)
	s.Require().NoError(err)
	s.Require().Equal(round1Info, gotR1)
	gotR1Count := k.GetRound1InfoCount(ctx, groupID)
	s.Require().Equal(uint64(1), gotR1Count)
}

func (s *KeeperTestSuite) TestDeleteRound1Info() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	memberID := tss.MemberID(1)
	round1Info := types.Round1Info{
		MemberID: memberID,
		CoefficientCommits: tss.Points{
			[]byte("point1"),
			[]byte("point2"),
		},
		OneTimePubKey:    []byte("OneTimePubKeySimple"),
		A0Signature:      []byte("A0SignatureSimple"),
		OneTimeSignature: []byte("OneTimeSignatureSimple"),
	}

	k.SetRound1Info(ctx, groupID, round1Info)

	got, err := k.GetRound1Info(ctx, groupID, memberID)
	s.Require().NoError(err)
	s.Require().Equal(round1Info, got)

	k.DeleteRound1Info(ctx, groupID, memberID)

	_, err = k.GetRound1Info(ctx, groupID, memberID)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestDeleteRound1Infos() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	memberID := tss.MemberID(1)
	round1Info := types.Round1Info{
		MemberID: memberID,
		CoefficientCommits: tss.Points{
			[]byte("point1"),
			[]byte("point2"),
		},
		OneTimePubKey:    []byte("OneTimePubKeySimple"),
		A0Signature:      []byte("A0SignatureSimple"),
		OneTimeSignature: []byte("OneTimeSignatureSimple"),
	}

	k.SetRound1Info(ctx, groupID, round1Info)

	got, err := k.GetRound1Info(ctx, groupID, memberID)
	s.Require().NoError(err)
	s.Require().Equal(round1Info, got)

	k.DeleteRound1Infos(ctx, groupID)

	_, err = k.GetRound1Info(ctx, groupID, memberID)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestGetRound1Infos() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	member1 := tss.MemberID(1)
	member2 := tss.MemberID(2)
	round1InfoMember1 := types.Round1Info{
		MemberID: member1,
		CoefficientCommits: tss.Points{
			[]byte("point1"),
			[]byte("point2"),
		},
		OneTimePubKey:    []byte("OneTimePubKeySimple"),
		A0Signature:      []byte("A0SignatureSimple"),
		OneTimeSignature: []byte("OneTimeSignatureSimple"),
	}
	round1InfoMember2 := types.Round1Info{
		MemberID: member2,
		CoefficientCommits: tss.Points{
			[]byte("point1"),
			[]byte("point2"),
		},
		OneTimePubKey:    []byte("OneTimePubKeySimple"),
		A0Signature:      []byte("A0SignatureSimple"),
		OneTimeSignature: []byte("OneTimeSignatureSimple"),
	}

	// Set round 1 infos
	k.AddRound1Info(ctx, groupID, round1InfoMember1)
	k.AddRound1Info(ctx, groupID, round1InfoMember2)

	got := k.GetRound1Infos(ctx, groupID)

	s.Require().Equal([]types.Round1Info{round1InfoMember1, round1InfoMember2}, got)
}

func (s *KeeperTestSuite) TestGetSetRound1InfoCount() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	count := uint64(5)

	// Set round 1 info count
	k.SetRound1InfoCount(ctx, groupID, count)

	got := k.GetRound1InfoCount(ctx, groupID)
	s.Require().Equal(uint64(5), got)
}

func (s *KeeperTestSuite) TestDeleteRound1InfoCount() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	count := uint64(5)

	// Set round 1 info count
	k.SetRound1InfoCount(ctx, groupID, count)

	// Delete round 1 info count
	k.DeleteRound1InfoCount(ctx, groupID)

	got := k.GetRound1InfoCount(ctx, groupID)
	s.Require().Empty(got)
}

func (s *KeeperTestSuite) TestGetSetAccumulatedCommit() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	index := uint64(1)
	commit := tss.Point([]byte("point"))

	// Set Accumulated Commit
	k.SetAccumulatedCommit(ctx, groupID, index, commit)

	// Get Accumulated Commit
	got := k.GetAccumulatedCommit(ctx, groupID, index)

	s.Require().Equal(commit, got)
}

func (s *KeeperTestSuite) TestDeleteAccumulatedCommit() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	index := uint64(1)
	commit := tss.Point([]byte("point"))

	// Set Accumulated Commit
	k.SetAccumulatedCommit(ctx, groupID, index, commit)

	// Delete Accumulated Commit
	k.DeleteAccumulatedCommit(ctx, groupID, index)

	// Get Accumulated Commit
	got := k.GetAccumulatedCommit(ctx, groupID, index)
	s.Require().Equal(tss.Point(nil), got)
}

func (s *KeeperTestSuite) TestDeleteAccumulatedCommits() {
	ctx, k := s.ctx, s.app.TSSKeeper
	groupID := tss.GroupID(1)
	index := uint64(1)
	commit := tss.Point([]byte("point"))

	// Set Accumulated Commit
	k.SetAccumulatedCommit(ctx, groupID, index, commit)

	// Delete Accumulated Commits
	k.DeleteAccumulatedCommits(ctx, groupID)

	// Get Accumulated Commit
	got := k.GetAccumulatedCommit(ctx, groupID, index)
	s.Require().Equal(tss.Point(nil), got)
}
