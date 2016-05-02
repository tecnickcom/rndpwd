package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// return a list of available routes
func index(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	sendResponse(rw, hr, ps, http.StatusOK, routes)
}

// check if the service is alive
func ping(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {
	sendResponse(rw, hr, ps, http.StatusOK, "ping")
}

// check if the service is alive
func password(rw http.ResponseWriter, hr *http.Request, ps httprouter.Params) {

	// decode URl parameters (if any)
	urlParams := appParams
	query := hr.URL.Query()
	if val, exist := query["quantity"]; exist {
		i, err := strconv.Atoi(val[0])
		if err == nil {
			urlParams.quantity = i
		}
	}
	if val, exist := query["length"]; exist {
		i, err := strconv.Atoi(val[0])
		if err == nil {
			urlParams.length = i
		}
	}
	if val, exist := query["charset"]; exist {
		urlParams.charset = val[0]
	}
	err := checkParams(urlParams)
	if err != nil {
		sendResponse(rw, hr, ps, http.StatusBadRequest, err)
		return
	}

	sendResponse(rw, hr, ps, http.StatusOK, getAllPassword(urlParams))
}
