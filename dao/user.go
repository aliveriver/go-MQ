package dao

import (
	"go-MQ/common"
	"go-MQ/entity"

	"github.com/sirupsen/logrus"
)

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
