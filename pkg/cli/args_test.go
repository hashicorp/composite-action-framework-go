package cli

import (
	"errors"
	"testing"

	"github.com/hashicorp/composite-action-framework-go/pkg/testhelpers/assert"
)

func TestArgList_parseArgs(t *testing.T) {

	type testOpts struct {
		A1, A2, A3 string
		V1         []string
	}
	mkArgs := func(args ...string) []string { return args }

	cases := []struct {
		desc    string
		setup   func(*testOpts, *ArgList)
		args    []string
		want    *testOpts
		wantErr error
	}{
		{
			"zero",
			func(opts *testOpts, al *ArgList) {},
			mkArgs(),
			&testOpts{},
			nil,
		},
		{
			"required_provided",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
			},
			mkArgs("value1"),
			&testOpts{A1: "value1"},
			nil,
		},
		{
			"required_missing",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
			},
			mkArgs(),
			&testOpts{},
			errors.New("required argument missing: name1"),
		},
		{
			"optional_provided",
			func(opts *testOpts, al *ArgList) {
				al.Optional(&opts.A1, "name1", "default1")
			},
			mkArgs("value1"),
			&testOpts{A1: "value1"},
			nil,
		},
		{
			"optional_missing",
			func(opts *testOpts, al *ArgList) {
				al.Optional(&opts.A1, "name1", "default1")
			},
			mkArgs(),
			&testOpts{A1: "default1"},
			nil,
		},
		{
			"required_optional_provided",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.Optional(&opts.A2, "name2", "default2")
			},
			mkArgs("val1", "val2"),
			&testOpts{A1: "val1", A2: "val2"},
			nil,
		},
		{
			"required_optional_missing",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.Optional(&opts.A2, "name2", "default2")
			},
			mkArgs("val1"),
			&testOpts{A1: "val1", A2: "default2"},
			nil,
		},
		{
			"required_variadic_1_provided",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 1)
			},
			mkArgs("val1"),
			&testOpts{V1: mkArgs("val1")},
			nil,
		},
		{
			"required_variadic_1_provided_extra",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 1)
			},
			mkArgs("val1", "val2", "val3"),
			&testOpts{V1: mkArgs("val1", "val2", "val3")},
			nil,
		},
		{
			"required_variadic_2_provided",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 2)
			},
			mkArgs("val1", "val2"),
			&testOpts{V1: mkArgs("val1", "val2")},
			nil,
		},
		{
			"required_variadic_2_provided_extra",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 2)
			},
			mkArgs("val1", "val2", "val3"),
			&testOpts{V1: mkArgs("val1", "val2", "val3")},
			nil,
		},
		{
			"required_variadic_1_all_missing",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 1)
			},
			mkArgs(),
			&testOpts{},
			errors.New("required name1 argument(s) missing; you must supply at least 1"),
		},
		{
			"required_variadic_2_all_missing",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 2)
			},
			mkArgs(),
			&testOpts{},
			errors.New("required name1 argument(s) missing; you must supply at least 2"),
		},
		{
			"required_variadic_2_one_missing",
			func(opts *testOpts, al *ArgList) {
				al.RequiredVariadic(&opts.V1, "name1", 2)
			},
			mkArgs("val1"),
			&testOpts{},
			errors.New("required name1 argument(s) missing; you must supply at least 2"),
		},
		{
			"optional_variadic_provided",
			func(opts *testOpts, al *ArgList) {
				al.OptionalVariadic(&opts.V1, "name1", "def1", "def2", "def3")
			},
			mkArgs("val1", "val2", "val3"),
			&testOpts{V1: mkArgs("val1", "val2", "val3")},
			nil,
		},
		{
			"optional_variadic_missing",
			func(opts *testOpts, al *ArgList) {
				al.OptionalVariadic(&opts.V1, "name1", "def1", "def2", "def3")
			},
			mkArgs(),
			&testOpts{V1: mkArgs("def1", "def2", "def3")},
			nil,
		},
		{
			"required_providedx2",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.Required(&opts.A2, "name2")
			},
			mkArgs("val1", "val2"),
			&testOpts{A1: "val1", A2: "val2"},
			nil,
		},
		{
			"required_provided_required_missing",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.Required(&opts.A2, "name2")
			},
			mkArgs("val1"),
			nil,
			errors.New("required argument missing: name2"),
		},
		{
			"required_provided_optional_provided",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.Optional(&opts.A2, "name2", "def2")
			},
			mkArgs("val1", "val2"),
			&testOpts{A1: "val1", A2: "val2"},
			nil,
		},
		{
			"required_provided_optional_missing",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.Optional(&opts.A2, "name2", "def2")
			},
			mkArgs("val1"),
			&testOpts{A1: "val1", A2: "def2"},
			nil,
		},
		{
			"required_provided_required_variadic_1_provided",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.RequiredVariadic(&opts.V1, "name2", 1)
			},
			mkArgs("val1", "val2"),
			&testOpts{A1: "val1", V1: mkArgs("val2")},
			nil,
		},
		{
			"required_provided_required_variadic_2_provided",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.RequiredVariadic(&opts.V1, "name2", 2)
			},
			mkArgs("val1", "val2", "val3"),
			&testOpts{A1: "val1", V1: mkArgs("val2", "val3")},
			nil,
		},
		{
			"required_provided_required_variadic_1_missing",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.RequiredVariadic(&opts.V1, "name2", 1)
			},
			mkArgs("val1"),
			nil,
			errors.New("required name2 argument(s) missing; you must supply at least 1"),
		},
		{
			"required_provided_required_variadic_2_missing_all",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.RequiredVariadic(&opts.V1, "name2", 2)
			},
			mkArgs("val1"),
			nil,
			errors.New("required name2 argument(s) missing; you must supply at least 2"),
		},
		{
			"required_provided_required_variadic_2_missing_one",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.RequiredVariadic(&opts.V1, "name2", 2)
			},
			mkArgs("val1", "val2"),
			nil,
			errors.New("required name2 argument(s) missing; you must supply at least 2"),
		},
		{
			"required_provided_optional_variadic_provided",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.OptionalVariadic(&opts.V1, "name2", "def1", "def2")
			},
			mkArgs("val1", "val2", "val3"),
			&testOpts{A1: "val1", V1: mkArgs("val2", "val3")},
			nil,
		},
		{
			"required_provided_optional_variadic_missing",
			func(opts *testOpts, al *ArgList) {
				al.Required(&opts.A1, "name1")
				al.OptionalVariadic(&opts.V1, "name2", "def1", "def2")
			},
			mkArgs("val1"),
			&testOpts{A1: "val1", V1: mkArgs("def1", "def2")},
			nil,
		},
	}

	for _, c := range cases {
		desc, setup, args, want, wantErr := c.desc, c.setup, c.args, c.want, c.wantErr
		t.Run(desc, func(t *testing.T) {
			argList := new(ArgList)
			opts := new(testOpts)
			setup(opts, argList)
			if gotErr := argList.parseArgs(args); gotErr != nil {
				if c.wantErr == nil {
					t.Fatalf("got unexpected error %q", gotErr)
				}
				got, want := gotErr.Error(), wantErr.Error()
				if got != want {
					t.Fatalf("got error %q; want %q", got, want)
				}
				return
			}
			if wantErr != nil {
				t.Fatalf("got nil error; want %s", c.wantErr)
			}
			assert.Equal(t, opts, want)
		})
	}

}
