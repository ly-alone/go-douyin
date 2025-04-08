package model

import (
	"errors"
	"fmt"
	"github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model/dbTable"
	"time"
)

func AddVideoLogData(uid int32, videoName, imgName, unixNanoStr string) (dbTable.VideoLog, error) {
	now := time.Now()
	date := now.Format("2006-01-02")
	dateTime, _ := time.Parse("2006-01-02", date)
	data := dbTable.VideoLog{UID: uid, Date: dateTime, AddVideoName: videoName, AddImgName: imgName, AddImgDirectoryName: unixNanoStr}
	db := conf.GetDB()
	res := db.Select("uid", "add_video_name", "add_img_name", "add_img_directory_name", "date").Create(&data)
	if res.Error != nil {
		return data, errors.New("add errror：")
	}
	return data, nil
}

func GetVideoLogData(uid int32) (int64, error) {
	now := time.Now()
	y, m, d := now.Date()
	date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
	db := conf.GetDB()
	var count int64
	model := dbTable.VideoLog{}
	res := db.Model(&model).Where("date = ?", date).Where("uid =?", uid).Count(&count)
	if res.Error != nil {
		fmt.Println(res.Error)
		return count, errors.New("add errror：")
	}
	return count, nil
}
