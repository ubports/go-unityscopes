package unityscope

import (
	"encoding/json"
)

type PreviewWidget map[string]interface{}

func NewPreviewWidget(id, widgetType string) PreviewWidget {
	return PreviewWidget{"id": id, "type": widgetType}
}

func (widget PreviewWidget) Id() string {
	return widget["id"].(string)
}

func (widget PreviewWidget) WidgetType() string {
	return widget["type"].(string)
}

func (widget PreviewWidget) AddAttributeValue(key string, value interface{}) {
	widget[key] = value
}

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
