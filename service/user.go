package service

import (
	"errors"
	"fmt"
	"go-MQ/common"
	"go-MQ/dao"
	"go-MQ/entity"
	"net"
	"regexp"
	"strings"
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
	if err := IsInvailEmail(user.Email); err != nil {
		logrus.Warn("用户注册邮箱不合法:", user.Email)
		return "", 0, entity.User{}, err
	}
	if err := IsInvailUserName(user.UserName); err != nil {
		logrus.Warn("用户注册用户名已被使用:", user.UserName)
		return "", 0, entity.User{}, err
	}
	if err := IsInvailPassword(user.Password); err != nil {
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

func (us UserService) LoginByEmail(email, password string) (entity.User, int64, string, error) {
	user, err := dao.UserDAOEntity.GetUserByEmail(email)
	if err != nil {
		logrus.Warn("用户登录通过邮箱获取用户失败:", err, email)
		return entity.User{}, 0, "", errors.New("invalid email")
	}
	//密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Warn("密码错误:", err, email)
		return entity.User{}, 0, "", errors.New("invalid password")
	}

	token, tokenExpiresAt, err := common.ReleaseToken(user.ID)
	if err != nil {
		logrus.Warn("生成token失败:", err, user.ID)
		return entity.User{}, 0, "", err
	}

	return *user, tokenExpiresAt, token, nil
}

func (us UserService) LoginByID(ID uint64, password string) (entity.User, int64, string, error) {
	user, err := dao.UserDAOEntity.GetUserByID(ID)
	if err != nil {
		logrus.Warn("用户登录通过ID获取用户失败:", err, ID)
		return entity.User{}, 0, "", errors.New("invalid email")
	}
	//密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Warn("密码错误:", err, ID)
		return entity.User{}, 0, "", errors.New("invalid password")
	}

	token, tokenExpiresAt, err := common.ReleaseToken(user.ID)
	if err != nil {
		logrus.Warn("生成token失败:", err, user.ID)
		return entity.User{}, 0, "", err
	}

	return *user, tokenExpiresAt, token, nil
}

func (us UserService) UpdateUserInfo(ID uint64) (entity.User, error) {

	return entity.User{}, errors.New("not implemented")
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsInvailEmail(email string) error {
	var count int64
	if err := common.GetDB().Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		logrus.Error("Failed to check email validity:", err, email)
		return err
	}
	if count > 0 {
		return errors.New("邮箱已被注册")
	}

	if len(email) < 3 || len(email) > 254 {
		return errors.New("邮箱长度不合法")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式错误")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("无法提取域名")
	}
	domain := parts[1]

	// 查询 MX 记录
	// net.LookupMX 会去询问 DNS 服务器："这个域名有邮件服务器吗？"
	mxRecords, err := net.LookupMX(domain)

	// 如果报错(如网络不通) 或者 找不到任何 MX 记录
	if err != nil || len(mxRecords) == 0 {
		// 这里的 err 可能会包含 "no such host" 等信息
		return errors.New("该邮箱域名不存在或无法接收邮件 (DNS MX lookup failed)")
	}

	return nil
}

func IsInvailUserName(username string) error {
	if len(username) > 10 {
		return errors.New("用户名过长")
	}
	return nil
}

func IsInvailPassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("密码长度应在6到20个字符之间")
	}
	return nil
}
