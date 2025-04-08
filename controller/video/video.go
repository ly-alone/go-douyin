package video

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liuyang/ly/model"
	"github.com/tidwall/gjson"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type VideoFrameCaptureParam struct {
	Fps string `form:"fps" binding:"required"` // 帧率参数
}

// VideoFrameCapture 视频取帧
func VideoFrameCapture(c *gin.Context) {
	userInfo := c.Value("user")
	userInfoJson, _ := json.Marshal(userInfo)
	userInfoNew := make(map[string]interface{})
	err := json.Unmarshal(userInfoJson, &userInfoNew)
	if err != nil {
		c.JSON(500, gin.H{"code": 40004, "msg": "系统错误"})
	}
	userData, _ := userInfoNew["userinfo"].(string)
	userId := gjson.Get(userData, "user_id").Int()
	//查询用户信息
	userContent := model.GetUserInfoById(int32(userId), "is_vip")
	if userContent.IsVip == 0 {
		//查询今日使用次数
		vCount, err := model.GetVideoLogData(int32(userId))
		if err != nil {
			c.JSON(400, gin.H{"code": 41000, "msg": "请稍后再试"})
			return
		}
		fmt.Println(vCount)
		if vCount >= 5 {
			c.JSON(400, gin.H{"code": 41001, "msg": "非vip用户一天只可五次"})
			return
		}
	}
	// 接收上传的文件
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(400, gin.H{"code": 40001, "msg": fmt.Sprintf("failed to get uploaded file: %v", err)})
		return
	}
	if file.Size > 1024*1024*15 {
		c.JSON(400, gin.H{"code": 40002, "msg": "非会员最大上传15M"})
		return
	}
	// 获取帧率参数
	var param VideoFrameCaptureParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(400, gin.H{"code": 40002, "msg": "missing fps parameter"})
		return
	}

	// 定义视频保存路径
	videoPath := "./uploads/"
	unixNano := time.Now().UnixNano()
	// 使用时间戳生成唯一子目录
	videoPath = fmt.Sprintf("%v%d/", videoPath, userId)

	// 创建保存目录
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		if err := os.MkdirAll(videoPath, 0755); err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": fmt.Sprintf("failed to create uploads directory: %v", err)})
			return
		}
	}

	// 保存上传的视频
	filePath := filepath.Join(videoPath, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": fmt.Sprintf("failed to save file: %v", err)})
		return
	}

	// 定义输出图片目录
	outputDir := "./image/"
	outputDir = fmt.Sprintf("%v%d/%d", outputDir, userId, unixNano)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": fmt.Sprintf("failed to create frames directory: %v", err)})
			return
		}
	}

	// 调用提取帧的函数
	if err := extractFrames(filePath, outputDir, param.Fps, userId); err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": fmt.Sprintf("extract frames failed: %v", err)})
	} else {
		//读取目录照片展示
		image, _ := getPhotos(outputDir)
		//记录日志表
		go addLog(userId, image, file.Filename, unixNano)
		c.JSON(200, gin.H{"code": 200, "msg": "Video frames extracted successfully!", "data": image})
	}
}

func addLog(userId int64, image []Photo, videoName string, unixNano int64) {
	imgNameArr := make([]string, len(image))
	for k, v := range image {
		imgNameArr[k] = v.Name
	}
	imgName, _ := json.Marshal(imgNameArr)
	unixNanoStr := fmt.Sprintf("%d", unixNano)
	userId32 := int32(userId)
	model.AddVideoLogData(userId32, videoName, string(imgName), unixNanoStr)
}

type Photo struct {
	Name string
	Path string
}

func getPhotos(dir string) ([]Photo, error) {
	var photos []Photo
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isImage(path) {
			photos = append(photos, Photo{
				Name: d.Name(),
				Path: path,
			})
		}
		return nil
	})
	return photos, err
}

// isImage 检查文件是否为常见图片格式
func isImage(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range extensions {
		if filepath.Ext(fileName) == ext {
			return true
		}
	}
	return false
}

func UrlVideoFrameCapture(c *gin.Context) {
	url, _ := c.Get("url")
	if url == nil {
		c.JSON(400, gin.H{"code": 400, "msg": "missing url parameter"})
	}
}

// extractFrames 提取视频帧
func extractFrames(videoPath string, outputDir, fps string, userId int64) error {
	year, month, day := time.Now().Date()
	date := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	pngName := fmt.Sprintf("%v%v_%d", "frame_", date, userId)
	// 定义输出图片的路径格式
	outputPattern := filepath.Join(outputDir, pngName+"_%04d.png")
	// 构建 ffmpeg 命令
	var cmd *exec.Cmd
	if fps == "0" {
		//cmd = exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=100", outputPattern)
		cmd = exec.Command("ffmpeg", "-i", videoPath, outputPattern)
	} else {
		//fps=10 代表的是每秒提取的帧数
		cmd = exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", outputPattern)
	}

	// 设置命令的输出，方便调试
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute ffmpeg: %v", err)
	}
	return nil
}
