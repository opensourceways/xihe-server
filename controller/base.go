package controller

import (
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	headerPrivateToken = "PRIVATE-TOKEN"
	headerSecWebsocket = "Sec-Websocket-Protocol"

	roleIndividuals = "individuals"
	fileReadme      = "README.md"
)

type baseController struct {
}

func (ctl baseController) newApiToken(ctx *gin.Context, pl interface{}) (
	string, error,
) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		return "", err
	}

	ac := &accessController{
		Expiry:     utils.Expiry(apiConfig.TokenExpiry),
		Role:       roleIndividuals,
		Payload:    pl,
		RemoteAddr: addr,
	}

	token, err := ac.newToken(apiConfig.TokenKey)
	if err != nil {
		return "", err
	}

	return ctl.encryptData(token)
}

func (ctl baseController) checkApiToken(ctx *gin.Context, token string, pl interface{}, refresh bool) (ok bool) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	b, err := ctl.decryptData(token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ac := accessController{
		Payload: pl,
	}

	if err := ac.initByToken(string(b), apiConfig.TokenKey); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ac.verify([]string{roleIndividuals}, addr); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, err),
		)
		return
	}

	ok = true

	if !refresh {
		return
	}

	if v, err := ac.refreshToken(apiConfig.TokenExpiry, apiConfig.TokenKey); err == nil {
		if v, err = ctl.encryptData(v); err == nil {
			token = v
		}
	}

	ctx.Header(headerPrivateToken, token)

	return
}

func (ctl baseController) checkNewUserApiToken(ctx *gin.Context) (
	pl newUserTokenPayload, ok bool,
) {
	token := ctx.GetHeader(headerPrivateToken)
	if token == "" {
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestHeader, "no token"),
		)

		return
	}

	ok = ctl.checkApiToken(ctx, token, &pl, false)

	return
}

func (ctl baseController) checkUserApiToken(
	ctx *gin.Context, allowVistor bool,
) (
	pl oldUserTokenPayload, visitor bool, ok bool,
) {
	token := ctx.GetHeader(headerPrivateToken)
	if token == "" {
		if allowVistor {
			visitor = true
			ok = true
		} else {
			ctx.JSON(
				http.StatusBadRequest,
				newResponseCodeMsg(errorBadRequestHeader, "no token"),
			)
		}

		return
	}

	ok = ctl.checkApiToken(ctx, token, &pl, true)

	return
}

func (ctl baseController) setRespToken(ctx *gin.Context, token string) {
	ctx.Header(headerPrivateToken, token)
}

func (ctl baseController) getRemoteAddr(ctx *gin.Context) (string, error) {
	ips := ctx.Request.Header.Get("x-forwarded-for")

	for _, item := range strings.Split(ips, ", ") {
		if net.ParseIP(item) != nil {
			return item, nil
		}
	}

	return "", errors.New("can not fetch client ip")
}

func (ctl baseController) encryptData(d string) (string, error) {
	t, err := encryptHelper.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (ctl baseController) decryptData(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelper.Decrypt(dst)
}

func (ctl baseController) getQueryParameter(ctx *gin.Context, key string) string {
	return ctx.Request.URL.Query().Get(key)
}

func (ctl baseController) sendRespWithInternalError(ctx *gin.Context, data responseData) {
	log.Errorf("code: %s, err: %s", data.Code, data.Msg)

	ctx.JSON(http.StatusInternalServerError, data)
}

func (ctl baseController) sendCodeMessage(ctx *gin.Context, code string, err error) {
	if code == "" {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(code, err))
	}
}

func (ctl baseController) sendBadRequest(ctx *gin.Context, data responseData) {
	ctx.JSON(http.StatusBadRequest, data)
}

func (ctl baseController) sendBadRequestBody(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, respBadRequestBody)
}

func (ctl baseController) sendBadRequestParam(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, newResponseCodeError(errorBadRequestParam, err))
}

func (ctl baseController) sendBadRequestParamWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(errorBadRequestParam, msg))
}

func (ctl baseController) sendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}

func (ctl baseController) sendRespOfPost(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, newResponseData(data))
}

func (ctl baseController) sendRespOfPut(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusAccepted, newResponseData(data))
}

func (ctl baseController) sendRespOfDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

func (ctl baseController) getListResourceParameter(
	ctx *gin.Context,
) (cmd app.ResourceListCmd, err error) {
	if v := ctl.getQueryParameter(ctx, "name"); v != "" {
		cmd.Name = v
	}

	if v := ctl.getQueryParameter(ctx, "repo_type"); v != "" {
		r := strings.Split(v, "+")
		for i := range r {
			var t domain.RepoType
			t, err = domain.NewRepoType(r[i])
			if err != nil {
				return
			}
			cmd.RepoType = append(cmd.RepoType, t)
		}
	}

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

	if v := ctl.getQueryParameter(ctx, "sort_by"); v != "" {
		if cmd.SortType, err = domain.NewSortType(v); err != nil {
			return
		}
	}

	return
}

func (ctl baseController) getListGlobalResourceParameter(
	ctx *gin.Context,
) (cmd app.GlobalResourceListCmd, err error) {
	v, err := ctl.getListResourceParameter(ctx)
	if err != nil {
		return
	}

	if s := ctl.getQueryParameter(ctx, "tags"); s != "" {
		tags := strings.Split(s, ",")
		if len(tags) > apiConfig.MaxTagsNumToSearchResource {
			err = errors.New("too many tags to search by")

			return
		}

		cmd.Tags = tags
	}

	if s := ctl.getQueryParameter(ctx, "tag_kinds"); s != "" {
		kinds := strings.Split(s, ",")
		if len(kinds) > apiConfig.MaxTagKindsNumToSearchResource {
			err = errors.New("too many tag kinds to search by")

			return
		}

		cmd.TagKinds = kinds
	}

	if s := ctl.getQueryParameter(ctx, "level"); s != "" {
		cmd.Level = domain.NewResourceLevel(s)
	}

	cmd.ResourceListOption = v.ResourceListOption
	cmd.SortType = v.SortType

	return
}

func (ctl *baseController) checkTokenForWebsocket(ctx *gin.Context) (
	pl oldUserTokenPayload, token string, ok bool,
) {
	token = ctx.GetHeader(headerSecWebsocket)
	if token == "" {
		//TODO delete
		log.Errorf("check token for ws, no token")

		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestHeader, "no token"),
		)

		return
	}

	ok = ctl.checkApiToken(ctx, token, &pl, false)

	return
}
