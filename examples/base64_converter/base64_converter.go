package main

import (
	"encoding/base64"
	"fmt"
	"os"

	flay "github.com/sandrolain/fylay"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Base64Converter handles the Base64 conversion application
type Base64Converter struct {
	builder    *flay.Builder
	window     fyne.Window
	inputArea  *widget.Entry
	outputArea *widget.Entry
	inputInfo  *widget.Label
	outputInfo *widget.Label
	status     *widget.Label
}

func main() {
	// Create the Fyne application with unique ID
	myApp := app.NewWithID("com.fylay.base64converter")
	window := myApp.NewWindow("Base64 Converter")
	window.Resize(fyne.NewSize(800, 600))

	// Create the event handler
	converter := &Base64Converter{
		window: window,
	}

	// Load the XML layout
	builder, content, err := flay.LoadLayoutFromFileWithHandler("base64_converter.xml", converter)
	if err != nil {
		dialog.ShowError(err, window)
		window.ShowAndRun()
		return
	}

	// Save the builder
	converter.builder = builder

	// Register onclick callbacks
	builder.On("onEncode", func(ctx *flay.EventContext) {
		converter.handleEncode()
	})
	builder.On("onDecode", func(ctx *flay.EventContext) {
		converter.handleDecode()
	})
	builder.On("onLoadFile", func(ctx *flay.EventContext) {
		converter.handleLoadFile()
	})
	builder.On("onSaveFile", func(ctx *flay.EventContext) {
		converter.handleSaveFile()
	})

	// Get references to widgets
	if w := builder.GetWidget("inputArea"); w != nil {
		converter.inputArea = w.(*widget.Entry)
		converter.inputArea.OnChanged = func(text string) {
			converter.updateInputInfo(text)
		}
	}
	if w := builder.GetWidget("outputArea"); w != nil {
		converter.outputArea = w.(*widget.Entry)
	}
	if w := builder.GetWidget("inputInfo"); w != nil {
		converter.inputInfo = w.(*widget.Label)
	}
	if w := builder.GetWidget("outputInfo"); w != nil {
		converter.outputInfo = w.(*widget.Label)
	}
	if w := builder.GetWidget("statusLabel"); w != nil {
		converter.status = w.(*widget.Label)
	}

	// Set the content and show the window
	window.SetContent(content)
	window.ShowAndRun()
}

// OnButtonTapped handles button tap events
func (c *Base64Converter) OnButtonTapped(id string) {
	switch id {
	case "encodeBtn":
		c.handleEncode()
	case "decodeBtn":
		c.handleDecode()
	case "loadFileBtn":
		c.handleLoadFile()
	case "saveFileBtn":
		c.handleSaveFile()
	}
}

// OnEntryChanged handles entry change events (required by EventHandler interface)
func (c *Base64Converter) OnEntryChanged(id, value string) {
	// Input changes are handled by direct OnChanged callback
}

// handleEncode converts input text to Base64
func (c *Base64Converter) handleEncode() {
	input := c.inputArea.Text
	if input == "" {
		c.setStatus("Error: Input is empty")
		return
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	fmt.Printf("encoded: %v\n", encoded)
	c.outputArea.SetText(encoded)
	c.updateOutputInfo(encoded)
	c.setStatus("Successfully encoded to Base64")
}

// handleDecode converts Base64 input to text
func (c *Base64Converter) handleDecode() {
	input := c.inputArea.Text
	if input == "" {
		c.setStatus("Error: Input is empty")
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		c.setStatus(fmt.Sprintf("Error: Invalid Base64 - %v", err))
		c.outputArea.SetText("")
		c.updateOutputInfo("")
		return
	}

	c.outputArea.SetText(string(decoded))
	c.updateOutputInfo(string(decoded))
	c.setStatus("Successfully decoded from Base64")
}

// handleLoadFile opens a file dialog, converts the content to Base64, and saves it
func (c *Base64Converter) handleLoadFile() {
	// Select file to load
	filepath, err := c.builder.ShowFileOpenDialog("Select File to Convert to Base64")
	if err != nil {
		c.setStatus(fmt.Sprintf("Error: %v", err))
		return
	}
	if filepath == "" {
		c.setStatus("File selection cancelled")
		return
	}

	// Read file content
	// #nosec G304 - File path comes from user selection via dialog
	data, err := os.ReadFile(filepath)
	if err != nil {
		c.setStatus(fmt.Sprintf("Error reading file: %v", err))
		return
	}

	c.setStatus(fmt.Sprintf("Loaded %d bytes, encoding...", len(data)))

	// Encode to Base64 directly without showing in UI
	encoded := base64.StdEncoding.EncodeToString(data)

	// Ask where to save the encoded result
	savepath, err := c.builder.ShowFileSaveDialog("Save Base64 Encoded File")
	if err != nil {
		c.setStatus(fmt.Sprintf("Error: %v", err))
		return
	}
	if savepath == "" {
		c.setStatus("Save cancelled")
		return
	}

	// Save encoded content
	// #nosec G306 - User data file with read/write permissions
	err = os.WriteFile(savepath, []byte(encoded), 0600)
	if err != nil {
		c.setStatus(fmt.Sprintf("Error saving file: %v", err))
		return
	}

	c.setStatus(fmt.Sprintf("Encoded %d bytes → %d bytes and saved to: %s", len(data), len(encoded), savepath))
}

// handleSaveFile opens a save dialog and saves the output content
func (c *Base64Converter) handleSaveFile() {
	output := c.outputArea.Text
	if output == "" {
		c.setStatus("Error: Output is empty")
		return
	}

	filepath, err := c.builder.ShowFileSaveDialog("Save Output to File")
	if err != nil {
		c.setStatus(fmt.Sprintf("Error: %v", err))
		return
	}
	if filepath == "" {
		c.setStatus("Save cancelled")
		return
	}

	// #nosec G306 - User data file with read/write permissions
	err = os.WriteFile(filepath, []byte(output), 0600)
	if err != nil {
		c.setStatus(fmt.Sprintf("Error saving file: %v", err))
		return
	}

	c.setStatus(fmt.Sprintf("Saved to: %s (%d bytes)", filepath, len(output)))
}

// updateInputInfo updates the input statistics label
func (c *Base64Converter) updateInputInfo(text string) {
	chars := len([]rune(text))
	bytes := len([]byte(text))
	c.inputInfo.SetText(fmt.Sprintf("Chars: %d | Bytes: %d", chars, bytes))
}

// updateOutputInfo updates the output statistics label
func (c *Base64Converter) updateOutputInfo(text string) {
	chars := len([]rune(text))
	bytes := len([]byte(text))
	c.outputInfo.SetText(fmt.Sprintf("Chars: %d | Bytes: %d", chars, bytes))
}

// setStatus updates the status label
func (c *Base64Converter) setStatus(message string) {
	c.status.SetText(message)
}
