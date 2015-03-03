package scopes_test

import (
	. "gopkg.in/check.v1"
	"testing"
)

type S struct{}

func init() {
	Suite(&S{})
}

func TestAll(t *testing.T) {
	TestingT(t)
}
