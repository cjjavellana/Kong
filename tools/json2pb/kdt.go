package main

import "regexp"

// Known Data Types

var (
	listenerDataTypeFields, _ = regexp.Compile("listener|proxy_protocol|reuseport|backlog=\\d+|deferred|ssl|ip|port|http2|bind")
)

func keyValuePairMessageDefinition() string {
	return `
import "google/protobuf/any.proto";

message KeyValuePair {
	string name = 1 [json_name="name"];
	string value = 2 [json_name="value"];
}
`
}

func listenerMessageDefinition() string {
	return `
message Listener {
	string listener = 1 [json_name="listener"];
	bool proxyProtocol = 2 [json_name="proxy_protocol"];
	bool reusePort = 3 [json_name="reuseport"];
	bool backlog = 4 [json_name="backlog"];
	bool deferred = 5 [json_name="deferred"];
	bool ssl = 6 [json_name="ssl"];
	string ip = 7 [json_name="ip"];
	int32 port = 8 [json_name="port"];
	bool http2 = 9 [json_name="http2"];
	bool bind = 10 [json_name="bind"];
}
`
}

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
