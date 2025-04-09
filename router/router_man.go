package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuyang/ly/controller/home"
	"github.com/liuyang/ly/controller/video"
	"github.com/liuyang/ly/public"
	"net/http"
)

// Router 结构体保存 Gin 引擎实例
type Router struct {
	Engine *gin.Engine
}

func AllRouter(engine *gin.Engine) {
	r := Router{Engine: engine}
	// 注册路由
	r.GetHttp()  // 注册 GET 路由
	r.PostHttp() // 注册 POST 路由
	r.Html()
	// 启动服务器
	r.Https()
}

// Html 加载所有html
func (r *Router) Html() {
	// 加载 HTML 模板文件
	r.Engine.LoadHTMLGlob("views/*.html")
	// 配置路由
	r.Engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.Engine.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.Engine.GET("/home", func(c *gin.Context) {
		// 渲染 HTML 页面
		c.HTML(http.StatusOK, "upload.html", nil)
	})
	r.Engine.GET("/video", func(c *gin.Context) {
		// 渲染 HTML 页面
		c.HTML(http.StatusOK, "video.html", nil)
	})
	r.Engine.GET("/user", func(c *gin.Context) {
		// 渲染 HTML 页面
		c.HTML(http.StatusOK, "user.html", nil)
	})
	r.Engine.GET("/comment", func(c *gin.Context) {
		// 渲染 HTML 页面
		c.HTML(http.StatusOK, "comment.html", nil)
	})
	r.Engine.Static("/static", "./static")
	r.Engine.Static("/image", "./image")

}

// GetHttp 定义 HTTP GET 路由
func (r *Router) GetHttp() {
	v1 := r.Engine.Group("tool")
	{
		v1.GET("comment/list", public.Page(), home.ShowComment)
		v1.GET("proxy", video.VideoProxy) // 这是你需要添加的路由
		//v1.GET("test", trading.Test)
	}
}

// PostHttp 定义 HTTP POST 路由
func (r *Router) PostHttp() {
	r.Engine.POST("tool/login", home.Login)
	r.Engine.POST("tool/register", home.Register)
	r.Engine.GET("tool/show/comment", home.ShowComment)
	v2 := r.Engine.Group("tool", public.ValidateJwt(), public.Count.Url())
	//v2 := r.Engine.Group("tool")
	{
		v2.POST("upload", video.VideoFrameCapture)
		v2.POST("add", video.RemoveWatermark)
		v2.POST("comment", home.Comment)
		//v2.GET("show/comment", home.ShowComment)
		v2.GET("user/info", home.UserInfo)
		v2.GET("user/img", home.Img)
	}
}

// Https 启动 HTTP 服务，监听 9999 端口
func (r *Router) Https() {
	r.Engine.Run(":9999")
}

// respondJSON 是一个辅助函数，用于发送统一格式的 JSON 响应
func (r *Router) respondJSON(c *gin.Context, status int, msg string, data interface{}) {
	c.JSON(status, gin.H{
		"msg":  msg,
		"data": data,
	})
}

// handleError 是一个统一的错误处理函数
func (r *Router) handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"msg":   "error",     // 错误提示信息
		"error": err.Error(), // 错误的具体描述
	})
}

// sendJSONResponse 是一个统一的函数，用来返回 JSON 格式的响应
func (r *Router) sendJSONResponse(c *gin.Context, status int, msg string, data interface{}) {
	c.JSON(status, gin.H{
		"msg":  msg,
		"data": data,
	})
}
