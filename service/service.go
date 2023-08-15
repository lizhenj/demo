package service

import (
	db "demo/database"
	r "demo/routers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	Get  = "get"
	Post = "post"
	Del  = "del"
)

func InitServices() {
	r.SetUpRouter("get", "/todo/get", todoGet)
	r.SetUpRouter("post", "/todo/post", todoPost)
	r.SetUpRouter("del", "/todo/del", todoDel)
}

func todoGet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	bases := c.Query("base")
	if db.InsertActor(id, bases) {
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"result": "error"})
	}
}

func todoPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": Post})
}

func todoDel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"method": Del})
}
