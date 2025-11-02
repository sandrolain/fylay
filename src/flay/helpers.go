package flay

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
