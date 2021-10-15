package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/url"
)

type ProxyConf struct {
	Pid  string `yaml:"pid"`
	Http struct {
		Listen   string                 `yaml:"listen"`
		Services map[string]ServiceConf `yaml:"services"`
		//first key is group name, second key is filter name, third key is setting column
		Filters map[string]map[string]map[string]interface{} `yaml:"filters"`
	} `yaml:"http"`
}

type ServiceConf struct {
	proxy    *ProxyConf
	Mode     string         `yaml:"mode"` //static mode: define value in yaml file, etcd mode: read value from etcd central repo
	Rule     string         `yaml:"rule"` //path selector
	Balancer string         `yaml:"balancer"`
	Servers  []*BackendInfo `yaml:"servers"`
	Filter   string         `yaml:"filter"`
}

func LoadProxyConf(conf io.Reader) (*ProxyConf, error) {
	var proxyConf ProxyConf

	confBuf, err := ioutil.ReadAll(conf)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(confBuf, &proxyConf)
	if err != nil {
		return nil, err
	}
	for _, srv := range proxyConf.Http.Services {
		srv.proxy = &proxyConf
		for _, backend := range srv.Servers {
			if err = backend.parse(); err != nil {
				log.Fatalln("parse server url error: ", backend.Url)
			}
		}
	}
	return &proxyConf, nil
}

type BackendInfo struct {
	Url    string `yaml:"url"`
	Weight int    `yaml:"weight"`

	Parsed struct {
		Scheme string
		Host   string
	}
}

func (bi *BackendInfo) parse() error {
	uri, err := url.Parse(bi.Url)
	if err != nil {
		return err
	}

	bi.Parsed.Scheme = uri.Scheme
	bi.Parsed.Host = uri.Host
	return nil
}
