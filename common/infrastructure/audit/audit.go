package audit

import (
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-audit-sync-sdk/audit"
	auditapi "github.com/opensourceways/xihe-audit-sync-sdk/audit/api"
	"github.com/opensourceways/xihe-audit-sync-sdk/httpclient"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
)

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
		e := xerrors.Errorf("call audit failed")
		return allerror.New(
			allerror.ErrorCodeCallAuditFailed,
			"", e)
	} else if resp.Result != "pass" {
		e := xerrors.Errorf("audit block")
		return allerror.New(
			allerror.ErrorCodeAuditBlock,
			"", e)
	}
	return nil
}
