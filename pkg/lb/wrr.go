package lb

import (
	"container/list"
	"log"
	"net/http"
	"waft/config"
)

//WRR algorithm from http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
//intro https://colobu.com/2016/12/04/smooth-weighted-round-robin-algorithm/ & https://github.com/smallnest/weighted/blob/master/roundrobin_weighted.go
type WRR struct {
	nopBalancer
	items []*config.BackendInfo
	n     int
	gcd   int
	maxW  int
	i     int
	cw    int
}

func (w *WRR) Init(conf *config.ServiceConf) {
	for _, srvConf := range conf.Servers {
		w.Add(srvConf)
	}
}

// All returns all items.
func (w *WRR) All() map[interface{}]int {
	m := make(map[interface{}]int)
	for _, i := range w.items {
		m[i.Url] = i.Weight
	}
	return m
}

// RemoveAll removes all weighted items.
func (w *WRR) RemoveAll() {
	w.items = w.items[:0]
	w.n = 0
	w.gcd = 0
	w.maxW = 0
	w.i = -1
	w.cw = 0
}

//Reset resets all current weights.
func (w *WRR) Reset() {
	w.i = -1
	w.cw = 0
}

func (w *WRR) Add(server *config.BackendInfo) {
	if server.Weight > 0 {
		if w.gcd == 0 {
			w.gcd = server.Weight
			w.maxW = server.Weight
			w.i = -1
			w.cw = 0
		} else {
			w.gcd = gcd(w.gcd, server.Weight)
			if w.maxW < server.Weight {
				w.maxW = server.Weight
			}
		}
	}
	w.items = append(w.items, server)
	w.n++
	log.Println("Add Server ", server.Url)
}

func gcd(x, y int) int {
	var t int
	for {
		t = (x % y)
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}

func (w *WRR) Select(req *http.Request, servers *list.List) *config.BackendInfo {
	if w.n == 0 {
		return nil
	}

	if w.n == 1 {
		return w.items[0]
	}

	for {
		w.i = (w.i + 1) % w.n
		if w.i == 0 {
			w.cw = w.cw - w.gcd
			if w.cw <= 0 {
				w.cw = w.maxW
				if w.cw == 0 {
					return nil
				}
			}
		}

		if w.items[w.i].Weight >= w.cw {
			log.Println("Selected Server is ", w.items[w.i].Url)
			return w.items[w.i]
		}
	}
}

func (w *WRR) Name() string {
	return "wrr"
}
