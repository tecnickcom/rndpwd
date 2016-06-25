package main

import (
	"time"
)

// JSend status codes
const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

// Response data format for HTTP
type Response struct {
	Service string      `json:"service"` // program name
	Version string      `json:"version"` // program version
	Time    time.Time   `json:"time"`    // timestamp
	Status  string      `json:"status"`  // status code (error|fail|success)
	Code    int         `json:"code"`    // HTTP status code
	Message string      `json:"message"` // error or status message
	Data    interface{} `json:"data"`    // data payload
}

// convert the HTTP status code into JSend status
func getStatus(code int) string {
	if code >= 500 {
		return StatusError
	}
	if code >= 400 {
		return StatusFail
	}
	return StatusSuccess
}
