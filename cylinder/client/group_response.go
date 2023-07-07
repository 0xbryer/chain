package client

import (
	"fmt"

	"github.com/bandprotocol/chain/v2/pkg/tss"
	"github.com/bandprotocol/chain/v2/x/tss/types"
)

// GroupResponse wraps the types.QueryGroupResponse to provide additional helper methods.
type GroupResponse struct {
	types.QueryGroupResponse
}

// NewGroupResponse creates a new instance of GroupResponse.
func NewGroupResponse(gr *types.QueryGroupResponse) *GroupResponse {
	return &GroupResponse{*gr}
}

// GetRound1Info retrieves the Round1Commitment for the specified member ID.
func (gr *GroupResponse) GetRound1Info(mid tss.MemberID) (types.Round1Info, error) {
	for _, data := range gr.Round1Infos {
		if data.MemberID == mid {
			return data, nil
		}
	}

	return types.Round1Info{}, fmt.Errorf("No Round1Info from MemberID(%d)", mid)
}

// GetRound2Info retrieves the Round1Commitment for the specified member ID.
func (gr *GroupResponse) GetRound2Info(mid tss.MemberID) (types.Round2Info, error) {
	for _, data := range gr.Round2Infos {
		if data.MemberID == mid {
			return data, nil
		}
	}

	return types.Round2Info{}, fmt.Errorf("No Round2Info from MemberID(%d)", mid)
}

// GetEncryptedSecretShare retrieves the encrypted secret share from member (Sender) to member (Receiver)
func (gr *GroupResponse) GetEncryptedSecretShare(senderID, receiverID tss.MemberID) (tss.Scalar, error) {
	r2Sender, err := gr.GetRound2Info(senderID)
	if err != nil {
		return nil, err
	}

	// Determine which slot of encrypted secret shares is for Receiver
	slot := types.FindMemberSlot(senderID, receiverID)

	if int(slot) >= len(r2Sender.EncryptedSecretShares) {
		return nil, fmt.Errorf("No encrypted secret share from MemberID(%d) to MemberID(%d)", senderID, receiverID)
	}

	return r2Sender.EncryptedSecretShares[slot], nil
}

// IsMember returns boolean to show if the address is the member in the group
func (gr *GroupResponse) IsMember(address string) bool {
	for _, member := range gr.Members {
		if member.Address == address {
			return true
		}
	}

	return false
}
