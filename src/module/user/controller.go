package user

import (
	"errors"
	"gorm.io/gorm"
	"logger"
)

import "main/model"

func SelectByUsername(username string) *model.User {
	db := model.Db
	u := model.User{}
	res := db.Where("username = ?", username).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logger.Debugf("Select by username err:" + "未查找到相关数据")
		return nil
	}
	return &u
}

func GetUserByPage(page, pageSize int) (int64, []*model.User) {
	db := model.Db
	var users []*model.User
	var total int64
	db.Model(model.User{}).Count(&total)
	res := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logger.Debugf("Get user by page err:" + "未查找到相关数据")
		return 0, nil
	}
	return total, users
}
