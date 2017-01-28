package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var startTime = time.Now()

// return a list of available routes
func indexHandler(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	stats.Increment("http.index.in")
	type configInfo struct {
		Charset  string `json:"charset"`  // characters to use to generate a password
		Length   int    `json:"length"`   // length of each password (number of characters or bytes)
		Quantity int    `json:"quantity"` // number of passwords to generate
	}
	type info struct {
		Duration float64    `json:"duration"` // elapsed time since last passwor drequest or service start
		Entries  Routes     `json:"routes"`   // available routes (http entry points)
		Config   configInfo `json:"config"`   // configuration parameters
	}
	sendResponse(rw, hr, ps, http.StatusOK, info{
		Duration: time.Since(startTime).Seconds(),
		Entries:  routes,
		Config:   configInfo{Charset: appParams.charset, Length: appParams.length, Quantity: appParams.quantity},
	})
	stats.Increment("http.index.out")
}

// returns the status of the service
func statusHandler(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	stats.Increment("http.status.in")
	type info struct {
		Duration float64 `json:"duration"` // elapsed time since last request in seconds
		Message  string  `json:"message"`  // error message
	}
	status := http.StatusOK
	message := "The service is healthy"
	sendResponse(rw, hr, ps, status, info{
		Duration: time.Since(startTime).Seconds(),
		Message:  message,
	})
	stats.Increment("http.status.out")
}

// password returns the requested password
func password(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	stats.Increment("http.password.in")
	startTime = time.Now()

	// decode URL parameters (if any)
	urlParams := appParams
	query := hr.URL.Query()
	if val, exist := query["charset"]; exist {
		urlParams.charset = val[0]
	}
	if val, exist := query["length"]; exist {
		i, err := strconv.Atoi(val[0])
		if err == nil {
			urlParams.length = i
		}
	}
	if val, exist := query["quantity"]; exist {
		i, err := strconv.Atoi(val[0])
		if err == nil {
			urlParams.quantity = i
		}
	}
	err := checkParams(urlParams)
	if err != nil {
		sendResponse(rw, hr, ps, http.StatusBadRequest, err.Error())
		return
	}

	type info struct {
		Passwords []string `json:"passwords"` // number of executed tests
		Duration  float64  `json:"duration"`  // password generation time
	}
	sendResponse(rw, hr, ps, http.StatusOK, info{
		Passwords: getAllPassword(urlParams),
		Duration:  time.Since(startTime).Seconds(),
	})
	stats.Increment("http.password.out")
}
