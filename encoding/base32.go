package encoding

import (
	"encoding/base32"
)

var rawBase32 = base32.StdEncoding.WithPadding(base32.NoPadding)

// func Base32Encode(data []byte) string {
// 	encoded := base32.StdEncoding.EncodeToString(data)
// 	return strings.TrimRight(encoded, "=")
// }

// func Base32DecodeNoPadding(input string) ([]byte, error) {
// 	if input == "" {
// 		return []byte{}, nil
// 	}

// 	paddingNeeded := (8 - (len(input) % 8)) % 8
// 	if paddingNeeded > 0 {
// 		input = input + strings.Repeat("=", paddingNeeded)
// 	}

// 	decoded, err := base32.StdEncoding.DecodeString(input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return decoded, nil
// }

func Base32Encode(data []byte) string {
	return rawBase32.EncodeToString(data)
}

func Base32Decode(input string) ([]byte, error) {
	return rawBase32.DecodeString(input)
}
