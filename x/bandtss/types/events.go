package types

const (
	EventTypeActivate               = "activate"
	EventTypeHeartbeat              = "heartbeat"
	EventTypeInactiveStatus         = "inactive_status"
	EventTypeGroupTransition        = "group_transition"
	EventTypeGroupTransitionSuccess = "group_transition_success"
	EventTypeGroupTransitionFailed  = "group_transition_failed"
	EventTypeSigningRequestCreated  = "bandtss_signing_request_created"

	AttributeKeyAddress                = "address"
	AttributeKeySigningID              = "bandtss_signing_id"
	AttributeKeyCurrentGroupID         = "current_group_id"
	AttributeKeyIncomingGroupID        = "incoming_group_id"
	AttributeKeyCurrentGroupSigningID  = "current_group_signing_id"
	AttributeKeyIncomingGroupSigningID = "incoming_group_signing_id"
	AttributeKeyTransitionStatus       = "transition_status"
	AttributeKeyExecTime               = "exec_time"
	AttributeKeyGroupID                = "group_id"
	AttributeKeyIncomingGroupPubKey    = "incoming_group_pub_key"
	AttributeKeyCurrentGroupPubKey     = "current_group_pub_key"
	AttributeKeyRAddress               = "r_address"
	AttributeKeySignature              = "signature"
)
