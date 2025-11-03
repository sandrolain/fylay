# Dialog Demo Example

This example demonstrates the integration of native system dialogs using the `flydialog` package.

## Features Demonstrated

- **Info Dialog**: Show informational messages
- **Error Dialog**: Display error messages  
- **Question Dialog**: Ask Yes/No questions
- **File Open Dialog**: Select files to open with filters
- **File Save Dialog**: Choose location to save files
- **Directory Select Dialog**: Pick a directory

## Running the Example

```bash
cd examples
go run dialogs_demo.go
```

## Code Structure

The example shows how to:

1. **Register Dialog Callbacks**: Use `builder.On()` to register callbacks that show dialogs
2. **Handle Dialog Results**: Process file paths, user choices, and errors
3. **Use File Filters**: Filter files by extension in open/save dialogs
4. **Read/Write Files**: Load and save file content using selected paths
5. **Show Feedback**: Display success/error messages to users

## Dialog Types Used

### Message Dialogs

```go
builder.On("onShowInfo", func(ctx *flay.EventContext) {
    builder.ShowInfoDialog("Information", "This is an info message!")
})

builder.On("onShowError", func(ctx *flay.EventContext) {
    builder.ShowErrorDialog("Error", "This is an error message!")
})

builder.On("onShowQuestion", func(ctx *flay.EventContext) {
    if builder.ShowQuestionDialog("Question", "Do you want to continue?") {
        // User clicked Yes
    } else {
        // User clicked No
    }
})
```

### File Dialogs

```go
builder.On("onOpenFile", func(ctx *flay.EventContext) {
    filepath, err := builder.ShowFileOpenDialog(
        "Select a file",
        flydialog.FileFilter{Description: "Text files", Pattern: "*.txt"},
        flydialog.FileFilter{Description: "Go files", Pattern: "*.go"},
    )
    if err != nil {
        // User cancelled
        return
    }
    
    // Read file content
    content, _ := os.ReadFile(filepath)
    // ... use content
})
```

## XML Layout

The layout (`dialogs_demo.xml`) defines:

- Buttons for each dialog type
- Result label to show dialog outcomes
- Text area for file content display

## Notes

- Dialogs are **modal** and block execution
- Use **native** system dialogs (not Fyne dialogs)
- Works on macOS, Windows, and Linux
- Requires the `github.com/sqweek/dialog` library

## See Also

- `/docs/DIALOGS.md` - Complete dialog documentation
- `base64_converter/` - Real-world example using file dialogs
