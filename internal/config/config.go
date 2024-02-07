package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	NameServer     string `yaml:"nameserver" envconfig:"NAMESERVER"`
    ApiKeyFile     string `envconfig:"API_KEY_FILE"`
	ApiSecretFile  string `envconfig:"API_SECRET_FILE"`
    Record struct {
        Domain string `yaml:"domain" envconfig:"DOMAIN"`
        Subdomain string `yaml:"subdomain" envconfig:"SUBDOMAIN"`
        TTL string `yaml:"ttl" envconfig:"TTL"`
        Id  string `yaml:"subdomainid" envconfig:"SUBDOMAIN_ID"`
    } `yaml:"record"`
	ApiCredentials struct {
		ApiKey string `yaml:"apikey" envconfig:"API_KEY"`
		Secret string `yaml:"apisecret" envconfig:"API_SECRET"`
	} `yaml:"apicredentials"`
}

func LoadConfig(path string) *Config {
	var cfg Config

	if path != "" {
		readFile(path, &cfg)
	}

	readEnv(&cfg)

	if cfg.ApiKeyFile != "" || cfg.ApiSecretFile != "" {
		loadCredentialsFromFile(&cfg)
	}

	return &cfg
}

func processError(err error) {
	log.Println(err)
	os.Exit(2)
}

func readFile(path string, cfg *Config) {
	f, err := os.Open(path)
	if err != nil {
		processError(err)
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

func loadCredentialsFromFile(cfg *Config) {
	if data, err := os.ReadFile(cfg.ApiKeyFile); err == nil {
		cfg.ApiCredentials.ApiKey = string(data)
	} else {
		log.Printf("Unable to open API_KEY_FILE %v", err)
	}

	if data, err := os.ReadFile(cfg.ApiSecretFile); err == nil {
		cfg.ApiCredentials.Secret = string(data)
	} else {
		log.Printf("Unable to open API_SECRET_FILE %v", err)
	}
}
