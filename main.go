package main

import (
	"github.com/gin-gonic/gin"
	conf "github.com/liuyang/ly/conf"
	rh "github.com/liuyang/ly/router"
)

func main() {
	conf.MysqlDb()
	conf.RedisDb()
	engine := gin.Default()
	rh.AllRouter(engine)
}
