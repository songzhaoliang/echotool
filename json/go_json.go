//go:build go_json
// +build go_json

package json

import (
	json "github.com/goccy/go-json"
)

var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
)
