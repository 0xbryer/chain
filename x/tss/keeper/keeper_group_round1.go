package keeper

import (
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/chain/v3/pkg/tss"
	"github.com/bandprotocol/chain/v3/x/tss/types"
)

// AddCoefficientCommits adds each coefficient commit into the accumulated commit of its index.
func (k Keeper) AddCoefficientCommits(
	ctx sdk.Context,
	groupID tss.GroupID,
	coefficientCommits tss.Points,
) error {
	for i, commit := range coefficientCommits {
		points := []tss.Point{commit}

		accCommit := k.GetAccumulatedCommit(ctx, groupID, uint64(i))
		if accCommit != nil {
			points = append(points, accCommit)
		}

		total, err := tss.SumPoints(points...)
		if err != nil {
			return types.ErrAddCoeffCommit.Wrapf("failed to sum points: %s", err)
		}

		k.SetAccumulatedCommit(ctx, groupID, uint64(i), total)
	}

	return nil
}

// ValidateRound1Info validates the round1Info of a group member.
func (k Keeper) ValidateRound1Info(
	ctx sdk.Context,
	group types.Group,
	round1Info types.Round1Info,
) error {
	// Check coefficients commit length
	if uint64(len(round1Info.CoefficientCommits)) != group.Threshold {
		return types.ErrInvalidLengthCoeffCommits
	}

	// Get dkg-context
	dkgContext, err := k.GetDKGContext(ctx, group.ID)
	if err != nil {
		return err
	}

	// Verify one time signature
	err = tss.VerifyOneTimeSignature(
		round1Info.MemberID,
		dkgContext,
		round1Info.OneTimeSignature,
		round1Info.OneTimePubKey,
	)
	if err != nil {
		return types.ErrVerifyOneTimeSignatureFailed.Wrapf("failed to verify one time signature: %v", err)
	}

	// Verify A0 signature
	err = tss.VerifyA0Signature(
		round1Info.MemberID,
		dkgContext,
		round1Info.A0Signature,
		round1Info.CoefficientCommits[0],
	)
	if err != nil {
		return types.ErrVerifyA0SignatureFailed.Wrapf("failed to verify A0 signature: %v", err)
	}

	return nil
}

// =====================================
// Round1Info store
// =====================================

// AddRound1Info adds the round1Info of a member in the store and increments the count of round1Info.
func (k Keeper) AddRound1Info(ctx sdk.Context, groupID tss.GroupID, round1Info types.Round1Info) {
	k.SetRound1Info(ctx, groupID, round1Info)

	count := k.GetRound1InfoCount(ctx, groupID)
	k.SetRound1InfoCount(ctx, groupID, count+1)
}

// SetRound1Info sets round1Info for a member of a group.
func (k Keeper) SetRound1Info(ctx sdk.Context, groupID tss.GroupID, round1Info types.Round1Info) {
	ctx.KVStore(k.storeKey).
		Set(types.Round1InfoStoreKey(groupID, round1Info.MemberID), k.cdc.MustMarshal(&round1Info))
}

// HasRound1Info checks if a round1Info of a member exists in the store.
func (k Keeper) HasRound1Info(ctx sdk.Context, groupID tss.GroupID, memberID tss.MemberID) bool {
	return ctx.KVStore(k.storeKey).Has(types.Round1InfoStoreKey(groupID, memberID))
}

// GetRound1Info retrieves round1Info of a member from the store.
func (k Keeper) GetRound1Info(ctx sdk.Context, groupID tss.GroupID, memberID tss.MemberID) (types.Round1Info, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.Round1InfoStoreKey(groupID, memberID))
	if bz == nil {
		return types.Round1Info{}, types.ErrRound1InfoNotFound.Wrapf(
			"failed to get round1Info for groupID: %d and memberID %d",
			groupID,
			memberID,
		)
	}

	var r1 types.Round1Info
	k.cdc.MustUnmarshal(bz, &r1)
	return r1, nil
}

// GetRound1InfoIterator gets an iterator over all round1Info of a group.
func (k Keeper) GetRound1InfoIterator(ctx sdk.Context, groupID tss.GroupID) dbm.Iterator {
	return storetypes.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.Round1InfosStoreKey(groupID))
}

// GetRound1Infos retrieves round1Infos for a group from the store.
func (k Keeper) GetRound1Infos(ctx sdk.Context, groupID tss.GroupID) []types.Round1Info {
	iterator := k.GetRound1InfoIterator(ctx, groupID)
	defer iterator.Close()

	var round1Infos []types.Round1Info
	for ; iterator.Valid(); iterator.Next() {
		var round1Info types.Round1Info
		k.cdc.MustUnmarshal(iterator.Value(), &round1Info)
		round1Infos = append(round1Infos, round1Info)
	}

	return round1Infos
}

// DeleteRound1Infos removes all round1Info associated with a specific group ID from the store.
func (k Keeper) DeleteRound1Infos(ctx sdk.Context, groupID tss.GroupID) {
	iterator := k.GetRound1InfoIterator(ctx, groupID)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		ctx.KVStore(k.storeKey).Delete(key)
	}

	k.DeleteRound1InfoCount(ctx, groupID)
}

// SetRound1InfoCount sets the count of round1Info for a group in the store.
func (k Keeper) SetRound1InfoCount(ctx sdk.Context, groupID tss.GroupID, count uint64) {
	ctx.KVStore(k.storeKey).Set(types.Round1InfoCountStoreKey(groupID), sdk.Uint64ToBigEndian(count))
}

// GetRound1InfoCount retrieves the count of round1Info for a group from the store.
func (k Keeper) GetRound1InfoCount(ctx sdk.Context, groupID tss.GroupID) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(types.Round1InfoCountStoreKey(groupID))
	return sdk.BigEndianToUint64(bz)
}

// DeleteRound1InfoCount remove the round1Info count data of a group from the store.
func (k Keeper) DeleteRound1InfoCount(ctx sdk.Context, groupID tss.GroupID) {
	ctx.KVStore(k.storeKey).Delete(types.Round1InfoCountStoreKey(groupID))
}

// =====================================
// Accumulated commit store
// =====================================

// GetAccumulatedCommitIterator gets an iterator over all accumulated commits of a group.
func (k Keeper) GetAccumulatedCommitIterator(ctx sdk.Context, groupID tss.GroupID) dbm.Iterator {
	return storetypes.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.AccumulatedCommitsStoreKey(groupID))
}

// SetAccumulatedCommit sets accumulated commit for a index of a group.
func (k Keeper) SetAccumulatedCommit(ctx sdk.Context, groupID tss.GroupID, index uint64, commit tss.Point) {
	ctx.KVStore(k.storeKey).Set(types.AccumulatedCommitStoreKey(groupID, index), commit)
}

// GetAccumulatedCommit retrieves accumulated commit of a index of the group from the store.
func (k Keeper) GetAccumulatedCommit(ctx sdk.Context, groupID tss.GroupID, index uint64) tss.Point {
	return ctx.KVStore(k.storeKey).Get(types.AccumulatedCommitStoreKey(groupID, index))
}

// GetAllAccumulatedCommits retrieves all accumulated commits of a group from the store.
func (k Keeper) GetAllAccumulatedCommits(ctx sdk.Context, groupID tss.GroupID) tss.Points {
	iterator := k.GetAccumulatedCommitIterator(ctx, groupID)
	defer iterator.Close()

	var commits tss.Points
	for ; iterator.Valid(); iterator.Next() {
		commits = append(commits, iterator.Value())
	}

	return commits
}

// DeleteAccumulatedCommit removes a accumulated commit of a index of the group from the store.
func (k Keeper) DeleteAccumulatedCommit(ctx sdk.Context, groupID tss.GroupID, index uint64) {
	ctx.KVStore(k.storeKey).Delete(types.AccumulatedCommitStoreKey(groupID, index))
}

// DeleteAccumulatedCommits removes all accumulated commit associated with a specific group ID from the store.
func (k Keeper) DeleteAccumulatedCommits(ctx sdk.Context, groupID tss.GroupID) {
	iterator := k.GetAccumulatedCommitIterator(ctx, groupID)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		ctx.KVStore(k.storeKey).Delete(key)
	}
}
