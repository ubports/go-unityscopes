package unityscope

/*
#cgo CXXFLAGS: -std=c++11
#cgo pkg-config: libunity-scopes
#include "shim.h"
*/
import "C"
import "runtime"

type Reply struct {
	r C.SharedPtrData
}

type Scope interface {
	Query(query string, reply *Reply, cancelled <-chan bool)
}

func finalizeReply(r Reply) {
	C.destroy_reply_ptr(&r.r[0])
}

//export callScopeQuery
func callScopeQuery(scope Scope, query *C.char, reply_data *C.uintptr_t, cancel <-chan bool) {
	reply := new(Reply)
	runtime.SetFinalizer(reply, finalizeReply)
	C.init_reply_ptr(&reply.r[0], reply_data)
	go scope.Query(C.GoString(query), reply, cancel)
}


//export makeCancelChannel
func makeCancelChannel() chan bool {
	return make(chan bool, 1)
}

//export sendCancelChannel
func sendCancelChannel(ch chan bool) {
	ch <- true
}
