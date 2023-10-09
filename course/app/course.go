package app

import (
	"strings"

	"github.com/opensourceways/xihe-server/course/domain/message"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"github.com/opensourceways/xihe-server/course/domain/user"
	projdomain "github.com/opensourceways/xihe-server/domain"
	projectrepo "github.com/opensourceways/xihe-server/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
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
	GetAssignment(*AsgGetCmd) (AsgDTO, error)
	AddPlayRecord(*RecordAddCmd) (string, error)
}

func NewCourseService(
	userCli user.User,
	projectRepo projectrepo.Project,

	courseRepo repository.Course,
	playerRepo repository.Player,
	workRepo repository.Work,
	recordRepo repository.Record,
	producer message.MessageProducer,
) *courseService {
	return &courseService{
		userCli:     userCli,
		projectRepo: projectRepo,

		courseRepo: courseRepo,
		playerRepo: playerRepo,
		workRepo:   workRepo,
		recordRepo: recordRepo,
		producer:   producer,
	}
}

type courseService struct {
	userCli     user.User
	projectRepo projectrepo.Project

	courseRepo repository.Course
	playerRepo repository.Player
	workRepo   repository.Work
	recordRepo repository.Record
	producer   message.MessageProducer
}

// List
func (s *courseService) List(cmd *CourseListCmd) (
	dtos []CourseSummaryDTO, err error,
) {
	if cmd.User != nil {
		return s.getCoursesUserApplied(cmd)
	}

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

func (s *courseService) getCoursesUserApplied(cmd *CourseListCmd) (
	dtos []CourseSummaryDTO, err error,
) {
	cs, err := s.playerRepo.FindCoursesUserApplied(cmd.User)
	if err != nil || len(cs) == 0 {
		return nil, err
	}

	return s.listCourses(&repository.CourseListOption{
		Status:    cmd.Status,
		Type:      cmd.Type,
		CourseIds: cs,
	})
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
		status := "not-finish"
		score := float32(0)
		if cmd.Status == nil {
			w, err := s.workRepo.GetWork(cmd.Cid, cmd.User, a[i].Id, nil)
			if err == nil {
				score = w.Score
				status = w.Status
			}
		} else if cmd.Status.IsFinished() {
			w, err := s.workRepo.GetWork(cmd.Cid, cmd.User, a[i].Id, cmd.Status)
			if err != nil {
				continue
			}
			score = w.Score
			status = w.Status
		} else {
			w, err := s.workRepo.GetWork(cmd.Cid, cmd.User, a[i].Id, nil)
			if err == nil {
				if w.Status != cmd.Status.WorkStatus() {
					continue
				}
			}
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

	resource, err := projdomain.NewResourceName(name[1])
	if err != nil {
		return
	}

	project, err := s.projectRepo.GetByName(cmd.User, resource)
	if err != nil {
		return
	}

	if dtos.RelatedProject, err = s.listRelatedProject(project, 1); err != nil {
		return
	}

	return
}

func (s *courseService) listRelatedProject(project projdomain.Project, count int) (
	dtos []ProjectSummuryDTO, err error,
) {
	dtos = make([]ProjectSummuryDTO, count)
	for i := range dtos {
		toProjectSummuryDTO(&project, &dtos[i])
	}

	return
}

func (s *courseService) GetAssignment(cmd *AsgGetCmd) (
	dto AsgDTO, err error,
) {
	p, err := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)
	if err != nil {
		return
	}

	c, err := s.courseRepo.FindCourse(cmd.Cid)
	if err != nil {
		return
	}

	if !c.IsApplyed(&p.Player) {
		return
	}

	asg, err := s.courseRepo.FindAssignment(cmd.Cid, cmd.AsgId)
	if err != nil {
		return
	}

	toAsgDTO(&asg, &c, &dto)

	return
}

func (s *courseService) AddPlayRecord(cmd *RecordAddCmd) (
	code string, err error,
) {
	// check phase
	course, err := s.courseRepo.FindCourse(cmd.Cid)
	if err != nil {
		return
	}

	// check permission
	player, err := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)

	if !course.IsApplyed(&player.Player) {
		code = errorNoPermission
		return
	}
	r := cmd.toRecord()
	if _, err = s.recordRepo.FindPlayRecord(&r); err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			err = s.recordRepo.AddPlayRecord(&r)
			if err != nil {
				return
			}
		}
	}

	a, err := s.recordRepo.FindPlayRecord(&r)
	if err != nil {
		return
	}

	err = s.recordRepo.UpdatePlayRecord(&r, a.Version)
	if err != nil {
		return
	}

	return
}
