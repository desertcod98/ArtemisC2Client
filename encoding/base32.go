package encoding

import (
	"encoding/base32"
	"strings"
)

func Base32Encode(data string) string {
	encoded := base32.StdEncoding.EncodeToString([]byte(data))
	return strings.TrimRight(encoded, "=")
}
