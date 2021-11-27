package main

import (
	"net/http"

	"github.com/rgood/go-controllers/pkg/controller"
)

type Response struct {
	Message string
}

func main() {
	cb := controller.NewControllerbuilder()

	http.HandleFunc(
		"/",
		cb.Handle(http.MethodGet, func(req *http.Request) (int, interface{}) {
			return 200, "Hello world!"
		}).Create(),
	)

	http.ListenAndServe(":8080", nil)
}
