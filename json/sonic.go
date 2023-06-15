//go:build sonic && avx && (linux || windows || darwin) && amd64
// +build sonic
// +build avx
// +build linux windows darwin
// +build amd64

package json

import (
	"github.com/bytedance/sonic"
)

var (
	json          = sonic.ConfigStd
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
)
