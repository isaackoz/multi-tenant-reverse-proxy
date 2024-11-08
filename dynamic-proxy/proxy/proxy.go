package proxy

import (
	"dynamic-proxy/config"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func New(cfg *config.Config) *ReverseProxy {
	// static for now, must be updated later
	target, _ := url.Parse("http://example.com")

	return &ReverseProxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
}

func (rp *ReverseProxy) Handle(w http.ResponseWriter, r *http.Request) {
	r.Host = rp.target.Host
	rp.proxy.ServeHTTP(w, r)
}
