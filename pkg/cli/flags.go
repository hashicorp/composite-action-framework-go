package cli

import (
	"flag"
	"reflect"
)

// Flags represents a type that sets options based on
// a set of command line flags.
type Flags interface {
	Flags(*flag.FlagSet)
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

type multiFlags []Flags

func (mf multiFlags) Flags(fs *flag.FlagSet) {
	for _, f := range mf {
		f.Flags(fs)
	}
}

var flagsType = reflect.TypeOf(new(Flags)).Elem()

func makeFlags(maybeFlags any) Flags {
	var mf multiFlags
	if f, ok := maybeFlags.(Flags); ok {
		mf = append(mf, f)
	}
	t := reflect.TypeOf(maybeFlags)
	if t.Kind() != reflect.Pointer {
		return mf
	}
	v := reflect.ValueOf(maybeFlags).Elem()
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return mf
	}

	for i := 0; i < t.NumField(); i++ {
		fieldVal := v.Field(i)
		if fieldVal.Kind() == reflect.Pointer {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		} else {
			fieldVal = fieldVal.Addr()
		}
		if fieldVal.Type().AssignableTo(flagsType) {
			mf = append(mf, fieldVal.Interface().(Flags))
		}
	}

	return mf
}
