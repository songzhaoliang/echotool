//go:build !jsoniter && !go_json
// +build !jsoniter,!go_json

package json

import (
	"encoding/json"
)

var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
)
