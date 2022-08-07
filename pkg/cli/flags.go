package cli

import (
	"flag"
)

// Flags represents a type that sets options based on
// a set of command line flags.
type Flags interface {
	Flags(*flag.FlagSet)
}

func FlagsAll(fs *flag.FlagSet, objs ...Flags) {
	for _, f := range objs {
		f.Flags(fs)
	}
}

func FlagFuncsAll(fs *flag.FlagSet, flagFuncs ...func(*flag.FlagSet)) {
	for _, f := range flagFuncs {
		f(fs)
	}
}

func createFlagSet(c *Command) *flag.FlagSet {
	var fs *flag.FlagSet
	if f := c.Flags(); f != nil {
		fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		f.Flags(fs)
	}
	return fs
}

func parseFlags(c *Command, args []string) ([]string, error) {
	if c.flagSet = createFlagSet(c); c.flagSet == nil {
		return args, nil
	}
	if err := c.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return c.flagSet.Args(), nil
}
