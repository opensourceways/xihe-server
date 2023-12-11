package domain

import (
	"errors"

	common "github.com/opensourceways/xihe-server/common/domain"
)

const (
	FieldEN = "English"
	FieldZH = "Chinese"
)

// PronotionName
type PromotionName interface {
	PromotionName() string
}

func NewPromotionName(v string) (PromotionName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return promotionName(v), nil
}

type promotionName string

func (r promotionName) PromotionName() string {
	return string(r)
}

// PromotionDuration
type PromotionDuration interface {
	PromotionDuration() string
}

func NewPromotionDuration(v string) (PromotionDuration, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return promotionDuration(v), nil
}

type promotionDuration string

func (r promotionDuration) PromotionDuration() string {
	return string(r)
}

// PromotionDesc
type PromotionDesc interface {
	PromotionDesc() string
}

func NewPromotionDesc(v string) (PromotionDesc, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return promotionDesc(v), nil
}

type promotionDesc string

func (r promotionDesc) PromotionDesc() string {
	return string(r)
}

// // Dones
// type Dones interface {
// 	Dones() []string
// }

// func NewDones(d []string) (Dones, error) {
// 	if len(d) <= 0 {
// 		return nil, errors.New("invalid dones input")
// 	}

// 	return dones(d), nil
// }

// func (r dones) Dones() []string {
// 	return r
// }

// type dones []string

// Sentence
type Sentence interface {
	Sentence(common.Language) string
	SentenceMap() map[string]string
	ENSentence() string
	ZHSentence() string
}

func NewSentence(en, zh string) (Sentence, error) {
	if en == "" || zh == "" {
		return nil, errors.New("empty value")
	}

	s := make(sentence, 2)
	s[FieldEN] = en
	s[FieldZH] = zh

	return s, nil
}

type sentence map[string]string

func (r sentence) Sentence(lang common.Language) string {
	switch lang.Language() {
	case FieldEN:
		return r.ENSentence()
	case FieldZH:
		return r.ZHSentence()
	}

	return r.ZHSentence()
}

func (r sentence) SentenceMap() map[string]string {
	m := make(map[string]string, 2)
	m[FieldEN] = r.ENSentence()
	m[FieldZH] = r.ZHSentence()
	return m
}

func (r sentence) ENSentence() string {
	return r[FieldEN]
}

func (r sentence) ZHSentence() string {
	return r[FieldZH]
}
