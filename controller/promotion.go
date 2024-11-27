package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/app"
	pc "github.com/opensourceways/xihe-server/promotion/controller"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

func AddRouterForPromotionController(
	rg *gin.RouterGroup,
	pros app.PromotionService,
	ps app.PointsService,
) {
	ctl := PromotionController{
		pros: pros,
		ps:   ps,
	}

	// promotion
	rg.GET("/v1/promotions", ctl.List)
	rg.GET("/v1/promotion/:id", ctl.Get)
	rg.POST("/v1/promotion/:id/apply", ctl.Apply)
	rg.GET("/v1/promotion/user/:account", ctl.GetUserRegistration)

	// user points
	rg.GET("/v1/promotion/:id/points/:account", ctl.GetUserPoints)
	rg.GET("/v1/promotion/:id/ranking", ctl.GetUserRanking)
}

type PromotionController struct {
	baseController

	pros app.PromotionService
	ps   app.PointsService
}

// @Summary		Apply
// @Description	apply the Promotion
// @Tags			Promotion
// @Param			id		path	string					true	"promotion id"
// @Param			body	body	pc.PromotionApplyReq	true	"body of applying"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/promotion/{id}/apply [post]
func (ctl *PromotionController) Apply(ctx *gin.Context) {
	req := pc.PromotionApplyReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "apply the promotion")

	cmd, err := req.ToCmd(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

		return
	}

	if code, err := ctl.pros.UserRegister(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary		GetUserRegistration
// @Description	get user registrater promotion
// @Tags			Promotion
// @Param			account	path	string	true	"username"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/promotion/user/{account} [get]
func (ctl *PromotionController) GetUserRegistration(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.Account != ctx.Param("account") {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(errors.New("cannot find this account")))

		return
	}

	u := pl.DomainAccount()
	if dto, code, err := ctl.pros.GetUserRegisterPromotion(u); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, dto)
	}
}

// @Summary		GetUserPoints
// @Description	get user points in promotion
// @Tags			Promotion
// @Param			account	path	string	true	"username"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/promotion/{promotion}/points/{account} [get]
func (ctl *PromotionController) GetUserPoints(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.Account != ctx.Param("account") {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(errors.New("cannot find this account")))

		return
	}

	lang, err := ctl.languageRuquested(ctx)
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.ps.GetPoints(
		&app.PointsCmd{
			Promotionid: ctx.Param("id"),
			User:        pl.DomainAccount(),
			Lang:        lang,
		},
	); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, dto)
	}
}

// @Summary		GetUserRanking
// @Description	get user points ranking in promotion
// @Tags			Promotion
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/promotion/{promotion}/ranking [get]
func (ctl *PromotionController) GetUserRanking(ctx *gin.Context) {
	if dto, err := ctl.ps.GetPointsRank(ctx.Param("id")); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, dto)
	}
}

type listPromotionsReq struct {
	Type     string `form:"type"      json:"type"`
	Status   string `form:"status"    json:"status"`
	Way      string `form:"way"       json:"way"`
	Tags     string `form:"tags"      json:"tags"`
	PageNo   int    `form:"page_no"   json:"page_no"   binding:"required,min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"required,min=1"`
}

func (req listPromotionsReq) toCmd(user types.Account) (*app.ListPromotionsCmd, error) {
	var err error

	cmd := &app.ListPromotionsCmd{
		User:     user,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}

	if cmd.Type, err = domain.NewPromotionType(req.Type); err != nil {
		return cmd, err
	}

	if cmd.Status, err = domain.NewPromotionStatus(req.Status); err != nil {
		return cmd, err
	}

	if cmd.Way, err = domain.NewPromotionWay(req.Way); err != nil {
		return cmd, err
	}

	if req.Tags != "" {
		cmd.Tags = strings.Split(req.Tags, ",")
	}

	return cmd, nil
}

// @Summary		List
// @Description	list promotions
// @Tags			Promotion
// @Accept			json
// @Success		200	{object}	app.PromotionsDTO
// @Failure		400	{object}	controller.responseData
// @Failure		500	{object}	controller.responseData
// @Router			/v1/promotions [get]
func (ctl *PromotionController) List(ctx *gin.Context) {
	req := listPromotionsReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	pl, vistor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "list promotions")

	var (
		cmd *app.ListPromotionsCmd
		err error
	)

	if vistor {
		cmd, err = req.toCmd(nil)
	} else {
		cmd, err = req.toCmd(pl.DomainAccount())
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))
		return
	}

	if promotionsDTO, err := ctl.pros.List(cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, promotionsDTO)
	}
}

// @Summary		Get
// @Description	get promotion
// @Tags			Promotion
// @Accept			json
// @Success		200	{object}	app.PromotionDTO
// @Failure		400	{object}	controller.responseData
// @Failure		404	{object}	controller.responseData
// @Failure		500	{object}	controller.responseData
// @Router			/v1/promotion/{id} [get]
func (ctl *PromotionController) Get(ctx *gin.Context) {
	pl, vistor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "get promotion")

	cmd := app.PromotionCmd{
		Id:   ctx.Param("id"),
		User: nil,
	}
	if !vistor {
		cmd.User = pl.DomainAccount()
	}

	if dto, err := ctl.pros.Get(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, dto)
	}
}
