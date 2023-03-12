package controller

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	resourceTagTypeModel         = "model"
	resourceTagTypeProject       = "project"
	resourceTagTypeDataset       = "dataset"
	resourceTagTypeGlobalModel   = "global_model"
	resourceTagTypeGlobalProject = "global_project"
	resourceTagTypeGlobalDataset = "global_dataset"
)

var (
	apiConfig     APIConfig
	encryptHelper utils.SymmetricEncryption
	log           *logrus.Entry
)

func Init(cfg *APIConfig, l *logrus.Entry) error {
	log = l
	apiConfig = *cfg

	e, err := utils.NewSymmetricEncryption(cfg.EncryptionKey, "")
	if err != nil {
		return err
	}

	encryptHelper = e

	return nil
}

type APIConfig struct {
	Tags                           Tags   `json:"tags"                        required:"true"`
	TokenKey                       string `json:"token_key"                   required:"true"`
	TokenExpiry                    int64  `json:"token_expiry"                required:"true"`
	EncryptionKey                  string `json:"encryption_key"              required:"true"`
	DefaultPassword                string `json:"default_password"            required:"true"`
	MaxTrainingRecordNum           int    `json:"max_training_record_num"     required:"true"`
	InferenceDir                   string `json:"inference_dir"`
	InferenceBootFile              string `json:"inference_boot_file"`
	InferenceTimeout               int    `json:"inference_timeout"`
	EvaluateTimeout                int    `json:"evaluate_timeout"`
	PodTimeout                     int    `json:"pod_timeout"`
	MaxPictureSizeToDescribe       int64  `json:"-"`
	MaxPictureSizeToVQA            int64  `json:"-"`
	MinSurvivalTimeOfEvaluate      int    `json:"min_survival_time_of_evaluate"`
	MinSurvivalTimeOfInference     int    `json:"min_survival_time_of_inference"`
	MaxTagsNumToSearchResource     int    `json:"max_tags_num_to_search_resource"`
	MaxTagKindsNumToSearchResource int    `json:"max_tag_kinds_num_to_search_resource"`
}

func (cfg *APIConfig) SetDefault() {
	if cfg.MinSurvivalTimeOfInference <= 0 {
		cfg.MinSurvivalTimeOfInference = 3600
	}

	if cfg.MinSurvivalTimeOfEvaluate <= 0 {
		cfg.MinSurvivalTimeOfEvaluate = 3600
	}

	if cfg.InferenceDir == "" {
		cfg.InferenceDir = "inference"
	}

	if cfg.InferenceBootFile == "" {
		cfg.InferenceBootFile = "inference/app.py"
	}

	if cfg.InferenceTimeout <= 0 {
		cfg.InferenceTimeout = 300
	}

	if cfg.EvaluateTimeout <= 0 {
		cfg.EvaluateTimeout = 300
	}

	if cfg.PodTimeout <= 0 {
		cfg.PodTimeout = 300
	}

	if cfg.MaxTagsNumToSearchResource <= 0 {
		cfg.MaxTagsNumToSearchResource = 5
	}

	if cfg.MaxTagKindsNumToSearchResource <= 0 {
		cfg.MaxTagKindsNumToSearchResource = 5
	}
}

func (cfg *APIConfig) Validate() (err error) {
	if _, err = domain.NewPassword(cfg.DefaultPassword); err != nil {
		return
	}

	if _, err = domain.NewDirectory(cfg.InferenceDir); err != nil {
		return
	}

	_, err = domain.NewFilePath(cfg.InferenceBootFile)

	return
}

type Tags struct {
	ModelTagDomains         []string `json:"model"            required:"true"`
	ProjectTagDomains       []string `json:"project"          required:"true"`
	DatasetTagDomains       []string `json:"dataset"          required:"true"`
	GlobalModelTagDomains   []string `json:"global_model"     required:"true"`
	GlobalProjectTagDomains []string `json:"global_project"   required:"true"`
	GlobalDatasetTagDomains []string `json:"global_dataset"   required:"true"`
}

func (t *Tags) getDomains(name string) []string {
	switch name {
	case resourceTagTypeModel:
		return t.ModelTagDomains

	case resourceTagTypeProject:
		return t.ProjectTagDomains

	case resourceTagTypeDataset:
		return t.DatasetTagDomains

	case resourceTagTypeGlobalModel:
		return t.GlobalModelTagDomains

	case resourceTagTypeGlobalProject:
		return t.GlobalProjectTagDomains

	case resourceTagTypeGlobalDataset:
		return t.GlobalDatasetTagDomains

	default:
		return nil
	}
}
