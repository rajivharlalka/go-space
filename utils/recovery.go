package utils

import "fmt"

func RecoverServer() {
	if r := recover(); r != nil {
		fmt.Println("Recovered", r)
	}
}
