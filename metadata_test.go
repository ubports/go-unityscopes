package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v1"
)

func (s *S) TestMyScope(c *C) {
	metadata := scopes.NewSearchMetadata(2, "us", "phone")

	// basic check
	c.Check(metadata.Locale(), Equals, "us")
	c.Check(metadata.FormFactor(), Equals, "phone")
	c.Check(metadata.Cardinality(), Equals, 2)
	c.Check(metadata.Location(), IsNil)
}
