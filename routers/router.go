package routers

import (
	"demo/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
)

func InitRouter(iad string) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger
	r = gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%v %v %v %v", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	go func() {
		r.Run(iad)
	}()

}

func SetUpRouter(method string, path string, acFunc func(*gin.Context)) {
	method = strings.ToUpper(method)

	log.Infof("SetUpRouter method: %s, path: %s", method, path)
	switch method {
	case "GET":
		r.GET(path, acFunc)
	case "POST":
		r.POST(path, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": method})
		})
	case "DEL":
		r.DELETE(path, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": method})
		})
	default:
		log.Errorf("SetUpRouter error unknown method: %s", method)
	}
}
