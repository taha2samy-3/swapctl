package main

import (
	"fmt"
	"os"

	"github.com/taha2samy/swapctl/internal/cli"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("Error: This tool must be run with root privileges (sudo).")
		os.Exit(1)
	}

	cli.StartInteractiveSession()
}
