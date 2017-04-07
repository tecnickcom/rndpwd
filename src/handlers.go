package main

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

var startTime = time.Now()

// index returns a list of available routes
func indexHandler(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	stats.Increment("http.index.in")
	defer stats.Increment("http.index.out")
	log.Debug("Handler: index")
	type info struct {
		Duration float64 `json:"duration"` // elapsed time since service start [seconds]
		Entries  Routes  `json:"routes"`   // available routes (http entry points)
	}
	sendResponse(rw, hr, ps, http.StatusOK, info{
		Duration: time.Since(startTime).Seconds(),
		Entries:  routes,
	})
}

// statusHandler returns the status of the service
func statusHandler(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	stats.Increment("http.status.in")
	defer stats.Increment("http.status.out")
	log.Debug("Handler: status")
	type info struct {
		Duration float64 `json:"duration"` // elapsed time since service start [seconds]
		Message  string  `json:"message"`  // error message
	}
	status := http.StatusOK
	message := "The service is healthy"
	sendResponse(rw, hr, ps, status, info{
		Duration: time.Since(startTime).Seconds(),
		Message:  message,
	})
}

// password returns the requested password
func passwordHandler(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	stats.Increment("http.password.in")
	defer stats.Increment("http.password.out")
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
}
