package api

import (
	"go-MQ/common"
	"go-MQ/entity"
	request "go-MQ/entity/request"
	"go-MQ/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserAPI struct{}

var UserAPIEntity = UserAPI{}

func (ua UserAPI) RegisterHandler(ctx *gin.Context) {
	var req request.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Warn("用户注册参数绑定失败:", err)
		common.SendErrorResponse(ctx, "invalid parameters")
		return
	}
	requser := entity.User{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
		Avatar:   req.Avatar,
	}
	token, tokenExpiresAt, user, err := service.UserServiceEntity.Register(requser, req.Code)
	if err != nil {
		logrus.Warn("用户注册失败:", err, user)
		common.SendErrorResponse(ctx, err.Error())
		return
	}
	response := request.RegisterUserResponse{
		ID:             user.ID,
		UserName:       user.UserName,
		Email:          user.Email,
		Avatar:         user.Avatar,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		LastActiveAt:   user.LastActiveAt,
		Token:          token,
		TokenExpiresAt: tokenExpiresAt,
	}
	common.SendSuccessResponse(ctx, "注册成功", response)
}

func (ua UserAPI) LoginHandler(ctx *gin.Context) {
	var (
		err            error
		user           entity.User
		tokenExpiresAt int64
		token          string
	)
	var req request.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Warn("用户登录参数绑定失败:", err)
		common.SendErrorResponse(ctx, "invalid parameters")
		return
	}
	if strings.IndexByte(req.UserName, '@') == -1 {
		userID, err := strconv.ParseUint(req.UserName, 10, 64)

		if err != nil {
			logrus.Println("string转uint64转换失败:", err)
			return
		}

		user, tokenExpiresAt, token, err = service.UserServiceEntity.LoginByID(userID, req.Password)
		if err != nil {
			logrus.Warn("用户登录失败:", err, req.UserName)
			common.SendErrorResponse(ctx, err.Error())
			return
		}
	} else {
		user, tokenExpiresAt, token, err = service.UserServiceEntity.LoginByEmail(req.UserName, req.Password)
		if err != nil {
			logrus.Warn("用户登录失败:", err, req.UserName)
			common.SendErrorResponse(ctx, err.Error())
			return
		}
	}

	response := request.RegisterUserResponse{
		ID:             user.ID,
		UserName:       user.UserName,
		Email:          user.Email,
		Avatar:         user.Avatar,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		LastActiveAt:   user.LastActiveAt,
		Token:          token,
		TokenExpiresAt: tokenExpiresAt,
	}
	common.SendSuccessResponse(ctx, "登录成功", response)
}
