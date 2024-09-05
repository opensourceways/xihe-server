/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain/privilege"
	userapp "github.com/opensourceways/xihe-server/user/app"
)

type action string

// const for privilege
const (
	// AllocNpu org who can alloc npu resource
	AllocNpu action = "npu"
	// Disable org who can disable model or space
	Disable action = "disable"
)

var disableObjType = []primitive.ObjType{primitive.ObjTypeSpace, primitive.ObjTypeModel, primitive.ObjTypeCodeRepo}

// NewAction returns an action type
func NewAction(action string) (action, error) {
	switch action {
	case string(AllocNpu):
		return AllocNpu, nil
	case string(Disable):
		return Disable, nil
	default:
		return "", fmt.Errorf("invalid action: %s", action)
	}
}

// PrivilegeOrgListOptions represents the options for listing organizations.
type PrivilegeOrgListOptions struct {
	User primitive.Account
	Type action
}

// PrivilegeOrg interface
type PrivilegeOrg interface {
	Contains(context.Context, primitive.Account) error
	List(context.Context, *PrivilegeOrgListOptions) ([]userapp.UserDTO, error)
	IsCanReadObj(action, primitive.ObjType) bool
}

// NewPrivilegeOrgService creates a new instance of PrivilegeOrg using
// the provided dependencies and configuration.
func NewPrivilegeOrgService(
	org OrgService,
	cfg privilege.PrivilegeConfig,
	t action,
) PrivilegeOrg {
	cfgs := make([]privilegeOrgConfig, 0)
	for _, c := range cfg.Orgs {
		a, err := newPrivilegeOrg(c)
		if err != nil {
			logrus.Warnf("invalid privilege org cfg: %v, %s", c, err.Error())
			continue
		}
		cfgs = append(cfgs, a)
	}

	if len(cfgs) == 0 {
		logrus.Warnf("empty privilege org config, privilege org will be ignored")
		return nil
	}

	logrus.Infof("init %s privilege org %v successfully", t, cfgs)

	return &privilegeOrg{
		cfgs:   cfgs,
		org:    org,
		Action: t,
	}
}

// privilegeOrgConfig represents the configuration for a privilege organization.
type privilegeOrgConfig struct {
	OrgId   string
	OrgName primitive.Account
}

// Create a new privilegeOrgConfig object based on the provided privilege.OrgIndex configuration.
func newPrivilegeOrg(cfg privilege.OrgIndex) (privilegeOrgConfig, error) {
	acc, err := primitive.NewAccount(cfg.OrgName)
	if err != nil {
		logrus.Warnf("invalid privileged org name: %s, privilege org will be ignored", cfg.OrgName)
		return privilegeOrgConfig{}, err
	}

	return privilegeOrgConfig{
		OrgId:   cfg.OrgId,
		OrgName: acc,
	}, nil
}

// privilegeOrg struct
type privilegeOrg struct {
	cfgs   []privilegeOrgConfig
	org    OrgService
	Action action
}

// Contains returns an error if the account is not in the privilege org.
func (p *privilegeOrg) Contains(ctx context.Context, account primitive.Account) error {
	if account == nil {
		e := fmt.Errorf("account is nil, cannot check privilege org")
		return allerror.NewInvalidParam(e.Error(), e)
	}

	for _, cfg := range p.cfgs {
		if err := cfg.contains(ctx, account, p.org); err == nil {
			return nil
		}
	}

	e := fmt.Errorf("account: %s not in %s privilege org", account.Account(), p.Action)
	return allerror.NewInvalidParam(e.Error(), e)
}

// Check if the privilegeOrgConfig contains the provided account within the specified organization.
func (p *privilegeOrgConfig) contains(ctx context.Context, account primitive.Account, org OrgService) error {
	o, err := org.GetByAccount(ctx, p.OrgName)
	if err != nil {
		logrus.Errorf("failed to get org: %s, %s", p.OrgName, err)
		return err
	}

	if o.Id != p.OrgId {
		e := fmt.Errorf("org id mismatch: actual: %s, config: %s", o.Id, p.OrgId)
		return allerror.New(allerror.ErrorCodePrivilegeOrgIdMismatch, e.Error(), e)
	}

	has := org.HasMember(ctx, primitive.CreateAccount(o.Account), account)
	if !has {
		e := fmt.Errorf("user: %s not in a privilegeorg", account.Account())
		return allerror.New(allerror.ErrorCodeNotInPrivilegeOrg, e.Error(), e)
	}

	return nil
}

// List returns the list of users in the privilege org.
func (p *privilegeOrg) List(ctx context.Context, l *PrivilegeOrgListOptions) ([]userapp.UserDTO, error) {
	if l == nil {
		e := fmt.Errorf("list options is nil")
		return nil, allerror.NewInvalidParam(e.Error(), e)
	}

	if l.Type != p.Action {
		return []userapp.UserDTO{}, nil
	}

	orgs := make([]userapp.UserDTO, 0)
	for _, cfg := range p.cfgs {
		o, err := p.org.GetByAccount(ctx, cfg.OrgName)
		if err != nil {
			logrus.Errorf("failed to get org: %s, %s", cfg.OrgName, err)
			return nil, err
		}

		if l.User != nil {
			has := p.org.HasMember(ctx, primitive.CreateAccount(o.Account), l.User)
			if !has {
				continue
			}
		}

		orgs = append(orgs, o)
	}

	return orgs, nil
}

// Check if the privilegeOrg allows reading the specified object for the given action.
func (p *privilegeOrg) IsCanReadObj(action action, obj primitive.ObjType) bool {
	var objs []primitive.ObjType
	if action == Disable {
		objs = disableObjType
	}

	for _, o := range objs {
		if o == obj {
			return true
		}
	}

	return false
}
