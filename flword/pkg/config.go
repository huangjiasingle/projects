package pkg

import (
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
)

type MySql struct {
	Dsn      string `json:"dsn"`
	Addr     string `json:addr`
	Interval int    `json:interval` //定时器执行的周期
	File     string `json:"file"`
	Smtp     string `json:"smtp"`
}

var (
	Config = &MySql{}
)

func InitConfig(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Critical("read the config file:%v err:%v", path, err)
	}

	if err = json.Unmarshal(contents, Config); err != nil {
		logs.Critical("Unmarshal the config file:%v err:%v", path, err)
	}
}
