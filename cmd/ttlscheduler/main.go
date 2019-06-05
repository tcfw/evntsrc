package main

import (
	"fmt"
	"os"

	cmd "github.com/tcfw/evntsrc/internal/ttlscheduler/cmd"
)

func main() {
	command := cmd.NewDefaultCommand()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
