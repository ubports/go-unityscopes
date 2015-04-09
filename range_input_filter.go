package scopes

import (
	"errors"
	"fmt"
	"reflect"
)

// RangeInputFilter is a range filter which allows a start and end value to be entered by user, and any of them is optional.
type RangeInputFilter struct {
	filterWithLabel
	StartLabel string
	EndLabel   string
	UnitLabel  string
}

// NewRangeInputFilter creates a new range input filter.
func NewRangeInputFilter(id, label, start_label, end_label, unit_label string) *RangeInputFilter {
	return &RangeInputFilter{
		filterWithLabel: filterWithLabel{
			filterBase: filterBase{
				Id:           id,
				DisplayHints: FilterDisplayDefault,
				FilterType:   "range_input",
			},
			Label: label,
		},
		StartLabel: start_label,
		EndLabel:   end_label,
		UnitLabel:  unit_label,
	}
}

// StartValue gets the start value of this filter from filter state object.
// If the value is not set for the filter it returns false as the second return statement,
// it returns true otherwise
func (f *RangeInputFilter) StartValue(state FilterState) (float64, bool) {
	var start float64
	var ok bool
	_, ok = state[f.Id]
	if ok {
		if reflect.TypeOf(state[f.Id]).Kind() == reflect.Slice {
			s := reflect.ValueOf(state[f.Id])
			if s.Len() != 2 {
				// something went really bad.
				// we should have just 2 values
				panic("RangeInputFilter:Values unexpected number of values found.")
			}
			start_interface := s.Index(0).Interface()
			// try to convert to float64 if it's a float64 or an integer
			if start_interface != nil {
				start, ok = start_interface.(float64)
				if !ok {
					var start_int int
					start_int, ok = start_interface.(int)
					if ok {
						start = float64(start_int)
					}
				}
			} else {
				// the value is undefined (nil)
				ok = false
			}
		}
	}
	return start, ok
}

// EndValue gets the end value of this filter from filter state object.
// If the value is not set for the filter it returns false as the second return statement,
// it returns true otherwise
func (f *RangeInputFilter) EndValue(state FilterState) (float64, bool) {
	var end float64
	var ok bool
	_, ok = state[f.Id]
	if ok {
		if reflect.TypeOf(state[f.Id]).Kind() == reflect.Slice {
			s := reflect.ValueOf(state[f.Id])
			if s.Len() != 2 {
				// something went really bad.
				// we should have just 2 values
				panic("RangeInputFilter:Values unexpected number of values found.")
			}
			end_interface := s.Index(1).Interface()
			// try to convert to float64 if it's a float64 or an integer
			if end_interface != nil {
				end, ok = end_interface.(float64)
				if !ok {
					var end_int int
					end_int, ok = end_interface.(int)
					if ok {
						end = float64(end_int)
					}
				}
			} else {
				// the value is undefined (nil)
				ok = false
			}
		}
	}
	return end, ok
}

func (f *RangeInputFilter) checkValidType(value interface{}) bool {
	if value != nil {
		switch value.(type) {
		case int:
		case float64:
			return true
		default:
			return false
		}
	}
	// we accept the nil value
	return true
}

func convertToFloat(value interface{}) float64 {
	if value != nil {
		fVal, ok := value.(float64)
		if !ok {
			iVal, ok := value.(int)
			if !ok {
				panic(fmt.Sprint("RangeInputFilter:convertToFloat unexpected type for given value %v", value))
			}
			return float64(iVal)
		}
		return fVal
	} else {
		panic("RangeInputFilter:convertToFloat nil values are not accepted")
	}
}

// UpdateState updates the value of the filter
func (f *RangeInputFilter) UpdateState(state FilterState, start, end interface{}) error {
	if !f.checkValidType(start) {
		return errors.New("RangeInputFilter:UpdateState: Bad type for start value. Valid types are int float64 and nil")
	}
	if !f.checkValidType(end) {
		return errors.New("RangeInputFilter:UpdateState: Bad type for end value. Valid types are int float64 and nil")
	}

	if start == nil && end == nil {
		// remove the state
		delete(state, f.Id)
		return nil
	}
	if start != nil && end != nil {
		fStart := convertToFloat(start)
		fEnd := convertToFloat(end)
		if fStart >= fEnd {
			return errors.New(fmt.Sprintf("RangeInputFilter::UpdateState(): start_value %v is greater or equal to end_value %v for filter %s", start, end, f.Id))
		}
	}
	state[f.Id] = []interface{}{start, end}
	return nil
}

func (f *RangeInputFilter) serializeFilter() interface{} {
	return map[string]interface{}{
		"start_label": f.StartLabel,
		"end_label":   f.EndLabel,
		"unit_label":  f.UnitLabel,
	}
}
