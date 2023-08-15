package main

import (
	"demo/database"
	"demo/log"
	"demo/routers"
	"demo/service"
	_ "demo/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/go-sql-driver/mysql"
)

var (
	signalC = make(chan os.Signal, 1)
	config  = make(map[string]string) //配置参数
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Infof("%v\n%s", err, string(debug.Stack()))
		}
		log.Info("demo exit!!")
		log.Close()
	}()

	//加载配置文件
	if buff, err := ioutil.ReadFile("config.json"); err != nil {
		panic(err)
	} else if err = json.Unmarshal(buff, &config); err != nil {
		panic(err)
	}

	//设置日志
	mysql.SetLogger(log.Logger)
	log.SetFile(fmt.Sprintf("%s/demo", config["logDir"]))
	log.Infof("%v", config)

	//启动数据库
	database.InitDb(config["gamedb"])

	//启动gin
	routers.InitRouter(config["ginIpaddress"])
	//初始化路由
	service.InitServices()

	log.Infof("----------demo runing---------")
	go func() {
		for {
			var s string
			fmt.Scanln(&s)
			switch s {
			case "exit":
				signalC <- os.Interrupt
			}
		}
	}()

	signal.Notify(signalC, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-signalC:
	}

	//程序进行退出的工作

}
