package main

import (
	"flag"
	"fmt"
	"github.com/helloh2o/lucky"
	"github.com/helloh2o/lucky/cache"
	"github.com/helloh2o/lucky/log"
	"github.com/helloh2o/lucky/xdb"
	"gorm.io/gorm"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"xvid/config"
	"xvid/constans"
	"xvid/entity"
	"xvid/router"
)

var confPath = flag.String("conf", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	cfg := config.Initialize(*confPath)
	// set log
	logger, err := log.New(cfg.LogLevel, cfg.LogFile, stdlog.LstdFlags|stdlog.Lshortfile)
	if err != nil {
		panic(err)
	}
	if logger.BaseFile != nil {
		lucky.SetLogOutput(io.MultiWriter(logger.BaseFile, os.Stdout))
	} else {
		lucky.SetLogOutput(os.Stdout)
	}
	if cfg.Pprof {
		go func() {
			addr := strings.Split(cfg.ListenAddr, ":")
			port, err := strconv.ParseInt(addr[1], 10, 64)
			if err != nil {
				log.Error("parse port error:%v", err)
			} else {
				err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port+1), nil)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	// 初始化Redis
	if cfg.RedisUrl != constans.EMPTY {
		if _, err := cache.OpenRedis(cfg.RedisUrl); err != nil {
			panic(err)
		}
		log.Release("redis cache running =>%s", cfg.RedisUrl)
	} else {
		log.Error("redis is not connected.")
	}
	if cfg.MysqlUrl != constans.EMPTY {
		// 初始化Mysql
		if _, err := xdb.OpenMysqlDB(xdb.MYSQL, cfg.MysqlUrl, &gorm.Config{}, cfg.IdleConnSize, cfg.MaxConnSize, entity.GetTables()...); err != nil {
			panic(err)
		}
		// QPS DB
		xdb.InitQpsDB(cfg.IdleConnSize, time.Second)
		log.Release("qps db running =>%s", cfg.MysqlUrl)
	} else {
		log.Error("mysql db is not connected.")
	}
	lucky.SetLogLv(cfg.LogLevel)
	lucky.EnableCrossOrigin()
	router.InitRouter(nil)
	// CPU 核心数
	if cfg.CpuNum > 0 {
		runtime.GOMAXPROCS(cfg.CpuNum)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	// 账号监听服务
	log.Fatal("%v", lucky.Run(cfg.ListenAddr))
}
