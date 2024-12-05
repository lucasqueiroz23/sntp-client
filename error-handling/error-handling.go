package errorHandling

import (
	"fmt"
	"os"
)

func LogErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, "erro: ", err)
	os.Exit(1)
}
