package repositoryimpl

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
)

type userInfo struct {
	DUser `bson:",inline"`

	Count int `bson:"count"`
}

func toUserRegInfoDoc(u *domain.UserRegInfo, doc *DUserRegInfo) {
	*doc = DUserRegInfo{
		Account:  u.Account.Account(),
		Name:     u.Name.Name(),
		City:     u.City.City(),
		Email:    u.Email.Email(),
		Phone:    u.Phone.Phone(),
		Identity: u.Identity.Identity(),
		Province: u.Province.Province(),
		Detail:   u.Detail,
		Version:  u.Version,
	}
}

func (doc *DUserRegInfo) toUserRegInfo(u *domain.UserRegInfo) (err error) {
	if u.Account, err = types.NewAccount(doc.Account); err != nil {
		return
	}

	if u.Name, err = domain.NewName(doc.Name); err != nil {
		return
	}

	if u.City, err = domain.NewCity(doc.City); err != nil {
		return
	}

	if u.Email, err = domain.NewEmail(doc.Email); err != nil {
		return
	}

	if u.Phone, err = domain.NewPhone(doc.Phone); err != nil {
		return
	}

	if u.Identity, err = domain.NewIdentity(doc.Identity); err != nil {
		return
	}

	if u.Province, err = domain.NewProvince(doc.Province); err != nil {
		return
	}

	u.Detail = doc.Detail

	u.Version = doc.Version

	return
}

func toUserDoc(u domain.User, doc *DUser) {
	*doc = DUser{
		Name:                    u.Account.Account(),
		Email:                   u.Email.Email(),
		AvatarId:                u.AvatarId.AvatarId(),
		PlatformToken:           u.PlatformToken,
		PlatformUserId:          u.PlatformUser.Id,
		PlatformUserNamespaceId: u.PlatformUser.NamespaceId,
		UserAgreement:           u.UserAgreement,
		CourseAgreement:         u.CourseAgreement,
		FinetuneAgreement:       u.FinetuneAgreement,
	}

	if u.Bio != nil {
		doc.Bio = u.Bio.Bio()
	}
}

func toUser(doc DUser, u *domain.User) (err error) {

	if u.Email, err = domain.NewEmail(doc.Email); err != nil {
		return
	}

	if u.Account, err = domain.NewAccount(doc.Name); err != nil {
		return
	}

	if u.Bio, err = domain.NewBio(doc.Bio); err != nil {
		return
	}

	if u.AvatarId, err = domain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	u.Id = doc.Id.Hex()
	u.Version = doc.Version
	u.PlatformToken = doc.PlatformToken
	u.PlatformUser.Id = doc.PlatformUserId
	u.PlatformUser.NamespaceId = doc.PlatformUserNamespaceId
	u.CourseAgreement = doc.CourseAgreement
	u.FinetuneAgreement = doc.FinetuneAgreement
	u.UserAgreement = doc.UserAgreement

	return
}

func toUserInfo(doc DUser, info *domain.UserInfo) (err error) {

	if info.Account, err = domain.NewAccount(doc.Name); err != nil {
		return
	}

	if info.AvatarId, err = domain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	return
}

func toFollowerUserInfo(doc DUser, info *domain.FollowerUserInfo) (err error) {

	if info.Account, err = domain.NewAccount(doc.Name); err != nil {
		return
	}

	if info.AvatarId, err = domain.NewAvatarId(doc.AvatarId); err != nil {
		return
	}

	if info.Bio, err = domain.NewBio(doc.Bio); err != nil {
		return
	}

	return
}
