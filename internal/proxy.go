package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"waft/config"
)

type Proxy struct {
	conf          *config.ProxyConf
	reverser      *httputil.ReverseProxy
	matcher       Matcher
	services      map[string]*Service
	//filter Life cycle management
	filterFactory *filterFactory
}

//init matcher
func initMatchers(conf *config.ProxyConf) Matcher {
	matcher := NewMatcher()

	for srvName, srvConf := range conf.Http.Services {
		if err := matcher.Add(srvName, srvConf.Rule); err != nil {
			log.Fatalln("matcher add error", err.Error())
		}
	}

	return matcher
}

//init service
func initServices(proxy *Proxy, conf *config.ProxyConf) map[string]*Service {
	services := make(map[string]*Service)
	for srvName, srvConf := range conf.Http.Services {
		services[srvName] = NewService(proxy, &srvConf)
	}
	return services
}

func NewProxy(conf *config.ProxyConf) *Proxy {
	proxy := &Proxy{
		conf:          conf,
		reverser:      &httputil.ReverseProxy{},
		matcher:       initMatchers(conf),
		filterFactory: createFilterFactory(conf),
	}
	proxy.services = initServices(proxy, conf)
	return proxy
}

func (p *Proxy) getService(req *http.Request) *Service {
	srvKey, err := p.matcher.Match(req)
	if err != nil {
		if err == ErrNotFound {
			//TODO service not found, return default service?
		}
		log.Fatalln("matcher error", err.Error())
	}

	service := p.services[srvKey]
	return service
}

func (p *Proxy) director(req *http.Request) {
	service := p.getService(req)
	backend := service.lb.Select(req, service.servers)
	req.URL.Scheme = backend.Parsed.Scheme
	req.URL.Host = backend.Parsed.Host
}

func (p *Proxy) modifyResp(resp *http.Response) error {
	service := p.getService(resp.Request)
	errFilter, code, err := service.doPostFilters(resp)
	resp.StatusCode = code
	if err != nil {
		return fmt.Errorf("error filter name:%s, error: %v", errFilter, err)
	}
	return err
}

func (p *Proxy) errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	service := p.getService(req)
	service.doPostErrFilters(rw, req, err)
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	service := p.getService(req)
	errFilter, code, err := service.doPreFilters(rw, req)
	rw.WriteHeader(code)
	if err != nil {
		p.errorHandler(rw, req, fmt.Errorf("error filter name:%s, error: %v", errFilter, err))
		return
	}

	p.reverser.ServeHTTP(rw, req)
}

func (p *Proxy) Run() error {
	p.reverser.Director = p.director
	p.reverser.ModifyResponse = p.modifyResp
	p.reverser.ErrorHandler = p.errorHandler
	log.Println(fmt.Sprintf("Proxy Server Started, Serve Port %s ...", p.conf.Http.Listen))
	return http.ListenAndServe(p.conf.Http.Listen, p)
}
