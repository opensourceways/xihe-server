package config

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

func SetDefault(items []interface{}) {
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func Validate(items []interface{}) error {
	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
