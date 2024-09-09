package main

import (
	"encoding/json"
	"fmt"
	"os"

	cstr "github.com/charles-m-knox/castopod-sub-token-retriever/pkg/cstr"
)

var conf cstr.Config

// LoadConfig reads from file f and applies sensible defaults to values not
// specifically set by the user.
func LoadConfig(f string) (cstr.Config, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return cstr.Config{}, fmt.Errorf("failed to load config from %v: %v", f, err)
	}

	var c cstr.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		return cstr.Config{}, fmt.Errorf("failed to unmarshal config from %v: %v", f, err)
	}

	return c, nil
}
