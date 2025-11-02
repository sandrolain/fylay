package fylay

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/sandrolain/fylay/core"
	"github.com/sandrolain/fylay/dialog"
)

// Constants for widget types and attributes
const (
	alignLeft     = core.AlignLeft
	alignRight    = core.AlignRight
	alignCenter   = core.AlignCenter
	attrValueTrue = core.AttrValueTrue
)

// Layout rappresenta un layout XML parsato
type Layout struct {
	XMLName xml.Name `xml:"Layout"`
	Styles  []Style  `xml:"Style"`
	Root    Element  `xml:",any"`
}

// Style rappresenta una regola di stile CSS
type Style struct {
	Selector   string            `xml:"selector,attr"`
	Properties map[string]string `xml:"-"`
	RawCSS     string            `xml:",innerxml"`
}

// Element rappresenta un elemento generico del layout
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

// EventContext contiene le informazioni di contesto di un evento
type EventContext struct {
	// EventName è il nome dell'evento (es. "onclick", "onchange")
	EventName string
	// Target è l'elemento che ha generato l'evento
	Target fyne.CanvasObject
	// TargetID è l'ID dell'elemento che ha generato l'evento
	TargetID string
	// Value è il valore associato all'evento (per Entry)
	Value string
}

// ButtonCallback è una funzione callback per eventi di pulsanti
type ButtonCallback func(ctx *EventContext)

// EntryCallback è una funzione callback per eventi Entry
type EntryCallback func(ctx *EventContext)

// Builder costruisce i widget Fyne dal layout
type Builder struct {
	styles          map[string]Style
	elements        map[string]fyne.CanvasObject
	eventHandler    EventHandler
	eventCallbacks  map[string]ButtonCallback
	entryCallbacks  map[string]EntryCallback
	bindingContext  *BindingContext
	templateContext *TemplateContext
}

// EventHandler gestisce gli eventi dei widget
type EventHandler interface {
	OnButtonTapped(id string)
	OnEntryChanged(id, value string)
}

// NewBuilder crea un nuovo builder
func NewBuilder() *Builder {
	return &Builder{
		styles:         make(map[string]Style),
		elements:       make(map[string]fyne.CanvasObject),
		eventCallbacks: make(map[string]ButtonCallback),
		entryCallbacks: make(map[string]EntryCallback),
	}
}

// SetEventHandler imposta l'handler degli eventi
func (b *Builder) SetEventHandler(handler EventHandler) {
	b.eventHandler = handler
}

// On registra una callback per un evento specifico (per pulsanti)
func (b *Builder) On(eventName string, callback ButtonCallback) {
	b.eventCallbacks[eventName] = callback
}

// OnEntry registra una callback per eventi Entry
func (b *Builder) OnEntry(eventName string, callback EntryCallback) {
	b.entryCallbacks[eventName] = callback
}

// LoadLayout carica un layout da un reader XML
func (b *Builder) LoadLayout(r io.Reader) (*Layout, error) {
	var layout Layout
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&layout); err != nil {
		return nil, fmt.Errorf("errore parsing XML: %w", err)
	}

	// Parse CSS styles
	for i := range layout.Styles {
		layout.Styles[i].Properties = parseCSS(layout.Styles[i].RawCSS)
		b.styles[layout.Styles[i].Selector] = layout.Styles[i]
	}

	return &layout, nil
}

// Build costruisce l'interfaccia dal layout
func (b *Builder) Build(layout *Layout) (fyne.CanvasObject, error) {
	return b.buildElement(layout.Root)
}

// GetElement restituisce un elemento per ID
func (b *Builder) GetElement(id string) fyne.CanvasObject {
	return b.elements[id]
}

// buildElement costruisce ricorsivamente un elemento
func (b *Builder) buildElement(elem Element) (fyne.CanvasObject, error) {
	// Calcola lo stile finale (combinando class, id e inline style)
	style := b.computeStyle(elem)

	var obj fyne.CanvasObject
	var err error

	switch elem.XMLName.Local {
	case "VBox":
		obj = b.buildVBox(elem, style)
	case "HBox":
		obj = b.buildHBox(elem, style)
	case "Grid":
		obj = b.buildGrid(elem, style)
	case "Border":
		obj = b.buildBorder(elem, style)
	case "Label":
		obj = b.buildLabel(elem, style)
	case "Button":
		obj = b.buildButton(elem, style)
	case "Entry":
		obj = b.buildEntry(elem, style)
	case "Rectangle":
		obj = b.buildRectangle(elem, style)
	case "Circle":
		obj = b.buildCircle(elem, style)
	case "Text":
		obj = b.buildText(elem, style)
	case "Spacer":
		obj = widget.NewLabel("")
	case "Checkbox":
		obj = b.buildCheckbox(elem, style)
	case "Select":
		obj = b.buildSelect(elem, style)
	case "ProgressBar":
		obj = b.buildProgressBar(elem, style)
	case "Slider":
		obj = b.buildSlider(elem, style)
	case "Image":
		obj = b.buildImage(elem, style)
	case "RadioGroup":
		obj = b.buildRadioGroup(elem, style)
	default:
		return nil, fmt.Errorf("tipo di elemento sconosciuto: %s", elem.XMLName.Local)
	}

	if obj != nil && elem.ID != "" {
		b.elements[elem.ID] = obj
	}

	return obj, err
}

// buildVBox costruisce un container verticale
func (b *Builder) buildVBox(elem Element, style map[string]string) fyne.CanvasObject {
	children := make([]fyne.CanvasObject, 0, len(elem.Children))
	for _, child := range elem.Children {
		if obj, err := b.buildElement(child); err == nil && obj != nil {
			children = append(children, obj)
		}
	}
	return container.NewVBox(children...)
}

// buildHBox costruisce un container orizzontale
func (b *Builder) buildHBox(elem Element, style map[string]string) fyne.CanvasObject {
	children := make([]fyne.CanvasObject, 0, len(elem.Children))
	for _, child := range elem.Children {
		if obj, err := b.buildElement(child); err == nil && obj != nil {
			children = append(children, obj)
		}
	}
	return container.NewHBox(children...)
}

// buildGrid costruisce un layout a griglia
func (b *Builder) buildGrid(elem Element, style map[string]string) fyne.CanvasObject {
	children := make([]fyne.CanvasObject, 0, len(elem.Children))
	for _, child := range elem.Children {
		if obj, err := b.buildElement(child); err == nil && obj != nil {
			children = append(children, obj)
		}
	}

	cols := 2
	if colsStr := elem.getAttr("columns"); colsStr != "" {
		if c, err := strconv.Atoi(colsStr); err == nil {
			cols = c
		}
	}

	return container.NewGridWithColumns(cols, children...)
}

// buildBorder costruisce un border layout
func (b *Builder) buildBorder(elem Element, style map[string]string) fyne.CanvasObject {
	var top, bottom, left, right, center fyne.CanvasObject

	for _, child := range elem.Children {
		obj, err := b.buildElement(child)
		if err != nil || obj == nil {
			continue
		}

		pos := child.getAttr("position")
		switch pos {
		case "top":
			top = obj
		case "bottom":
			bottom = obj
		case "left":
			left = obj
		case "right":
			right = obj
		case "center", "":
			center = obj
		}
	}

	return container.NewBorder(top, bottom, left, right, center)
}

// buildLabel costruisce una label
func (b *Builder) buildLabel(elem Element, style map[string]string) fyne.CanvasObject {
	text := elem.Text
	if text == "" {
		text = strings.TrimSpace(elem.Content)
	}

	label := widget.NewLabel(text)

	if align := style["text-align"]; align != "" {
		switch align {
		case alignCenter:
			label.Alignment = fyne.TextAlignCenter
		case alignRight:
			label.Alignment = fyne.TextAlignTrailing
		case alignLeft:
			label.Alignment = fyne.TextAlignLeading
		}
	}

	if bold := style["font-weight"]; bold == "bold" {
		label.TextStyle.Bold = true
	}

	if italic := style["font-style"]; italic == "italic" {
		label.TextStyle.Italic = true
	}

	return label
}

// buildButton costruisce un pulsante
func (b *Builder) buildButton(elem Element, style map[string]string) fyne.CanvasObject {
	text := elem.Text
	if text == "" {
		text = strings.TrimSpace(elem.Content)
	}

	// Verifica se c'è un attributo onclick
	onclick := elem.getAttr("onclick")

	var btn *widget.Button
	btn = widget.NewButton(text, func() {
		// Prima prova a chiamare la callback registrata
		if onclick != "" {
			if callback, ok := b.eventCallbacks[onclick]; ok {
				ctx := &EventContext{
					EventName: onclick,
					Target:    btn,
					TargetID:  elem.ID,
				}
				callback(ctx)
				return
			}
		}

		// Altrimenti usa l'EventHandler tradizionale
		if b.eventHandler != nil && elem.ID != "" {
			b.eventHandler.OnButtonTapped(elem.ID)
		}
	})

	return btn
}

// buildEntry costruisce un campo di input
func (b *Builder) buildEntry(elem Element, style map[string]string) fyne.CanvasObject {
	entry := widget.NewEntry()

	if placeholder := elem.getAttr("placeholder"); placeholder != "" {
		entry.PlaceHolder = placeholder
	}

	if elem.getAttr("password") == "true" {
		entry.Password = true
	}

	if elem.getAttr("multiline") == "true" {
		entry.MultiLine = true
	}

	// Verifica se c'è un attributo onchange
	onchange := elem.getAttr("onchange")

	entry.OnChanged = func(value string) {
		// Prima prova a chiamare la callback registrata
		if onchange != "" {
			if callback, ok := b.entryCallbacks[onchange]; ok {
				ctx := &EventContext{
					EventName: onchange,
					Target:    entry,
					TargetID:  elem.ID,
					Value:     value,
				}
				callback(ctx)
				return
			}
		}

		// Altrimenti usa l'EventHandler tradizionale
		if b.eventHandler != nil && elem.ID != "" {
			b.eventHandler.OnEntryChanged(elem.ID, value)
		}
	}

	return entry
}

// buildRectangle costruisce un rettangolo
func (b *Builder) buildRectangle(elem Element, style map[string]string) fyne.CanvasObject {
	rect := canvas.NewRectangle(parseColor(style["background-color"]))

	if w := style["width"]; w != "" {
		if width, err := parseSize(w); err == nil {
			rect.SetMinSize(fyne.NewSize(width, rect.MinSize().Height))
		}
	}

	if h := style["height"]; h != "" {
		if height, err := parseSize(h); err == nil {
			rect.SetMinSize(fyne.NewSize(rect.MinSize().Width, height))
		}
	}

	return rect
}

// buildCircle costruisce un cerchio
func (b *Builder) buildCircle(elem Element, style map[string]string) fyne.CanvasObject {
	circle := canvas.NewCircle(parseColor(style["background-color"]))

	if w := style["width"]; w != "" {
		if width, err := parseSize(w); err == nil {
			circle.Resize(fyne.NewSize(width, width))
		}
	}

	return circle
}

// buildText costruisce un testo canvas
func (b *Builder) buildText(elem Element, style map[string]string) fyne.CanvasObject {
	text := elem.Text
	if text == "" {
		text = strings.TrimSpace(elem.Content)
	}

	txt := canvas.NewText(text, parseColor(style["color"]))

	if size := style["font-size"]; size != "" {
		if s, err := parseSize(size); err == nil {
			txt.TextSize = s
		}
	}

	if align := style["text-align"]; align != "" {
		switch align {
		case alignCenter:
			txt.Alignment = fyne.TextAlignCenter
		case alignRight:
			txt.Alignment = fyne.TextAlignTrailing
		case alignLeft:
			txt.Alignment = fyne.TextAlignLeading
		}
	}

	return txt
}

// parseCSS analizza una stringa CSS e restituisce una mappa di proprietà
func parseCSS(css string) map[string]string {
	props := make(map[string]string)

	// Rimuovi eventuali { }
	css = strings.TrimSpace(css)
	css = strings.Trim(css, "{}")

	// Splitta per punto e virgola
	declarations := strings.Split(css, ";")
	for _, decl := range declarations {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}

		parts := strings.SplitN(decl, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			props[key] = value
		}
	}

	return props
}

// parseSize converte una stringa dimensione in float32
func parseSize(sizeStr string) (float32, error) {
	sizeStr = strings.TrimSpace(sizeStr)

	// Rimuovi "px" se presente
	sizeStr = strings.TrimSuffix(sizeStr, "px")

	val, err := strconv.ParseFloat(sizeStr, 32)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}

// getAttr ottiene un attributo dall'elemento
func (e *Element) getAttr(name string) string {
	for _, attr := range e.Attributes {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// ShowInfoDialog mostra un dialog informativo
func (b *Builder) ShowInfoDialog(title, message string) {
	dialog.Info(title, message)
}

// ShowErrorDialog mostra un dialog di errore
func (b *Builder) ShowErrorDialog(title, message string) {
	dialog.Error(title, message)
}

// ShowQuestionDialog mostra un dialog con domanda Yes/No
func (b *Builder) ShowQuestionDialog(title, message string) bool {
	return dialog.Question(title, message)
}

// ShowFileOpenDialog mostra un dialog per aprire file
func (b *Builder) ShowFileOpenDialog(title string, filters ...dialog.FileFilter) (string, error) {
	return dialog.FileOpen(title, filters...)
}

// ShowFileSaveDialog mostra un dialog per salvare file
func (b *Builder) ShowFileSaveDialog(title string, filters ...dialog.FileFilter) (string, error) {
	return dialog.FileSave(title, filters...)
}

// ShowDirSelectDialog mostra un dialog per selezionare directory
func (b *Builder) ShowDirSelectDialog(title string) (string, error) {
	return dialog.DirSelect(title)
}
