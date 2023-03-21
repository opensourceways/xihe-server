package controller

import (
	"github.com/opensourceways/xihe-server/course/app"
	"github.com/opensourceways/xihe-server/course/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type StudentApplyRequest struct {
	Name     string            `json:"name"`
	City     string            `json:"city"`
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	Identity string            `json:"identity"`
	Province string            `json:"province"`
	Detail   map[string]string `json:"detail"`
}

type AddCourseRelatedProjectRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"project_name"`
}

func (req *AddCourseRelatedProjectRequest) ToInfo() (
	owner types.Account, name types.ResourceName, err error,
) {
	if owner, err = types.NewAccount(req.Owner); err != nil {
		return
	}

	name, err = types.NewResourceName(req.Name)

	return
}

func (req *StudentApplyRequest) toCmd(cid string, user types.Account) (cmd app.PlayerApplyCmd, err error) {
	cmd.CourseId = cid

	if cmd.Name, err = domain.NewStudentName(req.Name); err != nil {
		return
	}

	if cmd.City, err = domain.NewCity(req.City); err != nil {
		return
	}

	if cmd.Email, err = types.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Phone, err = domain.NewPhone(req.Phone); err != nil {
		return
	}

	if cmd.Identity, err = domain.NewStudentIdentity(req.Identity); err != nil {
		return
	}

	if cmd.Province, err = domain.NewProvince(req.Province); err != nil {
		return
	}

	cmd.Detail = req.Detail
	cmd.Account = user

	err = cmd.Validate()

	return
}

func toGetCmd(cid string, user types.Account) (cmd app.CourseGetCmd) {
	cmd.User = user
	cmd.Cid = cid

	return
}

type submissionDetail struct {
	AvatarId string `json:"avatar_id"`

	*app.RelateProjectDTO
}
