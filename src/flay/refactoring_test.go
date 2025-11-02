package flay

import (
	"bytes"
	"image/color"
	"testing"
)

// TestStyleCascade verifica la precedenza degli stili (class < id < inline)
func TestStyleCascade(t *testing.T) {
	xml := `
		<Layout>
			<Style selector=".red-text">color: red;</Style>
			<Style selector="#special">color: blue;</Style>
			<VBox class="red-text" id="special" style="color: green;">
				<Label>Test</Label>
			</VBox>
		</Layout>
	`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(bytes.NewReader([]byte(xml)))
	if err != nil {
		t.Fatalf("Failed to parse layout: %v", err)
	}

	// Test cascade: inline > id > class
	elem := layout.Root
	style := builder.computeStyle(elem)

	// Inline style should win
	if style["color"] != "green" {
		t.Errorf("Expected inline style 'green', got '%s'", style["color"])
	}

	// Test without inline style - ID should win
	elem.Style = ""
	style = builder.computeStyle(elem)
	if style["color"] != "blue" {
		t.Errorf("Expected ID style 'blue', got '%s'", style["color"])
	}

	// Test without inline and ID - class should apply
	elem.ID = ""
	style = builder.computeStyle(elem)
	if style["color"] != "red" {
		t.Errorf("Expected class style 'red', got '%s'", style["color"])
	}
}

// TestColorParserStrategies verifies each color parser independently
func TestColorParserStrategies(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected color.Color
	}{
		// Named colors
		{"Named: black", "black", color.Black},
		{"Named: white", "white", color.White},
		{"Named: red", "red", color.RGBA{R: 255, A: 255}},
		{"Named: green", "green", color.RGBA{G: 255, A: 255}},
		{"Named: blue", "blue", color.RGBA{B: 255, A: 255}},
		{"Named: yellow", "yellow", color.RGBA{R: 255, G: 255, A: 255}},
		{"Named: cyan", "cyan", color.RGBA{G: 255, B: 255, A: 255}},
		{"Named: magenta", "magenta", color.RGBA{R: 255, B: 255, A: 255}},

		// Hex 6 digits
		{"Hex6: #FF0000", "#FF0000", color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{"Hex6: #00FF00", "#00FF00", color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		{"Hex6: #0000FF", "#0000FF", color.RGBA{R: 0, G: 0, B: 255, A: 255}},
		{"Hex6: #FFFFFF", "#FFFFFF", color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{"Hex6: #000000", "#000000", color.RGBA{R: 0, G: 0, B: 0, A: 255}},

		// Hex 3 digits
		{"Hex3: #F00", "#F00", color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{"Hex3: #0F0", "#0F0", color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		{"Hex3: #00F", "#00F", color.RGBA{R: 0, G: 0, B: 255, A: 255}},
		{"Hex3: #FFF", "#FFF", color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{"Hex3: #000", "#000", color.RGBA{R: 0, G: 0, B: 0, A: 255}},

		// RGB format
		{"RGB: rgb(255, 0, 0)", "rgb(255, 0, 0)", color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{"RGB: rgb(0, 255, 0)", "rgb(0, 255, 0)", color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		{"RGB: rgb(0, 0, 255)", "rgb(0, 0, 255)", color.RGBA{R: 0, G: 0, B: 255, A: 255}},
		{"RGB: with spaces", "rgb( 128 , 64 , 32 )", color.RGBA{R: 128, G: 64, B: 32, A: 255}},

		// Edge cases
		{"Empty string", "", color.Black},
		{"Whitespace", "  ", color.Black},
		{"Unknown color", "foobar", color.Black},
		{"Invalid hex", "#GGGGGG", color.Black},
		{"Invalid RGB", "rgb(300, 0, 0)", color.Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseColor(tt.input)
			if result != tt.expected {
				t.Errorf("parseColor(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNamedColorParser tests the NamedColorParser specifically
func TestNamedColorParser(t *testing.T) {
	parser := &NamedColorParser{}

	tests := []struct {
		input    string
		canParse bool
		expected color.Color
	}{
		{"red", true, color.RGBA{R: 255, A: 255}},
		{"RED", true, color.RGBA{R: 255, A: 255}}, // Case insensitive
		{"  blue  ", true, color.RGBA{B: 255, A: 255}},
		{"#FF0000", false, color.Black},
		{"unknown", false, color.Black},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if parser.CanParse(tt.input) != tt.canParse {
				t.Errorf("CanParse(%q) = %v, want %v", tt.input, parser.CanParse(tt.input), tt.canParse)
			}

			if tt.canParse {
				result, err := parser.Parse(tt.input)
				if err != nil {
					t.Errorf("Parse(%q) returned error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestHexColorParser tests the HexColorParser specifically
func TestHexColorParser(t *testing.T) {
	parser := &HexColorParser{}

	tests := []struct {
		input    string
		canParse bool
		hasError bool
	}{
		{"#FF0000", true, false},
		{"#F00", true, false},
		{"#GGGGGG", true, true}, // Can parse but returns error
		{"#12", false, false},
		{"red", false, false},
		{"FF0000", false, false}, // No #
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if parser.CanParse(tt.input) != tt.canParse {
				t.Errorf("CanParse(%q) = %v, want %v", tt.input, parser.CanParse(tt.input), tt.canParse)
			}

			if tt.canParse {
				_, err := parser.Parse(tt.input)
				if (err != nil) != tt.hasError {
					t.Errorf("Parse(%q) error = %v, wantError %v", tt.input, err, tt.hasError)
				}
			}
		})
	}
}

// TestRGBColorParser tests the RGBColorParser specifically
func TestRGBColorParser(t *testing.T) {
	parser := &RGBColorParser{}

	tests := []struct {
		input    string
		canParse bool
		hasError bool
	}{
		{"rgb(255, 0, 0)", true, false},
		{"rgb(0, 0, 0)", true, false},
		{"rgb(128, 64, 32)", true, false},
		{"rgb(300, 0, 0)", true, true},    // Out of range
		{"rgb(255, 0)", true, true},       // Too few values
		{"rgb(255, 0, 0, 0)", true, true}, // Too many values
		{"#FF0000", false, false},
		{"red", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if parser.CanParse(tt.input) != tt.canParse {
				t.Errorf("CanParse(%q) = %v, want %v", tt.input, parser.CanParse(tt.input), tt.canParse)
			}

			if tt.canParse {
				_, err := parser.Parse(tt.input)
				if (err != nil) != tt.hasError {
					t.Errorf("Parse(%q) error = %v, wantError %v", tt.input, err, tt.hasError)
				}
			}
		})
	}
}

// TestApplyStyleMethods tests individual style application methods
func TestApplyStyleMethods(t *testing.T) {
	builder := NewBuilder()

	// Setup test styles
	builder.styles = map[string]Style{
		".red": {
			Selector:   ".red",
			Properties: map[string]string{"color": "red", "font-size": "14"},
		},
		".large": {
			Selector:   ".large",
			Properties: map[string]string{"font-size": "20"},
		},
		"#unique": {
			Selector:   "#unique",
			Properties: map[string]string{"color": "blue"},
		},
	}

	t.Run("applyClassStyles - single class", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyClassStyles(&style, "red")

		if style["color"] != "red" {
			t.Errorf("Expected color 'red', got '%s'", style["color"])
		}
		if style["font-size"] != "14" {
			t.Errorf("Expected font-size '14', got '%s'", style["font-size"])
		}
	})

	t.Run("applyClassStyles - multiple classes", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyClassStyles(&style, "red large")

		// 'large' should override 'red' font-size
		if style["color"] != "red" {
			t.Errorf("Expected color 'red', got '%s'", style["color"])
		}
		if style["font-size"] != "20" {
			t.Errorf("Expected font-size '20', got '%s'", style["font-size"])
		}
	})

	t.Run("applyClassStyles - empty", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyClassStyles(&style, "")

		if len(style) != 0 {
			t.Errorf("Expected empty style, got %v", style)
		}
	})

	t.Run("applyIDStyles", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyIDStyles(&style, "unique")

		if style["color"] != "blue" {
			t.Errorf("Expected color 'blue', got '%s'", style["color"])
		}
	})

	t.Run("applyIDStyles - empty", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyIDStyles(&style, "")

		if len(style) != 0 {
			t.Errorf("Expected empty style, got %v", style)
		}
	})

	t.Run("applyInlineStyles", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyInlineStyles(&style, "color: green; font-weight: bold")

		if style["color"] != "green" {
			t.Errorf("Expected color 'green', got '%s'", style["color"])
		}
		if style["font-weight"] != "bold" {
			t.Errorf("Expected font-weight 'bold', got '%s'", style["font-weight"])
		}
	})

	t.Run("applyInlineStyles - empty", func(t *testing.T) {
		style := make(map[string]string)
		builder.applyInlineStyles(&style, "")

		if len(style) != 0 {
			t.Errorf("Expected empty style, got %v", style)
		}
	})
}
