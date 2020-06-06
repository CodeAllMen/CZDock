package models

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config 内容站配置
type Config struct {
	Service map[string]ServiceInfo
}

type ServiceInfo struct {
	ServiceID           string `yaml:"service_id" orm:"pk;column(service_id)"`
	ContentUrl          string `yaml:"content_url"`
	ServerUrl           string `yaml:"server_url"`
	UserApiUrl          string `yaml:"user_api_url"`
	MerchantPassword    string `yaml:"merchant_password"`
	NotificationUrl     string `yaml:"notification_url"`
	StartSubReturnUrl   string `yaml:"start_sub_return_url"`
	MerchantId          string `yaml:"merchant_id"`
	ServiceName         string `yaml:"service_name"`
	ServerOrder         string `yaml:"server_order"`
	LimitSubNum         int    `yaml:"limit_sub_num"`
	ShortCode           string `yaml:"short_code"`
	PromptContentArgsCs string `yaml:"prompt_content_args_cs"`
	PromptContentArgsEn string `yaml:"prompt_content_args_en"`
	DockUrl             string `yaml:"dock_url"`
	Keyword             string `yaml:"keyword"`
	Price               string `yaml:"price"`
	Country             string `yaml:"country"`
	Operator            string `yaml:"operator"`
	CampID              int    `yaml:"camp_id"`
}

const (
	WapIdentifyUser int = iota + 1
	GetUser
	WapAuthorize
	GetSubscription
	CloseSubscription
)

var ServiceData = make(map[string]ServiceInfo)

func (server *ServiceInfo) TableName() string {
	return "server_info"
}

func InitServiceConfig() {
	filename, _ := filepath.Abs("resource/config/conf.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}
	ServiceData = config.Service
}

type CommandParameter struct {
	Types          int
	TrackID        string
	IP             string
	Uid            string
	SessionID      string
	SubscriptionId string
}
