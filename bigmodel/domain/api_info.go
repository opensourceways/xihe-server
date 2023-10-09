package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type ApiInfo struct {
	Id       string
	Name     string
	Endpoint string
	Doc      types.URL
}
