package bigmodels

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

var reTimestamp = regexp.MustCompile("/[1-9][0-9]{9,}/")

type wukongInfo struct {
	cli         obsService
	cfg         WuKong
	maxBatch    int
	endpoints   chan string
	endpointsHF chan string
}

func newWuKongInfo(cfg *Config) (wukongInfo, error) {
	v := &cfg.WuKong

	cli, err := initOBS(&v.OBSAuthInfo)
	if err != nil {
		return wukongInfo{}, err
	}

	info := wukongInfo{
		cli:      cli,
		cfg:      *v,
		maxBatch: utils.LCM(v.SampleCount, v.SampleNum) / v.SampleNum,
	}

	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.WuKong)
	eshf, _ := ce.parse(ce.WuKongHF)

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	// init endpoints_hf
	info.endpointsHF = make(chan string, len(eshf))
	for _, e := range eshf {
		info.endpointsHF <- e
	}

	return info, nil
}

func (s *service) GenWuKongSampleNums(batchNum int) []int {
	cfg := &s.wukongInfo.cfg
	num := cfg.SampleNum
	count := cfg.SampleCount

	i := ((batchNum % s.wukongInfo.maxBatch) * num) % count

	r := make([]int, num)
	for j := 0; j < num; j++ {
		v := i + j
		if v >= count {
			v -= count
		}
		r[j] = v + 1
	}

	return r
}

func (s *service) GetWuKongSampleId() string {
	return s.wukongInfo.cfg.SampleId
}

func (s *service) GenPicturesByWuKong(
	user types.Account, desc *domain.WuKongPictureMeta, estype string,
) (map[string]string, error) {
	if err := s.check.check(desc.Desc.WuKongPictureDesc()); err != nil {
		return nil, err
	}

	var v []string

	f := func(e string) (err error) {
		v, err = s.genPicturesByWuKong(e, user, desc)

		return
	}

	// select endpoints
	var es chan string
	switch estype {
	case string(domain.BigmodelWuKong):
		es = s.wukongInfo.endpoints
	case string(domain.BigmodelWuKongHF):
		es = s.wukongInfo.endpointsHF
	}

	if err := s.doIfFree(es, f); err != nil {
		return nil, err
	}

	info := &s.wukongInfo
	bucket := info.cfg.Bucket
	expiry := info.cfg.DownloadExpiry

	r := map[string]string{}
	for _, p := range v {
		l, err := info.cli.genFileDownloadURL(bucket, p, expiry)
		if err != nil {
			return nil, err
		}
		r[p] = l
	}

	checkUrls := make([]string, len(r))

	var i int
	for _, v := range r {
		checkUrls[i] = v
		i++
	}
	if err := s.check.checkImages(checkUrls); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) genPicturesByWuKong(
	endpoint string, user types.Account, desc *domain.WuKongPictureMeta,
) ([]string, error) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return nil, err
	}

	opt := wukongRequest{
		Style: desc.Style,
		Desc:  desc.Desc.WuKongPictureDesc(),
		User:  user.Account(),
	}
	body, err := libutils.JsonMarshal(&opt)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	var r wukongResponse
	if _, err = s.hc.ForwardTo(req, &r); err != nil {
		return nil, err
	}

	if r.Status == 200 {
		return wukongLinksAC(r.Output), nil
	}

	return nil, errors.New(r.Msg)
}

func wukongLinksAC(v []string) []string {
	for i := range v {
		if strings.HasPrefix(v[i], "https") {
			v[i] = strings.Split(v[i], ".com/")[1]
		}
	}

	return v
}

func (s *service) MoveWuKongPictureToDir(dst, src string) error {
	info := &s.wukongInfo

	return info.cli.copyObject(info.cfg.Bucket, dst, src)
}

func (s *service) DeleteWuKongPicture(p string) error {
	info := &s.wukongInfo

	return info.cli.deleteObject(info.cfg.Bucket, p)
}

func (s *service) GenWuKongPictureLink(p string) (string, error) {
	info := &s.wukongInfo

	return info.cli.genFileDownloadURL(
		info.cfg.Bucket, p, info.cfg.DownloadExpiry,
	)
}

func (s *service) GenWuKongLinkFromOBSPath(obspath string) (link string) {
	cfg := s.wukongInfo.cfg
	endpoint := strings.Replace(cfg.Endpoint, "https://", "", 1)

	return fmt.Sprintf("https://%s.%s/%s", cfg.Bucket, endpoint, obspath)
}

func (s *service) CheckWuKongPictureTempToLike(user types.Account, p string) (
	meta domain.WuKongPictureMeta, path string, err error,
) {
	if meta, err = s.parseWuKongPictureMetaData(user, p); err == nil {
		v := user.Account()
		path = strings.Replace(p, v, v+"/like", 1)
	}

	return
}

func (s *service) CheckWuKongPicturePublicToLike(user types.Account, p string) (
	path string, err error,
) {
	v := user.Account()
	path = strings.Replace(p, "AI-gallery/gallery", "generate", 1)
	path = strings.Replace(path, "generate/"+v, "generate/"+v+"/like", 1)
	return
}

func (s *service) CheckWuKongPictureToPublic(user types.Account, p string) (
	meta domain.WuKongPictureMeta, path string, err error,
) {
	if strings.Contains(p, "/like") {
		path = strings.Replace(p, "/like", "", 1)
		path = strings.Replace(path, "generate", "AI-gallery/gallery", 1)
		return
	}
	if meta, err = s.parseWuKongPictureMetaData(user, p); err != nil {
		return
	}
	path = strings.Replace(p, "generate", "AI-gallery/gallery", 1)

	return
}

func (s *service) parseWuKongPictureMetaData(user types.Account, p string) (
	meta domain.WuKongPictureMeta, err error,
) {
	t := reTimestamp.FindString(p)
	if t == "" {
		err = errors.New("invalid path")

		return
	}

	v := strings.Split(p, "/"+user.Account()+t)
	if len(v) != 2 {
		err = errors.New("invalid path")

		return
	}

	desc := ""
	v = strings.Split(v[1], "/")
	switch len(v) {
	case 1:
		desc = v[0]
	case 2:
		meta.Style = strings.TrimSpace(v[0])
		desc = v[1]
	default:
		err = errors.New("invalid path")

		return
	}

	desc = strings.TrimSpace(strings.Split(desc, "-")[0])
	if meta.Style != "" {
		desc = strings.TrimSpace(strings.TrimSuffix(desc, meta.Style))
	}

	meta.Desc, err = domain.NewWuKongPictureDesc(desc)

	return
}

type wukongRequest struct {
	Style string `json:"style"`
	Desc  string `json:"desc"`
	User  string `json:"user_name"`
}

type wukongResponse struct {
	Status int      `json:"status"`
	Output []string `json:"output_image_url"`
	Msg    string   `json:"msg"`
}
