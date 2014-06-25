package main

import (
	"github.com/KunBetter/Router"
	"log"
	"net/http"
)

func ids(s *Router.Params) ([]byte, int) {
	return []byte(s.Value), 200
}

func main() {
	router := Router.GlobalRouter()
	router.Get("/ids/:id", Router.HandlerFunc(ids))
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
