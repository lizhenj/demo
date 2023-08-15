package service

import (
	"demo/log"
	r "demo/routers"
)

func InitServices() {
	r.SetUpRouter("Get", "/todo/get")
	r.SetUpRouter("post", "/todo/post")
	r.SetUpRouter("del", "/todo/del")
	log.Infof("InitServices success!!")
}
