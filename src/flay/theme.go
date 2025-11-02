package flay

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"gopkg.in/yaml.v3"
)

// ThemeConfig rappresenta la configurazione del tema da YAML
type ThemeConfig struct {
	Name    string             `yaml:"name"`
	Variant string             `yaml:"variant"` // "light" or "dark"
	Colors  ThemeColors        `yaml:"colors"`
	Fonts   ThemeFonts         `yaml:"fonts"`
	Sizes   map[string]float32 `yaml:"sizes"`
}

// ThemeColors contiene i colori del tema
type ThemeColors struct {
	Primary     string `yaml:"primary"`
	Background  string `yaml:"background"`
	Foreground  string `yaml:"foreground"`
	Button      string `yaml:"button"`
	Disabled    string `yaml:"disabled"`
	Error       string `yaml:"error"`
	Focus       string `yaml:"focus"`
	Hover       string `yaml:"hover"`
	Input       string `yaml:"input"`
	Placeholder string `yaml:"placeholder"`
	Pressed     string `yaml:"pressed"`
	ScrollBar   string `yaml:"scrollbar"`
	Shadow      string `yaml:"shadow"`
}

// ThemeFonts contiene le configurazioni dei font
type ThemeFonts struct {
	Regular   string  `yaml:"regular"`
	Bold      string  `yaml:"bold"`
	Italic    string  `yaml:"italic"`
	Monospace string  `yaml:"monospace"`
	Size      float32 `yaml:"size"`
}

// CustomTheme implementa fyne.Theme
type CustomTheme struct {
	config *ThemeConfig
	base   fyne.Theme
	colors map[string]color.Color
}

// LoadThemeFromYAML loads a theme from a YAML file
func LoadThemeFromYAML(filepath string) (fyne.Theme, error) {
	data, err := os.ReadFile(filepath) //nolint:gosec // Theme path is from user config, intentional
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var config ThemeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse theme YAML: %w", err)
	}

	return NewCustomTheme(&config), nil
}

// NewCustomTheme crea un nuovo tema custom
func NewCustomTheme(config *ThemeConfig) *CustomTheme {
	ct := &CustomTheme{
		config: config,
		colors: make(map[string]color.Color),
	}

	// Usa il tema base di Fyne (light o dark)
	// Note: Using deprecated methods until Fyne v3 provides alternative
	if config.Variant == "dark" {
		ct.base = theme.DarkTheme() //nolint:staticcheck // Will update when Fyne v3 is released
	} else {
		ct.base = theme.LightTheme() //nolint:staticcheck // Will update when Fyne v3 is released
	}

	// Parse colors
	ct.parseColors()

	return ct
}

// parseColors converte le stringhe colore in color.Color
func (ct *CustomTheme) parseColors() {
	if ct.config.Colors.Primary != "" {
		ct.colors["primary"] = parseColor(ct.config.Colors.Primary)
	}
	if ct.config.Colors.Background != "" {
		ct.colors["background"] = parseColor(ct.config.Colors.Background)
	}
	if ct.config.Colors.Foreground != "" {
		ct.colors["foreground"] = parseColor(ct.config.Colors.Foreground)
	}
	if ct.config.Colors.Button != "" {
		ct.colors["button"] = parseColor(ct.config.Colors.Button)
	}
	if ct.config.Colors.Disabled != "" {
		ct.colors["disabled"] = parseColor(ct.config.Colors.Disabled)
	}
	if ct.config.Colors.Error != "" {
		ct.colors["error"] = parseColor(ct.config.Colors.Error)
	}
	if ct.config.Colors.Focus != "" {
		ct.colors["focus"] = parseColor(ct.config.Colors.Focus)
	}
	if ct.config.Colors.Hover != "" {
		ct.colors["hover"] = parseColor(ct.config.Colors.Hover)
	}
	if ct.config.Colors.Input != "" {
		ct.colors["input"] = parseColor(ct.config.Colors.Input)
	}
	if ct.config.Colors.Placeholder != "" {
		ct.colors["placeholder"] = parseColor(ct.config.Colors.Placeholder)
	}
	if ct.config.Colors.Pressed != "" {
		ct.colors["pressed"] = parseColor(ct.config.Colors.Pressed)
	}
	if ct.config.Colors.ScrollBar != "" {
		ct.colors["scrollbar"] = parseColor(ct.config.Colors.ScrollBar)
	}
	if ct.config.Colors.Shadow != "" {
		ct.colors["shadow"] = parseColor(ct.config.Colors.Shadow)
	}
}

// Color implementa fyne.Theme
func (ct *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Check custom colors first
	switch name {
	case theme.ColorNamePrimary:
		if c, ok := ct.colors["primary"]; ok {
			return c
		}
	case theme.ColorNameBackground:
		if c, ok := ct.colors["background"]; ok {
			return c
		}
	case theme.ColorNameForeground:
		if c, ok := ct.colors["foreground"]; ok {
			return c
		}
	case theme.ColorNameButton:
		if c, ok := ct.colors["button"]; ok {
			return c
		}
	case theme.ColorNameDisabled:
		if c, ok := ct.colors["disabled"]; ok {
			return c
		}
	case theme.ColorNameError:
		if c, ok := ct.colors["error"]; ok {
			return c
		}
	case theme.ColorNameFocus:
		if c, ok := ct.colors["focus"]; ok {
			return c
		}
	case theme.ColorNameHover:
		if c, ok := ct.colors["hover"]; ok {
			return c
		}
	case theme.ColorNameInputBackground:
		if c, ok := ct.colors["input"]; ok {
			return c
		}
	case theme.ColorNamePlaceHolder:
		if c, ok := ct.colors["placeholder"]; ok {
			return c
		}
	case theme.ColorNamePressed:
		if c, ok := ct.colors["pressed"]; ok {
			return c
		}
	case theme.ColorNameScrollBar:
		if c, ok := ct.colors["scrollbar"]; ok {
			return c
		}
	case theme.ColorNameShadow:
		if c, ok := ct.colors["shadow"]; ok {
			return c
		}
	}

	// Fallback to base theme
	return ct.base.Color(name, variant)
}

// Font implementa fyne.Theme
func (ct *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Use base theme fonts for now
	// Future: support custom fonts from config
	return ct.base.Font(style)
}

// Icon implementa fyne.Theme
func (ct *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return ct.base.Icon(name)
}

// Size implementa fyne.Theme
func (ct *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	// Check custom sizes
	if ct.config.Sizes != nil {
		if size, ok := ct.config.Sizes[string(name)]; ok {
			return size
		}
	}

	// Check font size
	if name == theme.SizeNameText && ct.config.Fonts.Size > 0 {
		return ct.config.Fonts.Size
	}

	// Fallback to base theme
	return ct.base.Size(name)
}

// ApplyThemeToApp applica il tema all'applicazione Fyne
func ApplyThemeToApp(app fyne.App, themePath string) error {
	customTheme, err := LoadThemeFromYAML(themePath)
	if err != nil {
		return err
	}

	app.Settings().SetTheme(customTheme)
	return nil
}
