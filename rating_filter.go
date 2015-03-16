package scopes

import (
	"sort"
)

// RatingFilter is a filter that allows for rating-based selection
type RatingFilter struct {
	filterWithOptions
	OnIcon  string
	OffIcon string
}

// NewRatingFilter creates a new rating filter.
func NewRatingFilter(id, label string) *RatingFilter {
	return &RatingFilter{
		filterWithOptions: filterWithOptions{
			filterWithLabel: filterWithLabel{
				filterBase: filterBase{
					Id:           id,
					DisplayHints: FilterDisplayDefault,
					FilterType:   "rating",
				},
				Label: label,
			},
		},
	}
}

// HasActiveRating checks if a rating option is active.
func (f *RatingFilter) HasActiveRating(state FilterState) bool {
	return f.HasActiveOption(state)
}

// HasActiveRating checks if a rating option is active.
func (f *RatingFilter) ActiveRating(state FilterState) []interface{} {
	return f.ActiveOptions(state)
}

// UpdateState updates the value of a particular option in the filter state.
func (f *RatingFilter) UpdateState(state FilterState, optionId string, active bool) {
	if !f.isValidOption(optionId) {
		panic("invalid option ID")
	}
	// If the state isn't in a form we expect, treat it as empty
	selected, _ := state[f.Id].([]string)
	sort.Strings(selected)
	pos := sort.SearchStrings(selected, optionId)
	if active {
		if pos == len(selected) {
			selected = append(selected, optionId)
		} else if pos < len(selected) && selected[pos] != optionId {
			selected = append(selected[:pos], append([]string{optionId}, selected[pos:]...)...)
		}
	} else {
		if pos < len(selected) {
			selected = append(selected[:pos], selected[pos+1:]...)
		}
	}
	state[f.Id] = selected
}

func (f *RatingFilter) serializeFilter() interface{} {
	return map[string]interface{}{
		"filter_type":   f.FilterType,
		"id":            f.Id,
		"display_hints": f.DisplayHints,
		"label":         f.Label,
		"options":       f.Options,
	}
}
