package aliutils

import (
	alidns "github.com/alibabacloud-go/alidns-20150109/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

type Account struct{
	AccessKey  		string 		`json:"access_key"`
	AccessSecret 	string		`json:"access_secret"`
	Region 	   		string		`json:"region,omitempty"`
}

type DomainConfig struct {
	DomainName	string
	RR			string
	Type	   	string
	Value	   	string
}


// func CreateClient(accountConfig, *Account) (Result *alidns.Client, err error) {
// 	config := &openapi.Config{
// 		// 您的AccessKey ID
// 		AccessKeyId: accountConfig.AccessKey,
// 		// 您的AccessKey Secret
// 		AccessKeySecret: accountConfig.AccessSecret,
// 	  }
// 	  // 访问的域名
// 	  // config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
// 	  // Result = &alidns.Client{}
// 	  Result, err = alidns.NewClient(config)
// 	  return Result, err
// }

func (account *Account) CreateClient() (result *alidns.Client, err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &account.AccessKey,
		// 您的AccessKey Secret
		AccessKeySecret: &account.AccessSecret,
	}
	result, err = alidns.NewClient(config)
	return result, err
}

// func AddDomain(client *openapi.Client, config *Account) (err error){
// 	addDomainRecordRequest := &alidns.AddDomainRecordRequest{
// 		DomainName: tea.String(config.DomainName),
// 		RR:			tea.String(config.RR),
// 		Type: 		tea.String(config.Type),
// 		Value: 		tea.String(config.Value),			
// 	}
// 	return client.addDomainRecord(addDomainRecordRequest)
// }

func (dnconfig *DomainConfig) AddDomain(client *alidns.Client) (err error) {
		addDomainRecordRequest := &alidns.AddDomainRecordRequest{
		DomainName: tea.String(dnconfig.DomainName),
		RR:			tea.String(dnconfig.RR),
		Type: 		tea.String(dnconfig.Type),
		Value: 		tea.String(dnconfig.Value),	
				
	}
	_, err = client.AddDomainRecord(addDomainRecordRequest)
	if err != nil {
		return err
	}
	return err
}