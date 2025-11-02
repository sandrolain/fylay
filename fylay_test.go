package fylay

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func TestParseCSS(t *testing.T) {
	tests := []struct {
		name     string
		css      string
		expected map[string]string
	}{
		{
			name: "Simple CSS",
			css:  "font-size: 16; color: red;",
			expected: map[string]string{
				"font-size": "16",
				"color":     "red",
			},
		},
		{
			name: "CSS with braces",
			css:  "{ font-weight: bold; text-align: center; }",
			expected: map[string]string{
				"font-weight": "bold",
				"text-align":  "center",
			},
		},
		{
			name: "CSS with spaces",
			css:  "  background-color: #FF0000  ;  width: 300  ",
			expected: map[string]string{
				"background-color": "#FF0000",
				"width":            "300",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCSS(tt.css)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d properties, got %d", len(tt.expected), len(result))
			}

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("For key %s: expected %s, got %s", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		r, g, b uint8
	}{
		{"Named color - red", "red", 255, 0, 0},
		{"Named color - blue", "blue", 0, 0, 255},
		{"Named color - green", "green", 0, 255, 0},
		{"Hex 6 digit", "#FF0000", 255, 0, 0},
		{"Hex 3 digit", "#F00", 255, 0, 0},
		{"RGB format", "rgb(128, 64, 32)", 128, 64, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := parseColor(tt.input)
			r, g, b, _ := color.RGBA()

			// RGBA returns uint32 values scaled to 0-65535, need to scale back
			//nolint:gosec // Intentional conversion for test
			actualR := uint8(r >> 8)
			//nolint:gosec // Intentional conversion for test
			actualG := uint8(g >> 8)
			//nolint:gosec // Intentional conversion for test
			actualB := uint8(b >> 8)

			if actualR != tt.r || actualG != tt.g || actualB != tt.b {
				t.Errorf("Expected RGB(%d, %d, %d), got RGB(%d, %d, %d)",
					tt.r, tt.g, tt.b, actualR, actualG, actualB)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float32
		hasError bool
	}{
		{"Simple number", "100", 100, false},
		{"With px", "200px", 200, false},
		{"Float", "15.5", 15.5, false},
		{"With spaces", "  50  ", 50, false},
		{"Invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseSize(tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %f, got %f", tt.expected, result)
				}
			}
		})
	}
}

func TestLoadLayout(t *testing.T) {
	layoutXML := `
<Layout>
	<Style selector=".title">
		font-size: 20;
		font-weight: bold;
	</Style>
	
	<VBox>
		<Label class="title">Test Title</Label>
		<Button id="testBtn">Click Me</Button>
	</VBox>
</Layout>
`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(strings.NewReader(layoutXML))

	if err != nil {
		t.Fatalf("Failed to load layout: %v", err)
	}

	if layout == nil {
		t.Fatal("Layout is nil")
	}

	// Verifica che gli stili siano stati parsati
	if len(builder.styles) == 0 {
		t.Error("Expected styles to be parsed")
	}

	if style, ok := builder.styles[".title"]; ok {
		if style.Properties["font-size"] != "20" {
			t.Errorf("Expected font-size 20, got %s", style.Properties["font-size"])
		}
	} else {
		t.Error("Style .title not found")
	}
}

func TestBuildLayout(t *testing.T) {
	layoutXML := `
<Layout>
	<VBox>
		<Label id="label1">Test Label</Label>
		<Button id="button1">Test Button</Button>
		<Entry id="entry1" placeholder="Test" />
	</VBox>
</Layout>
`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(strings.NewReader(layoutXML))
	if err != nil {
		t.Fatalf("Failed to load layout: %v", err)
	}

	content, err := builder.Build(layout)
	if err != nil {
		t.Fatalf("Failed to build layout: %v", err)
	}

	if content == nil {
		t.Fatal("Content is nil")
	}

	// Verifica che il contenuto sia un VBox
	if _, ok := content.(*fyne.Container); !ok {
		t.Error("Expected container")
	}

	// Verifica che gli elementi siano stati registrati
	if builder.GetElement("label1") == nil {
		t.Error("Element label1 not found")
	}

	if builder.GetElement("button1") == nil {
		t.Error("Element button1 not found")
	}

	if builder.GetElement("entry1") == nil {
		t.Error("Element entry1 not found")
	}

	// Verifica i tipi
	if _, ok := builder.GetElement("label1").(*widget.Label); !ok {
		t.Error("label1 is not a Label")
	}

	if _, ok := builder.GetElement("button1").(*widget.Button); !ok {
		t.Error("button1 is not a Button")
	}

	if _, ok := builder.GetElement("entry1").(*widget.Entry); !ok {
		t.Error("entry1 is not an Entry")
	}
}

func TestStylePrecedence(t *testing.T) {
	layoutXML := `
<Layout>
	<Style selector=".myClass">
		font-size: 10;
	</Style>
	
	<Style selector="#myElement">
		font-size: 20;
	</Style>
	
	<VBox>
		<Label id="myElement" class="myClass" style="font-size: 30;">Test</Label>
	</VBox>
</Layout>
`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(strings.NewReader(layoutXML))
	if err != nil {
		t.Fatalf("Failed to load layout: %v", err)
	}

	// Costruisci per testare compute style
	_, err = builder.Build(layout)
	if err != nil {
		t.Fatalf("Failed to build: %v", err)
	}

	// Il test verifica che il build funzioni - lo stile inline dovrebbe avere precedenza
}

func TestEventHandler(t *testing.T) {
	layoutXML := `
<Layout>
	<VBox>
		<Button id="testBtn">Test</Button>
		<Entry id="testEntry" />
	</VBox>
</Layout>
`

	handler := &testEventHandler{}

	builder := NewBuilder()
	builder.SetEventHandler(handler)

	layout, err := builder.LoadLayout(strings.NewReader(layoutXML))
	if err != nil {
		t.Fatalf("Failed to load layout: %v", err)
	}

	_, err = builder.Build(layout)
	if err != nil {
		t.Fatalf("Failed to build: %v", err)
	}

	// Verifica che l'handler sia stato impostato
	if builder.eventHandler == nil {
		t.Error("Event handler not set")
	}
}

// Implementazione dummy dell'EventHandler per i test
type testEventHandler struct {
	buttonTapped bool
	entryChanged bool
	lastEntryID  string
	lastValue    string
}

func (h *testEventHandler) OnButtonTapped(id string) {
	h.buttonTapped = true
}

func (h *testEventHandler) OnEntryChanged(id, value string) {
	h.entryChanged = true
	h.lastEntryID = id
	h.lastValue = value
}
