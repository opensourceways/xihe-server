package message

import(
	common "github.com/opensourceways/xihe-server/common/domain/message"
)

type msgFolloweing struct {
	com      common.MsgNormal
	Follower string `json:"follower"`
}
