package audit

import (
	"github.com/opensourceways/xihe-audit-sync-sdk/audit"
	"github.com/opensourceways/xihe-audit-sync-sdk/httpclient"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"

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
		e := xerrors.Errorf("fail to moderate")
		return allerror.New(
			allerror.ErrorCodeFailToModerate,
			"", e)
	} else if resp.Result != "pass" {
		e := xerrors.Errorf("fail to moderate")
		return allerror.New(
			allerror.ErrorCodeModerateUnpass,
			"", e)
	}
	return nil
}
