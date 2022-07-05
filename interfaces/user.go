package interfaces

import (
	"net/http"

	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/opensourceways/xihe-server/application"
	"github.com/opensourceways/xihe-server/domain/entity"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure"
	"github.com/opensourceways/xihe-server/util"

	"github.com/gin-gonic/gin"
)

//Users struct defines the dependencies that will be used
type User struct {
	app *application.UserApp
}

//Users constructor
func NewUser(repo repository.UserRepository) *User {
	app := application.NewUserApp(repo)
	return &User{
		app: app,
	}
}

// @Summary CheckLogin
// @Description CheckLogin
// @Tags  Authing
// @Accept json
// @Produce json
// @Router /v1/user/checkLogin [get]
func (u *User) CheckLogin(c *gin.Context) {
	state, _ := util.RandString(16)
	nonce, _ := util.RandString(16)
	infrastructure.SetCallbackCookie(c.Writer, c.Request, "state", state)
	infrastructure.SetCallbackCookie(c.Writer, c.Request, "nonce", nonce)
	c.Redirect(http.StatusFound, infrastructure.OIDCConfig.AuthCodeURL(state, oidc.Nonce(nonce)))
}

// @Summary AuthingCallback
// @Description login success callback
// @Tags  Authing
// @Accept json
// @Produce json
// @Router /v1/user/callback [get]
func (u *User) AuthingCallback(c *gin.Context) {
	state, err := c.Request.Cookie("state")
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(util.CodeStatusClientError, "Cookie Error", err.Error()))
		return
	}
	if c.Request.URL.Query().Get("state") != state.Value {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusClientError, "state Error", err.Error()),
		)
		return
	}
	result, err := u.app.GetTokenFromAuthing(c.Request.URL.Query().Get("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(util.CodeStatusClientError, "Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary FindUser
// @Description FindUser
// @Tags  Authing
// @Param	body		body 	infrastructure.FindUserRequest	true		"email username phone"
// @Accept json
// @Produce json
// @Router /v1/user/findUser [get]
func (u *User) FindUser(c *gin.Context) {
	var findUserRequest infrastructure.FindUserRequest
	err := c.BindJSON(&findUserRequest)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "BindQuery error", err.Error()),
		)
		return
	}
	thisUser, err := u.app.FindUser(&findUserRequest)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "FindUser error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", thisUser))
}

// @Summary GetCurrentUser
// @Description 获取用户资料, 在请求的request的header中必须带有accessToken
// @Tags  Authing
// @Accept json
// @Produce json
// @Router /v1/user/getCurrentUser [get]
func (u *User) GetCurrentUser(c *gin.Context) {
	accessToken := c.GetHeader("accessToken")
	if len(accessToken) < 10 {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "accessToken inValidate ", accessToken),
		)
		return
	}
	userDetail, err := u.app.GetCurrentUser(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "GetCurrentUser error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", userDetail))
}

// @Summary UpdatePhone
// @Description 更换手机号，需要新号码和旧号码都发验证码
// @Tags  Authing
// @Param	newphone		query 	string	true		"new phone "
// @Param	newcode		query 	string	true		"new code of new phone"
// @Param	oldphone		query 	string	true		"old phone "
// @Param	oldcode		query 	string	true		"old code  of old phone "
// @Accept json
// @Produce json
// @Router /v1/user/updatePhone [put]
func (u *User) UpdatePhone(c *gin.Context) {
	newphone := c.Query("newphone")
	newcode := c.Query("newcode")
	oldphone := c.Query("oldphone")
	oldcode := c.Query("oldcode")
	thisUser, err := u.app.UpdatePhone(newphone, newcode, oldphone, oldcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "UpdatePhone error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", thisUser))
}

// @Summary SendSmsCode
// @Description 发送验证码
// @Tags  Authing
// @Param	phone		query 	string	true		"new phone "
// @Accept json
// @Produce json
// @Router /v1/user/sendSmsCode [get]
func (u *User) SendSmsCode(c *gin.Context) {
	phoneNum := c.Query("phone")
	result, err := infrastructure.AuthingDefaultUserClient.SendSmsCode(phoneNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "SendSmsCode error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary BindPhone
// @Description 绑定手机号
// @Tags  Authing
// @Param	phone		query 	string	true		"  phone "
// @Accept json
// @Produce json
// @Router /v1/user/bindPhone [put]
func (u *User) BindPhone(c *gin.Context) {
	phone := c.Query("phone")
	userid := c.Keys["id"].(string)
	if len(userid) == 0 {
		c.JSON(http.StatusBadRequest,
			util.ExportData(util.CodeStatusServerError, "token error", nil),
		)
		return
	}
	if len(phone) == 0 {
		c.JSON(http.StatusBadRequest,
			util.ExportData(util.CodeStatusServerError, "phone error", phone),
		)
		return
	}
	var updateUserInput model.UpdateUserInput
	updateUserInput.Phone = &phone
	result, err := infrastructure.AuthingManageClient.UpdateUser(userid, updateUserInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "BindPhone error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary SendEmailToResetPswd
// @Description 发送重置密码的邮件，内含验证码
// @Tags  Authing
// @Param	email		query 	string	true		"  email "
// @Accept json
// @Produce json
// @Router /v1/user/sendEmailToResetPswd [get]
func (u *User) SendEmailToResetPswd(c *gin.Context) {
	email := c.Query("email")
	result, err := infrastructure.AuthingDefaultUserClient.SendEmail(email, model.EnumEmailSceneResetPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "SendEmailToResetPswd error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary SendEmailToVerifyEmail
// @Description 发送验证邮件,
// @Tags  Authing
// @Param	email		query 	string	true		"  email "
// @Accept json
// @Produce json
// @Router /v1/user/sendEmailToVerifyEmail [get]
func (u *User) SendEmailToVerifyEmail(c *gin.Context) {
	email := c.Query("email")
	result, err := infrastructure.AuthingDefaultUserClient.SendEmail(email, model.EnumEmailSceneVerifyEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "SendEmailToVerifyEmail error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary ResetPasswordByEmailCode
// @Description 通过邮箱验证码来重置密码,
// @Tags  Authing
// @Param	email		query 	string	true		"  email "
// @Param	code		query 	string	true		"  code "
// @Param	newpswd		query 	string	true		"  newpswd "
// @Accept json
// @Produce json
// @Router /v1/user/resetPasswordByEmailCode [put]
func (u *User) ResetPasswordByEmailCode(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")
	newpswd := c.Query("newpswd")
	result, err := infrastructure.AuthingDefaultUserClient.ResetPasswordByEmailCode(email, code, newpswd)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "ResetPasswordByEmailCode error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary UpdateProfile
// @Description 修改用户资料， 需要修改哪个就填哪个。不填的不修改
// @Tags  Authing
// @Param	id		path 	string	true		"  user id "
// @Param	username		query 	string	false		"  username 用户名"
// @Param	nickname		query 	string	false		"  nickname 昵称"
// @Param	name		query 	string	false		"  name 真实姓名"
// @Param	photo		query 	string	false		"  photo 头像 "
// @Param	company		query 	string	false		"  company  公司 "
// @Param	givenName		query 	string	false		"  givenName "
// @Param	middleName		query 	string	false		"  middleName "
// @Param	profile		query 	string	false		"  Profile Url "
// @Param	gender		query 	string	false		"  gender   性别, M（Man） 表示男性、F（Female） 表示女性、未知表示 U（Unknown）"
// @Param	preferredUsername		query 	string	false		"  preferredUsername "
// @Param	website		query 	string	false		"  website 个人网站"
// @Param	address		query 	string	false		"  address 详细地址"
// @Param	birthdate		query 	string	false		"  birthdate 生日 "
// @Param	streetAddress		query 	string	false		"  streetAddress 街道地址"
// @Param	postalCode		query 	string	false		"  postalCode 邮编"
// @Param	city		query 	string	false		"  city 城市"
// @Param	province		query 	string	false		"  province 省份 "
// @Param	country		query 	string	false		"  country 国家"
// @Accept json
// @Produce json
// @Router /v1/user/updateProfile/{id} [put]
func (u *User) UpdateProfile(c *gin.Context) {

	subID := c.Param("id")
	var updateUserInput entity.User
	err := c.BindQuery(&updateUserInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "UpdateProfile BindJSON error", err.Error()),
		)
		return
	}
	authingUpdateUserInput, err := updateUserInput.ExportToAuthingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "UpdateProfile ExportTo error", err.Error()),
		)
		return
	}
	result, err := infrastructure.AuthingManageClient.UpdateUser(subID, *authingUpdateUserInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			util.ExportData(util.CodeStatusServerError, "UpdateProfile error", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}
