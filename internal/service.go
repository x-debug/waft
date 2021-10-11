package internal

import (
	"container/list"
	"waft/config"
	"waft/pkg/lb"
)

type Service struct {
	proxy *Proxy
	conf  *config.ServiceConf
	lb    lb.Balancer
	//list of Filter
	filters *list.List
	//list of *BackendInfo
	servers *list.List
}

func NewService(proxy *Proxy, conf *config.ServiceConf) *Service {
	service := &Service{proxy: proxy, conf: conf}
	service.lb = lb.CreateBalancer(conf.Balancer, conf)
	service.filters = proxy.filterFactory.find(conf.Filter)
	service.servers = list.New()
	for _, srvConf := range conf.Servers {
		service.servers.PushBack(srvConf)
	}
	return service
}
