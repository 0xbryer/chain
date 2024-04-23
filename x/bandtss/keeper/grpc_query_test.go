package keeper_test

import (
	"fmt"

	"github.com/bandprotocol/chain/v2/x/bandtss/types"
	tsstypes "github.com/bandprotocol/chain/v2/x/tss/types"
)

func (s *KeeperTestSuite) TestGRPCQueryMembers() {
	ctx, q := s.ctx, s.queryClient

	s.SetupGroup(tsstypes.GROUP_STATUS_ACTIVE)

	var req types.QueryMembersRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		postTest func(res *types.QueryMembersResponse, err error)
	}{
		{
			"success",
			func() {
				req = types.QueryMembersRequest{
					IsActive: true,
				}
			},
			true,
			func(res *types.QueryMembersResponse, err error) {
				s.Require().NoError(err)
				s.Require().NotNil(res)
				s.Require().Len(res.Members, 2)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()

			res, err := q.Members(ctx, &req)
			if tc.expPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}

			tc.postTest(res, err)
		})
	}
}
