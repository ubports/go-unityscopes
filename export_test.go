package scopes

// This file exports certain private functions for use by tests.

func NewTestingResult() *Result {
	return newTestingResult()
}

func NewTestingScopeMetadata(json_data string) ScopeMetadata {
	return newScopeMetadata(json_data)
}
