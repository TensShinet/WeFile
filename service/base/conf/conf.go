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
	Redis   common.RedisConfig   `yaml:"redis"`
	BaseAPI BaseAPIConfig        `yaml:"base_api"`
}

type BaseAPIConfig struct {
	Address        string `yaml:"address"`          // 监听地址
	Salt           string `yaml:"salt"`             // 密码盐值
	SessionSecrete string `yaml:"session_secrete"`  // session secrete
	SessionMaxAge  int    `yaml:"session_max_age"`  // session 存活时间 单位分钟
	SessionName    string `yaml:"session_name"`     // cookie 中 session name
	FileAPIAddress string `yaml:"file_api_address"` // file api 地址
}

var (
	once   sync.Once
	c      = &Config{}
	logger = logging.GetLogger("base_service_conf")
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
