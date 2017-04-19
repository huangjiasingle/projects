package pkg

import (
	"github.com/astaxie/beego/logs"
)

func InitLog() {
	logs.SetLogger(logs.AdapterFile, Config.File)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterMail, Config.Smtp)
	logs.SetLevel(7)
}
