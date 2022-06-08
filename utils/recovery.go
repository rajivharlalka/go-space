package utils

import (
	"fmt"
	"runtime/debug"
)

func RecoverServer() {
	if r := recover(); r != nil {
		fmt.Println("Recovered", r)
		fmt.Println("stacktrace", string(debug.Stack()))
	}
}
