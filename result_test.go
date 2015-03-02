package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v1"
)

func (s *S) TestResultSetURI(c *C) {
	r := scopes.NewTestingResult()
	r.SetURI("http://example.com")
	c.Check(r.URI(), Equals, "http://example.com")

	var uri string
	c.Check(r.Get("uri", &uri), IsNil)
	c.Check(uri, Equals, "http://example.com")
}
