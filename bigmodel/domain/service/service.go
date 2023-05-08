package service

import (
	"net/url"
	"strings"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type BigModelService interface {
	IsLike(*domain.WuKongPicture, types.Account) (bool, string, error)
	IsPublic(*domain.WuKongPicture) (bool, error)
	IsDigg(types.Account, []string) bool
	LinkLikePublic(string, types.Account) (bool, string, bool, error)
}

type bigModelService struct {
	fm            bigmodel.BigModel
	wukongPicture repository.WuKongPicture
}

func NewBigModelService(
	fm bigmodel.BigModel,
	wukongPicture repository.WuKongPicture,
) BigModelService {
	return &bigModelService{
		fm:            fm,
		wukongPicture: wukongPicture,
	}
}

func (s *bigModelService) IsLike(
	p *domain.WuKongPicture,
	user types.Account,
) (isLike bool, id string, err error) {
	pics, _, err := s.wukongPicture.ListLikesByUserName(user)
	if err != nil {
		return
	}

	for _, pic := range pics {
		var likePath string
		likePath, err = s.fm.CheckWuKongPicturePublicToLike(user, p.OBSPath)
		if err != nil {
			return
		}

		if pic.OBSPath == likePath {
			return true, pic.Id, nil
		}
	}

	return
}

func (s *bigModelService) IsPublic(
	p *domain.WuKongPicture,
) (bool, error) {
	pics, _, err := s.wukongPicture.ListPublicsByUserName(p.Owner)
	if err != nil {
		return false, err
	}

	for _, pic := range pics {
		_, publicPath, err := s.fm.CheckWuKongPictureToPublic(p.Owner, p.OBSPath)
		if err != nil {
			return false, err
		}

		if pic.OBSPath == publicPath {
			return true, nil
		}
	}

	return false, nil
}

func (s *bigModelService) IsDigg(
	user types.Account,
	diggs []string,
) bool {
	for _, username := range diggs {
		if user.Account() == username {
			return true
		}
	}

	return false
}

func (s *bigModelService) LinkLikePublic(link string, user types.Account) (
	isLike bool, likeId string, isPublic bool, err error,
) {
	obspath := toOBSPath(link)

	p := domain.WuKongPicture{
		OBSPath: obspath,
		Owner:   user,
	}

	if isLike, likeId, err = s.IsLike(&p, user); err != nil {
		return
	}

	if isPublic, err = s.IsPublic(&p); err != nil {
		return
	}

	return
}

func toOBSPath(link string) (obspath string) {
	u, _ := url.QueryUnescape(link)

	t := strings.Split(u, ".ovaijisuan.com:443/")[1]
	obspath = strings.Split(t, "?AWSAccessKeyId")[0]

	return
}
