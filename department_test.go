package scopes

import (
	go_check "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestDepartment(t *testing.T) { go_check.TestingT(t) }

type TestDepartmentSuite struct{}

var _ = go_check.Suite(&TestDepartmentSuite{})

func (s *TestDepartmentSuite) TestDepartment(c *go_check.C) {
	query := NewCannedQuery("scope", "query_string", "department_string")
	department, err := NewDepartment("department_string2", query, "TEST_DEPARTMENT")
	c.Assert(err == nil, go_check.Equals, true)

	department.SetAlternateLabel("test_alternate_label")
	c.Assert(department.AlternateLabel(), go_check.Equals, "test_alternate_label")
	c.Assert(department.Id(), go_check.Equals, "department_string2")
	c.Assert(department.Label(), go_check.Equals, "TEST_DEPARTMENT")

	department.SetHasSubdepartments(true)
	c.Assert(department.HasSubdepartments(), go_check.Equals, true)

	department.SetHasSubdepartments(false)
	c.Assert(department.HasSubdepartments(), go_check.Equals, false)

	department2, err2 := NewDepartment("sub_department_string", query, "TEST_SUB_DEPARTMENT")
	c.Assert(err2 == nil, go_check.Equals, true)
	department2.SetAlternateLabel("test_alternate_label_2")

	department3, err3 := NewDepartment("sub_department_2_string", query, "TEST_SUB_DEPARTMENT_2")
	c.Assert(err3 == nil, go_check.Equals, true)
	department3.SetAlternateLabel("test_alternate_label_3")

	subdepartments := department.Subdepartments()
	c.Assert(len(subdepartments), go_check.Equals, 0)
	c.Assert(department.HasSubdepartments(), go_check.Equals, false)

	department.AddSubdepartment(department2)
	department.AddSubdepartment(department3)

	subdepartments = department.Subdepartments()
	c.Assert(len(subdepartments), go_check.Equals, 2)
	c.Assert(department.HasSubdepartments(), go_check.Equals, true)

	// verify that the values are correct in all subdepartments
	c.Assert(subdepartments[0].Id(), go_check.Equals, "sub_department_string")
	c.Assert(subdepartments[0].Label(), go_check.Equals, "TEST_SUB_DEPARTMENT")
	c.Assert(subdepartments[0].AlternateLabel(), go_check.Equals, "test_alternate_label_2")
	c.Assert(subdepartments[1].Id(), go_check.Equals, "sub_department_2_string")
	c.Assert(subdepartments[1].Label(), go_check.Equals, "TEST_SUB_DEPARTMENT_2")
	c.Assert(subdepartments[1].AlternateLabel(), go_check.Equals, "test_alternate_label_3")

	sub_depts := make([]*Department, 0)
	department.SetSubdepartments(sub_depts)

	subdepartments = department.Subdepartments()
	c.Assert(len(subdepartments), go_check.Equals, 0)
	c.Assert(department.HasSubdepartments(), go_check.Equals, false)

	sub_depts = append(sub_depts, department2)
	sub_depts = append(sub_depts, department3)

	department.SetSubdepartments(sub_depts)
	subdepartments = department.Subdepartments()
	c.Assert(len(subdepartments), go_check.Equals, 2)
	c.Assert(department.HasSubdepartments(), go_check.Equals, true)

	c.Assert(subdepartments[0].Id(), go_check.Equals, "sub_department_string")
	c.Assert(subdepartments[0].Label(), go_check.Equals, "TEST_SUB_DEPARTMENT")
	c.Assert(subdepartments[0].AlternateLabel(), go_check.Equals, "test_alternate_label_2")
	c.Assert(subdepartments[1].Id(), go_check.Equals, "sub_department_2_string")
	c.Assert(subdepartments[1].Label(), go_check.Equals, "TEST_SUB_DEPARTMENT_2")
	c.Assert(subdepartments[1].AlternateLabel(), go_check.Equals, "test_alternate_label_3")

	stored_query := department.Query()
	c.Assert(stored_query.ScopeID(), go_check.Equals, "scope")
	c.Assert(stored_query.DepartmentID(), go_check.Equals, "department_string")
	c.Assert(stored_query.QueryString(), go_check.Equals, "query_string")
	c.Assert(query.q != nil, go_check.Equals, true)

	// test finishing department
	finalizeDepartment(department)
	c.Assert(department.d[0] == 0, go_check.Equals, true)
	c.Assert(department.d[1] == 0, go_check.Equals, true)
}

func (s *TestDepartmentSuite) TestDepartmentDifferentCreation(c *go_check.C) {
	query := NewCannedQuery("scope", "query_string", "department_string")
	department, err := NewDepartmentWithoutId(query, "TEST_DEPARTMENT")

	c.Assert(err == nil, go_check.Equals, true)
	c.Assert(department.Id(), go_check.Equals, "department_string")
	c.Assert(department.Label(), go_check.Equals, "TEST_DEPARTMENT")

	// test finishing department
	finalizeDepartment(department)
	c.Assert(department.d[0] == 0, go_check.Equals, true)
	c.Assert(department.d[1] == 0, go_check.Equals, true)
}

func (s *TestDepartmentSuite) TestDepartmentEmptyLabel(c *go_check.C) {
	query := NewCannedQuery("scope", "query_string", "department_string")
	department, err := NewDepartmentWithoutId(query, "")
	c.Assert(err != nil, go_check.Equals, true)
	c.Assert(department == nil, go_check.Equals, true)

	department2, err2 := NewDepartment("dept_id", query, "")
	c.Assert(err2 != nil, go_check.Equals, true)
	c.Assert(department2 == nil, go_check.Equals, true)
}
