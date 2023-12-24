package main

import (
	"os"

	"github.com/Denis-Carlos-Farias/go-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
