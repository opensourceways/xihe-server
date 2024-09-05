/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package domain provides domain organization and configuration for the app service.
package domain

import (
	"encoding/json"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
)

const (
	userJoinTypeRequest = "request"
	userJoinTypeInvite  = "invite"
)

// userJoinEvent
type userJoinEvent struct {
	OrgName  string `json:"org_name"`
	UserName string `json:"user_name"`
	Type     string `json:"type"`
}

// Message returns the JSON representation of the userJoinEvent.
func (e *userJoinEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewUserJoinEventByRequest creates a new userJoinEvent instance with the given Organization.
func NewUserJoinEventByRequest(org primitive.Account, user primitive.Account) userJoinEvent {
	return userJoinEvent{
		OrgName:  org.Account(),
		UserName: user.Account(),
		Type:     userJoinTypeRequest,
	}
}

// NewUserJoinEventByInvite creates a new userJoinEvent instance with the given Organization.
func NewUserJoinEventByInvite(org primitive.Account, user primitive.Account) userJoinEvent {
	return userJoinEvent{
		OrgName:  org.Account(),
		UserName: user.Account(),
		Type:     userJoinTypeInvite,
	}
}

type requestRejectEvent struct {
	Org       string `json:"org"`
	Passed    bool   `json:"passed"`
	Requester string `json:"requester"`
}

func (e *requestRejectEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewOrgRequestRejectEvent(org primitive.Account, requester primitive.Account) requestRejectEvent {
	return requestRejectEvent{
		Org:       org.Account(),
		Passed:    false,
		Requester: requester.Account(),
	}
}

// userRemoveEvent
type userRemoveEvent struct {
	OrgName  string `json:"org_name"`
	UserName string `json:"user_name"`
}

// Message returns the JSON representation of the userRemoveEvent.
func (e *userRemoveEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewUserRemoveEvent creates a new userRemoveEvent instance with the given Organization.
func NewUserRemoveEvent(a *OrgRemoveMemberCmd) userRemoveEvent {
	return userRemoveEvent{
		OrgName:  a.Org.Account(),
		UserName: a.Account.Account(),
	}
}

// orgDeleteEvent
type orgDeleteEvent struct {
	OrgName string `json:"org_name"`
}

// Message returns the JSON representation of the orgDeleteEvent.
func (e *orgDeleteEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewOrgDeleteEvent creates a new orgDeleteEvent instance with the given Organization.
func NewOrgDeleteEvent(a *OrgDeletedCmd) orgDeleteEvent {
	return orgDeleteEvent{
		OrgName: a.Name.Account(),
	}
}

type orgInviteEvent struct {
	Org     string `json:"org"`
	Inviter string `json:"inviter"`
	Invitee string `json:"invitee"`
}

func (o orgInviteEvent) Message() ([]byte, error) {
	return json.Marshal(o)
}

func NewOrgInviteEvent(org primitive.Account, inviter primitive.Account, invitee primitive.Account,
) orgInviteEvent {
	return orgInviteEvent{
		Org:     org.Account(),
		Inviter: inviter.Account(),
		Invitee: invitee.Account(),
	}
}

type orgRequest struct {
	Org       string `json:"org"`
	Requester string `json:"requester"`
}

func (o orgRequest) Message() ([]byte, error) {
	return json.Marshal(o)
}

func NewOrgRequestEvent(org primitive.Account, user primitive.Account) orgRequest {
	return orgRequest{
		Org:       org.Account(),
		Requester: user.Account(),
	}
}
