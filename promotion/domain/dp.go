package domain

import (
	"errors"
	"strings"
	"time"

	common "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	FieldEN = "English"
	FieldZH = "Chinese"

	promotionStatusOver       = "over"
	promotionStatusPreparing  = "preparing"
	promotionStatusInProgress = "in-progress"
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
	PromotionStatus() (string, error)
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

func (r promotionDuration) PromotionStatus() (string, error) {
	start, end, err := r.durationTime()
	if err != nil {
		return "", errors.New("parse duration error")
	}

	now := time.Now().Unix()
	if now < start {
		return promotionStatusPreparing, nil
	}
	if now >= start && now <= end {
		return promotionStatusInProgress, nil
	}
	if now > end {
		return promotionStatusOver, nil
	}

	return "", errors.New("promotion status internal error")
}

func (r promotionDuration) durationTime() (int64, int64, error) {
	t := strings.Split(string(r), "-")

	start, err := utils.ToUnixTimeLayout2(t[0])
	if err != nil {
		return 0, 0, err
	}

	end, err := utils.ToUnixTimeLayout2(t[1])
	if err != nil {
		return 0, 0, err
	}

	return start.Unix(), end.Unix(), nil
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
