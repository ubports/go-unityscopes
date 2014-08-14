package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"runtime"
	"unsafe"
)

// SearchMetadata holds additional metadata about the search request.
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

// Locale returns the expected locale for the search request.
func (metadata *SearchMetadata) Locale() string {
	locale := C.search_metadata_get_locale(metadata.m)
	defer C.free(unsafe.Pointer(locale))
	return C.GoString(locale)
}

// FormFactor returns the form factor for the search request.
func (metadata *SearchMetadata) FormFactor() string {
	formFactor := C.search_metadata_get_form_factor(metadata.m)
	defer C.free(unsafe.Pointer(formFactor))
	return C.GoString(formFactor)
}

// Cardinality returns the desired number of results for the search request.
func (metadata *SearchMetadata) Cardinality() int {
	return int(C.search_metadata_get_cardinality(metadata.m))
}

type Location struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Altitude           float64 `json:"altitude"`
	AreaCode           string  `json:"area_code"`
	City               string  `json:"city"`
	CountryCode        string  `json:"country_code"`
	CountryName        string  `json:"country_name"`
	HorizontalAccuracy float64 `json:"horizontal_accuracy"`
	VerticalAccuracy   float64 `json:"vertical_accuracy"`
	RegionCode         string  `json:"region_code"`
	RegionName         string  `json:"region_name"`
	ZipPostalCode      string  `json:"zip_postal_code"`
}

func (metadata *SearchMetadata) Location() *Location {
	locData := C.search_metadata_get_location(metadata.m)
	if locData == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(locData))
	locString := C.GoString(locData)
	var location Location
	if err := json.Unmarshal([]byte(locString), &location); err != nil {
		panic(err)
	}
	return &location
}

// ActionMetadata holds additional metadata about the preview request
// or result activation.
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

// Locale returns the expected locale for the preview or result activation.
func (metadata *ActionMetadata) Locale() string {
	locale := C.action_metadata_get_locale(metadata.m)
	defer C.free(unsafe.Pointer(locale))
	return C.GoString(locale)
}

// FormFactor returns the form factor for the preview or result activation.
func (metadata *ActionMetadata) FormFactor() string {
	formFactor := C.action_metadata_get_form_factor(metadata.m)
	defer C.free(unsafe.Pointer(formFactor))
	return C.GoString(formFactor)
}
