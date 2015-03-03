package main

import (
	"log"
	"scopes"
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

var scope_interface scopes.Scope

func SearchDepartment(root *scopes.Department, id string) *scopes.Department {
	sub_depts := root.Subdepartments()
	for _,element := range sub_depts {
		if element.Id() == id {
			return element
		}
	}
	return nil
}

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
	root_department := s.CreateMainDepartments(query, metadata, reply)

	if query.DepartmentID() == "Rock" {
//		root_department = s.AddRockSubdepartments(query, metadata, reply, root_department)
	} else if query.DepartmentID() == "Soul" {
//		root_department = s.AddSoulSubdepartments(query, metadata, reply, root_department)
//		active_dep := SearchDepartment(root_department, "Soul")
//		if active_dep != nil {
//			department, _ := scopes.NewDepartment("Motown", query, "Motown Soul")
//			active_dep.AddSubdepartment(department)
//			
//			department2, _ := scopes.NewDepartment("NewSoul", query, "New Soul")
//			active_dep.AddSubdepartment(department2)
//		}
	}
	
	reply.RegisterDepartments(root_department)
	
	return s.AddEmptyQueryResults(reply)
}

func (s *MyScope) SetScopeBase(base *scopes.ScopeBase) {
	s.base = base
}

func (s *MyScope) CreateMainDepartments(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply) *scopes.Department {
	department, _ := scopes.NewDepartment("", query, "Browse Music")
	department.SetAlternateLabel("Browse Music Alt")
	
	// ROCK MUSIC
	department2, _ := scopes.NewDepartment("Rock", query, "Rock Music")
	department2.SetAlternateLabel("Rock Music Alt")
	// add Rock subdepartments
	rock1, _ := scopes.NewDepartment("60s", query, "Rock from the 60s")
	department2.AddSubdepartment(rock1)
	
	rock2, _ := scopes.NewDepartment("70s", query, "Rock from the 70s")
	department2.AddSubdepartment(rock2)

	// SOUL MUSIC
	department3, _ := scopes.NewDepartment("Soul", query, "Soul Music")
	department3.SetAlternateLabel("Soul Music Alt")
	
	// add Soul subdepartments
	soul1, _ := scopes.NewDepartment("Motown", query, "Motown Soul")
	department3.AddSubdepartment(soul1)
	
	soul2, _ := scopes.NewDepartment("NewSoul", query, "New Soul")
	department3.AddSubdepartment(soul2)
	
	// Add top subdepartments
	department.AddSubdepartment(department2)
	department.AddSubdepartment(department3)
	
	return department
}

func (s *MyScope) AddRockSubdepartments(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, root_dept *scopes.Department) *scopes.Department {
	active_dep := SearchDepartment(root_dept, "Rock")
	if active_dep != nil {
		department, _ := scopes.NewDepartment("60s", query, "Rock from the 60s")
		active_dep.AddSubdepartment(department)
		
		department2, _ := scopes.NewDepartment("70s", query, "Rock from the 70s")
		active_dep.AddSubdepartment(department2)
	}
	
	return root_dept
}

func (s *MyScope) AddSoulSubdepartments(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, root_dept *scopes.Department) *scopes.Department {
	active_dep := SearchDepartment(root_dept, "Soul")
	if active_dep != nil {
		department, _ := scopes.NewDepartment("Motown", query, "Motown Soul")
		active_dep.AddSubdepartment(department)
		
		department2, _ := scopes.NewDepartment("NewSoul", query, "New Soul")
		active_dep.AddSubdepartment(department2)
	}
	
	return root_dept
}

func (s *MyScope) AddEmptyQueryResults(reply *scopes.SearchReply) error {
	cat := reply.RegisterCategory("category", "Category", "", searchCategoryTemplate)

	result := scopes.NewCategorisedResult(cat)
	result.SetURI("http://localhost/")
	result.SetDndURI("http://localhost_dnduri")
	result.SetTitle("TEST")
	result.SetArt("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
	result.Set("test_value_bool", true)
	result.Set("test_value_string", "test_value")
	result.Set("test_value_int", 1999)
	result.Set("test_value_float", 1.999)
	if err := reply.Push(result); err != nil {
		return err
	}

	result.SetURI("http://localhost2/")
	result.SetDndURI("http://localhost_dnduri2")
	result.SetTitle("TEST2")
	result.SetArt("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
	result.Set("test_value_bool", false)
	result.Set("test_value_string", "test_value2")
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

func (s *MyScope) AddRockResults(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, root_dept *scopes.Department) error {
	
	cat := reply.RegisterCategory("category", "Category", "", searchCategoryTemplate)

	result := scopes.NewCategorisedResult(cat)
	result.SetURI("http://localhost/")
	result.SetDndURI("http://localhost_dnduri")
	result.SetTitle("TEST")
	result.SetArt("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
	result.Set("test_value_bool", true)
	result.Set("test_value_string", "test_value")
	result.Set("test_value_int", 1999)
	result.Set("test_value_float", 1.999)
	if err := reply.Push(result); err != nil {
		return err
	}

	result.SetURI("http://localhost2/")
	result.SetDndURI("http://localhost_dnduri2")
	result.SetTitle("TEST2")
	result.SetArt("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
	result.Set("test_value_bool", false)
	result.Set("test_value_string", "test_value2")
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

func main() {
	var sc MyScope
	scope_interface = &sc

	if err := scopes.Run(&MyScope{}); err != nil {
		log.Fatalln(err)
	}
}
