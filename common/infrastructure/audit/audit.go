package audit

import (
	"github.com/opensourceways/xihe-audit-sync-sdk/audit"
	"github.com/opensourceways/xihe-audit-sync-sdk/httpclient"
	"golang.org/x/xerrors"

	auditapi "github.com/opensourceways/xihe-audit-sync-sdk/audit/api"
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
		return xerrors.Errorf("fail to moderate")
	} else if resp.Result != "pass" {
		return xerrors.Errorf("moderate unpass")
	}
	return nil
}
