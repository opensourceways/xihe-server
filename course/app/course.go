package app

import (
	"strings"

	"github.com/opensourceways/xihe-server/course/domain/repository"
	"github.com/opensourceways/xihe-server/course/domain/user"
	projdomain "github.com/opensourceways/xihe-server/domain"
	projectrepo "github.com/opensourceways/xihe-server/domain/repository"
)

type CourseService interface {
	// player
	Apply(*PlayerApplyCmd) (code string, err error)

	// course
	List(*CourseListCmd) ([]CourseSummaryDTO, error)
	Get(*CourseGetCmd) (CourseDTO, error)
	AddReleatedProject(*CourseAddReleatedProjectCmd) (string, error)
	ListAssignments(*AsgListCmd) ([]AsgWorkDTO, error)
	GetSubmissions(*GetSubmissionCmd) (RelateProjectDTO, error)
	GetCertification(*CourseGetCmd) (CertInfoDTO, error)
}

func NewCourseService(
	userCli user.User,
	projectRepo projectrepo.Project,

	courseRepo repository.Course,
	playerRepo repository.Player,
	workRepo repository.Work,
) *courseService {
	return &courseService{
		userCli:     userCli,
		projectRepo: projectRepo,

		courseRepo: courseRepo,
		playerRepo: playerRepo,
		workRepo:   workRepo,
	}
}

type courseService struct {
	userCli     user.User
	projectRepo projectrepo.Project

	courseRepo repository.Course
	playerRepo repository.Player
	workRepo   repository.Work
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

func (s *courseService) ListAssignments(cmd *AsgListCmd) (
	dtos []AsgWorkDTO, err error,
) {

	a, err := s.courseRepo.FindAssignments(cmd.Cid)

	if err != nil || len(a) == 0 {
		return
	}

	dtos = make([]AsgWorkDTO, len(a))
	j := 0
	for i := 0; i < len(a); i++ {

		w, err := s.workRepo.GetWork(cmd.Cid, cmd.User, a[i].Id, cmd.Status)
		status := w.Status
		score := w.Score
		if err != nil {
			return nil, err
		}

		if cmd.Status != nil && cmd.Status.WorkStatus() != status {
			continue

		}

		toAsgWorkDTO(&a[i], score, status, &dtos[j])
		j++
	}
	return dtos[:j], err

}

func (s *courseService) GetSubmissions(cmd *GetSubmissionCmd) (
	dtos RelateProjectDTO, err error,
) {

	p, err := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)
	if err != nil {
		return
	}

	repo := p.Player.RelatedProject
	if repo == "" {
		return
	}
	name := strings.Split(repo, "/")

	resorce, err := projdomain.NewResourceName(name[1])
	if err != nil {
		return
	}

	project, err := s.projectRepo.GetByName(cmd.User, resorce)
	if err != nil {
		return
	}

	toRelateProjectDTO(&project, &dtos)
	return
}
