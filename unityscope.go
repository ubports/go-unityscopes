package unityscope

/*
#cgo CXXFLAGS: -std=c++11
#cgo pkg-config: libunity-scopes
#include <stdlib.h>
#include "shim.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"runtime"
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

type SearchReply struct {
	r C.SharedPtrData
}

func finalizeSearchReply(reply *SearchReply) {
	C.destroy_search_reply_ptr(&reply.r[0])
}

func (reply *SearchReply) Finished() {
	C.search_reply_finished(&reply.r[0])
}

func (reply *SearchReply) Error(err error) {
	errString := err.Error()
	C.search_reply_error(&reply.r[0], unsafe.Pointer(&errString))
}

func (reply *SearchReply) RegisterCategory(id, title, icon, template string) *Category {
	cat := new(Category)
	runtime.SetFinalizer(cat, finalizeCategory)
	C.search_reply_register_category(&reply.r[0], unsafe.Pointer(&id), unsafe.Pointer(&title), unsafe.Pointer(&icon), unsafe.Pointer(&template), &cat.c[0])
	return cat
}

func (reply *SearchReply) Push(result *CategorisedResult) error {
	var errorString *C.char = nil
	C.search_reply_push(&reply.r[0], result.result, &errorString)
	return checkError(errorString)
}

type PreviewReply struct {
	r C.SharedPtrData
}

func finalizePreviewReply(reply *PreviewReply) {
	C.destroy_search_reply_ptr(&reply.r[0])
}

func (reply *PreviewReply) Finished() {
	C.preview_reply_finished(&reply.r[0])
}

func (reply *PreviewReply) Error(err error) {
	errString := err.Error()
	C.preview_reply_error(&reply.r[0], unsafe.Pointer(&errString))
}

func (reply *PreviewReply) PushWidgets(widgets []PreviewWidget) error {
	widget_data := make([]string, len(widgets))
	for i, w := range widgets {
		data, err := w.data()
		if err != nil {
			return err
		}
		widget_data[i] = string(data)
	}
	var errorString *C.char = nil
	C.preview_reply_push_widgets(&reply.r[0], unsafe.Pointer(&widget_data[0]), C.int(len(widget_data)), &errorString)
	return checkError(errorString)
}

func (reply *PreviewReply) PushAttr(attr string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	json_value := string(data)
	var errorString *C.char = nil
	C.preview_reply_push_attr(&reply.r[0], unsafe.Pointer(&attr), unsafe.Pointer(&json_value), &errorString)
	return checkError(errorString)
}


type Category struct {
	c C.SharedPtrData
}

func finalizeCategory(cat *Category) {
	C.destroy_category_ptr(&cat.c[0])
}

// Scope defines the interface that scope implementations must implement
type Scope interface {
	Search(query string, reply *SearchReply, cancelled <-chan bool) error
	Preview(result *Result, reply *PreviewReply, cancelled <-chan bool) error
}

//export callScopeSearch
func callScopeSearch(scope Scope, query *C.char, reply_data *C.uintptr_t, cancel <-chan bool) {
	reply := new(SearchReply)
	runtime.SetFinalizer(reply, finalizeSearchReply)
	C.init_search_reply_ptr(&reply.r[0], reply_data)

	go func() {
		err := scope.Search(C.GoString(query), reply, cancel)
		if err != nil {
			reply.Error(err)
			return
		}
		reply.Finished()
	}()
}

//export callScopePreview
func callScopePreview(scope Scope, res uintptr, reply_data *C.uintptr_t, cancel <-chan bool) {
	result := new(Result)
	runtime.SetFinalizer(result, finalizeResult)
	result.result = unsafe.Pointer(res)

	reply := new(PreviewReply)
	runtime.SetFinalizer(reply, finalizePreviewReply)
	C.init_preview_reply_ptr(&reply.r[0], reply_data)

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
 run until the program exits.  For example:

   func main() {
       scope := ...
       unityscope.Run("myscope", os.Args[1], scope)
   }
*/
func Run(scopeName, runtimeConfig string, scope Scope) {
	cScopeName := C.CString(scopeName)
	defer C.free(unsafe.Pointer(cScopeName))
	cRuntimeConfig := C.CString(runtimeConfig)
	defer C.free(unsafe.Pointer(cRuntimeConfig))

	C.run_scope(cScopeName, cRuntimeConfig, unsafe.Pointer(&scope))
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
