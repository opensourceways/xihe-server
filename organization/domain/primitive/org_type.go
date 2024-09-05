/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "errors"

const (
	typeEnterprise        = "企业"
	typeSchool            = "学校"
	typeResearch          = "研究机构"
	typePublicInstitution = "事业单位"
	typeFoundation        = "基金会"
	typeOthers            = "其他"
)

// CertificateOrgType is an interface representing a org type
type CertificateOrgType interface {
	CertificateOrgType() string
}

// ValidateOrgType validates the org type
func ValidateOrgType(v string) bool {
	return v == typeEnterprise || v == typeSchool || v == typeResearch ||
		v == typePublicInstitution || v == typeFoundation || v == typeOthers
}

// NewCertificateOrgType creates a new OrgType instance from a string value.
func NewCertificateOrgType(v string) (CertificateOrgType, error) {
	if !ValidateOrgType(v) {
		return nil, errors.New("invalid org type")
	}

	return certificateOrgType(v), nil
}

// CreateCertificateOrgType creates a new OrgType instance from a string value.
func CreateCertificateOrgType(v string) CertificateOrgType {
	return certificateOrgType(v)
}

// certificateOrgType represents the type of organization
type certificateOrgType string

// CertificateOrgType represents the type of organization
func (o certificateOrgType) CertificateOrgType() string {
	return string(o)
}
