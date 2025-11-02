package dialog

import (
	"testing"
)

func TestDialogConfig(t *testing.T) {
	tests := []struct {
		name   string
		config DialogConfig
	}{
		{
			name: "Info dialog config",
			config: DialogConfig{
				Type:    DialogTypeInfo,
				Title:   "Test Info",
				Message: "This is a test",
			},
		},
		{
			name: "Error dialog config",
			config: DialogConfig{
				Type:    DialogTypeError,
				Title:   "Test Error",
				Message: "This is an error",
			},
		},
		{
			name: "Question dialog config",
			config: DialogConfig{
				Type:    DialogTypeQuestion,
				Title:   "Test Question",
				Message: "Are you sure?",
			},
		},
		{
			name: "FileOpen dialog config",
			config: DialogConfig{
				Type:        DialogTypeFileOpen,
				Title:       "Open File",
				DefaultPath: "/tmp",
				Filters: []FileFilter{
					{Description: "Text files", Pattern: "*.txt"},
				},
			},
		},
		{
			name: "FileSave dialog config",
			config: DialogConfig{
				Type:        DialogTypeFileSave,
				Title:       "Save File",
				DefaultPath: "/tmp",
				Filters: []FileFilter{
					{Description: "JSON files", Pattern: "*.json"},
				},
			},
		},
		{
			name: "DirSelect dialog config",
			config: DialogConfig{
				Type:        DialogTypeDirSelect,
				Title:       "Select Directory",
				DefaultPath: "/tmp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Type == "" {
				t.Error("DialogType should not be empty")
			}
			if tt.config.Title == "" {
				t.Error("Title should not be empty")
			}
		})
	}
}

func TestFileFilter(t *testing.T) {
	filter := FileFilter{
		Description: "Text files",
		Pattern:     "*.txt",
	}

	if filter.Description == "" {
		t.Error("Description should not be empty")
	}
	if filter.Pattern == "" {
		t.Error("Pattern should not be empty")
	}
}

func TestDialogTypes(t *testing.T) {
	types := []DialogType{
		DialogTypeInfo,
		DialogTypeError,
		DialogTypeQuestion,
		DialogTypeFileOpen,
		DialogTypeFileSave,
		DialogTypeDirSelect,
	}

	for _, dt := range types {
		if dt == "" {
			t.Error("DialogType should not be empty")
		}
	}
}
