package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	cmd "github.com/tcfw/evntsrc/pkg/storer/cmd"
)

func main() {
	command := cmd.NewDefaultCommand()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
