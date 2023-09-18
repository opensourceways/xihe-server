package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type WuKongInferenceStartEvent struct {
	Account types.Account
	Desc    WuKongPictureDesc
	Style   string
	EsStyle string
}

type WuKongInferenceErrorEvent struct {
	Account types.Account
	TaskId  uint64
	ErrMsg  string
}

type WuKongAsyncTaskStartEvent struct {
	Account types.Account
	TaskId  uint64
}

type WuKongAsyncInferenceFinishEvent struct {
	Account types.Account
	TaskId  uint64
	Links   map[string]string
}

type BigModelAccessLogEvent struct {
	Account      types.Account
	BigModelType BigmodelType
}

type WuKongPicturePublicizedEvent struct {
	Account types.Account
}

type WuKongPictureLikedEvent struct {
	Account types.Account
}
