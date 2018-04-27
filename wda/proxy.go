package wda

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Proxy struct {
	rp    *httputil.ReverseProxy
	mutex sync.Mutex
}

func NewProxy(target *url.URL) *Proxy {
	return &Proxy{
		rp: httputil.NewSingleHostReverseProxy(target),
	}
}

func (w *Proxy) SetUrl(u string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.rp = httputil.NewSingleHostReverseProxy(&url.URL{Host: u, Scheme: "http"})
}

func (w *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w.rp.ServeHTTP(rw, req)
}
