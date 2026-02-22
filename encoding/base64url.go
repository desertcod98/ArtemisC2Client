package encoding

import (
	"encoding/base64"
)

var rawBase64Url = base64.URLEncoding.WithPadding(base64.NoPadding)

func Base64UrlEncode(data []byte) string {
	return rawBase64Url.EncodeToString(data)
}

func Base64UrlDecode(input string) ([]byte, error) {
	return rawBase64Url.DecodeString(input)
}
