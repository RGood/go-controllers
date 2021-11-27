package main

import (
	"net/http"

	"github.com/rgood/go-controllers/pkg/controller"
)

type Response struct {
	Message string
}

func main() {

	http.HandleFunc(
		"/",
		controller.NewControllerbuilder().
			Before(func(req *http.Request, next func(req *http.Request) (int, interface{})) (int, interface{}) {
				println("In Before Function: " + req.Method + " " + req.URL.String())
				return next(req)
			}).
			Handle(http.MethodGet, func(req *http.Request) (int, interface{}) {
				println("In Handler Function: " + req.Method + " " + req.URL.String())
				return 200, req.Method + " " + req.URL.String()
			}).
			Create(),
	)

	http.ListenAndServe(":8080", nil)
}
