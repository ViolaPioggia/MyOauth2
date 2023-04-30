package dao

import (
	"MyOauth2/config/global"
	"MyOauth2/config/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

var u model.User

func Authentication(ctx context.Context, username int64, password string) (string, error) {
	err, flag := FindUser(ctx, username)
	if err != nil {
		return "", err
	} else if flag == false {
		return "", errors.New("用户不存在")
	}
	err, flag = CheckPassword(ctx, username, password)
	if err != nil && flag {
		return "", err
	} else if !flag {
		return "", errors.New("密码错误")
	} else {
		return string(username), nil
	}
}

func FindUser(ctx context.Context, username int64) (error, bool) {
	db := global.MysqlDB
	u.Username = username
	result := db.WithContext(ctx).Where("username=?", username).First(&model.User{})
	flag := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if flag == false && result != nil {
		return nil, true
	} else if flag == true {
		return nil, false
	} else {
		return result.Error, false
	}
}

func CheckPassword(ctx context.Context, username int64, password string) (error, bool) {
	db := global.MysqlDB
	u.Username = username
	u.Password = password
	//result := db.Where(&model.User{
	//	Username: username,
	//	Password: password,
	//}).First(&model.User{})
	result := db.WithContext(ctx).Where(&u).First(&model.User{})
	flag := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if result.Error == nil {
		return nil, true
	} else if flag {
		return nil, false
	} else {
		return result.Error, false
	}
}
