package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/subcommands"
)

type uploadCmd struct {
	title string
	sha   string
	enp   string
}

func (*uploadCmd) Name() string     { return "upload" }
func (*uploadCmd) Synopsis() string { return "upload an inspection file" }
func (*uploadCmd) Usage() string {
	return `upload <filename>:
	Upload an inspection.
  `
}

func (p *uploadCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.title, "title", "default", "Inspection title - must set when uploading multiple inspections for a single SHA")
	f.StringVar(&p.sha, "work-sha", "", "The hash of the commit being analyzed.")
	f.StringVar(&p.enp, "endpoint", "https://methodci.org/api/user/ci/data", "The anaysis endpoints.")
}

func (p *uploadCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Println("Expects exactly one filename argument")
		return subcommands.ExitUsageError
	}

	filename := fs.Arg(0)

	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.IsDir() {
		fmt.Printf("File '%s' not found\n", filename)
		return subcommands.ExitFailure
	}

	if p.sha == "" {
		p.sha = scanenv(
			"GITHUB_SHA",
			"DRONE_COMMIT",
			"TRAVIS_COMMIT",
			"CI_COMMIT_SHA",
		)
	}

	if p.sha == "" {
		log.Println("unable to determine the working commit")
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
