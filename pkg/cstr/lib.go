package cstr

import (
	"fmt"
	"strings"

	uuid "git.cmcode.dev/cmcode/uuid"
)

type Config struct {
	// Connection string for the Castopod mariadb database.
	SQLConnectionString string `json:"sqlConnectionString"`
}

// NewUUID returns a double uuid with the hyphens removed, leading to a
// 64-character string. Castopod appears to use something like this for its
// token.
func NewUUID() string {
	return strings.ReplaceAll(fmt.Sprintf("%v%v", uuid.New(), uuid.New()), "-", "")
}
