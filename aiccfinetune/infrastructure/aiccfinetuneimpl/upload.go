package aiccfinetuneimpl

import "io"

var cs *service

func Init(cfg *Config) error {
	obs, err := initOBS(&cfg.OBSConfig)
	if err != nil {
		return err
	}

	cs = &service{
		obs: obs,
	}

	return nil
}

func NewAICCUploadService() *service {
	return cs
}

type service struct {
	obs obsService
}

func (cs *service) UploadAICC(data io.Reader, path string) error {
	return cs.obs.createObject(data, path)
}
