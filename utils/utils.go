package utils

import (
	"crypto/md5"
	"encoding/base64"
)

// GetMD5B64 was shamelessly stolen from github.com/alphagov/paas-rds-broker
func GetMD5B64(text string, maxLength int) string {
	md5 := md5.Sum([]byte(text))
	encoded := base64.URLEncoding.EncodeToString(md5[:])
	if len(encoded) > maxLength {
		return encoded[0:maxLength]
	} else {
		return encoded
	}
}
