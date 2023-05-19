package service

import (
	"net/url"
	"sort"
	"strings"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type BigModelService interface {
	// wukong
	IsLike(*domain.WuKongPicture, types.Account) (bool, string, error)
	IsPublic(*domain.WuKongPicture) (bool, string, error)
	IsDigg(types.Account, []string) bool
	LinkLikePublic(string, types.Account) (LinkLikePublicOpt, error)
	IsPathCotain(string, []domain.WuKongPicture) bool

	// luojia
	LatestLuojiaList([]domain.LuoJiaRecord) domain.LuoJiaRecord
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
		likePath, err = s.fm.CheckWuKongPicturePublicToLike(user, p.OBSPath.OBSPath())
		if err != nil {
			return
		}

		if pic.OBSPath.OBSPath() == likePath {
			return true, pic.Id, nil
		}
	}

	return
}

func (s *bigModelService) IsPublic(
	p *domain.WuKongPicture,
) (isPublic bool, publicId string, err error) {
	pics, _, err := s.wukongPicture.ListPublicsByUserName(p.Owner)
	if err != nil {
		return
	}

	for _, pic := range pics {
		var publicPath string
		_, publicPath, err = s.fm.CheckWuKongPictureToPublic(p.Owner, p.OBSPath.OBSPath())
		if err != nil {
			return
		}

		if pic.OBSPath.OBSPath() == publicPath {
			isPublic = true
			publicId = pic.Id

			return
		}
	}

	return
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
	opt LinkLikePublicOpt, err error,
) {
	obspath := toOBSPath(link)
	op, err := domain.NewOBSPath(obspath)
	if err != nil {
		return
	}

	p := domain.WuKongPicture{
		OBSPath: op,
		Owner:   user,
	}

	if opt.IsLike, opt.LikeId, err = s.IsLike(&p, user); err != nil {
		return
	}

	if opt.IsPublic, opt.PublicId, err = s.IsPublic(&p); err != nil {
		return
	}

	return
}

func (s *bigModelService) IsPathCotain(path string, v []domain.WuKongPicture) bool {
	for i := range v {
		if v[i].OBSPath.OBSPath() == path {
			return true
		}
	}

	return false
}

func toOBSPath(link string) (obspath string) {
	u, _ := url.QueryUnescape(link)

	t := strings.Split(u, ".ovaijisuan.com:443/")[1]
	obspath = strings.Split(t, "?AWSAccessKeyId")[0]

	return
}

func (s *bigModelService) LatestLuojiaList(v []domain.LuoJiaRecord) (r domain.LuoJiaRecord) {
	sort.Slice(v, func(i, j int) bool {
		return v[i].CreatedAt < v[j].CreatedAt
	})

	if len(v) > 0 {
		return v[len(v)-1]
	}

	return
}
