package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForBigModelController(
	rg *gin.RouterGroup,
	bm bigmodel.BigModel,
	luojia repository.LuoJia,
	wukong repository.WuKong,
	wukongPicture repository.WuKongPicture,
	sender message.Sender,
) {
	ctl := BigModelController{
		s: app.NewBigModelService(bm, luojia, wukong, wukongPicture, sender),
	}

	rg.POST("/v1/bigmodel/describe_picture", ctl.DescribePicture)
	rg.POST("/v1/bigmodel/single_picture", ctl.GenSinglePicture)
	rg.POST("/v1/bigmodel/multiple_pictures", ctl.GenMultiplePictures)
	rg.POST("/v1/bigmodel/vqa_upload_picture", ctl.VQAUploadPicture)
	rg.POST("/v1/bigmodel/luojia_upload_picture", ctl.LuoJiaUploadPicture)
	rg.POST("/v1/bigmodel/ask", ctl.Ask)
	rg.POST("/v1/bigmodel/pangu", ctl.PanGu)
	rg.POST("/v1/bigmodel/luojia", ctl.LuoJia)
	rg.POST("/v1/bigmodel/codegeex", ctl.CodeGeex)
	rg.POST("/v1/bigmodel/wukong", ctl.WuKong)
	rg.PUT("/v1/bigmodel/wukong", ctl.AddLike)
	rg.PUT("/v1/bigmodel/wukong/link", ctl.GenDownloadURL)
	rg.DELETE("/v1/bigmodel/wukong/:id", ctl.CancelLike)
	rg.GET("/v1/bigmodel/wukong/samples/:batch", ctl.GenWuKongSamples)
	rg.GET("/v1/bigmodel/wukong/pictures", ctl.WuKongPictures)
	rg.GET("/v1/bigmodel/wukong", ctl.ListLike)
	rg.GET("/v1/bigmodel/luojia", ctl.ListLuoJiaRecord)
}

type BigModelController struct {
	baseController

	s app.BigModelService
}

// @Title DescribePicture
// @Description describe a picture
// @Tags  BigModel
// @Param	picture		formData 	file	true	"picture"
// @Accept json
// @Success 201 {object} describePictureResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/describe_picture [post]
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

// @Title GenSinglePicture
// @Description generate a picture based on a text
// @Tags  BigModel
// @Param	body	body 	pictureGenerateRequest	true	"body of generating picture"
// @Accept json
// @Success 201 {object} pictureGenerateResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/single_picture [post]
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

// @Title GenMultiplePictures
// @Description generate multiple pictures based on a text
// @Tags  BigModel
// @Param	body	body 	pictureGenerateRequest	true	"body of generating picture"
// @Accept json
// @Success 201 {object} multiplePicturesGenerateResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/multiple_pictures [post]
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

// @Title Ask
// @Description ask question based on a picture
// @Tags  BigModel
// @Param	body	body 	questionAskRequest	true	"body of ask question"
// @Accept json
// @Success 201 {object} questionAskResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/ask [post]
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

// @Title PanGu
// @Description pan-gu big model
// @Tags  BigModel
// @Param	body	body 	panguRequest	true	"body of pan-gu"
// @Accept json
// @Success 201 {object} panguResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/pangu [post]
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

// @Title LuoJia
// @Description luo-jia big model
// @Tags  BigModel
// @Accept json
// @Success 201 {object} luojiaResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/luojia [post]
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

// @Title ListLuoJiaRecord
// @Description list luo-jia big model records
// @Tags  BigModel
// @Accept json
// @Success 200 {object} app.LuoJiaRecordDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/luojia [get]
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

// @Title CodeGeex
// @Description codegeex big model
// @Tags  BigModel
// @Param	body	body 	CodeGeexRequest		true	"codegeex body"
// @Accept json
// @Success 201 {object} app.CodeGeexDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/codegeex [post]
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

// @Title VQAUploadPicture
// @Description upload a picture for vqa
// @Tags  BigModel
// @Param	picture		formData 	file	true	"picture"
// @Accept json
// @Success 201 {object} pictureUploadResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/vqa_upload_picture [post]
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

// @Title LuoJiaUploadPicture
// @Description upload a picture for luo-jia
// @Tags  BigModel
// @Param	picture		formData 	file	true	"picture"
// @Accept json
// @Success 201 {object} pictureUploadResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/luojia_upload_picture [post]
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

// @Title GenWuKongSamples
// @Description gen wukong samples
// @Tags  BigModel
// @Param	batch	path 	int	true	"batch num"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong/samples/{batch} [get]
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

// @Title WuKongPictures
// @Description list wukong pictures
// @Tags  BigModel
// @Param	count_per_page	query	int	false	"count per page"
// @Param	page_num	query	int	false	"page num which starts from 1"
// @Accept json
// @Success 200 {object} app.WuKongPicturesDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong/pictures [get]
func (ctl *BigModelController) WuKongPictures(ctx *gin.Context) {
	cmd := app.WuKongPicturesListCmd{}

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

		return
	}

	if err := f(); err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.WuKongPictures(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title WuKong
// @Description generates pictures by WuKong
// @Tags  BigModel
// @Param	body	body 	wukongRequest	true	"body of wukong"
// @Accept json
// @Success 201 {object} wukongPicturesGenerateResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong [post]
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

// @Title AddLike
// @Description add like to wukong picture
// @Tags  BigModel
// @Param	body	body 	wukongAddLikeRequest	true	"body of wukong"
// @Accept json
// @Success 202 {object} wukongAddLikeResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong [put]
func (ctl *BigModelController) AddLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := wukongAddLikeRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	cmd := req.toCmd(pl.DomainAccount())
	if pid, code, err := ctl.s.AddLikeToWuKongPicture(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, wukongAddLikeResp{pid})
	}
}

// @Title CancelLike
// @Description cancel like on wukong picture
// @Tags  BigModel
// @Param	id	path 	string	true	"picture id"
// @Accept json
// @Success 204
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong/{id} [delete]
func (ctl *BigModelController) CancelLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	err := ctl.s.CancelLikeOnWuKongPicture(
		pl.DomainAccount(), ctx.Param("id"),
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfDelete(ctx)
	}
}

// @Title ListLike
// @Description list wukong pictures user liked
// @Tags  BigModel
// @Accept json
// @Success 200 {object} app.UserLikedWuKongPictureDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong [get]
func (ctl *BigModelController) ListLike(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, err := ctl.s.ListLikedWuKongPictures(pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title GenDownloadURL
// @Description generate download url of wukong picture
// @Tags  BigModel
// @Param	body	body 	wukongPictureLink	true	"body of wukong"
// @Accept json
// @Success 202 {object} wukongPictureLink
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong/link [put]
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

	link, code, err := ctl.s.ReGenerateDownloadURLOfWuKongPicture(
		pl.DomainAccount(), req.Link,
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, wukongPictureLink{link})
	}
}
