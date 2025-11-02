package fylay

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// BindingContext manages data bindings for the layout
type BindingContext struct {
	data     map[string]binding.DataItem
	widgets  map[string]fyne.CanvasObject
	bindings map[string][]string // widget ID -> bound data keys
}

// NewBindingContext creates a new binding context
func NewBindingContext() *BindingContext {
	return &BindingContext{
		data:     make(map[string]binding.DataItem),
		widgets:  make(map[string]fyne.CanvasObject),
		bindings: make(map[string][]string),
	}
}

// BindString binds a string data item to a key
func (bc *BindingContext) BindString(key string, value string) binding.String {
	if existingData, ok := bc.data[key]; ok {
		if strData, ok := existingData.(binding.String); ok {
			_ = strData.Set(value) //nolint:errcheck // Ignore error on set
			return strData
		}
	}

	strData := binding.NewString()
	_ = strData.Set(value) //nolint:errcheck // Ignore error on set
	bc.data[key] = strData
	return strData
}

// BindInt binds an integer data item to a key
func (bc *BindingContext) BindInt(key string, value int) binding.Int {
	if existingData, ok := bc.data[key]; ok {
		if intData, ok := existingData.(binding.Int); ok {
			_ = intData.Set(value) //nolint:errcheck // Ignore error on set
			return intData
		}
	}

	intData := binding.NewInt()
	_ = intData.Set(value) //nolint:errcheck // Ignore error on set
	bc.data[key] = intData
	return intData
}

// BindFloat binds a float data item to a key
func (bc *BindingContext) BindFloat(key string, value float64) binding.Float {
	if existingData, ok := bc.data[key]; ok {
		if floatData, ok := existingData.(binding.Float); ok {
			_ = floatData.Set(value) //nolint:errcheck // Ignore error on set
			return floatData
		}
	}

	floatData := binding.NewFloat()
	_ = floatData.Set(value) //nolint:errcheck // Ignore error on set
	bc.data[key] = floatData
	return floatData
}

// BindBool binds a boolean data item to a key
func (bc *BindingContext) BindBool(key string, value bool) binding.Bool {
	if existingData, ok := bc.data[key]; ok {
		if boolData, ok := existingData.(binding.Bool); ok {
			_ = boolData.Set(value) //nolint:errcheck // Ignore error on set
			return boolData
		}
	}

	boolData := binding.NewBool()
	_ = boolData.Set(value) //nolint:errcheck // Ignore error on set
	bc.data[key] = boolData
	return boolData
}

// GetBinding retrieves a binding by key
func (bc *BindingContext) GetBinding(key string) (binding.DataItem, bool) {
	data, ok := bc.data[key]
	return data, ok
}

// GetString retrieves a string binding value
func (bc *BindingContext) GetString(key string) (string, error) {
	data, ok := bc.data[key]
	if !ok {
		return "", fmt.Errorf("binding not found: %s", key)
	}

	strData, ok := data.(binding.String)
	if !ok {
		return "", fmt.Errorf("binding is not a string: %s", key)
	}

	return strData.Get()
}

// GetInt retrieves an integer binding value
func (bc *BindingContext) GetInt(key string) (int, error) {
	data, ok := bc.data[key]
	if !ok {
		return 0, fmt.Errorf("binding not found: %s", key)
	}

	intData, ok := data.(binding.Int)
	if !ok {
		return 0, fmt.Errorf("binding is not an int: %s", key)
	}

	return intData.Get()
}

// GetFloat retrieves a float binding value
func (bc *BindingContext) GetFloat(key string) (float64, error) {
	data, ok := bc.data[key]
	if !ok {
		return 0, fmt.Errorf("binding not found: %s", key)
	}

	floatData, ok := data.(binding.Float)
	if !ok {
		return 0, fmt.Errorf("binding is not a float: %s", key)
	}

	return floatData.Get()
}

// GetBool retrieves a boolean binding value
func (bc *BindingContext) GetBool(key string) (bool, error) {
	data, ok := bc.data[key]
	if !ok {
		return false, fmt.Errorf("binding not found: %s", key)
	}

	boolData, ok := data.(binding.Bool)
	if !ok {
		return false, fmt.Errorf("binding is not a bool: %s", key)
	}

	return boolData.Get()
}

// RegisterWidget registers a widget with an ID for binding
func (bc *BindingContext) RegisterWidget(id string, w fyne.CanvasObject) {
	bc.widgets[id] = w
}

// GetWidget retrieves a registered widget by ID
func (bc *BindingContext) GetWidget(id string) (fyne.CanvasObject, bool) {
	w, ok := bc.widgets[id]
	return w, ok
}

// BindWidgetToData binds a widget to a data key
func (bc *BindingContext) BindWidgetToData(widgetID string, dataKey string) error {
	w, ok := bc.widgets[widgetID]
	if !ok {
		return fmt.Errorf("widget not found: %s", widgetID)
	}

	data, ok := bc.data[dataKey]
	if !ok {
		return fmt.Errorf("binding not found: %s", dataKey)
	}

	// Bind based on widget and data type
	switch wid := w.(type) {
	case *widget.Label:
		if strData, ok := data.(binding.String); ok {
			wid.Bind(strData)
		} else {
			return fmt.Errorf("label requires string binding")
		}

	case *widget.Entry:
		if strData, ok := data.(binding.String); ok {
			wid.Bind(strData)
		} else {
			return fmt.Errorf("entry requires string binding")
		}

	case *widget.Check:
		if boolData, ok := data.(binding.Bool); ok {
			wid.Bind(boolData)
		} else {
			return fmt.Errorf("checkbox requires bool binding")
		}

	case *widget.Slider:
		if floatData, ok := data.(binding.Float); ok {
			wid.Bind(floatData)
		} else {
			return fmt.Errorf("slider requires float binding")
		}

	default:
		return fmt.Errorf("unsupported widget type for binding: %T", wid)
	}

	// Track binding
	bc.bindings[widgetID] = append(bc.bindings[widgetID], dataKey)

	return nil
}

// ParseBindAttribute parses a bind attribute value (e.g., "user.name")
// Returns the key to use for the binding
func ParseBindAttribute(bindAttr string) string {
	// For now, simple dot notation: "user.name" -> "user.name"
	// Future: could support more complex expressions
	return strings.TrimSpace(bindAttr)
}

// SetBindingContext sets the binding context on the builder
func (b *Builder) SetBindingContext(ctx *BindingContext) {
	b.bindingContext = ctx
}

// GetBindingContext returns the builder's binding context
func (b *Builder) GetBindingContext() *BindingContext {
	if b.bindingContext == nil {
		b.bindingContext = NewBindingContext()
	}
	return b.bindingContext
}
