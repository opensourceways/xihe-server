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

	// luojia
	LuoJiaUploadPicture(f io.Reader, u types.Account) error
	LuoJia(string) (string, error)
	LuoJiaHF(io.Reader) (string, error)

	// ai detector
	AIDetector(domain.AIDetectorInput) (bool, error)

	// baichuan2
	BaiChuan(*domain.BaiChuanInput) (string, string, error)

	// glm2
	GLM2(chan string, *domain.GLM2Input) error

	// llama2
	LLAMA2(chan string, *domain.LLAMA2Input) error

	// skywork
	SkyWork(chan string, *domain.SkyWorkInput) error
}
