package domain

import (
	"errors"

	common "github.com/opensourceways/xihe-server/common/domain"
)

const (
	FieldEN = "English"
	FieldZH = "Chinese"

	PromotionStatusOver       = "over"
	PromotionStatusPreparing  = "preparing"
	PromotionStatusInProgress = "in-progress"

	originMindSpore = "MindSpore"
	originCSDN      = "CSDN"
	originOther     = "other"

	promotionWayOnline  = "online"
	promotionWayOffline = "offline"

	promotionTypeClockin      = "clockin"
	promotionTypeTrainingCamp = "trainingCamp"
	promotionTypeForum        = "forum"
	promotionTypeHackathon    = "hackathon"
	promotionTypeMsg          = "MSG"
	promotionTypeMindcon      = "MindCon"
	promotionTypeOther        = "other"
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

type PromotionStatus interface {
	PromotionStatus() string
}

type promotionStatus string

func NewPromotionStatus(v string) (PromotionStatus, error) {
	switch v {
	case PromotionStatusPreparing, PromotionStatusInProgress, PromotionStatusOver:
		return promotionStatus(v), nil
	case "":
		return nil, nil
	}

	return nil, errors.New("unsupported promotion status")
}

func (s promotionStatus) PromotionStatus() string {
	return string(s)
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

type Origin interface {
	Oringn() string
}

func NewOrigin(v string) (Origin, error) {
	switch v {
	case originMindSpore, originCSDN, originOther:
		return origin(v), nil
	}

	return nil, errors.New("invalid origin")
}

type origin string

func (o origin) Oringn() string {
	return string(o)
}

type PromotionWay interface {
	PromotionWay() string
}

func NewPromotionWay(v string) (PromotionWay, error) {
	switch v {
	case promotionWayOnline, promotionWayOffline:
		return promotionWay(v), nil
	case "":
		return nil, nil
	}

	return nil, errors.New("unsupported promotion way")
}

type promotionWay string

func (w promotionWay) PromotionWay() string {
	return string(w)
}

type PromotionType interface {
	PromotionType() string
}

func NewPromotionType(v string) (PromotionType, error) {
	switch v {
	case promotionTypeClockin, promotionTypeTrainingCamp, promotionTypeForum, promotionTypeHackathon, promotionTypeMsg,
		promotionTypeMindcon, promotionTypeOther:
		return promotionType(v), nil
	case "":
		return nil, nil
	}

	return nil, errors.New("unsupported promotion type")
}

type promotionType string

func (t promotionType) PromotionType() string {
	return string(t)
}
