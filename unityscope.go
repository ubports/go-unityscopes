package unityscope

// #cgo CXXFLAGS: -std=c++11
// #cgo pkg-config: libunity-scopes
// #include "shim.h"
import "C"

type Reply struct {
	r C.SharedPtrData
}

type Scope interface {
	Query(query string, reply *Reply, cancelled <-chan bool)
}
