package dao

import (
	"errors"
	"fmt"
	"go-MQ/common"
	"go-MQ/entity"
	"net"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type UserDAO struct{}

var UserDAOEntity = UserDAO{}

func (dao UserDAO) AddUser(user *entity.User) error {
	if err := common.GetDB().Create(user).Error; err != nil {
		logrus.Error("Failed to add user:", err, user)
		return err
	}
	return nil
}

func (dao UserDAO) GetUserByName(username string) (*entity.User, error) {
	var user entity.User
	if err := common.GetDB().Where("user_name = ?", username).First(&user).Error; err != nil {
		logrus.Error("Failed to get user by name:", err, username)
		return nil, err
	}
	return &user, nil
}

func (dao UserDAO) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := common.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		logrus.Error("Failed to get user by email:", err, email)
		return nil, err
	}
	return &user, nil
}

func (dao UserDAO) GetUserByID(ID uint64) (*entity.User, error) {
	var user entity.User
	if err := common.GetDB().Where("id = ?", ID).First(&user).Error; err != nil {
		logrus.Error("Failed to get user by ID:", err, ID)
		return nil, err
	}
	return &user, nil
}

func (dao UserDAO) UpdateUser(userID uint64, datas map[string]interface{}) error {
	if err := common.GetDB().Model(&entity.User{}).Where("id = ?", userID).Updates(datas).Error; err != nil {
		logrus.Error("Failed to update user:", err, userID, datas)
		return err
	}
	return nil
}

func (dao UserDAO) DeleteUser(userID uint64) error {
	if err := common.GetDB().Where("id = ?", userID).Delete(&entity.User{}).Error; err != nil {
		logrus.Error("Failed to delete user:", err, userID)
		return err
	}
	return nil
}

func (dao UserDAO) IsInvailEmail(email string) error {
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

func (dao UserDAO) IsInvailUserName(username string) error {
	if len(username) > 10 {
		return errors.New("用户名过长")
	}
	return nil
}

func (dao UserDAO) IsInvailPassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("密码长度应在6到20个字符之间")
	}
	return nil
}
