package watchimpl

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type Config struct {
	Time TimeConfig `json:"time" required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Time,
	}
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	items := cfg.configItems()

	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type TimeConfig struct {
	ScanTime    int64 `json:"scan_time"`
	TriggerTime int64 `json:"trigger_time"`
}
