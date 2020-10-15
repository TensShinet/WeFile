package conf

import (
	"github.com/TensShinet/WeFile/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"time"
)

type Config struct {
	LogLevel string        `yaml:"log_level"` // log 级别
	Etcd     EtcdConfig    `yaml:"etcd"`      // etcd 配置
	DB       DBConfig      `yaml:"db"`        // DB 配置
	Redis    RedisConfig   `yaml:"redis"`     // redis 配置
	NodeID   int           `yaml:"node_id"`   // 节点 id
	JWT      JWTConfig     `yaml:"jwt"`       // jwt 配置
	Service  ServiceConfig `yaml:"service"`   // service 配置
	BaseAPI  BaseAPIConfig `yaml:"base_api"`  // base api 配置
}

type BaseAPIConfig struct {
	Address        string `yaml:"address"`         // 监听地址
	Salt           string `yaml:"salt"`            // 密码盐值
	SessionSecrete string `yaml:"session_secrete"` // session secrete
	SessionMaxAge  int    `yaml:"session_max_age"` // session 存活时间 单位分钟
	SessionName    string `yaml:"session_name"`    // cookie 中 session name
}

type ServiceConfig struct {
	RegisterTTL      int `yaml:"register_ttl"`      // 服务注册超时时间 单位秒
	RegisterInterval int `yaml:"register_interval"` // 报告状态的时间间隔 单位秒
}

// TODO: JWT 详细配置
type JWTConfig struct {
	ValidTime int    `yaml:"valid_time"` // jwt token 有效时间 单位秒
	Secret    string `yaml:"secret"`     // jwt 秘钥
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
	Network  string        `yaml:"network"`
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

		if Conf.Service.RegisterInterval == 0 {
			Conf.Service.RegisterInterval = 5
		}

		if Conf.Service.RegisterTTL == 0 {
			Conf.Service.RegisterTTL = 5
		}

	})
}

func GetConfig() *Config {
	return Conf
}
