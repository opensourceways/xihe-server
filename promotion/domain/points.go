package domain

import types "github.com/opensourceways/xihe-server/domain"

type UserPoints struct {
	User  types.Account
	Total int
	Items []Item
}

type Item struct {
	Id       string
	TaskName Sentence
	Descs    Sentence
	Date     string
	Points   int
}


