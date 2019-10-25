package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&uploadCmd{}, "")

	flag.Parse()
	os.Exit(int(subcommands.Execute(context.Background())))
}

func scanenv(vars ...string) string {
	val := ""
	for _, v := range vars {
		val = os.Getenv(v)
		if val != "" {
			break
		}
	}

	return val
}
