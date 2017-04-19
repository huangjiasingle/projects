package main

import (
	"flag"
	"runtime"

	"sync/cmd/app"
	"sync/cmd/app/options"
	"sync/pkg/util"
)

var (
	path = flag.String("config", "config.json", "please use --config=config file's path")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	options.InitConfig(*path)
	options.InitLog()
	util.InitDB()
	util.InitRedisClient()
	app.Run()
	select {}
}
