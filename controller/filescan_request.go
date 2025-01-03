package controller

import (
	"strconv"

	"github.com/opensourceways/xihe-server/filescan/app"
	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
	"github.com/opensourceways/xihe-server/user/domain"
)

// ReqToUpdateFileScan is the request of updating a file scan.
type ReqToUpdateFileScan struct {
	ModerationResult string `json:"moderation_result" required:"true"`
}

// ToCmdToUpdateFileScan converts the request to the command.
func (r *ReqToUpdateFileScan) ToCmdToUpdateFileScan(id string) (cmd app.CmdToUpdateFileScan, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	cmd.Id = int64(idInt)

	cmd.ModerationResult, err = primitive.NewFileModerationResult(r.ModerationResult)
	if err != nil {
		return
	}
	return cmd, nil
}

type Repository struct {
	RepoId   int64  `json:"repo_id"`
	Owner    string `json:"owner"`
	Branch   string `json:"branch"`
	RepoName string `json:"repo_name"`
}

type CreateFileScanListReq struct {
	Repository

	Added []string `json:"added"`
}

func (r CreateFileScanListReq) toCmd() (app.CreateFileScanListCmd, error) {
	cmd := app.CreateFileScanListCmd{
		RepoId:   r.RepoId,
		Branch:   r.Branch,
		RepoName: r.RepoName,
		Added:    r.Added,
	}

	var err error
	cmd.Owner, err = domain.NewAccount(r.Owner)

	return cmd, err
}

type RemoveFileScanListReq struct {
	Repository

	Removed []string `json:"removed"`
}

func (r RemoveFileScanListReq) toCmd() (app.RemoveFileScanListCmd, error) {
	return app.RemoveFileScanListCmd{
		RepoId:  r.RepoId,
		Removed: r.Removed,
	}, nil
}

type ModifyFileScanListReq struct {
	Repository

	Modified []string `json:"modified"`
}

func (r ModifyFileScanListReq) toCmd() (app.LaunchModerationCmd, error) {
	cmd := app.LaunchModerationCmd{
		RepoId:   r.RepoId,
		Branch:   r.Branch,
		RepoName: r.RepoName,
		Modified: r.Modified,
	}

	var err error
	cmd.Owner, err = domain.NewAccount(r.Owner)

	return cmd, err
}
