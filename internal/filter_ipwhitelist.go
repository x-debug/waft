package internal

import (
	"errors"
	"net/http"
)

const filterNameIpWhite = "ipWhiteList"

var ErrSourceRange = errors.New("source range parameter must pass")

type IpWhiteListFilter struct {
	baseFilter
	sourceRange []string
}

func (f *IpWhiteListFilter) Init(setting map[string]interface{}) error {
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

func (f *IpWhiteListFilter) Name() string {
	return filterNameIpWhite
}

func (f *IpWhiteListFilter) Pre(rw http.ResponseWriter, req *http.Request) (statusCode int, err error) {
	return 0, nil
}
