package scopes

import (
	go_check "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestMyScope(t *testing.T) { go_check.TestingT(t) }

type TestMyScopeSuite struct{}

var _ = go_check.Suite(&TestMyScopeSuite{})

func (s *TestMyScopeSuite) TestMyScope(c *go_check.C) {
	metadata := NewSearchMetadata("us", "phone")

	// basic check
	c.Assert(metadata.Locale(), go_check.Equals, "us")
	c.Assert(metadata.FormFactor(), go_check.Equals, "phone")
	c.Assert(metadata.Cardinality(), go_check.Equals, 0)
	c.Assert(metadata.Location() == nil, go_check.Equals, true)

	// test finishing metadata
	finalizeSearchMetadata(metadata)
	c.Assert(metadata.m == nil, go_check.Equals, true)
}
