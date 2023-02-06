package bigmodels

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

var reTimestamp = regexp.MustCompile("/[1-9][0-9]{9,}/")

type wukongInfo struct {
	cli       obsService
	cfg       WuKong
	maxBatch  int
	endpoints chan string
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

	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
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
	user domain.Account, desc *domain.WuKongPictureMeta,
) (map[string]string, error) {
	if err := s.check.check(desc.Desc.WuKongPictureDesc()); err != nil {
		return nil, err
	}

	var v []string

	f := func(e string) (err error) {
		v, err = s.genPicturesByWuKong(e, user, desc)

		return
	}

	if err := s.doIfFree(s.wukongInfo.endpoints, f); err != nil {
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

	return r, nil
}

func (s *service) genPicturesByWuKong(
	endpoint string, user domain.Account, desc *domain.WuKongPictureMeta,
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
		return r.Output, nil
	}

	return nil, errors.New(r.Msg)
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
	// fmt.Printf("s.wukongInfo.cfg.Endpoint: %v\n", s.wukongInfo.cfg.Endpoint)
	return fmt.Sprintf("https://%s.%s/%s", cfg.Bucket, cfg.Endpoint, obspath)
}

func (s *service) CheckWuKongPictureTempToLike(user domain.Account, p string) (
	meta domain.WuKongPictureMeta, path string, err error,
) {
	if meta, err = s.parseWuKongPictureMetaData(user, p); err == nil {
		v := user.Account()
		path = strings.Replace(p, v, v+"/like", 1)
	}

	return
}

func (s *service) CheckWuKongPicturePublicToLike(user domain.Account, p string) (
	path string, err error,
) {
	v := user.Account()
	path = strings.Replace(p, "AI-gallery/gallery", "generate", 1)
	path = strings.Replace(path, "generate/"+v, "generate/"+v+"/like", 1)
	return
}

func (s *service) CheckWuKongPictureToPublic(user domain.Account, p string) (
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

func (s *service) parseWuKongPictureMetaData(user domain.Account, p string) (
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
