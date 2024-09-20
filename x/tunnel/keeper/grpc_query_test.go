package keeper_test

import (
	"github.com/bandprotocol/chain/v2/x/tunnel/types"
)

func (s *KeeperTestSuite) TestGRPCQueryTunnels() {
	ctx, k, q := s.ctx, s.keeper, s.queryServer

	tunnel1 := types.Tunnel{
		ID: 1,
	}
	tunnel2 := types.Tunnel{
		ID: 2,
	}
	k.SetTunnel(ctx, tunnel1)
	k.SetTunnel(ctx, tunnel2)

	resp, err := q.Tunnels(ctx, &types.QueryTunnelsRequest{})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.Tunnels, 2)
	s.Require().Equal(tunnel1, *resp.Tunnels[0])
	s.Require().Equal(tunnel2, *resp.Tunnels[1])
}

func (s *KeeperTestSuite) TestGRPCQueryTunnel() {
	ctx, k, q := s.ctx, s.keeper, s.queryServer

	tunnel := types.Tunnel{
		ID: 1,
	}
	k.SetTunnel(ctx, tunnel)

	resp, err := q.Tunnel(ctx, &types.QueryTunnelRequest{
		TunnelId: 1,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(tunnel, resp.Tunnel)
}

func (s *KeeperTestSuite) TestGRPCQueryPackets() {
	ctx, k, q := s.ctx, s.keeper, s.queryServer

	tunnel := types.Tunnel{
		ID:         1,
		NonceCount: 2,
	}
	r := types.TSSRoute{
		DestinationChainID:         "1",
		DestinationContractAddress: "0x123",
	}
	err := tunnel.SetRoute(&r)
	s.Require().NoError(err)

	k.SetTunnel(ctx, tunnel)

	packet1 := types.Packet{
		TunnelID: 1,
		Nonce:    1,
	}
	packet2 := types.Packet{
		TunnelID: 1,
		Nonce:    2,
	}
	err = packet1.SetPacketContent(&types.TSSPacketContent{
		SigningID:                  1,
		DestinationChainID:         r.DestinationChainID,
		DestinationContractAddress: r.DestinationContractAddress,
	})
	s.Require().NoError(err)
	err = packet2.SetPacketContent(&types.TSSPacketContent{
		SigningID:                  2,
		DestinationChainID:         r.DestinationChainID,
		DestinationContractAddress: r.DestinationContractAddress,
	})
	s.Require().NoError(err)
	k.SetPacket(ctx, packet1)
	k.SetPacket(ctx, packet2)

	resp, err := q.Packets(ctx, &types.QueryPacketsRequest{
		TunnelId: 1,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.Packets, 2)
	s.Require().Equal(packet1, *resp.Packets[0])
	s.Require().Equal(packet2, *resp.Packets[1])
}

func (s *KeeperTestSuite) TestGRPCQueryPacket() {
	ctx, k, q := s.ctx, s.keeper, s.queryServer

	// set tunnel
	tunnel := types.Tunnel{
		ID:         1,
		NonceCount: 2,
	}
	r := types.TSSRoute{
		DestinationChainID:         "1",
		DestinationContractAddress: "0x123",
	}
	err := tunnel.SetRoute(&r)
	s.Require().NoError(err)
	k.SetTunnel(ctx, tunnel)

	packet1 := types.Packet{
		TunnelID: 1,
		Nonce:    1,
	}
	err = packet1.SetPacketContent(&types.TSSPacketContent{
		SigningID:                  1,
		DestinationChainID:         r.DestinationChainID,
		DestinationContractAddress: r.DestinationContractAddress,
	})
	s.Require().NoError(err)
	k.SetPacket(ctx, packet1)

	res, err := q.Packet(ctx, &types.QueryPacketRequest{
		TunnelId: 1,
		Nonce:    1,
	})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(packet1, *res.Packet)
}
