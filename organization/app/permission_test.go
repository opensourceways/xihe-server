/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for handling organization-related operations.
package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
	orgdomain "github.com/opensourceways/xihe-server/organization/domain"
	"github.com/opensourceways/xihe-server/organization/domain/permission"
)

type stubOrg struct{}

var stubMember = &orgdomain.OrgMember{
	Id:       primitive.CreateIdentity(1),
	OrgName:  primitive.CreateAccount("org"),
	Username: primitive.CreateAccount("XXXX"),
	Role:     primitive.NewAdminRole(),
}

var stub1Member = &orgdomain.OrgMember{
	Id:       primitive.CreateIdentity(1),
	OrgName:  primitive.CreateAccount("member"),
	Username: primitive.CreateAccount("ffff"),
	Role:     primitive.NewWriteRole(),
}

// Add adds a new organization member and returns it without any errors.
func (s *stubOrg) Add(o *orgdomain.OrgMember) (orgdomain.OrgMember, error) {
	return *o, nil
}

// Save saves the given organization member and returns it without any errors.
func (s *stubOrg) Save(o *orgdomain.OrgMember) (orgdomain.OrgMember, error) {
	return *o, nil
}

// Delete deletes the given organization member without returning any errors.
func (s *stubOrg) Delete(context.Context, *orgdomain.OrgMember) error {
	return nil
}

// DeleteByOrg deletes all organization members of the specified organization without returning any errors.
func (s *stubOrg) DeleteByOrg(primitive.Account) error {
	return nil
}

// GetByOrg retrieves all organization members of the specified organization and returns them without any errors.
func (s *stubOrg) GetByOrg(o *orgdomain.OrgListMemberCmd) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember, *stub1Member}, nil
}

// GetByOrgAndRole retrieves all organization members of the specified organization with the specified role
// and returns them without any errors.
func (s *stubOrg) GetByOrgAndRole(u string, r primitive.Role) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember}, nil
}

// GetByOrgAndUser retrieves the organization member of the specified organization and user
// and returns it without any errors.
func (s *stubOrg) GetByOrgAndUser(ctx context.Context, org, user string) (orgdomain.OrgMember, error) {
	if org == "org" && user == "XXXX" {
		return *stubMember, nil
	} else if org == "member" && user == "ffff" {
		return *stub1Member, nil
	}
	return orgdomain.OrgMember{}, fmt.Errorf("not found")
}

// GetByUser retrieves all organization members of the specified user and returns them without any errors.
func (s *stubOrg) GetByUser(string) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember}, nil

}

// GetByUserAndRoles retrieves all organization members of the specified user and returns them without any errors.
func (s *stubOrg) GetByUserAndRoles(primitive.Account, []primitive.Role) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember}, nil
}

// TestPermCheck is a test function for the permission check functionality.
func TestPermCheck(t *testing.T) {
	type testdata struct {
		user    primitive.Account
		org     primitive.Account
		objType primitive.ObjType
		op      primitive.Action
	}

	results := []bool{
		false,
		false,
		false,
		true,
		false,
		false,
		true,
		true,
	}

	tests := []testdata{
		{
			user:    nil,
			org:     nil,
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("123"),
			org:     nil,
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("123"),
			org:     nil,
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("XXXX"),
			org:     primitive.CreateAccount("org"),
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionCreate,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionDelete,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionCreate,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionWrite,
		},
	}
	var cfg permission.Config
	cfg.Permissions = []permission.PermObject{
		{
			ObjectType: string(primitive.ObjTypeOrg),
			Rules: []permission.Rule{
				{
					Role:      primitive.NewAdminRole().Role(),
					Operation: []string{"write", "read", "create", "delete"},
				},
			},
		},
		{
			ObjectType: string(primitive.ObjTypeMember),
			Rules: []permission.Rule{
				{
					Role:      primitive.NewWriteRole().Role(),
					Operation: []string{"write", "read"},
				},
			},
		},
	}
	app := NewPermService(&cfg, &stubOrg{})

	for i, test := range tests {
		err := app.Check(context.Background(), test.user, test.org, test.objType, test.op)
		if (err == nil) != results[i] {
			t.Errorf("case num %d valid result is %v ", i, err)
		}

	}
}
