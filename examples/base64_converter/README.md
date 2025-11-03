# Base64 Converter Example

This example demonstrates a Base64 encoder/decoder application built with Fylay.

## Features

- **Text Conversion**: Convert text to/from Base64 encoding in the UI
- **Direct File Encoding**: Encode files directly to Base64 without loading into memory (handles large files)
- **File Save**: Save conversion results to files
- **Statistics**: Real-time character and byte count for input and output
- **User-friendly Interface**: Clear separation between input and output areas with side-by-side encode/decode buttons
- **Status Feedback**: Status bar showing operation results and errors

## Usage

### Running the Application

From the example directory:

```bash
cd examples/base64_converter
go run base64_converter.go
```

Or from the project root:

```bash
go run examples/base64_converter/base64_converter.go
```

Note: The application must be run from the `examples/base64_converter` directory or the XML file path must be adjusted accordingly.

### Operations

1. **Encode to Base64** (Text):
   - Enter text in the input area
   - Click "Encode to Base64 ▼" button
   - The Base64 encoded result appears in the output area

2. **Decode from Base64** (Text):
   - Enter Base64 encoded text in the input area
   - Click "Decode from Base64 ▲" button
   - The decoded result appears in the output area

3. **Encode File to Base64** (Direct conversion):
   - Click "Encode File to Base64" button
   - Select a file from the file dialog
   - File is immediately encoded without loading into the UI (prevents freezing with large files)
   - Choose where to save the Base64 encoded result
   - The encoded file is saved directly

4. **Save to File**:
   - After encoding/decoding text, click "Save to File" button
   - Choose the destination in the file dialog
   - The output content will be saved to the selected file

## XML Layout Features

The example showcases:

- Multi-line text entry widgets
- Custom CSS-like styling for layout and appearance
- Event handlers using `onclick` attributes in XML
- Callback registration with `builder.On()` method
- Read-only output area
- Dynamic content updates
- Flexible box layout (VBox and HBox)
- Style selectors for reusable styling

## Code Structure

- `base64_converter.xml`: UI layout definition with styling
- `base64_converter.go`: Application logic and event handlers
- `README.md`: This documentation file

## Learning Points

This example demonstrates:

1. How to create complex multi-widget layouts
2. Handling file I/O operations with Fyne dialogs
3. Implementing bidirectional conversion logic
4. Providing user feedback through status updates
5. Managing widget state and content updates
6. Using `onclick` callbacks registered with `builder.On()`
7. Applying CSS-like styles with `<Style>` selectors
