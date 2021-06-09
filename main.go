package main

import (
	"alidns-utilgo/aliutils"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	Init		bool
	Add			bool
	Del			bool
	Update		bool
	List		bool
	U			string
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
	flag.BoolVar(&Init, "init", false, "初始化账户配置")
	flag.BoolVar(&Add, "add", false, "添加解析记录")
	flag.BoolVar(&Del, "del", false, "删除解析记录")
	flag.BoolVar(&Update, "update", false, "更新域名解析")
	flag.BoolVar(&List, "list", false, "获取域名所有解析")
	flag.StringVar(&U, "u", "", "需要更新的字段已经值，如：-u value=bar将修改解析主机记录为bar，多个字段值以,分割")
	flag.StringVar(&DN, "dn", "", "需要解析的完整域名, 如: bar.foo.com, 当指定--list参数时，dn为不包含'主机记录'的域名时，如：foo.com， 则获取所有该域名的解析记录")
	flag.StringVar(&Type, "type", "", "解析记录类型,default A,参见:https://help.aliyun.com/document_detail/29805.html?spm=api-workbench.API%20Document.0.0.4fbd1e0fFdFBGG")
	flag.StringVar(&Type, "t", "", "与--type相同")
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
	if Init {
		createAccountConfig()
	}
	accountConfig, err := initAccountConfig()
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}
	client, err := accountConfig.CreateClient()
	if err != nil {
		fmt.Println("初始化客户端失败", err)
		os.Exit(1)
	}
	if Add {
		adrconfig := initDomainConfig()
		if err := adrconfig.AddDomainRecord(client); err != nil {
			fmt.Println("添加解析失败：", err)
			os.Exit(1)
		}
		// 输出添加解析记录成功信息
		fmt.Printf("域名:    %s\n", adrconfig.DomainName)
		fmt.Printf("记录类型: %s\n", adrconfig.Type)
		fmt.Printf("主机记录: %s\n", adrconfig.RR)
		fmt.Printf("记录值:   %s\n", adrconfig.Value)
		os.Exit(0)
	}
	if Del {
		dsdrconfig := initDelSubDomainRecordsConfig()
		if result, err := dsdrconfig.DelSubDomainRecords(client); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}else {
			// 输出删除解析记录成功信息
			fmt.Printf("\033[33m已删除解析记录个数: %s\033[0m\n", *result.Body.TotalCount)
			fmt.Printf("域名:%s\n", dsdrconfig.DomainName)
			fmt.Printf("记录类型: %s\n", dsdrconfig.Type)
			fmt.Printf("主机记录: %s\n", dsdrconfig.RR)
			fmt.Printf("记录值:   %s\n", dsdrconfig.Value)
		}
		os.Exit(0)
	}
	if Update {
		ldrconfig := beforeUpdateDomainRecordConfig()
		result, err := ldrconfig.ListDomainRecords(client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// 得到目标子域名RecordId
		if *result.Body.TotalCount == 0 {
			fmt.Printf("未查找到主机记录为: %v, 解析类型为: %v 的解析记录\n",DN, Type)
			os.Exit(1)
		}
		if *result.Body.TotalCount > 1 {
			fmt.Printf("根据提供的DN: %v 查找多过个解析记录,请使用-t参数指定解析记录类型尝试确定一个解析记录",DN)
			os.Exit(1)
		}
		RecordId := result.Body.DomainRecords.Record[0].RecordId
		udrconfig := initUpdateDomainRecordConfig(ldrconfig)
		if result, err := udrconfig.UpdateDomainRecords(client, RecordId); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}else {
			fmt.Println(result)
			fmt.Printf("域名:%s\n", ldrconfig.DomainName)
			fmt.Printf("记录类型: %-20s %s %s\n", ldrconfig.TypeKeyWord, "->", udrconfig.Type)
			fmt.Printf("主机记录: %-20s %s %s\n", ldrconfig.RRKeyWord, "->", udrconfig.RR)
			fmt.Printf("记录值:   %-20s %s %s\n", ldrconfig.ValueKeyWord, "->", udrconfig.Value)
		}
		os.Exit(0)
	}
	if List {
		ldconfig, err := initListDomainConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(DN) != 0 {
			if result, err := ldconfig.DoListDomainRecords(client); err != nil {	
				fmt.Println(err)
				os.Exit(1)
			}else{
				fmt.Printf("\033[33m域名: %s, 解析记录: %d\033[0m\n", ldconfig.DomainName, len(result))
				fmt.Printf("%-30s%-10s%-20s%-20s%-10s%-20s%-20s\n", "主机记录", "记录类型", "解析线路", "记录值", "TTL", "状态", "备注")
				for _, v := range result {
					if v.Remark == nil {
						v.SetRemark(" ")
					}
					// 格式化输出，将以下输出内容同上面的标题对其
					// 4: 标题长度
					// 30: 标题占位符长度
					// len(xx): 内容长度
					// NullLength: 通过公式计算得出内容需要加上NullLength个空字符才能填满占位符
					RRNullLength := 4 + 30 - len(*v.RR)
					if RRNullLength < 0 {
						RRNullLength = 0
					}
					fmt.Printf("%-30s", *v.RR + strings.Repeat(" ", RRNullLength))
					TypeNullLength := 4 + 10 - len(*v.Type)
					fmt.Printf("%-10s", *v.Type + strings.Repeat(" ", TypeNullLength))
					LineNullLength := 4 + 20 - len(*v.Line)
					fmt.Printf("%-20s", *v.Line + strings.Repeat(" ", LineNullLength))
					ValueNullLength := 3 + 20 - len(*v.Value)
					if ValueNullLength < 0 {
						ValueNullLength = 0
					}
					fmt.Printf("%-20s", *v.Value + strings.Repeat(" ", ValueNullLength))
					TTLNullLength := 10 - len(strconv.FormatInt(*v.TTL, 10))
					fmt.Printf("%-10s", strconv.FormatInt(*v.TTL, 10) + strings.Repeat(" ", TTLNullLength))
					StatusNullLength := 2 + 20 - len(*v.Status)
					fmt.Printf("%-20s", *v.Status + strings.Repeat(" ", StatusNullLength))
					fmt.Printf("%-20s\n", *v.Remark)
				}
			}
		}else{
			// 获取账户下所有域名
			// if result, err := ldconfig.ListDomains(client); err != nil {
				if result, err := ldconfig.DoListDomains(client); err != nil {	
				fmt.Println(err)
				os.Exit(1)
			}else {
				// 格式化输出
				fmt.Printf("%-30s%-20s%-10s%-10s\n", "域名", "标签", "记录数", "创建时间")
				for _, v := range result {
					DRNullLength := 2 + 30 - len(*v.DomainName)
					if DRNullLength < 0 {
						DRNullLength = 0
					}
					fmt.Printf("%-30s", *v.DomainName + strings.Repeat(" ", DRNullLength))
					var tags []string
					for _, v := range v.Tags.Tag {
						tags = append(tags, *v.Key)
					}
					tagsString := strings.Join(tags, ",")
					TagNullLength := 2 + 20 - len(tagsString)
					if TagNullLength < 0 {
						TagNullLength = 0
					}
					fmt.Printf("%-20s", tagsString + strings.Repeat(" ", TagNullLength))
					RecordCountNullLength := 3 + 10 -len(strconv.FormatInt(*v.RecordCount, 10))
					fmt.Printf("%-10s", strconv.FormatInt(*v.RecordCount, 10) + strings.Repeat(" ", RecordCountNullLength))
					fmt.Printf("%-10s\n", *v.CreateTime)
				}
			}
		}
		os.Exit(0)
	}
	flag.PrintDefaults()
}

// 初始化域名对应的配置文件，提供一个域名解析的全部必要字段
func initDomainConfig() *aliutils.DomainConfig{
	if len(DN) == 0 {
		fmt.Println("请输入完整的子域名，如: bar.foo.com")
		os.Exit(1)
	}
	if len(Type) == 0 {
		Type = "A"
	}
	if len(Value) == 0 {
		fmt.Println("请指定记录值")
		os.Exit(1)
	}
	if Type == "A" {
		ip := net.ParseIP(Value)
		if ip == nil {
			fmt.Println("提供的IP地址不合法")
			os.Exit(1)
		}
	}
	dn := strings.Split(DN, ".")
	if len(dn) <= 2 {
		fmt.Println("请输入完整的子域名，如: bar.foo.com")
		os.Exit(1)
	}
	dnconfig := &aliutils.DomainConfig{
		RR: 		strings.Join(dn[:len(dn) - 2], "."),
		DomainName: strings.Join(dn[len(dn) - 2:], "."),
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
		ldconfig := &aliutils.ListDomainConfig{
			KeyWord: Search,
		}
		return ldconfig, nil
	}
	return checkDn()
}

// 初始化删除主机记录对应的解析记录的参数配置
func initDelSubDomainRecordsConfig() *aliutils.DomainConfig {
	if len(DN) == 0 {
		fmt.Println("请输入完整的子域名，如: bar.foo.com")
		os.Exit(1)
	}
	dn := strings.Split(DN, ".")
	if len(dn) <= 2 {
		fmt.Println("请输入完整的子域名，如: bar.foo.com")
		os.Exit(1)
	}
	if len(Type) == 0 {
		Type = "A"
	}
	// 将用户提供的DN按照'.'进行分割成数组，默认数组的最后两个元素为域名(即一级域名和二级域名)，其余的为子域名
	dsdrconfig := &aliutils.DomainConfig{
		DomainName: strings.Join(dn[len(dn) - 2:], "."),
		RR: 		strings.Join(dn[:len(dn) - 2], "."),
		Type: 		Type,
	}
	return dsdrconfig
}

// 根据用户提供的DN确定更新目标子域名
func beforeUpdateDomainRecordConfig() *aliutils.ListDomainConfig {
	if len(DN) == 0 {
		fmt.Println("请输入完整的子域名，如: bar.foo.com")
		os.Exit(1)
	}
	if len(U) == 0 {
		fmt.Println("请指定需要更新的字段和值，如value=bar")
		os.Exit(1)
	}
	dn := strings.Split(DN, ".")
	if len(dn) <= 2 {
		fmt.Println("请输入完整的子域名，如: bar.foo.com")
		os.Exit(1)
	}
	if len(Type) == 0 {
		Type = "A"
	}
	// 由于查找的子域名必须精确，设置搜索模式为EXACT，并设置搜索关键字KeyWord进行精确查找，RRKeyWord在查找中将不起作用
	ldrconfig := &aliutils.ListDomainConfig{
		DomainName: 	strings.Join(dn[len(dn) - 2:], "."),
		RRKeyWord:		strings.Join(dn[:len(dn) - 2], "."),
		KeyWord: 		strings.Join(dn[:len(dn) - 2], "."),
		TypeKeyWord: 	Type,
		ValueKeyWord: 	Value,
		SearchMode:		"EXACT",
		PageNumber: 	1,
	}
	return ldrconfig
}

// 初始化更新域名解析参数配置
func initUpdateDomainRecordConfig(ldrconfig *aliutils.ListDomainConfig) *aliutils.DomainConfig {
	u := strings.Split(U, ",")
	udrconfig := &aliutils.DomainConfig{}
	for _, kv := range u {
		if !strings.Contains(kv, "=") || strings.Count(kv, "=") > 1 {
			fmt.Println("需要更新的字段和值输入格式错误")
			os.Exit(1)
		}
		kvArray := strings.Split(kv, "=")
		if strings.ToTitle(kvArray[0]) == "RR" {
			udrconfig.RR = kvArray[1]
		}else if strings.ToTitle(kvArray[0]) == "VALUE" {
			udrconfig.Value = kvArray[1]
		}else if strings.ToTitle(kvArray[0]) == "TYPE" {
			udrconfig.Type = kvArray[1]
		}
	}
	// 检查各项参数是否为空，将空设置为默认值
	if len(udrconfig.RR) == 0 {
		udrconfig.RR = ldrconfig.RRKeyWord
	}
	if len(udrconfig.Value) == 0 {
		udrconfig.Value = ldrconfig.ValueKeyWord
	}
	if len(udrconfig.Type) == 0 {
		udrconfig.Type = ldrconfig.TypeKeyWord
	}
	return udrconfig
}

func checkDn() (*aliutils.ListDomainConfig, error) {
	domianLen := len(strings.Split(DN, "."))
	if domianLen >= 3 {
		// 此时的域名类型bar.foo.com
		dn := strings.Split(DN, ".")
		if len(Search) == 0 {
			// 用户未指定search
			ldnconfig := &aliutils.ListDomainConfig{
				RRKeyWord: 		strings.Join(dn[:len(dn) - 2], "."),
				DomainName: 	strings.Join(dn[len(dn) - 2:], "."),
				TypeKeyWord: 	Type,
				ValueKeyWord: 	Value,
			}
			return ldnconfig, nil
		}else{
			// 用户指定了search
			ldnconfig := &aliutils.ListDomainConfig{
				RRKeyWord: 		strings.Join(dn[:len(dn) - 2], "."), // 所搜时将失效
				DomainName: 	strings.Join(dn[len(dn) - 2:], "."),
				TypeKeyWord: 	Type,
				ValueKeyWord: 	Value,
				KeyWord: 		Search,
			}
			return ldnconfig, nil
		}
	}else if domianLen == 2 {
		// 此时的域名类型foo.com
		// 当用户此时未指定search关键字，将获取foo.com下所有的解析,注意：解析类型为A的所有解析
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
func initAccountConfig() (*aliutils.Account, error) {
	accountConfig := &aliutils.Account{}
	if len(AccessKey) != 0 {
		accountConfig.AccessKey = AccessKey
	}
	if len(AccessSecret) != 0 {
		accountConfig.AccessSecret = AccessSecret
	}
	if len(Region) != 0 {
		accountConfig.Region = Region
	}
	if !(len(accountConfig.AccessKey) >0 && len(accountConfig.AccessSecret) > 0 && len(accountConfig.Region) > 0) {
		var configFile string
		if len(Config) != 0 {
			configFile = Config
		}else{
			home := os.Getenv("HOME")
			configFile = home + "/.alidns-utilgo/config.json"
		}
		ac, err := loadAccountConfig(configFile)
		if err != nil {
			return ac, err
		}
		accountConfig = ac
	}
	if len(accountConfig.AccessKey) == 0 {
		return accountConfig, errors.New("请提供access_key")
	}
	if len(accountConfig.AccessSecret) == 0 {
		return accountConfig, errors.New("请提供access_secret")
	}
	if len(accountConfig.Region) == 0 {
		return accountConfig, errors.New("请提供region")
	}
	return accountConfig, nil
}

// 从配置文件中读取账户信息
func loadAccountConfig(config string) (*aliutils.Account, error) {
	bytes, err := ioutil.ReadFile(config)
	if err != nil {
		err := fmt.Sprintf("指定的配置文件不存在:%v\n", config)
		return &aliutils.Account{}, errors.New(err)
	}
	account := &aliutils.Account{}
	if err := json.Unmarshal(bytes, account); err != nil {
		err := fmt.Sprintf("解析配置文件发生错误:%v\n", config)
		return account, errors.New(err)
	}
	return account, nil
}

func createAccountConfig(){
	defaultDirPath := os.Getenv("HOME") + "/.alidns-utilgo"
	if _, err := os.Stat(defaultDirPath); err != nil {
		os.Mkdir(defaultDirPath, 0644)
	}
	defaultFilePath := defaultDirPath + "/config.json"
	var ak, as, re string
	fmt.Print("输入access_key:")
	fmt.Scanln(&ak)
	if len(ak) == 0 {
		fmt.Println("未输入access_key")
		os.Exit(1)
	}
	fmt.Print("输入access_secret:")
	fmt.Scanln(&as)
	if len(as) == 0 {
		fmt.Println("未输入access_secret")
		os.Exit(1)
	}
	fmt.Print("输入region[缺省值cn-hangzhou]:")
	fmt.Scanln(&re)
	if len(re) == 0 {
		re = "cn-hangzhou"
	}
	accountConfig := &aliutils.Account{
		AccessKey: 		ak,
		AccessSecret: 	as,
		Region: 		re,
	}
	data, err := json.MarshalIndent(accountConfig, "", "    ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fd, err := os.OpenFile(defaultFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := fd.Write(data); err != nil {
		fmt.Println("写入文件错误")
		os.Exit(1)
	}
	fd.Close()
	fmt.Println("初始化完成")
	os.Exit(0)
}
