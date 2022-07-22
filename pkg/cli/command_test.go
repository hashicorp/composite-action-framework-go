package cli

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/composite-action-framework-go/pkg/testhelpers/assert"
)

// args constructs a slice like os.Args, setting the first
// arg to the empty string, which represents the command
// name used to call the CLI.
func args(a ...string) []string {
	return append([]string{""}, a...)
}

type testFlags struct {
	flag1, flag2 bool
}

func (o *testFlags) Flags(fs *flag.FlagSet) {
	fs.BoolVar(&o.flag1, "flag1", false, "flag1 desc")
	fs.BoolVar(&o.flag2, "flag2", false, "flag2 desc")
}

type testArgs struct {
	args []string
}

func (a *testArgs) ParseArgs(args []string) error {
	a.args = args
	return nil
}

type testEnv struct {
	home string
}

func (e *testEnv) ReadEnv() error {
	e.home = "/test/home"
	return nil
}

type testOpts struct {
	testFlags
	testArgs
	testEnv
}

type envFlagArgsOpts struct {
	env            string
	envArg         string
	envArgInit     string
	envFlag        string
	envFlagArg     string
	envFlagInit    string
	envFlagArgInit string

	flag        string
	flagArg     string
	flagInit    string
	flagArgInit string

	arg     string
	argInit string

	init string
}

func (efa *envFlagArgsOpts) Init() error {
	efa.envArgInit = "init"
	efa.envFlagInit = "init"
	efa.envFlagArgInit = "init"
	efa.flagInit = "init"
	efa.flagArgInit = "init"
	efa.argInit = "init"
	efa.init = "init"
	return nil
}

func (efa *envFlagArgsOpts) ReadEnv() error {
	efa.env = "env"
	efa.envFlag = "env"
	efa.envArg = "env"
	efa.envFlagArg = "env"
	return nil
}

func (efa *envFlagArgsOpts) Flags(fs *flag.FlagSet) {
	fs.StringVar(&efa.envFlag, "envFlag", "", "env overridden by flag")
	fs.StringVar(&efa.envFlagArg, "envFlagArg", "", "env overridden by flag and arg")
	fs.StringVar(&efa.envFlagInit, "envFlagInit", "", "env overridden by flag and init")
	fs.StringVar(&efa.envFlagArgInit, "envFlagArgInit", "", "env overridden by flag, arg, and init")

	fs.StringVar(&efa.flag, "flag", "", "flag only")
	fs.StringVar(&efa.flagArg, "flagArg", "", "flag overridden by arg")
	fs.StringVar(&efa.flagInit, "flagInit", "", "flag overridden by init")
	fs.StringVar(&efa.flagInit, "flagArgInit", "", "flag overridden by arg and init")

}

func (efa *envFlagArgsOpts) ParseArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly 1 arg required")
	}
	arg := args[0]
	efa.envArg = arg
	efa.envArgInit = arg
	efa.envFlagArg = arg
	efa.envFlagArgInit = arg

	efa.flagArg = arg
	efa.flagArgInit = arg

	efa.arg = arg
	efa.argInit = arg

	return nil
}

func testCLI() (*Command, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	write := func(a ...any) error {
		s := make([]string, len(a))
		for i, item := range a {
			s[i] = fmt.Sprint(item)
		}
		_, err := buf.WriteString(strings.Join(s, ", "))
		return err
	}
	root := RootCommand("root", "root command",
		LeafCommand("leaf", "leaf command", func(None) error {
			return write("leaf")
		}),
		LeafCommand("leaf2", "leaf command 2", func(flags *testFlags) error {
			return write("leaf2", flags.flag1, flags.flag2)
		}),
		RootCommand("root2", "root command 2",
			LeafCommand("leaf3", "leaf command 3", func(None) error {
				return write("leaf3")
			}),
		),
		LeafCommand("leaf4", "leaf command 4", func(a *testArgs) error {
			return write("leaf4", strings.Join(a.args, ", "))
		}),
		LeafCommand("leaf5", "leaf command 5", func(o *testOpts) error {
			return write("leaf5", o.flag1, o.flag2, strings.Join(o.args, ", "))
		}),
		LeafCommand("leaf6", "leaf command 6", func(e *testEnv) error {
			return write("leaf6", e.home)
		}),
		LeafCommand("leaf7", "leaf command 7", func(o *testOpts) error {
			return write("leaf7", o.flag1, o.flag2, strings.Join(o.args, ", "), o.home)
		}),
		LeafCommand("leaf8", "leaf command 8", func(o *envFlagArgsOpts) error {
			return write("leaf8",
				o.env,
				o.envArg,
				o.envArgInit,
				o.envFlag,
				o.envFlagArg,
				o.envFlagInit,
				o.envFlagArgInit,
				o.flag,
				o.flagArg,
				o.flagInit,
				o.flagArgInit,
				o.arg,
				o.argInit,
				o.init,
			)
		}),
	)

	root.SetStdout(buf)

	return root, buf
}

func TestCommand_ok(t *testing.T) {

	cases := []struct {
		args []string
		want string
	}{
		{
			args(),
			"",
		},
		{
			args("-h"), `
root - root command

Subcommands:

  leaf   leaf command
  leaf2  leaf command 2
  root2  root command 2
  leaf4  leaf command 4
  leaf5  leaf command 5
  leaf6  leaf command 6
  leaf7  leaf command 7
  leaf8  leaf command 8
			`,
		},
		{
			args("leaf2", "-h"), `
leaf2 - leaf command 2

  -flag1
        flag1 desc
  -flag2
        flag2 desc
		`,
		},
		{
			args("leaf"),
			"leaf",
		},
		{
			args("leaf2"),
			"leaf2, false, false",
		},
		{
			args("leaf2", "-flag1"),
			"leaf2, true, false",
		},
		{
			args("leaf2", "-flag2"),
			"leaf2, false, true",
		},
		{
			args("leaf2", "-flag1", "-flag2"),
			"leaf2, true, true",
		},
		{
			args("root2"),
			"",
		},
		{
			args("root2", "leaf3"),
			"leaf3",
		},
		{
			args("leaf4"),
			"leaf4, ",
		},
		{
			args("leaf4", "hello"),
			"leaf4, hello",
		},
		{
			args("leaf4", "hello", "world"),
			"leaf4, hello, world",
		},
		{
			args("leaf5", "hello", "world"),
			"leaf5, false, false, hello, world",
		},
		{
			args("leaf5", "-flag1", "hello", "world"),
			"leaf5, true, false, hello, world",
		},
		{
			args("leaf5", "-flag2", "hello", "world"),
			"leaf5, false, true, hello, world",
		},
		{
			args("leaf5", "-flag1", "-flag2", "hello", "world"),
			"leaf5, true, true, hello, world",
		},
		{
			args("leaf6"),
			"leaf6, /test/home",
		},
		{
			args("leaf7", "-flag1", "-flag2=false", "hello", "world"),
			"leaf7, true, false, hello, world, /test/home",
		},
		// Test that we apply env, flags, args, init in that order.
		{
			args("leaf8",
				"-envFlag=flag",
				"-envFlagArg=flag",
				"-envFlagInit=flag",
				"-envFlagArgInit=flag",
				"-flag=flag",
				"-flagArg=flag",
				"-flagInit=flag",
				"-flagArgInit=flag",
				"arg",
			),
			"leaf8, env, arg, init, flag, arg, init, init, flag, arg, init, init, arg, init, init",
		},
	}

	for _, c := range cases {
		args, want := c.args, c.want
		t.Run("", func(t *testing.T) {
			c, buf := testCLI()
			if err := c.Execute(args); err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			want = strings.TrimSpace(want)
			got = strings.TrimSpace(got)

			want = strings.ReplaceAll(want, "\t", "    ")
			got = strings.ReplaceAll(got, "\t", "    ")

			if got != want {
				t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
				assert.Equal(t, got, want)
			}
		})
	}
}
