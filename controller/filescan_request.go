package controller

import (
	"strconv"

	"github.com/opensourceways/xihe-server/filescan/app"
	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
	"github.com/opensourceways/xihe-server/user/domain"
)

// ReqToUpdateFileScan is the request of updating a file scan.
type ReqToUpdateFileScan struct {
	// SensitiveItem    string `json:"sensitive_item"    required:"true"`
	ModerationResult string `json:"moderation_result" required:"true"`
}

// ToCmdToUpdateFileScan converts the request to the command.
func (r *ReqToUpdateFileScan) ToCmdToUpdateFileScan(id string) (cmd app.CmdToUpdateFileScan, err error) {
	// if r.Status == "" && r.ModerationResult == "" {
	// 	err = errors.New("need status or moderation result parameter")
	// 	return
	// }

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	cmd.Id = int64(idInt)

	// cmd.ModerationStatus, err = primitive.NewFileModerationStatus(r.Status)
	// if err != nil {
	// 	return
	// }

	// cmd.SensitiveItem, err = primitive.NewSensitiveItemResult(r.SensitiveItem)
	// if err != nil {
	// 	return
	// }

	// cmd.Virus = r.Virus

	cmd.ModerationResult, err = primitive.NewFileModerationResult(r.ModerationResult)
	if err != nil {
		return
	}
	return cmd, nil
}

type Repository struct {
	RepoID   int64  `json:"repo_id"`
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

type RemoveFileScansReq struct {
	Repository

	Removed []string `json:"removed"`
}

func (r RemoveFileScansReq) toCmd() (app.RemoveFileScanCmd, error) {
	return app.RemoveFileScanCmd{
		RepoID:  r.RepoID,
		Removed: r.Removed,
	}, nil
}
