package sender

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	bandtsstypes "github.com/bandprotocol/chain/v3/x/bandtss/types"
	"github.com/bandprotocol/chain/v3/x/tss/types"
)

// GetMsgDetail represents the detail string of a message for logging.
func GetMsgDetail(msg sdk.Msg) (detail string) {
	switch t := msg.(type) {
	case *types.MsgSubmitDKGRound1:
		detail = fmt.Sprintf("Type: %s, GroupID: %d", sdk.MsgTypeURL(t), t.GroupID)
	case *types.MsgSubmitDKGRound2:
		detail = fmt.Sprintf("Type: %s, GroupID: %d", sdk.MsgTypeURL(t), t.GroupID)
	case *types.MsgConfirm:
		detail = fmt.Sprintf("Type: %s, GroupID: %d", sdk.MsgTypeURL(t), t.GroupID)
	case *types.MsgComplain:
		detail = fmt.Sprintf("Type: %s, GroupID: %d", sdk.MsgTypeURL(t), t.GroupID)
	case *types.MsgSubmitDEs:
		detail = fmt.Sprintf("Type: %s", sdk.MsgTypeURL(t))
	case *types.MsgSubmitSignature:
		detail = fmt.Sprintf("Type: %s, SigningID: %d", sdk.MsgTypeURL(t), t.SigningID)
	case *bandtsstypes.MsgHeartbeat:
		detail = fmt.Sprintf("Type: %s", sdk.MsgTypeURL(t))
	default:
		detail = "Type: Unknown"
	}

	return detail
}

// GetMsgDetails extracts the detail from SDK messages.
func GetMsgDetails(msgs ...sdk.Msg) (details []string) {
	for _, msg := range msgs {
		details = append(details, GetMsgDetail(msg))
	}

	return details
}
