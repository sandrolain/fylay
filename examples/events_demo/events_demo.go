package main

import (
	"fmt"

	flay "github.com/sandrolain/fylay"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// EventsDemo gestisce l'applicazione demo
type EventsDemo struct {
	builder    *flay.Builder
	window     fyne.Window
	clickCount int
}

// Implementazione EventHandler tradizionale (opzionale)
func (d *EventsDemo) OnButtonTapped(id string) {
	fmt.Printf("EventHandler: Button tapped: %s\n", id)

	if id == "traditionalBtn" {
		d.updateLabel("mixedLabel", "Clicked usando EventHandler tradizionale!")
	}
}

func (d *EventsDemo) OnEntryChanged(id, value string) {
	fmt.Printf("EventHandler: Entry changed: %s = %s\n", id, value)
}

// Helper per aggiornare label
func (d *EventsDemo) updateLabel(id, text string) {
	if elem := d.builder.GetElement(id); elem != nil {
		if label, ok := elem.(*widget.Label); ok {
			label.SetText(text)
		}
	}
}

// Helper per ottenere valore entry
func (d *EventsDemo) getEntryValue(id string) string {
	if elem := d.builder.GetElement(id); elem != nil {
		if entry, ok := elem.(*widget.Entry); ok {
			return entry.Text
		}
	}
	return ""
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("github.com/sandrolain/fylay - Events Demo")
	window.Resize(fyne.NewSize(600, 500))

	demo := &EventsDemo{window: window}

	// Crea il builder
	builder := flay.NewBuilder()

	// Imposta l'EventHandler tradizionale (opzionale)
	builder.SetEventHandler(demo)

	// ===== REGISTRA CALLBACK PER BUTTON (onclick) =====

	builder.On("onHelloClick", func(ctx *flay.EventContext) {
		fmt.Printf("Callback: Hello button clicked! Event: %s, TargetID: %s\n", ctx.EventName, ctx.TargetID)
		demo.updateLabel("messageLabel", "👋 Ciao! Benvenuto in Fylay!")
	})

	builder.On("onCountClick", func(ctx *flay.EventContext) {
		demo.clickCount++
		msg := fmt.Sprintf("🔢 Hai cliccato %d volte!", demo.clickCount)
		fmt.Printf("Callback: Count = %d, Event: %s\n", demo.clickCount, ctx.EventName)
		demo.updateLabel("messageLabel", msg)
	})

	builder.On("onResetClick", func(ctx *flay.EventContext) {
		demo.clickCount = 0
		fmt.Printf("Callback: Reset clicked! TargetID: %s\n", ctx.TargetID)
		demo.updateLabel("messageLabel", "🔄 Counter resettato a 0")
	})

	builder.On("onMixedClick", func(ctx *flay.EventContext) {
		fmt.Printf("Callback: Mixed approach button clicked! Event: %s\n", ctx.EventName)
		demo.updateLabel("mixedLabel", "Clicked usando Callback registrata con On()!")
	})

	// ===== REGISTRA CALLBACK PER ENTRY (onchange) =====

	builder.OnEntry("onNameChange", func(ctx *flay.EventContext) {
		fmt.Printf("Callback: Name changed to: %s (Event: %s, TargetID: %s)\n", ctx.Value, ctx.EventName, ctx.TargetID)

		if ctx.Value == "" {
			demo.updateLabel("statusLabel", "Il nome è vuoto...")
		} else {
			msg := fmt.Sprintf("✓ Nome: %s", ctx.Value)
			demo.updateLabel("statusLabel", msg)
		}
	})

	builder.OnEntry("onEmailChange", func(ctx *flay.EventContext) {
		fmt.Printf("Callback: Email changed to: %s (Event: %s)\n", ctx.Value, ctx.EventName)

		name := demo.getEntryValue("nameInput")

		if name != "" && ctx.Value != "" {
			msg := fmt.Sprintf("✓ %s <%s>", name, ctx.Value)
			demo.updateLabel("statusLabel", msg)
		} else if ctx.Value != "" {
			msg := fmt.Sprintf("✓ Email: %s", ctx.Value)
			demo.updateLabel("statusLabel", msg)
		}
	})

	// Carica il layout
	builder, content, err := flay.LoadLayoutFromFileWithHandler("events_demo.xml", demo)
	if err != nil {
		panic(fmt.Sprintf("Errore caricamento layout: %v", err))
	}

	demo.builder = builder

	window.SetContent(content)
	window.ShowAndRun()
}
