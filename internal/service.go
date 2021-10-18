package internal

import (
	"container/list"
	"waft/config"
	"waft/pkg/lb"
)

//Service every domain can abstract a service,
//It is a logical concept that can be used to configure domain name rules and specify filters,
//One proxy can correspond to multiple services
type Service struct {
	proxy *Proxy
	conf  *config.ServiceConf
	lb    lb.Balancer
	//list of Filter
	filters *list.List
	//list of *BackendInfo
	servers *list.List
}

//NewService create new service
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
