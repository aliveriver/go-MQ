package service

import (
	"errors"
	"go-MQ/common"
	"go-MQ/dao"
	"go-MQ/entity"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

var UserServiceEntity = UserService{}

func (us UserService) Register(user entity.User, code string) (string, int64, entity.User, error) {
	if user.Password == "" || user.UserName == "" {
		logrus.Warn("用户注册缺少信息")
		return "", 0, entity.User{}, errors.New("username or password cannot be empty")
	}
	if err := dao.UserDAOEntity.IsInvailEmail(user.Email); err != nil {
		logrus.Warn("用户注册邮箱不合法:", user.Email)
		return "", 0, entity.User{}, err
	}
	if err := dao.UserDAOEntity.IsInvailUserName(user.UserName); err != nil {
		logrus.Warn("用户注册用户名已被使用:", user.UserName)
		return "", 0, entity.User{}, err
	}
	if err := dao.UserDAOEntity.IsInvailPassword(user.Password); err != nil {
		logrus.Warn("用户注册密码不合法:", user.Password)
		return "", 0, entity.User{}, err
	}
	//TODO: 校验验证码

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("密码加密失败:", err)
		return "", 0, entity.User{}, err
	}
	user.Password = string(hasedPassword)

	user.LastActiveAt = time.Now().UnixMilli()

	if err := dao.UserDAOEntity.AddUser(&user); err != nil {
		logrus.Error("用户注册失败:", err, user)
		return "", 0, entity.User{}, err
	}

	token, tokenExpiresAt, err := common.ReleaseToken(user.ID)
	if err != nil {
		return "", 0, entity.User{}, err
	}
	return token, tokenExpiresAt, user, nil
}

func (us UserService) LoginByEmail(email, password string) (error, entity.User, int64, string) {
	user, err := dao.UserDAOEntity.GetUserByEmail(email)
	if err != nil {
		logrus.Warn("用户登录通过邮箱获取用户失败:", err, email)
		return errors.New("invalid email"), entity.User{}, 0, ""
	}
	//密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Warn("密码错误:", err, email)
		return errors.New("invalid password"), entity.User{}, 0, ""
	}

	token, tokenExpiresAt, err := common.ReleaseToken(user.ID)
	if err != nil {
		logrus.Warn("生成token失败:", err, user.ID)
		return err, entity.User{}, 0, ""
	}

	return nil, *user, tokenExpiresAt, token
}

func (us UserService) LoginByID(ID uint64, password string) (error, entity.User, int64, string) {
	user, err := dao.UserDAOEntity.GetUserByID(ID)
	if err != nil {
		logrus.Warn("用户登录通过ID获取用户失败:", err, ID)
		return errors.New("invalid email"), entity.User{}, 0, ""
	}
	//密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Warn("密码错误:", err, ID)
		return errors.New("invalid password"), entity.User{}, 0, ""
	}

	token, tokenExpiresAt, err := common.ReleaseToken(user.ID)
	if err != nil {
		logrus.Warn("生成token失败:", err, user.ID)
		return err, entity.User{}, 0, ""
	}

	return nil, *user, tokenExpiresAt, token
}
