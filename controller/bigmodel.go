package controller

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/bigmodel/app"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func AddRouterForBigModelController(
	rg *gin.RouterGroup,
	s app.BigModelService,
) {
	ctl := BigModelController{
		s: s,
	}

	rg.POST("/v1/bigmodel/describe_picture", ctl.DescribePicture)
	rg.POST("/v1/bigmodel/describe_picture_hf", ctl.DescribePictureHF)
	rg.POST("/v1/bigmodel/single_picture", ctl.GenSinglePicture)
	rg.POST("/v1/bigmodel/multiple_pictures", ctl.GenMultiplePictures)
	rg.POST("/v1/bigmodel/vqa_upload_picture", ctl.VQAUploadPicture)
	rg.POST("/v1/bigmodel/luojia_upload_picture", ctl.LuoJiaUploadPicture)
	rg.POST("/v1/bigmodel/ask", ctl.Ask)
	rg.POST("/v1/bigmodel/ask_hf", ctl.AskHF)
	rg.POST("/v1/bigmodel/pangu", ctl.PanGu)
	rg.POST("/v1/bigmodel/luojia", ctl.LuoJia)
	rg.POST("/v1/bigmodel/luojia_hf", ctl.LuoJiaHF)
	rg.POST("/v1/bigmodel/codegeex", ctl.CodeGeex)
	rg.POST("/v1/bigmodel/wukong", ctl.WuKong)
	rg.POST("/v1/bigmodel/wukong_hf", ctl.WuKongHF)
	rg.POST("/v1/bigmodel/wukong_icbc", ctl.WuKongICBC)
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
	rg.GET("/v1/bigmodel/luojia", ctl.ListLuoJiaRecord)
	rg.POST("/v1/bigmodel/ai_detector", ctl.AIDetector)
}

type BigModelController struct {
	baseController

	s app.BigModelService
}

//	@Title			DescribePicture
//	@Description	describe a picture
//	@Tags			BigModel
//	@Param			picture	formData	file	true	"picture"
//	@Accept			json
//	@Success		201	{object}		describePictureResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/describe_picture [post]
func (ctl *BigModelController) DescribePicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	f, err := ctx.FormFile("picture")
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if f.Size > apiConfig.MaxPictureSizeToDescribe {
		ctl.sendBadRequestParamWithMsg(ctx, "too big picture")

		return
	}

	p, err := f.Open()
	if err != nil {
		ctl.sendBadRequestParamWithMsg(ctx, "can't get picture")

		return
	}

	defer p.Close()

	v, err := ctl.s.DescribePicture(pl.DomainAccount(), p, f.Filename, f.Size)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, describePictureResp{v})
	}
}

//	@Title			DescribePicture
//	@Description	describe a picture for hf
//	@Tags			BigModel
//	@Param			picture	formData	file	true	"picture"
//	@Accept			json
//	@Success		201	{object}		describePictureResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/describe_picture_hf [post]
func (ctl *BigModelController) DescribePictureHF(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	f, err := ctx.FormFile("picture")
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if f.Size > apiConfig.MaxPictureSizeToDescribe {
		ctl.sendBadRequestParamWithMsg(ctx, "too big picture")

		return
	}

	p, err := f.Open()
	if err != nil {
		ctl.sendBadRequestParamWithMsg(ctx, "can't get picture")

		return
	}

	defer p.Close()

	cmd := app.DescribePictureCmd{
		User:    pl.DomainAccount(),
		Picture: p,
		Name:    f.Filename,
		Length:  f.Size,
	}

	v, err := ctl.s.DescribePictureHF(&cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, describePictureResp{v})
	}
}

//	@Title			GenSinglePicture
//	@Description	generate a picture based on a text
//	@Tags			BigModel
//	@Param			body	body	pictureGenerateRequest	true	"body of generating picture"
//	@Accept			json
//	@Success		201	{object}		pictureGenerateResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/single_picture [post]
func (ctl *BigModelController) GenSinglePicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := pictureGenerateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if err := req.validate(); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, code, err := ctl.s.GenPicture(pl.DomainAccount(), req.Desc)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, pictureGenerateResp{v})
	}
}

//	@Title			GenMultiplePictures
//	@Description	generate multiple pictures based on a text
//	@Tags			BigModel
//	@Param			body	body	pictureGenerateRequest	true	"body of generating picture"
//	@Accept			json
//	@Success		201	{object}		multiplePicturesGenerateResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/multiple_pictures [post]
func (ctl *BigModelController) GenMultiplePictures(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := pictureGenerateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if err := req.validate(); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, code, err := ctl.s.GenPictures(pl.DomainAccount(), req.Desc)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, multiplePicturesGenerateResp{v})
	}
}

//	@Title			Ask
//	@Description	ask question based on a picture
//	@Tags			BigModel
//	@Param			body	body	questionAskRequest	true	"body of ask question"
//	@Accept			json
//	@Success		201	{object}		questionAskResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/ask [post]
func (ctl *BigModelController) Ask(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := questionAskRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	q, f, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, code, err := ctl.s.Ask(
		pl.DomainAccount(), q,
		filepath.Join(pl.Account, f),
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, questionAskResp{v})
	}
}

//	@Title			AskHF
//	@Description	vqa for hf
//	@Tags			BigModel
//	@Param			picture	formData	file	true	"picture"
//	@Accept			json
//	@Success		201	{object}		describePictureResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/ask_hf [post]
func (ctl *BigModelController) AskHF(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	req := questionAskHFReq{}

	// get picture
	f, err := ctx.FormFile("picture")
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if f.Size > apiConfig.MaxPictureSizeToDescribe {
		ctl.sendBadRequestParamWithMsg(ctx, "too big picture")

		return
	}

	p, err := f.Open()
	if err != nil {
		ctl.sendBadRequestParamWithMsg(ctx, "can't get picture")

		return
	}

	defer p.Close()

	req.Picture = p

	// get question
	req.Question = ctx.PostForm("question")

	var cmd app.VQAHFCmd
	if cmd, err = req.toCmd(); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, code, err := ctl.s.VQAHF(&cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, questionAskResp{v})
	}
}

//	@Title			PanGu
//	@Description	pan-gu big model
//	@Tags			BigModel
//	@Param			body	body	panguRequest	true	"body of pan-gu"
//	@Accept			json
//	@Success		201	{object}		panguResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/pangu [post]
func (ctl *BigModelController) PanGu(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := panguRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	v, code, err := ctl.s.PanGu(pl.DomainAccount(), req.Question)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, panguResp{v})
	}
}

//	@Title			LuoJia
//	@Description	luo-jia big model
//	@Tags			BigModel
//	@Accept			json
//	@Success		201	{object}		luojiaResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/luojia [post]
func (ctl *BigModelController) LuoJia(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, err := ctl.s.LuoJia(pl.DomainAccount()); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, luojiaResp{v})
	}
}

//	@Title			LuoJiaHF
//	@Description	luojia for hf
//	@Tags			BigModel
//	@Param			picture	formData	file	true	"picture"
//	@Accept			json
//	@Success		201	{object}		describePictureResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/luojia_hf [post]
func (ctl *BigModelController) LuoJiaHF(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	req := luojiaHFReq{}

	// get picture
	f, err := ctx.FormFile("picture")
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if f.Size > apiConfig.MaxPictureSizeToDescribe {
		ctl.sendBadRequestParamWithMsg(ctx, "too big picture")

		return
	}

	p, err := f.Open()
	if err != nil {
		ctl.sendBadRequestParamWithMsg(ctx, "can't get picture")

		return
	}

	defer p.Close()

	req.Picture = p

	var cmd app.LuoJiaHFCmd
	if cmd, err = req.toCmd(); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, err := ctl.s.LuoJiaHF(&cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, luojiaResp{v})
	}
}

//	@Title			ListLuoJiaRecord
//	@Description	list luo-jia big model records
//	@Tags			BigModel
//	@Accept			json
//	@Success		200	{object}		app.LuoJiaRecordDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/luojia [get]
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

//	@Title			CodeGeex
//	@Description	codegeex big model
//	@Tags			BigModel
//	@Param			body	body	CodeGeexRequest	true	"codegeex body"
//	@Accept			json
//	@Success		201	{object}		app.CodeGeexDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/codegeex [post]
func (ctl *BigModelController) CodeGeex(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := CodeGeexRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	v, code, err := ctl.s.CodeGeex(pl.DomainAccount(), &cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, v)
	}
}

//	@Title			VQAUploadPicture
//	@Description	upload a picture for vqa
//	@Tags			BigModel
//	@Param			picture	formData	file	true	"picture"
//	@Accept			json
//	@Success		201	{object}		pictureUploadResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/vqa_upload_picture [post]
func (ctl *BigModelController) VQAUploadPicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	f, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	if f.Size > apiConfig.MaxPictureSizeToVQA {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "too big picture",
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

	if err := ctl.s.VQAUploadPicture(p, pl.DomainAccount(), f.Filename); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(pictureUploadResp{f.Filename}))
	}
}

//	@Title			LuoJiaUploadPicture
//	@Description	upload a picture for luo-jia
//	@Tags			BigModel
//	@Param			picture	formData	file	true	"picture"
//	@Accept			json
//	@Success		201	{object}		pictureUploadResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/luojia_upload_picture [post]
func (ctl *BigModelController) LuoJiaUploadPicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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

	if err := ctl.s.LuoJiaUploadPicture(p, pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(pictureUploadResp{f.Filename}))
	}
}

//	@Title			GenWuKongSamples
//	@Description	gen wukong samples
//	@Tags			BigModel
//	@Param			batch	path	int	true	"batch num"
//	@Accept			json
//	@Success		201
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/samples/{batch} [get]
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

//	@Title			WuKong
//	@Description	generates pictures by WuKong
//	@Tags			BigModel
//	@Param			body	body	wukongRequest	true	"body of wukong"
//	@Accept			json
//	@Success		201	{object}		wukongPicturesGenerateResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong [post]
func (ctl *BigModelController) WuKong(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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

//	@Title			WuKong
//	@Description	generates pictures by WuKong-hf
//	@Tags			BigModel
//	@Param			body	body	wukongHFRequest	true	"body of wukong"
//	@Accept			json
//	@Success		201	{object}				wukongPicturesGenerateResp
//	@Failure		500	system_error			system	error
//	@Failure		404	bigmodel_sensitive_info	picture	error
//	@Router			/v1/bigmodel/wukong_hf [post]
func (ctl *BigModelController) WuKongHF(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	req := wukongHFRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, code, err := ctl.s.WuKongHF(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, wukongPicturesGenerateResp{v})
	}
}

//	@Title			WuKong
//	@Description	generates pictures by WuKong-icbc
//	@Tags			BigModel
//	@Param			body	body	wukongICBCRequest	true	"body of wukong"
//	@Accept			json
//	@Success		201	{object}				wukongPicturesGenerateResp
//	@Failure		500	system_error			system	error
//	@Failure		404	bigmodel_sensitive_info	picture	error
//	@Router			/v1/bigmodel/wukong_icbc [post]
func (ctl *BigModelController) WuKongICBC(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	req := wukongICBCRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, code, err := ctl.s.WuKong(cmd.User, &cmd.WuKongCmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, wukongPicturesGenerateResp{v})
	}
}

//	@Title			WuKong
//	@Description	send async wukong request task
//	@Tags			BigModel
//	@Param			body	body	wukongRequest	true	"body of wukong"
//	@Accept			json
//	@Success		201	{object}		wukongPicturesGenerateResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong_async [post]
func (ctl *BigModelController) WuKongAsync(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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
		ctl.sendRespOfPost(ctx, "")
	}
}

//	@Title			WuKong
//	@Description	get wukong rank
//	@Tags			BigModel
//	@Accept			json
//	@Success		200	{object}		app.WuKongRankDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/rank [get]
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
			ws.WriteJSON(newResponseError(err))

			log.Errorf("get rank failed: get status, err:%s", err.Error())

			return
		} else {
			ws.WriteJSON(newResponseData(dto))
		}

		log.Debugf("info dto:%v", dto)

		if dto.Rank == 0 {
			log.Debug("task done")

			return
		}

		time.Sleep(time.Second)
	}
}

//	@Title			WuKong
//	@Description	get last finished task
//	@Tags			BigModel
//	@Accept			json
//	@Success		200	{object}		app.WuKongRankDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/task [get]
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

//	@Title			AddLike
//	@Description	add like to wukong picture
//	@Tags			BigModel
//	@Accept			json
//	@Success		202	{object}		wukongAddLikeResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/like [post]
func (ctl *BigModelController) AddLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	reqTemp := wukongAddLikeFromTempRequest{}
	reqPublic := wukongAddLikeFromPublicRequest{}

	errTemp := ctx.ShouldBindBodyWith(&reqTemp, binding.JSON)
	errPublic := ctx.ShouldBindBodyWith(&reqPublic, binding.JSON)
	if errTemp != nil && errPublic != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)
		return
	}

	if errTemp == nil {
		cmd := reqTemp.toCmd(pl.DomainAccount())
		if pid, code, err := ctl.s.AddLikeFromTempPicture(&cmd); err != nil {
			ctl.sendCodeMessage(ctx, code, err)
		} else {
			ctl.sendRespOfPost(ctx, wukongAddLikeResp{pid})
		}
	}

	if errPublic == nil {
		cmd := reqPublic.toCmd(pl.DomainAccount())
		if pid, code, err := ctl.s.AddLikeFromPublicPicture(&cmd); err != nil {
			ctl.sendCodeMessage(ctx, code, err)
		} else {
			ctl.sendRespOfPost(ctx, wukongAddLikeResp{pid})
		}
	}
}

//	@Title			CancelLike
//	@Description	cancel like on wukong picture
//	@Tags			BigModel
//	@Param			id	path	string	true	"picture id"
//	@Accept			json
//	@Success		204
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/like/{id} [delete]
func (ctl *BigModelController) CancelLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	err := ctl.s.CancelLike(
		pl.DomainAccount(), ctx.Param("id"),
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfDelete(ctx)
	}
}

//	@Title			CancelPublic
//	@Description	cancel public on wukong picture
//	@Tags			BigModel
//	@Param			id	path	string	true	"picture id"
//	@Accept			json
//	@Success		204
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/public/{id} [delete]
func (ctl *BigModelController) CancelPublic(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	err := ctl.s.CancelPublic(
		pl.DomainAccount(), ctx.Param("id"),
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfDelete(ctx)
	}
}

//	@Title			ListLike
//	@Description	list wukong pictures user liked
//	@Tags			BigModel
//	@Accept			json
//	@Success		200	{object}		app.WuKongLikeDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong [get]
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

//	@Title			AddDigg
//	@Description	add digg to wukong picture
//	@Tags			BigModel
//	@Accept			json
//	@Success		202	{object}		wukongDiggResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/digg [post]
func (ctl *BigModelController) AddDigg(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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

//	@Title			CancelDigg
//	@Description	delete digg to wukong picture
//	@Tags			BigModel
//	@Param			body	body	wukongCancelDiggPublicRequest	true	"body of wukong"
//	@Accept			json
//	@Success		202	{object}		wukongDiggResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/digg [delete]
func (ctl *BigModelController) CancelDigg(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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

//	@Title			AddPublic
//	@Description	add public to wukong picture
//	@Tags			BigModel
//	@Accept			json
//	@Success		202	{object}		wukongAddPublicResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/public [post]
func (ctl *BigModelController) AddPublic(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	reqTemp := wukongAddPublicFromTempRequest{}
	reqPublic := wukongAddPublicFromLikeRequest{}

	errTemp := ctx.ShouldBindBodyWith(&reqTemp, binding.JSON)
	errPublic := ctx.ShouldBindBodyWith(&reqPublic, binding.JSON)
	if errTemp != nil && errPublic != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)
		return
	}

	if errTemp == nil {
		cmd := reqTemp.toCmd(pl.DomainAccount())
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

//	@Title			ListPublic
//	@Description	list wukong pictures user publiced
//	@Tags			BigModel
//	@Accept			json
//	@Success		200	{object}		app.WuKongPublicDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/public [get]
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

//	@Title			GetPublicGlobal
//	@Description	list all wukong pictures publiced
//	@Tags			BigModel
//	@Accept			json
//	@Success		200	{object}		app.WuKongPublicDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/publics [get]
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

//	@Title			GenDownloadURL
//	@Description	generate download url of wukong picture
//	@Tags			BigModel
//	@Param			body	body	wukongPictureLink	true	"body of wukong"
//	@Accept			json
//	@Success		202	{object}		wukongPictureLink
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/link [put]
func (ctl *BigModelController) GenDownloadURL(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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

//	@Title			AI Detector
//	@Description	detecte if text generate by ai
//	@Tags			BigModel
//	@Param			body	body	aiDetectorReq	true	"body of ai detector"
//	@Accept			json
//	@Success		202	{object}		aiDetectorResp
//	@Failure		500	system_error	system	error
//	@Router			/v1/bigmodel/wukong/link [put]
func (ctl *BigModelController) AIDetector(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := aiDetectorReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
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
