package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// SearchMetadata holds additional metadata about the search request
type SearchMetadata struct {
	m unsafe.Pointer
}

func finalizeSearchMetadata(metadata *SearchMetadata) {
	if metadata.m != nil {
		C.destroy_search_metadata(metadata.m)
	}
	metadata.m = nil
}

func makeSearchMetadata(m unsafe.Pointer) *SearchMetadata {
	metadata := new(SearchMetadata)
	runtime.SetFinalizer(metadata, finalizeSearchMetadata)
	metadata.m = m
	return metadata
}

func (metadata *SearchMetadata) Locale() string {
	locale := C.search_metadata_get_locale(metadata.m)
	defer C.free(unsafe.Pointer(locale))
	return C.GoString(locale)
}

func (metadata *SearchMetadata) FormFactor() string {
	formFactor := C.search_metadata_get_form_factor(metadata.m)
	defer C.free(unsafe.Pointer(formFactor))
	return C.GoString(formFactor)
}

func (metadata *SearchMetadata) Cardinality() int {
	return int(C.search_metadata_get_cardinality(metadata.m))
}

// ActionMetadata holds additional metadata about the preview request
// or result activation
type ActionMetadata struct {
	m unsafe.Pointer
}

func finalizeActionMetadata(metadata *ActionMetadata) {
	if metadata.m != nil {
		C.destroy_action_metadata(metadata.m)
	}
	metadata.m = nil
}

func makeActionMetadata(m unsafe.Pointer) *ActionMetadata {
	metadata := new(ActionMetadata)
	runtime.SetFinalizer(metadata, finalizeActionMetadata)
	metadata.m = m
	return metadata
}

func (metadata *ActionMetadata) Locale() string {
	locale := C.action_metadata_get_locale(metadata.m)
	defer C.free(unsafe.Pointer(locale))
	return C.GoString(locale)
}

func (metadata *ActionMetadata) FormFactor() string {
	formFactor := C.action_metadata_get_form_factor(metadata.m)
	defer C.free(unsafe.Pointer(formFactor))
	return C.GoString(formFactor)
}
