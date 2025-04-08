package main

import (
	"fmt"
	"github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model/dbTable"
	"log"
	"os"
	"sync"
	"time"
)

var uploadsPath = "./../uploads"
var imagePath = "./../image"

func main() {
	// 获取当前时间
	now := time.Now()

	// 获取当前时间的前三天日期
	threeDaysAgo := now.AddDate(0, 0, -1)

	// 格式化为 YYYY-MM-DD 格式
	formattedDate := threeDaysAgo.Format("2006-01-02")

	// 将 formattedDate 转换为 time.Time 类型
	parsedDate, err := time.Parse("2006-01-02", formattedDate)
	if err != nil {
		log.Fatal("时间解析错误:", err)
	}

	// 获取数据库连接
	conf.MysqlDb()
	db := conf.GetDB()
	if db == nil {
		fmt.Println("数据库连接失败")
		return
	}
	var videoLogs []dbTable.VideoLog
	if err := db.Select("uid, add_img_directory_name,add_video_name").Where("date < ?", parsedDate).Find(&videoLogs); err.Error != nil {
		fmt.Println(err)
		return
	}
	// 输出查询到的记录
	wg := sync.WaitGroup{}
	wg.Add(len(videoLogs))
	for _, vLog := range videoLogs {
		fmt.Println(vLog.AddImgDirectoryName)
		go delImgData(vLog.AddImgDirectoryName, vLog.UID, &wg)
		go delUploadsData(vLog.UID, vLog.AddVideoName, &wg)
	}
	wg.Wait()
	fmt.Println("success")
}

// 删除目录
func delUploadsData(userPath int32, filesName string, wg *sync.WaitGroup) {
	// 目录路径
	dir := fmt.Sprintf("%v/%d", uploadsPath, userPath)
	fmt.Println(dir)
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("目录不存在")
		return
	}
	// 目录存在，检查是否有 111.mp4 文件
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("无法读取目录:", err)
		return
	}

	// 遍历目录中的文件，检查是否有 111.mp4
	for _, file := range files {
		if !file.IsDir() && file.Name() == filesName {
			// 找到文件并删除
			filePath := dir + "/" + file.Name()
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("删除文件失败:", err)
				return
			}
			fmt.Println("文件已删除：" + filesName)
			return
		}
	}
	fmt.Println("未找到文件: " + filesName)
}

func delImgData(userPath string, uid int32, wg *sync.WaitGroup) {
	defer wg.Done()
	// 目录路径
	dir := fmt.Sprintf("%v/%d", imagePath, uid)
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("目录不存在")
		return
	}
	// 目录存在
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("无法读取目录:", err)
		return
	}
	// 遍历目录中的文件，直接删除目录
	for _, file := range files {
		if file.Name() == userPath {
			err := os.RemoveAll(dir + "/" + file.Name())
			if err != nil {
				fmt.Println("删除目录失败:", err)
			} else {
				fmt.Println("目录及其内容已删除", dir+"/"+file.Name())
			}
		}
	}
}
