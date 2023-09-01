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
	// common
	GetIdleEndpoint(bid string) (c int, err error)
	CheckText(content string) error
	CheckImages(urls []string) error

	// wukong
	GetWuKongSampleId() string
	GenWuKongSampleNums(int) []int
	GenPicturesByWuKong(types.Account, *domain.WuKongPictureMeta, string) (map[string]string, error)
	DeleteWuKongPicture(string) error
	GenWuKongPictureLink(p string) (string, error)
	MoveWuKongPictureToDir(string, string) error
	GenWuKongLinkFromOBSPath(string) string
	CheckWuKongPictureTempToLike(types.Account, string) (domain.WuKongPictureMeta, string, error)
	CheckWuKongPicturePublicToLike(types.Account, string) (string, error)
	CheckWuKongPictureToPublic(types.Account, string) (domain.WuKongPictureMeta, string, error)

	// taichu
	DescribePicture(io.Reader, string, int64, string) (string, error)
	GenPicture(types.Account, string) (string, error)
	GenPictures(types.Account, string) ([]string, error)
	Ask(domain.Question, string) (string, error)
	VQAUploadPicture(f io.Reader, u types.Account, fileName string) error
	AskHF(f io.Reader, u types.Account, ask string) (string, error)

	// luojia
	LuoJiaUploadPicture(f io.Reader, u types.Account) error
	LuoJia(string) (string, error)
	LuoJiaHF(io.Reader) (string, error)

	// pangu
	PanGu(string) (string, error)

	// codegeex
	CodeGeex(*CodeGeexReq) (CodeGeexResp, error)

	// ai detector
	AIDetector(domain.AIDetectorInput) (bool, error)

	// baichuan2
	BaiChuan(*domain.BaiChuanInput) (string, string, error)
}
