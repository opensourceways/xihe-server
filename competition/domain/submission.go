package domain

import (
	"fmt"
	"io"
	"strconv"

	"github.com/opensourceways/xihe-server/competition/domain/uploader"
	"github.com/opensourceways/xihe-server/utils"
)

// SubmissionUpdatingInfo
type SubmissionUpdatingInfo struct {
	Index  WorkIndex
	Phase  CompetitionPhase
	Id     string
	Status string
	Score  float32
}

// SubmissionMessage
type SubmissionMessage struct {
	PlayerId      string `json:"pid"`
	CompetitionId string `json:"cid"`
	Phase         string `json:"phase"`
	Id            string `json:"id"`
	OBSPath       string `json:"obs_path"`
}

// Submission
type Submission struct {
	Id       string
	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}

func (info *Submission) isSuccess() bool {
	return info.Status == competitionSubmissionStatusSuccess
}

// PhaseSubmission
type PhaseSubmission struct {
	Phase CompetitionPhase

	Submission
}

// SubmissionService
type SubmissionService struct {
	uploader uploader.SubmissionFileUploader
}

func NewSubmissionService(v uploader.SubmissionFileUploader) SubmissionService {
	return SubmissionService{v}
}

func (s *SubmissionService) Submit(
	w *Work, phase CompetitionPhase, fileName string, data io.Reader,
) (PhaseSubmission, error) {
	now := utils.Now()

	obspath := fmt.Sprintf(
		"%s/%s_%s",
		w.submissionOBSPathPrefix(phase),
		strconv.FormatInt(now, 10), fileName,
	)
	if err := s.uploader.Upload(data, obspath); err != nil {
		return PhaseSubmission{}, err
	}

	return PhaseSubmission{
		Submission: Submission{
			SubmitAt: now,
			OBSPath:  obspath,
			Status:   "calculating",
		},
		Phase: phase,
	}, nil
}
