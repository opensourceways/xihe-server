package bigmodels

import (
	"errors"
	"net/url"
	"strings"
)

type Config struct {
	OBS        OBSConfig   `json:"obs"             required:"true"`
	Cloud      CloudConfig `json:"cloud"           required:"true"`
	WuKong     WuKong      `json:"wukong"          required:"true"`
	Endpoints  Endpoints   `json:"endpoints"       required:"true"`
	Moderation Moderation  `json:"moderation"      required:"true"`
	CloudGY    CloudConfig `json:"auth_gy"         required:"true"`

	MaxPictureSizeToDescribe int64 `json:"max_picture_size_to_describe"`
	MaxPictureSizeToVQA      int64 `json:"max_picture_size_to_vqa"`
}

func (cfg *Config) SetDefault() {
	cfg.WuKong.setDefault()

	if cfg.MaxPictureSizeToDescribe <= 0 {
		cfg.MaxPictureSizeToDescribe = 2 << 21
	}

	if cfg.MaxPictureSizeToVQA <= 0 {
		cfg.MaxPictureSizeToVQA = 2 << 21
	}
}

func (cfg *Config) Validate() error {
	if err := cfg.WuKong.validate(); err != nil {
		return err
	}

	return cfg.Endpoints.validate()
}

type OBSConfig struct {
	OBSAuthInfo

	VQABucket    string `json:"vqa_bucket"             required:"true"`
	LuoJiaBucket string `json:"luo_jia_bucket"         required:"true"`
}

type OBSAuthInfo struct {
	Endpoint  string `json:"endpoint"                  required:"true"`
	AccessKey string `json:"access_key"                required:"true"`
	SecretKey string `json:"secret_key"                required:"true"`
}

type CloudConfig struct {
	Domain       string `json:"domain"                 required:"true"`
	User         string `json:"user"                   required:"true"`
	Password     string `json:"password"               required:"true"`
	Project      string `json:"project"                required:"true"`
	AuthEndpoint string `json:"auth_endpoint"          required:"true"`
}

type Endpoints struct {
	VQA              string `json:"vqa"                required:"true"`
	VQAHF            string `json:"vqa_hf"             required:"true"`
	Pangu            string `json:"pangu"              required:"true"`
	LuoJia           string `json:"luojia"             required:"true"`
	LuoJiaHF         string `json:"luojia_hf"          required:"true"`
	WuKong           string `json:"wukong"             required:"true"`
	WuKong4IMG       string `json:"wukong_4img"        required:"true"`
	WuKongHF         string `json:"wukong_hf"          required:"true"`
	WuKongUser       string `json:"wukong_user"        required:"true"`
	CodeGeex         string `json:"codegeex"           required:"true"`
	DescPicture      string `json:"desc_picture"       required:"true"`
	DescPictureHF    string `json:"desc_picture_hf"    required:"true"`
	SinglePicture    string `json:"single_picture"     required:"true"`
	MultiplePictures string `json:"multiple_pictures"  required:"true"`
	AIDetector       string `json:"ai_detector"        required:"true"`
	BaiChuan         string `json:"baichuan"           required:"true"`
	GLM2             string `json:"glm"                required:"true"`
	LLAMA2           string `json:"llama"              required:"true"`
	SkyWork          string `json:"skywork"            required:"true"`
	IFlytekspark     string `json:"iflytekspark"       required:"true"`
	IFlyteksparkLong string `json:"iflytekspark_long"  required:"true"`
}

func (e *Endpoints) validate() (err error) {
	if _, err = e.parse(e.VQA); err != nil {
		return
	}

	if _, err = e.parse(e.Pangu); err != nil {
		return
	}

	if _, err = e.parse(e.LuoJia); err != nil {
		return
	}

	if _, err = e.parse(e.WuKong); err != nil {
		return
	}

	if _, err = e.parse(e.WuKongHF); err != nil {
		return
	}

	if _, err = e.parse(e.DescPicture); err != nil {
		return
	}

	if _, err = e.parse(e.DescPictureHF); err != nil {
		return
	}

	if _, err = e.parse(e.SinglePicture); err != nil {
		return
	}

	if _, err = e.parse(e.MultiplePictures); err != nil {
		return
	}

	if _, err = e.parse(e.AIDetector); err != nil {
		return
	}

	if _, err = e.parse(e.BaiChuan); err != nil {
		return
	}

	if _, err = e.parse(e.GLM2); err != nil {
		return
	}

	if _, err = e.parse(e.GLM2); err != nil {
		return
	}

	if _, err = e.parse(e.LLAMA2); err != nil {
		return
	}

	if _, err = e.parse(e.SkyWork); err != nil {
		return
	}

	return
}

func (e *Endpoints) parse(s string) ([]string, error) {
	v := strings.Split(strings.Trim(s, ","), ",")

	for _, i := range v {
		if _, err := url.Parse(i); err != nil {
			return nil, errors.New("invalid url")
		}
	}

	if len(v) == 0 {
		return nil, errors.New("missing endpoints")
	}

	return v, nil
}

type Moderation struct {
	Endpoint    string `json:"endpoint"       required:"true"`
	AccessKey   string `json:"access_key"     required:"true"`
	SecretKey   string `json:"secret_key"     required:"true"`
	IAMEndpoint string `json:"iam_endpoint"   required:"true"`
	Region      string `json:"region"         required:"true"`
}

type WuKong struct {
	WuKongSample
	CloudConfig
	OBSAuthInfo

	Bucket string `json:"bucket"             required:"true"`

	// DownloadExpiry specifies the timeout to download a obs file.
	// The unit is second.
	DownloadExpiry int `json:"download_expiry"`
}

type WuKongSample struct {
	SampleId    string `json:"sample_id"     required:"true"`
	SampleNum   int    `json:"sample_num"    required:"true"`
	SampleCount int    `json:"sample_count"  required:"true"`
}

func (cfg *WuKong) setDefault() {
	if cfg.DownloadExpiry <= 0 {
		cfg.DownloadExpiry = 3600
	}
}

func (cfg *WuKong) validate() error {
	if cfg.SampleNum > cfg.SampleCount {
		return errors.New("make sure that sample_num <= sample_count")
	}

	if cfg.SampleNum <= 0 {
		return errors.New("invalid sample_num")
	}

	if cfg.SampleCount <= 0 {
		return errors.New("invalid sample_count")
	}

	return nil
}

type ApiService struct {
	TokenExpire string
}
