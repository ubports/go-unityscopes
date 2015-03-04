package main

import (
	"launchpad.net/go-unityscopes/v1"
	"log"
)

const searchCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "grid",
    "card-size": "small"
  },
  "components": {
    "title": "title",
    "art":  "art",
    "subtitle": "username"
  }
}`

// SCOPE ***********************************************************************

var scope_interface scopes.Scope

type MyScope struct {
	base *scopes.ScopeBase
}

func (s *MyScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	widget := scopes.NewPreviewWidget("foo", "text")
	widget.AddAttributeValue("text", "Hello")
	reply.PushWidgets(widget)
	return nil
}

func (s *MyScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	root_department := s.CreateDepartments(query, metadata, reply)
	reply.RegisterDepartments(root_department)
	
	return s.AddQueryResults(reply, query.QueryString())
}

func (s *MyScope) SetScopeBase(base *scopes.ScopeBase) {
	s.base = base
}

// RESULTS *********************************************************************

func (s *MyScope) AddQueryResults(reply *scopes.SearchReply, query string) error {
	cat := reply.RegisterCategory("category", "Category", "", searchCategoryTemplate)

	result := scopes.NewCategorisedResult(cat)
	result.SetURI("http://localhost/" + query)
	result.SetDndURI("http://localhost_dnduri" + query)
	result.SetTitle("TEST" + query)
	result.SetArt("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
	result.Set("test_value_bool", true)
	result.Set("test_value_string", "test_value"+query)
	result.Set("test_value_int", 1999)
	result.Set("test_value_float", 1.999)
	if err := reply.Push(result); err != nil {
		return err
	}

	result.SetURI("http://localhost2/" + query)
	result.SetDndURI("http://localhost_dnduri2" + query)
	result.SetTitle("TEST2")
	result.SetArt("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
	result.Set("test_value_bool", false)
	result.Set("test_value_string", "test_value2"+query)
	result.Set("test_value_int", 2000)
	result.Set("test_value_float", 2.100)

	// add a variant map value
	m := make(map[string]interface{})
	m["value1"] = 1
	m["value2"] = "string_value"
	result.Set("test_value_map", m)

	// add a variant array value
	l := make([]interface{}, 0)
	l = append(l, 1999)
	l = append(l, "string_value")
	result.Set("test_value_array", l)
	if err := reply.Push(result); err != nil {
		return err
	}

	return nil
}

// DEPARTMENTS *****************************************************************

func SearchDepartment(root *scopes.Department, id string) *scopes.Department {
	sub_depts := root.Subdepartments()
	for _, element := range sub_depts {
		if element.Id() == id {
			return element
		}
	}
	return nil
}

func (s *MyScope) GetRockSubdepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {
	active_dep, err := scopes.NewDepartment("Rock", query, "Rock Music")
	if err == nil {
		active_dep.SetAlternateLabel("Rock Music Alt")
		department, _ := scopes.NewDepartment("60s", query, "Rock from the 60s")
		active_dep.AddSubdepartment(department)

		department2, _ := scopes.NewDepartment("70s", query, "Rock from the 70s")
		active_dep.AddSubdepartment(department2)
	}

	return active_dep
}

func (s *MyScope) GetSoulSubdepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {
	active_dep, err := scopes.NewDepartment("Soul", query, "Soul Music")
	if err == nil {
		active_dep.SetAlternateLabel("Soul Music Alt")
		department, _ := scopes.NewDepartment("Motown", query, "Motown Soul")
		active_dep.AddSubdepartment(department)

		department2, _ := scopes.NewDepartment("New Soul", query, "New Soul")
		active_dep.AddSubdepartment(department2)
	}

	return active_dep
}

func (s *MyScope) CreateDepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {
	department, _ := scopes.NewDepartment("", query, "Browse Music")
	department.SetAlternateLabel("Browse Music Alt")

	rock_dept := s.GetRockSubdepartments(query, metadata, reply)
	if rock_dept != nil {
		department.AddSubdepartment(rock_dept)
	}

	soul_dept := s.GetSoulSubdepartments(query, metadata, reply)
	if soul_dept != nil {
		department.AddSubdepartment(soul_dept)
	}

	return department
}

// MAIN ************************************************************************

func main() {
	var sc MyScope
	scope_interface = &sc

	if err := scopes.Run(&MyScope{}); err != nil {
		log.Fatalln(err)
	}
}
