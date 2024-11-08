package util_http

import (
	"net/http"
	"slices"
)

type (
	middleware func(http.Handler) http.Handler
	router     struct {
		*http.ServeMux
		chain []middleware
	}
)

func NewRouter(mx ...middleware) *router {
	return &router{ServeMux: &http.ServeMux{}, chain: mx}
}

func (r *router) Use(mx ...middleware) {
	r.chain = append(r.chain, mx...)
}

func (r *router) Group(fn func(r *router)) {
	fn(&router{ServeMux: r.ServeMux, chain: slices.Clone(r.chain)})
}

func (r *router) Get(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodGet, path, fn, mx)
}

func (r *router) Post(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodPost, path, fn, mx)
}

func (r *router) Put(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodPut, path, fn, mx)
}

func (r *router) Delete(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodDelete, path, fn, mx)
}

func (r *router) Head(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodHead, path, fn, mx)
}

func (r *router) Options(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodOptions, path, fn, mx)
}

func (r *router) handle(method, path string, fn http.HandlerFunc, mx []middleware) {
	r.Handle(method+" "+path, r.wrap(fn, mx))
}

func (r *router) wrap(fn http.HandlerFunc, mx []middleware) (out http.Handler) {
	out, mx = http.Handler(fn), append(slices.Clone(r.chain), mx...)

	slices.Reverse(mx)

	for _, m := range mx {
		out = m(out)
	}

	return
}
