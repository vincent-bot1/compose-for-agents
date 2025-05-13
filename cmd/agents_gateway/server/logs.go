package server

import (
	"fmt"
	"os"
)

func log(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}
