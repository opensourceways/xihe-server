package application

import (
	"context"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/opensourceways/xihe-server/domain/entity"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure"
	"github.com/opensourceways/xihe-server/util"
)

type UserApp struct {
	repo repository.UserRepository
}

func NewUserApp(repo repository.UserRepository) *UserApp {
	app := new(UserApp)
	app.repo = repo
	return app
}

func (f *UserApp) Save(item *entity.User) (*entity.User, error) {

	return f.repo.Save(item)
}

func (f *UserApp) FindUser(item *infrastructure.FindUserRequest) (*model.User, error) {
	thisUser, err := infrastructure.AuthingManageClient.FindUser(infrastructure.Change2FindUserRequest(item))
	if err != nil {
		return nil, err
	}
	return thisUser, nil
}

func (f *UserApp) GetCurrentUser(accessToken string) (string, error) {
	authingClient := authentication.NewClient(util.GetConfig().AuthingConfig.AppID, util.GetConfig().AuthingConfig.AppSecret)
	userDetail, err := authingClient.GetUserInfoByAccessToken(accessToken)
	if err != nil {
		return "", err
	}
	return userDetail, nil
}

func (f *UserApp) UpdatePhone(newphone, newcode, oldphone, oldcode string) (*model.User, error) {
	thisUser, err := infrastructure.AuthingDefaultUserClient.UpdatePhone(newphone, newcode, &oldphone, &oldcode)
	if err != nil {
		return nil, err
	}
	return thisUser, nil
}

func (f *UserApp) BindPhone(userid string, updateUserInput model.UpdateUserInput) (*model.User, error) {
	result, err := infrastructure.AuthingManageClient.UpdateUser(userid, updateUserInput)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *UserApp) SendSmsCode(phoneNum string) (interface{}, error) {
	result, err := infrastructure.AuthingDefaultUserClient.SendSmsCode(phoneNum)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (f *UserApp) SendEmailToResetPswd(email string) (interface{}, error) {
	result, err := infrastructure.AuthingDefaultUserClient.SendEmail(email, model.EnumEmailSceneResetPassword)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (f *UserApp) SendEmailToVerifyEmail(email string) (interface{}, error) {
	result, err := infrastructure.AuthingDefaultUserClient.SendEmail(email, model.EnumEmailSceneVerifyEmail)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (f *UserApp) ResetPasswordByEmailCode(email, code, newpswd string) (interface{}, error) {
	result, err := infrastructure.AuthingDefaultUserClient.ResetPasswordByEmailCode(email, code, newpswd)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (f *UserApp) UpdateUser(subID string, updateUserInput entity.User) (interface{}, error) {
	authingUpdateUserInput, err := updateUserInput.ExportToAuthingData()
	if err != nil {
		return nil, err
	}
	result, err := infrastructure.AuthingManageClient.UpdateUser(subID, *authingUpdateUserInput)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *UserApp) GetTokenFromAuthing(code string) (interface{}, error) {
	oauth2Token, err := infrastructure.OIDCConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	userInfo, err := infrastructure.GetUserInfoByToken(oauth2Token.AccessToken)
	if err != nil {

		return nil, err
	}

	token, err := infrastructure.GetJwtString(util.GetConfig().JwtConfig.Expire, userInfo.Sub, userInfo.Name, userInfo.ExternalID)
	if err != nil {

		return nil, err
	}
	result := &struct {
		AccessToken  string                           `json:"accessToken"`
		RefreshToken string                           `json:"refreshToken"`
		Token        string                           `json:"token"`
		User         *infrastructure.AuthingLoginUser `json:"user"`
	}{}
	result.User = userInfo
	result.AccessToken = oauth2Token.AccessToken
	result.RefreshToken = oauth2Token.RefreshToken
	result.Token = token

	return result, nil
}
