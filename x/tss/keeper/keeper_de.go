package keeper

import (
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bandprotocol/chain/v3/x/tss/types"
)

// EnqueueDEs insert multiple DE objects for a given address in the context's KVStore
// them into the queue. It returns an error if the DE size exceeds the maximum limit.
func (k Keeper) EnqueueDEs(ctx sdk.Context, address sdk.AccAddress, des []types.DE) error {
	deQueue := k.GetDEQueue(ctx, address)
	cnt := deQueue.Tail - deQueue.Head
	if cnt+uint64(len(des)) > k.GetParams(ctx).MaxDESize {
		return types.ErrDEReachMaxLimit.Wrapf("DE size exceeds %d", k.GetParams(ctx).MaxDESize)
	}

	for i, de := range des {
		k.SetDE(ctx, address, deQueue.Tail+uint64(i), de)
	}

	deQueue.Tail += uint64(len(des))
	k.SetDEQueue(ctx, address, deQueue)
	return nil
}

// DequeueDE retrieves a DE object from the context's KVStore for a given address and remove
// from the queue. Returns an error if no DE in the queue.
func (k Keeper) DequeueDE(ctx sdk.Context, address sdk.AccAddress) (types.DE, error) {
	deQueue := k.GetDEQueue(ctx, address)
	if deQueue.Head >= deQueue.Tail {
		return types.DE{}, types.ErrDENotFound.Wrapf("DE not found for address %s", address)
	}

	de, err := k.GetDE(ctx, address, deQueue.Head)
	if err != nil {
		return types.DE{}, err
	}
	k.DeleteDE(ctx, address, deQueue.Head)

	deQueue.Head += 1
	k.SetDEQueue(ctx, address, deQueue)
	return de, nil
}

// DequeueDEs dequeues DEs from the selected members. It takes a list of member IDs (mids)
// and members information (members) and returns the list of selected DEs ordered by selected members.
func (k Keeper) DequeueDEs(ctx sdk.Context, members []types.Member) ([]types.DE, error) {
	des := make([]types.DE, 0, len(members))
	for _, member := range members {
		accMember, err := sdk.AccAddressFromBech32(member.Address)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid account address: %s", err)
		}

		de, err := k.DequeueDE(ctx, accMember)
		if err != nil {
			return nil, err
		}
		des = append(des, de)
	}

	return des, nil
}

// ====================================
// DE Queue store
// ====================================

// SetDEQueue sets the DE's queue information of a given address.
func (k Keeper) SetDEQueue(ctx sdk.Context, address sdk.AccAddress, deQueue types.DEQueue) {
	ctx.KVStore(k.storeKey).Set(types.DEQueueStoreKey(address), k.cdc.MustMarshal(&deQueue))
}

// GetDEQueue retrieves the DE's queue information of a given address.
func (k Keeper) GetDEQueue(ctx sdk.Context, address sdk.AccAddress) types.DEQueue {
	bz := ctx.KVStore(k.storeKey).Get(types.DEQueueStoreKey(address))
	if bz == nil {
		return types.NewDEQueue(0, 0)
	}

	var deQueue types.DEQueue
	k.cdc.MustUnmarshal(bz, &deQueue)
	return deQueue
}

// ====================================
// DE store
// ====================================

// SetDE sets a DE object in the context's KVStore for a given address at the given index.
func (k Keeper) SetDE(ctx sdk.Context, address sdk.AccAddress, index uint64, de types.DE) {
	ctx.KVStore(k.storeKey).Set(types.DEStoreKey(address, index), k.cdc.MustMarshal(&de))
}

// GetDE retrieves the DE's of a given address at the given index.
func (k Keeper) GetDE(ctx sdk.Context, address sdk.AccAddress, index uint64) (types.DE, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.DEStoreKey(address, index))
	if bz == nil {
		return types.DE{}, types.ErrDENotFound.Wrapf("DE not found for address %s", address)
	}

	var de types.DE
	k.cdc.MustUnmarshal(bz, &de)
	return de, nil
}

// HasDE function checks if the DE exists in the store.
func (k Keeper) HasDE(ctx sdk.Context, address sdk.AccAddress) bool {
	deQueue := k.GetDEQueue(ctx, address)
	return deQueue.Tail > deQueue.Head
}

func (k Keeper) DeleteDE(ctx sdk.Context, address sdk.AccAddress, index uint64) {
	ctx.KVStore(k.storeKey).Delete(types.DEStoreKey(address, index))
}

// GetDEQueueIterator function gets an iterator over all de queue from the context's KVStore
func (k Keeper) GetDEQueueIterator(ctx sdk.Context) dbm.Iterator {
	return storetypes.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.DEQueueStoreKeyPrefix)
}
