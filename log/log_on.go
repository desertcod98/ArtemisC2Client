//go:build debug

package log

import "fmt"

func Log(a ...any) {
	fmt.Println(a...)
}
