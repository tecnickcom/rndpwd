package main

import (
	"log"
	"os"
)

// Log levels
const (
	NONE int = iota - 1
	EMERGENCY
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

// logLevelCodes list log level codes by name (revese map of logLevelNames)
var logLevelCodes = map[string]int{
	"NONE":      NONE,      // Disable log
	"EMERGENCY": EMERGENCY, // System is unusable
	"ALERT":     ALERT,     // Should be corrected immediately
	"CRITICAL":  CRITICAL,  // Critical conditions
	"ERROR":     ERROR,     // Error conditions
	"WARNING":   WARNING,   // May indicate that an error will occur if action is not taken
	"NOTICE":    NOTICE,    // Events that are unusual, but not error conditions
	"INFO":      INFO,      // Normal operational messages that require no action
	"DEBUG":     DEBUG,     // Information useful to developers for debugging the application
}

// logLevelNames list log level names by code (revese map of logLevelCodes)
var logLevelNames = map[int]string{
	NONE:      "NONE",      // Disable log
	EMERGENCY: "EMERGENCY", // System is unusable
	ALERT:     "ALERT",     // Should be corrected immediately
	CRITICAL:  "CRITICAL",  // Critical conditions
	ERROR:     "ERROR",     // Error conditions
	WARNING:   "WARNING",   // May indicate that an error will occur if action is not taken
	NOTICE:    "NOTICE",    // Events that are unusual, but not error conditions
	INFO:      "INFO",      // Normal operational messages that require no action
	DEBUG:     "DEBUG",     // Information useful to developers for debugging the application
}

// logLevel is the current reporting log level
var logLevelCode = INFO

// logOutput set the log output
var logOutput = os.Stderr

// Log using the the specified level
func Log(level int, format string, v ...interface{}) {
	levelName, ok := logLevelNames[level]
	if !ok || (level > logLevelCode) {
		return
	}
	prefix := "[" + levelName + "] [" + ServiceName + "] "
	logger := log.New(logOutput, prefix, log.LstdFlags|log.LUTC)
	if level > CRITICAL {
		logger.Printf(format, v...)
	} else {
		logger.Fatalf(format, v...)
	}
}
