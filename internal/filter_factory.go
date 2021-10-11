package internal

import (
	"container/list"
	"fmt"
	"log"
	"waft/config"
)

type filterFactory struct {
	//list of []Filter
	filters map[string]*list.List
}

func initializeFilters(name string, setting map[string]interface{}) Filter {
	var filter Filter
	switch name {
	case filterNameIpWhite:
		filter = &IpWhiteListFilter{}
	case filterNameRateLimit:
		filter = &RateLimitFilter{}
	}

	if err := filter.Init(setting); err != nil {
		log.Fatalln(fmt.Sprintf("filter %s initialize error", name), err.Error())
	}
	return filter
}

func createFilterFactory(conf *config.ProxyConf) *filterFactory {
	filterMap := make(map[string]*list.List)
	for group, filters := range conf.Http.Filters {
		scopeFilters := list.New()
		for name, setting := range filters {
			scopeFilters.PushBack(initializeFilters(name, setting))
		}
		filterMap[group] = scopeFilters
	}
	return &filterFactory{filters: filterMap}
}

func (ff *filterFactory) find(groupName string) *list.List {
	return ff.filters[groupName]
}
