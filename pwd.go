package sdstudio

import (
	"os"
)

func Root() (path string) {
	path, _ = os.Getwd()
	return
}
