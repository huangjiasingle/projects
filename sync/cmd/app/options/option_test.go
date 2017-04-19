package options

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	path := `E:\workspace\projects\src\sync\cmd\config.json`
	InitConfig(path)
	if GlobalConfig.MySql.Dsn == "" && GlobalConfig.Redis.Addr == "192.168.0.40:6379" && GlobalConfig.Redis.Password == "" && GlobalConfig.Redis.DB == 0 && GlobalConfig.Redis.RequiredPW == true && GlobalConfig.Redis.PoolSize == 100 {
		t.Log("load config successed!")
	}
}
