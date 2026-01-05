// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"flag"
)

// Flags represents a type that sets options based on
// a set of command line flags.
type Flags interface {
	Flags(*flag.FlagSet)
}

type FlagHider interface {
	HideFlags() []string
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
	if fh := c.flagHider; fh != nil {
		for _, name := range fh.HideFlags() {
			c.hideFlagsFromSynopsis[name] = nil
		}
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
