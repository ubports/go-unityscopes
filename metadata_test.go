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

	c.Check(stored_location, DeepEquals, &location)
}

func metadata_test_bad_marshall(test int) int {
	return test
}

func (s *S) TestActionMetadata(c *C) {
	metadata := scopes.NewActionMetadata("us", "phone")

	// basic check
	c.Check(metadata.Locale(), Equals, "us")
	c.Check(metadata.FormFactor(), Equals, "phone")

	var scope_data interface{}
	metadata.ScopeData(&scope_data)
	c.Check(scope_data, IsNil)

	err := metadata.SetScopeData([]string{"test1", "test2", "test3"})
	c.Check(err, IsNil)

	err = metadata.ScopeData(&scope_data)
	c.Check(err, IsNil)
	c.Check(scope_data, DeepEquals, []interface{}{"test1", "test2", "test3"})

	err = metadata.ScopeData(metadata_test_bad_marshall)
	c.Check(err, Not(Equals), nil)

	err = metadata.SetScopeData(metadata_test_bad_marshall)
	c.Check(err, Not(Equals), nil)

	tuple1 := make(map[string]interface{})
	tuple1["id"] = "open"
	tuple1["label"] = "Open"
	tuple1["uri"] = "application:///tmp/non-existent.desktop"

	tuple2 := make(map[string]interface{})
	tuple2["id"] = "download"
	tuple2["label"] = "Download"

	tuple3 := make(map[string]interface{})
	tuple3["id"] = "hide"
	tuple3["label"] = "Hide"

	err = metadata.SetScopeData([]interface{}{tuple1, tuple2, tuple3})
	c.Check(err, IsNil)

	err = metadata.ScopeData(&scope_data)
	c.Check(err, IsNil)
	c.Check(scope_data, DeepEquals, []interface{}{tuple1, tuple2, tuple3})
}

func (s *S) TestHints(c *C) {
	metadata := scopes.NewActionMetadata("us", "phone")

	var value interface{}

	// we still have no hints
	err := metadata.Hints(&value)
	c.Check(err, IsNil)
	c.Check(value, DeepEquals, map[string]interface{}{})

	err = metadata.SetHint("test_1", "value_1")
	c.Check(err, IsNil)

	err, ok := metadata.GetHint("test_1", &value)
	c.Check(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(value, Equals, "value_1")

	err, ok = metadata.GetHint("test_1_not_exists", &value)
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)

	err = metadata.Hints(&value)
	expected_results := make(map[string]interface{})
	expected_results["test_1"] = "value_1"
	c.Check(expected_results, DeepEquals, value)

	err = metadata.SetHint("test_2", "value_2")
	c.Check(err, IsNil)

	expected_results["test_2"] = "value_2"
	err = metadata.Hints(&value)
	c.Check(err, IsNil)
	c.Check(expected_results, DeepEquals, value)

	err = metadata.SetHint("test_3", []interface{}{"value_3_1", "value_3_2"})
	c.Check(err, IsNil)

	expected_results["test_3"] = []interface{}{"value_3_1", "value_3_2"}
	err = metadata.Hints(&value)
	c.Check(err, IsNil)
	c.Check(expected_results, DeepEquals, value)

	err = metadata.Hints(metadata_test_bad_marshall)
	c.Check(err, Not(Equals), nil)

	err = metadata.SetHint("bad_hint", metadata_test_bad_marshall)
	c.Check(err, Not(Equals), nil)
}
