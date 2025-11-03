package main

import (
	"fmt"

	flay "github.com/sandrolain/fylay"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// HelloApp gestisce l'applicazione di esempio
type HelloApp struct {
	builder *flay.Builder
	window  fyne.Window
}

// OnButtonTapped gestisce i click sui pulsanti
func (a *HelloApp) OnButtonTapped(id string) {
	if id == "greetBtn" {
		// Ottieni il nome dall'input
		name := "Straniero"
		if entry := a.builder.GetElement("nameInput"); entry != nil {
			if e, ok := entry.(*widget.Entry); ok && e.Text != "" {
				name = e.Text
			}
		}

		message := fmt.Sprintf("Ciao, %s! 👋\n\nBenvenuto nel mondo di Fylay!\n\nQuesta interfaccia è stata creata usando un semplice file XML.", name)
		dialog.ShowInformation("Saluti!", message, a.window)
	}
}

// OnEntryChanged gestisce i cambiamenti nei campi di input
func (a *HelloApp) OnEntryChanged(id, value string) {
	// Non necessario per questo esempio semplice
}

func main() {
	// Crea l'applicazione Fyne
	myApp := app.New()
	window := myApp.NewWindow("Hello Fylay")
	window.Resize(fyne.NewSize(400, 300))

	// Crea l'handler degli eventi
	appHandler := &HelloApp{window: window}

	// Carica il layout XML
	builder, content, err := flay.LoadLayoutFromFileWithHandler("hello.xml", appHandler)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Salva il builder nell'handler
	appHandler.builder = builder

	// Imposta il contenuto e mostra la finestra
	window.SetContent(content)
	window.ShowAndRun()
}
