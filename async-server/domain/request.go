package domain

import (
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type WuKongRequest struct {
	User      types.Account
	TaskType  TaskType
	Style     string
	Desc      bigmodeldomain.WuKongPictureDesc
	CreatedAt commondomain.Time
}
