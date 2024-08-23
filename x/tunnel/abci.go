package tunnel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/chain/v2/x/tunnel/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
}

func EndBlocker(ctx sdk.Context, k *keeper.Keeper) {
	// Generate packets to be sent
	packets := k.GeneratePackets(ctx)
	for _, packet := range packets {
		// Send packet to the destination route and store the packet
		k.HandlePacket(ctx, packet)
	}
}
