/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "fmt"

// ObjType represents the type of object.
type ObjType string

const (
	ObjTypeUser     ObjType = "user"
	ObjTypeOrg      ObjType = "organization"
	ObjTypeModel    ObjType = "model"
	ObjTypeDataset  ObjType = "dataset"
	ObjTypeSpace    ObjType = "space"
	ObjTypeMember   ObjType = "member"
	ObjTypeInvite   ObjType = "invite"
	ObjTypeCodeRepo ObjType = "codeRepo"

	TokenPermWrite string = "write"
	TokenPermRead  string = "read"

	ActionRead Action = iota
	ActionWrite
	ActionDelete
	ActionCreate
)

// Action represents the action to be performed on an object.
type Action int

// String method for Action type.
func (a Action) String() string {
	switch a {
	case ActionRead:
		return "read"
	case ActionWrite:
		return "write"
	case ActionDelete:
		return "delete"
	case ActionCreate:
		return "create"
	default:
		return ""
	}
}

// IsModification checks if the action is a modification.
func (a Action) IsModification() bool {
	return a == ActionDelete || a == ActionWrite
}

type tokenPerm string

// TokenPerm returns the string representation of the token permission.
func (r tokenPerm) TokenPerm() string {
	return string(r)
}

// PermissionAllow checks if the given token has the required permission.
func (t tokenPerm) PermissionAllow(expect TokenPerm) bool {
	if expect.TokenPerm() == TokenPermRead {
		return true
	}

	if expect.TokenPerm() == TokenPermWrite {
		return t.TokenPerm() == TokenPermWrite
	}

	return false
}

// TokenPerm interface for token permission.
type TokenPerm interface {
	TokenPerm() string
	PermissionAllow(expect TokenPerm) bool
}

// NewTokenPerm creates a new TokenPerm instance.
func NewTokenPerm(v string) (TokenPerm, error) {
	if v != TokenPermWrite && v != TokenPermRead {
		return nil, fmt.Errorf("invalid permission(%s) , can only be %s/%s", v, TokenPermWrite, TokenPermRead)
	}

	return tokenPerm(v), nil
}

// NewReadPerm creates a read permission token.
func NewReadPerm() TokenPerm {
	return tokenPerm(TokenPermRead)
}

// NewWritePerm creates a write permission token.
func NewWritePerm() TokenPerm {
	return tokenPerm(TokenPermWrite)
}

// CreateTokenPerm creates a token with the given permission.
func CreateTokenPerm(v string) TokenPerm {
	return tokenPerm(v)
}
