package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"waft/config"
)

//Proxy core data structure, mainly used to abstract a proxy server
type Proxy struct {
	conf     *config.ProxyConf
	reverser *httputil.ReverseProxy
	matcher  Matcher
	services map[string]*Service
	//filter Life cycle management
	filterFactory *filterFactory
}

//init matcher
func initMatchers(conf *config.ProxyConf) Matcher {
	matcher := NewMatcher()

	for srvName, srvConf := range conf.HTTP.Services {
		if err := matcher.Add(srvName, srvConf.Rule); err != nil {
			log.Fatalln("matcher add error", err.Error())
		}
	}

	return matcher
}

//init service
func initServices(proxy *Proxy, conf *config.ProxyConf) map[string]*Service {
	services := make(map[string]*Service)
	for srvName, srvConf := range conf.HTTP.Services {
		services[srvName] = NewService(proxy, &srvConf)
	}
	return services
}

//NewProxy create proxy server
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
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
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
	errFilter, _, err := service.doPreFilters(rw, req)
	//rw.WriteHeader(code)
	if err != nil {
		p.errorHandler(rw, req, fmt.Errorf("error filter name:%s, error: %v", errFilter, err))
		return
	}

	p.reverser.ServeHTTP(rw, req)
}

//Run proxy start, it is always running and will not exit
func (p *Proxy) Run() error {
	p.reverser.Director = p.director
	p.reverser.ModifyResponse = p.modifyResp
	p.reverser.ErrorHandler = p.errorHandler
	log.Println(fmt.Sprintf("Proxy Server Started, Serve Port %s ...", p.conf.HTTP.Listen))
	return http.ListenAndServe(p.conf.HTTP.Listen, p)
}
