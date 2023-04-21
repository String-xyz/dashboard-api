package main

import (
	"os"
)

func main() {
	var script string
	if len(os.Args) > 1 {
		script = os.Args[1]
	}

	// if script == "something" {
	// scripts.Something()
	// }
}
