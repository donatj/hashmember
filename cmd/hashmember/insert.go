package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/donatj/hashmember"
	"github.com/google/subcommands"
)

type insertCmd struct{}

func (*insertCmd) Name() string     { return "insert" }
func (*insertCmd) Synopsis() string { return "insert a value into a hashmember file" }
func (p *insertCmd) Usage() string {
	return fmt.Sprintf(`insert <filename> <value>:
	%s
`, p.Synopsis())
}

func (p *insertCmd) SetFlags(f *flag.FlagSet) {}
func (p *insertCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		log.Println("Expects exactly two aguments - a filename and a value to insert")
		return subcommands.ExitUsageError
	}

	hm, err := loadHashmember(fs.Arg(0))
	if errors.Is(err, hashmember.ErrUnhandledVersion) {
		log.Println("Unable to read hashmember file. It appears to have been created with a newer version of hashmember.")

		return subcommands.ExitFailure
	} else if err != nil {
		log.Printf("Failed to read: %s", err)

		return subcommands.ExitFailure
	}

	hm.Insert([]byte(fs.Arg(1)))

	err = saveHashmember(fs.Arg(0), hm)
	if err != nil {
		log.Printf("Failed to write: %s", err)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
