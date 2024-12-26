package app

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/agreement/app"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/domain/authing"
	platform "github.com/opensourceways/xihe-server/domain/platform"
	typerepo "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/login"
	"github.com/opensourceways/xihe-server/user/domain/message"
	pointsPort "github.com/opensourceways/xihe-server/user/domain/points"
	"github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type UserService interface {
	// user
	Create(*UserCreateCmd) (UserDTO, error)
	CreatePlatformAccount(*CreatePlatformAccountCmd) (PlatformInfoDTO, error)
	UpdatePlateformInfo(*UpdatePlateformInfoCmd) error
	UpdatePlateformToken(*UpdatePlateformTokenCmd) error
	NewPlatformAccountWithUpdate(*CreatePlatformAccountCmd) error
	UpdatePlatformAccountWithUpdate(int, *CreatePlatformAccountCmd) error
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error

	UpdateAgreement(u domain.Account, t app.AgreementType) error
	PrivacyRevoke(domain.Account) error
	AgreementRevoke(u domain.Account, t app.AgreementType) error

	UserInfo(domain.Account) (UserInfoDTO, error)
	GetByAccount(domain.Account) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)

	AddFollowing(*domain.FollowerInfo) error
	RemoveFollowing(*domain.FollowerInfo) error
	ListFollowing(*FollowsListCmd) (FollowsDTO, error)

	AddFollower(*domain.FollowerInfo) error
	RemoveFollower(*domain.FollowerInfo) error
	ListFollower(*FollowsListCmd) (FollowsDTO, error)

	RefreshGitlabToken(*RefreshTokenCmd) error

	ModifyInfo(ctx context.Context, cmd CmdToModifyInfo, user domain.Account) (string, error)
}

// ps: platform user service
func NewUserService(
	repo repository.User,
	ps platform.User,
	sender message.MessageProducer,
	points pointsPort.Points,
	encryption utils.SymmetricEncryption,
	auth authing.User,
	login login.Login,
) UserService {
	return userService{
		ps:         ps,
		repo:       repo,
		sender:     sender,
		points:     points,
		encryption: encryption,
		auth:       auth,
		login:      login,
	}
}

type userService struct {
	ps         platform.User
	repo       repository.User
	sender     message.MessageProducer
	points     pointsPort.Points
	encryption utils.SymmetricEncryption
	auth       authing.User
	login      login.Login
}

func (s userService) Create(cmd *UserCreateCmd) (dto UserDTO, err error) {
	v := cmd.toUser()
	// set agreement
	v.UserAgreement = app.GetCurrentUserAgree()
	v.CourseAgreement = app.GetCurrentCourseAgree()
	v.FinetuneAgreement = app.GetCurrentFinetuneAgree()

	// update user
	u, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	toUserDTO(&u, &dto)

	_ = s.sender.AddOperateLogForNewUser(u.Account)

	_ = s.sender.SendUserSignedUpEvent(&domain.UserSignedUpEvent{
		Account: cmd.Account,
	})

	return
}

func (s userService) UserInfo(account domain.Account) (dto UserInfoDTO, err error) {
	if dto.UserDTO, err = s.GetByAccount(account); err != nil {
		return
	}

	dto.Points, err = s.points.Points(account)

	return
}

func (s userService) UpdateAgreement(u domain.Account, t app.AgreementType) (err error) {
	user, err := s.repo.GetByAccount(u)
	if err != nil {
		return
	}

	// update agreement
	switch t {
	case app.Course:
		ver := app.GetCurrentCourseAgree()
		if user.CourseAgreement == ver {
			logrus.Infoln("no need update course agreement")
			return nil
		}
		user.CourseAgreement = ver
	case app.Finetune:
		ver := app.GetCurrentFinetuneAgree()
		if user.FinetuneAgreement == ver {
			logrus.Infoln("no need update finetune agreement")
			return nil
		}
		user.FinetuneAgreement = ver
	case app.User:
		ver := app.GetCurrentUserAgree()
		if user.UserAgreement == ver {
			logrus.Infoln("no need update user agreement")
			return nil
		}
		user.UserAgreement = ver
	default:
		str := fmt.Sprintf("Invalid agreement type :%s", t)
		logrus.Error(str)
		return fmt.Errorf("%s", str)
	}

	// update userinfo
	_, err = s.repo.Save(&user)

	return
}

// PrivacyRevoke revokes the privacy settings for a user.
func (s userService) PrivacyRevoke(user domain.Account) error {
	userInfo, err := s.repo.GetByAccount(user)
	if err != nil {
		if typerepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			return allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			return xerrors.Errorf("failed to get user: %w", err)
		}
	}

	userInfo.RevokePrivacy()
	if _, err = s.repo.Save(&userInfo); err != nil {
		return allerror.New(allerror.ErrorCodeRevokePrivacyFailed, "",
			xerrors.Errorf("failed to save user: %w", err))
	}

	return nil
}

func (s userService) AgreementRevoke(u domain.Account, t app.AgreementType) error {
	user, err := s.repo.GetByAccount(u)
	if err != nil {
		e := xerrors.Errorf("user %s not found: %w", user.Account.Account(), err)
		return allerror.New(allerror.ErrorCodeUserNotFound, "", e)
	}

	// revoke agreement
	switch t {
	case app.Course:
		user.CourseAgreement = ""
	case app.Finetune:
		user.FinetuneAgreement = ""
	default:
		err = errors.New("invalid agreement type")
		return allerror.New(allerror.ErrorCodeRevokeAgreementFailed, "",
			xerrors.Errorf("failed to save user: %w", err))
	}

	// update userinfo
	if _, err = s.repo.Save(&user); err != nil {
		return allerror.New(allerror.ErrorCodeRevokeAgreementFailed, "",
			xerrors.Errorf("failed to save user: %w", err))
	}

	return nil
}

func (s userService) GetByAccount(account domain.Account) (dto UserDTO, err error) {
	v, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	if v.PlatformToken.CreateAt == 0 && v.PlatformToken.Token != "" {
		if t, err := s.ps.GetToken(v.PlatformUser.Id); err == nil {
			v.PlatformToken.CreateAt = t.CreateAt
			//try our best to update the create_at filed
			logrus.Infof("will update token create_at for %s", v.Account.Account())
			_, _ = s.repo.Save(&v)
		} else {
			logrus.Warnf("get token for %s failed: %s", v.Account.Account(), err)
		}
	}

	if v.PlatformToken.Token != "" {
		token := v.PlatformToken.Token
		v.PlatformToken.Token, err = s.decryptToken(token)
		if err != nil {
			return
		}
	}

	toUserDTO(&v, &dto)

	return
}

func (s userService) GetByFollower(owner, follower domain.Account) (
	dto UserDTO, isFollower bool, err error,
) {
	v, isFollower, err := s.repo.GetByFollower(owner, follower)
	if err != nil {
		return
	}

	toUserDTO(&v, &dto)

	return
}

func (s userService) NewPlatformAccountWithUpdate(cmd *CreatePlatformAccountCmd) (err error) {
	// create platform account
	dto, err := s.CreatePlatformAccount(cmd)
	if err != nil {
		return
	}

	// update user information
	updatecmd := &UpdatePlateformInfoCmd{
		PlatformInfoDTO: dto,
		User:            cmd.Account,
		Email:           cmd.Email,
	}

	for i := 0; i <= 5; i++ {
		if err = s.UpdatePlateformInfo(updatecmd); err != nil {
			if !typerepo.IsErrorConcurrentUpdating(err) {
				return
			}
		} else {
			break
		}
	}

	return
}

func (s userService) CreatePlatformAccount(cmd *CreatePlatformAccountCmd) (dto PlatformInfoDTO, err error) {
	// create platform account
	pu, err := s.ps.New(platform.UserOption{
		Email:    cmd.Email,
		Name:     cmd.Account,
		Password: cmd.Password,
	})
	if err != nil {
		return
	}

	dto.PlatformUser = pu

	// apply token
	token, err := s.ps.NewToken(pu)
	if err != nil {
		return
	}

	eToken, err := s.encryptToken(token.Token)
	if err != nil {
		return
	}

	dto.PlatformToken = domain.PlatformToken{
		Token:    eToken,
		CreateAt: token.CreateAt,
	}

	return
}

func (s userService) UpdatePlateformInfo(cmd *UpdatePlateformInfoCmd) (err error) {
	// get userinfo
	u, err := s.repo.GetByAccount(cmd.User)
	if err != nil {
		return
	}

	// update some data
	u.PlatformUser = cmd.PlatformUser
	u.PlatformToken = cmd.PlatformToken
	u.Email = cmd.Email

	// update userinfo
	if _, err = s.repo.Save(&u); err != nil {
		return
	}

	return
}

func (s userService) UpdatePlateformToken(cmd *UpdatePlateformTokenCmd) (err error) {
	// get userinfo
	u, err := s.repo.GetByAccount(cmd.User)
	if err != nil {
		return
	}

	// update token
	u.PlatformToken = cmd.PlatformToken

	// update userinfo
	if _, err = s.repo.Save(&u); err != nil {
		return
	}

	return
}

func (s userService) encryptToken(d string) (string, error) {
	t, err := s.encryption.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (s userService) decryptToken(d string) (string, error) {
	tb, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}

	dtoken, err := s.encryption.Decrypt(tb)
	if err != nil {
		return "", err
	}

	return string(dtoken), nil
}

func (s userService) RefreshGitlabToken(cmd *RefreshTokenCmd) (err error) {
	token, err := s.ps.RefreshToken(cmd.Id)
	if err != nil {
		return
	}

	eToken, err := s.encryptToken(token.Token)
	if err != nil {
		return
	}

	updatecmd := &UpdatePlateformTokenCmd{
		User: cmd.Account,
		PlatformToken: domain.PlatformToken{
			Token:    eToken,
			CreateAt: token.CreateAt,
		},
	}

	for i := 0; i <= 5; i++ {
		if err = s.UpdatePlateformToken(updatecmd); err != nil {
			if !typerepo.IsErrorConcurrentUpdating(err) {
				return
			}
		} else {
			break
		}
	}

	return
}

func (s userService) ModifyInfo(ctx context.Context, cmd CmdToModifyInfo, user domain.Account,
) (errMsgCode string, err error) {
	//校验验证码
	info, err := s.login.GetAccessAndIdToken(user)
	if err != nil {
		return
	}

	if info.UserId == "" {
		errMsgCode = errorNoUserId
		err = errors.New("cannot read user id")

		return
	}

	if errMsgCode, err = s.auth.VerifyBindEmail(cmd.Account, cmd.Code, info.UserId); err != nil {
		if !isCodeUserDuplicateBind(errMsgCode) {
			return
		}
	}

	if errMsgCode, err = s.auth.VerifyBindEmail(cmd.OldAccount, cmd.OldCode, info.UserId); err != nil {
		if !isCodeUserDuplicateBind(errMsgCode) {
			return
		}
	}
	//获取用户信息
	u, err := s.repo.GetByAccount(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeUserNotFound, "",
				xerrors.Errorf("user %s not found: %w", user.Account(), err))
		} else {
			err = allerror.New(allerror.ErrorCodeUserNotFound, "",
				xerrors.Errorf("failed to get user: %w", err))
		}
		return
	}

	id, err := strconv.Atoi(u.PlatformUser.Id)
	if err != nil {
		err = xerrors.Errorf("failed to get user: %w", err)
		return
	}

	//检查是否变化
	if err = cmd.IsChange(&u); err != nil {
		logrus.Errorf("nothing changed:%s", err.Error())
		return
	}

	//修改用户信息
	u.Email, err = domain.NewEmail(cmd.Account)
	if err != nil {
		e := xerrors.Errorf("failed to create email: %w", err)
		err = allerror.New(allerror.ErrorFailedToCreateEmail, "", e)
		return
	}
	//保存用户信息
	if u, err = s.repo.Save(&u); err != nil {
		e := xerrors.Errorf("failed to update user info: %w", err)
		err = allerror.New(allerror.ErrorFailedToUpdateUserInfo, "", e)
		return
	}
	//更新git用户信息
	if u.Email != nil && u.Email.Email() != "" {
		if err = s.UpdatePlatformAccountWithUpdate(id, &CreatePlatformAccountCmd{
			Account: u.Account,
			Email:   u.Email,
		}); err != nil {
			u.Email, _ = domain.NewEmail(cmd.OldAccount)
			u, err = s.repo.Save(&u)
			if err != nil {
				e := xerrors.Errorf("failed to update user info and fall back info too: %w", err)
				err = allerror.New(allerror.ErrorFailedToUpdateUserInfo, "", e)
				return
			}
			e := xerrors.Errorf("failed to update git user info: %w", err)
			err = allerror.New(allerror.ErrorFailedToUPdateGitUserInfo, "", e)
			return
		}
	}
	return
}

func (s userService) UpdatePlatformAccountWithUpdate(id int, cmd *CreatePlatformAccountCmd) (err error) {
	// create platform account
	dto, err := s.UpdatePlatformAccount(id, cmd)
	if err != nil {
		return
	}

	// update user information
	updatecmd := &UpdatePlateformInfoCmd{
		PlatformInfoDTO: dto,
		User:            cmd.Account,
		Email:           cmd.Email,
	}

	for i := 0; i <= 5; i++ {
		if err = s.UpdatePlateformInfo(updatecmd); err != nil {
			if !typerepo.IsErrorConcurrentUpdating(err) {
				return
			}
		} else {
			break
		}
	}

	return
}

func (s userService) UpdatePlatformAccount(id int, cmd *CreatePlatformAccountCmd) (dto PlatformInfoDTO, err error) {
	// create platform account
	pu, err := s.ps.Update(id, platform.UserOption{
		Email:    cmd.Email,
		Name:     cmd.Account,
		Password: cmd.Password,
	})
	if err != nil {
		return
	}

	dto.PlatformUser = pu

	// apply token
	token, err := s.ps.NewToken(pu)
	if err != nil {
		return
	}

	eToken, err := s.encryptToken(token.Token)
	if err != nil {
		return
	}

	dto.PlatformToken = domain.PlatformToken{
		Token:    eToken,
		CreateAt: token.CreateAt,
	}

	return
}
