package main

import "regexp"

// Known Data Types

const ListenerDataTypeFields = "listener|proxy_protocol|reuseport|backlog=\\d+|deferred|ssl|ip|port|http2|bind"

var (
	listenerDataTypeFields, _ = regexp.Compile(ListenerDataTypeFields)
)

// Checks whether map `m` has a `name` and `value` members
func isKeyValuePairMap(m map[string]interface{}) bool {
	return m["name"] != nil && m["value"] != nil
}

// checks whether map `m` represents a `Listener` data type
// {
// 	"listener": "0.0.0.0:8005",
//  "proxy_protocol": false,
// 	"reuseport": false,
//  "backlog=%d+": false,
//  "deferred": false,
// 	"ssl": false,
// 	"ip": "0.0.0.0",
//  "port": 8005,
// 	"http2": false,
// 	"bind": false
//  }
func isListener(m map[string]interface{}) bool {
	for k := range m {
		if !listenerDataTypeFields.Match([]byte(k)) {
			return false
		}
	}

	return true
}
