package domain

import (
	"errors"

	libutil "github.com/opensourceways/community-robot-lib/utils"
)

const (
	identityStudent   = "student"
	identityTeacher   = "teacher"
	identityDeveloper = "developer"
)

// Name
type Name interface {
	Name() string
}

func NewName(v string) (Name, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return name(v), nil
}

type name string

func (r name) Name() string {
	return string(r)
}

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	if v == "" || !libutil.IsValidEmail(v) {
		return nil, errors.New("invalid email")
	}

	return dpEmail(v), nil
}

type dpEmail string

func (r dpEmail) Email() string {
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
