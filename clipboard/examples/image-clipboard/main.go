// Package main demonstrates image clipboard operations using Phoenix clipboard library.
package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

func main() {
	fmt.Println("Phoenix Clipboard - Image Support Demo")
	fmt.Println("=======================================")
	fmt.Println()

	// Demonstrate ImageCodec functionality
	demonstrateImageCodec()
}

func demonstrateImageCodec() {
	codec := service.NewImageCodec()

	fmt.Println("1. Creating Test Image (10x10 red square)...")
	img := createTestImage()
	fmt.Printf("   Created image with dimensions: %dx%d\n",
		img.Bounds().Dx(), img.Bounds().Dy())
	fmt.Println()

	fmt.Println("2. Encoding to PNG...")
	pngData, err := codec.EncodePNG(img)
	if err != nil {
		fmt.Printf("   Error encoding PNG: %v\n", err)
		return
	}
	fmt.Printf("   PNG size: %d bytes\n", len(pngData))
	fmt.Println()

	fmt.Println("3. Detecting format from PNG data...")
	format, err := codec.DetectFormat(pngData)
	if err != nil {
		fmt.Printf("   Error detecting format: %v\n", err)
		return
	}
	fmt.Printf("   Detected format: %s\n", format)
	fmt.Println()

	fmt.Println("4. Converting PNG to JPEG...")
	jpegData, err := codec.ConvertFormat(pngData, value.MIMETypeImageJPEG)
	if err != nil {
		fmt.Printf("   Error converting to JPEG: %v\n", err)
		return
	}
	fmt.Printf("   JPEG size: %d bytes\n", len(jpegData))
	fmt.Println()

	fmt.Println("5. Converting PNG to GIF...")
	gifData, err := codec.ConvertFormat(pngData, value.MIMETypeImageGIF)
	if err != nil {
		fmt.Printf("   Error converting to GIF: %v\n", err)
		return
	}
	fmt.Printf("   GIF size: %d bytes\n", len(gifData))
	fmt.Println()

	fmt.Println("6. Decoding PNG back to image...")
	decoded, err := codec.DecodePNG(pngData)
	if err != nil {
		fmt.Printf("   Error decoding PNG: %v\n", err)
		return
	}
	fmt.Printf("   Decoded image dimensions: %dx%d\n",
		decoded.Bounds().Dx(), decoded.Bounds().Dy())
	fmt.Println()

	fmt.Println("7. Round-trip test (Encode -> Decode -> Encode)...")
	roundtrip, err := codec.EncodePNG(decoded)
	if err != nil {
		fmt.Printf("   Error in round-trip: %v\n", err)
		return
	}
	fmt.Printf("   Round-trip successful! Size: %d bytes\n", len(roundtrip))
	fmt.Println()

	fmt.Println("8. Testing all supported MIME types...")
	supportedTypes := []value.MIMEType{
		value.MIMETypeImagePNG,
		value.MIMETypeImageJPEG,
		value.MIMETypeImageGIF,
		value.MIMETypeImageBMP,
	}
	for _, mimeType := range supportedTypes {
		fmt.Printf("   - %s: IsImage=%v, IsBinary=%v\n",
			mimeType, mimeType.IsImage(), mimeType.IsBinary())
	}
	fmt.Println()

	fmt.Println("9. Saving test images to files...")
	if err := os.WriteFile("test_output.png", pngData, 0644); err != nil {
		fmt.Printf("   Error saving PNG: %v\n", err)
	} else {
		fmt.Println("   ✓ Saved: test_output.png")
	}

	if err := os.WriteFile("test_output.jpg", jpegData, 0644); err != nil {
		fmt.Printf("   Error saving JPEG: %v\n", err)
	} else {
		fmt.Println("   ✓ Saved: test_output.jpg")
	}

	if err := os.WriteFile("test_output.gif", gifData, 0644); err != nil {
		fmt.Printf("   Error saving GIF: %v\n", err)
	} else {
		fmt.Println("   ✓ Saved: test_output.gif")
	}
	fmt.Println()

	fmt.Println("Demo completed successfully!")
	fmt.Println()
	fmt.Println("Note: Full clipboard operations (copy/paste images) require")
	fmt.Println("      platform-specific native provider implementation.")
	fmt.Println("      This demo shows the ImageCodec functionality that powers")
	fmt.Println("      image format conversion and detection.")
}

func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, red)
		}
	}
	return img
}
