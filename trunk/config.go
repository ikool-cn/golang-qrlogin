package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	RedirectUrl             string `json:"redirect_url"`
	Listen                  string `json:"listen"`
	Timeout                 int    `json:"timeout"`
	CreatedOrganizationsUrl string `json:"create_organization_url"`
}

func NewConfigDefault() *Config {
	cfg := new(Config)
	cfg.RedirectUrl = "http://www.abc.com/"
	cfg.Listen = ":8090"
	cfg.Timeout = 30
	return cfg
}

func NewConfigFromFile(filename string) *Config {
	data, _ := ioutil.ReadFile(filename)
	return NewConfig(data)
}

func NewConfig(data []byte) *Config {
	cfg := NewConfigDefault()
	err := json.Unmarshal(data, cfg)
	if err != nil {
		log.Fatal("config is not a valid json")
	}
	return cfg
}
