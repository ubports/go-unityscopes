package scopes_test

import (
	"testing"
	. "gopkg.in/check.v1"
)

type S struct {}

func init() {
	Suite(&S{})
}

func TestAll(t *testing.T) {
	TestingT(t)
}
