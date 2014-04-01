package unityscope

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// Result represents a result from the scope
type Result struct {
	result unsafe.Pointer
}

func finalizeResult(res *Result) {
	if res.result != nil {
		C.destroy_result(res.result)
	}
	res.result = nil
}

// SetURI sets the "uri" attribute of the result.
func (res *Result) SetURI(uri string) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	C.result_set_uri(res.result, cUri)
}

// SetTitle sets the "title" attribute of the result.
func (res *Result) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.result_set_title(res.result, cTitle)
}

// SetArt sets the "art" attribute of the result.
func (res *Result) SetArt(art string) {
	cArt := C.CString(art)
	defer C.free(unsafe.Pointer(cArt))
	C.result_set_art(res.result, cArt)
}

// SetDndURI sets the "dnd_uri" attribute of the result.
func (res *Result) SetDndURI(uri string) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	C.result_set_dnd_uri(res.result, cUri)
}

// CategorisedResult represents a result linked to a particular category.
type CategorisedResult struct {
	Result
}

// NewCategorisedResult creates a new empty result linked to the given
// category.
func NewCategorisedResult(category *Category) *CategorisedResult {
	res := new(CategorisedResult)
	runtime.SetFinalizer(res, finalizeCategorisedResult)
	res.result = C.new_categorised_result(&category.c[0])
	return res
}

func finalizeCategorisedResult(res *CategorisedResult) {
	finalizeResult(&res.Result)
}
