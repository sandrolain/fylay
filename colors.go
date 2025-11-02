package fylay

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// ColorParser defines the interface for color parsing strategies
type ColorParser interface {
	CanParse(colorStr string) bool
	Parse(colorStr string) (color.Color, error)
}

// NamedColorParser handles named colors (black, white, red, etc.)
type NamedColorParser struct{}

func (p *NamedColorParser) CanParse(colorStr string) bool {
	normalized := strings.ToLower(strings.TrimSpace(colorStr))
	namedColors := []string{"black", "white", "red", "green", "blue", "yellow", "cyan", "magenta"}
	for _, name := range namedColors {
		if normalized == name {
			return true
		}
	}
	return false
}

func (p *NamedColorParser) Parse(colorStr string) (color.Color, error) {
	switch strings.ToLower(strings.TrimSpace(colorStr)) {
	case "black":
		return color.Black, nil
	case "white":
		return color.White, nil
	case "red":
		return color.RGBA{R: 255, A: 255}, nil
	case "green":
		return color.RGBA{G: 255, A: 255}, nil
	case "blue":
		return color.RGBA{B: 255, A: 255}, nil
	case "yellow":
		return color.RGBA{R: 255, G: 255, A: 255}, nil
	case "cyan":
		return color.RGBA{G: 255, B: 255, A: 255}, nil
	case "magenta":
		return color.RGBA{R: 255, B: 255, A: 255}, nil
	default:
		return color.Black, fmt.Errorf("unknown named color: %s", colorStr)
	}
}

// HexColorParser handles hex colors (#RGB and #RRGGBB)
type HexColorParser struct{}

func (p *HexColorParser) CanParse(colorStr string) bool {
	colorStr = strings.TrimSpace(colorStr)
	if !strings.HasPrefix(colorStr, "#") {
		return false
	}
	hex := colorStr[1:]
	return len(hex) == 3 || len(hex) == 6
}

func (p *HexColorParser) Parse(colorStr string) (color.Color, error) {
	colorStr = strings.TrimSpace(colorStr)
	if !strings.HasPrefix(colorStr, "#") {
		return color.Black, fmt.Errorf("hex color must start with #")
	}

	hex := colorStr[1:]
	var r, g, b uint8

	if len(hex) == 3 {
		// #RGB format - expand to #RRGGBB
		if _, err := fmt.Sscanf(hex[0:1]+hex[0:1], "%02x", &r); err != nil {
			return color.Black, fmt.Errorf("invalid red component: %w", err)
		}
		if _, err := fmt.Sscanf(hex[1:2]+hex[1:2], "%02x", &g); err != nil {
			return color.Black, fmt.Errorf("invalid green component: %w", err)
		}
		if _, err := fmt.Sscanf(hex[2:3]+hex[2:3], "%02x", &b); err != nil {
			return color.Black, fmt.Errorf("invalid blue component: %w", err)
		}
	} else if len(hex) == 6 {
		// #RRGGBB format
		if _, err := fmt.Sscanf(hex[0:2], "%02x", &r); err != nil {
			return color.Black, fmt.Errorf("invalid red component: %w", err)
		}
		if _, err := fmt.Sscanf(hex[2:4], "%02x", &g); err != nil {
			return color.Black, fmt.Errorf("invalid green component: %w", err)
		}
		if _, err := fmt.Sscanf(hex[4:6], "%02x", &b); err != nil {
			return color.Black, fmt.Errorf("invalid blue component: %w", err)
		}
	} else {
		return color.Black, fmt.Errorf("invalid hex color length: %d", len(hex))
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// RGBColorParser handles rgb(r, g, b) format
type RGBColorParser struct{}

func (p *RGBColorParser) CanParse(colorStr string) bool {
	colorStr = strings.TrimSpace(colorStr)
	return strings.HasPrefix(colorStr, "rgb(") && strings.HasSuffix(colorStr, ")")
}

func (p *RGBColorParser) Parse(colorStr string) (color.Color, error) {
	colorStr = strings.TrimSpace(colorStr)
	if !strings.HasPrefix(colorStr, "rgb(") || !strings.HasSuffix(colorStr, ")") {
		return color.Black, fmt.Errorf("invalid rgb format")
	}

	// Extract values
	values := strings.TrimPrefix(colorStr, "rgb(")
	values = strings.TrimSuffix(values, ")")
	parts := strings.Split(values, ",")

	if len(parts) != 3 {
		return color.Black, fmt.Errorf("rgb requires 3 values, got %d", len(parts))
	}

	var r, g, b uint8

	if v, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil && v >= 0 && v <= 255 {
		r = uint8(v) //nolint:gosec // Already validated v <= 255
	} else {
		return color.Black, fmt.Errorf("invalid red value: %s", parts[0])
	}

	if v, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil && v >= 0 && v <= 255 {
		g = uint8(v) //nolint:gosec // Already validated v <= 255
	} else {
		return color.Black, fmt.Errorf("invalid green value: %s", parts[1])
	}

	if v, err := strconv.Atoi(strings.TrimSpace(parts[2])); err == nil && v >= 0 && v <= 255 {
		b = uint8(v) //nolint:gosec // Already validated v <= 255
	} else {
		return color.Black, fmt.Errorf("invalid blue value: %s", parts[2])
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// Global color parsers registry
var colorParsers = []ColorParser{
	&NamedColorParser{},
	&HexColorParser{},
	&RGBColorParser{},
}

// parseColor converts a color string to color.Color using strategy pattern
func parseColor(colorStr string) color.Color {
	colorStr = strings.TrimSpace(colorStr)

	if colorStr == "" {
		return color.Black
	}

	// Try each parser in order
	for _, parser := range colorParsers {
		if parser.CanParse(colorStr) {
			if c, err := parser.Parse(colorStr); err == nil {
				return c
			}
		}
	}

	// Fallback to black if no parser succeeds
	return color.Black
}
