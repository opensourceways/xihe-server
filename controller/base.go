package controller

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	common "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	PrivateToken       = "PRIVATE-TOKEN"
	csrfToken          = "CSRF-Token" // #nosec G101 -- this is a false positive
	headerLanguage     = "Expect-Language"
	Token              = "token"
	encodeUsername     = "encode-username"
	headerSecWebsocket = "Sec-Websocket-Protocol"
	PayLoad            = "PAYLOAD"

	roleIndividuals = "individuals"
	fileReadme      = "README.md"
	fileApp         = "app.py"

	visitorPrefix = "visitor"
)

type baseController struct {
}

func (ctl baseController) newApiToken(ctx *gin.Context, pl interface{}) (
	token string, csrftoken string, err error,
) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		return "", "", err
	}

	ac := &accessController{
		Expiry:     utils.Expiry(apiConfig.TokenExpiry),
		Role:       roleIndividuals,
		Payload:    pl,
		RemoteAddr: addr,
	}

	t, err := ac.newToken(apiConfig.TokenKey)
	if err != nil {
		return "", "", err
	}

	// clear token and csrf-token memory
	defer utils.ClearStringMemory(t)

	if token, err = ctl.encryptData(t); err != nil {
		return
	}

	if csrftoken, err = ctl.encryptDataForCSRF(t); err != nil {
		return
	}

	return
}

func (ctl baseController) checkToken(
	ctx *gin.Context, token string, pl interface{},
) (ac accessController, tokenbyte []byte, ok bool) {
	tokenbyte, err := ctl.decryptData(token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ac = accessController{
		Payload: pl,
	}

	if err := ac.initByToken(string(tokenbyte), apiConfig.TokenKey); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ac.verify([]string{roleIndividuals}); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, err),
		)
		return
	}

	ok = true

	return
}

func (ctl baseController) checkCSRFToken(
	ctx *gin.Context, tokenbyte []byte, csrftoken string,
) (ok bool) {
	csrfbyte, err := ctl.decryptDataForCSRF(csrftoken)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}
	defer utils.ClearByteArrayMemory(csrfbyte)

	if !verifyCSRFToken(tokenbyte, csrfbyte) {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, errors.New("token not allowed")),
		)

		return
	}

	ok = true

	return
}

func (ctl baseController) refreshDoubleToken(ac accessController) (token, csrftoken string) {
	decode, err := ac.refreshToken(apiConfig.TokenExpiry, apiConfig.TokenKey)
	if err == nil {
		if t, err := ctl.encryptData(decode); err == nil {
			token = t
		}
	}

	if len(decode) != 0 {
		if t, err := ctl.encryptDataForCSRF(decode); err == nil {
			csrftoken = t
		}
	}

	return
}

func (ctl baseController) checkApiToken(
	ctx *gin.Context, token string, csrftoken string, pl interface{}, refresh bool,
) bool {
	ac, tokenbyte, ok := ctl.checkToken(ctx, token, pl)
	if !ok {
		logrus.Warnf("check token error")
		return ok
	}
	defer utils.ClearByteArrayMemory(tokenbyte)

	if ok = ctl.checkCSRFToken(ctx, tokenbyte, csrftoken); !ok {
		logrus.Warnf("check csrf token error")
		return ok
	}

	if !refresh {
		return ok
	}

	token, csrftoken = ctl.refreshDoubleToken(ac)
	defer utils.ClearStringMemory(token, csrftoken)

	payload, ok := ac.Payload.(*oldUserTokenPayload)
	if !ok {
		logrus.Warnf("payload assert error")
		return ok
	}

	if err := ctl.setRespToken(ctx, token, csrftoken, payload.Account, apiConfig.SessionDomain); err != nil {
		logrus.Debugf("set resp token error: %s", err.Error())
	}

	return ok
}

func (ctl baseController) checkUserApiToken(
	ctx *gin.Context, allowVistor bool,
) (
	pl *oldUserTokenPayload, visitor bool, ok bool,
) {
	pl, visitor, ok = ctl.checkUserApiTokenBase(ctx, allowVistor, true)
	if !ok {
		if pl == nil || pl.Account == "" {
			prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check api token")
		} else {
			prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_SYSTEM, "check api token")
		}
	}

	return pl, visitor, ok
}

func (ctl baseController) checkUserApiTokenNoRefresh(
	ctx *gin.Context, allowVistor bool,
) (
	pl *oldUserTokenPayload, visitor bool, ok bool,
) {
	pl, visitor, ok = ctl.checkUserApiTokenBase(ctx, allowVistor, false)
	if !ok {
		if pl == nil || pl.Account == "" {
			prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check api token")
		} else {
			prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_SYSTEM, "check api token")
		}
	}

	return pl, visitor, ok
}

func (ctl baseController) checkUserApiTokenBase(
	ctx *gin.Context, allowVistor bool, refresh bool,
) (
	pl *oldUserTokenPayload, visitor bool, ok bool,
) {
	token, err := ctl.getCookieToken(ctx)
	if err != nil {
		return
	}

	csrftoken, err := ctl.getCSRFToken(ctx)
	if err != nil {
		return
	}

	pl = new(oldUserTokenPayload)

	if token == "" || csrftoken == "" {
		// try best to grab userinfo
		if token != "" {
			_, _, _ = ctl.checkToken(ctx, token, pl)
			fmt.Printf("pl: %+v\n", pl)
		}
		if allowVistor {
			visitor = true
			ok = true
		} else {
			ctx.JSON(
				http.StatusUnauthorized,
				newResponseCodeMsg(errorBadRequestHeader, "no token"),
			)
		}

		return
	}

	ok = ctl.checkApiToken(ctx, token, csrftoken, pl, refresh)

	// set payload address in context
	ctx.Set(PayLoad, pl)

	return
}

func (ctl baseController) setRespCookieToken(ctx *gin.Context, token, username, domain string) error {
	// encrypt username
	u, err := ctl.encryptData(username)
	if err != nil {
		return err
	}

	// insert redis
	if err = ctl.newRepo().Insert(u, token); err != nil {
		return err
	}

	// set expire time for old token
	var oldusername = ""
	var ok bool
	o, exist := ctx.Get(encodeUsername)
	if exist {
		if oldusername, ok = o.(string); !ok {
			return errors.New("encode username illegal")
		}
	}

	if err = ctl.newRepo().Expire(oldusername, 3); err != nil {
		return err
	}

	// set cookie
	setCookie(ctx, PrivateToken, u, domain, true,
		utils.ExpiryReduceSecond(apiConfig.TokenExpiry), http.SameSiteLaxMode)

	return nil
}

func (ctl baseController) setRespCSRFToken(ctx *gin.Context, token, domain string) {
	setCookie(ctx, csrfToken, token, domain, false,
		utils.ExpiryReduceSecond(apiConfig.TokenExpiry), http.SameSiteStrictMode)
}

func setCookie(ctx *gin.Context, key, val, domain string, httpOnly bool,
	expireTime time.Time, sameSite http.SameSite) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     "/",
		Expires:  expireTime,
		HttpOnly: httpOnly,
		Secure:   true,
		SameSite: sameSite,
		Domain:   domain,
	}

	http.SetCookie(ctx.Writer, cookie)
}

func (ctl baseController) setRespToken(ctx *gin.Context, token, csrftoken, username, domain string) error {
	ctl.setRespCSRFToken(ctx, csrftoken, domain)

	if apiConfig.LocalDomainCookie {
		logrus.Info("local set cookie domain empty")

		ctl.setRespCSRFToken(ctx, csrftoken, "")
	}

	if err := ctl.setRespCookieToken(ctx, token, username, domain); err != nil {
		return err
	}
	if apiConfig.LocalDomainCookie {
		logrus.Info("local set cookie domain empty")

		if err := ctl.setRespCookieToken(ctx, token, username, ""); err != nil {
			return err
		}
	}

	return nil
}

func (ctl *baseController) checkTokenForWebsocket(
	ctx *gin.Context, allowVistor bool,
) (
	pl *oldUserTokenPayload, csrftoken string, visitor, ok bool,
) {
	csrftoken = ctl.getTokenForWebsocket(ctx)

	if strings.HasPrefix(csrftoken, visitorPrefix) {
		return pl, csrftoken, true, true
	}

	if csrftoken == "" {
		if allowVistor {
			visitor = true
			ok = true
		} else {
			ctx.JSON(
				http.StatusUnauthorized,
				newResponseCodeMsg(errorBadRequestHeader, "no token"),
			)
			prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check token for web socket")
		}

		return
	}

	pl = new(oldUserTokenPayload)
	ok = ctl.checkCSRFTokenForWebSocket(ctx, csrftoken, pl)
	if !ok {
		if pl.Account == "" {
			prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check token for web socket")
		} else {
			prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_SYSTEM, "check token for web socket")
		}

	}

	// set payload address in context
	ctx.Set(PayLoad, pl)

	return
}

func (ctl *baseController) getCookieToken(ctx *gin.Context) (token string, err error) { // TODO add return value exist
	// get encode username
	u, err := getCookieValue(ctx, PrivateToken)
	if err != nil {
		logrus.Warnf("failed to get cookie token %s", err)
		return "", nil
	}

	fmt.Printf("private token ==================================u: %+v\n", u)

	// insert encode username to context
	ctx.Set(encodeUsername, u)

	// get token from redis
	token, err = ctl.newRepo().Get(u)
	if err != nil {
		logrus.Warnf("failed to get cookie token from redis %s", err)
		return "", nil
	}

	return
}

func (ctl *baseController) getCSRFToken(ctx *gin.Context) (string, error) {
	v := ctx.GetHeader(csrfToken)
	return v, nil
}

func getCookieValue(ctx *gin.Context, key string) (string, error) {
	cookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return "", nil
	}

	return cookie.Value, nil
}

func (ctl *baseController) getTokenForWebsocket(ctx *gin.Context) (csrftoken string) {
	return ctx.GetHeader(headerSecWebsocket)
}

func (ctl *baseController) languageRuquested(ctx *gin.Context) (common.Language, error) {
	v := common.NewLanguage(ctx.GetHeader(headerLanguage))
	if v == nil {
		return nil, errors.New("unknown language")
	}

	return v, nil
}

func (ctl baseController) checkCSRFTokenForWebSocket(
	ctx *gin.Context, csrftoken string, pl interface{},
) (ok bool) {
	ctokenbyte, err := ctl.decryptDataForCSRF(csrftoken)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}
	defer utils.ClearByteArrayMemory(ctokenbyte)

	ac := accessController{
		Payload: pl,
	}

	if err := ac.initByToken(string(ctokenbyte), apiConfig.TokenKey); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ac.verify([]string{roleIndividuals}); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, err),
		)
		return
	}

	ok = true

	// set payload address in context
	ctx.Set(PayLoad, pl)

	return
}

func (ctl baseController) getRemoteAddr(ctx *gin.Context) (string, error) {
	ips := ctx.ClientIP()

	for _, item := range strings.Split(ips, ", ") {
		if net.ParseIP(item) != nil {
			return item, nil
		}
	}

	return "", errors.New("can not fetch client ip")
}

// crypt for token
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

// crypt for csrftoken
func (ctl baseController) encryptDataForCSRF(d string) (string, error) {
	t, err := encryptHelperCSRF.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (ctl baseController) decryptDataForCSRF(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelperCSRF.Decrypt(dst)
}

// crypt for database
func (ctl baseController) encryptDataForToken(d string) (string, error) {
	t, err := encryptHelperToken.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (ctl baseController) decryptDataForToken(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelperToken.Decrypt(dst)
}

func (ctl baseController) cleanCookie(ctx *gin.Context, domain string) {
	setCookie(ctx, PrivateToken, "", domain, true, time.Now().AddDate(0, 0, -1), http.SameSiteLaxMode)
	setCookie(ctx, PrivateToken, "", "", true, time.Now().AddDate(0, 0, -1), http.SameSiteLaxMode)
	t, ok := ctx.Get(encodeUsername)
	if !ok {
		logrus.Warnf("context get encode username failed")
	}

	u, ok2 := t.(string)
	if !ok2 {
		logrus.Warnf("encode username convert error")
	}

	_ = ctl.newRepo().Expire(u, 0)

	setCookie(ctx, csrfToken, "", domain, false, time.Now().AddDate(0, 0, -1), http.SameSiteStrictMode)
	setCookie(ctx, csrfToken, "", "", false, time.Now().AddDate(0, 0, -1), http.SameSiteStrictMode)
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

func (ctl baseController) newRepo() repositories.Access {
	return repositories.NewAccessRepo(int(apiConfig.TokenExpiry - 10))
}

func (ctl baseController) checkBigmodelApiToken(ctx *gin.Context) (user string, ok bool) {
	v := ctx.GetHeader(Token)
	deToken, err := ctl.decryptData(v)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "invalid token",
		))
		prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check bigmodel api token")
		return
	}
	defer utils.ClearByteArrayMemory(deToken)

	strs := strings.Split(string(deToken), "+")
	user = strs[0]

	time, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "invalid token",
		))

		if user == "" {
			prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check bigmodel api token")
		} else {
			prepareOperateLog(ctx, user, OPERATE_TYPE_SYSTEM, "check bigmodel api token")
		}

		return
	}

	if utils.Now()-time > 5184000 {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "token expire",
		))

		if user == "" {
			prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "check bigmodel api token")
		} else {
			prepareOperateLog(ctx, user, OPERATE_TYPE_SYSTEM, "check bigmodel api token")
		}

		return
	}

	return user, true
}
