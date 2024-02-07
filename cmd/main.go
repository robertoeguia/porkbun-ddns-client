package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/robertoeguia/porkbun-ddns-client/internal/config"
	"github.com/robertoeguia/porkbun-ddns-client/internal/dnsutil"
)

const (
    ipv4hostname = "https://api-ipv4.porkbun.com/api/json/v3"
    apiHostname = "https://porkbun.com/api/json/v3"
)

type porkbunUpdateRecord struct {
	SecretAPIKey	string	`json:"secretapikey"`
	APIKey			string	`json:"apikey"`
	Name			string	`json:"name,omitempty"`
	TTL				string	`json:"ttl,omitempty"`
	Type			string	`json:"type,omitempty"`
	Content			string	`json:"content,omitempty"`
}

type porkbunDDNSClient struct {
    cfg *config.Config
    restClient *resty.Client
}

func main() {
    ddnsClient := &porkbunDDNSClient{}

    cfg := config.LoadConfig("./config/config.yaml")
    client := resty.New()

    ddnsClient.restClient = resty.New()
    ddnsClient.cfg = cfg

    dnsutil.SetNameserver(cfg.NameServer)

    porkbunEndpoint := fmt.Sprintf("%v/dns/edit/%v/%v", apiHostname , cfg.Record.Domain, cfg.Record.Id)
    data := &porkbunUpdateRecord{
        SecretAPIKey: cfg.ApiCredentials.Secret,
        APIKey: cfg.ApiCredentials.ApiKey,
        Name: cfg.Record.Subdomain,
        TTL: cfg.Record.TTL,
        Type: "A",
        Content: "n/a",
    }

    record := cfg.Record.Domain
    if cfg.Record.Subdomain != "" {
        record = fmt.Sprintf("%v.%v",cfg.Record.Subdomain,cfg.Record.Domain)
    }

    log.Printf("Getting DNS record: %v", record)
    dnsRecord,err := dnsutil.GetARecord(record)
    if err != nil {
        log.Panic("Error getting dns record",err)
    }
    log.Printf("DNS Record: %v", dnsRecord)


    log.Printf("Getting Public IP address")
    body := &porkbunUpdateRecord {
        SecretAPIKey: cfg.ApiCredentials.Secret,
        APIKey: cfg.ApiCredentials.ApiKey,
    }

    requestBody,_ := json.Marshal(body)
    response,_ := client.NewRequest().EnableTrace().
                    SetBody(requestBody).Post(fmt.Sprintf("%v/ping",ipv4hostname))
    
    var result map[string]string
    json.Unmarshal(response.Body(),&result)
    log.Printf("Public IP: %v",result["yourIp"])


    if dnsRecord == result["yourIp"] {
        log.Println("IP addresses match. Doing nothing")
        return;
    }

    log.Println("Updating DNS record")
    data.Content = result["yourIp"]

    response,_ = client.NewRequest().EnableTrace().
                    SetBody(data).Post(porkbunEndpoint)

    if response.StatusCode() == 400 {
        log.Fatalln("Error updating DNS record", response)
    } else {
        log.Printf("DNS record update successful: %v", response)
    }
}