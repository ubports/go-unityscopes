package scopes

import (
	"encoding/json"
)

type PreviewWidget map[string]interface{}

/*
NewPreviewWidget creates a preview widget with the given name and type.

Widget type specific attributes can be set directly with
AddAttributeValue(), or mapped to result attributes with
AddAttributeMapping().

A list of available widget types and their associated attributes is
available here:

http://developer.ubuntu.com/api/scopes/sdk-14.04/preview_20widget_20types/
*/
func NewPreviewWidget(id, widgetType string) PreviewWidget {
	return PreviewWidget{"id": id, "type": widgetType}
}

// Id returns the name of this widget.
func (widget PreviewWidget) Id() string {
	return widget["id"].(string)
}

// WidgetType returns the type of this widget.
func (widget PreviewWidget) WidgetType() string {
	return widget["type"].(string)
}

// AddAttributeValue sets a widget attribute to a particular value.
func (widget PreviewWidget) AddAttributeValue(key string, value interface{}) {
	widget[key] = value
}

// AddAttributeMapping maps a widget attribute to a named result attribute.
func (widget PreviewWidget) AddAttributeMapping(key, fieldName string) {
	var components map[string]interface{}
	if comp, ok := widget["components"]; ok {
		components = comp.(map[string]interface{})
	} else {
		components = make(map[string]interface{})
		widget["components"] = components
	}
	components[key] = fieldName
}

func (widget PreviewWidget) data() ([]byte, error) {
	return json.Marshal(widget)
}
