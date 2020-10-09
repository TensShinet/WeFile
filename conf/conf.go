package conf

import (
	"github.com/TensShinet/WeFile/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	LogLevel string     `yaml:"log_level"` // log 级别
	Etcd     EtcdConfig `yaml:"etcd"`      // etcd 配置
	NodeID   int        `yaml:"node_id"`   // 节点 id
}

// TODO:etcd 详细配置
type EtcdConfig struct {
	EndPoints []string `yaml:"end_points"`
}

var (
	logger = logging.GetLogger("conf")
	Conf   *Config
	once   sync.Once
	// default 配置
	//
	// 没有为 0 的节点
	defaultNodeID        = 1
	defaultEtcdEndPoints = []string{"127.0.0.1:2379"}
)

func GetConfig(filepath string) *Config {
	once.Do(func() {
		if filepath == "" {
			filepath = "conf.yml"
		}
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			logger.Errorf("read file failed, for the reason:%v", err.Error())
		}
		Conf = &Config{}
		if err := yaml.Unmarshal(file, Conf); err != nil {
			logger.Errorf("read file failed, for the reason:%v", err.Error())
		}
		if Conf.NodeID == 0 {
			Conf.NodeID = defaultNodeID
		}

		if Conf.Etcd.EndPoints == nil {
			Conf.Etcd.EndPoints = defaultEtcdEndPoints
		}
	})
	return Conf
}
