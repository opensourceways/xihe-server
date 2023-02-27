package domain

import (
	course "github.com/opensourceways/xihe-server/course/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type UserRegInfo struct {
	Account  types.Account
	Name     Name
	City     City
	Email    Email
	Phone    Phone
	Identity Identity
	Province Province
	Detail   map[string]string
}

func (r *UserRegInfo) NewRegFromStudent(s *course.Student) (err error) {
	r.Account = s.Account

	r.Name, err = NewName(s.Name.StudentName())
	if err != nil {
		return
	}

	r.City, err = NewCity(s.City.City())
	if err != nil {
		return
	}

	r.Email, err = NewEmail(s.Email.Email())
	if err != nil {
		return
	}

	r.Phone, err = NewPhone(s.Phone.Phone())
	if err != nil {
		return
	}

	r.Identity, err = NewIdentity(s.Identity.StudentIdentity())
	if err != nil {
		return
	}

	r.Province, err = NewProvince(s.Province.Province())
	if err != nil {
		return
	}

	r.Detail = s.Detail

	return
}
