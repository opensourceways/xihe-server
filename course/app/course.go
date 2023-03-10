package app

import (
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"github.com/opensourceways/xihe-server/course/domain/user"
)

type CourseService interface {
	// player
	Apply(*PlayerApplyCmd) (code string, err error)

	// course
	List(*CourseListCmd) ([]CourseSummaryDTO, error)
	Get(*CourseGetCmd) (CourseDTO, error)
	AddReleatedProject(*CourseAddReleatedProjectCmd) (string, error)
}

func NewCourseService(
	userCli user.User,

	courseRepo repository.Course,
	playerRepo repository.Player,
) *courseService {
	return &courseService{
		userCli:    userCli,
		courseRepo: courseRepo,
		playerRepo: playerRepo,
	}
}

type courseService struct {
	userCli user.User

	courseRepo repository.Course
	playerRepo repository.Player
}

// List
func (s *courseService) List(cmd *CourseListCmd) (
	dtos []CourseSummaryDTO, err error,
) {
	return s.listCourses(&repository.CourseListOption{
		Status: cmd.Status,
		Type:   cmd.Type,
	})

}

func (s *courseService) listCourses(opt *repository.CourseListOption) (
	dtos []CourseSummaryDTO, err error,
) {
	v, err := s.courseRepo.FindCourses(opt)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]CourseSummaryDTO, len(v))
	for i := range v {
		n, err := s.playerRepo.PlayerCount(v[i].Id)
		if err != nil {
			return nil, err
		}

		toCourseSummaryDTO(&v[i], n, &dtos[i])
	}

	return
}

func (s *courseService) Get(cmd *CourseGetCmd) (dto CourseDTO, err error) {
	c, err := s.courseRepo.FindCourse(cmd.Cid)
	if err != nil {
		return
	}

	count, err := s.playerRepo.PlayerCount(c.Id)
	if err != nil {
		return
	}

	if cmd.User != nil {
		p, _ := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)
		if c.IsApplyed(&p.Player) {
			dto.toCourseDTO(&c, true, count)

			return
		}

	}

	dto.toCourseNoVideoDTO(&c, false, count)

	return
}
