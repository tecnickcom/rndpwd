package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// start the HTTP server
func startServer(address string) error {
	Log(INFO, "starting %s %s http server", ServiceName, ServiceVersion)
	router := httprouter.New()

	// set error handlers
	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, hr *http.Request) { // 404
		sendResponse(rw, hr, nil, http.StatusNotFound, "invalid end point")
	})
	router.MethodNotAllowed = http.HandlerFunc(func(rw http.ResponseWriter, hr *http.Request) { // 405
		sendResponse(rw, hr, nil, http.StatusMethodNotAllowed, "the request cannot be routed")
	})
	router.PanicHandler = func(rw http.ResponseWriter, hr *http.Request, p interface{}) { // 500
		sendResponse(rw, hr, nil, http.StatusInternalServerError, "internal error")
	}

	// index handler
	router.GET("/", index)

	// set end points and handlers
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handle)
	}

	Log(INFO, "http server listening at '%s'", address)
	return fmt.Errorf("unable to start the HTTP server: %v", http.ListenAndServe(address, router))
}

// send the HTTP response in JSON format
func sendResponse(hw http.ResponseWriter, hr *http.Request, ps httprouter.Params, code int, data interface{}) {
	hw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	hw.Header().Set("Pragma", "no-cache")
	hw.Header().Set("Expires", "0")
	hw.Header().Set("Content-Type", "application/json")
	hw.WriteHeader(code)

	response := Response{
		Service: ServiceName,
		Version: ServiceVersion + "-" + ServiceRelease,
		Time:    time.Now().UTC(),
		Status:  getStatus(code),
		Code:    code,
		Message: http.StatusText(code),
		Data:    data,
	}

	// log request
	Log(INFO, "%s\t%s\t%d", hr.Method, hr.RequestURI, code)

	// send response as JSON
	if err := json.NewEncoder(hw).Encode(response); err != nil {
		Log(ERROR, "%v", err)
	}
}
