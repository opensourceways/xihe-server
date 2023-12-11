package domain

import types "github.com/opensourceways/xihe-server/domain"

type UserPoints struct {
	User    types.Account
	Total   int
	Items   []Item
	Version int
}

type Item struct {
	TaskId   string
	TaskName Sentence
	Descs    Sentence
	Date     int64
	Points   int
}
