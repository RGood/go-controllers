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
			Handle(http.MethodGet, func(req *http.Request) (int, interface{}) {
				return 200, req.Method + " " + req.URL.String()
			}).
			Create(),
	)

	http.ListenAndServe(":8080", nil)
}
