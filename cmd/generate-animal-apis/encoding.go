package main

import (
	"crypto/sha1"
	"encoding/base32"
	"strings"
)

// Encode a name as its base32-encoded SHA1 hash to normalize its length.
// The result is 32 character string (160 bits/5 bits-per-character).
// Return values are lower-cased for readability and consistency with a
// preference for lower-cased IDs.
func encodeName(name string) string {
	hash := sha1.New()
	hash.Write([]byte(name))
	return strings.ToLower(base32.StdEncoding.EncodeToString(hash.Sum(nil)))
}
