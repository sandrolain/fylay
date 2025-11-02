package main

import (
	"fylay/src/flay"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// AppEventHandler gestisce gli eventi dell'applicazione
type AppEventHandler struct {
	builder *flay.Builder
}

func (h *AppEventHandler) OnButtonTapped(id string) {
	log.Printf("Pulsante premuto: %s", id)

	switch id {
	case "submitBtn":
		// Ottieni l'entry username
		if entry := h.builder.GetElement("username"); entry != nil {
			if e, ok := entry.(*widget.Entry); ok {
				log.Printf("Username: %s", e.Text)
			}
		}
	case "clearBtn":
		if entry := h.builder.GetElement("username"); entry != nil {
			if e, ok := entry.(*widget.Entry); ok {
				e.SetText("")
			}
		}
	}
}

func (h *AppEventHandler) OnEntryChanged(id, value string) {
	log.Printf("Entry '%s' cambiata: %s", id, value)
}

func main() {
	// Crea l'applicazione
	myApp := app.New()
	myWindow := myApp.NewWindow("Fylay Demo")
	myWindow.Resize(fyne.NewSize(600, 400))

	// Layout XML di esempio
	layoutXML := `
<Layout>
	<Style selector=".header">
		font-weight: bold;
		text-align: center;
		font-size: 20;
	</Style>
	
	<Style selector=".primary-btn">
		background-color: #0066cc;
	</Style>
	
	<Style selector="#username">
		width: 300;
	</Style>
	
	<VBox>
		<Label class="header">Benvenuto in Fylay</Label>
		
		<Rectangle style="background-color: #cccccc; height: 2;" />
		
		<VBox style="padding: 20;">
			<Label style="font-weight: bold;">Inserisci i tuoi dati:</Label>
			
			<Entry id="username" placeholder="Username" />
			
			<Entry id="password" placeholder="Password" password="true" />
			
			<HBox>
				<Button id="submitBtn" class="primary-btn">Invia</Button>
				<Button id="clearBtn">Cancella</Button>
			</HBox>
		</VBox>
		
		<Border>
			<VBox position="center">
				<Grid columns="2">
					<Label style="font-weight: bold;">Nome:</Label>
					<Label>Esempio</Label>
					
					<Label style="font-weight: bold;">Età:</Label>
					<Label>25</Label>
				</Grid>
			</VBox>
		</Border>
	</VBox>
</Layout>
`

	// Crea il builder
	builder := flay.NewBuilder()

	// Imposta l'event handler
	handler := &AppEventHandler{builder: builder}
	builder.SetEventHandler(handler)

	// Carica il layout
	layout, err := builder.LoadLayout(strings.NewReader(layoutXML))
	if err != nil {
		log.Fatalf("Errore nel caricamento del layout: %v", err)
	}

	// Costruisci l'interfaccia
	content, err := builder.Build(layout)
	if err != nil {
		log.Fatalf("Errore nella costruzione dell'interfaccia: %v", err)
	}

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
