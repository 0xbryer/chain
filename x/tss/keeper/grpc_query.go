package keeper

import (
	"context"

	"github.com/bandprotocol/chain/v2/x/tss/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) Group(goCtx context.Context, req *types.QueryGroupRequest) (*types.QueryGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	group, found := k.GetGroup(ctx, req.GroupId)
	if !found {
		return &types.QueryGroupResponse{}, sdkerrors.Wrapf(types.ErrGroupNotFound, "groupID: %d")
	}

	return &types.QueryGroupResponse{
		Group: &group,
	}, nil
}

func (k Querier) Members(goCtx context.Context, req *types.QueryMembersRequest) (*types.QueryMembersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	members, found := k.GetMembers(ctx, req.GroupId)
	if !found {
		return &types.QueryMembersResponse{}, sdkerrors.Wrapf(types.ErrMemberNotFound, "groupID: %d")
	}

	return &types.QueryMembersResponse{
		Members: members,
	}, nil
}

func (k Querier) IsGrantee(
	goCtx context.Context,
	req *types.QueryIsGranteeRequest,
) (*types.QueryIsGranteeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	granter, err := sdk.AccAddressFromBech32(req.GranterAddress)
	if err != nil {
		return &types.QueryIsGranteeResponse{}, err
	}
	grantee, err := sdk.AccAddressFromBech32(req.GranteeAddress)
	if err != nil {
		return &types.QueryIsGranteeResponse{}, err
	}

	return &types.QueryIsGranteeResponse{
		IsGrantee: k.Keeper.IsGrantee(ctx, granter, grantee),
	}, nil
}
