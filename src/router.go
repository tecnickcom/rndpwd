package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// start the HTTP server
func startServer(address string) {
	log.Printf("starting %s %s http server", AppName, AppVersion)
	router := httprouter.New()

	// set error handlers
	router.NotFound = http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
		// 404
		sendResponse(hw, hr, nil, http.StatusNotFound, "invalid end point")
	})
	router.MethodNotAllowed = http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
		// 405
		sendResponse(hw, hr, nil, http.StatusMethodNotAllowed, "the request cannot be routed")
	})
	router.PanicHandler = func(hw http.ResponseWriter, hr *http.Request, p interface{}) {
		// 500
		sendResponse(hw, hr, nil, http.StatusInternalServerError, "internal error")
	}

	// index handler
	router.GET("/", index)

	// set end points and handlers
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handle)
	}

	log.Printf("http server listening at '%s'", address)

	// start the http server
	log.Fatal(http.ListenAndServe(address, router))
}

// send the HTTP response in JSON format
func sendResponse(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params, code int, data interface{}) {
	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Expires", "0")
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	response := Response{
		Program: AppName,
		Version: AppVersion,
		Time:    time.Now().UTC(),
		Status:  getStatus(code),
		Code:    code,
		Message: http.StatusText(code),
		Data:    data,
	}

	// log request
	log.Printf("%s\t%s\t%d", hr.Method, hr.RequestURI, code)

	// send response as JSON
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println(err)
	}
}
