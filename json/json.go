//go:build !jsoniter && !go_json && !(sonic && avx && (linux || windows || darwin) && amd64)
// +build !jsoniter
// +build !go_json
// +build !sonic !avx !linux,!windows,!darwin !amd64

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
