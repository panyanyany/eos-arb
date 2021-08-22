package log_util

import (
	"fmt"
	"os"
	"path"

	"eos-arb/util/ip_util"

	"github.com/cihub/seelog"
)

// custom formatter，显示进程名
func createProcessFormatter(params string) seelog.FormatterFunc {
	return func(message string, level seelog.LogLevel, context seelog.LogContextInterface) interface{} {
		return os.Args[0]
	}
}
func createShortProcessFormatter(params string) seelog.FormatterFunc {
	return func(message string, level seelog.LogLevel, context seelog.LogContextInterface) interface{} {
		_, filename := path.Split(os.Args[0])
		return filename
	}
}

var ip string

func createIpFormatter(params string) seelog.FormatterFunc {
	return func(message string, level seelog.LogLevel, context seelog.LogContextInterface) interface{} {
		if ip != "" {
			return ip
		}
		ip, err := ip_util.GetMyIP()
		if err != nil {
			ip = fmt.Sprintf("<IP:%v>", err)
			return ip
		}
		return ip
	}
}

func SetupSeelog() {
	var err error
	// custom formatter，显示进程名
	err = seelog.RegisterCustomFormatter("Process", createProcessFormatter)
	if err != nil {
		panic(err)
	}
	err = seelog.RegisterCustomFormatter("Proc", createShortProcessFormatter)
	if err != nil {
		panic(err)
	}
	err = seelog.RegisterCustomFormatter("IP", createIpFormatter)
	if err != nil {
		panic(err)
	}
	// seelog配置，用于输出调试信息
	logger, err := seelog.LoggerFromConfigAsFile("config/seelog.xml")
	if err != nil {
		panic(err)
	}
	err = seelog.UseLogger(logger)
	if err != nil {
		panic(err)
	}
}
