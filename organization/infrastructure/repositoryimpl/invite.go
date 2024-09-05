/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryimpl provides implementations of repository interfaces for the organization domain.
package repositoryimpl

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	postgresql "github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/organization/domain"
	"github.com/opensourceways/xihe-server/organization/domain/repository"
	"github.com/sirupsen/logrus"
)

// NewInviteRepo creates a new instance of inviteRepoImpl.
func NewInviteRepo(db postgresql.Impl) repository.Approve {
	if err := postgresql.DB().Table(db.TableName()).AutoMigrate(&Approve{}); err != nil {
		return nil
	}

	return &inviteRepoImpl{Impl: db}
}

type inviteRepoImpl struct {
	postgresql.Impl
}

// ListInvitation lists the invitations based on the provided command.
func (impl *inviteRepoImpl) ListInvitation(cmd *domain.OrgInvitationListCmd) (approves []domain.Approve, err error) {
	var v []Approve

	query := impl.DB().Where(impl.EqualQuery(fieldType), domain.InviteTypeInvite)

	if cmd.Org != nil {
		query = query.Where(impl.EqualQuery(fieldOrg), cmd.Org.Account())
	}

	if cmd.Invitee != nil {
		query = query.Where(impl.EqualQuery(fieldInvitee), cmd.Invitee.Account())
	}

	if cmd.Inviter != nil {
		query = query.Where(impl.EqualQuery(fieldInviter), cmd.Inviter.Account())
	}

	if cmd.Status != "" {
		query = query.Where(impl.EqualQuery(fieldStatus), cmd.Status)
	}

	err = query.Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, err
	}

	approves = make([]domain.Approve, len(v))
	for i := range v {
		approves[i] = toApprove(&v[i])
	}

	return
}

func (impl *inviteRepoImpl) Count(user primitive.Account) (c int64, err error) {
	query := impl.DB().Where(impl.EqualQuery(fieldInviter), user.Account())
	query = query.Where(impl.EqualQuery(fieldStatus), domain.ApproveStatusPending)

	err = query.Count(&c).Error
	if err != nil {
		return 0, xerrors.Errorf("failed to count: %w", err)
	}

	return
}

// AddInvite adds a new invite to the database.
func (impl *inviteRepoImpl) AddInvite(o *domain.Approve) (new domain.Approve, isUpdate bool, err error) {
	o.Id = primitive.CreateIdentity(primitive.GetId())
	do := toApproveDoc(o)

	// Update existed record or add a new record, keep only one approve record
	if isUpdate, err = impl.saveAndKeepOneApprove(&do); err != nil {
		return
	}

	new = toApprove(&do)
	return
}

// SaveInvite saves an existing invite in the database.
func (impl *inviteRepoImpl) SaveInvite(o *domain.Approve) (new domain.Approve, err error) {
	do := toApproveDoc(o)
	do.Version += 1

	tmpDo := &Approve{}
	tmpDo.ID = o.Id.Integer()
	v := impl.DB().Model(
		tmpDo,
	).Clauses(clause.Returning{}).Where(
		impl.EqualQuery(fieldVersion), o.Version,
	).Select(`*`).Omit("created_at").Updates(&do) // should not update created_at

	if v.Error != nil {
		err = v.Error
		return
	}

	if v.RowsAffected == 0 {
		err = commonrepo.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
		return
	}
	new = toApprove(tmpDo)

	return
}

// DeleteInviteAndReqByOrg deletes invite and request records associated with the given organization account.
func (impl *inviteRepoImpl) DeleteInviteAndReqByOrg(acc primitive.Account) error {
	return impl.DB().Where(impl.EqualQuery(fieldOrg), acc.Account()).Delete(&Approve{}).Error
}

// AddRequest adds a new member request and returns the created request.
func (impl *inviteRepoImpl) AddRequest(r *domain.MemberRequest) (new domain.MemberRequest, isUpdate bool, err error) {
	r.Id = primitive.CreateIdentity(primitive.GetId())
	do := toRequestDoc(r)

	// Update existed record or add a new record, keep only one approve record
	if isUpdate, err = impl.saveAndKeepOneApprove(&do); err != nil {
		return
	}

	new = toMemberRequest(&do)

	return
}

// SaveRequest updates an existing member request and returns the updated request.
func (impl *inviteRepoImpl) SaveRequest(r *domain.MemberRequest) (new domain.MemberRequest, err error) {
	do := toRequestDoc(r)
	do.Version += 1

	tmpDo := &Approve{}
	tmpDo.ID = r.Id.Integer()
	v := impl.DB().Model(
		tmpDo,
	).Clauses(clause.Returning{}).Where(
		impl.EqualQuery(fieldVersion), r.Version,
	).Select(`*`).Omit("created_at").Updates(&do) // should not update created_at

	if v.Error != nil {
		err = v.Error
		return
	}

	if v.RowsAffected == 0 {
		err = commonrepo.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
		return
	}

	return toMemberRequest(tmpDo), nil

}

// ListRequests lists member requests based on the provided command criteria.
func (impl *inviteRepoImpl) ListRequests(cmd *domain.OrgMemberReqListCmd) (rs []domain.MemberRequest, err error) {
	var v []Approve

	query := impl.DB().Where(impl.EqualQuery(fieldType), domain.InviteTypeRequest)

	if cmd.Org != nil {
		query = query.Where(impl.EqualQuery(fieldOrg), cmd.Org.Account())
	}

	if cmd.Requester != nil {
		query = query.Where(impl.EqualQuery(fieldInvitee), cmd.Requester.Account())
	}

	if cmd.Status != "" {
		query = query.Where(impl.EqualQuery(fieldStatus), cmd.Status)
	}

	err = query.Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, err
	}

	rs = make([]domain.MemberRequest, len(v))
	for i := range v {
		rs[i] = toMemberRequest(&v[i])
	}

	return
}

func (impl *inviteRepoImpl) GetOneApply(userName string, orgName string) (res []domain.MemberRequest, err error) {
	var v []Approve
	query := impl.DB().Where(impl.EqualQuery(fieldType), domain.InviteTypeRequest)
	query = query.Where(impl.EqualQuery(fieldOrg), orgName)
	query = query.Where(impl.EqualQuery(fieldUser), userName)
	logrus.Info(fmt.Sprintf("%s DESC", fieldCreatedAt))
	err = query.Order(fmt.Sprintf("%s DESC", fieldCreatedAt)).Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, err
	}
	res = make([]domain.MemberRequest, 1)
	res[0] = toMemberRequest(&v[0])
	return res, nil
}

func (impl *inviteRepoImpl) GetInvite(userName, orgName string) (res domain.MemberRequest, err error) {
	var v Approve
	query := impl.DB().Where(impl.EqualQuery(fieldOrg), orgName)
	query = query.Where(impl.EqualQuery(fieldUser), userName)
	err = query.Order(fmt.Sprintf("%s DESC", fieldUpdateAt)).Find(&v).Error
	if err != nil {
		return
	}
	res = toMemberRequest(&v)
	return
}

func (impl *inviteRepoImpl) ListPagnation(cmd *domain.OrgMemberReqListCmd) ([]domain.MemberRequest, int, error) {
	var v []Approve
	query := impl.DB().Where(impl.EqualQuery(fieldType), domain.InviteTypeRequest)
	if cmd.Org != nil {
		query = query.Where(impl.EqualQuery(fieldOrg), cmd.Org.Account())
	}

	if cmd.Requester != nil {
		query = query.Where(impl.EqualQuery(fieldInvitee), cmd.Requester.Account())
	}

	if cmd.Status != "" {
		query = query.Where(impl.EqualQuery(fieldStatus), cmd.Status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := 0
	if cmd.PageNum > 0 && cmd.PageSize > 0 {
		offset = (cmd.PageNum - 1) * cmd.PageSize
	}
	if offset > 0 {
		query = query.Limit(cmd.PageSize).Offset(offset)
	} else {
		query = query.Limit(cmd.PageSize)
	}
	err := query.Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, 0, err
	}

	res := make([]domain.MemberRequest, len(v))
	for i := range v {
		res[i] = toMemberRequest(&v[i])
	}
	return res, int(total), nil
}

// saveAndKeepOneApprove save a new member request or invite, keep only one record.
func (impl *inviteRepoImpl) saveAndKeepOneApprove(ap *Approve) (isUpdate bool, err error) {
	if ap == nil {
		err = errors.New("approve is nil")
		return
	}
	// Build query to check for an existing invitation
	query := impl.DB().Where(impl.EqualQuery(fieldUser), ap.Username).
		Where(impl.EqualQuery(fieldOrg), ap.Orgname).
		Where(impl.EqualQuery(fieldInviter), ap.Inviter).
		Where(impl.EqualQuery(fieldStatus), ap.Status).
		Where(impl.EqualQuery(fieldType), ap.Type)

	// Attempt to find an existing record
	var existingApprove Approve
	err = query.First(&existingApprove).Error
	if err == nil {
		// Found an existing record, update it
		ap.ID = existingApprove.ID
		err = impl.DB().Model(&Approve{CommonModel: postgresql.CommonModel{ID: existingApprove.ID}}).
			Clauses(clause.Returning{}).Updates(&ap).Error

		isUpdate = true

		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// An error occurred other than "record not found", return the error
		return
	}

	err = impl.DB().Clauses(clause.Returning{}).Create(&ap).Error

	return
}

// UpdateAllApproveStatus updates all requests and invites status with pending status.
func (impl *inviteRepoImpl) UpdateAllApproveStatus(user, org primitive.Account, status domain.ApproveStatus) error {
	return impl.DB().Model(&Approve{}).Clauses(clause.Returning{}).
		Where(impl.EqualQuery(fieldUser), user).
		Where(impl.EqualQuery(fieldOrg), org).
		Where(impl.EqualQuery(fieldStatus), domain.ApproveStatusPending).
		Update(fieldStatus, status).Error
}
