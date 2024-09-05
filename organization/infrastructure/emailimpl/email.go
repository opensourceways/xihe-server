/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package emailimpl provides the implementation of the email interface.
package emailimpl

import (
	"fmt"

	"github.com/opensourceways/xihe-server/organization/domain"
	"github.com/opensourceways/xihe-server/organization/domain/email"
	"github.com/opensourceways/xihe-server/organization/domain/primitive"
)

// Email is an interface for sending organization certificates.
type Email interface {
	Send(receiver []string, subject, content string) error
}

// NewEmailImpl creates a new instance of the email implementation.
func NewEmailImpl(e Email, receiver []string) email.Email {
	return &emailImpl{
		instance: e,
		receiver: receiver,
	}
}

type emailImpl struct {
	instance Email
	receiver []string
}

// Send sends an organization certificate.
func (impl emailImpl) Send(cert domain.OrgCertificate, image primitive.Image) error {
	body := impl.buildEmailBody(cert, image)

	return impl.instance.Send(impl.receiver, "openmind组织认证审核", body)
}

func (impl emailImpl) buildEmailBody(cert domain.OrgCertificate, image primitive.Image) string {
	template := `
<html>
<body>
<h3>申请组织名</h3>
<p>%s</p>
<h3>认证组织类型</h3>
<p>%s</p>
<h3>认证组织名称</h3>
<p>%s</p>
<h3>申请人身份</h3>
<p>%s</p>
<h3>联系电话</h3>
<p>%s</p>
<h3>统一社会信用码/组织机构代码</h3>
<p>%s</p>
<h3>认证组织证件</h3>
<img src="data:image/%s;base64,%s">
</body>
</html>
`
	return fmt.Sprintf(template,
		cert.OrgName.Account(),
		cert.CertificateOrgType.CertificateOrgType(),
		cert.CertificateOrgName.AccountFullname(),
		cert.Identity.Identity(),
		cert.Phone.PhoneNumber(),
		cert.UnifiedSocialCreditCode.USCC(),
		image.ImageType(),
		image.ContentOfBase64(),
	)
}
