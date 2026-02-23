# Image Clipboard Example

This example demonstrates the image clipboard functionality in Phoenix clipboard library.

## Features Demonstrated

1. **Image Creation** - Creating test images programmatically
2. **Format Detection** - Detecting image format from bytes (magic bytes)
3. **Image Encoding** - Encoding images to PNG, JPEG, GIF formats
4. **Image Decoding** - Decoding images from various formats
5. **Format Conversion** - Converting between PNG, JPEG, GIF
6. **MIME Type Support** - Working with image MIME types
7. **Round-trip Testing** - Encoding and decoding to verify data integrity

## Supported Image Formats

- PNG (image/png)
- JPEG (image/jpeg)
- GIF (image/gif)
- BMP (image/bmp) - detection only

## Building and Running

```bash
cd examples/image-clipboard
go build
./image-clipboard    # On Linux/macOS
image-clipboard.exe  # On Windows
```

## Example Output

```
Phoenix Clipboard - Image Support Demo
=======================================

1. Creating Test Image (10x10 red square)...
   Created image with dimensions: 10x10

2. Encoding to PNG...
   PNG size: 119 bytes

3. Detecting format from PNG data...
   Detected format: image/png

4. Converting PNG to JPEG...
   JPEG size: 632 bytes

5. Converting PNG to GIF...
   GIF size: 126 bytes

6. Decoding PNG back to image...
   Decoded image dimensions: 10x10

7. Round-trip test (Encode -> Decode -> Encode)...
   Round-trip successful! Size: 119 bytes

8. Testing all supported MIME types...
   - image/png: IsImage=true, IsBinary=true
   - image/jpeg: IsImage=true, IsBinary=true
   - image/gif: IsImage=true, IsBinary=true
   - image/bmp: IsImage=true, IsBinary=true

9. Saving test images to files...
   ✓ Saved: test_output.png
   ✓ Saved: test_output.jpg
   ✓ Saved: test_output.gif

Demo completed successfully!
```

## Generated Files

Running the demo creates three image files in the current directory:

- `test_output.png` - PNG format (smallest for simple graphics)
- `test_output.jpg` - JPEG format (lossy compression)
- `test_output.gif` - GIF format (supports animation, limited colors)

## Architecture

This demo uses the **ImageCodec** domain service from the Phoenix clipboard library:

```
clipboard/
├── internal/
│   ├── domain/
│   │   ├── service/
│   │   │   └── image_codec.go      # Image encoding/decoding service
│   │   └── value/
│   │       └── mime_type.go        # MIME type value objects
│   └── application/
│       └── clipboard_manager.go    # Application service
└── clipboard.go                    # Public API
```

## Implementation Details

### ImageCodec Service

The `ImageCodec` is a domain service providing:

- **EncodePNG(img)** - Encode image to PNG format
- **EncodeJPEG(img, quality)** - Encode image to JPEG with quality control
- **EncodeGIF(img)** - Encode image to GIF format
- **DecodePNG(data)** - Decode PNG image data
- **DecodeJPEG(data)** - Decode JPEG image data
- **DecodeGIF(data)** - Decode GIF image data
- **Decode(data)** - Auto-detect format and decode
- **DetectFormat(data)** - Detect image format from magic bytes
- **ConvertFormat(data, targetFormat)** - Convert between formats

### Test Coverage

The ImageCodec service has extensive test coverage, exceeding project targets.

## Clipboard Integration (Future)

Full clipboard copy/paste operations require platform-specific native provider implementation:

```go
// Future API (not yet implemented in providers)
clipboard := clipboard.New()

// Copy image to clipboard
imageData, _ := os.ReadFile("photo.png")
clipboard.WriteImagePNG(imageData)

// Paste image from clipboard
imageData, mimeType, _ := clipboard.ReadImage()
os.WriteFile("pasted.png", imageData, 0644)
```

## Notes

- This demo uses the Go standard library `image`, `image/png`, `image/jpeg`, `image/gif` packages
- No external image processing libraries are required
- Image format detection uses magic bytes for fast identification
- Quality parameter for JPEG encoding defaults to 90 if out of range (1-100)

## Related Examples

- [basic](../basic/) - Basic clipboard text operations
- [formats](../formats/) - Multiple text formats (plain text, HTML, RTF)
- [osc52](../osc52/) - Remote clipboard via OSC 52 escape sequences
