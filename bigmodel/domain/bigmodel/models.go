package bigmodel

import (
	"io"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type CodeGeexReq struct {
	Lang    string
	Content string
}

type CodeGeexResp struct {
	Result string `json:"result"`
	Finish string `json:"finish"`
}

type BigModel interface {
	DescribePicture(io.Reader, string, int64) (string, error)
	GenPicture(types.Account, string) (string, error)
	GenPictures(types.Account, string) ([]string, error)
	Ask(domain.Question, string) (string, error)
	VQAUploadPicture(f io.Reader, u types.Account, fileName string) error
	LuoJiaUploadPicture(f io.Reader, u types.Account) error
	PanGu(string) (string, error)
	LuoJia(string) (string, error)
	CodeGeex(*CodeGeexReq) (CodeGeexResp, error)
	GetWuKongSampleId() string
	GenWuKongSampleNums(int) []int
	GenPicturesByWuKong(types.Account, *domain.WuKongPictureMeta) (map[string]string, error)
	DeleteWuKongPicture(string) error
	GenWuKongPictureLink(p string) (string, error)
	MoveWuKongPictureToDir(string, string) error
	GenWuKongLinkFromOBSPath(string) string
	CheckWuKongPictureTempToLike(types.Account, string) (domain.WuKongPictureMeta, string, error)
	CheckWuKongPicturePublicToLike(types.Account, string) (string, error)
	CheckWuKongPictureToPublic(types.Account, string) (domain.WuKongPictureMeta, string, error)
	GetIdleEndpoint(bid string) (c int, err error)
}
