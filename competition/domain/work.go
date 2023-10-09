package domain

import (
	"fmt"

	"github.com/opensourceways/xihe-server/utils"
)

type WorkIndex struct {
	PlayerId      string
	CompetitionId string
}

func NewWorkIndex(cid, pid string) WorkIndex {
	return WorkIndex{
		PlayerId:      pid,
		CompetitionId: cid,
	}
}

// Work
type Work struct {
	WorkIndex

	PlayerName string

	Repo        string
	Final       []Submission
	Preliminary []Submission
}

func NewWork(cid string, p *Player) Work {
	return Work{
		WorkIndex:  NewWorkIndex(cid, p.Id),
		PlayerName: p.Name(),
	}
}

func (w *Work) Submissions(phase CompetitionPhase) []Submission {
	if phase.IsPreliminary() {
		return w.Preliminary
	}

	if phase.IsFinal() {
		return w.Final
	}

	return nil
}

func (w *Work) BestOne(phase CompetitionPhase, order CompetitionScoreOrder) (
	r *Submission,
) {
	submissions := w.Submissions(phase)
	for i := range submissions {
		item := &submissions[i]

		if !item.isSuccess() {
			continue
		}

		if r == nil || order.IsBetterThanB(item.Score, r.Score) {
			r = item
		}
	}

	return
}

func (w *Work) submissionOBSPathPrefix(phase CompetitionPhase) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		w.CompetitionId,
		phase.CompetitionPhase(),
		w.PlayerId,
	)
}

func (w *Work) HasSubmittedToday(phase CompetitionPhase) bool {
	today := utils.Date()
	submissions := w.Submissions(phase)
	for i := range submissions {
		if utils.ToDate(submissions[i].SubmitAt) == today {
			return true
		}
	}

	return false
}

func (w *Work) NewSubmissionMessage(s *PhaseSubmission) WorkSubmittedEvent {
	return WorkSubmittedEvent{
		Id:            s.Submission.Id,
		Phase:         s.Phase.CompetitionPhase(),
		OBSPath:       s.Submission.OBSPath,
		PlayerId:      w.WorkIndex.PlayerId,
		CompetitionId: w.WorkIndex.CompetitionId,
	}
}

func (w *Work) UpdateSubmission(info *SubmissionUpdatingInfo) *Submission {
	submissions := w.Submissions(info.Phase)
	for i := range submissions {
		if item := &submissions[i]; item.Id == info.Id {
			item.Status = info.Status
			item.Score = info.Score

			return &submissions[i]
		}
	}

	return nil
}
