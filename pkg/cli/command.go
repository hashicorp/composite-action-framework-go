package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Command represents a command in the CLI graph.
// Don't construct Commands manually, instead use the RootCommand
// and LeafCommand functions to construct root and leaf commands.
type Command struct {
	name  string
	desc  string
	run   func() error
	flags Flags
	args  Args
	env   Env
	init  Init
	subs  []*Command

	// Runtime
	flagSet        *flag.FlagSet
	stdout, stderr io.Writer
	stdin          io.Reader
}

func (c *Command) Name() string                { return c.name }
func (c *Command) Description() string         { return c.desc }
func (c *Command) Run() func() error           { return c.run }
func (c *Command) Flags() Flags                { return c.flags }
func (c *Command) Args() Args                  { return c.args }
func (c *Command) Env() Env                    { return c.env }
func (c *Command) Init() Init                  { return c.init }
func (c *Command) Subcommands() []*Command     { return c.subs }
func (c *Command) Execute(args []string) error { return runCLI(c, args) }

func getSubCommand(parent *Command, name string) (*Command, bool) {
	for _, c := range parent.Subcommands() {
		if c.Name() == name {
			return c, true
		}
	}
	return nil, false
}

// RootCommand is a command that only contains subcommands and doesn't do anything
// by itself.
func RootCommand(name, desc string, subcommands ...*Command) *Command {
	c := &Command{
		name: name,
		desc: desc,
		subs: subcommands,

		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  os.Stdin,
	}
	c.run = func() error {
		return c.printHelp(c.stderr)
	}
	return c
}

// None can be used as the parameter for leaf commands that have no options.
type None = *any

// LeafCommand is a command that runs something.
// The run function can accept a opts argument, which should be a pointer to a
// type that implements Flags or Args, or both. If no flags or args are needed,
// it can accept None instead.
//
// When the command is run, a new instance of opts (*T) is created first, then its
// flags and args are handled. The resultant opts is passed to the run function.
//
// If opts implements Flags, then its flags are registered with the flag set.
// The Flag set is parsed before the args.
// If opts implements Args, then its ParseArgs method is called on args remaining
// after the flag set is parsed.
// If opts implements Env, then its ReadEnv method is called to populate it with
// config from the environment.
// The run function is called after flags and args have been parsed, and passed
// the resultant opts.
func LeafCommand[T any](name, desc string, run func(opts *T) error) *Command {
	opts, optionSet := makeOptionSet[T]()
	return &Command{
		name:  name,
		desc:  desc,
		env:   optionSet.env,
		flags: optionSet.flags,
		args:  optionSet.args,
		init:  optionSet.init,
		run:   func() error { return run(opts) },

		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  os.Stdin,
	}
}

type optionSet struct {
	flags Flags
	args  Args
	env   Env
	init  Init
}

func makeOptionSet[T any]() (*T, optionSet) {
	opts := new(T)
	// It's ok for all/eny of flags, args, env to be nil.
	var os optionSet
	os.flags = makeFlags(opts)
	os.args, _ = any(opts).(Args)
	os.env, _ = any(opts).(Env)
	os.init, _ = any(opts).(Init)

	return opts, os
}

func (c *Command) printHelp(w io.Writer) error {
	fmt.Fprintf(w, "%s - %s\n\n", c.name, c.desc)
	fs := createFlagSet(c)
	if fs != nil {
		fs.SetOutput(w)
		fs.PrintDefaults()
	}
	if len(c.subs) == 0 {
		return nil
	}
	fmt.Fprint(w, "Subcommands:\n\n")
	return TabWrite(w, c.subs, func(c *Command) string {
		return fmt.Sprintf("\t%s\t%s", c.Name(), c.Description())
	})
}

func (c *Command) SetStdout(w io.Writer) {
	c.stdout = w
	for _, s := range c.subs {
		s.SetStdout(w)
	}
}

func (c *Command) SetStderr(w io.Writer) {
	c.stderr = w
	for _, s := range c.subs {
		s.SetStderr(w)
	}
}

func (c *Command) SetStdin(r io.Reader) {
	c.stdin = r
	for _, s := range c.subs {
		s.SetStdin(r)
	}
}
