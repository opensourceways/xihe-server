/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for handling organization-related operations.
package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	perm "github.com/opensourceways/xihe-server/organization/domain/permission"
	"github.com/opensourceways/xihe-server/organization/domain/repository"
)

// NewPermService creates a new PermService instance with the given configuration and org member.
func NewPermService(cfg *perm.Config, org repository.OrgMember) *PermService {
	p := &PermService{
		org: org,
	}
	p.Init(cfg)

	return p
}

func initActioin(actions []string) (bitmap uint64) {
	for _, action := range actions {
		switch action {
		case "write":
			bitmap |= 1 << primitive.ActionWrite
		case "read":
			bitmap |= 1 << primitive.ActionRead
		case "delete":
			bitmap |= 1 << primitive.ActionDelete
		case "create":
			bitmap |= 1 << primitive.ActionCreate
		default:
			logrus.Fatalf("invalid action: %s", action)
		}
	}

	return
}

func checkAction(bitmap uint64, action primitive.Action) bool {
	return bitmap&(1<<action) != 0
}

// Init initializes the permission service with the given configuration.
func (p *PermService) Init(cfg *perm.Config) {
	p.permissions = make(map[primitive.ObjType]map[primitive.Role]uint64)
	for _, permission := range cfg.Permissions {
		r := make(map[primitive.Role]uint64)
		for _, rule := range permission.Rules {
			role, err := primitive.NewRole(rule.Role)
			if err != nil {
				logrus.Fatalf("invalid role: %s", rule.Role)
			}

			r[role] = initActioin(rule.Operation)
		}

		p.permissions[primitive.ObjType(permission.ObjectType)] = r
	}
}

type PermService struct {
	permissions map[primitive.ObjType]map[primitive.Role]uint64

	org repository.OrgMember
}

func (p *PermService) doCheckPerm(role primitive.Role, objType primitive.ObjType, op primitive.Action) bool {
	if v, ok := p.permissions[objType][role]; ok {
		if checkAction(v, op) {
			return true
		}
	}

	return false
}

// Check checks if user can operate on organization's resource of a specified type.
func (p *PermService) Check(
	ctx context.Context,
	user primitive.Account,
	org primitive.Account,
	objType primitive.ObjType,
	op primitive.Action,
) error {
	if user == nil {
		e := xerrors.Errorf("user is nil")
		return allerror.NewNoPermission(e.Error(), e)
	}

	if org == nil {
		e := xerrors.Errorf("org is nil")
		return allerror.NewNoPermission(e.Error(), e)
	}

	m, err := p.org.GetByOrgAndUser(ctx, org.Account(), user.Account())
	if err != nil {
		e := xerrors.Errorf(
			"%s does not have a valid role in %s", user.Account(), org.Account(),
		)
		return allerror.NewNoPermission(e.Error(), e)
	}

	ok := p.doCheckPerm(m.Role, objType, op)
	res := "cannot"
	if ok {
		res = "can"
	}

	logrus.Debugf(
		"user %s (role %s) %s do %d on %s:%s",
		user.Account(), m.Role, res, op, org.Account(), objType,
	)

	if !ok {
		e := xerrors.Errorf(
			"%s %s %s permission denied", user.Account(), op.String(), string(objType),
		)
		return allerror.NewNoPermission(e.Error(), e)
	}

	return nil
}
