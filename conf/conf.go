package conf

import (
	"github.com/TensShinet/WeFile/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"time"
)

type Config struct {
	LogLevel string      `yaml:"log_level"` // log 级别
	Etcd     EtcdConfig  `yaml:"etcd"`      // etcd 配置
	DB       DBConfig    `yaml:"db"`        // DB 配置
	Redis    RedisConfig `yaml:"redis"`     // redis 配置
	NodeID   int         `yaml:"node_id"`   // 节点 id
}

// TODO: DB 详细配置
type DBConfig struct {
	MySQL MySQLConfig `yaml:"mysql"` // mysql 配置
}

// TODO:etcd 详细配置
type EtcdConfig struct {
	EndPoints []string `yaml:"end_points"`
}

type MySQLConfig struct {
	DSN               string        `yaml:"dsn"` // mysql dsn "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	MaxIdleConnection int           `yaml:"max_idle_connection"`
	MaxOpenConnection int           `yaml:"max_open_connection"`
	ConnMaxLifetime   time.Duration `yaml:"conn_max_lifetime"`
	Enabled           bool          `yaml:"enabled"` // 是否启用 mysql
}

type RedisConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Conn     string        `yaml:"conn"`
	Password string        `yaml:"password"`
	DBNum    int           `yaml:"db_num"`
	Timeout  int           `yaml:"timeout"`
	Sentinel redisSentinel `yaml:"sentinel"`
}

type redisSentinel struct {
	Enabled bool     `yaml:"enabled"`
	Master  string   `yaml:"master"`
	Nodes   []string `yaml:"nodes"`
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

func Init(filepath string) {
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
}

func GetConfig() *Config {
	return Conf
}
