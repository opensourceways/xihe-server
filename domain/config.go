package domain

import (
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"
)

var DomainConfig Config

func Init(cfg *Config) {
	DomainConfig = *cfg
}

type Config struct {
	covers           sets.Set[string]
	protocols        sets.Set[string]
	projectType      sets.Set[string]
	trainingPlatform sets.Set[string]
	avatarURL        sets.Set[string]

	MaxBioLength          int `json:"max_bio_length"`
	MaxNameLength         int `json:"max_name_length"`
	MinNameLength         int `json:"min_name_length"`
	MaxTitleLength        int `json:"max_title_length"`
	MinTitleLength        int `json:"min_title_length"`
	MaxDescLength         int `json:"max_desc_length"`
	MaxNicknameLength     int `json:"max_nickname_length"`
	MaxRelatedResourceNum int `json:"max_related_resource_num"`

	Covers           []string `json:"covers"            required:"true"`
	Protocols        []string `json:"protocols"         required:"true"`
	ProjectType      []string `json:"project_type"      required:"true"`
	TrainingPlatform []string `json:"training_platform" required:"true"`
	AvatarURL        []string `json:"avatar_url"        required:"true"`

	MaxTrainingNameLength int `json:"max_training_name_length"`
	MinTrainingNameLength int `json:"min_training_name_length"`
	MaxTrainingDescLength int `json:"max_training_desc_length"`

	MaxFinetuneNameLength int `json:"max_finetune_name_length"`
	MinFinetuneNameLength int `json:"min_finetune_name_length"`

	WuKongPictureMaxDescLength int `json:"wukong_picture_max_desc_length"`

	// Key is the finetue model name
	Finetunes map[string]FinetuneParameterConfig `json:"finetunes"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxNameLength <= 0 {
		cfg.MaxNameLength = 50
	}

	if cfg.MinNameLength <= 0 {
		cfg.MinNameLength = 5
	}

	if cfg.MaxTitleLength <= 0 {
		cfg.MaxTitleLength = 100
	}

	if cfg.MaxDescLength <= 0 {
		cfg.MaxDescLength = 200
	}

	if cfg.MaxRelatedResourceNum <= 0 {
		cfg.MaxRelatedResourceNum = 5
	}

	if cfg.MaxNicknameLength == 0 {
		cfg.MaxNicknameLength = 20
	}

	if cfg.MaxBioLength == 0 {
		cfg.MaxBioLength = 200
	}

	if cfg.MaxTrainingNameLength == 0 {
		cfg.MaxTrainingNameLength = 50
	}

	if cfg.MinTrainingNameLength == 0 {
		cfg.MinTrainingNameLength = 5
	}

	if cfg.MaxFinetuneNameLength == 0 {
		cfg.MaxFinetuneNameLength = 25
	}

	if cfg.MinFinetuneNameLength == 0 {
		cfg.MinFinetuneNameLength = 1
	}

	if cfg.MaxTrainingDescLength == 0 {
		cfg.MaxTrainingDescLength = 100
	}

	if cfg.WuKongPictureMaxDescLength <= 0 {
		cfg.WuKongPictureMaxDescLength = 75
	}

	if cfg.Finetunes == nil {
		cfg.Finetunes = map[string]FinetuneParameterConfig{}
	}
}

func (r *Config) Validate() error {
	if r.MaxNameLength < (r.MinNameLength + 10) {
		return errors.New("invalid name length")
	}

	r.covers = sets.New[string](r.Covers...)
	r.protocols = sets.New[string](r.Protocols...)
	r.projectType = sets.New[string](r.ProjectType...)
	r.trainingPlatform = sets.New[string](r.TrainingPlatform...)
	r.avatarURL = sets.New[string](r.AvatarURL...)

	return nil
}

func (cfg *Config) hasCover(v string) bool {
	return cfg.covers.Has(v)
}

func (cfg *Config) hasProtocol(v string) bool {
	return cfg.protocols.Has(v)
}

func (cfg *Config) hasProjectType(v string) bool {
	return cfg.projectType.Has(v)
}

func (cfg *Config) hasPlatform(v string) bool {
	return cfg.trainingPlatform.Has(v)
}

func (cfg *Config) HasAvatarURL(v string) bool {
	return cfg.avatarURL.Has(v)
}

type FinetuneParameterConfig struct {
	Tasks           []string `json:"tasks"           required:"true"`
	Hyperparameters []string `json:"hyperparameters" required:"true"`
}
