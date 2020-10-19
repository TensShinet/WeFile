package conf

import (
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	Service common.ServiceConfig `yaml:"service"`
	FileAPI FileAPIConfig        `yaml:"file_api"`
}

type FileAPIConfig struct {
	Address        string `yaml:"address"`          // 监听地址
	LocalTempStore string `yaml:"local_temp_store"` // 暂存地址
}

var (
	once   sync.Once
	c      = &Config{}
	logger = logging.GetLogger("file_service_conf")
)

func Init(filepath string) {
	once.Do(func() {
		if filepath == "" {
			filepath = "conf.yml"
		}
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			logger.Errorf("read file failed, for the reason:%v", err.Error())
		}
		if err := yaml.Unmarshal(file, c); err != nil {
			logger.Errorf("read file failed, for the reason:%v", err.Error())
		}
		c.Service.Init()
		// 是否开启 debug 模式
		logging.SetGlobalLevel(logging.GetLevel(c.Service.LogLevel))
	})
}

func GetConfig() *Config {
	return c
}
