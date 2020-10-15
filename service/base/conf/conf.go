package conf

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
)

func Init() {
	conf.Init("base_conf.yml")
	config := conf.GetConfig()
	logging.SetGlobalLevel(logging.GetLevel(config.LogLevel))
}
