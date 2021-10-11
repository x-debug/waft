package lb

import (
	"container/list"
	"math/rand"
	"net/http"
	"waft/config"
)

type randomBalancer struct {
	nopBalancer
}

func (r *randomBalancer) Select(req *http.Request, servers *list.List) *config.BackendInfo {
	serversArr := make([]*config.BackendInfo, 0)

	for srv := servers.Front(); srv != nil; srv = srv.Next() {
		serversArr = append(serversArr, srv.Value.(*config.BackendInfo))
	}

	return serversArr[rand.Intn(len(serversArr))]
}

func (r *randomBalancer) Name() string {
	return "random"
}
