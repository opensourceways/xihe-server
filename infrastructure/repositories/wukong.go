package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongPictureListOptionDO = repository.WuKongPictureListOption

type WuKongMapper interface {
	ListSamples(string, []int) ([]string, error)
}

func NewWuKongRepository(mapper WuKongMapper) repository.WuKong {
	return wukong{mapper}
}

type wukong struct {
	mapper WuKongMapper
}

func (impl wukong) ListSamples(sid string, nums []int) ([]string, error) {
	return impl.mapper.ListSamples(sid, nums)
}

type WuKongPictureMetaDO struct {
	Style string
	Desc  string
}

func (do *WuKongPictureMetaDO) toWuKongPictureMeta() (
	meta domain.WuKongPictureMeta, err error,
) {
	meta.Style = do.Style
	meta.Desc, err = domain.NewWuKongPictureDesc(do.Desc)

	return
}
