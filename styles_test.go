package fylay

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// TestApplyMinSizeToWidgets tests that width/height styles are applied to widgets
func TestApplyMinSizeToWidgets(t *testing.T) {
	// Create test app for widget testing
	_ = test.NewApp()

	tests := []struct {
		name     string
		xml      string
		elemID   string
		wantSize fyne.Size
	}{
		{
			name: "Button with width and height",
			xml: `<Layout>
				<Button id="btn1" style="width: 200px; height: 50px">Click Me</Button>
			</Layout>`,
			elemID:   "btn1",
			wantSize: fyne.NewSize(200, 50),
		},
		{
			name: "Entry with CSS class",
			xml: `<Layout>
				<Style selector=".wide-entry">
					width: 300px;
					height: 100px;
				</Style>
				<Entry id="entry1" class="wide-entry" />
			</Layout>`,
			elemID:   "entry1",
			wantSize: fyne.NewSize(300, 100),
		},
		{
			name: "Select with ID style",
			xml: `<Layout>
				<Style selector="#mySelect">
					width: 250px;
				</Style>
				<Select id="mySelect">
					<Option>Option 1</Option>
				</Select>
			</Layout>`,
			elemID: "mySelect",
			// Only width is set, height will be from widget's default
		},
		{
			name: "Slider with inline style",
			xml: `<Layout>
				<Slider id="slider1" style="width: 400px" min="0" max="100" value="50" />
			</Layout>`,
			elemID: "slider1",
		},
		{
			name: "Button with min-width and min-height",
			xml: `<Layout>
				<Button id="btn2" style="min-width: 250px; min-height: 60px">Min Size Button</Button>
			</Layout>`,
			elemID:   "btn2",
			wantSize: fyne.NewSize(250, 60),
		},
		{
			name: "Entry with min-height from CSS class",
			xml: `<Layout>
				<Style selector=".tall-entry">
					min-height: 150px;
				</Style>
				<Entry id="entry2" class="tall-entry" />
			</Layout>`,
			elemID:   "entry2",
			wantSize: fyne.NewSize(0, 150), // Only height is set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			layout, err := builder.LoadLayout(strings.NewReader(tt.xml))
			if err != nil {
				t.Fatalf("Failed to load layout: %v", err)
			}

			_, err = builder.Build(layout)
			if err != nil {
				t.Fatalf("Failed to build layout: %v", err)
			}

			elem := builder.GetElement(tt.elemID)
			if elem == nil {
				t.Fatalf("Element %s not found", tt.elemID)
			}

			minSize := elem.MinSize()

			// Check width if specified
			if tt.wantSize.Width > 0 && minSize.Width < tt.wantSize.Width {
				t.Errorf("Width = %v, want >= %v", minSize.Width, tt.wantSize.Width)
			}

			// Check height if specified
			if tt.wantSize.Height > 0 && minSize.Height < tt.wantSize.Height {
				t.Errorf("Height = %v, want >= %v", minSize.Height, tt.wantSize.Height)
			}
		})
	}
}

// TestApplyMinSizeToCanvasObjects tests that width/height styles work for canvas objects
func TestApplyMinSizeToCanvasObjects(t *testing.T) {
	// Create test app
	_ = test.NewApp()

	tests := []struct {
		name     string
		xml      string
		elemID   string
		wantSize fyne.Size
	}{
		{
			name: "Rectangle with width and height",
			xml: `<Layout>
				<Rectangle id="rect1" style="width: 100px; height: 50px; background-color: red;" />
			</Layout>`,
			elemID:   "rect1",
			wantSize: fyne.NewSize(100, 50),
		},
		{
			name: "Circle with size",
			xml: `<Layout>
				<Circle id="circle1" style="width: 75px; background-color: blue;" />
			</Layout>`,
			elemID:   "circle1",
			wantSize: fyne.NewSize(75, 75),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			layout, err := builder.LoadLayout(strings.NewReader(tt.xml))
			if err != nil {
				t.Fatalf("Failed to load layout: %v", err)
			}

			_, err = builder.Build(layout)
			if err != nil {
				t.Fatalf("Failed to build layout: %v", err)
			}

			elem := builder.GetElement(tt.elemID)
			if elem == nil {
				t.Fatalf("Element %s not found", tt.elemID)
			}

			minSize := elem.MinSize()

			// Special case for Circle which uses Resize instead of MinSize
			if tt.elemID == "circle1" {
				size := elem.Size()
				if size.Width < tt.wantSize.Width {
					t.Errorf("Width = %v, want >= %v", size.Width, tt.wantSize.Width)
				}
				if size.Height < tt.wantSize.Height {
					t.Errorf("Height = %v, want >= %v", size.Height, tt.wantSize.Height)
				}
				return
			}

			// Canvas objects should have exact sizes
			if minSize.Width < tt.wantSize.Width {
				t.Errorf("Width = %v, want >= %v", minSize.Width, tt.wantSize.Width)
			}

			if minSize.Height < tt.wantSize.Height {
				t.Errorf("Height = %v, want >= %v", minSize.Height, tt.wantSize.Height)
			}
		})
	}
}

// TestStylePrecedenceWithSizes tests that inline styles override class and ID styles for sizes
func TestStylePrecedenceWithSizes(t *testing.T) {
	// Create test app
	_ = test.NewApp()

	xml := `<Layout>
		<Style selector=".btn-class">
			width: 100px;
			height: 40px;
		</Style>
		<Style selector="#btn1">
			width: 150px;
		</Style>
		<Button id="btn1" class="btn-class" style="width: 200px">Test</Button>
	</Layout>`

	builder := NewBuilder()
	layout, err := builder.LoadLayout(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Failed to load layout: %v", err)
	}

	_, err = builder.Build(layout)
	if err != nil {
		t.Fatalf("Failed to build layout: %v", err)
	}

	btn := builder.GetElement("btn1")
	if btn == nil {
		t.Fatal("Button not found")
	}

	minSize := btn.MinSize()

	// Inline style (200px) should override ID style (150px) and class style (100px)
	if minSize.Width < 200 {
		t.Errorf("Width = %v, want >= 200 (inline style should take precedence)", minSize.Width)
	}

	// Height should come from class style (40px) since it's not overridden
	if minSize.Height < 40 {
		t.Errorf("Height = %v, want >= 40 (from class style)", minSize.Height)
	}
}
