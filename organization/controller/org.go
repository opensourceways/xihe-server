/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	commonctl "github.com/opensourceways/xihe-server/common/controller"
	"github.com/opensourceways/xihe-server/common/controller/middleware"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	orgapp "github.com/opensourceways/xihe-server/organization/app"
	"github.com/opensourceways/xihe-server/organization/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userctl "github.com/opensourceways/xihe-server/user/controller"
)

// AddRouterForOrgController adds routes for organization-related operations to the given router group.
func AddRouterForOrgController(
	rg *gin.RouterGroup,
	org orgapp.OrgService,
	cert orgapp.OrgCertificateService,
	user userapp.UserService,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	m middleware.UserMiddleWare,
	rl middleware.RateLimiter,
	p middleware.PrivacyCheck,
	npuGatekeeper orgapp.PrivilegeOrg,
	disable orgapp.PrivilegeOrg,
) {
	ctl := OrgController{
		m:              m,
		org:            org,
		user:           user,
		npuGatekeeper:  npuGatekeeper,
		disable:        disable,
		orgCertificate: cert,
	}

	rg.PUT("/v1/organization/:name", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Update)
	rg.POST("/v1/organization", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Create)
	rg.GET("/v1/organization/:name", m.Optional, rl.CheckLimit, ctl.Get)
	rg.GET("/v1/organization", m.Optional, rl.CheckLimit, ctl.List)
	rg.POST("/v1/organization/:name", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Leave)
	rg.DELETE("/v1/organization/:name", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Delete)
	rg.HEAD("/v1/name", m.Read, rl.CheckLimit, ctl.Check)

	rg.POST("/v1/invite", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.InviteMember)
	rg.PUT("/v1/invite", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.AcceptInvite)
	rg.GET("/v1/invite", m.Read, rl.CheckLimit, ctl.ListInvitation)
	rg.DELETE("/v1/invite", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.RemoveInvitation)

	rg.POST("/v1/request", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.RequestMember)
	rg.PUT("/v1/request", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.ApproveRequest)
	rg.GET("/v1/request", m.Read, rl.CheckLimit, ctl.ListRequests)
	rg.GET("/v1/request/only/:username/:orgname", m.Optional, rl.CheckLimit, ctl.GetOnlyRequest)
	rg.DELETE("/v1/request", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.RemoveRequest)

	rg.DELETE("/v1/organization/:name/member", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.RemoveMember)
	rg.GET("/v1/organization/:name/member", m.Optional, rl.CheckLimit, ctl.ListMember)
	rg.PUT("/v1/organization/:name/member", m.Write,
		userctl.CheckMail(ctl.m, ctl.user, sl), l.Write, rl.CheckLimit, ctl.EditMember)

	rg.POST("/v1/organization/:name/certificate", m.Write, l.Write, rl.CheckLimit, ctl.Certificate)
	rg.GET("/v1/organization/:name/certificate", m.Optional, rl.CheckLimit, ctl.GetCertification)
	rg.GET("/v1/organization/:name/certificate/check", m.Write, rl.CheckLimit, ctl.CertificateCheck)

	rg.GET("/v1/account/:name", p.CheckName, m.Optional, rl.CheckLimit, ctl.GetUser)
	rg.GET("/v1/user/privilege", m.Read, ctl.GetPrivilege)
}

// OrgController is a struct that contains the necessary dependencies for organization-related operations.
type OrgController struct {
	m              middleware.UserMiddleWare
	org            orgapp.OrgService
	user           userapp.UserService
	npuGatekeeper  orgapp.PrivilegeOrg
	disable        orgapp.PrivilegeOrg
	orgCertificate orgapp.OrgCertificateService
}

// @Summary  Update
// @Description  update org basic info
// @Tags     Organization
// @Param    name  path  string                     true  "name" MaxLength(40)
// @Param    body  body  orgBasicInfoUpdateRequest  true  "body of new organization"
// @Accept   json
// @Security Bearer
// @Success  202  {object}  commonctl.ResponseData{data=userapp.UserDTO,msg=string,code=string}
// @Router   /v1/organization/{name} [put]
func (ctl *OrgController) Update(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("update basic info of %s", ctx.Param("name")))

	var req orgBasicInfoUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(user, ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if o, err := ctl.org.UpdateBasicInfo(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, o)
	}
}

// @Summary  Get
// @Description  get organization info
// @Tags     Organization
// @Param    name  path  string  true  "name" MaxLength(40)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=userapp.UserDTO,msg=string,code=string}
// @Router   /v1/organization/{name} [get]
func (ctl *OrgController) Get(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if o, err := ctl.org.GetByAccount(ctx.Request.Context(), orgName); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, o)
	}
}

// @Summary   User or organization info
// @Description  get organization or user info
// @Tags     Organization
// @Param    name  path  string  true  "name of the user of organization" MaxLength(40)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=userapp.UserDTO,msg=string,code=string}
// @Failure  404  "user not found"
// @Failure  400  {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Router   /v1/account/{name} [get]
func (ctl *OrgController) GetUser(ctx *gin.Context) {
	name, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.m.GetUser(ctx)

	if o, err := ctl.org.GetOrgOrUser(ctx.Request.Context(), user, name); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, o)
	}
}

// @Summary  Check
// @Description  Check the name is available
// @Tags     Organization
// @Param    name  query  string  true  "the name to be check whether it's usable" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  200  "name is valid"
// @Failure  409  "name is invalid"
// @Router   /v1/name [head]
func (ctl *OrgController) Check(ctx *gin.Context) {

	var req reqToCheckName

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	name, err := req.toAccount()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if ctl.org.CheckName(ctx.Request.Context(), name) {
		ctx.JSON(http.StatusOK, nil)
	} else {
		ctx.JSON(http.StatusConflict, nil)
	}
}

// @Summary  List
// @Description  get organization info
// @Tags     Organization
// @Param    owner     query  string  false  "filter by owner" MaxLength(40)
// @Param    username  query  string  false  "filter by username" MaxLength(40)
// @Param    roles     query  []string  false  "filter by roles" Enums(read, write,admin)
// @Param    page_num  query  int     false    "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Param    search    query string false   "filter by name or fullname" MaxLength(40)
// @Param    cert_type  query  string  false  "filter by org type" MaxLength(40)
// @Param    sort_by  query  string  false  "need sort" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=userapp.UserPaginationDTO,msg=string,code=string}
// @Failure  400  {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Router   /v1/organization [get]
func (ctl *OrgController) List(ctx *gin.Context) {
	var req orgListRequest

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}
	me := ctl.m.GetUser(ctx)
	cmd, err := req.toCmd()
	if req.Username != "" && cmd.User == nil {
		cmd.User = me
	}
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}
	if err := ctl.org.IsSortInit(ctx); err != nil {
		logrus.Errorf("sort failed: %v", err)
	}
	listOption := &orgapp.OrgListOptions{
		Owner:    cmd.Owner,
		Member:   cmd.User,
		Roles:    cmd.Roles,
		Page:     cmd.PageNUm,
		PageSize: cmd.CountPerPage,
		Search:   cmd.Search,
		OrgType:  cmd.OrgType,
	}
	var os []userapp.BaseUserDTO
	var total int
	if req.SortBy == "" {
		os, total, err = ctl.org.OrgGetNoSort(ctx.Request.Context(), listOption)
	} else if cmd.OrgType != nil {
		os, total, err = ctl.org.OrgByTypeInfoList(ctx.Request.Context(), listOption)
	} else {
		os, total, err = ctl.org.OrgInfoList(ctx.Request.Context(), listOption)
	}
	if err != nil {
		commonctl.SendError(ctx, err)
	}
	res := ctl.org.ResOrgInfoList(os, total)
	commonctl.SendRespOfGet(ctx, res)
}

// @Summary  Create
// @Description  create a new organization
// @Tags     Organization
// @Param    body  body  orgCreateRequest  true  "body of new organization"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData{data=userapp.UserDTO,msg=string,code=string}
// @Router   /v1/organization [post]
func (ctl *OrgController) Create(ctx *gin.Context) {
	middleware.SetAction(ctx, "create organization")

	var req orgCreateRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd.Owner = user

	o, err := ctl.org.Create(ctx.Request.Context(), &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, o)
	}
}

// @Summary  Delete
// @Description  delete a organization
// @Tags     Organization
// @Param    name  path  string  true  "name" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/organization/{name} [delete]
func (ctl *OrgController) Delete(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("delete organization %s", ctx.Param("name")))

	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	err = ctl.org.Delete(ctx.Request.Context(), &domain.OrgDeletedCmd{
		Actor: user,
		Name:  orgName,
	})
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  ListMember
// @Description  list organization members
// @Tags     Organization
// @Param    name  path  string  true  "name" MaxLength(40)
// @Param    username  query  string  false  "filter by username" MaxLength(40)
// @Param    role  query  string  false  "filter by role" Enums(read, write,admin)
// @Accept   json
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData{data=[]app.MemberDTO,msg=string,code=string}
// @Router   /v1/organization/{name}/member [get]
func (ctl *OrgController) ListMember(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var req orgListMemberRequest
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}
	cmd, err := req.toCmd(orgName)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if members, err := ctl.org.ListMember(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, members)
	}
}

// @Summary  EditMember
// @Description Edit a member to the organization's role
// @Tags     Organization
// @Param    body  body  OrgMemberEditRequest  true  "body of new member"
// @Param    name  path  string  true  "name" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  202 {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/organization/{name}/member [put]
func (ctl *OrgController) EditMember(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("edit member of %s", ctx.Param("name")))

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgMemberEditRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd(ctx.Param("name"), user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx,
		fmt.Sprintf("edit member %s to be %s of %s", req.User, req.Role, cmd.Org.Account()))

	if _, err = ctl.org.EditMember(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  RemoveMember
// @Description Remove a member from a organization
// @Tags     Organization
// @Param    body  body  orgMemberRemoveRequest  true  "body of the removed member"
// @Param    name  path  string  true  "name" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/organization/{name}/member [delete]
func (ctl *OrgController) RemoveMember(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("remove member of %s", ctx.Param("name")))

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req orgMemberRemoveRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd(ctx.Param("name"), user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx,
		fmt.Sprintf("remove member %s from %s", req.User, cmd.Org.Account()))

	if err = ctl.org.RemoveMember(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Leave
// @Description  leave the organization
// @Tags     Organization
// @Param    name  path  string  true  "name" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/organization/{name} [post]
func (ctl *OrgController) Leave(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("leave organization of %s", ctx.Param("name")))

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("leave organization %s", orgName))

	err = ctl.org.RemoveMember(ctx.Request.Context(), &domain.OrgRemoveMemberCmd{
		Actor:   user,
		Org:     orgName,
		Account: user,
	})
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  InviteMember
// @Description Send invitation to a user to join the organization
// @Tags     Organization
// @Param    body  body  OrgInviteMemberRequest  true  "body of the invitation"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData{data=app.ApproveDTO,msg=string,code=string}
// @Router   /v1/invite [post]
func (ctl *OrgController) InviteMember(ctx *gin.Context) {
	middleware.SetAction(ctx, "invite member")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgInviteMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.org.InviteMember(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, dto)
	}
}

// @Summary  ListInvitation
// @Description List invitation of the organization
// @Tags     Organization
// @Param    org_name   query  string  false  "organization name" MaxLength(40)
// @Param    invitee    query  string  false  "invitee name" MaxLength(40)
// @Param    inviter    query  string  false  "inviter name" MaxLength(40)
// @Param    status     query  string  false  "invitation status, can be: pending/approved/rejected" Enums(pending, approved,rejected)
// @Param    page_size  query  int     false  "page size" Mininum(1)
// @Param    page       query  int     false  "page index" Mininum(1)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=[]app.ApproveDTO,msg=string,code=string}
// @Router   /v1/invite [get]
func (ctl *OrgController) ListInvitation(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgListInviteRequest

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd := req.toCmd(user)
	if err := cmd.Validate(); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var dtos []orgapp.ApproveDTO
	var err error
	if cmd.Org != nil {
		dtos, err = ctl.org.ListInvitationByOrg(ctx.Request.Context(), user, cmd.Org, cmd.Status)
	} else if cmd.Invitee != nil {
		dtos, err = ctl.org.ListInvitationByInvitee(ctx.Request.Context(), user, cmd.Invitee, cmd.Status)
	} else if cmd.Inviter != nil {
		dtos, err = ctl.org.ListInvitationByInviter(ctx.Request.Context(), user, cmd.Inviter, cmd.Status)
	}
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, dtos)
	}
}

// @Summary  RevokeMember
// @Description Revoke member request of the organization
// @Tags     Organization
// @Param    body  body  OrgRevokeMemberReqRequest  true  "body of the member request"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/request [delete]
func (ctl *OrgController) RemoveRequest(ctx *gin.Context) {
	middleware.SetAction(ctx, "remove member")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgRevokeMemberReqRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if _, err = ctl.org.CancelReqMember(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  RequestMember
// @Description Request to be a member of the organization
// @Tags     Organization
// @Param    body  body  OrgReqMemberRequest  true  "body of the member request"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData{data=orgapp.MemberRequestDTO,msg=string,code=string}
// @Router   /v1/request [post]
func (ctl *OrgController) RequestMember(ctx *gin.Context) {
	middleware.SetAction(ctx, "request to be a member")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgReqMemberRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.org.RequestMember(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, dto)
	}
}

// @Summary  ApproveRequest
// @Description Approve a user's member request of the organization
// @Tags     Organization
// @Param    body  body  OrgApproveMemberRequest  true  "body of the accept"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/request [put]
func (ctl *OrgController) ApproveRequest(ctx *gin.Context) {
	middleware.SetAction(ctx, "approve request to be a member")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgApproveMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if _, err = ctl.org.ApproveRequest(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  ListRequests
// @Description  List requests of the organization
// @Tags     Organization
// @Param    org_name   query  string  false  "organization name" MaxLength(40)
// @Param    requester  query  string  false  "invitee name" MaxLength(40)
// @Param    status     query  string  false  "invitation status, can be: pending/approved/rejected" Enums(pending, approved,rejected)
// @Param    page_size  query  int     false  "page size" Mininum(1)
// @Param    page       query  int     false  "page index" Mininum(1)
// @Accept   json
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData{data=app.MemberPagnationDTO,msg=string,code=string}
// @Router   /v1/request [get]
func (ctl *OrgController) ListRequests(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgListMemberReqRequest

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dtos, err := ctl.org.ListMemberReq(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, dtos)
	}
}

// @Summary GetOnlyRequest
// @Description Search one member Record
// @Tags Organization
// @Param username path string true "username" MaxLength(40)
// @param orgname path string true "orgname" MaxLength(40)
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData{data=[]domain.MemberRequest,msg=string,code=string}
// @Router   /v1/request/only/{username}/{orgname} [get]
func (ctl *OrgController) GetOnlyRequest(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}
	userName := ctx.Param("username")
	orgName := ctx.Param("orgname")
	if userName == "" || orgName == "" {
		e := fmt.Errorf("missing parameters")
		err := allerror.NewNoPermission(e.Error(), e)
		commonctl.SendBadRequestBody(ctx, err)
	}
	res, err := ctl.org.GetOnlyApply(userName, orgName)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, res)
	}
}

// @Summary  RevokeInvitation
// @Description Revoke invitation of the organization
// @Tags     Organization
// @Param    body  body  OrgRevokeInviteRequest  true  "body of the invitation"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/invite [delete]
func (ctl *OrgController) RemoveInvitation(ctx *gin.Context) {
	middleware.SetAction(ctx, "remove invite")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgRevokeInviteRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if _, err = ctl.org.RevokeInvite(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  AcceptInvite
// @Description Accept invite of the organization
// @Tags     Organization
// @Param    body  body  OrgAcceptMemberRequest  true  "body of the invitation"
// @Accept   json
// @Security Bearer
// @Success  202 {object}  commonctl.ResponseData{data=app.ApproveDTO,msg=string,code=string}
// @Router   /v1/invite [put]
func (ctl *OrgController) AcceptInvite(ctx *gin.Context) {
	middleware.SetAction(ctx, "accept invite")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgAcceptMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if a, err := ctl.org.AcceptInvite(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, a)
	}
}

// @Summary  List User's privilege organization info
// @Description List User's privilege organization info
// @Tags     Organization
// @Param    type   query  string  true  "privilege type, can be: npu/disable" Enums(npu, disable)
// @Param    user   query  string  false  "user name to filter the organizations which contain the user"
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData{data=[]userapp.UserDTO,msg=string,code=string}
// @Router   /v1/user/privilege [get]
func (ctl *OrgController) GetPrivilege(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req PrivilegeOption
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	// we don't need to check if user is nil
	user, _ = primitive.NewAccount(req.User)

	t, err := orgapp.NewAction(req.Type)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	opt := &orgapp.PrivilegeOrgListOptions{
		User: user,
		Type: t,
	}

	var orgs []userapp.UserDTO
	switch req.Type {
	case "npu":
		if ctl.npuGatekeeper != nil {
			orgs, err = ctl.npuGatekeeper.List(ctx.Request.Context(), opt)
		}
	case "disable":
		if ctl.disable != nil {
			orgs, err = ctl.disable.List(ctx.Request.Context(), opt)
		}
	default:
		e := fmt.Errorf("invalid privilege type: %s", req.Type)
		err = allerror.NewInvalidParam(e.Error(), e)
	}
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, orgs)
	}
}

// @Summary  Certification
// @Description organization certification
// @Tags     Organization
// @Param    name  path  string  true  "organization name" MaxLength(40)
// @Param    body  body orgCertificateRequest  true  "body of certificate"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/organization/{name}/certificate [post]
func (ctl *OrgController) Certificate(ctx *gin.Context) {
	middleware.SetAction(ctx, "organization certification")

	var req orgCertificateRequest
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(ctx.Param("name"), user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err = ctl.orgCertificate.Certificate(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}

// @Summary  Certification
// @Description get organization certification
// @Tags     Organization
// @Param    name  path  string  true  "name" MaxLength(40)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=app.OrgCertificateDTO,msg=string,code=string}
// @Router   /v1/organization/{name}/certificate [get]
func (ctl *OrgController) GetCertification(ctx *gin.Context) {
	middleware.SetAction(ctx, "get certification")

	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.m.GetUser(ctx)
	if v, err := ctl.orgCertificate.GetCertification(ctx.Request.Context(), orgName, user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// @Summary  Certification
// @Description check organization certification
// @Tags     Organization
// @Param    name  path  string  true  "name" MaxLength(40)
// @Param    certificate_org_name   query  string  false  "certificate organization name" MaxLength(100)
// @Param    unified_social_credit_code   query  string  false  "the unified social credit code" MaxLength(100)
// @Param    phone   query  string  false  "phone number" MaxLength(16)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=bool,msg=string,code=string}
// @Router   /v1/organization/{name}/certificate/check [get]
func (ctl *OrgController) CertificateCheck(ctx *gin.Context) {
	middleware.SetAction(ctx, "check certification")

	var req orgCertificateCheckRequest
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if b, err := ctl.orgCertificate.DuplicateCheck(ctx.Request.Context(), cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, b)
	}
}
