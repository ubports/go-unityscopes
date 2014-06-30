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
func NewDepartment(departmentID string, query *CannedQuery, label string) (*Department, error) {
	dept := makeDepartment()
	var errorString *C.char = nil
	C.new_department(unsafe.Pointer(&departmentID), query.q, unsafe.Pointer(&label), &dept.d[0], &errorString)
	if err := checkError(errorString); err != nil {
		return nil, err
	}
	return dept, nil
}

func finalizeDepartment(dept *Department) {
	C.destroy_department_ptr(&dept.d[0])
}

// AddSubdepartment adds a new child department to this department.
func (dept *Department) AddSubdepartment(child *Department) {
	C.department_add_subdepartment(&dept.d[0], &child.d[0])
}

// SetAlternateLabel sets the alternate label for this department.
//
// This should express the plural form of the normal label.  For
// example, if the normal label is "Books", then the alternate label
// should be "All Books".
//
// The alternate label only needs to be provided for the current
// department.
func (dept *Department) SetAlternateLabel(label string) {
	C.department_set_alternate_label(&dept.d[0], unsafe.Pointer(&label))
}

// SetHasSubdepartments sets whether this department has subdepartments.
//
// It is not necessary to call this if AddSubdepartment has been
// called.  It intended for cases where subdepartments have not been
// specified but the shell should still act as if it has them.
func (dept *Department) SetHasSubdepartments(subdepartments bool) {
	var cSubdepts C.int
	if subdepartments {
		cSubdepts = 1
	} else {
		cSubdepts = 0
	}
	C.department_set_has_subdepartments(&dept.d[0], cSubdepts)
}
