package lb

import (
	"testing"
	"waft/config"
)

func selectServer(t *testing.T, wrr *WRR) {
	server := wrr.Select(nil, nil)
	if server == nil {
		t.Error("WRR select error")
	}
	t.Logf("selected server is %s", server.Url)
}

func TestWRR_Select(t *testing.T) {
	wrr := &WRR{}

	services := make([]*config.BackendInfo, 0)
	services = append(services, &config.BackendInfo{Url: "A", Weight: 4})
	services = append(services, &config.BackendInfo{Url: "B", Weight: 3})
	services = append(services, &config.BackendInfo{Url: "C", Weight: 2})
	srvConf := &config.ServiceConf{Servers: services}
	wrr.Init(srvConf)
	for i := 0; i < 9; i++ {
		selectServer(t, wrr)
	}
}
