package main

import (
	"fmt"
	"os"

	flay "github.com/sandrolain/fylay"
	"github.com/sandrolain/fylay/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// DialogsDemo gestisce l'applicazione demo per i dialog
type DialogsDemo struct {
	builder *flay.Builder
	window  fyne.Window
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("github.com/sandrolain/fylay - Dialogs Demo")
	window.Resize(fyne.NewSize(500, 400))

	demo := &DialogsDemo{window: window}

	// Crea il builder
	builder := flay.NewBuilder()

	// ===== REGISTRA CALLBACK PER I DIALOG =====

	builder.On("onShowInfo", func(ctx *flay.EventContext) {
		fmt.Println("Showing info dialog")
		builder.ShowInfoDialog("Information", "This is an info message!")
	})

	builder.On("onShowError", func(ctx *flay.EventContext) {
		fmt.Println("Showing error dialog")
		builder.ShowErrorDialog("Error", "This is an error message!")
	})

	builder.On("onShowQuestion", func(ctx *flay.EventContext) {
		fmt.Println("Showing question dialog")
		result := builder.ShowQuestionDialog("Question", "Do you want to continue?")
		if result {
			demo.updateLabel("resultLabel", "✓ You clicked Yes")
		} else {
			demo.updateLabel("resultLabel", "✗ You clicked No")
		}
	})

	builder.On("onOpenFile", func(ctx *flay.EventContext) {
		fmt.Println("Opening file dialog")
		filepath, err := builder.ShowFileOpenDialog(
			"Select a file",
			dialog.FileFilter{Description: "Text files", Pattern: "*.txt"},
			dialog.FileFilter{Description: "Go files", Pattern: "*.go"},
			dialog.FileFilter{Description: "All files", Pattern: "*"},
		)
		if err != nil {
			demo.updateLabel("resultLabel", "File selection cancelled")
			return
		}
		demo.updateLabel("resultLabel", fmt.Sprintf("Selected: %s", filepath))

		// Leggi e mostra il contenuto
		// #nosec G304 - File path comes from user selection via dialog
		content, err := os.ReadFile(filepath)
		if err != nil {
			demo.updateLabel("resultLabel", fmt.Sprintf("Error reading file: %v", err))
			return
		}
		demo.updateTextArea("contentArea", string(content))
	})

	builder.On("onSaveFile", func(ctx *flay.EventContext) {
		fmt.Println("Opening save dialog")
		filepath, err := builder.ShowFileSaveDialog(
			"Save file",
			dialog.FileFilter{Description: "Text files", Pattern: "*.txt"},
			dialog.FileFilter{Description: "All files", Pattern: "*"},
		)
		if err != nil {
			demo.updateLabel("resultLabel", "Save cancelled")
			return
		}

		// Salva il contenuto
		content := demo.getTextAreaValue("contentArea")
		// #nosec G306 - User data file with read/write permissions
		err = os.WriteFile(filepath, []byte(content), 0600)
		if err != nil {
			builder.ShowErrorDialog("Error", fmt.Sprintf("Error saving file: %v", err))
			return
		}

		demo.updateLabel("resultLabel", fmt.Sprintf("Saved to: %s", filepath))
		builder.ShowInfoDialog("Success", "File saved successfully!")
	})

	builder.On("onSelectDir", func(ctx *flay.EventContext) {
		fmt.Println("Opening directory selection dialog")
		dirpath, err := builder.ShowDirSelectDialog("Select a directory")
		if err != nil {
			demo.updateLabel("resultLabel", "Directory selection cancelled")
			return
		}
		demo.updateLabel("resultLabel", fmt.Sprintf("Selected directory: %s", dirpath))
	})

	// Carica il layout
	builder, content, err := flay.LoadLayoutFromFileWithHandler("dialogs_demo.xml", demo)
	if err != nil {
		panic(fmt.Sprintf("Errore caricamento layout: %v", err))
	}

	demo.builder = builder

	window.SetContent(content)
	window.ShowAndRun()
}

// Helper per aggiornare label
func (d *DialogsDemo) updateLabel(id, text string) {
	if elem := d.builder.GetElement(id); elem != nil {
		if label, ok := elem.(*widget.Label); ok {
			label.SetText(text)
		}
	}
}

// Helper per aggiornare text area
func (d *DialogsDemo) updateTextArea(id, text string) {
	if elem := d.builder.GetElement(id); elem != nil {
		if entry, ok := elem.(*widget.Entry); ok {
			entry.SetText(text)
		}
	}
}

// Helper per ottenere valore text area
func (d *DialogsDemo) getTextAreaValue(id string) string {
	if elem := d.builder.GetElement(id); elem != nil {
		if entry, ok := elem.(*widget.Entry); ok {
			return entry.Text
		}
	}
	return ""
}

// OnButtonTapped implementa EventHandler (opzionale)
func (d *DialogsDemo) OnButtonTapped(id string) {
	// Non usato in questo esempio
}

// OnEntryChanged implementa EventHandler (opzionale)
func (d *DialogsDemo) OnEntryChanged(id, value string) {
	// Non usato in questo esempio
}
