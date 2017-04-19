package options

import (
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
)

//Config is the application's config
type Config struct {
	*MySql   `json:"mysql"`
	*Redis   `json:"redis"`
	Interval int64  `json:"interval"`
	File     string `json:"file"`
	Smtp     string `json:"smtp"`
}

//Mysql is the mysql's config
type MySql struct {
	Dsn string `json:"dsn"`
}

//Redis is the redis's config
type Redis struct {
	//host:port
	Addr     string `json:"addr"`
	Password string `json:"password"`
	//require password or not
	RequiredPW bool `requiredPW`
	//select the number's redis db
	DB int `json:"db"`
	//redis connection pool size
	PoolSize int `json:"poolSize"`
}

var (
	GlobalConfig = &Config{}
)

func InitConfig(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error("read the config file:%v err:%v", path, err)
	}

	if err = json.Unmarshal(contents, GlobalConfig); err != nil {
		logs.Critical("Unmarshal the config file:%v err:%v", path, err)
	}
}

func InitLog() {
	logs.SetLogger(logs.AdapterFile, GlobalConfig.File)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterMail, GlobalConfig.Smtp)
	logs.SetLevel(7)
}
