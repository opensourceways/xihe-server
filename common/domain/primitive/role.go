/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
)

const (
	Read  orgRole = "read"  // orgRoleRead in read team
	Write orgRole = "write" // orgRoleWrit in write team
	Admin orgRole = "admin" // orgRoleAdmin in owner team
)

// Role is an interface representing a org role.
type Role interface {
	Role() string
}

func RoleValidate(role string) bool {
	v := orgRole(role)

	return v == Read || v == Write || v == Admin
}

// Role creates a new Role instance from a string value.
func NewRole(v string) (Role, error) {
	if !RoleValidate(v) {
		return nil, errors.New("invalid role")
	}

	return orgRole(v), nil
}

// Role creates a new Role instance from a string value.
func NewReadRole() Role {
	return Read
}

// Role creates a new Role instance from a string value.
func NewWriteRole() Role {
	return Write
}

// Role creates a new Role instance from a string value.
func NewAdminRole() Role {
	return Admin
}

// CreateRole creates a new Role instance directly from a string value.
func CreateRole(v string) Role {
	return orgRole(v)
}

// orgRole represents the role of a user in an organization.
type orgRole string

// Role returns the string representation of the description.
func (r orgRole) Role() string {
	return string(r)
}
