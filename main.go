package main

import (
	"alidns-utilgo/aliutils"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)


var (
	Add			bool
	Del			bool
	DN 	   		string
	Type	   	string
	Value	   	string
	AccessKey  	string
	AccessSecret string
	Region 	   	string
	Config		string
)

func init(){
	flag.BoolVar(&Add, "add", false, "添加解析")
	flag.BoolVar(&Del, "del", false, "删除解析")
	flag.StringVar(&DN, "dn", "", "需要解析的完整域名, 如: www.baidu.com")
	flag.StringVar(&Type, "type", "A", "解析记录类型,参见:https://help.aliyun.com/document_detail/29805.html?spm=api-workbench.API%20Document.0.0.4fbd1e0fFdFBGG")
	flag.StringVar(&Type, "t", "A", "与--type相同")
	flag.StringVar(&Value, "value", "", "解析记录值")
	flag.StringVar(&Value, "v", "", "与--value相同")
	flag.StringVar(&AccessKey, "accesskey", "", "指定access_key")
	flag.StringVar(&AccessSecret, "accesssecret", "", "指定access_secret")
	flag.StringVar(&Region, "region", "", "指定Region")
	flag.StringVar(&Config, "config", "", "手动指定配置文件")
	flag.StringVar(&Config, "c", "", "与--config相同")
}


func main(){
	flag.Parse()
	client, err := initAccountConfig().CreateClient()
	if err != nil {
		fmt.Println("初始化客户端失败", err)
		os.Exit(1)
	}
	if Add {
		if err := initDomainConfig().AddDomain(client); err != nil {
			fmt.Println("添加解析失败：", err)
		}
		//todo: 展示成功解析的域名信息
	}
	if Del {
		
	}
}

func initDomainConfig() *aliutils.DomainConfig{
	if len(DN) == 0 {
		fmt.Println("请指定完整域名，如: www.baidu.com")
		os.Exit(1)
	}
	// if len(Type) == 0 {
	// 	fmt.Println("请指定解析记录类型")
	// 	os.Exit(1)
	// }
	if len(Value) == 0 {
		fmt.Println("请指定记录值")
		os.Exit(1)
	}
	ip := net.ParseIP(Value)
	if ip == nil {
		fmt.Println("提供的IP地址不合法")
		os.Exit(1)
	}
	dn := strings.SplitN(DN, ".", 2)
	dnconfig := &aliutils.DomainConfig{
		RR: 		dn[0],
		DomainName: dn[1],
		Type: 		Type,
		Value: 		Value,	
	}
	return dnconfig	
}

// 	初始化账户信息，优先从命令行获取，命令行未指定则从配置文件中获取
func initAccountConfig() *aliutils.Account {
	var configFile string
	if len(Config) != 0 {
		configFile = Config
	}else{
		home := os.Getenv("HOME")
		configFile = home + "/.alidns-utilgo/config.json"
	}
	accountConfig := loadAccountConfig(configFile)
	if len(AccessKey) != 0 {
		accountConfig.AccessKey = AccessKey
	}
	if len(AccessSecret) != 0 {
		accountConfig.AccessSecret = AccessSecret
	}
	if len(Region) != 0 {
		accountConfig.Region = Region
	}
	if len(accountConfig.AccessKey) == 0 {
		fmt.Println("请提供access_key")
		os.Exit(1)
	}
	if len(accountConfig.AccessSecret) == 0 {
		fmt.Println("请提供access_secret")
		os.Exit(1)
	}
	if len(accountConfig.Region) == 0 {
		fmt.Println("请提供region")
		os.Exit(1)
	}
	return accountConfig
}

// 从配置文件中读取账户信息
func loadAccountConfig(config string) *aliutils.Account {
	bytes, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println("指定的配置文件不存在:", config)
		os.Exit(1)
	}
	account := &aliutils.Account{}
	if err := json.Unmarshal(bytes, account); err != nil {
		fmt.Println("解析配置文件发生错误")
		os.Exit(2)
	}
	return account
}

func AddDNS() {

}