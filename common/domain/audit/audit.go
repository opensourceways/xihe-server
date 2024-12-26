package audit

type AuditService interface {
	TextAudit(content, contentType string) error
}
