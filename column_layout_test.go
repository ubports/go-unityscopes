package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v1"
)

func (s *S) TestColumnLayout(c *C) {
	layout := scopes.NewColumnLayout(3)

	c.Check(layout.Size(), Equals, 0)
	c.Check(layout.NumberOfColumns(), Equals, 3)

	err := layout.AddColumn([]string{"column_1", "column_2"})
	c.Check(err, IsNil)

	c.Check(layout.Size(), Equals, 1)
	c.Check(layout.NumberOfColumns(), Equals, 3)

	col, err := layout.Column(0)
	c.Check(err, IsNil)

	c.Check(len(col), Equals, 2)
	c.Check(col[0], Equals, "column_1")
	c.Check(col[1], Equals, "column_2")

	// add another column
	err = layout.AddColumn([]string{"column_3", "column_4", "column_5"})
	c.Check(err, IsNil)

	col, err = layout.Column(1)
	c.Check(err, IsNil)

	c.Check(len(col), Equals, 3)
	c.Check(col[0], Equals, "column_3")
	c.Check(col[1], Equals, "column_4")
	c.Check(col[2], Equals, "column_5")

	// check for a bad column
	_, err = layout.Column(2)
	c.Check(err, Not(Equals), nil)

	// now add the last column
	err = layout.AddColumn([]string{"column_6"})
	c.Check(err, IsNil)

	col, err = layout.Column(2)
	c.Check(err, IsNil)

	c.Check(len(col), Equals, 1)
	c.Check(col[0], Equals, "column_6")

	// try to add more columns ... should obtain an error
	err = layout.AddColumn([]string{"column_3", "column_4", "column_5"})
	c.Check(err, Not(Equals), nil)

	// check size again
	c.Check(layout.Size(), Equals, 3)
	c.Check(layout.NumberOfColumns(), Equals, 3)
}
