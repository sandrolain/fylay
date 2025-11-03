package fylay

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
			check.OnChanged = func(checked bool) {
				ctx := &EventContext{
					EventName: onchange,
					Target:    check,
					TargetID:  elem.ID,
					Value:     strconv.FormatBool(checked),
				}
				callback(ctx)
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

	// Register widget with ID before applying styles
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, check)
		b.widgets[elem.ID] = check // Original widget
	}

	// Apply common styles (width, height) - may wrap in container
	styled := applyMinSize(check, style)

	// Store the final styled version
	if elem.ID != "" {
		b.elements[elem.ID] = styled
	}

	return styled
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
			sel.OnChanged = func(value string) {
				ctx := &EventContext{
					EventName: onchange,
					Target:    sel,
					TargetID:  elem.ID,
					Value:     value,
				}
				callback(ctx)
			}
		}
	}

	// Handle data binding
	if bindAttr := elem.getAttr("bind"); bindAttr != "" {
		ctx := b.GetBindingContext()
		key := ParseBindAttribute(bindAttr)
		strData := ctx.BindString(key, selected)
		sel.Bind(strData)
	}

	// Register widget with ID before applying styles
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, sel)
		b.widgets[elem.ID] = sel // Original widget
	}

	// Apply common styles (width, height) - may wrap in container
	styled := applyMinSize(sel, style)

	// Store the final styled version
	if elem.ID != "" {
		b.elements[elem.ID] = styled
	}

	return styled
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

	// Register widget with ID before applying styles
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, progress)
		b.widgets[elem.ID] = progress // Original widget
	}

	// Apply common styles (width, height) - may wrap in container
	styled := applyMinSize(progress, style)

	// Store the final styled version
	if elem.ID != "" {
		b.elements[elem.ID] = styled
	}

	return styled
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
			slider.OnChanged = func(value float64) {
				ctx := &EventContext{
					EventName: onchange,
					Target:    slider,
					TargetID:  elem.ID,
					Value:     strconv.FormatFloat(value, 'f', -1, 64),
				}
				callback(ctx)
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

	// Register widget with ID before applying styles
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, slider)
		b.widgets[elem.ID] = slider // Original widget
	}

	// Apply common styles (width, height) - may wrap in container
	styled := applyMinSize(slider, style)

	// Store the final styled version
	if elem.ID != "" {
		b.elements[elem.ID] = styled
	}

	return styled
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

	// Register widget with ID before applying styles
	if elem.ID != "" {
		b.widgets[elem.ID] = img // Original widget
	}

	// Apply common styles (width, height) if not already set by attributes - may wrap in container
	var styled fyne.CanvasObject = img
	if elem.getAttr("width") == "" || elem.getAttr("height") == "" {
		styled = applyMinSize(img, style)
	}

	// Store the final styled version
	if elem.ID != "" {
		b.elements[elem.ID] = styled
	}

	return styled
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
			radio.OnChanged = func(value string) {
				ctx := &EventContext{
					EventName: onchange,
					Target:    radio,
					TargetID:  elem.ID,
					Value:     value,
				}
				callback(ctx)
			}
		}
	}

	// Note: RadioGroup doesn't support direct binding in Fyne
	// Binding would need to be done through OnChanged callback

	// Register widget with ID before applying styles
	if elem.ID != "" {
		b.GetBindingContext().RegisterWidget(elem.ID, radio)
		b.widgets[elem.ID] = radio // Original widget
	}

	// Apply common styles (width, height) - may wrap in container
	styled := applyMinSize(radio, style)

	// Store the final styled version
	if elem.ID != "" {
		b.elements[elem.ID] = styled
	}

	return styled
}
