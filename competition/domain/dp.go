package domain

import (
	"errors"
	"net/url"

	"github.com/opensourceways/xihe-server/utils"
)

const (
	competitionTeamRoleLeader = "leader"

	competitionTypeChallenge = "challenge"

	competitionPhaseFinal       = "final"
	competitionPhasePreliminary = "preliminary"

	competitionStatusOver       = "over"
	competitionStatusPreparing  = "preparing"
	competitionStatusInProgress = "in-progress"

	competitionIdentityStudent   = "student"
	competitionIdentityTeacher   = "teacher"
	competitionIdentityDeveloper = "developer"

	competitionSubmissionStatusSuccess = "success"

	competitionTagElectricity = "electricity"
	competitionTagLearn       = "learn"
	competitionTagChallenge   = "challenge"
)

var (
	CompetitionPhaseFinal       = competitionPhase("final")
	CompetitionPhasePreliminary = competitionPhase("preliminary")
)

// CompetitionType
type CompetitionType interface {
	CompetitionType() string
}

func NewCompetitionType(v string) (CompetitionType, error) {
	if v == "" || v == competitionTypeChallenge {
		return competitionType(v), nil
	}

	return nil, errors.New("invalid competition type")
}

type competitionType string

func (r competitionType) CompetitionType() string {
	return string(r)
}

// CompetitionPhase
type CompetitionPhase interface {
	CompetitionPhase() string
	IsFinal() bool
	IsPreliminary() bool
}

func NewCompetitionPhase(v string) (CompetitionPhase, error) {
	if v == competitionPhasePreliminary || v == competitionPhaseFinal {
		return competitionPhase(v), nil
	}

	return nil, errors.New("invalid competition phase")
}

type competitionPhase string

func (r competitionPhase) CompetitionPhase() string {
	return string(r)
}

func (r competitionPhase) IsFinal() bool {
	return string(r) == competitionPhaseFinal
}

func (r competitionPhase) IsPreliminary() bool {
	return string(r) == competitionPhasePreliminary
}

// CompetitionStatus
type CompetitionStatus interface {
	CompetitionStatus() string
	IsOver() bool
}

func NewCompetitionStatus(v string) (CompetitionStatus, error) {
	b := v == competitionStatusPreparing ||
		v == competitionStatusInProgress ||
		v == competitionStatusOver

	if b {
		return competitionStatus(v), nil
	}

	return nil, errors.New("invalid competition status")
}

type competitionStatus string

func (r competitionStatus) CompetitionStatus() string {
	return string(r)
}

func (r competitionStatus) IsOver() bool {
	return string(r) == competitionStatusOver
}

// CompetitionName
type CompetitionName interface {
	CompetitionName() string
}

func NewCompetitionName(v string) (CompetitionName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionName(v), nil
}

type competitionName string

func (r competitionName) CompetitionName() string {
	return string(r)
}

// CompetitionDesc
type CompetitionDesc interface {
	CompetitionDesc() string
}

func NewCompetitionDesc(v string) (CompetitionDesc, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionDesc(v), nil
}

type competitionDesc string

func (r competitionDesc) CompetitionDesc() string {
	return string(r)
}

// CompetitionDuration
type CompetitionDuration interface {
	CompetitionDuration() string
}

func NewCompetitionDuration(v string) (CompetitionDuration, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionDuration(v), nil
}

type competitionDuration string

func (r competitionDuration) CompetitionDuration() string {
	return string(r)
}

// CompetitionBonus
type CompetitionBonus interface {
	CompetitionBonus() int
}

func NewCompetitionBonus(v int) (CompetitionBonus, error) {
	if v < 0 {
		return nil, errors.New("invalid bonus")
	}

	return competitionBonus(v), nil
}

type competitionBonus int

func (r competitionBonus) CompetitionBonus() int {
	return int(r)
}

type CompetitionHost interface {
	CompetitionHost() string
}

func NewCompetitionHost(v string) (CompetitionHost, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionHost(v), nil
}

type competitionHost string

func (r competitionHost) CompetitionHost() string {
	return string(r)
}

// Forum
type Forum interface {
	Forum() string
}

func NewForum(v string) (Forum, error) {
	if v == "" {
		return dpForum(v), nil
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpForum(v), nil
}

type dpForum string

func (r dpForum) Forum() string {
	return string(r)
}

// Winners
type Winners interface {
	Winners() string
}

func NewWinners(v string) (Winners, error) {
	if v == "" {
		return dpWinners(v), nil
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpWinners(v), nil
}

type dpWinners string

func (r dpWinners) Winners() string {
	return string(r)
}

// URL
type URL interface {
	URL() string
}

func NewURL(v string) (URL, error) {
	if v == "" {
		return nil, errors.New("empty url")
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpURL(v), nil
}

type dpURL string

func (r dpURL) URL() string {
	return string(r)
}

// Phone
type Phone interface {
	Phone() string
}

func NewPhone(v string) (Phone, error) {
	if len(v) > 0 && !utils.IsChinesePhone(v) {
		return nil, errors.New("invalid phone number")
	}

	return phone(v), nil
}

type phone string

func (r phone) Phone() string {
	return string(r)
}

// CompetitionIdentity
type CompetitionIdentity interface {
	CompetitionIdentity() string
}

func NewcompetitionIdentity(v string) (CompetitionIdentity, error) {
	b := v == competitionIdentityStudent ||
		v == competitionIdentityTeacher ||
		v == competitionIdentityDeveloper ||
		v == ""

	if !b {
		return nil, errors.New("invalid competition identity")
	}

	return competitionIdentity(v), nil
}

type competitionIdentity string

func (r competitionIdentity) CompetitionIdentity() string {
	return string(r)
}

// Province
type Province interface {
	Province() string
}

func NewProvince(v string) (Province, error) {
	v = utils.XSSFilter(v)

	if utils.StrLen(v) > 15 {
		return nil, errors.New("invalid province")
	}

	return province(v), nil
}

type province string

func (r province) Province() string {
	return string(r)
}

// City
type City interface {
	City() string
}

func NewCity(v string) (City, error) {
	v = utils.XSSFilter(v)

	if utils.StrLen(v) > 20 {
		return nil, errors.New("invalid city")
	}

	return city(v), nil
}

type city string

func (r city) City() string {
	return string(r)
}

// CompetitorName
type CompetitorName interface {
	CompetitorName() string
}

func NewCompetitorName(v string) (CompetitorName, error) {
	v = utils.XSSFilter(v)

	if v == "" || len(v) > 30 {
		return nil, errors.New("invalid competitor name")
	}

	return competitorName(v), nil
}

type competitorName string

func (r competitorName) CompetitorName() string {
	return string(r)
}

// TeamName
type TeamName interface {
	TeamName() string
}

func NewTeamName(v string) (TeamName, error) {
	v = utils.XSSFilter(v)

	if v == "" || utils.StrLen(v) > 15 || len(v) > 40 {
		return nil, errors.New("invalid team name")
	}

	return teamName(v), nil
}

type teamName string

func (r teamName) TeamName() string {
	return string(r)
}

// TeamRole
type TeamRole interface {
	TeamRole() string
	IsLeader() bool
}

func NewTeamRole(v string) (TeamRole, error) {
	if v == "" || v == competitionTeamRoleLeader {
		return teamRole(v), nil
	}

	return nil, errors.New("invalid team role")
}

type teamRole string

func (r teamRole) TeamRole() string {
	return string(r)
}

func (r teamRole) IsLeader() bool {
	return string(r) == competitionTeamRoleLeader
}

func TeamLeaderRole() string {
	return competitionTeamRoleLeader
}

// CompetitionTag
type CompetitionTag interface {
	CompetitionTag() string
}

func NewCompetitionTag(v string) (CompetitionTag, error) {
	b := v == competitionTagChallenge ||
		v == competitionTagElectricity ||
		v == competitionTagLearn

	if b {
		return competitionTag(v), nil
	}

	return nil, errors.New("invalid competition tags")
}

type competitionTag string

func (r competitionTag) CompetitionTag() string {
	return string(r)
}

const (
	languageEN = "en"
	languageCN = "cn"
)

// Language
type Language interface {
	Language() string
}

func NewLanguage(v string) (Language, error) {
	b := v == languageEN ||
		v == languageCN

	if b {
		return dpLanguage(v), nil
	}

	return nil, errors.New("invalid competition status")
}

type dpLanguage string

func (r dpLanguage) Language() string {
	return string(r)
}
