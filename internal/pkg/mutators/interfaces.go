package mutators

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Mutator is a common interface for all mutators
type Mutator interface {
	Mutate(Log) (string, error)
}

// Log
type Log struct {
	AppName        string            `json:"app_name"`
	Client         string            `json:"client"`
	Facility       int               `json:"facility"`
	Hostname       string            `json:"hostname"`
	Message        string            `json:"message"`
	MsgID          string            `json:"msg_id"`
	Priority       int               `json:"priority"`
	ProcId         string            `json:"proc_id"`
	Severity       int               `json:"severity"`
	StructuredData map[string]string `json:"structured_data"`
	Timestamp      time.Time         `json:"timestamp"`
	TLSPeer        string            `json:"tls_peer"`
	Version        int               `json:"version"`
}

// NewLog creates a Log instance
func NewLog(logParts map[string]interface{}) Log {
	var sd map[string]string
	if logParts["structured_data"] == nil {
		sd = make(map[string]string)
	} else {
		sd = parseStructuredData(logParts["structured_data"].(string))
	}
	return Log{
		AppName:        getString(logParts["app_name"], "-"),
		Client:         getString(logParts["client"], "-"),
		Facility:       logParts["facility"].(int),
		Hostname:       getString(logParts["hostname"], "-"),
		Message:        getString(logParts["message"], getString(logParts["content"], "-")),
		MsgID:          getString(logParts["msg_id"], "-"),
		Priority:       logParts["priority"].(int),
		ProcId:         getString(logParts["proc_id"], "-"),
		Severity:       logParts["severity"].(int),
		StructuredData: sd,
		Timestamp:      logParts["timestamp"].(time.Time),
		TLSPeer:        logParts["tls_peer"].(string),
		Version:        getInt(logParts["version"], 1),
	}
}

func getString(v interface{}, def string) string {
	if v == nil {
		return def
	} else {
		return v.(string)
	}
}

func getInt(v interface{}, def int) int {
	if v == nil {
		return def
	} else {
		return v.(int)
	}
}

func parseStructuredData(s string) map[string]string {
	m := make(map[string]string)

	replacer := strings.NewReplacer("[", "", "]", "")
	s = replacer.Replace(s)
	items := strings.Split(s, " ")

	for _, i := range items {
		at := strings.Index(i, "@")
		if at >= 0 {
			kv := strings.Split(i, "@")
			if len(kv) < 2 {
				log.Error().Msgf("Error parsing structured data item: '%v'", i)
			} else {
				m[kv[0]] = kv[1]
			}
		}

		equal := strings.Index(i, "=")
		if equal >= 0 {
			kv := strings.Split(i, "=")
			if len(kv) < 2 {
				log.Error().Msgf("Error parsing structured data item: '%v'", i)
			} else {
				m[kv[0]] = strings.Replace(kv[1], "\"", "", -1)
			}
		}
	}
	return m
}
