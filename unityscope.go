package scopes

/*
#cgo CXXFLAGS: -std=c++11
#cgo pkg-config: libunity-scopes
#include <stdlib.h>
#include "shim.h"
*/
import "C"
import (
	"errors"
	"os"
	"sync"
	"unsafe"
)

func checkError(errorString *C.char) (err error) {
	if errorString != nil {
		err = errors.New(C.GoString(errorString))
		C.free(unsafe.Pointer(errorString))
	}
	return
}

// Category represents a search result category.
type Category struct {
	c C.SharedPtrData
}

func finalizeCategory(cat *Category) {
	C.destroy_category_ptr(&cat.c[0])
}

// Scope defines the interface that scope implementations must implement
type Scope interface {
	Search(query *CannedQuery, reply *SearchReply, cancelled <-chan bool) error
	Preview(result *Result, reply *PreviewReply, cancelled <-chan bool) error
}

//export callScopeSearch
func callScopeSearch(scope Scope, queryData uintptr, replyData *C.uintptr_t, cancel <-chan bool) {
	query := makeCannedQuery(unsafe.Pointer(queryData))
	reply := makeSearchReply(replyData)

	go func() {
		err := scope.Search(query, reply, cancel)
		if err != nil {
			reply.Error(err)
			return
		}
		reply.Finished()
	}()
}

//export callScopePreview
func callScopePreview(scope Scope, res uintptr, replyData *C.uintptr_t, cancel <-chan bool) {
	result := makeResult(res)
	reply := makePreviewReply(replyData)

	go func() {
		err := scope.Preview(result, reply, cancel)
		if err != nil {
			reply.Error(err)
			return
		}
		reply.Finished()
	}()
}

/*
Run will initialise the scope runtime and make a scope availble.  It
is intended to be called from the program's main function, and will
run until the program exits.
*/
func Run(scopeName string, scope Scope) error {
	if len(os.Args) < 3 {
		return errors.New("Expected to find runtime and scope config command line arguments")
	}
	runtimeConfig := os.Args[1]
	scopeConfig := os.Args[2]

	var errorString *C.char = nil
	C.run_scope(unsafe.Pointer(&scopeName), unsafe.Pointer(&runtimeConfig), unsafe.Pointer(&scopeConfig), unsafe.Pointer(&scope), &errorString)
	return checkError(errorString)
}

var (
	cancelChannels = make(map[chan bool] bool)
	cancelChannelsLock sync.Mutex
)

//export makeCancelChannel
func makeCancelChannel() chan bool {
	ch := make(chan bool, 1)
	cancelChannelsLock.Lock()
	cancelChannels[ch] = true
	cancelChannelsLock.Unlock()
	return ch
}

//export sendCancelChannel
func sendCancelChannel(ch chan bool) {
	ch <- true
}

//export releaseCancelChannel
func releaseCancelChannel(ch chan bool) {
	cancelChannelsLock.Lock()
	delete(cancelChannels, ch)
	cancelChannelsLock.Unlock()
}
