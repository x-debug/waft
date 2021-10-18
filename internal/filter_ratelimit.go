package internal

import "net/http"

const filterNameRateLimit = "rateLimit"

//RateLimitFilter rate limit filter
type RateLimitFilter struct {
	baseFilter
	average int
	burst   int
}

//Init init filter
func (f *RateLimitFilter) Init(setting map[string]interface{}) error {
	f.average = 100
	f.burst = 50

	if value, ok := setting["average"]; ok {
		f.average = value.(int)
	}

	if value, ok := setting["burst"]; ok {
		f.burst = value.(int)
	}
	return nil
}

//Name filter name
func (f *RateLimitFilter) Name() string {
	return filterNameRateLimit
}

//Pre front filter, execute before the request, usually used to modify the request object
func (f *RateLimitFilter) Pre(rw http.ResponseWriter, req *http.Request) (statusCode int, err error) {
	return 0, nil
}
