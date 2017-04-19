package logger

import (
	"os"

	"github.com/op/go-logging"
)

func NewLog() *logging.Logger {
	log := logging.MustGetLogger("")
	format := logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05} %{level:.5s} > %{message} [%{shortfile}] %{color:reset}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(-1, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	return log
}
