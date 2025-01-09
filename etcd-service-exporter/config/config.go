package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

var (
	conf = new(Config)
	once sync.Once
)

type Etcd struct {
	Endpoints   []string `mapstructure:"endpoints"`
	DialTimeout int      `mapstructure:"dialTimeout"`
}

type Config struct {
	*Etcd `mapstructure:"etcd"`
}

func mustInitConfig() {
	configFilePath := flag.String("config", "../conf/config.yaml", "配置文件路径")
	viper.SetConfigFile(*configFilePath)

	viper.SetDefault("etcd.dialTimeout", 5)
	viper.SetDefault("etcd.endpoints", []string{"localhost:2379"})

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}

	// 将配置文件解析到结构体中
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("unmarshal to conf failed, err:%v", err))
	}
}

// Get 用于获取配置文件
func Get() *Config {
	once.Do(mustInitConfig)
	return conf
}
