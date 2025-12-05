// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strings"
)

type Args interface {
	ParseArgs([]string) error
}

type ArgDefiner interface {
	Args(*ArgList)
}

func parseArgs(c *Command, args []string) error {
	if a := c.Args(); a != nil {
		return a.ParseArgs(args)
	}
	if a := c; a != nil {
		l := makeArgList(c)
		return l.parseArgs(args)
	}
	if len(args) != 0 {
		return ErrNoArgsAllowed
	}
	return nil
}

func makeArgList(c *Command) ArgList {
	if c.argDefiner == nil {
		return nil
	}
	al := new(ArgList)
	c.argDefiner.Args(al)
	return *al
}

type Arg struct {
	name string

	required   bool
	defaultVal string
	val        *string

	variadic    bool
	listVal     *[]string
	defaultVals []string
	minVals     int
}

type ArgList []Arg

func (al *ArgList) parseArgs(args []string) error {
	for i, a := range *al {
		if a.variadic {
			return a.parseVariadic(args[i:])
		}
		if len(args) < i+1 {
			if a.required {
				return fmt.Errorf("required argument missing: %s", a.name)
			}
			*a.val = a.defaultVal
			return nil
		} else {
			*a.val = args[i]
		}
	}
	return nil
}

func (a Arg) parseVariadic(args []string) error {
	if !a.required {
		if len(args) != 0 {
			(*a.listVal) = args
		} else {
			(*a.listVal) = a.defaultVals
		}
		return nil
	}
	if len(args) < a.minVals {
		return fmt.Errorf("required %s argument(s) missing; you must supply at least %d", a.name, a.minVals)
	}
	(*a.listVal) = args
	return nil
}

func (al *ArgList) assertLast() {
	if len(*al) == 0 {
		return
	}
	last := (*al)[len(*al)-1]
	if last.variadic {
		panic("can't put any further args after a variadic arg")
	}
	if !last.required {
		panic("can't put any further args after an optional arg")
	}
}

func (al *ArgList) add(a Arg) {
	al.assertLast()
	a.name = strings.ToUpper(a.name)
	(*al) = append(*al, a)
}

func (al *ArgList) Required(val *string, name string) {
	al.add(Arg{val: val, name: name, required: true})
}

func (al *ArgList) Optional(val *string, name, defaultVal string) {
	al.add(Arg{val: val, name: name, defaultVal: defaultVal})
}

func (al *ArgList) OptionalVariadic(vals *[]string, name string, defaultVals ...string) {
	al.add(Arg{variadic: true, listVal: vals, name: name, defaultVals: defaultVals})
}

func (al *ArgList) RequiredVariadic(vals *[]string, name string, minimumVals int) {
	if minimumVals < 1 {
		panic(fmt.Sprintf("cannot require %d values; must be > 0", minimumVals))
	}
	al.add(Arg{variadic: true, listVal: vals, name: name, required: true, minVals: minimumVals})
}
