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
}

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

func (ldnconfig *ListDomainConfig) ListDomain(client *alidns.Client) (result *alidns.DescribeDomainRecordsResponse, err error){
	describeDomainRecordsRequest := &alidns.DescribeDomainRecordsRequest{
		DomainName: 	tea.String(ldnconfig.DomainName),
		RRKeyWord: 		tea.String(ldnconfig.RRKeyWord),
		TypeKeyWord: 	tea.String(ldnconfig.TypeKeyWord),
		ValueKeyWord: 	tea.String(ldnconfig.ValueKeyWord),
		KeyWord:		tea.String(ldnconfig.KeyWord),
	}
	result, err = client.DescribeDomainRecords(describeDomainRecordsRequest)
	return result, err
}