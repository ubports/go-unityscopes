package scopes

import (
	go_check "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestCannedQuery(t *testing.T) { go_check.TestingT(t) }

type TestCannedQuerySuite struct{}

var _ = go_check.Suite(&TestCannedQuerySuite{})

func (s *TestCannedQuerySuite) TestQuery(c *go_check.C) {
	query := NewCannedQuery("scope", "query_string", "department_string")

	// basic check
	c.Assert(query.ScopeID(), go_check.Equals, "scope")
	c.Assert(query.DepartmentID(), go_check.Equals, "department_string")
	c.Assert(query.QueryString(), go_check.Equals, "query_string")
	c.Assert(query.q != nil, go_check.Equals, true)

	// verify uri
	c.Assert(query.ToURI(), go_check.Equals, "scope://scope?q=query%5Fstring&dep=department%5Fstring")

	// check setters
	query.SetDepartmentID("department_id")
	c.Assert(query.DepartmentID(), go_check.Equals, "department_id")

	query.SetQueryString("new_query_value")
	c.Assert(query.QueryString(), go_check.Equals, "new_query_value")

	// TODO FilterState setter

	// test finishing query
	finalizeCannedQuery(query)
	c.Assert(query.q == nil, go_check.Equals, true)
}
