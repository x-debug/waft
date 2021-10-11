package internal

import (
	"net/http"
)

type filterContext struct {
	rw     http.ResponseWriter
	req    *http.Request
	outReq *http.Request
}

type Filter interface {
	Name() string

	Init(map[string]interface{}) error
	Pre(rw http.ResponseWriter, req *http.Request) (statusCode int, err error)
	Post(*http.Response) (statusCode int, err error)
	PostErr(http.ResponseWriter, *http.Request, error)
}

type baseFilter struct{}

func (filter *baseFilter) Init(setting map[string]interface{}) error {
	return nil
}

func (filter *baseFilter) Pre(rw http.ResponseWriter, req *http.Request) (statusCode int, err error) {
	return http.StatusOK, nil
}

func (filter *baseFilter) Post(resp *http.Response) (statusCode int, err error) {
	return http.StatusOK, nil
}

func (filter *baseFilter) PostErr(rw http.ResponseWriter, req *http.Request, err error) {

}

func (serv *Service) doPreFilters(rw http.ResponseWriter, req *http.Request) (filterName string, statusCode int, err error) {
	for iter := serv.filters.Front(); iter != nil; iter = iter.Next() {
		f, _ := iter.Value.(Filter)
		filterName = f.Name()

		statusCode, err = f.Pre(rw, req)
		if nil != err {
			return filterName, statusCode, err
		}
	}

	return "", http.StatusOK, nil
}

func (serv *Service) doPostFilters(resp *http.Response) (filterName string, statusCode int, err error) {
	for iter := serv.filters.Back(); iter != nil; iter = iter.Prev() {
		f, _ := iter.Value.(Filter)

		statusCode, err = f.Post(resp)
		if nil != err {
			return filterName, statusCode, err
		}
	}

	return "", http.StatusOK, nil
}

func (serv *Service) doPostErrFilters(rw http.ResponseWriter, req *http.Request, err error) {
	for iter := serv.filters.Back(); iter != nil; iter = iter.Prev() {
		f, _ := iter.Value.(Filter)

		f.PostErr(rw, req, err)
	}
}
