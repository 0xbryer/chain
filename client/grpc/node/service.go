package node

import (
	context "context"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterNodeService registers the node gRPC service on the provided gRPC router.
func RegisterNodeService(clientCtx client.Context, server gogogrpc.Server) {
	RegisterServiceServer(server, NewQueryServer(clientCtx))
}

// RegisterGRPCGatewayRoutes mounts the node gRPC service's GRPC-gateway routes
// on the given mux object.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	_ = RegisterServiceHandlerClient(context.Background(), mux, NewServiceClient(clientConn))
}

var _ ServiceServer = queryServer{}

type queryServer struct {
	clientCtx client.Context
}

func NewQueryServer(clientCtx client.Context) ServiceServer {
	return queryServer{
		clientCtx: clientCtx,
	}
}

func (s queryServer) ChainID(ctx context.Context, _ *QueryChainIDRequest) (*QueryChainIDResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return &QueryChainIDResponse{
		ChainID: sdkCtx.ChainID(),
	}, nil
}

func (s queryServer) EVMValidators(
	ctx context.Context,
	_ *QueryEVMValidatorsRequest,
) (*QueryEVMValidatorsResponse, error) {
	node, err := s.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// Get top 100 validators for now
	var page int = 1
	var perPage int = 100
	validators, err := node.Validators(context.Background(), nil, &page, &perPage)
	if err != nil {
		return nil, err
	}

	evmValidatorsResponse := QueryEVMValidatorsResponse{}
	evmValidatorsResponse.BlockHeight = validators.BlockHeight
	evmValidatorsResponse.Validators = []*ValidatorMinimal{}

	for _, validator := range validators.Validators {
		pubKeyBytes, ok := validator.PubKey.(secp256k1.PubKey)
		if !ok {
			return nil, fmt.Errorf("can't get validator public key")
		}

		if pubkey, err := crypto.DecompressPubkey(pubKeyBytes[:]); err != nil {
			return nil, err
		} else {
			evmValidatorsResponse.Validators = append(
				evmValidatorsResponse.Validators,
				&ValidatorMinimal{
					Address:     crypto.PubkeyToAddress(*pubkey).String(),
					VotingPower: validator.VotingPower,
				},
			)
		}
	}

	return &evmValidatorsResponse, nil
}
