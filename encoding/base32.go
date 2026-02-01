package encoding

import (
	"encoding/base32"
	"strings"
)

func Base32Encode(data []byte) string {
	encoded := base32.StdEncoding.EncodeToString(data)
	return strings.TrimRight(encoded, "=")
}
