package controller

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ControllerBuilder struct {
	controllerFuncs map[string]func(*http.Request) (int, interface{})
	defaultFunc     func(*http.Request) (int, interface{})
	beforeFunc      func(*http.Request, func(*http.Request) (int, interface{})) (int, interface{})
	headers         map[string]string
}

func NewControllerbuilder() *ControllerBuilder {
	return &ControllerBuilder{
		controllerFuncs: map[string]func(*http.Request) (int, interface{}){},
		defaultFunc: func(r *http.Request) (int, interface{}) {
			return 404, struct{ Error string }{
				Error: "Method not implemented",
			}
		},
	}
}

func (cb *ControllerBuilder) Before(beforeFunc func(*http.Request, func(*http.Request) (int, interface{})) (int, interface{})) *ControllerBuilder {
	cb.beforeFunc = beforeFunc
	return cb
}

func (cb *ControllerBuilder) Handle(method string, f func(*http.Request) (int, interface{})) *ControllerBuilder {
	cb.controllerFuncs[method] = f
	return cb
}

func (cb *ControllerBuilder) Default(f func(*http.Request) (int, interface{})) *ControllerBuilder {
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

func writeResponse(w http.ResponseWriter, code int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error parsing response body"))
	} else {
		w.WriteHeader(code)
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
			if cb.beforeFunc != nil {
				code, res := cb.beforeFunc(req, controllerFunc)
				writeResponse(w, code, res)
			} else {
				code, res := controllerFunc(req)
				writeResponse(w, code, res)
			}
		} else {
			code, res := cb.defaultFunc(req)
			writeResponse(w, code, res)
		}
	}
}
