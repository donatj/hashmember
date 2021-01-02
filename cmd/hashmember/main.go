package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	_ "github.com/seiflotfy/cuckoofilter"
)

func init() {
	flag.Parse()
}

func main() {
	subcommands.Register(&initCmd{}, "setup")

	subcommands.Register(&insertCmd{}, "filter")
	subcommands.Register(&lookupCmd{}, "filter")

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
