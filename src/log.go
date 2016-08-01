package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(new(logJSONFormatter))

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

type logJSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

// Format is a custom JSON formatter for the logs
func (f *logJSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = log.DefaultTimestampFormat
	}

	data["program"] = ServiceName
	data["version"] = ServiceVersion
	data["release"] = ServiceRelease
	data["time"] = entry.Time.Format(timestampFormat)
	data["timestamp"] = time.Now().UTC().UnixNano()
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

// prefixFieldClashes avoid overwrite default fields
func prefixFieldClashes(data log.Fields) {
	fields := [...]string{"program", "version", "release", "time", "timestamp", "msg", "level"}
	for i := range fields {
		if val, ok := data[fields[i]]; ok {
			data["fields."+fields[i]] = val
		}
	}
}
