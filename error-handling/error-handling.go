package errorHandling

import (
	"fmt"
	"os"
)

func LogErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
