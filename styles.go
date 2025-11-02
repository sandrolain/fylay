package fylay

import "strings"

// applyClassStyles applies CSS class styles to the style map
func (b *Builder) applyClassStyles(style *map[string]string, classes string) {
	if classes == "" {
		return
	}

	classList := strings.Split(classes, " ")
	for _, class := range classList {
		class = strings.TrimSpace(class)
		if class == "" {
			continue
		}

		selector := "." + class
		if s, ok := b.styles[selector]; ok {
			for k, v := range s.Properties {
				(*style)[k] = v
			}
		}
	}
}

// applyIDStyles applies CSS ID styles to the style map
func (b *Builder) applyIDStyles(style *map[string]string, id string) {
	if id == "" {
		return
	}

	selector := "#" + id
	if s, ok := b.styles[selector]; ok {
		for k, v := range s.Properties {
			(*style)[k] = v
		}
	}
}

// applyInlineStyles applies inline CSS styles to the style map
func (b *Builder) applyInlineStyles(style *map[string]string, inlineStyle string) {
	if inlineStyle == "" {
		return
	}

	inlineStyles := parseCSS(inlineStyle)
	for k, v := range inlineStyles {
		(*style)[k] = v
	}
}

// computeStyle calculates the final style for an element
// This applies styles in order of precedence: class < id < inline
func (b *Builder) computeStyle(elem Element) map[string]string {
	style := make(map[string]string)

	// Apply styles from classes (lowest precedence)
	b.applyClassStyles(&style, elem.Class)

	// Apply styles from ID (medium precedence)
	b.applyIDStyles(&style, elem.ID)

	// Apply inline styles (highest precedence)
	b.applyInlineStyles(&style, elem.Style)

	return style
}
