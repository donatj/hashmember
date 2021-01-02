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

type lookupCmd struct {
	json bool
}

func (*lookupCmd) Name() string     { return "lookup" }
func (*lookupCmd) Synopsis() string { return "Lookup a value into a hashmember file" }
func (p *lookupCmd) Usage() string {
	return fmt.Sprintf(`lookup <filename> <value>:
	%s
`, p.Synopsis())
}

func (p *lookupCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.json, "json", false, "output json")
}
func (p *lookupCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		log.Println("Expects exactly two aguments - a filename and a value to lookup")
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

	known := hm.Lookup([]byte(fs.Arg(1)))

	if p.json {
		if known {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}

		return subcommands.ExitSuccess
	}

	if known {
		return ExitKnown
	}

	return subcommands.ExitSuccess
}

// ExitKnown is returned when a lookup of a hashmember is known
var ExitKnown subcommands.ExitStatus = 5
