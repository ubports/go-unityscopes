package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v1"
)

func (s *S) TestRatingFilter(c *C) {
	filter1 := scopes.NewRatingFilter("f1", "Options")
	c.Check("f1", Equals, filter1.Id)
	c.Check("Options", Equals, filter1.Label)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)

	filter1.DisplayHints = scopes.FilterDisplayPrimary
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")

	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayPrimary)
	c.Check(2, Equals, len(filter1.Options))
	c.Check("1", Equals, filter1.Options[0].Id)
	c.Check("Option 1", Equals, filter1.Options[0].Label)
	c.Check("2", Equals, filter1.Options[1].Id)
	c.Check("Option 2", Equals, filter1.Options[1].Label)

	// verify the list of options
	c.Check(len(filter1.Options), Equals, 2)
	c.Check(filter1.Options, DeepEquals, []scopes.FilterOption{scopes.FilterOption{"1", "Option 1"}, scopes.FilterOption{"2", "Option 2"}})
}

func (s *S) TestRatingFilterMultiSelection(c *C) {
	filter1 := scopes.NewRatingFilter("f1", "Options")
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 3")

	fstate := make(scopes.FilterState)

	// enable option1 & option2
	filter1.UpdateState(fstate, "1", true)
	filter1.UpdateState(fstate, "2", true)
	_, ok := fstate["f1"]
	c.Check(ok, Equals, true)

	c.Check(filter1.HasActiveRating(fstate), Equals, true)
	active := filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 2)
	c.Check(active, DeepEquals, []interface{}{"1", "2"})

	// disable option1
	filter1.UpdateState(fstate, "1", false)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 1)
	c.Check(active[0], Equals, "2")

	// disable option2
	filter1.UpdateState(fstate, "2", false)
	c.Check(0, Equals, len(filter1.ActiveOptions(fstate)))

	filter1.UpdateState(fstate, "3", true)
	filter1.UpdateState(fstate, "1", true)

	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 2)
	c.Check(active, DeepEquals, []interface{}{"1", "3"})

	// add existing item
	filter1.UpdateState(fstate, "1", true)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 2)
	c.Check(active, DeepEquals, []interface{}{"1", "3"})

	// add in the middle
	filter1.UpdateState(fstate, "2", true)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 3)
	c.Check(active, DeepEquals, []interface{}{"1", "2", "3"})

	// erase in the middle
	filter1.UpdateState(fstate, "2", false)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 2)
	c.Check(active, DeepEquals, []interface{}{"1", "3"})

	filter1.UpdateState(fstate, "2", true)

	// erase at the beginning
	filter1.UpdateState(fstate, "1", false)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 2)
	c.Check(active, DeepEquals, []interface{}{"2", "3"})

	filter1.UpdateState(fstate, "1", true)

	// erase at the end
	filter1.UpdateState(fstate, "3", false)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 2)
	c.Check(active, DeepEquals, []interface{}{"1", "2"})

	filter1.UpdateState(fstate, "1", false)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 1)
	c.Check(active, DeepEquals, []interface{}{"2"})

	filter1.UpdateState(fstate, "2", false)
	active = filter1.ActiveRating(fstate)
	c.Check(len(active), Equals, 0)
}

func (s *S) TestRatingFilterBadOption(c *C) {
	filter1 := scopes.NewRatingFilter("f1", "Options")
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 3")

	fstate := make(scopes.FilterState)

	c.Assert(func() { filter1.UpdateState(fstate, "5", true) }, PanicMatches, "invalid option ID")
	c.Assert(func() { filter1.UpdateState(fstate, "5", false) }, PanicMatches, "invalid option ID")
}
