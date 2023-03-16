package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	Account  types.Account
	Name     StudentName
	City     City
	Email    types.Email
	Phone    Phone
	Identity StudentIdentity
	Province Province
	Detail   map[string]string
}

type Player struct {
	Student

	Id             string
	CourseId       string
	CreatedAt      CourseTime
	RelatedProject string
}

func (p *Player) CreateToday() (err error) {
	if p.CreatedAt, err = NewCourseTime(utils.Now()); err != nil {
		return
	}

	return
}

func (p *Player) NewId() {
	if p.Id == "" {
		p.Id = primitive.NewObjectID().Hex()
	}
}
