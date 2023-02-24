package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Command represents a command in the CLI graph.
// Don't construct Commands manually, instead use the RootCommand
// and LeafCommand functions to construct root and leaf commands.
type Command struct {
	name, desc, help string
	run              func() error
	optionSet
	subs   []*Command
	parent *Command

	hideFlagsFromSynopsis map[string]any

	// Runtime
	flagSet        *flag.FlagSet
	stdout, stderr io.Writer
	stdin          io.Reader
}

func (c *Command) Name() string        { return c.name }
func (c *Command) Description() string { return c.desc }
func (c *Command) Help() string {
	return strings.TrimSpace(fmt.Sprintf("%s\n\n%s", c.Usage(), strings.TrimSpace(c.help)))
}
func (c *Command) Run() func() error       { return c.run }
func (c *Command) Flags() Flags            { return c.flags }
func (c *Command) Args() Args              { return c.args }
func (c *Command) Env() Env                { return c.env }
func (c *Command) Init() Init              { return c.init }
func (c *Command) Subcommands() []*Command { return c.subs }

func (c *Command) Path() []string {
	curr := []string{c.Name()}
	for c.parent != nil {
		c = c.parent
		curr = append([]string{c.Name()}, curr...)
	}
	return curr
}

func (c *Command) PathString() string {
	return strings.Join(c.Path(), " ")
}

// Execute should be called on the root command, it is the starting point for evaluating
// args and routing to the requested command.
func (c *Command) Execute(args []string) error { return runCLI(c, args) }

func (c *Command) WithHelp(h string) *Command { c.help = h; return c }

func (c *Command) Usage() string {
	buf := &bytes.Buffer{}
	fs := createFlagSet(c)
	if fs != nil {
		fs.SetOutput(buf)
		fs.Usage()
		return strings.TrimSpace(buf.String())
	}
	return fmt.Sprintf("Usage of %s:", c.name)
}

func (c *Command) Synopsis() string {
	buf := &bytes.Buffer{}
	fs := createFlagSet(c)
	if fs != nil {
		fs.VisitAll(func(f *flag.Flag) {
			if _, ok := c.hideFlagsFromSynopsis[f.Name]; ok {
				return
			}
			if f.DefValue == "true" || f.DefValue == "false" {
				fmt.Fprintf(buf, "[-%s] ", f.Name)
			} else if f.DefValue == "" {
				fmt.Fprintf(buf, "[-%s=%s] ", f.Name, strings.ToUpper(f.Name))
			} else {
				fmt.Fprintf(buf, "[-%s=%s (%s)] ", f.Name, strings.ToUpper(f.Name), f.DefValue)
			}
		})
	}
	if argList := makeArgList(c); argList != nil {
		for _, a := range argList {
			if a.required && !a.variadic {
				fmt.Fprintf(buf, "<%s>", a.name)
			} else if !a.required && !a.variadic {
				fmt.Fprintf(buf, "[%s (%s)]", a.name, a.defaultVal)
			} else if a.required && a.variadic {
				fmt.Fprintf(buf, "<")
				for i := 0; i < a.minVals; i++ {
					fmt.Fprintf(buf, "%s%d, ", a.name, i)
				}
				fmt.Fprintf(buf, "...>")
			} else if !a.required && a.variadic {
				fmt.Fprintf(buf, "[%s...](%s)", a.name, strings.Join(a.defaultVals, " "))
			} else {
				panic("logical error with arg handling; please alert the maintainers")
			}
		}
	}
	return buf.String()
}

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
		name:                  name,
		desc:                  desc,
		subs:                  subcommands,
		stdout:                os.Stdout,
		stderr:                os.Stderr,
		stdin:                 os.Stdin,
		hideFlagsFromSynopsis: map[string]any{},
	}
	c.run = func() error {
		return c.printHelp(c.stderr)
	}
	// Track commands' parents so we can generate the command's full path.
	for _, s := range subcommands {
		s.parent = c
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
		name:                  name,
		desc:                  desc,
		optionSet:             optionSet,
		run:                   func() error { return run(opts) },
		stdout:                os.Stdout,
		stderr:                os.Stderr,
		stdin:                 os.Stdin,
		hideFlagsFromSynopsis: map[string]any{},
	}
}

type optionSet struct {
	flags      Flags
	flagHider  FlagHider
	args       Args
	argDefiner ArgDefiner
	env        Env
	init       Init
}

func makeOptionSet[T any]() (*T, optionSet) {
	opts := new(T)
	// It's ok for all/eny of flags, args, env, init to be nil.
	var os optionSet
	os.flags, _ = any(opts).(Flags)
	os.flagHider, _ = any(opts).(FlagHider)
	os.args, _ = any(opts).(Args)
	os.argDefiner, _ = any(opts).(ArgDefiner)
	os.env, _ = any(opts).(Env)
	os.init, _ = any(opts).(Init)

	if os.args != nil && os.argDefiner != nil {
		panic("opts cannot implement both Args and ArgDefiner")
	}

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
