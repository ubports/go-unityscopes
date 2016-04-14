package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"reflect"
	"unsafe"
)

func strData(s string) C.StrData {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return C.StrData{
		data: (*C.char)(unsafe.Pointer(h.Data)),
		length: C.long(h.Len),
	}
}
