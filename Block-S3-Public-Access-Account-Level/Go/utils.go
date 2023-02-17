/* utility funcs to be used by main.go */

package main

import (
	"fmt"
	"strings"
)

func parseResults(status string) bool {
	return strings.Contains(status, "false")
}

func printOut(message string) {
	fmt.Printf("success: %v\n", message)
}