package scopes_test

import (
	"errors"
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v1"
)

type unserializable struct{}

func (u *unserializable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Can not marshal to JSON")
}

func (u *unserializable) UnmarshalJSON(data []byte) error {
	return errors.New("Can not unmarshal from JSON")
}

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

	// try to pass a non-pointer object
	var errorJsonUnserialize unserializable
	err = metadata.ScopeData(errorJsonUnserialize)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "json: Unmarshal(non-pointer scopes_test.unserializable)")

	// try to use an unserializable object
	// We should get an error
	err = metadata.ScopeData(&errorJsonUnserialize)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "Can not unmarshal from JSON")
}

func (s *S) TestActionMetadataHints(c *C) {
	metadata := scopes.NewActionMetadata("us", "phone")

	var value interface{}

	// we still have no hints
	err := metadata.Hints(&value)
	c.Check(err, IsNil)
	c.Check(value, DeepEquals, map[string]interface{}{})

	err = metadata.SetHint("test_1", "value_1")
	c.Check(err, IsNil)

	err = metadata.GetHint("test_1", &value)
	c.Check(err, IsNil)
	c.Check(value, Equals, "value_1")

	err = metadata.GetHint("test_1_not_exists", &value)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "ActionMetadata:GetHint() value not found for key [test_1_not_exists]")

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

	// pass non-pointer
	var errorJsonUnserialize unserializable
	err = metadata.Hints(errorJsonUnserialize)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "json: Unmarshal(non-pointer scopes_test.unserializable)")

	// pass non-serializable object
	err = metadata.Hints(&errorJsonUnserialize)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "Can not unmarshal from JSON")

	err = metadata.SetHint("bad_hint", &errorJsonUnserialize)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "json: error calling MarshalJSON for type *scopes_test.unserializable: Can not marshal to JSON")
}
