package service

import (
	"context"
	"go-MQ/tool"

	"github.com/sirupsen/logrus"
)

type EmailService struct{}

var EmailServiceEntity = EmailService{}

func (es EmailService) SendEmailCode(ctx context.Context, email string) error {
	if err := IsInvailEmail(email); err != nil {
		logrus.Warn("发送邮箱验证码失败，邮箱不合法:", email)
		return err
	}
	_, code, err := tool.SendVerificationEmail(email)
	if err != nil {
		logrus.Warn("发送邮箱验证码失败:", email, err)
		return err
	}
	if err := tool.SaveEmailCode(ctx, email, code); err != nil {
		logrus.Warn("保存邮箱验证码失败:", email, err)
		return err
	}
	return nil
}
