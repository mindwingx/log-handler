package utils

import (
	"fmt"
	"github.com/google/uuid"
)

func PanicHandler() (res string) {
	rec := recover()
	if rec != nil {
		res = fmt.Sprintf("RECOVER: %s", rec)
	}

	return
}

func NewUuid() string {
	val, _ := uuid.NewRandom()
	return val.String()
}
