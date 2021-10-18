package lb

import (
	"container/list"
	"net/http"
	"waft/config"
)

//Balancer Load balancing is used to select one machine among multiple machines
type Balancer interface {
	Init(conf *config.ServiceConf)
	Select(req *http.Request, services *list.List) *config.BackendInfo
	Name() string
}

type nopBalancer struct {
	conf *config.ServiceConf
}

func (r *nopBalancer) Init(conf *config.ServiceConf) {
	r.conf = conf
}

func (r *nopBalancer) Select(req *http.Request, servers *list.List) *config.BackendInfo {
	return nil
}

func (r *nopBalancer) Name() string {
	return "nop"
}

var balancerTable = map[string]Balancer{
	"random": &randomBalancer{},
	"wrr":    &WRR{},
}

//CreateBalancer create new balancer
func CreateBalancer(name string, conf *config.ServiceConf) Balancer {
	var (
		balancer Balancer
		ok       bool
	)

	if balancer, ok = balancerTable[name]; !ok {
		return nil
	}
	balancer.Init(conf)
	return balancer
}
