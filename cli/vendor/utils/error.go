package utils

import (
	"fmt"
	"os"
)

func PrintError(error string) {
	fmt.Println(error)
	os.Exit(1)
}
