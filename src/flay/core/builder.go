// Package core contains the core types and builder for Fylay.
// This package provides the fundamental structures for parsing XML layouts
// and building Fyne UI components.
package core

import (
	"encoding/xml"
	"fmt"
	"io"

	"fyne.io/fyne/v2"
)

// Layout represents a parsed XML layout
type Layout struct {
	XMLName xml.Name `xml:"Layout"`
	Styles  []Style  `xml:"Style"`
	Root    Element  `xml:",any"`
}

// Style represents a CSS style rule
type Style struct {
	Selector   string            `xml:"selector,attr"`
	Properties map[string]string `xml:"-"`
	RawCSS     string            `xml:",innerxml"`
}

// Element represents a generic layout element
type Element struct {
	XMLName    xml.Name
	ID         string     `xml:"id,attr"`
	Class      string     `xml:"class,attr"`
	Style      string     `xml:"style,attr"`
	Text       string     `xml:"text,attr"`
	Attributes []xml.Attr `xml:",any,attr"`
	Children   []Element  `xml:",any"`
	Content    string     `xml:",chardata"`
}

// GetAttr returns the value of an attribute by name
func (e *Element) GetAttr(name string) string {
	for _, attr := range e.Attributes {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// EventHandler handles widget events
type EventHandler interface {
	OnButtonTapped(id string)
	OnEntryChanged(id, value string)
}

// Builder builds Fyne widgets from layouts
type Builder struct {
	styles         map[string]Style
	elements       map[string]fyne.CanvasObject
	eventHandler   EventHandler
	eventCallbacks map[string]func()
	entryCallbacks map[string]func(string)
}

// NewBuilder creates a new builder instance
func NewBuilder() *Builder {
	return &Builder{
		styles:         make(map[string]Style),
		elements:       make(map[string]fyne.CanvasObject),
		eventCallbacks: make(map[string]func()),
		entryCallbacks: make(map[string]func(string)),
	}
}

// SetEventHandler sets the event handler for the builder
func (b *Builder) SetEventHandler(handler EventHandler) {
	b.eventHandler = handler
}

// On registers a callback for a specific event (for buttons)
func (b *Builder) On(eventName string, callback func()) {
	b.eventCallbacks[eventName] = callback
}

// OnEntry registers a callback for Entry events
func (b *Builder) OnEntry(eventName string, callback func(string)) {
	b.entryCallbacks[eventName] = callback
}

// LoadLayout loads a layout from an XML reader
func (b *Builder) LoadLayout(r io.Reader) (*Layout, error) {
	var layout Layout
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&layout); err != nil {
		return nil, fmt.Errorf("XML parsing error: %w", err)
	}

	// Parse CSS styles - Note: this requires parseCSS function from styles package
	// This will be refactored when we move the actual implementation
	for i := range layout.Styles {
		// Placeholder: actual parsing will be done by styles package
		layout.Styles[i].Properties = make(map[string]string)
		b.styles[layout.Styles[i].Selector] = layout.Styles[i]
	}

	return &layout, nil
}

// GetElement returns an element by ID
func (b *Builder) GetElement(id string) fyne.CanvasObject {
	return b.elements[id]
}

// GetStyles returns the loaded styles
func (b *Builder) GetStyles() map[string]Style {
	return b.styles
}

// GetEventCallback returns a registered event callback
func (b *Builder) GetEventCallback(name string) func() {
	return b.eventCallbacks[name]
}

// GetEntryCallback returns a registered entry callback
func (b *Builder) GetEntryCallback(name string) func(string) {
	return b.entryCallbacks[name]
}

// RegisterElement registers an element with the builder
func (b *Builder) RegisterElement(id string, obj fyne.CanvasObject) {
	if id != "" {
		b.elements[id] = obj
	}
}
