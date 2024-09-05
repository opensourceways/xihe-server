/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

const (
	statusPassed     = "passed"
	statusProcessing = "processing"
	statusFailed     = "failed"
)

// CertificateStatus represents the status of certificate
type CertificateStatus interface {
	CertificateStatus() string
}

// NewProcessingStatus create new processing status
func NewProcessingStatus() CertificateStatus {
	return certificateStatus(statusProcessing)
}

// NewPassedStatus create new passed status
func NewPassedStatus() CertificateStatus {
	return certificateStatus(statusPassed)
}

// CreateCertificateStatus create new certificate status
func CreateCertificateStatus(v string) CertificateStatus {
	return certificateStatus(v)
}

type certificateStatus string

// CertificateStatus returns the status of certificate
func (c certificateStatus) CertificateStatus() string {
	return string(c)
}
