package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"runtime"
	"unsafe"
)

// SearchMetadata holds additional metadata about the search request.
type SearchMetadata struct {
	m *C._SearchMetadata
}

func finalizeSearchMetadata(metadata *SearchMetadata) {
	if metadata.m != nil {
		C.destroy_search_metadata(metadata.m)
	}
	metadata.m = nil
}

func makeSearchMetadata(m *C._SearchMetadata) *SearchMetadata {
	metadata := new(SearchMetadata)
	runtime.SetFinalizer(metadata, finalizeSearchMetadata)
	metadata.m = m
	return metadata
}

// NewSearchMetadata creates a new SearchMetadata with the given locale and
// form_factor
func NewSearchMetadata(cardinality int, locale, form_factor string) *SearchMetadata {
	return makeSearchMetadata(C.new_search_metadata(C.int(cardinality),
		unsafe.Pointer(&locale),
		unsafe.Pointer(&form_factor)))
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

// we use this type to reimplement the marshaller interface in order to make values
// like 1.0 not being converted as 1 (integer).
type Float float64

func (n Float) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%f", n)), nil
}

type Location struct {
	Latitude           Float  `json:"latitude"`
	Longitude          Float  `json:"longitude"`
	Altitude           Float  `json:"altitude"`
	AreaCode           string `json:"area_code"`
	City               string `json:"city"`
	CountryCode        string `json:"country_code"`
	CountryName        string `json:"country_name"`
	HorizontalAccuracy Float  `json:"horizontal_accuracy"`
	VerticalAccuracy   Float  `json:"vertical_accuracy"`
	RegionCode         string `json:"region_code"`
	RegionName         string `json:"region_name"`
	ZipPostalCode      string `json:"zip_postal_code"`
}

func (metadata *SearchMetadata) Location() *Location {
	var length C.int
	locData := C.search_metadata_get_location(metadata.m, &length)
	if locData == nil {
		return nil
	}
	defer C.free(locData)
	var location Location
	if err := json.Unmarshal(C.GoBytes(locData, length), &location); err != nil {
		panic(err)
	}
	return &location
}

// SetLocation sets the location
func (metadata *SearchMetadata) SetLocation(location *Location) error {
	data, err := json.Marshal(location)
	if err != nil {
		return err
	}
	json_value := string(data)
	var errorString *C.char = nil
	C.search_metadata_set_location(metadata.m, unsafe.Pointer(&json_value), &errorString)
	return checkError(errorString)
}

// ActionMetadata holds additional metadata about the preview request
// or result activation.
type ActionMetadata struct {
	m *C._ActionMetadata
}

func finalizeActionMetadata(metadata *ActionMetadata) {
	if metadata.m != nil {
		C.destroy_action_metadata(metadata.m)
	}
	metadata.m = nil
}

func makeActionMetadata(m *C._ActionMetadata) *ActionMetadata {
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

// ScopeData decodes the stored scope data into the given variable.
//
// Scope data is either set by the shell when calling a preview
// action, or set by the scope through an ActivationResponse object.
func (metadata *ActionMetadata) ScopeData(v interface{}) error {
	var dataLength C.int
	scopeData := C.action_metadata_get_scope_data(metadata.m, &dataLength)
	defer C.free(scopeData)
	return json.Unmarshal(C.GoBytes(scopeData, dataLength), v)
}
