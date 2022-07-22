package cli

import (
	"fmt"
	"testing"
)

func TestMultiplexInterface_doer1(t *testing.T) {
	thing := &doer1{}
	multi := multiplexInterface[Doer](thing, &multiDoer{})

	multi.Do()
	t.Fail()
}

func TestMultiplexInterface_compoundDoer(t *testing.T) {
	thing := &compoundDoer{}
	multi := multiplexInterface[Doer](thing, &multiDoer{})

	multi.Do()
	t.Fail()
}

type Doer interface {
	Do()
}

type multiDoer []Doer

func (md *multiDoer) Collect(d Doer) {
	*md = append(*md, d)
}

func (md *multiDoer) Do() {
	for _, d := range *md {
		d.Do()
	}
}

type doer1 struct{}

func (d1 *doer1) Do() { fmt.Println("doer1") }

type doer2 struct{}

func (d2 *doer2) Do() { fmt.Println("doer2") }

type compoundDoer struct {
	Doer1 doer1
	Doer2 doer2
	name  string
}

func (cd *compoundDoer) Do() { fmt.Println("compoundDoer") }
