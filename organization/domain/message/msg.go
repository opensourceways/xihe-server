/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package message provides functionality for sending and handling event messages.
package message

// EventMessage is an interface that represents an event message.
type EventMessage interface {
	Message() ([]byte, error)
}

// OrganizationMessage is an interface that defines a method for sending a space app created event.
type OrganizationMessage interface {
	SendComputilityUserJoinEvent(EventMessage) error
	SendComputilityUserRemoveEvent(EventMessage) error
	SendComputilityOrgDeleteEvent(EventMessage) error
	SendOrgInviteEvent(EventMessage) error
	SendOrgRequestEvent(EventMessage) error
	SendOrgRequestRejectEvent(EventMessage) error
}
