package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// Department represents a section of the a scope's results.  A
// department can have sub-departments.
type Department struct {
	d C.SharedPtrData
}

func makeDepartment() *Department {
	dept := new(Department)
	runtime.SetFinalizer(dept, finalizeDepartment)
	return dept
}

// NewDepartment creates a new department using the given canned query.
func NewDepartment(query *CannedQuery, label string) (*Department, error) {
	dept := makeDepartment()
	var errorString *C.char = nil
	C.new_department(query.q, unsafe.Pointer(&label), &dept.d[0], &errorString)
	if err := checkError(errorString); err != nil {
		return nil, err
	}
	return dept, nil
}

func finalizeDepartment(dept *Department) {
	C.destroy_department_ptr(&dept.d[0])
}

func (dept *Department) AddSubdepartment(child *Department) {
	C.department_add_subdepartment(&dept.d[0], &child.d[0])
}

func (dept *Department) SetAlternateLabel(label string) {
	C.department_set_alternate_label(&dept.d[0], unsafe.Pointer(&label))
}

func (dept *Department) SetHasSubdepartments(subdepartments bool) {
	var cSubdepts C.int
	if subdepartments {
		cSubdepts = 1
	} else {
		cSubdepts = 0
	}
	C.department_set_has_subdepartments(&dept.d[0], cSubdepts)
}
