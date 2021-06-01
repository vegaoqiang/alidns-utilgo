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
	List		bool
	DN 	   		string
	Type	   	string
	Value	   	string
	Search		string
	AccessKey  	string
	AccessSecret string
	Region 	   	string
	Config		string
)

func init(){
	flag.BoolVar(&Add, "add", false, "添加解析")
	flag.BoolVar(&Del, "del", false, "删除解析")
	flag.BoolVar(&List, "list", false, "获取域名所有解析")
	flag.StringVar(&DN, "dn", "", "需要解析的完整域名, 如: bar.foo.com, 当指定--list参数时，dn为不包含'主机记录'的域名时，如：foo.com， 则获取所有该域名的解析记录")
	flag.StringVar(&Type, "type", "A", "解析记录类型,参见:https://help.aliyun.com/document_detail/29805.html?spm=api-workbench.API%20Document.0.0.4fbd1e0fFdFBGG")
	flag.StringVar(&Type, "t", "A", "与--type相同")
	flag.StringVar(&Value, "value", "", "解析记录值")
	flag.StringVar(&Value, "v", "", "与--value相同")
	flag.StringVar(&Search, "search", "", "指定搜索域名解析的关键字")
	flag.StringVar(&Search, "s", "", "与--search相同")
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
		if err := initDomainConfig().AddDomainRecord(client); err != nil {
			fmt.Println("添加解析失败：", err)
		}
		//todo: 展示成功解析的域名信息
	}
	if List {
		ldnconfig, err := initListDomainConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(DN) != 0 {
			if result, err := ldnconfig.ListDomainRecord(client); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}else{
				fmt.Println(result.Body.DomainRecords.Record)
			}
		}else{
			if result, err := ldnconfig.ListDomain(client); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}else {
				fmt.Println(result)
			}
		}
		
	}
}

func initDomainConfig() *aliutils.DomainConfig{
	if len(DN) == 0 {
		fmt.Println("请指定完整域名，如: bar.foo.com")
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
	if domainLen := len(strings.Split(dn[1], ".")); domainLen < 2 {
		// 默认dn[0]为三级域名，dn[1]为二级域名，验证dn[1]是否由一个"."组成
		// 如果dn[1]不包含一个".",则用户提供的DN不是一个合法的参数,如用户提供
		// foo.com，去掉dn[0],则dn[1]为com，则该DN不合法
		fmt.Printf("提供的DN参数%v不正确，请按照如下：bar.foo.com提供", DN)
		os.Exit(1)
	}
	dnconfig := &aliutils.DomainConfig{
		RR: 		dn[0],
		DomainName: dn[1],
		Type: 		Type,
		Value: 		Value,	
	}
	return dnconfig	
}

// 根据用户提供的DN获取DomainName和RRKeyWord等信息，其中，当用指定了search参数是，RRKeyWord将失效
// 以search指定的为准
func initListDomainConfig() (*aliutils.ListDomainConfig, error) {
	if len(DN) == 0 {
		// 用户未指定DN，该情况将查询域名列表，同时可提供搜索关键字
		ldnconfig := &aliutils.ListDomainConfig{
			KeyWord: Search,
		}
		return ldnconfig, nil
	}
	domianLen := len(strings.Split(DN, "."))
	if domianLen >= 3 {
		// 此时的域名类型bar.foo.com
		dn := strings.SplitN(DN, ".", 2)
		if len(Search) == 0 {
			// 用户未指定search
			ldnconfig := &aliutils.ListDomainConfig{
				RRKeyWord: 		dn[0],
				DomainName: 	dn[1],
				TypeKeyWord: 	Type,
				ValueKeyWord: 	Value,
			}
			return ldnconfig, nil
		}else{
			// 用户指定了search
			ldnconfig := &aliutils.ListDomainConfig{
				RRKeyWord: 		dn[0], // 所搜时将失效
				DomainName: 	dn[1],
				TypeKeyWord: 	Type,
				ValueKeyWord: 	Value,
				KeyWord: 		Search,
			}
			return ldnconfig, nil
		}
	}else if domianLen == 2 {
		// 此时的域名类型foo.com
		// 当用户此时未指定search关键字，将获取foo.com下所有的解析
		ldnconfig := &aliutils.ListDomainConfig{
			DomainName: 	DN,
			TypeKeyWord: 	Type,
			ValueKeyWord: 	Value,
			KeyWord: 		Search,
		}
		return ldnconfig, nil
	}
	return &aliutils.ListDomainConfig{}, fmt.Errorf("提供的域名%v不正确", DN)
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
