package main

import (
	"bytes"
	"fmt"
	"ginWeb/config"
	_ "ginWeb/docs"
	"ginWeb/log"
	"ginWeb/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
)

var logger *zap.Logger

// @title Swagger Example API
// @version 0.0.1
// @description  云e办后台管理系统
// @BasePath /
func main() {
	app := gin.New()
	pprof.Register(app)
	app.Use(gin.Recovery())
	app.Use(dolog())
	//设置Validator提示中文
	util.SetValidatorZH()
	//初始化日志
	logger = log.Logger
	//初始化路由
	initRouter(app)

	//初始化redis
	util.CreateRedis()
	logger.Debug("========服务器启动成功==========")
	err := app.Run(":" + strconv.Itoa(config.AppConfig.AppInfo.Port))

	if err != nil {
		fmt.Println("......服务器启动失败..............")
	}

}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

//controller 的日志记录
func dolog() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		method := c.Request.Method
		url := c.Request.RequestURI
		if "GET" == method {
			logger.Info("=====客户端请求:", zap.String("方式method:", method), zap.String("路径url:", url))
		} else {
			data, _ := c.GetRawData()
			logger.Info("=====客户端请求:", zap.String("方式method:", method),
				zap.String("路径url:", url),
				zap.String("参数params:", string(data)),
			)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		}
		c.Next()
		logger.Info("=====服务端返回数据:", zap.String("方式method:", method),
			zap.String("路径url:", url), zap.String("数据data:", w.body.String()))
	}
}
