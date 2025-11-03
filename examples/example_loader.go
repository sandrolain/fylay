package main

import (
	"fmt"
	"log"

	flay "github.com/sandrolain/fylay"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// ExampleHandler gestisce gli eventi dell'applicazione
type ExampleHandler struct {
	builder *flay.Builder
	window  fyne.Window
}

func (h *ExampleHandler) OnButtonTapped(id string) {
	log.Printf("Pulsante premuto: %s", id)

	switch id {
	case "loginBtn":
		// Ottieni username e password
		username := h.getEntryText("usernameField")
		password := h.getEntryText("passwordField")

		if username != "" && password != "" {
			log.Printf("Login con: %s / %s", username, password)
			h.showDialog("Login effettuato!", fmt.Sprintf("Benvenuto %s!", username))
		} else {
			h.showDialog("Errore", "Inserisci username e password")
		}

	case "cancelBtn":
		log.Println("Login annullato")
		h.clearEntry("usernameField")
		h.clearEntry("passwordField")

	case "saveBtn":
		// Form di registrazione
		firstName := h.getEntryText("firstName")
		lastName := h.getEntryText("lastName")
		email := h.getEntryText("email")

		msg := fmt.Sprintf("Dati salvati:\nNome: %s\nCognome: %s\nEmail: %s",
			firstName, lastName, email)
		h.showDialog("Salvataggio", msg)

	case "resetBtn":
		// Reset form
		h.clearEntry("firstName")
		h.clearEntry("lastName")
		h.clearEntry("email")
		h.clearEntry("phone")
		h.clearEntry("street")
		h.clearEntry("city")
		h.clearEntry("zip")
		h.clearEntry("notes")

	case "menuHome", "menuStats", "menuSettings":
		h.showDialog("Menu", fmt.Sprintf("Navigazione a: %s", id))
	}
}

func (h *ExampleHandler) OnEntryChanged(id, value string) {
	log.Printf("Campo '%s' modificato: %s", id, value)
}

// Helper methods
func (h *ExampleHandler) getEntryText(id string) string {
	if element := h.builder.GetElement(id); element != nil {
		if entry, ok := element.(*widget.Entry); ok {
			return entry.Text
		}
	}
	return ""
}

func (h *ExampleHandler) clearEntry(id string) {
	if element := h.builder.GetElement(id); element != nil {
		if entry, ok := element.(*widget.Entry); ok {
			entry.SetText("")
		}
	}
}

func (h *ExampleHandler) showDialog(title, message string) {
	dialog := widget.NewModalPopUp(
		widget.NewLabel(message),
		h.window.Canvas(),
	)
	dialog.Show()
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("github.com/sandrolain/fylay Examples")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Scegli quale esempio caricare
	// Cambia il path per testare diversi esempi
	examples := []string{
		"examples/login.xml",
		"examples/dashboard.xml",
		"examples/form.xml",
	}

	// Carica il primo esempio (login)
	examplePath := examples[0]

	handler := &ExampleHandler{window: myWindow}

	builder, content, err := flay.LoadLayoutFromFileWithHandler(examplePath, handler)
	if err != nil {
		log.Fatalf("Errore nel caricamento del layout: %v", err)
	}

	handler.builder = builder

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
