/*
Package cmd provides CLI functionality.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintError(err error) {
	fmt.Fprintf(os.Stdout, "%s %s\n", color.HiRedString("Error:"), err.Error())
}
