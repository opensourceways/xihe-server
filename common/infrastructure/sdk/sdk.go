package sdk

import (
	"github.com/opensourceways/xihe-audit-sync-sdk/httpclient"
)

// Init is for http client init
func Init(cfg *Config) {
	httpclient.Init(cfg)
}

// Config is for http client config
type Config = httpclient.Config
