package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/bandprotocol/chain/v2/pkg/tss"
	"github.com/bandprotocol/chain/v2/x/tss/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	authzKeeper       types.AuthzKeeper
	rollingseedKeeper types.RollingseedKeeper

	router    *types.Router
	hooks     types.TSSHooks
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	authzKeeper types.AuthzKeeper,
	rollingseedKeeper types.RollingseedKeeper,
	rtr *types.Router,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		paramSpace:        paramSpace,
		authzKeeper:       authzKeeper,
		rollingseedKeeper: rollingseedKeeper,
		router:            rtr,
		authority:         authority,
	}
}

// GetAuthority returns the x/tss module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// SetGroupCount sets the number of group count to the given value.
func (k Keeper) SetGroupCount(ctx sdk.Context, count uint64) {
	ctx.KVStore(k.storeKey).Set(types.GroupCountStoreKey, sdk.Uint64ToBigEndian(count))
}

// GetGroupCount returns the current number of all groups ever existed.
func (k Keeper) GetGroupCount(ctx sdk.Context) uint64 {
	return sdk.BigEndianToUint64(ctx.KVStore(k.storeKey).Get(types.GroupCountStoreKey))
}

// GetNextGroupID increments the group count and returns the current number of groups.
func (k Keeper) GetNextGroupID(ctx sdk.Context) tss.GroupID {
	groupNumber := k.GetGroupCount(ctx)
	k.SetGroupCount(ctx, groupNumber+1)
	return tss.GroupID(groupNumber + 1)
}

// CheckIsGrantee checks if the granter granted permissions to the grantee.
func (k Keeper) CheckIsGrantee(ctx sdk.Context, granter sdk.AccAddress, grantee sdk.AccAddress) bool {
	for _, msg := range types.GetTSSGrantMsgTypes() {
		cap, _ := k.authzKeeper.GetAuthorization(
			ctx,
			grantee,
			granter,
			msg,
		)

		if cap == nil {
			return false
		}
	}

	return true
}

// CreateNewGroup creates a new group in the store and returns the id of the group.
func (k Keeper) CreateNewGroup(ctx sdk.Context, group types.Group) tss.GroupID {
	group.ID = k.GetNextGroupID(ctx)
	group.CreatedHeight = uint64(ctx.BlockHeight())
	k.SetGroup(ctx, group)

	return group.ID
}

// GetGroup retrieves a group from the store.
func (k Keeper) GetGroup(ctx sdk.Context, groupID tss.GroupID) (types.Group, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.GroupStoreKey(groupID))
	if bz == nil {
		return types.Group{}, types.ErrGroupNotFound.Wrapf("failed to get group with groupID: %d", groupID)
	}

	group := types.Group{}
	k.cdc.MustUnmarshal(bz, &group)
	return group, nil
}

// MustGetGroup returns the group for the given ID. Panics error if not exists.
func (k Keeper) MustGetGroup(ctx sdk.Context, groupID tss.GroupID) types.Group {
	group, err := k.GetGroup(ctx, groupID)
	if err != nil {
		panic(err)
	}
	return group
}

// SetGroup set a group in the store.
func (k Keeper) SetGroup(ctx sdk.Context, group types.Group) {
	ctx.KVStore(k.storeKey).Set(types.GroupStoreKey(group.ID), k.cdc.MustMarshal(&group))
}

// GetGroupsIterator gets an iterator all group.
func (k Keeper) GetGroupsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.GroupStoreKeyPrefix)
}

// GetGroups retrieves all group of the store.
func (k Keeper) GetGroups(ctx sdk.Context) []types.Group {
	var groups []types.Group
	iterator := k.GetGroupsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var group types.Group
		k.cdc.MustUnmarshal(iterator.Value(), &group)
		groups = append(groups, group)
	}
	return groups
}

// DeleteGroup removes the group from the store.
func (k Keeper) DeleteGroup(ctx sdk.Context, groupID tss.GroupID) {
	ctx.KVStore(k.storeKey).Delete(types.GroupStoreKey(groupID))
}

// SetDKGContext sets DKG context for a group in the store.
func (k Keeper) SetDKGContext(ctx sdk.Context, groupID tss.GroupID, dkgContext []byte) {
	ctx.KVStore(k.storeKey).Set(types.DKGContextStoreKey(groupID), dkgContext)
}

// GetDKGContext retrieves DKG context of a group from the store.
func (k Keeper) GetDKGContext(ctx sdk.Context, groupID tss.GroupID) ([]byte, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.DKGContextStoreKey(groupID))
	if bz == nil {
		return nil, errors.Wrapf(types.ErrDKGContextNotFound, "failed to get dkg-context with groupID: %d", groupID)
	}
	return bz, nil
}

// DeleteDKGContext removes the DKG context data of a group from the store.
func (k Keeper) DeleteDKGContext(ctx sdk.Context, groupID tss.GroupID) {
	ctx.KVStore(k.storeKey).Delete(types.DKGContextStoreKey(groupID))
}

// SetMember sets a member of a group in the store.
func (k Keeper) SetMember(ctx sdk.Context, member types.Member) {
	ctx.KVStore(k.storeKey).Set(types.MemberOfGroupKey(member.GroupID, member.ID), k.cdc.MustMarshal(&member))
}

// SetMembers sets members of a group in the store.
func (k Keeper) SetMembers(ctx sdk.Context, members []types.Member) {
	for _, member := range members {
		k.SetMember(ctx, member)
	}
}

// GetMemberByAddress function retrieves a member of a group from the store by using address.
func (k Keeper) GetMemberByAddress(ctx sdk.Context, groupID tss.GroupID, address string) (types.Member, error) {
	members, err := k.GetGroupMembers(ctx, groupID)
	if err != nil {
		return types.Member{}, err
	}

	for _, member := range members {
		if member.Verify(address) {
			return member, nil
		}
	}

	return types.Member{}, errors.Wrapf(
		types.ErrMemberNotFound,
		"failed to get member with groupID: %d and address: %s",
		groupID,
		address,
	)
}

// GetMember function retrieves a member of a group from the store.
func (k Keeper) GetMember(ctx sdk.Context, groupID tss.GroupID, memberID tss.MemberID) (types.Member, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.MemberOfGroupKey(groupID, memberID))
	if bz == nil {
		return types.Member{}, errors.Wrapf(
			types.ErrMemberNotFound,
			"failed to get member with groupID: %d and memberID: %d",
			groupID,
			memberID,
		)
	}

	member := types.Member{}
	k.cdc.MustUnmarshal(bz, &member)
	return member, nil
}

// MustGetMember returns the member for the given groupID and memberID. Panics error if not exists.
func (k Keeper) MustGetMember(ctx sdk.Context, groupID tss.GroupID, memberID tss.MemberID) types.Member {
	member, err := k.GetMember(ctx, groupID, memberID)
	if err != nil {
		panic(err)
	}
	return member
}

// GetGroupMembersIterator gets an iterator over all members of a group.
func (k Keeper) GetGroupMembersIterator(ctx sdk.Context, groupID tss.GroupID) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.MembersStoreKey(groupID))
}

// GetGroupMembers retrieves all members of a group from the store.
func (k Keeper) GetGroupMembers(ctx sdk.Context, groupID tss.GroupID) ([]types.Member, error) {
	var members []types.Member
	iterator := k.GetGroupMembersIterator(ctx, groupID)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var member types.Member
		k.cdc.MustUnmarshal(iterator.Value(), &member)
		members = append(members, member)
	}
	if len(members) == 0 {
		return nil, errors.Wrapf(types.ErrMemberNotFound, "failed to get members with groupID: %d", groupID)
	}
	return members, nil
}

// GetMembers retrieves all members from store.
func (k Keeper) GetMembers(ctx sdk.Context) []types.Member {
	var members []types.Member
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.MemberStoreKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var member types.Member
		k.cdc.MustUnmarshal(iterator.Value(), &member)
		members = append(members, member)
	}

	return members
}

// DeleteGroupMembers removes all members in the group
func (k Keeper) DeleteGroupMembers(ctx sdk.Context, groupID tss.GroupID) error {
	members, err := k.GetGroupMembers(ctx, groupID)
	if err != nil {
		return err
	}

	for _, member := range members {
		k.DeleteMember(ctx, member)
	}

	return nil
}

// DeleteMember removes a member
func (k Keeper) DeleteMember(ctx sdk.Context, member types.Member) {
	ctx.KVStore(k.storeKey).Delete(types.MemberOfGroupKey(member.GroupID, member.ID))
}

// MustGetMembers retrieves all members of a group from the store. Panics error if not exists.
func (k Keeper) MustGetMembers(ctx sdk.Context, groupID tss.GroupID) []types.Member {
	members, err := k.GetGroupMembers(ctx, groupID)
	if err != nil {
		panic(err)
	}
	return members
}

// GetAvailableMembers retrieves all active members of a group from the store.
func (k Keeper) GetAvailableMembers(ctx sdk.Context, groupID tss.GroupID) ([]types.Member, error) {
	var activeMembers []types.Member
	iterator := k.GetGroupMembersIterator(ctx, groupID)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var member types.Member
		k.cdc.MustUnmarshal(iterator.Value(), &member)

		if member.IsActive {
			activeMembers = append(activeMembers, member)
		}
	}

	// Filter members that have DE left
	filteredMembers, err := k.FilterMembersHaveDE(ctx, activeMembers)
	if err != nil {
		return nil, err
	}

	if len(filteredMembers) == 0 {
		return nil, types.ErrNoActiveMember.Wrapf("no active member in groupID: %d", groupID)
	}
	return filteredMembers, nil
}

// SetLastExpiredGroupID sets the last expired group ID in the store.
func (k Keeper) SetLastExpiredGroupID(ctx sdk.Context, groupID tss.GroupID) {
	ctx.KVStore(k.storeKey).Set(types.LastExpiredGroupIDStoreKey, sdk.Uint64ToBigEndian(uint64(groupID)))
}

// SetMemberIsActive sets a boolean flag represent activeness of the user.
func (k Keeper) SetMemberIsActive(ctx sdk.Context, groupID tss.GroupID, address sdk.AccAddress, status bool) error {
	members := k.MustGetMembers(ctx, groupID)
	for _, m := range members {
		if m.Address == address.String() {
			m.IsActive = status
			k.SetMember(ctx, m)
			return nil
		}
	}

	return types.ErrMemberNotFound.Wrapf(
		"failed to set member active status with groupID: %d and address: %s",
		groupID,
		address,
	)
}

// ActivateMember sets a boolean flag represent activeness of the user to true.
func (k Keeper) ActivateMember(ctx sdk.Context, groupID tss.GroupID, address sdk.AccAddress) error {
	return k.SetMemberIsActive(ctx, groupID, address, true)
}

// DeactivateMember sets a boolean flag represent activeness of the user to false.
func (k Keeper) DeactivateMember(ctx sdk.Context, groupID tss.GroupID, address sdk.AccAddress) error {
	return k.SetMemberIsActive(ctx, groupID, address, false)
}

// GetLastExpiredGroupID retrieves the last expired group ID from the store.
func (k Keeper) GetLastExpiredGroupID(ctx sdk.Context) tss.GroupID {
	bz := ctx.KVStore(k.storeKey).Get(types.LastExpiredGroupIDStoreKey)
	return tss.GroupID(sdk.BigEndianToUint64(bz))
}

// HandleExpiredGroups cleans up expired groups and removes them from the store.
func (k Keeper) HandleExpiredGroups(ctx sdk.Context) {
	// Get the current group ID to start processing from
	currentGroupID := k.GetLastExpiredGroupID(ctx) + 1

	// Get the last group ID in the store
	lastGroupID := tss.GroupID(k.GetGroupCount(ctx))

	// Get the group signature creating period
	creatingPeriod := k.GetParams(ctx).CreatingPeriod

	// Process each group starting from currentGroupID
	for ; currentGroupID <= lastGroupID; currentGroupID++ {
		// Get the group
		group := k.MustGetGroup(ctx, currentGroupID)

		// Check if the group is still within the expiration period
		if group.CreatedHeight+creatingPeriod > uint64(ctx.BlockHeight()) {
			break
		}

		// Check group is not active
		if group.Status != types.GROUP_STATUS_ACTIVE && group.Status != types.GROUP_STATUS_FALLEN {
			// Handle the hooks before setting group to be expired.
			// this shouldn't return any error.
			if err := k.Hooks().BeforeSetGroupExpired(ctx, group); err != nil {
				panic(err)
			}

			// Update group status
			group.Status = types.GROUP_STATUS_EXPIRED
			k.SetGroup(ctx, group)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeExpiredGroup,
					sdk.NewAttribute(types.AttributeKeyGroupID, fmt.Sprintf("%d", group.ID)),
				),
			)
		}

		// Cleanup all interim data associated with the group
		k.DeleteAllDKGInterimData(ctx, currentGroupID)

		// Set the last expired group ID to the current group ID
		k.SetLastExpiredGroupID(ctx, currentGroupID)
	}
}

// AddPendingProcessGroup adds a new pending process group to the store.
func (k Keeper) AddPendingProcessGroup(ctx sdk.Context, groupID tss.GroupID) {
	pgs := k.GetPendingProcessGroups(ctx)
	pgs = append(pgs, groupID)
	k.SetPendingProcessGroups(ctx, types.PendingProcessGroups{
		GroupIDs: pgs,
	})
}

// SetPendingProcessGroups sets the given pending process groups in the store.
func (k Keeper) SetPendingProcessGroups(ctx sdk.Context, pgs types.PendingProcessGroups) {
	ctx.KVStore(k.storeKey).Set(types.PendingProcessGroupsStoreKey, k.cdc.MustMarshal(&pgs))
}

// GetPendingProcessGroups retrieves the list of pending process groups from the store.
// It returns an empty list if the key does not exist in the store.
func (k Keeper) GetPendingProcessGroups(ctx sdk.Context) []tss.GroupID {
	bz := ctx.KVStore(k.storeKey).Get(types.PendingProcessGroupsStoreKey)
	if len(bz) == 0 {
		// Return an empty list if the key does not exist in the store.
		return []tss.GroupID{}
	}
	pgs := types.PendingProcessGroups{}
	k.cdc.MustUnmarshal(bz, &pgs)
	return pgs.GroupIDs
}

// HandleProcessGroup handles the pending process group based on its status.
// It updates the group status and emits appropriate events.
func (k Keeper) HandleProcessGroup(ctx sdk.Context, groupID tss.GroupID) {
	group := k.MustGetGroup(ctx, groupID)
	switch group.Status {
	case types.GROUP_STATUS_ROUND_1:
		group.Status = types.GROUP_STATUS_ROUND_2
		group.PubKey = k.GetAccumulatedCommit(ctx, groupID, 0)
		k.SetGroup(ctx, group)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRound1Success,
				sdk.NewAttribute(types.AttributeKeyGroupID, fmt.Sprintf("%d", groupID)),
				sdk.NewAttribute(types.AttributeKeyStatus, group.Status.String()),
			),
		)
	case types.GROUP_STATUS_ROUND_2:
		group.Status = types.GROUP_STATUS_ROUND_3
		k.SetGroup(ctx, group)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRound2Success,
				sdk.NewAttribute(types.AttributeKeyGroupID, fmt.Sprintf("%d", groupID)),
				sdk.NewAttribute(types.AttributeKeyStatus, group.Status.String()),
			),
		)
	case types.GROUP_STATUS_FALLEN:
		group.Status = types.GROUP_STATUS_FALLEN
		k.SetGroup(ctx, group)

		// Handle the hooks when group creation is fallen; this shouldn't return any error.
		if err := k.Hooks().AfterCreatingGroupFailed(ctx, group); err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRound3Failed,
				sdk.NewAttribute(types.AttributeKeyGroupID, fmt.Sprintf("%d", groupID)),
				sdk.NewAttribute(types.AttributeKeyStatus, group.Status.String()),
			),
		)
	case types.GROUP_STATUS_ROUND_3:
		// Get members to check malicious
		members := k.MustGetMembers(ctx, group.ID)
		if !types.Members(members).HaveMalicious() {
			group.Status = types.GROUP_STATUS_ACTIVE
			k.SetGroup(ctx, group)

			// Handle the hooks when group is ready. this shouldn't return any error.
			if err := k.Hooks().AfterCreatingGroupCompleted(ctx, group); err != nil {
				panic(err)
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeRound3Success,
					sdk.NewAttribute(types.AttributeKeyGroupID, fmt.Sprintf("%d", groupID)),
					sdk.NewAttribute(types.AttributeKeyStatus, group.Status.String()),
				),
			)
		} else {
			group.Status = types.GROUP_STATUS_FALLEN
			k.SetGroup(ctx, group)

			// Handle the hooks when group creation is fallen; this shouldn't return any error.
			if err := k.Hooks().AfterCreatingGroupFailed(ctx, group); err != nil {
				panic(err)
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeRound3Failed,
					sdk.NewAttribute(types.AttributeKeyGroupID, fmt.Sprintf("%d", groupID)),
					sdk.NewAttribute(types.AttributeKeyStatus, group.Status.String()),
				),
			)
		}
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Hooks gets the hooks for tss *Keeper {
func (k *Keeper) Hooks() types.TSSHooks {
	if k.hooks == nil {
		return types.MultiTSSHooks{}
	}

	return k.hooks
}

// SetHooks Set the hooks for the tss keeper.
func (k *Keeper) SetHooks(sh types.TSSHooks) {
	if k.hooks != nil {
		panic("cannot set hooks twice")
	}

	k.hooks = sh
}
