package config

import (
	"fmt"
	"github.com/helloh2o/lucky/log"
	v12ctx "github.com/kataras/iris/v12/context"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type Conf struct {
	// mysql 数据库连接
	ServerId     string `yaml:"server_id"  json:"server_id"`
	MysqlUrl     string `yaml:"mysql_url" json:"mysql_url"`
	CnfDBUrl     string `yaml:"cnf_db_url" json:"cnf_db_url"`
	IdleConnSize int    `yaml:"idle_conn_size" json:"idle_conn_size"`
	MaxConnSize  int    `yaml:"max_conn_size"  json:"max_conn_size"`
	// Redis 连接地址
	RedisUrl string `yaml:"redis_url" json:"redis_url"`
	// 服务监听地址
	ListenAddr string `yaml:"listen_addr" json:"listen_addr"`
	// 日志级别
	LogLevel string `yaml:"log_level" json:"log_level"`
	LogFile  string `yaml:"log_file" json:"log_file"`
	// WEB HOST [域名连接，生成图片连接，apk连接]
	WebHost string `yaml:"web_host" json:"web_host"`
	// 签名验证
	Signature string `yaml:"signature" json:"signature"`
	// CPU 数量
	CpuNum int `yaml:"cpu_num" json:"cpu_num"`
	// pprof
	Pprof bool `yaml:"pprof" json:"pprof"`
	// nats queue url
	NatsUrl string `yaml:"nats_url" json:"nats_url"`
	// 是否维护中
	Maintenance bool `yaml:"maintenance" json:"maintenance"`
	// 维护中能测试的IP
	MaintenanceIp string `yaml:"maintenance_ip" json:"maintenance_ip"`
}

var (
	c     = new(Conf)
	cPath string
	cLock sync.RWMutex
)

type ConfServer struct {
	Id     uint32 `yaml:"id"`     // 网关ID（llk=1）
	Type   uint32 `yaml:"type"`   // 类型 (1聊天，2游戏)
	Host   string `yaml:"host"`   // 目标服务器
	Port   int32  `yaml:"port"`   // 目标端口
	Scheme string `yaml:"scheme"` // 连接协议 (ws,wss,http,https,tcp,udp,quic)
	Proto  string `yaml:"proto"`  // 消息协议 (json & protobuf)
}

// 初始化配置
func Initialize(path string) *Conf {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		panic(err)
	}
	cPath = path
	return c
}

// ReloadConf 重新加载配置
func ReloadConf(ctx *v12ctx.Context) {
	cLock.Lock()
	defer cLock.Unlock()
	var temp Conf
	f, err := os.Open(cPath)
	if err != nil {
		log.Error("reload read conf err:%v", err)
		ctx.WriteString(err.Error())
		return
	} else {
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&temp)
		if err != nil {
			log.Error("reload decode conf err:%v", err)
			ctx.WriteString(err.Error())
			return
		} else {
			c = &temp
		}
	}
	// 日志级别
	log.SetLogLevelDefault(c.LogLevel)
	msg := fmt.Sprintf("===> reload ok, log level:%s <===\n", c.LogLevel)
	log.Release("%s", msg)
	ctx.WriteString(msg)
}

func Get() *Conf {
	cLock.RLock()
	rt := c
	defer cLock.RUnlock()
	return rt
}
