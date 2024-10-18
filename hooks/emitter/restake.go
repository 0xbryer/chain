package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/chain/v3/hooks/common"
	"github.com/bandprotocol/chain/v3/x/restake/types"
)

func (h *Hook) updateRestakeStake(ctx sdk.Context, stakerAddr string) {
	addr := sdk.MustAccAddressFromBech32(stakerAddr)
	stake := h.restakeKeeper.GetStake(ctx, addr)
	h.Write("SET_RESTAKE_STAKE", common.JsDict{
		"staker": stakerAddr,
		"coins":  stake.Coins.String(),
	})
}

func (h *Hook) updateRestakeVault(ctx sdk.Context, key string) {
	vault, _ := h.restakeKeeper.GetVault(ctx, key)
	h.Write("SET_RESTAKE_VAULT", common.JsDict{
		"key":       vault.Key,
		"is_active": vault.IsActive,
	})
}

func (h *Hook) updateRestakeLock(ctx sdk.Context, stakerAddr string, key string) {
	addr := sdk.MustAccAddressFromBech32(stakerAddr)
	lock, found := h.restakeKeeper.GetLock(ctx, addr, key)
	if !found {
		h.Write("REMOVE_RESTAKE_LOCK", common.JsDict{
			"staker_address": addr,
			"key":            key,
		})

		return
	}

	h.Write("SET_RESTAKE_LOCK", common.JsDict{
		"staker_address": addr,
		"key":            key,
		"power":          lock.Power.String(),
	})
}

// handleRestakeEventCreateVault implements emitter handler for EventCreateVault.
func (h *Hook) handleRestakeEventCreateVault(ctx sdk.Context, evMap common.EvMap) {
	h.updateRestakeVault(ctx, evMap[types.EventTypeCreateVault+"."+types.AttributeKeyKey][0])
}

// handleRestakeEventDeactivateVault implements emitter handler for EventDeactivateVault.
func (h *Hook) handleRestakeEventDeactivateVault(ctx sdk.Context, evMap common.EvMap) {
	h.updateRestakeVault(ctx, evMap[types.EventTypeDeactivateVault+"."+types.AttributeKeyKey][0])
}

// handleRestakeEventLockPower implements emitter handler for EventLockPower.
func (h *Hook) handleRestakeEventLockPower(ctx sdk.Context, evMap common.EvMap) {
	stakerAddr := evMap[types.EventTypeLockPower+"."+types.AttributeKeyStaker][0]
	key := evMap[types.EventTypeLockPower+"."+types.AttributeKeyKey][0]
	h.updateRestakeLock(ctx, stakerAddr, key)
}

// handleRestakeEventStake implements emitter handler for EventStake.
func (h *Hook) handleRestakeEventStake(ctx sdk.Context, evMap common.EvMap) {
	h.updateRestakeStake(ctx, evMap[types.EventTypeStake+"."+types.AttributeKeyStaker][0])
}

// handleRestakeEventUnstake implements emitter handler for EventUnstake.
func (h *Hook) handleRestakeEventUnstake(ctx sdk.Context, evMap common.EvMap) {
	h.updateRestakeStake(ctx, evMap[types.EventTypeUnstake+"."+types.AttributeKeyStaker][0])
}
