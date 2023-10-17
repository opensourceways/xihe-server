package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
)

// @Title			Create
// @Description	add a following
// @Tags			Following
// @Param			body	body	followingCreateRequest	true	"body of creating following"
// @Accept			json
// @Success		201
// @Failure		400	bad_request_body	can't	parse		request		body
// @Failure		401	bad_request_param	some	parameter	of			body		is		invalid
// @Failure		402	not_allowed			can't	add			yourself	as			your	following
// @Failure		403	resource_not_exists	the		target		of			following	does	not	exist
// @Failure		404	duplicate_creating	add		following	again
// @Failure		500	system_error		system	error
// @Router			/v1/user/following [post]
func (ctl *UserController) AddFollowing(ctx *gin.Context) {
	req := followingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	user, err := domain.NewAccount(req.Account)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if pl.isMyself(user) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't add yourself as your following",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "add a following")

	if _, err = ctl.repo.GetByAccount(user); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	f := &userdomain.FollowerInfo{
		User:     user,
		Follower: pl.DomainAccount(),
	}

	if err := ctl.s.AddFollowing(f); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

// @Title			Delete
// @Description	remove a following
// @Tags			Following
// @Param			account	path	string	true	"the account of following"
// @Accept			json
// @Success		204
// @Failure		400	bad_request_param	invalid	account
// @Failure		401	not_allowed			can't	remove	yourself	from	your	following
// @Failure		500	system_error		system	error
// @Router			/v1/user/following/{account} [delete]
func (ctl *UserController) RemoveFollowing(ctx *gin.Context) {
	user, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if pl.isMyself(user) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "invalid operation",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "remove a following")

	f := &userdomain.FollowerInfo{
		User:     user,
		Follower: pl.DomainAccount(),
	}

	if err := ctl.s.RemoveFollowing(f); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Title			List
// @Description	list followings
// @Tags			Following
// @Param			account	path	string	true	"the account the followings belong to"
// @Accept			json
// @Success		200	{object}		app.FollowsDTO
// @Failure		500	system_error	system	error
// @Router			/v1/user/following/{account} [get]
func (ctl *UserController) ListFollowing(ctx *gin.Context) {
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd, ok := ctl.genListFollowsCmd(ctx, account)
	if !ok {
		return
	}

	if data, err := ctl.s.ListFollowing(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

func (ctl *UserController) genListFollowsCmd(
	ctx *gin.Context, user domain.Account,
) (cmd userapp.FollowsListCmd, ok bool) {
	var err error

	if v := ctl.getQueryParameter(ctx, "count_per_page"); v != "" {
		if cmd.CountPerPage, err = strconv.Atoi(v); err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	if v := ctl.getQueryParameter(ctx, "page_num"); v != "" {
		if cmd.PageNum, err = strconv.Atoi(v); err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if !visitor {
		cmd.Follower = pl.DomainAccount()
	}

	cmd.User = user

	return
}
