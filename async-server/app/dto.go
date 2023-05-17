package app

import (
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type WuKongCmd struct {
	User  types.Account
	Desc  bigmodeldomain.WuKongPictureDesc
	Style string
}
