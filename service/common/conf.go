package common

import "time"

// 所有配置的基础配置
type ServiceConfig struct {
	LogLevel         string     `yaml:"log_level"`         // log 级别
	NodeID           int        `yaml:"node_id"`           // 节点 id
	RegisterTTL      int        `yaml:"register_ttl"`      // 服务注册超时时间 单位秒
	RegisterInterval int        `yaml:"register_interval"` // 报告状态的时间间隔 单位秒
	Etcd             EtcdConfig `yaml:"etcd"`              // etcd 配置
}

// TODO:etcd 详细配置
type EtcdConfig struct {
	EndPoints []string `yaml:"end_points"`
}

// redis 配置
type RedisConfig struct {
	Network  string        `yaml:"network"`
	Enabled  bool          `yaml:"enabled"`
	Conn     string        `yaml:"conn"`
	Password string        `yaml:"password"`
	DBNum    int           `yaml:"db_num"`
	Timeout  int           `yaml:"timeout"`
	Sentinel redisSentinel `yaml:"sentinel"`
}

// redis 集群配置
type redisSentinel struct {
	Enabled bool     `yaml:"enabled"`
	Master  string   `yaml:"master"`
	Nodes   []string `yaml:"nodes"`
}

// mysql 配置
type MySQLConfig struct {
	DSN               string        `yaml:"dsn"` // mysql dsn "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	MaxIdleConnection int           `yaml:"max_idle_connection"`
	MaxOpenConnection int           `yaml:"max_open_connection"`
	ConnMaxLifetime   time.Duration `yaml:"conn_max_lifetime"`
	Enabled           bool          `yaml:"enabled"` // 是否启用 mysql
}

var (
	// default 配置
	//
	// 没有为 0 的节点
	defaultNodeID        = 1
	defaultEtcdEndPoints = []string{"127.0.0.1:2379"}
)

func (conf *ServiceConfig) Init() {
	if conf.NodeID == 0 {
		conf.NodeID = defaultNodeID
	}
	if conf.Etcd.EndPoints == nil {
		conf.Etcd.EndPoints = defaultEtcdEndPoints
	}
	if conf.RegisterInterval == 0 {
		conf.RegisterInterval = 5
	}

	if conf.RegisterTTL == 0 {
		conf.RegisterTTL = 5
	}
}
