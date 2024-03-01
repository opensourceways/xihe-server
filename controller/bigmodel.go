package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/bigmodel/app"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	chBufferSize = 1000
)

func AddRouterForBigModelController(
	rg *gin.RouterGroup,
	s app.BigModelService,
	us userapp.RegService,
) {
	ctl := BigModelController{
		s:  s,
		us: us,
	}

	// luojia
	rg.POST("/v1/bigmodel/luojia_upload_picture", ctl.LuoJiaUploadPicture)
	rg.POST("/v1/bigmodel/luojia", ctl.LuoJia)
	rg.GET("/v1/bigmodel/luojia", ctl.ListLuoJiaRecord)

	// wukong
	rg.POST("/v1/bigmodel/wukong", ctl.WuKong)
	rg.POST("/v1/bigmodel/wukong_async", ctl.WuKongAsync)
	rg.GET("/v1/bigmodel/wukong/rank", ctl.WuKongRank)
	rg.GET("/v1/bigmodel/wukong/task", ctl.WuKongLastFinisedTask)
	rg.POST("/v1/bigmodel/wukong/like", ctl.AddLike)
	rg.POST("/v1/bigmodel/wukong/public", ctl.AddPublic)
	rg.GET("/v1/bigmodel/wukong/public", ctl.ListPublic)
	rg.GET("/v1/bigmodel/wukong/publics", ctl.GetPublicsGlobal)
	rg.PUT("/v1/bigmodel/wukong/link", ctl.GenDownloadURL)
	rg.DELETE("/v1/bigmodel/wukong/like/:id", ctl.CancelLike)
	rg.DELETE("/v1/bigmodel/wukong/public/:id", ctl.CancelPublic)
	rg.GET("/v1/bigmodel/wukong/samples/:batch", ctl.GenWuKongSamples)
	rg.GET("/v1/bigmodel/wukong", ctl.ListLike)
	rg.POST("/v1/bigmodel/wukong/digg", ctl.AddDigg)
	rg.DELETE("/v1/bigmodel/wukong/digg", ctl.CancelDigg)

	// others
	rg.POST("/v1/bigmodel/ai_detector", ctl.AIDetector)
	rg.POST("/v1/bigmodel/baichuan2_7b_chat", ctl.BaiChuan)
	rg.POST("/v1/bigmodel/glm2_6b", ctl.GLM2)
	rg.POST("/v1/bigmodel/llama2_7b", ctl.LLAMA2)
	rg.POST("/v1/bigmodel/skywork_13b", ctl.SkyWork)
	rg.POST("/v1/bigmodel/iflytekspark", ctl.IFlytekSpark)

	// api apply
	rg.POST("/v1/bigmodel/api/apply/:model", ctl.ApplyApi)
	rg.GET("/v1/bigmodel/api/get", ctl.GetUserApplyRecord)
	rg.GET("/v1/bigmodel/api/apply/:model", ctl.IsApplied)
	rg.POST("/v1/bigmodel/api/wukong", ctl.WukongAPI)
	rg.GET("/v1/bigmodel/apiinfo/get/:model", ctl.GetApiInfo)
	rg.GET("/v1/bigmodel/api/refresh/:model", ctl.RefreshApiToken)
}

type BigModelController struct {
	baseController

	s  app.BigModelService
	us userapp.RegService
}

// @Title			LuoJia
// @Description	luo-jia big model
// @Tags			BigModel
// @Accept			json
// @Success		201	{object}		luojiaResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/luojia [post]
func (ctl *BigModelController) LuoJia(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "launch LuoJia bigmodel")

	if v, err := ctl.s.LuoJia(pl.DomainAccount()); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, luojiaResp{v})
	}
}

// @Title			ListLuoJiaRecord
// @Description	list luo-jia big model records
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.LuoJiaRecordDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/luojia [get]
func (ctl *BigModelController) ListLuoJiaRecord(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, err := ctl.s.ListLuoJiaRecord(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(v))
	}
}

// @Title			LuoJiaUploadPicture
// @Description	upload a picture for luo-jia
// @Tags			BigModel
// @Param			picture	formData	file	true	"picture"
// @Accept			json
// @Success		201	{object}		pictureUploadResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/luojia_upload_picture [post]
func (ctl *BigModelController) LuoJiaUploadPicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "upload picture to LuoJia bigmodel")

	f, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	p, err := f.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get picture",
		))

		return
	}

	defer p.Close()

	if !utils.IsPictureName(f.Filename) {
		ctl.sendBadRequestParamWithMsg(ctx, "image format not allowed")

		return
	}

	if err := ctl.s.LuoJiaUploadPicture(p, pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(pictureUploadResp{f.Filename}))
	}
}

// @Title			GenWuKongSamples
// @Description	gen wukong samples
// @Tags			BigModel
// @Param			batch	path	int	true	"batch num"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/samples/{batch} [get]
func (ctl *BigModelController) GenWuKongSamples(ctx *gin.Context) {
	if _, _, ok := ctl.checkUserApiToken(ctx, false); !ok {
		return
	}

	i, err := strconv.Atoi(ctx.Param("batch"))
	if err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.GenWuKongSamples(i); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			WuKong
// @Description	generates pictures by WuKong
// @Tags			BigModel
// @Param			body	body	wukongRequest	true	"body of wukong"
// @Accept			json
// @Success		201	{object}		wukongPicturesGenerateResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong [post]
func (ctl *BigModelController) WuKong(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "launch wukong bigmodel")

	req := wukongRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, code, err := ctl.s.WuKong(pl.DomainAccount(), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, wukongPicturesGenerateResp{v})
	}
}

// @Title			WuKong
// @Description	send async wukong request task
// @Tags			BigModel
// @Param			body	body	wukongRequest	true	"body of wukong"
// @Accept			json
// @Success		201	{object}		wukongPicturesGenerateResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong_async [post]
func (ctl *BigModelController) WuKongAsync(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "launch wukong bigmodel task")

	req := wukongRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if code, err := ctl.s.WuKongInferenceAsync(pl.DomainAccount(), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		utils.DoLog("", pl.Account, "create wukong picture generate task",
			fmt.Sprintf("imageQuantity: %d", req.ImgQuantity), "success")

		ctl.sendRespOfPost(ctx, "")
	}
}

// @Title			WuKong
// @Description	get wukong rank
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.WuKongRankDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/rank [get]
func (ctl *BigModelController) WuKongRank(ctx *gin.Context) {
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{csrftoken},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == csrftoken
		},
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//TODO delete
		log.Errorf("update ws failed, err:%s", err.Error())

		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	for i := 0; i < apiConfig.PodTimeout; i++ {
		dto, err := ctl.s.GetWuKongWaitingTaskRank(pl.DomainAccount())
		if err != nil {
			_ = ws.WriteJSON(newResponseError(err))

			log.Errorf("get rank failed: get status, err:%s", err.Error())

			return
		} else {
			_ = ws.WriteJSON(newResponseData(dto))
		}

		log.Debugf("info dto:%v", dto)

		if dto.Rank == 0 {
			log.Debug("task done")

			return
		}

		time.Sleep(time.Second)
	}
}

// @Title			WuKong
// @Description	get last finished task
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.WuKongRankDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/task [get]
func (ctl *BigModelController) WuKongLastFinisedTask(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, code, err := ctl.s.GetWuKongLastTaskResp(pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			AddLike
// @Description	add like to wukong picture
// @Tags			BigModel
// @Accept			json
// @Success		202	{object}		wukongAddLikeResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/like [post]
func (ctl *BigModelController) AddLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "add like to picture generated by wukong bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	reqTemp := wukongAddLikeFromTempRequest{}
	reqPublic := wukongAddLikeFromPublicRequest{}

	errTemp := ctx.ShouldBindBodyWith(&reqTemp, binding.JSON)
	errPublic := ctx.ShouldBindBodyWith(&reqPublic, binding.JSON)
	if errTemp != nil && errPublic != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)
		return
	}

	if errTemp == nil {
		cmd, err := reqTemp.toCmd(pl.DomainAccount())
		if err != nil {
			ctl.sendBadRequestParam(ctx, err)

			return
		}

		if pid, code, err := ctl.s.AddLikeFromTempPicture(&cmd); err != nil {
			ctl.sendCodeMessage(ctx, code, err)
		} else {
			ctl.sendRespOfPost(ctx, wukongAddLikeResp{pid})
		}
	}

	if errPublic == nil {
		cmd, err := reqPublic.toCmd(pl.DomainAccount())
		if err != nil {
			ctl.sendBadRequestParam(ctx, err)

			return
		}

		if pid, code, err := ctl.s.AddLikeFromPublicPicture(&cmd); err != nil {
			ctl.sendCodeMessage(ctx, code, err)
		} else {
			ctl.sendRespOfPost(ctx, wukongAddLikeResp{pid})
		}
	}
}

// @Title			CancelLike
// @Description	cancel like on wukong picture
// @Tags			BigModel
// @Param			id	path	string	true	"picture id"
// @Accept			json
// @Success		204
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/like/{id} [delete]
func (ctl *BigModelController) CancelLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "cancel like on wukong picture")

	err := ctl.s.CancelLike(
		pl.DomainAccount(), ctx.Param("id"),
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {

		utils.DoLog("", pl.Account, "delete wukong like picture",
			fmt.Sprintf("pictureid: %s", ctx.Param("id")), "success")

		ctl.sendRespOfDelete(ctx)
	}
}

// @Title			CancelPublic
// @Description	cancel public on wukong picture
// @Tags			BigModel
// @Param			id	path	string	true	"picture id"
// @Accept			json
// @Success		204
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/public/{id} [delete]
func (ctl *BigModelController) CancelPublic(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "stop publicizing picture generated by wukong bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	err := ctl.s.CancelPublic(
		pl.DomainAccount(), ctx.Param("id"),
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {

		utils.DoLog("", pl.Account, "delete wukong public picture",
			fmt.Sprintf("pictureid: %s", ctx.Param("id")), "success")

		ctl.sendRespOfDelete(ctx)
	}
}

// @Title			ListLike
// @Description	list wukong pictures user liked
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.WuKongLikeDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong [get]
func (ctl *BigModelController) ListLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, err := ctl.s.ListLikes(pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			AddDigg
// @Description	add digg to wukong picture
// @Tags			BigModel
// @Accept			json
// @Success		202	{object}		wukongDiggResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/digg [post]
func (ctl *BigModelController) AddDigg(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "add digg to picture generated by wukong bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := wukongAddDiggPublicRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)
		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		return
	}

	if count, err := ctl.s.DiggPicture(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, wukongDiggResp{count})
	}
}

// @Title			CancelDigg
// @Description	delete digg to wukong picture
// @Tags			BigModel
// @Param			body	body	wukongCancelDiggPublicRequest	true	"body of wukong"
// @Accept			json
// @Success		202	{object}		wukongDiggResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/digg [delete]
func (ctl *BigModelController) CancelDigg(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "stop digging to picture generated by wukong bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := wukongCancelDiggPublicRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)
		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		return
	}

	if count, err := ctl.s.CancelDiggPicture(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, wukongDiggResp{count})
	}
}

// @Title			AddPublic
// @Description	add public to wukong picture
// @Tags			BigModel
// @Accept			json
// @Success		202	{object}		wukongAddPublicResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/public [post]
func (ctl *BigModelController) AddPublic(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "publicized picture generated by wukong bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	reqTemp := wukongAddPublicFromTempRequest{}
	reqPublic := wukongAddPublicFromLikeRequest{}

	errTemp := ctx.ShouldBindBodyWith(&reqTemp, binding.JSON)
	errPublic := ctx.ShouldBindBodyWith(&reqPublic, binding.JSON)
	if errTemp != nil && errPublic != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	if errTemp == nil {
		cmd, err := reqTemp.toCmd(pl.DomainAccount())
		if err != nil {
			ctl.sendBadRequestParam(ctx, err)

			return
		}

		if pid, code, err := ctl.s.AddPublicFromTempPicture(&cmd); err != nil {
			ctl.sendCodeMessage(ctx, code, err)
		} else {
			ctl.sendRespOfPost(ctx, wukongAddPublicResp{pid})
		}

		return
	}

	if errPublic == nil {
		cmd := reqPublic.toCmd(pl.DomainAccount())

		if pid, code, err := ctl.s.AddPublicFromLikePicture(&cmd); err != nil {
			ctl.sendCodeMessage(ctx, code, err)
		} else {
			ctl.sendRespOfPost(ctx, wukongAddPublicResp{pid})
		}

		return
	}
}

// @Title			ListPublic
// @Description	list wukong pictures user publiced
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.WuKongPublicDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/public [get]
func (ctl *BigModelController) ListPublic(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, err := ctl.s.ListPublics(pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			GetPublicGlobal
// @Description	list all wukong pictures publiced
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.WuKongPublicDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/publics [get]
func (ctl *BigModelController) GetPublicsGlobal(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd := app.WuKongListPublicGlobalCmd{}

	f := func() (err error) {
		if v := ctl.getQueryParameter(ctx, "count_per_page"); v != "" {
			if cmd.CountPerPage, err = strconv.Atoi(v); err != nil {
				return
			}
			if cmd.CountPerPage > 100 || cmd.CountPerPage <= 0 {
				err = errors.New("bad count_per_page")
				return
			}
		}

		if v := ctl.getQueryParameter(ctx, "page_num"); v != "" {
			if cmd.PageNum, err = strconv.Atoi(v); err != nil {
				return
			}
		}

		if v := ctl.getQueryParameter(ctx, "level"); v != "" {
			cmd.Level = domain.NewWuKongPictureLevel(v)
		}

		cmd.User = pl.DomainAccount()

		return
	}

	if err := f(); err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err := cmd.Validate(); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, err := ctl.s.GetPublicsGlobal(&cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			GenDownloadURL
// @Description	generate download url of wukong picture
// @Tags			BigModel
// @Param			body	body	wukongPictureLink	true	"body of wukong"
// @Accept			json
// @Success		202	{object}		wukongPictureLink
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/wukong/link [put]
func (ctl *BigModelController) GenDownloadURL(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "generate download url of wukong picture"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := wukongPictureLink{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	link, code, err := ctl.s.ReGenerateDownloadURL(
		pl.DomainAccount(), req.Link,
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, wukongPictureLink{link})
	}
}

// @Title			AI Detector
// @Description	detecte if text generate by ai
// @Tags			BigModel
// @Param			body	body	aiDetectorReq	true	"body of ai detector"
// @Accept			json
// @Success		202	{object}		aiDetectorResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/ai_detector [post]
func (ctl *BigModelController) AIDetector(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "launch AI Detector bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := aiDetectorReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	code, ismachine, err := ctl.s.AIDetector(&cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, aiDetectorResp{ismachine})
	}
}

// @Title			BaiChuan
// @Description	conversational AI
// @Tags			BigModel
// @Param			body	body	baichuanReq	true	"body of baichuan"
// @Accept			json
// @Success		202	{object}		app.BaiChuanDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/baichuan2_7b_chat [post]
func (ctl *BigModelController) BaiChuan(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "launch baichuan bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := baichuanReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	code, dto, err := ctl.s.BaiChuan(&cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, dto)
	}
}

// @Title			GLM
// @Description	conversational AI
// @Tags			BigModel
// @Param			body	body	glm2Request	true	"body of glm2"
// @Accept			json
// @Success		202	{object}		string
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/glm2_6b [post]
func (ctl *BigModelController) GLM2(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "launch glm2 bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := glm2Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	ch := make(chan string, chBufferSize)
	cmd, err := req.toCmd(ch, pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	code, err := ctl.s.GLM2(&cmd)
	if err != nil {
		ctx.Stream(func(w io.Writer) bool {
			if code == app.ErrorBigModelRecourseBusy {
				ctx.SSEvent("message", "access overload, please try again later")
			} else if code == app.ErrorBigModelSensitiveInfo {
				ctx.SSEvent("message", "I cannot answer such questions")
			}

			close(ch)

			return false
		})

		return
	}

	ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-ch; ok {
			if msg == "done" {
				ctx.SSEvent("status", "done")
			} else {
				ctx.SSEvent("message", msg)
			}

			return true
		}

		return false
	})
}

// @Title			LLAMA2
// @Description	conversational AI
// @Tags			BigModel
// @Param			body	body	llama2Request	true	"body of llama2"
// @Accept			json
// @Success		202	{object}		string
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/llama2_7b [post]
func (ctl *BigModelController) LLAMA2(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "launch llama2 bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := llama2Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	ch := make(chan string, chBufferSize)
	cmd, err := req.toCmd(ch, pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	code, err := ctl.s.LLAMA2(&cmd)
	if err != nil {
		ctx.Stream(func(w io.Writer) bool {
			if code == app.ErrorBigModelRecourseBusy {
				ctx.SSEvent("message", "access overload, please try again later")
			} else if code == app.ErrorBigModelSensitiveInfo {
				ctx.SSEvent("message", "I cannot answer such questions")
			}

			close(ch)

			return false
		})

		return
	}

	ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-ch; ok {
			if msg == "done" {
				ctx.SSEvent("status", "done")
			} else {
				ctx.SSEvent("message", msg)
			}

			return true
		}

		return false
	})
}

// @Title			SkyWork
// @Description	conversational AI
// @Tags			BigModel
// @Param			body	body	skyWorkRequest	true	"body of skywork"
// @Accept			json
// @Success		202	{object}		string
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/skywork_13b [post]
func (ctl *BigModelController) SkyWork(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "launch skywork bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := skyWorkRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	ch := make(chan string)
	cmd, err := req.toCmd(ch, pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	code, err := ctl.s.SkyWork(&cmd)
	if err != nil {
		ctx.Stream(func(w io.Writer) bool {
			if code == app.ErrorBigModelRecourseBusy {
				ctx.SSEvent("message", "access overload, please try again later")
			} else if code == app.ErrorBigModelSensitiveInfo {
				ctx.SSEvent("message", "I cannot answer such questions")
			}

			close(ch)

			return false
		})

		return
	}

	ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-ch; ok {
			if msg == "done" {
				ctx.SSEvent("status", "done")
				close(ch)
			} else {
				ctx.SSEvent("message", msg)
			}
			return true
		}
		return false
	})
}

// @Title			IFlytekSpark
// @Description	conversational AI
// @Tags			BigModel
// @Param			body	body	iflyteksparkRequest	true	"body of iflytekspark"
// @Accept			json
// @Success		202	{object}		string
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/iflytekspark [post]
func (ctl *BigModelController) IFlytekSpark(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	desc := "launch iflytekspark bigmodel"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	req := iflyteksparkRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	ch := make(chan string, chBufferSize)
	cmd, err := req.toCmd(ch, pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	code, err := ctl.s.IFlytekSpark(&cmd)
	if err != nil {
		ctx.Stream(func(w io.Writer) bool {
			if code == app.ErrorBigModelRecourseBusy {
				ctx.SSEvent("message", "access overload, please try again later")
			} else if code == app.ErrorBigModelSensitiveInfo {
				ctx.SSEvent("message", "I cannot answer such questions")
			}

			close(ch)

			return false
		})

		return
	}

	ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-ch; ok {
			if msg == "done" {
				ctx.SSEvent("status", "done")
			} else {
				ctx.SSEvent("message", msg)
			}

			return true
		}

		return false
	})
}

// @Title			ApplyApi
// @Description	generates pictures by WuKong-hf
// @Tags			BigModel
// @Param			body	body	applyApiReq	true	"body of wukong"
// @Accept			json
// @Success		201	{object}				newApiTokenResp
// @Failure		500	system_error			system	error
// @Router			/v1/bigmodel/api/apply/{model} [post]
func (ctl *BigModelController) ApplyApi(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	desc := "apply wukong api"
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, desc)

	b, err := ctl.s.IsApplyModel(pl.DomainAccount(), model)
	if b || err != nil {
		ctl.sendBadRequestParamWithMsg(ctx, "already applied")
		return
	}

	req := applyApiReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	if !req.Agreement {
		ctl.sendBadRequestParamWithMsg(ctx, "do not sign the agreement")
		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	v := pl.Account + "+" + strconv.FormatInt(utils.Now(), 10)
	token, err := ctl.encryptData(v)
	if err != nil {
		return
	}
	enToken, err := ctl.encryptDataForToken(token)
	if err != nil {
		return
	}

	if err = ctl.us.UpsertUserRegInfo(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	if err := ctl.s.ApplyApi(pl.DomainAccount(), model, enToken); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Title			WuKongAPI
// @Description	generates pictures by WuKong
// @Tags			BigModel
// @Param			body	body	wukongRequest	true	"body of wukong"
// @Accept			json
// @Success		201	{object}		wukongPicturesGenerateResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/api/{model} [post]
func (ctl *BigModelController) WukongAPI(ctx *gin.Context) {
	user, ok := ctl.checkBigmodelApiToken(ctx)
	if !ok {
		return
	}
	ac, err := types.NewAccount(user)
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	desc := "launch wukong bigmodel by api"
	prepareOperateLog(ctx, user, OPERATE_TYPE_USER, desc)

	model, err := domain.NewModelName("wukong")
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	r, err := ctl.s.GetApplyRecordByModel(ac, model)
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	deToken, err := ctl.decryptDataForToken(r)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	if ctx.GetHeader(Token) != string(deToken) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "invalid token",
		))
		return
	}

	req := wukongApiRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	if v, code, err := ctl.s.WukongApi(ac, model, &cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, wukongPicturesGenerateResp{v})
	}
}

// @Title			GetUserApplyRecord
// @Description	get user apply record
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.ApiApplyRecordDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/api/get/ [get]
func (ctl *BigModelController) GetUserApplyRecord(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, err := ctl.s.GetApplyRecordByUser(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			IsApplied
// @Description	is user applied for api
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		isApplyResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/api/apply/{model} [get]
func (ctl *BigModelController) IsApplied(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	if b, err := ctl.s.IsApplyModel(pl.DomainAccount(), model); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, isApplyResp{b})
	}
}

// @Title			GetApiInfo
// @Description	get api info
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		app.ApiInfoDTO
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/apiinfo/get/{model} [get]
func (ctl *BigModelController) GetApiInfo(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}
	b, err := ctl.s.IsApplyModel(pl.DomainAccount(), model)

	if err != nil || !b {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	if v, err := ctl.s.GetApiInfo(model); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title			RefreshApiToken
// @Description	refresh api token
// @Tags			BigModel
// @Accept			json
// @Success		200	{object}		newApiTokenResp
// @Failure		500	system_error	system	error
// @Router			/v1/bigmodel/api/refresh/{model} [get]
func (ctl *BigModelController) RefreshApiToken(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
		return
	}

	v := pl.Account + "+" + strconv.FormatInt(utils.Now(), 10)
	newToken, err := ctl.encryptData(v)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
		return
	}
	enNewToken, err := ctl.encryptDataForToken(newToken)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
		return
	}

	if date, err := ctl.s.UpdateToken(pl.DomainAccount(), model, enNewToken); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, newApiTokenResp{newToken, date})
	}
}
