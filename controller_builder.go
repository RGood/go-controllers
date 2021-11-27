package controller

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ControllerBuilder struct {
	controllerFuncs map[string]func(*http.Request) Response
	defaultFunc     func(*http.Request) Response
	headers         map[string]string
}

type Response struct {
	code int
	body interface{}
}

func NewControllerbuilder() *ControllerBuilder {
	return &ControllerBuilder{
		controllerFuncs: map[string]func(*http.Request) Response{},
		defaultFunc: func(r *http.Request) Response {
			return Response{
				code: 404,
				body: struct{ Error string }{
					Error: "Method not implemented",
				},
			}
		},
	}
}

func (cb *ControllerBuilder) Handle(method string, f func(*http.Request) Response) *ControllerBuilder {
	cb.controllerFuncs[method] = f
	return cb
}

func (cb *ControllerBuilder) Default(f func(*http.Request) Response) *ControllerBuilder {
	cb.defaultFunc = f
	return cb
}

func (cb *ControllerBuilder) SetHeaders(h map[string]string) *ControllerBuilder {
	cb.headers = h
	return cb
}

func (cb *ControllerBuilder) getOptions() []string {
	methods := []string{}

	for k, _ := range cb.controllerFuncs {
		methods = append(methods, k)
	}

	return methods
}

func writeResponse(w http.ResponseWriter, res Response) {
	b, err := json.Marshal(res.body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error parsing response body"))
	} else {
		w.WriteHeader(res.code)
		w.Write(b)
	}

}

func (cb *ControllerBuilder) Create() func(http.ResponseWriter, *http.Request) {
	options := cb.getOptions()
	options = append(options, http.MethodOptions)
	optionsText := strings.Join(options, ", ")

	return func(w http.ResponseWriter, req *http.Request) {
		for k, v := range cb.headers {
			w.Header().Set(k, v)
		}

		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", optionsText)
			return
		}

		controllerFunc, ok := cb.controllerFuncs[req.Method]
		if ok {
			writeResponse(w, controllerFunc(req))
		} else {
			writeResponse(w, cb.defaultFunc(req))
		}
	}
}
