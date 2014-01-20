package unityscope

/*
#cgo CXXFLAGS: -std=c++11
#cgo pkg-config: libunity-scopes
#include <stdlib.h>
#include "shim.h"
*/
import "C"
import (
	"errors"
	"log"
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

type Reply struct {
	r C.SharedPtrData
}

func (reply *Reply) Finished() {
	C.reply_finished(&reply.r[0])
}

func (reply *Reply) RegisterCategory(id, title, icon string) *Category {
	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	cIcon := C.CString(icon)
	defer C.free(unsafe.Pointer(cIcon))

	cat := new(Category)
	runtime.SetFinalizer(cat, finalizeCategory)
	C.reply_register_category(&reply.r[0], cId, cTitle, cIcon, &cat.c[0])
	return cat
}

func (reply *Reply) Push(result *CategorisedResult) error {
	var errorString *C.char = nil
	C.reply_push(&reply.r[0], result.result, &errorString)
	return checkError(errorString)
}

type Category struct {
	c C.SharedPtrData
}

func finalizeCategory(cat *Category) {
	C.destroy_category_ptr(&cat.c[0])
}

type CategorisedResult struct {
	result unsafe.Pointer
}

func NewCategorisedResult(category *Category) *CategorisedResult {
	res := new(CategorisedResult)
	runtime.SetFinalizer(res, finalizeCategorisedResult)
	res.result = C.new_categorised_result(&category.c[0])
	return res
}

func finalizeCategorisedResult(res *CategorisedResult) {
	if res.result != nil {
		C.destroy_categorised_result(res.result)
	}
	res.result = nil
}

func (res *CategorisedResult) SetURI(uri string) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	C.categorised_result_set_uri(res.result, cUri)
}

func (res *CategorisedResult) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.categorised_result_set_title(res.result, cTitle)
}

func (res *CategorisedResult) SetArt(art string) {
	cArt := C.CString(art)
	defer C.free(unsafe.Pointer(cArt))
	C.categorised_result_set_art(res.result, cArt)
}

func (res *CategorisedResult) SetDndURI(uri string) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	C.categorised_result_set_dnd_uri(res.result, cUri)
}

// Scope defines the interface that 
type Scope interface {
	Query(query string, reply *Reply, cancelled <-chan bool) error
}

func finalizeReply(reply *Reply) {
	C.destroy_reply_ptr(&reply.r[0])
}

//export callScopeQuery
func callScopeQuery(scope Scope, query *C.char, reply_data *C.uintptr_t, cancel <-chan bool) {
	reply := new(Reply)
	runtime.SetFinalizer(reply, finalizeReply)
	C.init_reply_ptr(&reply.r[0], reply_data)
	go func() {
		err := scope.Query(C.GoString(query), reply, cancel)
		if err != nil {
			log.Println(err)
		}
		reply.Finished()
	}()
}

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
