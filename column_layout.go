package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"runtime"
	"unsafe"
)

// ColumnLayout is used represent different representations of a widget.
// Depending on the device your applications runs you can have several predefined
// column layouts in order to represent your view in the way it fits better the
// aspect ratio.
type ColumnLayout struct {
	c *C._ColumnLayout
}

func finalizeColumnLayout(layout *ColumnLayout) {
	if layout.c != nil {
		C.destroy_column_layout(layout.c)
	}
	layout.c = nil
}

func makeColumnLayout(c *C._ColumnLayout) *ColumnLayout {
	layout := new(ColumnLayout)
	runtime.SetFinalizer(layout, finalizeColumnLayout)
	layout.c = c
	return layout
}

// NewColumnLayout Creates a layout definition that expects num_of_columns columns to be added with ColumnLayout.AddColumn.
func NewColumnLayout(num_columns int) *ColumnLayout {
	return makeColumnLayout(C.new_column_layout(C.int(num_columns)))
}

// AddColumn adds a new column and assigns widgets to it.
// ColumnLayout expects exactly the number of columns passed to the constructor to be created with the
// AddColumn method.
func (layout *ColumnLayout) AddColumn(columns []string) error {
	api_columns := make([]*C._GoString, len(columns))
	for i := 0; i < len(columns); i++ {
		api_columns[i] = (*C._GoString)(unsafe.Pointer((&columns[i])))
	}
	var errorString *C.char = nil
	C.column_layout_add_column(layout.c, (**C._GoString)(unsafe.Pointer(&api_columns[0])), C.int(len(columns)), &errorString)

	if err := checkError(errorString); err != nil {
		return err
	}
	return nil
}

// NumberOfColumns gets the number of columns expected by this layout as specified in the constructor.
func (layout *ColumnLayout) NumberOfColumns() int {
	return int(C.column_layout_number_of_columns(layout.c))
}

// Size gets the current number of columns in this layout.
func (layout *ColumnLayout) Size() int {
	return int(C.column_layout_size(layout.c))
}

func (layout *ColumnLayout) Column(column int) ([]string, error) {
	var (
		length      C.int
		errorString *C.char = nil
	)

	var value []string
	data := C.column_layout_column(layout.c, C.int(column), &length, &errorString)
	if err := checkError(errorString); err != nil {
		return nil, err
	}
	defer C.free(data)

	err := json.Unmarshal(C.GoBytes(data, length), &value)

	return value, err
}
