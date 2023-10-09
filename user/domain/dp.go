package domain

import (
	"errors"
	"net/url"
	"strings"

	libutil "github.com/opensourceways/community-robot-lib/utils"
	codomain "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	identityStudent   = "student"
	identityTeacher   = "teacher"
	identityDeveloper = "developer"
)

// DomainValue
type DomainValue interface {
	DomainValue() string
}

func IsSameDomainValue(a, b DomainValue) bool {
	if a == nil && b == nil {
		return true
	}

	if a != nil && b != nil {
		return a.DomainValue() == b.DomainValue()
	}

	return false
}

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !utils.IsUserName(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// Password
type Password interface {
	Password() string
}

func NewPassword(s string) (Password, error) {
	if n := len(s); n < 8 || n > 20 {
		return nil, errors.New("invalid password")
	}

	part := make([]bool, 4)

	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			part[0] = true
		} else if c >= 'A' && c <= 'Z' {
			part[1] = true
		} else if c >= '0' && c <= '9' {
			part[2] = true
		} else {
			part[3] = true
		}
	}

	i := 0
	for _, b := range part {
		if b {
			i++
		}
	}
	if i < 3 {
		return nil, errors.New(
			"the password must includes three of lowercase, uppercase, digital and special character",
		)
	}

	return dpPassword(s), nil
}

type dpPassword string

func (r dpPassword) Password() string {
	return string(r)
}

// Bio
type Bio interface {
	Bio() string

	DomainValue
}

func NewBio(v string) (Bio, error) {
	if v == "" {
		return dpBio(v), nil
	}

	if utils.StrLen(v) > codomain.DomainConfig.MaxBioLength {
		return nil, errors.New("invalid bio")
	}

	return dpBio(v), nil
}

type dpBio string

func (r dpBio) Bio() string {
	return string(r)
}

func (r dpBio) DomainValue() string {
	return string(r)
}

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	if v != "" && !libutil.IsValidEmail(v) {
		return nil, errors.New("invalid email")
	}

	return dpEmail(v), nil
}

type dpEmail string

func (r dpEmail) Email() string {
	return string(r)
}

// AvatarId
type AvatarId interface {
	AvatarId() string

	DomainValue
}

func NewAvatarId(v string) (AvatarId, error) {
	if v == "" {
		return dpAvatarId(v), nil
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid avatar")
	}

	if !codomain.DomainConfig.HasAvatarURL(v) {
		v = codomain.DomainConfig.AvatarURL[0]
	}

	return dpAvatarId(v), nil
}

type dpAvatarId string

func (r dpAvatarId) AvatarId() string {
	return string(r)
}

func (r dpAvatarId) DomainValue() string {
	return string(r)
}

// Name
type Name interface {
	Name() string
}

func NewName(v string) (Name, error) {
	v = utils.XSSFilter(v)

	if v == "" || utils.StrLen(v) > 30 {
		return nil, errors.New("invalid name")
	}

	return name(v), nil
}

type name string

func (r name) Name() string {
	return string(r)
}

// City
type City interface {
	City() string
}

func NewCity(v string) (City, error) {
	return city(v), nil
}

type city string

func (r city) City() string {
	return string(r)
}

// Phone
type Phone interface {
	Phone() string
}

func NewPhone(v string) (Phone, error) {
	return phone(v), nil
}

type phone string

func (r phone) Phone() string {
	return string(r)
}

// Identity
type Identity interface {
	Identity() string
}

func NewIdentity(v string) (Identity, error) {
	b := v == identityStudent ||
		v == identityTeacher ||
		v == identityDeveloper ||
		v == ""

	if !b {
		return nil, errors.New("invalid competition identity")
	}

	return identity(v), nil
}

type identity string

func (r identity) Identity() string {
	return string(r)
}

// Province
type Province interface {
	Province() string
}

func NewProvince(v string) (Province, error) {
	return province(v), nil
}

type province string

func (r province) Province() string {
	return string(r)
}
