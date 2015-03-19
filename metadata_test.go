package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v1"
)

func (s *S) TestMetadataBasic(c *C) {
	metadata := scopes.NewSearchMetadata(2, "us", "phone")

	// basic check
	c.Check(metadata.Locale(), Equals, "us")
	c.Check(metadata.FormFactor(), Equals, "phone")
	c.Check(metadata.Cardinality(), Equals, 2)
	c.Check(metadata.Location(), IsNil)
}

func (s *S) TestSetLocation(c *C) {
	metadata := scopes.NewSearchMetadata(2, "us", "phone")
	location := scopes.Location{1.1, 2.1, 0.0, "EU", "Barcelona", "es", "Spain", 1.1, 1.1, "BCN", "BCN", "08080"}

	// basic check
	c.Check(metadata.Location(), IsNil)

	// set the location
	err := metadata.SetLocation(&location)
	c.Check(err, IsNil)

	stored_location := metadata.Location()
	c.Assert(stored_location, Not(Equals), nil)

	//	c.Check(stored_location, DeepEquals, &location)
}
