package flay

import (
	"bytes"
	"testing"
)

// BenchmarkComputeStyle measures performance of style computation
func BenchmarkComputeStyle(b *testing.B) {
	xml := `
		<Layout>
			<Style selector=".text">color: red; font-size: 14px;</Style>
			<Style selector=".large">font-size: 20px; font-weight: bold;</Style>
			<Style selector="#unique">color: blue; background: white;</Style>
			<VBox class="text large" id="unique" style="padding: 10px;">
				<Label>Test</Label>
			</VBox>
		</Layout>
	`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(bytes.NewReader([]byte(xml)))
	if err != nil {
		b.Fatalf("Failed to load layout: %v", err)
	}

	elem := layout.Root

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.computeStyle(elem)
	}
}

// BenchmarkComputeStyleClassOnly measures style computation with only classes
func BenchmarkComputeStyleClassOnly(b *testing.B) {
	xml := `
		<Layout>
			<Style selector=".text">color: red;</Style>
			<VBox class="text">
				<Label>Test</Label>
			</VBox>
		</Layout>
	`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(bytes.NewReader([]byte(xml)))
	if err != nil {
		b.Fatalf("Failed to load layout: %v", err)
	}

	elem := layout.Root

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.computeStyle(elem)
	}
}

// BenchmarkComputeStyleInlineOnly measures style computation with only inline styles
func BenchmarkComputeStyleInlineOnly(b *testing.B) {
	xml := `
		<Layout>
			<VBox style="color: red; font-size: 14px; padding: 10px;">
				<Label>Test</Label>
			</VBox>
		</Layout>
	`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(bytes.NewReader([]byte(xml)))
	if err != nil {
		b.Fatalf("Failed to load layout: %v", err)
	}

	elem := layout.Root

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.computeStyle(elem)
	}
}

// BenchmarkParseColor measures performance of color parsing
func BenchmarkParseColor(b *testing.B) {
	colors := []string{
		"red",
		"#FF0000",
		"#F00",
		"rgb(255, 0, 0)",
		"blue",
		"#00FF00",
		"rgb(0, 255, 0)",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, color := range colors {
			_ = parseColor(color)
		}
	}
}

// BenchmarkParseColorNamed measures performance of named color parsing
func BenchmarkParseColorNamed(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseColor("red")
	}
}

// BenchmarkParseColorHex6 measures performance of 6-digit hex parsing
func BenchmarkParseColorHex6(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseColor("#FF0000")
	}
}

// BenchmarkParseColorHex3 measures performance of 3-digit hex parsing
func BenchmarkParseColorHex3(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseColor("#F00")
	}
}

// BenchmarkParseColorRGB measures performance of RGB parsing
func BenchmarkParseColorRGB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseColor("rgb(255, 0, 0)")
	}
}

// BenchmarkParseCSS measures performance of CSS parsing
func BenchmarkParseCSS(b *testing.B) {
	css := "color: red; font-size: 14px; padding: 10px; margin: 5px; background: white; border: 1px solid black;"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseCSS(css)
	}
}

// BenchmarkBuildLayout measures overall layout building performance
func BenchmarkBuildLayout(b *testing.B) {
	xml := `
		<Layout>
			<Style selector=".text">color: red; font-size: 14px;</Style>
			<Style selector="#header">font-size: 20px; font-weight: bold;</Style>
			<VBox>
				<Label id="header">Title</Label>
				<Label class="text">Content 1</Label>
				<Label class="text">Content 2</Label>
				<HBox>
					<Button>OK</Button>
					<Button>Cancel</Button>
				</HBox>
			</VBox>
		</Layout>
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewBuilder()
		layout, err := builder.LoadLayout(bytes.NewReader([]byte(xml)))
		if err != nil {
			b.Fatalf("Failed to load layout: %v", err)
		}
		_, _ = builder.Build(layout)
	}
}
