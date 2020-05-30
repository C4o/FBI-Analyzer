package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	C   string
	Cfg Config
)

type Config struct {
	// redis配置
	RedAddr string `yaml:"redis"`
	RedPass string `yaml:"password"`
	DB      int    `yaml:"db"`
	// kafka配置
	Broker  string   `yaml:"broker"`
	GroupID string   `yaml:"groupid"`
	Topic   []string `yaml:"topic"`
	Offset  string   `yaml:"offset"`
	// 日志配置
	Path string `yaml:"path"`
}

func Load() error {

	conf := new(Config)
	file, err := ioutil.ReadFile(C)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		return err
	}
	Cfg = *conf
	return nil
}

func New(cfg string) error {

	C = cfg
	if err := Load(); err != nil {
		log.Printf("[ERROR] load config file error : %v", err)
		return err
	}
	return nil
}
