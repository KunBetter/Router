package Router

import (
	"net/http"
)

type Router struct {
	Paths map[string]*PathTrie
}

func GlobalRouter() *Router {
	return &Router{
		Paths: make(map[string]*PathTrie),
	}
}

//match the correct router.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	match, p := r.Paths[req.Method].MatchPath(req.URL.Path)
	if match {
		p.Process(w, req)
	}
}

func (r *Router) Insert(method, path string, handler HandlerFunc) {
	_, exist := r.Paths[method]
	if !exist {
		r.Paths[method] = NewPathTrie()
	}
	r.Paths[method].AddPath(path, handler)
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
