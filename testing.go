package scopes

// #include "shim.h"
import "C"
import (
	"encoding/json"
)

// These functions are used by tests.  They are not part of a
// *_test.go file because they make use of cgo.

func newTestingResult() *Result {
	return makeResult(C.new_testing_result())
}

func newScopeMetadata(json_data string) ScopeMetadata {
	var scopeMetadata ScopeMetadata
	if err := json.Unmarshal([]byte(json_data), &scopeMetadata); err != nil {
		panic(err)
	}

	return scopeMetadata
}
