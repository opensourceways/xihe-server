/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics defines the topic names for message adapter operations.
type Topics struct {
	ComputilityUserJoined  string `json:"org_user_joined" required:"true"`
	ComputilityUserRemoved string `json:"org_user_removed" required:"true"`
	ComputilityOrgDeleted  string `json:"org_deleted" required:"true"`
	OrgInvite              string `json:"org_invite" required:"true"`
	OrgRequest             string `json:"org_request" required:"true"`
	OrgRequestReject       string `json:"org_request_reject" required:"true"`
}
