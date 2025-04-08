package model

import (
	"errors"
	"fmt"
	"github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model/dbTable"
	"math/rand"
	"time"
)

func GetUserInfoByName(name string) dbTable.User {
	var user dbTable.User
	db := conf.GetDB()
	db.Where("name = ?", name).First(&user)
	return user
}

func GetUserInfoById(id int32, fields ...string) *dbTable.User {
	var user *dbTable.User
	db := conf.GetDB()
	db.Select(fields).Where("id = ?", id).First(&user)
	return user
}

func AddUserData(name, paw, email string) (dbTable.User, error) {
	times := time.Now().Unix()
	ramdInt := rand.Intn(17) + 1
	img := fmt.Sprintf("/static/img/%d.png", ramdInt)
	data := dbTable.User{Name: name, Pwd: paw, Email: email, CTime: int32(times), Img: img}
	db := conf.GetDB()
	res := db.Select("name", "pwd", "email", "c_time", "img").Create(&data)
	if res.Error != nil {
		return data, errors.New("add errror：")
	}
	return data, nil
}
