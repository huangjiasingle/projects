package app

import (
	"time"

	"sync/cmd/app/options"
	"sync/pkg/types"
	"sync/pkg/util"

	"github.com/astaxie/beego/logs"
)

var (
	firstStart = true
)

func Run() {
	logs.Info("the sync application is begin runing")
	helper := new(types.Helper)
	tick := time.NewTicker(time.Second * time.Duration(options.GlobalConfig.Interval))
	var syncTime time.Time
	var err error
	for {
		select {
		case <-tick.C:
			if firstStart {
				firstStart = false
				syncTm, _ := util.Client.Get("sync-time")
				if syncTm != "" {
					syncTime, err = time.Parse("2006-01-02 15:04:05", syncTm)
					if err != nil {
						logs.Error("pasre last sync time string to time err: %v", err)
						continue
					}
				}
			} else {
				syncTime = time.Now()
			}

			if _, err = util.Client.Set("sync-time", syncTime.Format("2006-01-02 15:04:05")); err != nil {
				logs.Error("set sync-time to redis err:%v", err)
				continue
			}
			err = helper.SyncDbToRedis(syncTime)
			if err != nil {
				logs.Error("sync the user to redis failed,the err:%v", err)
				continue
			}
			logs.Info("sync the user to redis successed")
		}
	}
}
