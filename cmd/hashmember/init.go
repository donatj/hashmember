package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/donatj/hashmember"
	"github.com/google/subcommands"
)

type initCmd struct{}

func (*initCmd) Name() string     { return "init" }
func (*initCmd) Synopsis() string { return "Initialize hashmember file" }
func (p *initCmd) Usage() string {
	return fmt.Sprintf(`init <filename>:
	%s
`, p.Synopsis())
}

func (p *initCmd) SetFlags(f *flag.FlagSet) {}
func (p *initCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Println("Expects exactly one filename to initilize")
		return subcommands.ExitUsageError
	}

	hm := hashmember.New()
	err := saveHashmember(fs.Arg(0), hm)
	if err != nil {
		log.Printf("Error encoding hashmember: %s", err)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
