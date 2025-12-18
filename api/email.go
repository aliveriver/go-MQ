package api

import (
	"go-MQ/common"
	entity "go-MQ/entity/request"
	"go-MQ/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type EmailAPi struct{}

var EmailAPIEntity = EmailAPi{}

func (ea EmailAPi) SendEmailCode(ctx *gin.Context) {
	var req entity.RequestEmailCode
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Warn("发送邮箱验证码参数绑定失败:", err)
		common.SendErrorResponse(ctx, "invalid parameters")
		return
	}
	c := ctx.Request.Context()
	if err := service.EmailServiceEntity.SendEmailCode(c, req.Email); err != nil {
		logrus.Warn("发送邮箱验证码失败:", err)
		common.SendErrorResponse(ctx, err.Error())
		return
	}
	common.SendSuccessResponse(ctx, "邮箱验证码发送成功", nil)
}
