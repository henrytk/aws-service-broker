package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hex(text string, maxLength int) string {
	md5 := md5.Sum([]byte(text))
	md5Hex := make([]byte, hex.EncodedLen(len(md5)))
	hex.Encode(md5Hex, md5[:])

	if len(md5Hex) > maxLength {
		return string(md5Hex[0:maxLength])
	} else {
		return string(md5Hex)
	}
}
