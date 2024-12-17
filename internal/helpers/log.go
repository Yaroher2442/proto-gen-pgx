package helpers

import (
	"fmt"
	"os"
)

func Log(str string) {
	os.Stderr.WriteString(fmt.Sprintln(str))
}
