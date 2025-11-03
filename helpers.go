package fylay

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
)

// LoadLayoutFromFile carica un layout da un file XML
func LoadLayoutFromFile(filepath string) (*Builder, fyne.CanvasObject, error) {
	builder := NewBuilder()

	file, err := os.Open(filepath) //nolint:gosec // Filepath is from user config, intentional
	if err != nil {
		return nil, nil, fmt.Errorf("errore apertura file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck // Defer close error can be ignored
	}()

	layout, err := builder.LoadLayout(file)
	if err != nil {
		return nil, nil, fmt.Errorf("errore caricamento layout: %w", err)
	}

	content, err := builder.Build(layout)
	if err != nil {
		return nil, nil, fmt.Errorf("errore costruzione interfaccia: %w", err)
	}

	return builder, content, nil
}

// applyMinSize applies width/height/min-width/min-height styles by wrapping the object if needed
func applyMinSize(obj fyne.CanvasObject, style map[string]string) fyne.CanvasObject {
	var width, height float32
	hasWidth := false
	hasHeight := false

	// Check for width or min-width
	if w := style["width"]; w != "" {
		if parsed, err := parseSize(w); err == nil {
			width = parsed
			hasWidth = true
		}
	} else if w := style["min-width"]; w != "" {
		if parsed, err := parseSize(w); err == nil {
			width = parsed
			hasWidth = true
		}
	}

	// Check for height or min-height
	if h := style["height"]; h != "" {
		if parsed, err := parseSize(h); err == nil {
			height = parsed
			hasHeight = true
		}
	} else if h := style["min-height"]; h != "" {
		if parsed, err := parseSize(h); err == nil {
			height = parsed
			hasHeight = true
		}
	}

	if !hasWidth && !hasHeight {
		return obj
	}

	// For objects with SetMinSize method (canvas objects)
	if sizable, ok := obj.(interface{ SetMinSize(fyne.Size) }); ok {
		currentSize := obj.MinSize()
		if !hasWidth {
			width = currentSize.Width
		}
		if !hasHeight {
			height = currentSize.Height
		}
		sizable.SetMinSize(fyne.NewSize(width, height))
		return obj
	}

	// For widgets, we need to use a container with min size
	currentSize := obj.MinSize()
	if !hasWidth {
		width = currentSize.Width
	}
	if !hasHeight {
		height = currentSize.Height
	}

	// Use a MaxLayout container with a sized rectangle behind
	rect := &fyne.Container{}
	rect.Resize(fyne.NewSize(width, height))

	return &fyne.Container{
		Layout:  &fixedSizeLayout{size: fyne.NewSize(width, height)},
		Objects: []fyne.CanvasObject{obj},
	}
}

// fixedSizeLayout is a simple layout that enforces a minimum size
type fixedSizeLayout struct {
	size fyne.Size
}

func (f *fixedSizeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return f.size
}

func (f *fixedSizeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, obj := range objects {
		obj.Resize(size)
		obj.Move(fyne.NewPos(0, 0))
	}
}

// LoadLayoutFromFileWithHandler carica un layout da file con event handler
func LoadLayoutFromFileWithHandler(filepath string, handler EventHandler) (*Builder, fyne.CanvasObject, error) {
	builder := NewBuilder()
	builder.SetEventHandler(handler)

	file, err := os.Open(filepath) //nolint:gosec // Filepath is from user config, intentional
	if err != nil {
		return nil, nil, fmt.Errorf("errore apertura file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck // Defer close error can be ignored
	}()

	layout, err := builder.LoadLayout(file)
	if err != nil {
		return nil, nil, fmt.Errorf("errore caricamento layout: %w", err)
	}

	content, err := builder.Build(layout)
	if err != nil {
		return nil, nil, fmt.Errorf("errore costruzione interfaccia: %w", err)
	}

	return builder, content, nil
}
