package domain

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/utils"
)

// FinetuneName
type FinetuneName interface {
	FinetuneName() string
}

func NewFinetuneName(v string) (FinetuneName, error) {
	v = utils.XSSFilter(v)

	max := DomainConfig.MaxFinetuneNameLength
	min := DomainConfig.MinFinetuneNameLength

	if n := utils.StrLen(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return finetuneName(v), nil
}

type finetuneName string

func (r finetuneName) FinetuneName() string {
	return string(r)
}
