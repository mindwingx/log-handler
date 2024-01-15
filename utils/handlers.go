package utils

import "fmt"

func PanicHandler() {
	rec := recover()
	if rec != nil {
		fmt.Println("RECOVER: ", rec)
	}
}
