package main

import (
	"log"
	"os"
)

// syslog error levels
var (
	// 0 - Emergency: System is unusable
	emergLog = log.New(os.Stderr, "[EMERG] [RNDPWD] ", log.LstdFlags|log.Lshortfile|log.LUTC)

	// 1 - Alert: Should be corrected immediately
	alertLog = log.New(os.Stderr, "[ALERT] [RNDPWD] ", log.LstdFlags|log.Lshortfile|log.LUTC)

	// 2 - Critical: Critical conditions
	critLog = log.New(os.Stderr, "[CRIT] [RNDPWD] ", log.LstdFlags|log.Lshortfile|log.LUTC)

	// 3 - Error: Error conditions
	errLog = log.New(os.Stderr, "[ERR] [RNDPWD] ", log.LstdFlags|log.Lshortfile|log.LUTC)

	// 4 - Warning: May indicate that an error will occur if action is not taken.
	warningLog = log.New(os.Stderr, "[WARNING] [RNDPWD] ", log.LstdFlags|log.LUTC)

	// 5 - Notice: Events that are unusual, but not error conditions.
	noticeLog = log.New(os.Stderr, "[NOTICE] [RNDPWD] ", log.LstdFlags|log.LUTC)

	// 6 - Informational: Normal operational messages that require no action.
	infoLog = log.New(os.Stdout, "[INFO] [RNDPWD] ", log.LstdFlags|log.LUTC)

	// 7 - Debug: Information useful to developers for debugging the application.
	debugLog = log.New(os.Stderr, "[DEBUG] [RNDPWD] ", log.LstdFlags|log.Lshortfile|log.LUTC)
)
