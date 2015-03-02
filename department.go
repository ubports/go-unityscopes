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
	d     C.SharedPtrData
	query *CannedQuery
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
	dept.query = query
	return dept, nil
}

// NewDepartment creates a new department using the given canned query.
func NewDepartmentWithoutId(query *CannedQuery, label string) (*Department, error) {
	dept := makeDepartment()
	var errorString *C.char = nil
	C.new_department(nil, query.q, unsafe.Pointer(&label), &dept.d[0], &errorString)
	if err := checkError(errorString); err != nil {
		return nil, err
	}
	dept.query = query
	return dept, nil
}

func finalizeDepartment(dept *Department) {
	C.destroy_department_ptr(&dept.d[0])
}

// AddSubdepartment adds a new child department to this department.
func (dept *Department) AddSubdepartment(child *Department) {
	C.department_add_subdepartment(&dept.d[0], &child.d[0])
}

// Get the identifier of this department.
func (dept *Department) Id() string {
	s := C.department_get_id(&dept.d[0])
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// Get the label of this department.
func (dept *Department) Label() string {
	s := C.department_get_label(&dept.d[0])
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// Get the canned query associated with this department.
func (dept *Department) Query() *CannedQuery {
	return dept.query
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

// AlternateLabel sets gets the alternate label for this department.
//
// This should express the plural form of the normal label.  For
// example, if the normal label is "Books", then the alternate label
// should be "All Books".
func (dept *Department) AlternateLabel() string {
	s := C.department_get_alternate_label(&dept.d[0])
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
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

// Check if this department has subdepartments or has_subdepartments flag is set
func (dept *Department) HasSubdepartments() bool {
	if C.department_has_subdepartments(&dept.d[0]) == 1 {
		return true
	} else {
		return false
	}
}

// Get list of sub-departments of this department.
func (dept *Department) Subdepartments() []*Department {
	var nb_subdepartments C.int
	var theCArray *C.SharedPtrData = C.department_get_subdepartments(&dept.d[0], &nb_subdepartments)
	defer C.free(unsafe.Pointer(theCArray))
	length := int(nb_subdepartments)
	// create a very big slice and then slice it to the number of subdepartments
	slice := (*[1 << 30]C.SharedPtrData)(unsafe.Pointer(theCArray))[:length:length]
	ptr_depts := make([]*Department, length)
	for i := 0; i < length; i++ {
		ptr_depts[i] = makeDepartment()
		ptr_depts[i].d = slice[i]
	}
	return ptr_depts
}

// Set sub-departments of this department.
func (dept *Department) SetSubdepartments(subdepartments []*Department) {
	api_depts := make([]*C.SharedPtrData, len(subdepartments))
	for i := 0; i < len(subdepartments); i++ {
		api_depts[i] = (*C.SharedPtrData)(unsafe.Pointer(&subdepartments[i].d))
	}
	if len(subdepartments) > 0 {
		C.department_set_subdepartments(&dept.d[0], (**C.SharedPtrData)(unsafe.Pointer(&api_depts[0])), C.int(len(subdepartments)))
	} else {
		C.department_set_subdepartments(&dept.d[0], (**C.SharedPtrData)(unsafe.Pointer(nil)), C.int(len(subdepartments)))
	}
}
