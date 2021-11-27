package controller

import (
	"net/http"
	"strings"
)

type GetController interface {
	Get(http.ResponseWriter, *http.Request)
}

type PostController interface {
	Post(http.ResponseWriter, *http.Request)
}

type PutController interface {
	Put(http.ResponseWriter, *http.Request)
}

type DeleteController interface {
	Delete(http.ResponseWriter, *http.Request)
}

func notImplemented(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Method Not Implemented"))
}

func getControllerMethods(controller interface{}) []string {
	methods := []string{}
	_, ok := controller.(GetController)
	if ok {
		methods = append(methods, http.MethodGet)
	}

	_, ok = controller.(PutController)
	if ok {
		methods = append(methods, http.MethodPut)
	}

	_, ok = controller.(PostController)
	if ok {
		methods = append(methods, http.MethodPost)
	}

	_, ok = controller.(DeleteController)
	if ok {
		methods = append(methods, http.MethodDelete)
	}

	return methods
}

func NewDefaultControllerFunc(controller interface{}) func(http.ResponseWriter, *http.Request) {
	return NewControllerFunc(controller, notImplemented)
}

func NewControllerFunc(controller interface{}, defaultFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			methods := getControllerMethods(controller)
			methods = append(methods, "OPTIONS")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			return
		}

		switch req.Method {
		case "GET":
			controller, ok := controller.(GetController)
			if ok {
				controller.Get(w, req)
			} else {
				defaultFunc(w, req)
			}
		case "POST":
			controller, ok := controller.(PostController)
			if ok {
				controller.Post(w, req)
			} else {
				defaultFunc(w, req)
			}
		case "PUT":
			controller, ok := controller.(PutController)
			if ok {
				controller.Put(w, req)
			} else {
				defaultFunc(w, req)
			}
		case "DELETE":
			controller, ok := controller.(DeleteController)
			if ok {
				controller.Delete(w, req)
			} else {
				defaultFunc(w, req)
			}
		default:
			defaultFunc(w, req)
		}
	}
}
