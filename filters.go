package scopes

import (
	"sort"
)

// Filter is implemented by all scope filter types.
type Filter interface {
	serializeFilter() interface{}
}

type FilterDisplayHints int

const (
	FilterDisplayDefault FilterDisplayHints = 0
	FilterDisplayPrimary FilterDisplayHints = 1 << iota
)

// FilterState represents the current state of a set of filters.
type FilterState map[string]interface{}

// OptionSelectorFilter is used to implement single-select or multi-select filters.
type OptionSelectorFilter struct {
	Id           string
	DisplayHints FilterDisplayHints
	Label        string
	MultiSelect  bool
	Options      []FilterOption
}

type FilterOption struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

// NewOptionSelectorFilter creates a new option filter.
func NewOptionSelectorFilter(id, label string, multiSelect bool) *OptionSelectorFilter {
	return &OptionSelectorFilter{
		Id:          id,
		Label:       label,
		MultiSelect: multiSelect,
	}
}

// AddOption adds a new option to the filter.
func (f *OptionSelectorFilter) AddOption(id, label string) {
	f.Options = append(f.Options, FilterOption{
		Id:    id,
		Label: label,
	})
}

// SetDisplayHints changes how the filter is displayed.
func (f *OptionSelectorFilter) SetDisplayHints(hints FilterDisplayHints) {
	f.DisplayHints = hints
}

func (f *OptionSelectorFilter) isValidOption(optionId string) bool {
	for _, o := range f.Options {
		if o.Id == optionId {
			return true
		}
	}
	return false
}

// HasActiveOption returns true if any of the filters options are active.
func (f *OptionSelectorFilter) HasActiveOption(state FilterState) bool {
	for _, optionId := range f.ActiveOptions(state) {
		if f.isValidOption(optionId) {
			return true
		}
	}
	return false
}

// ActiveOptions returns the filter's active options from the filter state.
func (f *OptionSelectorFilter) ActiveOptions(state FilterState) []string {
	selected, _ := state[f.Id].([]string)
	return selected
}

// UpdateState updates the value of a particular option in the filter state.
func (f *OptionSelectorFilter) UpdateState(state FilterState, optionId string, active bool) {
	if !f.isValidOption(optionId) {
		panic("invalid option ID")
	}
	// For single-select filters, clear the previous state when
	// setting a new active option.
	if active && !f.MultiSelect {
		delete(state, f.Id)
	}
	// If the state isn't in a form we expect, treat it as empty
	selected, _ := state[f.Id].([]string)
	sort.Strings(selected)
	pos := sort.SearchStrings(selected, optionId)
	if active {
		if pos < len(selected) && selected[pos] != optionId {
			selected = append(selected[:pos], append([]string{optionId}, selected[pos:]...)...)
		}
	} else {
		if pos < len(selected) {
			selected = append(selected[:pos], selected[pos+1:]...)
		}
	}
	state[f.Id] = selected
}

func (f *OptionSelectorFilter) serializeFilter() interface{} {
	return map[string]interface{}{
		"filter_type":   "option_selector",
		"id":            f.Id,
		"display_hints": f.DisplayHints,
		"label":         f.Label,
		"multi_select":  f.MultiSelect,
		"options":       f.Options,
	}
}
