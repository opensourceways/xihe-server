package inferenceimpl

type Config struct {
	ContainerManagerEndpoint string `json:"endpoint"  required:"true"`

	// unit is second
	SurvivalTimeForNormal int `json:"survival_time_for_normal"`

	// unit is second
	SurvivalTimeForOfficial int `json:"survival_time_for_official"`

	ProjectTagsForOfficial []string `json:"project_tags_for_official" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.SurvivalTimeForNormal <= 0 {
		cfg.SurvivalTimeForNormal = 5 * 3600
	}

	if cfg.SurvivalTimeForOfficial <= 0 {
		cfg.SurvivalTimeForOfficial = 12 * 3600
	}

}
