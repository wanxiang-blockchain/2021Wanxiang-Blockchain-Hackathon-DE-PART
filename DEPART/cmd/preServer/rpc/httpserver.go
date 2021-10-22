package rpc

import (
	"net/http"
	"preServer/conf"
	"preServer/utils"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	conf conf.HttpConfig

	router *gin.Engine
}

func NewHttpServer(config conf.HttpConfig) *HttpServer {
	return &HttpServer{conf: config}
}

func (s *HttpServer) Start() {
	s.router = makeRouter()
	s.router.Run(s.conf.Ip + ":" + s.conf.Port)
}

func makeRouter() *gin.Engine {
	router := gin.New()
	gin.Default()

	//router.Use(gin.Recovery())
	router.NoMethod(HandleNotFound)
	router.NoRoute(HandleNotFound)
	router.Use(ErrHandler())
	urls(router)

	return router
}

func urls(router *gin.Engine) {
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/uploadID", parseUploadID)
	router.POST("/reEncryption", parseReEncryption)
}

func parseUploadID(c *gin.Context) {
	var rec *utils.Record
	var err error
	contentType := c.Request.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&rec)
	}
	if err != nil || rec == nil {
		c.String(http.StatusBadRequest, "request bind err")
		return
	}
	err = uploadId(rec)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")

}

func parseReEncryption(c *gin.Context) {
	var rk *utils.ReKey
	var err error
	contentType := c.Request.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&rk)
	}
	if err != nil || rk == nil {
		c.String(http.StatusBadRequest, "request bind err")
		return
	}
	result, err := reEncryption(rk)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
