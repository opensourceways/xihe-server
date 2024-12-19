package audit

import (
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-audit-sync-sdk/audit"
	auditapi "github.com/opensourceways/xihe-audit-sync-sdk/audit/api"
	"github.com/opensourceways/xihe-audit-sync-sdk/httpclient"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
)

const auditPass = "pass"

var instance *auditImpl

// NewModerationService returns the singleton instance of the moderation service.
func NewAuditService() *auditImpl {
	return instance
}

// Init initializes the Moderation agent with the specified configuration.
func Init(cfg *Config) {
	httpclient.Init(cfg)

	instance = &auditImpl{
		cfg: cfg,
	}
}

type Config = httpclient.Config

type auditImpl struct {
	cfg *Config
}

func (a *auditImpl) TextAudit(content, contentType string) error {
	var resp audit.ModerationDTO
	resp, _, err := auditapi.Text(content, contentType)
	if err != nil {
		e := xerrors.New("call audit failed")
		return allerror.New(
			allerror.ErrorCodeCallAuditFailed, resp.Result, e)
	} else if resp.Result != auditPass {
		e := xerrors.New("audit block")
		return allerror.New(
			allerror.ErrorCodeAuditBlock, resp.Result, e)
	}
	return nil
}
