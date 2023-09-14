package config

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type configItems interface {
	ConfigItems() []interface{}
}

func SetDefault(cfg interface{}) {
	if f, ok := cfg.(configSetDefault); ok {
		f.SetDefault()

		return
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			SetDefault(items[i])
		}
	}
}

func Validate(cfg interface{}) error {
	if f, ok := cfg.(configValidate); ok {
		return f.Validate()
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			if err := Validate(items[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
