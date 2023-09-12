package domain

import (
	"fmt"
	"io"
	"strconv"

	"github.com/opensourceways/xihe-server/competition/domain/uploader"
	"github.com/opensourceways/xihe-server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubmissionUpdatingInfo
type SubmissionUpdatingInfo struct {
	Index  WorkIndex
	Phase  CompetitionPhase
	Id     string
	Status string
	Score  float32
}

// WorkSubmittedEvent
type WorkSubmittedEvent struct {
	Id            string `json:"id"`
	Phase         string `json:"phase"`
	OBSPath       string `json:"obs_path"`
	PlayerId      string `json:"pid"`
	CompetitionId string `json:"cid"`
}

// Submission
type Submission struct {
	Id       string
	Status   string
	OBSPath  string
	SubmitAt int64
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
			Id:       primitive.NewObjectID().Hex(),
			SubmitAt: now,
			OBSPath:  obspath,
			Status:   "calculating",
		},
		Phase: phase,
	}, nil
}
