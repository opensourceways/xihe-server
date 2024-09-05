/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for handling organization-related operations.
package app

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/common/domain/trace"
	git "github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/organization/domain"
	"github.com/opensourceways/xihe-server/organization/domain/message"
	orgprimitive "github.com/opensourceways/xihe-server/organization/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain/repository"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func errOrgNotFound(msg string, err error) error {
	if msg == "" {
		msg = "not found"
	}

	return allerror.NewNotFound(allerror.ErrorCodeOrganizationNotFound, msg, err)
}

// OrgService is an interface that defines the methods for organization-related operations.
type OrgService interface {
	Create(context.Context, *domain.OrgCreatedCmd) (userapp.UserDTO, error)
	Delete(context.Context, *domain.OrgDeletedCmd) error
	UpdateBasicInfo(context.Context, *domain.OrgUpdatedBasicInfoCmd) (userapp.UserDTO, error)
	GetByAccount(context.Context, primitive.Account) (userapp.UserDTO, error)
	GetOrgOrUser(context.Context, primitive.Account, primitive.Account) (userapp.UserDTO, error)
	ListAccount(*userrepo.ListOption) ([]userapp.UserDTO, error)

	CheckName(context.Context, primitive.Account) bool
	GetByOwner(context.Context, primitive.Account, primitive.Account) ([]userapp.UserDTO, error)
	GetByUser(context.Context, primitive.Account, primitive.Account) ([]userapp.UserDTO, error)
	List(context.Context, *OrgListOptions) ([]userapp.UserDTO, int, error)
	OrgInfoList(ctx context.Context, opt *OrgListOptions) ([]userapp.BaseUserDTO, int, error)
	OrgGetNoSort(ctx context.Context, opt *OrgListOptions) ([]userapp.BaseUserDTO, int, error)
	OrgByTypeInfoList(ctx context.Context, opt *OrgListOptions) ([]userapp.BaseUserDTO, int, error)
	ResOrgInfoList([]userapp.BaseUserDTO, int) userapp.UserPaginationDTO
	SortAllOrg(ctx context.Context) error
	IsSortInit(ctx context.Context) error
	HasMember(context.Context, primitive.Account, primitive.Account) bool
	InviteMember(context.Context, *domain.OrgInviteMemberCmd) (ApproveDTO, error)
	RequestMember(context.Context, *domain.OrgRequestMemberCmd) (MemberRequestDTO, error)
	CancelReqMember(context.Context, *domain.OrgCancelRequestMemberCmd) (MemberRequestDTO, error)
	ApproveRequest(context.Context, *domain.OrgApproveRequestMemberCmd) (MemberRequestDTO, error)
	AcceptInvite(context.Context, *domain.OrgAcceptInviteCmd) (ApproveDTO, error)
	RevokeInvite(context.Context, *domain.OrgRemoveInviteCmd) (ApproveDTO, error)
	GetOnlyApply(string, string) ([]domain.MemberRequest, error)
	ListMemberReq(context.Context, *domain.OrgMemberReqListCmd) (MemberPagnationDTO, error)
	ListInvitationByInvitee(
		context.Context, primitive.Account, primitive.Account, domain.ApproveStatus) ([]ApproveDTO, error)
	ListInvitationByInviter(
		context.Context, primitive.Account, primitive.Account, domain.ApproveStatus) ([]ApproveDTO, error)
	ListInvitationByOrg(
		context.Context, primitive.Account, primitive.Account, domain.ApproveStatus) ([]ApproveDTO, error)
	AddMember(context.Context, *domain.OrgAddMemberCmd) error
	RemoveMember(context.Context, *domain.OrgRemoveMemberCmd) error
	EditMember(context.Context, *domain.OrgEditMemberCmd) (MemberDTO, error)
	ListMember(context.Context, *domain.OrgListMemberCmd) ([]MemberDTO, error)
	GetMemberByUserAndOrg(context.Context, primitive.Account, primitive.Account) (MemberDTO, error)
}

// NewOrgService creates a new instance of the OrgService.
func NewOrgService(
	user userapp.UserService,
	repo userrepo.User,
	member repository.OrgMember,
	invite repository.Approve,
	perm *PermService,
	cfg *domain.Config,
	git git.User,
	message message.OrganizationMessage,
	cert repository.Certificate,
	trace trace.Trace,
) OrgService {
	return &orgService{
		user:         user,
		repo:         repo,
		member:       member,
		perm:         perm,
		invite:       invite,
		defaultRole:  primitive.CreateRole(cfg.DefaultRole),
		inviteExpiry: cfg.InviteExpiry,
		config:       cfg,
		git:          git,
		message:      message,
		certificate:  cert,
		tracer:       trace,
	}
}

type orgService struct {
	config       *domain.Config
	inviteExpiry int64
	defaultRole  primitive.Role
	user         userapp.UserService
	repo         userrepo.User
	member       repository.OrgMember
	invite       repository.Approve
	perm         *PermService
	git          git.User
	message      message.OrganizationMessage
	certificate  repository.Certificate
	tracer       trace.Trace
}

// Create creates a new organization with the given command and returns the created organization as a UserDTO.
func (org *orgService) Create(ctx context.Context, cmd *domain.OrgCreatedCmd) (o userapp.UserDTO, err error) {
	orgTemp, err := cmd.ToOrg()
	if err != nil {
		return
	}

	if !org.repo.CheckName(ctx, cmd.Name) {
		e := fmt.Errorf("name %s is already been taken", cmd.Name.Account())
		err = allerror.New(allerror.ErrorNameAlreadyBeenTaken, e.Error(), e)
		return
	}

	if err = org.orgCountCheck(ctx, cmd.Owner); err != nil {
		return
	}

	owner, err := org.repo.GetByAccount(ctx, cmd.Owner)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedGetOwnerInfo, "failed to get owner info", err)
		return
	}

	pl, err := org.user.GetPlatformUser(ctx, orgTemp.Owner)
	if err != nil {
		err = allerror.New(allerror.ErrorFailGetPlatformUser, err.Error(), err)
		return
	}

	err = pl.CreateOrg(orgTemp)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedCreateOrg, "failed to create org", err)
		return
	}

	orgTemp.DefaultRole = org.defaultRole
	orgTemp.AllowRequest = false
	orgTemp.OwnerId = owner.Id

	*orgTemp, err = org.repo.AddOrg(ctx, orgTemp)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedCreateToOrg, "failed to create to org", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	_, err = org.member.Add(&domain.OrgMember{
		OrgName:  cmd.Name,
		OrgId:    orgTemp.Id,
		Username: cmd.Owner,
		FullName: owner.Fullname,
		UserId:   owner.Id,
		Role:     primitive.NewAdminRole(),
	})
	if err != nil {
		err = allerror.New(allerror.ErrorFailSaveOrgMember, "failed to create to org", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	o = ToDTO(orgTemp)

	return
}

func (org *orgService) orgCountCheck(ctx context.Context, owner primitive.Account) error {
	total, err := org.repo.GetOrgCountByOwner(ctx, owner)
	if err != nil {
		return err
	}

	if total >= org.config.MaxCountPerOwner {
		return allerror.NewCountExceeded("org count exceed", fmt.Errorf("org count(now:%d max:%d) exceed",
			total, org.config.MaxCountPerOwner))
	}

	return nil
}

// GetByAccount retrieves an organization by its account and returns it as a UserDTO.
func (org *orgService) GetByAccount(ctx context.Context, acc primitive.Account) (dto userapp.UserDTO, err error) {
	o, err := org.repo.GetOrgByName(ctx, acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", acc.Account()),
				fmt.Errorf("org %s not found, %w",
					acc.Account(), err))
		}

		return
	}

	dto = ToDTO(&o)
	return
}

// GetOrgOrUser retrieves either an organization or a user by their account and returns it as a UserDTO.
func (org *orgService) GetOrgOrUser(
	ctx context.Context, actor, acc primitive.Account) (dto userapp.UserDTO, err error) {
	u, err := org.repo.GetByAccount(ctx, acc)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		return
	} else if err == nil {
		dto = userapp.NewUserDTO(&u, actor)
		return
	}

	o, err := org.repo.GetOrgByName(ctx, acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeUserNotFound, fmt.Sprintf("org %s not found", acc.Account()),
				fmt.Errorf("org %s not found, %w", acc.Account(), err))
		}

		return
	}

	dto = ToDTO(&o)
	return
}

// ListAccount lists organizations based on the provided options and returns them as a slice of UserDTOs.
func (org *orgService) ListAccount(opt *userrepo.ListOption) (dtos []userapp.UserDTO, err error) {
	return
}

// Delete deletes an organization based on the provided command and returns an error if any occurs.
func (org *orgService) Delete(ctx context.Context, cmd *domain.OrgDeletedCmd) error {
	err := org.perm.Check(ctx, cmd.Actor, cmd.Name, primitive.ObjTypeOrg, primitive.ActionDelete)
	if err != nil {
		return err
	}
	o, err := org.repo.GetOrgByName(ctx, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return err
	}

	pl, err := org.user.GetPlatformUser(ctx, o.Owner)
	if err != nil {
		return allerror.New(allerror.ErrorFailGetPlatformUser, "failed to get platform user",
			fmt.Errorf("failed to get platform user, %w", err))
	}

	can, err := pl.CanDelete(cmd.Name)
	if err != nil {
		return allerror.New(allerror.ErrorAccountCannotDeleteTheOrg, "can't delete the org",
			fmt.Errorf("%s can't delete the org, %w", cmd.Name.Account(), err))
	}

	if !can {
		e := fmt.Errorf("can't delete the organization, while some repos still existed")
		return allerror.New(allerror.ErrorCodeOrgExistResource, e.Error(), e)
	}

	err = org.member.DeleteByOrg(o.Account)
	if err != nil {
		return allerror.New(allerror.ErrorBaseCase, "failed to delete org member",
			fmt.Errorf("failed to delete org member, %w", err))
	}

	err = org.invite.DeleteInviteAndReqByOrg(o.Account)
	if err != nil {
		return allerror.New(allerror.ErrorBaseCase, "failed to delete org invite",
			fmt.Errorf("failed to delete org invite, %w", err))
	}

	err = org.git.DeleteOrg(o.Account)
	if err != nil {
		return allerror.New(allerror.ErrorBaseCase, "failed to delete git org",
			fmt.Errorf("failed to delete git org, %w", err))
	}

	err = org.repo.DeleteOrg(ctx, &o)
	if err != nil {
		return err
	}

	err = org.certificate.DeleteByOrgName(cmd.Name)
	if err != nil {
		logrus.Errorf("delete certificate of %s failed", cmd.Name.Account())
	}

	event := domain.NewOrgDeleteEvent(cmd)
	err = org.message.SendComputilityOrgDeleteEvent(&event)
	if err != nil {
		e := xerrors.Errorf("send message to delete computility org failed: %w", err)
		err = allerror.New(allerror.ErrorMsgPublishFailed, "delete computility org failed", e)
	}

	return err
}

// UpdateBasicInfo updates the basic information of an organization based on the provided command
// and returns the updated organization as a UserDTO.
func (org *orgService) UpdateBasicInfo(
	ctx context.Context, cmd *domain.OrgUpdatedBasicInfoCmd) (dto userapp.UserDTO, err error) {
	if cmd == nil {
		err = allerror.New(allerror.ErrorSystemError, "system error", err)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	err = org.perm.Check(ctx, cmd.Actor, cmd.OrgName, primitive.ObjTypeOrg, primitive.ActionWrite)
	if err != nil {
		return
	}

	o, err := org.repo.GetOrgByName(ctx, cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.OrgName.Account()),
				fmt.Errorf("org %s not found: %w", cmd.OrgName.Account(), err))
		}

		return
	}

	change := cmd.ToOrg(&o)
	if change {
		o, err = org.repo.SaveOrg(ctx, &o)
		if err != nil {
			err = allerror.New(allerror.ErrorFailedToSaveOrg, "failed to save org",
				fmt.Errorf("failed to save org, %w", err))
			return
		}
		dto = ToDTO(&o)
		return
	}

	err = allerror.New(allerror.ErrorNothingChanged, "nothing changed",
		fmt.Errorf("nothing changed when update basic info %v", cmd))
	return
}

// GetByOwner retrieves organizations owned by the specified account and returns them as a slice of UserDTOs.
func (org *orgService) GetByOwner(
	ctx context.Context, actor, acc primitive.Account) (orgs []userapp.UserDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	orgs, _, err = org.List(ctx, &OrgListOptions{
		Owner: acc,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
	}

	return
}

// GetByUser retrieves organizations associated with a user.
func (org *orgService) GetByUser(
	ctx context.Context, actor, acc primitive.Account) (orgs []userapp.UserDTO, err error) {
	if acc == nil {
		e := fmt.Errorf("account is nil")
		err = allerror.New(allerror.ErrorSystemError, "account is nil", e)
		return
	}

	members, err := org.member.GetByUser(acc.Account())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	orgs = make([]userapp.UserDTO, len(members))
	for i := range members {
		o, e := org.repo.GetOrgByName(ctx, members[i].OrgName)
		if e != nil {
			e := fmt.Errorf("failed to get org when get org by user: %w", e)
			err = allerror.New(allerror.ErrorFailedToGetOrg, e.Error(), e)
			return
		}
		orgs[i] = ToDTO(&o)
	}

	return
}

// List retrieves a list of organizations based on the provided options.
func (org *orgService) List(ctx context.Context, l *OrgListOptions) (orgs []userapp.UserDTO, total int, err error) {
	if l == nil {
		e := fmt.Errorf("list options is nil")
		return nil, 0, allerror.New(allerror.ErrorSystemError, e.Error(), e)
	}
	orgs = []userapp.UserDTO{}

	var orgIDs []int64
	if l.Member != nil {
		orgIDs, err = org.getOrgIDsByUserAndRoles(l.Member, l.Roles)
		if err != nil || len(orgIDs) == 0 {
			return
		}
	}

	listOption := &userrepo.ListPageOrgOption{
		OrgIDs:   orgIDs,
		Owner:    l.Owner,
		PageNum:  l.Page,
		PageSize: l.PageSize,
	}
	os, total, err := org.repo.GetOrgPageList(ctx, listOption)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	orgs = make([]userapp.UserDTO, len(os))
	for i := range os {
		orgs[i] = ToDTO(&os[i])
	}

	return orgs, int(total), nil
}

func (org *orgService) OrgInfoList(ctx context.Context,
	opt *OrgListOptions) (orgs []userapp.BaseUserDTO, total int, err error) {
	var orgIDs []int64
	if opt == nil {
		e := fmt.Errorf("list options is nil")
		return nil, 0, allerror.New(allerror.ErrorSystemError, e.Error(), e)
	}
	if opt.Member != nil {
		orgIDs, err = org.getOrgIDsByUserAndRoles(opt.Member, opt.Roles)
		if err != nil || len(orgIDs) == 0 {
			return []userapp.BaseUserDTO{}, 0, err
		}
	}
	listOption := &userrepo.ListPageOrgOption{
		OrgIDs:   orgIDs,
		Owner:    opt.Owner,
		PageNum:  opt.Page,
		PageSize: opt.PageSize,
		Search:   opt.Search,
		OrgType:  opt.OrgType,
	}
	orgCertInfo, _,
		err := org.certificate.FindList(ctx, listOption.PageNum, listOption.PageSize, opt.OrgType)
	if err != nil {
		return []userapp.BaseUserDTO{}, 0, err
	}
	var names []string
	orgMap := map[string]string{}
	if len(orgCertInfo) != 0 {
		for _, value := range orgCertInfo {
			names = append(names, value.OrgName.Account())
			orgMap[value.OrgName.Account()] = value.CertificateOrgType.CertificateOrgType()
		}
		listOption.Names = names
	}
	resOrg, total, err := org.repo.GetOrgInfoList(ctx, listOption)
	if err != nil || len(resOrg) == 0 {
		return []userapp.BaseUserDTO{}, 0, err
	}
	for index := range resOrg {
		if utils.StrIn(resOrg[index].Account.Account(), names) {
			resOrg[index].OrgType = orgMap[resOrg[index].Account.Account()]
		} else {
			resOrg[index].OrgType = ""
		}
		orgs = append(orgs, ToBaseDTO(&resOrg[index]))
	}
	return orgs, total, nil
}

func (org *orgService) OrgGetNoSort(ctx context.Context,
	opt *OrgListOptions) (orgs []userapp.BaseUserDTO, total int, err error) {
	var orgIDs []int64
	if opt == nil {
		e := fmt.Errorf("list options is nil")
		return nil, 0, allerror.New(allerror.ErrorSystemError, e.Error(), e)
	}
	if opt.Member != nil {
		orgIDs, err = org.getOrgIDsByUserAndRoles(opt.Member, opt.Roles)
		if err != nil || len(orgIDs) == 0 {
			return []userapp.BaseUserDTO{}, 0, err
		}
	}
	listOption := &userrepo.ListPageOrgOption{
		OrgIDs:   orgIDs,
		Owner:    opt.Owner,
		PageNum:  opt.Page,
		PageSize: opt.PageSize,
		Search:   opt.Search,
		OrgType:  opt.OrgType,
	}
	res, total, err := org.repo.GetOrgInfoList(ctx, listOption)
	if err != nil || len(res) == 0 {
		return []userapp.BaseUserDTO{}, 0, err
	}
	for index := range res {
		orgCertInfo, err := org.certificate.Find(ctx, repository.FindOption{OrgName: res[index].Account})
		if err != nil || orgCertInfo.Status.CertificateStatus() != orgprimitive.NewPassedStatus().CertificateStatus() {
			res[index].OrgType = ""
		} else {
			res[index].OrgType = orgCertInfo.CertificateOrgType.CertificateOrgType()
		}
		orgs = append(orgs, ToBaseDTO(&res[index]))
	}
	return orgs, total, nil
}

func (org *orgService) OrgByTypeInfoList(ctx context.Context,
	opt *OrgListOptions) ([]userapp.BaseUserDTO, int, error) {
	var orgIDs []int64
	if opt.Member != nil {
		orgIDs, err := org.getOrgIDsByUserAndRoles(opt.Member, opt.Roles)
		if err != nil || len(orgIDs) == 0 {
			return []userapp.BaseUserDTO{}, 0, err
		}
	}
	listOption := &userrepo.ListPageOrgOption{
		OrgIDs:   orgIDs,
		Owner:    opt.Owner,
		PageNum:  opt.Page,
		PageSize: opt.PageSize,
		Search:   opt.Search,
		OrgType:  opt.OrgType,
	}
	orgCertInfo, _,
		err := org.certificate.FindList(ctx, listOption.PageNum, listOption.PageSize, opt.OrgType)
	if err != nil {
		return []userapp.BaseUserDTO{}, 0, err
	}
	if len(orgCertInfo) == 0 {
		return []userapp.BaseUserDTO{}, 0, nil
	}
	var names []string
	orgMap := map[string]string{}
	for _, value := range orgCertInfo {
		names = append(names, value.OrgName.Account())
		orgMap[value.OrgName.Account()] = value.CertificateOrgType.CertificateOrgType()
	}
	if len(names) == 0 {
		return []userapp.BaseUserDTO{}, 0, nil
	}
	res, _, err := org.repo.GetOrgInfoListByNames(ctx, listOption, names)
	logrus.Info(res)
	if err != nil || len(res) == 0 {
		return []userapp.BaseUserDTO{}, 0, err
	}
	if len(res) == 0 {
		return []userapp.BaseUserDTO{}, 0, nil
	}
	var orgs []userapp.BaseUserDTO
	for index := range res {
		if res[index].Account == nil {
			return []userapp.BaseUserDTO{}, 0, err
		}
		res[index].OrgType = orgMap[res[index].Account.Account()]
		orgs = append(orgs, ToBaseDTO(&res[index]))
	}
	return orgs, len(orgs), nil
}

func (org *orgService) GetAllCertName(ctx context.Context) ([]string, error) {
	orgCertInfoList, err := org.certificate.FinAllName(ctx)
	if err != nil {
		return nil, err
	}
	return orgCertInfoList, nil
}

func (org *orgService) ResOrgInfoList(dos []userapp.BaseUserDTO, total int) (res userapp.UserPaginationDTO) {
	if len(dos) == 1 {
		res.Total = total
		res.Labels = dos
		return res
	}
	if len(dos) == 0 {
		res.Labels = []userapp.BaseUserDTO{}
		res.Total = total
		return res
	}
	res.Total = total
	res.Labels = dos
	return res
}

func (org *orgService) SortAllOrg(ctx context.Context) error {
	data, err := org.repo.GetAllOrgInfoList(ctx)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	maxTotal := 0
	var orgs []domain.Organization
	for index := range data {
		total := data[index].ModelNum + data[index].SpaceNum + data[index].DatasetNum
		if total > maxTotal {
			maxTotal = total
		}
		data[index].Score = total
		orgs = append(orgs, data[index].TOUser())
	}
	certNames, err := org.GetAllCertName(ctx)
	if err != nil {
		return err
	}
	var res []domain.OrgUpdate
	for index := range orgs {
		if utils.StrIn(orgs[index].Account.Account(), certNames) {
			orgs[index].Score = orgs[index].Score + maxTotal + 1
		}
		res = append(res, domain.OrgUpdate{Name: orgs[index].Account.Account(),
			Score: orgs[index].Score})
	}
	err = org.repo.Updates(ctx, res)
	if err != nil {
		return err
	}
	return nil
}

func (org *orgService) IsSortInit(ctx context.Context) error {
	certNames, err := org.GetAllCertName(ctx)
	if err != nil {
		err = xerrors.Errorf("get certName failed: %w", err)
		return err
	}
	ok, err := org.repo.GetAllScore(ctx, certNames)
	if err != nil {
		err = xerrors.Errorf("Score select failed: %w", err)
		return err
	}
	if !ok {
		err = org.SortAllOrg(ctx)
		if err != nil {
			err = xerrors.Errorf("sort function failed: %w", err)
			return err
		}
	}
	return nil
}

func (org *orgService) getOrgIDsByUserAndRoles(user primitive.Account,
	roles []primitive.Role) (orgIDs []int64, err error) {
	members, err := org.member.GetByUserAndRoles(user, roles)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
		return
	}
	for _, mem := range members {
		orgIDs = append(orgIDs, mem.OrgId.Integer())
	}
	return
}

// ListMember retrieves a list of members for a given organization.
func (org *orgService) ListMember(ctx context.Context, cmd *domain.OrgListMemberCmd) (dtos []MemberDTO, err error) {
	if cmd == nil || cmd.Org == nil {
		e := fmt.Errorf("org account is nil")
		err = allerror.New(allerror.ErrorSystemError, e.Error(), e)
		return
	}

	o, err := org.GetByAccount(ctx, cmd.Org)
	if err != nil {
		return
	}

	members, e := org.member.GetByOrg(cmd)
	if e != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	dtos = make([]MemberDTO, len(members))
	for i := range members {
		dtos[i] = ToMemberDTO(&members[i])
		dtos[i].OrgName = o.Account
		dtos[i].OrgFullName = o.Fullname
		u, err := org.user.GetUserAvatarId(context.Background(), members[i].Username)
		if err != nil {
			logrus.Errorf("list org members err: get avatar id error: %v", err)

			continue
		}
		dtos[i].AvatarId = u.AvatarId
	}

	return
}

// AddMember adds a new member to an organization.
func (org *orgService) AddMember(ctx context.Context, cmd *domain.OrgAddMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		e := fmt.Errorf("failed to validate cmd: %w", err)
		return allerror.New(allerror.ErrorFailedToValidateCmd, e.Error(), e)
	}

	o, err := org.repo.GetOrgByName(ctx, cmd.Org)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToGetOrgInfo, "failed to get org info",
			fmt.Errorf("failed to get org info: %w", err))
	}

	memberInfo, err := org.repo.GetByAccount(ctx, cmd.User)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToGetMemberInfo, "failed to get member info",
			fmt.Errorf("failed to get member info: %w", err))
	}

	m := cmd.ToMember(memberInfo)

	pl, err := org.user.GetPlatformUser(ctx, cmd.Actor)
	if err != nil {
		return allerror.New(allerror.ErrorFailGetPlatformUser,
			"failed to get platform user for adding member",
			fmt.Errorf("failed to get platform user for adding member: %w", err))
	}

	err = pl.AddMember(&o, &m)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToAddMemberToOrg, fmt.Sprintf("failed to add member:%s to org:%s",
			m.Username.Account(), o.Account.Account()), fmt.Errorf("failed to add member:%s to org:%s err: %w",
			m.Username.Account(), o.Account.Account(), err))
	}

	_, err = org.member.Add(&m)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToSaveMemberForAddingMember, "failed to save member for adding member",
			fmt.Errorf("failed to save member for adding member: %w", err))
	}

	return nil
}

func (org *orgService) canEditMember(cmd *domain.OrgEditMemberCmd) (err error) {
	return org.canRemoveMember(&domain.OrgRemoveMemberCmd{
		Org:     cmd.Org,
		Account: cmd.Account,
		Actor:   cmd.Actor,
		Msg:     "",
	})
}

func (org *orgService) members(orgName primitive.Account) ([]domain.OrgMember, int, error) {
	members, err := org.member.GetByOrg(&domain.OrgListMemberCmd{Org: orgName})
	if err != nil {
		e := fmt.Errorf("failed to get members by org name: %s err: %w", orgName, err)
		err = allerror.New(allerror.ErrorFailedToGetMembersByOrgName, e.Error(), e)
		return []domain.OrgMember{}, 0, err
	}

	return members, len(members), nil

}

func (org *orgService) getOwners(orgName primitive.Account) ([]domain.OrgMember, error) {
	members, err := org.member.GetByOrg(&domain.OrgListMemberCmd{Org: orgName, Role: primitive.Admin})
	if err != nil {
		e := fmt.Errorf("failed to get members by org name: %s err: %w", orgName, err)
		err = allerror.New(allerror.ErrorFailedToGetMembersByOrgName, e.Error(), e)
		return []domain.OrgMember{}, err
	}

	if len(members) == 0 {
		e := fmt.Errorf("no owners found in org %s", orgName.Account())
		err = allerror.NewInvalidParam(e.Error(), e)
		return []domain.OrgMember{}, err
	}

	return members, nil
}

func (org *orgService) canRemoveMember(cmd *domain.OrgRemoveMemberCmd) (err error) {
	// check if this is the only owner
	members, count, err := org.members(cmd.Org)
	if err != nil {
		return err
	}
	if count == 1 {
		e := fmt.Errorf("the org has only one member")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	if count == 0 {
		e := fmt.Errorf("the org has no member")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	member := cmd.ToMember()

	ownerCount := 0
	removeOwner := false
	can := false
	for _, m := range members {
		if m.Role == primitive.Admin {
			ownerCount++
			if m.Username == member.Username {
				removeOwner = true
				can = true
			}
		}
		if m.Username == member.Username {
			can = true
		}
	}

	if ownerCount == 1 && removeOwner {
		e := fmt.Errorf("the only owner can not be removed")
		err = allerror.New(allerror.ErrorOnlyOwnerCanNotBeRemoved, e.Error(), e)
		return
	}

	if !can {
		e := fmt.Errorf("the member is not in the org")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	return
}

// RemoveMember removes a member from an organization.
func (org *orgService) RemoveMember(ctx context.Context, cmd *domain.OrgRemoveMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		e := fmt.Errorf("failed to validate cmd: %w", err)
		return allerror.New(allerror.ErrorFailedToValidateCmd, e.Error(), e)
	}

	err = org.canRemoveMember(cmd)
	if err != nil {
		e := fmt.Errorf("failed to validate cmd: %w", err)
		return allerror.New(allerror.ErrorFailedToRemoveMember, e.Error(), e)
	}

	o, err := org.repo.GetOrgByName(ctx, cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()), err)
		}

		return err
	}

	_, err = org.repo.GetByAccount(ctx, cmd.Actor)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, fmt.Sprintf("user %s not existed",
				cmd.Actor.Account()), fmt.Errorf("user %s not existed: %w", cmd.Actor.Account(), err))
		}

		return err
	}

	if cmd.Actor.Account() != cmd.Account.Account() {
		_, err = org.repo.GetByAccount(ctx, cmd.Account)
		if err != nil {
			if commonrepo.IsErrorResourceNotExists(err) {
				err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, fmt.Sprintf("user %s not existed",
					cmd.Account.Account()), fmt.Errorf("user %s not existed: %w", cmd.Actor.Account(), err))
			}

			return err
		}
		err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionDelete)
		if err != nil {
			return err
		}
	} else {
		err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionRead)
		if err != nil {
			return err
		}
	}

	owners, err := org.getOwners(cmd.Org)
	if err != nil {
		e := fmt.Errorf("failed to get owners of org when add new member: %s err: %w", cmd.Org.Account(), err)
		return allerror.NewInvalidParam(e.Error(), e)
	}

	pl, err := org.user.GetPlatformUser(ctx, owners[0].Username)
	if err != nil {
		return allerror.New(allerror.ErrorFailGetPlatformUser, err.Error(), err)
	}

	m, err := org.member.GetByOrgAndUser(ctx, cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		e := fmt.Errorf("failed to get member when remove member by org %s and user %s, %w",
			cmd.Org.Account(), cmd.Account.Account(), err)
		return allerror.NewInvalidParam("failed to remove member", e)
	}

	err = pl.RemoveMember(&o, &m)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToDeleteGitMember, "failed to delete git member",
			fmt.Errorf("failed to delete git member, %w", err))
	}

	err = org.member.Delete(ctx, &m)
	if err != nil {
		_ = pl.AddMember(&o, &m)
		return allerror.New(allerror.ErrorFailedToDeleteMember, "failed to delete member",
			fmt.Errorf("failed to delete member, %w", err))
	}

	// when owner is removed, a new owner must be set
	if cmd.Account == o.Owner {
		o.Owner = cmd.Actor
		_, err = org.repo.SaveOrg(ctx, &o)
		if err != nil {
			e := fmt.Errorf("failed to change owner of org, %w", err)
			return allerror.New(allerror.ErrorFailedToChangeOwnerOfOrg, e.Error(), e)
		}
	}

	event := domain.NewUserRemoveEvent(cmd)
	err = org.message.SendComputilityUserRemoveEvent(&event)
	if err != nil {
		e := xerrors.Errorf("send message to user removed from computility org failed: %w", err)
		err = allerror.New(allerror.ErrorMsgPublishFailed, "user remove computility failed", e)
	}

	return err
}

// EditMember edits the role of a member in an organization.
func (org *orgService) EditMember(ctx context.Context, cmd *domain.OrgEditMemberCmd) (dto MemberDTO, err error) {
	if err = org.canEditMember(cmd); err != nil {
		return
	}

	err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionWrite)
	if err != nil {
		return
	}

	m, err := org.member.GetByOrgAndUser(ctx, cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		err = fmt.Errorf("failed to get member when edit member by org:%s and user:%s err: %w",
			cmd.Org.Account(), cmd.Account.Account(), err)
		return
	}

	o, err := org.repo.GetOrgByName(ctx, cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()), err)
		}

		return
	}

	pl, err := org.user.GetPlatformUser(ctx, cmd.Actor)
	if err != nil {
		err = fmt.Errorf("failed to get platform user: %w", err)
		return
	}

	if m.Role != cmd.Role {
		origRole := m.Role
		m.Role = cmd.Role
		err = pl.EditMemberRole(&o, origRole, &m)
		if err != nil {
			err = fmt.Errorf("failed to edit git member: %w", err)
			return
		}

		m, err = org.member.Save(&m)
		if err != nil {
			err = fmt.Errorf("failed to save member: %w", err)
			return
		}
		dto = ToMemberDTO(&m)
	} else {
		logrus.Warn("role not changed")
	}

	return
}

// InviteMember invites a new member to an organization.
func (org *orgService) InviteMember(ctx context.Context, cmd *domain.OrgInviteMemberCmd) (dto ApproveDTO, err error) {
	if org.IsInvite(ctx, cmd.Org, cmd.Account) {
		e := fmt.Errorf("the user is invited too")
		err = allerror.New(allerror.ErrorInviteBadRequest, e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.HasMember(ctx, cmd.Org, cmd.Account) {
		e := fmt.Errorf("the user is already a member of the org")
		err = allerror.New(allerror.ErrorUserAlreadyInOrg, e.Error(), e)
		return
	}

	invitee, err := org.repo.GetByAccount(ctx, cmd.Account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "invitee not found", err)
		}

		return
	}
	if invitee.Email.Email() == "" {
		e := fmt.Errorf("user email failed")
		err = allerror.New(allerror.ErrorUserBadEmail, e.Error(), e)
		return
	}

	inviter, err := org.repo.GetByAccount(ctx, cmd.Actor)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "inviter not found", err)
		}

		return
	}

	o, err := org.repo.GetOrgByName(ctx, cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "organization not found", err)
		}

		return
	}

	err = org.canInvite(ctx, cmd)
	if err != nil {
		return
	}

	invite := cmd.ToApprove(org.inviteExpiry)
	invite.InviterId = inviter.Id
	invite.UserId = invitee.Id
	invite.OrgId = o.Id

	var isUpdate bool
	*invite, isUpdate, err = org.invite.AddInvite(invite)
	if err != nil {
		err = fmt.Errorf("failed to save member: %w", err)
		return
	}

	if !isUpdate {
		event := domain.NewOrgInviteEvent(cmd.Org, inviter.Account, invitee.Account)
		if err = org.message.SendOrgInviteEvent(event); err != nil {
			logrus.Errorf("send org invite message failed: %s", err.Error())
		}
	}

	dto = ToApproveDTO(ctx, invite, org.user)

	return
}

func (org *orgService) canInvite(ctx context.Context, cmd *domain.OrgInviteMemberCmd) error {
	err := org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionCreate)
	if err != nil {
		return err
	}

	c, err := org.invite.Count(cmd.Actor)
	if err != nil {
		return xerrors.Errorf("failed to count invite: %w", err)
	}

	logrus.Infof("invite count: %d max: %d", c, org.config.MaxInviteCount)

	if c >= org.config.MaxInviteCount {
		return allerror.NewCountExceeded("exceed max invite count",
			xerrors.Errorf("invite count(now:%d max:%d) exceed", c, org.config.MaxInviteCount))
	}

	return nil
}

// HasMember returns true if the user is already a member of the organization.
func (org *orgService) HasMember(ctx context.Context, o, user primitive.Account) bool {
	_, err := org.member.GetByOrgAndUser(ctx, o.Account(), user.Account())
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		logrus.Errorf("failed to get member when check existence by org:%s and user:%s, %s",
			o.Account(), user.Account(), err)
		return true
	}

	if err == nil {
		logrus.Warnf("the user %s is already a member of the org %s", user.Account(), o.Account())
		return true
	}

	return false
}

func (org *orgService) IsInvite(ctx context.Context, o, user primitive.Account) bool {
	data, err := org.invite.GetInvite(user.Account(), o.Account())
	if err != nil {
		return false
	}
	if data.Status == domain.ApproveStatusPending {
		return true
	}
	return false
}

// RequestMember sends a membership request to join an organization.
func (org *orgService) RequestMember(
	ctx context.Context, cmd *domain.OrgRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		e := fmt.Errorf("invalid param for request member")
		err = allerror.NewInvalidParam(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.HasMember(ctx, cmd.Org, cmd.Actor) {
		e := fmt.Errorf(" user %s is already a member of the org %s", cmd.Actor.Account(), cmd.Org.Account())
		err = allerror.New(allerror.ErrorUserAccountIsAlreadyAMemberOfOrgAccount, e.Error(), e)
		return
	}

	requester, err := org.repo.GetByAccount(ctx, cmd.Actor)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "requester not found", err)
		}

		return

	}

	o, err := org.repo.GetOrgByName(ctx, cmd.Org)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "organization not found", err)
		}

		return
	}

	if !o.AllowRequest {
		e := fmt.Errorf("org not allow request member")
		err = allerror.New(allerror.ErrorOrgNotAllowRequestMember, e.Error(), e)
		return
	}

	if utf8.RuneCountInString(cmd.Msg) > int(org.config.MaxStrSize) {
		e := fmt.Errorf("character length exceeds limit")
		err = allerror.New(allerror.ErrorOrgNotAllowRequestMember, e.Error(), e)
		return
	}

	request := cmd.ToMemberRequest(o.DefaultRole)
	request.OrgId = o.Id
	request.UserId = requester.Id

	approve, isUPdate, err := org.invite.AddRequest(request)
	if err != nil {
		return
	}

	if !isUPdate {
		event := domain.NewOrgRequestEvent(cmd.Org, cmd.Actor)
		if err = org.message.SendOrgRequestEvent(event); err != nil {
			logrus.Errorf("send org request message failed: %s", err.Error())
		}
	}

	dto = ToMemberRequestDTO(ctx, &approve, org.user)

	return
}

// AcceptInvite accept the invitation the admin sent to me
func (org *orgService) AcceptInvite(ctx context.Context, cmd *domain.OrgAcceptInviteCmd) (dto ApproveDTO, err error) {
	if cmd == nil {
		e := fmt.Errorf("invalid param for cancel request member")
		err = allerror.NewInvalidParam(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.HasMember(ctx, cmd.Org, cmd.Actor) {
		e := fmt.Errorf("the user %s is already a member of the org %s", cmd.Actor.Account(), cmd.Org.Account())
		err = allerror.New(allerror.ErrorUserAccountIsAlreadyAMemberOfOrgAccount, e.Error(), e)
		return
	}

	// list all invitations sent to myself in the org
	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Invitee: cmd.Actor,
		OrgNormalCmd: domain.OrgNormalCmd{
			Org: cmd.Org,
		},
		Status: domain.ApproveStatusPending,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("the %s's invitation to org %s not found",
				cmd.Actor.Account(), cmd.Org.Account()), err)
		}

		return
	}

	if len(o) == 0 {
		err = allerror.New(allerror.ErrorNoInvitationFound, "", errors.New("no invitation found"))

		return
	}

	approve := o[0]

	if cmd.Actor.Account() != approve.Username.Account() {
		e := fmt.Errorf("can't accept other's invitation")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	// should change invitation's status when invitation has expired
	if approve.ExpireAt < utils.Now() {
		approve.Status = domain.ApproveStatusRejected
		_, err = org.invite.SaveInvite(&approve)
		if err != nil {
			err = allerror.NewExpired("failed to update expired invitation status", err)
			return
		}
		e := fmt.Errorf("the invitation has expired")
		err = allerror.NewExpired(e.Error(), e)
		return
	}

	owners, err := org.getOwners(cmd.Org)
	if err != nil {
		e := fmt.Errorf("failed to get owners of org when add new member: %s err: %w", cmd.Org.Account(), err)
		err = allerror.NewInvalidParam(e.Error(), e)
		return
	}

	approve.By = cmd.Actor.Account()
	approve.Status = domain.ApproveStatusApproved
	approve.Msg = cmd.Msg

	invite, err := org.invite.SaveInvite(&approve)
	if err != nil {
		return
	}

	err = org.AddMember(ctx, &domain.OrgAddMemberCmd{
		Actor:  owners[0].Username,
		Org:    cmd.Org,
		OrgId:  approve.OrgId,
		User:   cmd.Actor,
		UserId: approve.UserId,
		Role:   approve.Role,
		Type:   domain.InviteTypeInvite,
	})

	if err != nil {
		return
	}

	dto = ToApproveDTO(ctx, &invite, org.user)

	// Update all requests and invites status pending to approved
	err = org.invite.UpdateAllApproveStatus(approve.Username, approve.OrgName, approve.Status)
	if err != nil {
		return
	}

	event := domain.NewUserJoinEventByInvite(approve.OrgName, approve.Username)
	err = org.message.SendComputilityUserJoinEvent(&event)
	if err != nil {
		e := xerrors.Errorf("send message to user join computility org failed: %w", err)
		err = allerror.New(allerror.ErrorMsgPublishFailed, "", e)
	}

	return
}

// ApproveRequest approve the request from the user outside the org
func (org *orgService) ApproveRequest(
	ctx context.Context, cmd *domain.OrgApproveRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		e := fmt.Errorf("invalid param for cancel request member")
		err = allerror.NewInvalidParam(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if cmd.Actor.Account() == cmd.Requester.Account() {
		e := fmt.Errorf("can't approve request from yourself")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionWrite)
	if err != nil {
		return
	}

	reqs, err := org.invite.ListRequests(cmd.ToListReqCmd())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("the %s's member request to org %s not found",
				cmd.Requester.Account(), cmd.Org.Account()), err)
		}

		return
	}

	if len(reqs) > 1 {
		err = fmt.Errorf("multiple requests found")
		return
	}

	if len(reqs) == 0 {
		err = fmt.Errorf("no request found")
		return
	}

	request := reqs[0]
	request.By = cmd.Actor.Account()
	request.Status = domain.ApproveStatusApproved
	request.Msg = cmd.Msg
	request.Member = cmd.Member
	role, _ := primitive.NewRole(cmd.Member.Account())

	_, err = org.invite.SaveRequest(&request)
	if err != nil {
		return
	}

	err = org.AddMember(ctx, &domain.OrgAddMemberCmd{
		Actor:  cmd.Actor,
		Org:    cmd.Org,
		OrgId:  request.OrgId,
		User:   cmd.Requester,
		UserId: request.UserId,
		Type:   domain.InviteTypeRequest,
		Role:   role,
	})
	if err != nil {
		return
	}

	// Update all requests and invites status pending to approved
	err = org.invite.UpdateAllApproveStatus(request.Username, request.OrgName, request.Status)
	if err != nil {
		return
	}

	event := domain.NewUserJoinEventByRequest(cmd.Org, cmd.Requester)
	err = org.message.SendComputilityUserJoinEvent(&event)
	if err != nil {
		e := xerrors.Errorf("send message to user join computility org failed: %w", err)
		err = allerror.New(allerror.ErrorMsgPublishFailed, "", e)
	}

	return
}

// CancelReqMember cancels a member request in an organization.
func (org *orgService) CancelReqMember(
	ctx context.Context, cmd *domain.OrgCancelRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		e := fmt.Errorf("invalid param for cancel request member")
		err = allerror.NewInvalidParam(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}
	// user can cancel the request by self
	// or admin can reject the request
	if cmd.Actor.Account() != cmd.Requester.Account() {
		err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionDelete)
		if err != nil {
			return
		}
	}

	o, err := org.invite.ListRequests(cmd.ToListReqCmd())
	if err != nil {
		return
	}

	if len(o) == 0 {
		err = fmt.Errorf("no request found")
		return
	}

	approve := o[0]
	approve.By = cmd.Actor.Account()
	approve.Status = domain.ApproveStatusRejected
	approve.Msg = cmd.Msg

	updatedRequest, err := org.invite.SaveRequest(&approve)
	if err != nil {
		err = fmt.Errorf("failed to remove invite: %w", err)
		return
	}

	if cmd.Actor.Account() != cmd.Requester.Account() {
		event := domain.NewOrgRequestRejectEvent(cmd.Org, cmd.Requester)
		if err = org.message.SendOrgRequestRejectEvent(&event); err != nil {
			logrus.Errorf("send org request reject message failed: %s", err.Error())
		}
	}

	dto = ToMemberRequestDTO(ctx, &updatedRequest, org.user)

	return
}

func (org *orgService) GetOnlyApply(userName string, orgName string) ([]domain.MemberRequest, error) {
	res, err := org.invite.GetOneApply(userName, orgName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListMemberReq lists the member requests for an organization.
func (org *orgService) ListMemberReq(
	ctx context.Context, cmd *domain.OrgMemberReqListCmd) (dtos MemberPagnationDTO, err error) {
	if cmd == nil {
		e := fmt.Errorf("invalid param for list member request")
		err = allerror.NewInvalidParam(e.Error(), e)
		return
	}

	if cmd.Actor == nil {
		e := fmt.Errorf("anno can not list requests")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	// 只有管理员可以查询组织内的申请
	if cmd.Org != nil {
		err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionRead)
		if err != nil {
			return
		}
	}

	// 不能列出其他人发出的申请
	if cmd.Requester != nil && cmd.Org == nil && cmd.Actor.Account() != cmd.Requester.Account() {
		e := fmt.Errorf("can't list requests from other people")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	reqs, total, err := org.invite.ListPagnation(cmd)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}
	result := make([]MemberRequestDTO, len(reqs))
	for i := range reqs {
		result[i] = ToMemberRequestDTO(ctx, &reqs[i], org.user)
	}
	dtos.Total = total
	dtos.Members = result
	return
}

// RevokeInvite revokes an organization invite.
func (org *orgService) RevokeInvite(ctx context.Context, cmd *domain.OrgRemoveInviteCmd) (dto ApproveDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}
	// user can revoke the invite by self
	// or admin can revoke the invite
	if cmd.Actor.Account() != cmd.Account.Account() {
		err = org.perm.Check(ctx, cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionDelete)
		if err != nil {
			return
		}
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		OrgNormalCmd: domain.OrgNormalCmd{
			Org:   cmd.Org,
			Actor: cmd.Actor,
		},
		Invitee: cmd.Account,
		Status:  domain.ApproveStatusPending,
	})

	if err != nil {
		return
	}

	if len(o) == 0 {
		err = allerror.New(allerror.ErrorNoInvitationFound, "", errors.New("no invitation found"))

		return
	}

	approve := o[0]
	approve.By = cmd.Actor.Account()
	approve.Status = domain.ApproveStatusRejected
	approve.Msg = cmd.Msg

	updatedInvite, err := org.invite.SaveInvite(&approve)
	if err != nil {
		err = fmt.Errorf("failed to remove invite: %w", err)
		return
	}

	dto = ToApproveDTO(ctx, &updatedInvite, org.user)

	return
}

// ListInvitationByOrg lists the invitations based on the org.
func (org *orgService) ListInvitationByOrg(ctx context.Context, actor, orgName primitive.Account,
	status domain.ApproveStatus) (dtos []ApproveDTO, err error) {
	if _, err = org.repo.GetOrgByName(ctx, orgName); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", orgName.Account()), err)
		}

		return
	}

	// permission check
	// check role when list invitations in a org
	err = org.perm.Check(ctx, actor, orgName, primitive.ObjTypeInvite, primitive.ActionRead)
	if err != nil {
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		OrgNormalCmd: domain.OrgNormalCmd{
			Org: orgName,
		},
		Status: status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", orgName.Account()), err)
		}

		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(ctx, &o[i], org.user)
	}

	return
}

// ListInvitationByInviter lists the invitations based on the inviter.
func (org *orgService) ListInvitationByInviter(ctx context.Context, actor, inviter primitive.Account,
	status domain.ApproveStatus) (dtos []ApproveDTO, err error) {
	// can't list other's sent invitations
	if inviter != actor {
		e := fmt.Errorf("can not list invitation sent by other user")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Inviter: inviter,
		Status:  status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", inviter), err)
		}

		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(ctx, &o[i], org.user)
	}

	return
}

// ListInvitationByInvitee lists the invitations based on the invitee.
func (org *orgService) ListInvitationByInvitee(ctx context.Context, actor, invitee primitive.Account,
	status domain.ApproveStatus) (dtos []ApproveDTO, err error) {
	// permission check
	// can't list other's received invitations
	if invitee != actor {
		e := fmt.Errorf("can not list invitation received by other user")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Invitee: invitee,
		Status:  status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", invitee), err)
		}

		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(ctx, &o[i], org.user)
	}

	return
}

// CheckName checks if the given name exists in the repository.
func (org *orgService) CheckName(ctx context.Context, name primitive.Account) bool {
	if name == nil {
		logrus.Error("name is nil")
		return false
	}

	return org.repo.CheckName(ctx, name)
}

// GetMemberByUserAndOrg retrieves the member information for a given user and organization.
func (org *orgService) GetMemberByUserAndOrg(
	ctx context.Context, u primitive.Account, o primitive.Account) (member MemberDTO, err error) {
	if u == nil {
		err = fmt.Errorf("user is nil")
		return
	}

	if o == nil {
		err = fmt.Errorf("org is nil")
		return
	}

	m, err := org.member.GetByOrgAndUser(ctx, o.Account(), u.Account())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s with user %s not found", o.Account(), u.Account()), err)
		}

		return
	}

	member = ToMemberDTO(&m)

	return
}
