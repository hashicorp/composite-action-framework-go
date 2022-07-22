package cli

import (
	"flag"
	"testing"
)

type actualFlags struct {
	val string
}

func (tf *actualFlags) Flags(fs *flag.FlagSet) {
	fs.StringVar(&tf.val, "test", "", "test flag")
}

type hasInnerFlagsPtr struct {
	Inner *actualFlags
}

type hasInnerFlagsValue struct {
	Inner actualFlags
}

type actualAndInnerFlags struct {
	Inner actualFlags
	Val   string
}

func (aif *actualAndInnerFlags) Flags(fs *flag.FlagSet) {
	fs.StringVar(&aif.Val, "actual", "", "actual test flag")
}

func TestMakeFlags_actual(t *testing.T) {
	tf := &actualFlags{}
	f := makeFlags(tf)
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	f.Flags(fs)
	if err := fs.Parse([]string{"-test=hi"}); err != nil {
		t.Fatal(err)
	}
	if tf.val != "hi" {
		t.Errorf("got tf.val=%q; want %q", tf.val, "hi")
	}
}

func TestMakeFlags_innerPtr(t *testing.T) {
	tf := &hasInnerFlagsPtr{}
	f := makeFlags(tf)
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	f.Flags(fs)
	if err := fs.Parse([]string{"-test=hi"}); err != nil {
		t.Fatal(err)
	}
	if tf.Inner.val != "hi" {
		t.Errorf("got tf.inner.val=%q; want %q", tf.Inner.val, "hi")
	}
}

func TestMakeFlags_innerVal(t *testing.T) {
	tf := &hasInnerFlagsPtr{}
	f := makeFlags(tf)
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	f.Flags(fs)
	if err := fs.Parse([]string{"-test=hi"}); err != nil {
		t.Fatal(err)
	}
	if tf.Inner.val != "hi" {
		t.Errorf("got tf.inner.val=%q; want %q", tf.Inner.val, "hi")
	}
}

func TestMakeFlags_innerAndActual(t *testing.T) {
	tf := &actualAndInnerFlags{}
	f := makeFlags(tf)
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	f.Flags(fs)
	if err := fs.Parse([]string{"-test=hi", "-actual=hello"}); err != nil {
		t.Fatal(err)
	}
	if tf.Inner.val != "hi" {
		t.Errorf("got tf.inner.al=%q; want %q", tf.Inner.val, "hi")
	}
	if tf.Val != "hello" {
		t.Errorf("got tf.val=%q; want %q", tf.Val, "hello")
	}
}
