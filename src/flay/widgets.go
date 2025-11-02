package flay

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// buildCheckbox costruisce un widget Checkbox
func (b *Builder) buildCheckbox(elem Element, style map[string]string) fyne.CanvasObject {
	label := elem.getAttr("label")
	if label == "" {
		label = elem.Content
	}

	checked := elem.getAttr("checked") == attrValueTrue

	check := widget.NewCheck(label, nil)
	check.Checked = checked

	// Handle onchange event
	onchange := elem.getAttr("onchange")
	if onchange != "" {
		if callback, ok := b.eventCallbacks[onchange]; ok {
			check.OnChanged = func(bool) {
				callback()
			}
		}
	}

	// Handle data binding
	if bindAttr := elem.getAttr("bind"); bindAttr != "" {
		ctx := b.GetBindingContext()
		key := ParseBindAttribute(bindAttr)
		boolData := ctx.BindBool(key, checked)
		check.Bind(boolData)
	}

	// Register widget
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, check)
	}

	return check
}

// buildSelect costruisce un widget Select (dropdown)
func (b *Builder) buildSelect(elem Element, style map[string]string) fyne.CanvasObject {
	var options []string
	selected := ""

	// Parse Option children
	for _, child := range elem.Children {
		if child.XMLName.Local == "Option" {
			value := child.getAttr("value")
			if value == "" {
				value = child.Content
			}
			options = append(options, value)

			if child.getAttr("selected") == attrValueTrue {
				selected = value
			}
		}
	}

	sel := widget.NewSelect(options, nil)
	if selected != "" {
		sel.SetSelected(selected)
	}

	// Handle onchange event
	onchange := elem.getAttr("onchange")
	if onchange != "" {
		if callback, ok := b.entryCallbacks[onchange]; ok {
			sel.OnChanged = callback
		}
	}

	// Handle data binding
	if bindAttr := elem.getAttr("bind"); bindAttr != "" {
		ctx := b.GetBindingContext()
		key := ParseBindAttribute(bindAttr)
		strData := ctx.BindString(key, selected)
		sel.Bind(strData)
	}

	// Register widget
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, sel)
	}

	return sel
}

// buildProgressBar costruisce un widget ProgressBar
func (b *Builder) buildProgressBar(elem Element, style map[string]string) fyne.CanvasObject {
	value := 0.0
	if v := elem.getAttr("value"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			value = parsed
		}
	}

	max := 1.0
	if m := elem.getAttr("max"); m != "" {
		if parsed, err := strconv.ParseFloat(m, 64); err == nil {
			max = parsed
		}
	}

	progress := widget.NewProgressBar()
	progress.Max = max
	progress.SetValue(value)

	// Handle data binding
	if bindAttr := elem.getAttr("bind"); bindAttr != "" {
		ctx := b.GetBindingContext()
		key := ParseBindAttribute(bindAttr)
		floatData := ctx.BindFloat(key, value)
		progress.Bind(floatData)
	}

	// Register widget
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, progress)
	}

	return progress
}

// buildSlider costruisce un widget Slider
func (b *Builder) buildSlider(elem Element, style map[string]string) fyne.CanvasObject {
	min := 0.0
	if m := elem.getAttr("min"); m != "" {
		if parsed, err := strconv.ParseFloat(m, 64); err == nil {
			min = parsed
		}
	}

	max := 100.0
	if m := elem.getAttr("max"); m != "" {
		if parsed, err := strconv.ParseFloat(m, 64); err == nil {
			max = parsed
		}
	}

	value := min
	if v := elem.getAttr("value"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			value = parsed
		}
	}

	step := 1.0
	if s := elem.getAttr("step"); s != "" {
		if parsed, err := strconv.ParseFloat(s, 64); err == nil {
			step = parsed
		}
	}

	slider := widget.NewSlider(min, max)
	slider.Value = value
	slider.Step = step

	// Handle onchange event
	onchange := elem.getAttr("onchange")
	if onchange != "" {
		if callback, ok := b.eventCallbacks[onchange]; ok {
			slider.OnChanged = func(float64) {
				callback()
			}
		}
	}

	// Handle data binding
	if bindAttr := elem.getAttr("bind"); bindAttr != "" {
		ctx := b.GetBindingContext()
		key := ParseBindAttribute(bindAttr)
		floatData := ctx.BindFloat(key, value)
		slider.Bind(floatData)
	}

	// Register widget
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, slider)
	}

	return slider
}

// buildImage costruisce un widget Image
func (b *Builder) buildImage(elem Element, style map[string]string) fyne.CanvasObject {
	src := elem.getAttr("src")
	if src == "" {
		// Return empty rectangle if no source
		rect := canvas.NewRectangle(nil)
		rect.SetMinSize(fyne.NewSize(100, 100))
		return rect
	}

	img, err := NewImageWidget(src)
	if err != nil {
		// Return error placeholder
		rect := canvas.NewRectangle(nil)
		rect.SetMinSize(fyne.NewSize(100, 100))
		return rect
	}

	// Apply size from attributes or style
	if width := elem.getAttr("width"); width != "" {
		if w, err := parseSize(width); err == nil {
			if height := elem.getAttr("height"); height != "" {
				if h, err := parseSize(height); err == nil {
					img.SetMinSize(fyne.NewSize(w, h))
				}
			}
		}
	}

	// Apply FillMode
	fillMode := elem.getAttr("fillMode")
	switch fillMode {
	case "contain":
		img.FillMode = canvas.ImageFillContain
	case "original":
		img.FillMode = canvas.ImageFillOriginal
	case "stretch":
		img.FillMode = canvas.ImageFillStretch
	default:
		img.FillMode = canvas.ImageFillContain
	}

	return img
}

// buildRadioGroup costruisce un widget RadioGroup
func (b *Builder) buildRadioGroup(elem Element, style map[string]string) fyne.CanvasObject {
	var options []string
	selected := ""

	// Parse Radio children
	for _, child := range elem.Children {
		if child.XMLName.Local == "Radio" {
			value := child.getAttr("value")
			if value == "" {
				value = child.Content
			}
			options = append(options, value)

			if child.getAttr("selected") == attrValueTrue {
				selected = value
			}
		}
	}

	radio := widget.NewRadioGroup(options, nil)
	if selected != "" {
		radio.SetSelected(selected)
	}

	// Handle onchange event
	onchange := elem.getAttr("onchange")
	if onchange != "" {
		if callback, ok := b.entryCallbacks[onchange]; ok {
			radio.OnChanged = callback
		}
	}

	// Note: RadioGroup doesn't support direct binding in Fyne
	// Binding would need to be done through OnChanged callback

	// Register widget
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, radio)
	}

	return radio
}
