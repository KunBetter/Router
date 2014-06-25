package Router

import (
	"net/http"
)

type Router struct {
	Paths map[string][]*Matcher
}

func GlobalRouter() *Router {
	return &Router{
		make(map[string][]*Matcher),
	}
}

//match the correct router.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, method := range r.Paths[req.Method] {
		if match, Params := method.Matching(req.URL.Path); match {
			method.Process(Params, w, req)
		}
	}
}

func (r *Router) Insert(method, path string, handler HandlerFunc) {
	r.Paths[method] = append(r.Paths[method], &Matcher{path, handler})
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.Insert("HEAD", path, handler)
	r.Insert("GET", path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.Insert("POST", path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.Insert("PUT", path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.Insert("DELETE", path, handler)
}

func (r *Router) Head(path string, handler HandlerFunc) {
	r.Insert("HEAD", path, handler)
}

func (r *Router) Options(path string, handler HandlerFunc) {
	r.Insert("OPTIONS", path, handler)
}
