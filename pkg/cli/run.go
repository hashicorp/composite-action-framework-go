// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package cli

import "fmt"

func runCLI(c *Command, args []string) error {
	if helpFunc := helpRequested(c, args); helpFunc != nil {
		return helpFunc()
	}
	if err := parseEnv(c); err != nil {
		return err
	}
	subArgs, err := parseFlags(c, args[1:])
	if err != nil {
		return err
	}
	if len(subArgs) == 0 {
		return run(c, nil)
	}
	sub := subArgs[0]
	if len(c.Subcommands()) == 0 {
		return run(c, subArgs)
	}
	sc, ok := getSubCommand(c, sub)
	if !ok {
		return fmt.Errorf("subcommand %q not found", sub)
	}
	return runCLI(sc, subArgs)
}

func run(c *Command, args []string) error {
	if c.Run() == nil {
		return ErrNotImplemented
	}
	if err := parseArgs(c, args); err != nil {
		return err
	}
	if err := initOpts(c); err != nil {
		return err
	}
	return c.Run()()
}

func helpRequested(c *Command, args []string) func() error {
	if len(args) < 2 {
		return nil
	}
	switch args[1] {
	default:
		return nil
	case "-h", "--h", "-help", "--help", "?":
	}
	return func() error { return c.printHelp(c.stdout) }
}
