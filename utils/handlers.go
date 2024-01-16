package utils

import "fmt"

func PanicHandler() (res string) {
	rec := recover()
	if rec != nil {
		res = fmt.Sprintf("RECOVER: %s", rec)
	}

	return
}
