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

type ListDomainConfig struct {
	DomainName		string
	KeyWord 		string
	RRKeyWord		string
	TypeKeyWord 	string
	ValueKeyWord 	string
	PageSize		string
}

// 初始化一个alidns客户端
func (account *Account) CreateClient() (result *alidns.Client, err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: 		&account.AccessKey,
		// 您的AccessKey Secret
		AccessKeySecret: 	&account.AccessSecret,
		RegionId: 		 	&account.Region,
	}
	result, err = alidns.NewClient(config)
	return result, err
}

// 添加一个新的主机解析记录
func (dnconfig *DomainConfig) AddDomainRecord(client *alidns.Client) (err error) {
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

// 列出域名下所有主机解析记录
func (ldrconfig *ListDomainConfig) ListDomainRecords(client *alidns.Client) (result *alidns.DescribeDomainRecordsResponse, err error){
	describeDomainRecordsRequest := &alidns.DescribeDomainRecordsRequest{
		DomainName: 	tea.String(ldrconfig.DomainName),
		RRKeyWord: 		tea.String(ldrconfig.RRKeyWord),
		TypeKeyWord: 	tea.String(ldrconfig.TypeKeyWord),
		ValueKeyWord: 	tea.String(ldrconfig.ValueKeyWord),
		KeyWord:		tea.String(ldrconfig.KeyWord),
		PageSize:		tea.Int64(500),
	}
	result, err = client.DescribeDomainRecords(describeDomainRecordsRequest)
	return result, err
}

// 列出账户下所有域名
func (ldconfig *ListDomainConfig) ListDomains(client *alidns.Client) (result *alidns.DescribeDomainsResponse, err error) {
	describeDomainsRequest := &alidns.DescribeDomainsRequest{
		KeyWord: tea.String(ldconfig.KeyWord),
		PageSize: tea.Int64(100),
	}
	result, err = client.DescribeDomains(describeDomainsRequest)
	return result, err
}

// 删除域名主机记录对应的解析记录
func (dsdrconfig *DomainConfig) DelSubDomainRecords(client *alidns.Client) (result *alidns.DeleteSubDomainRecordsResponse, err error) {
	delSubDomainRecordsRequest := &alidns.DeleteSubDomainRecordsRequest{
		DomainName: tea.String(dsdrconfig.DomainName),
		RR: 		tea.String(dsdrconfig.RR),
		Type: 		tea.String(dsdrconfig.Type),
	}
	result, err = client.DeleteSubDomainRecords(delSubDomainRecordsRequest)
	return result, err
}