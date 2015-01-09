package scopes

import (
	"sort"
)

type Filter interface {
	serializeFilter() interface{}
}

type FilterDisplayHints int

const (
	FilterDisplayDefault FilterDisplayHints = 0
	FilterDisplayPrimary FilterDisplayHints = 1 << iota
)

type FilterState map[string]interface{}

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

func NewOptionSelectorFilter(id, label string, multiSelect bool) *OptionSelectorFilter {
	return &OptionSelectorFilter{
		Id:          id,
		Label:       label,
		MultiSelect: multiSelect,
	}
}

func (f *OptionSelectorFilter) AddOption(id, label string) {
	f.Options = append(f.Options, FilterOption{
		Id:    id,
		Label: label,
	})
}

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

func (f *OptionSelectorFilter) HasActiveOption(state FilterState) bool {
	for _, optionId := range f.ActiveOptions(state) {
		if f.isValidOption(optionId) {
			return true
		}
	}
	return false
}

func (f *OptionSelectorFilter) ActiveOptions(state FilterState) []string {
	selected, _ := state[f.Id].([]string)
	return selected
}

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
