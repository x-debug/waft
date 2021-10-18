package internal

import (
	"errors"
	"net/http"
)

const filterNameIPWhite = "ipWhiteList"

//ErrSourceRange IP White List source range config error
var ErrSourceRange = errors.New("source range parameter must pass")

//IPWhiteListFilter Restriction can only be released within these ip ranges
type IPWhiteListFilter struct {
	baseFilter
	sourceRange []string
}

//Init init filter
func (f *IPWhiteListFilter) Init(setting map[string]interface{}) error {
	if value, ok := setting["sourceRange"]; ok {
		ranges := make([]string, 0)
		for _, val := range value.([]interface{}) {
			ranges = append(ranges, val.(string))
		}
		f.sourceRange = ranges
	} else {
		return ErrSourceRange
	}
	return nil
}

//Name filter name
func (f *IPWhiteListFilter) Name() string {
	return filterNameIPWhite
}

//Pre front filter, execute before the request, usually used to modify the request object
func (f *IPWhiteListFilter) Pre(rw http.ResponseWriter, req *http.Request) (statusCode int, err error) {
	return 0, nil
}
