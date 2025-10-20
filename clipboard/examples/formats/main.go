// Package main demonstrates clipboard content format detection and handling.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/clipboard/api"
	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
	"github.com/phoenix-tui/phoenix/clipboard/domain/value"
)

//nolint:gocyclo,cyclop // Example code demonstrates multiple format handling scenarios
func main() {
	fmt.Println("Phoenix Clipboard - Multiple Formats Example")
	fmt.Println("=============================================")
	fmt.Println()

	// Create clipboard instance
	clipboard, err := api.New()
	if err != nil {
		log.Fatalf("Failed to create clipboard: %v", err)
	}

	if !clipboard.IsAvailable() {
		log.Fatal("Clipboard is not available")
	}

	fmt.Printf("Using provider: %s\n", clipboard.GetProviderName())
	fmt.Println()

	// Example 1: Plain text
	fmt.Println("1. Plain Text Format")
	fmt.Println("--------------------")

	plainText := "Hello, World!"
	err = clipboard.Write(plainText)
	if err != nil {
		log.Fatalf("Failed to write plain text: %v", err)
	}
	fmt.Printf("✓ Written: %s\n", plainText)

	readText, err := clipboard.Read()
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}
	fmt.Printf("✓ Read: %s\n", readText)
	fmt.Println()

	// Example 2: Multi-line text
	fmt.Println("2. Multi-line Text")
	fmt.Println("------------------")

	multilineText := `Line 1
Line 2
Line 3`
	err = clipboard.Write(multilineText)
	if err != nil {
		log.Fatalf("Failed to write multiline text: %v", err)
	}
	fmt.Printf("✓ Written %d lines\n", 3)

	readMultiline, err := clipboard.Read()
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}
	fmt.Printf("✓ Read:\n%s\n", readMultiline)
	fmt.Println()

	// Example 3: Unicode and emoji
	fmt.Println("3. Unicode and Emoji")
	fmt.Println("--------------------")

	unicodeText := "Hello 世界! 🚀 Phoenix TUI 🎉"
	err = clipboard.Write(unicodeText)
	if err != nil {
		log.Fatalf("Failed to write unicode text: %v", err)
	}
	fmt.Printf("✓ Written: %s\n", unicodeText)

	readUnicode, err := clipboard.Read()
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}
	fmt.Printf("✓ Read: %s\n", readUnicode)
	fmt.Println()

	// Example 4: Domain model usage
	fmt.Println("4. Domain Model Usage")
	fmt.Println("---------------------")

	content, err := model.NewTextContent("Domain model content")
	if err != nil {
		log.Fatalf("Failed to create content: %v", err)
	}

	fmt.Printf("Content type: %s\n", content.MIMEType())
	fmt.Printf("Content encoding: %s\n", content.Encoding())
	fmt.Printf("Content size: %d bytes\n", content.Size())
	fmt.Printf("Is text: %v\n", content.IsText())
	fmt.Printf("Is binary: %v\n", content.IsBinary())
	fmt.Println()

	// Example 5: Binary content (conceptual)
	fmt.Println("5. Binary Content (Conceptual)")
	fmt.Println("------------------------------")

	binaryData := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F} // "Hello" in bytes
	binaryContent, err := model.NewBinaryContent(binaryData)
	if err != nil {
		log.Fatalf("Failed to create binary content: %v", err)
	}

	fmt.Printf("Binary content type: %s\n", binaryContent.MIMEType())
	fmt.Printf("Binary content encoding: %s\n", binaryContent.Encoding())
	fmt.Printf("Binary content size: %d bytes\n", binaryContent.Size())
	fmt.Printf("Is text: %v\n", binaryContent.IsText())
	fmt.Printf("Is binary: %v\n", binaryContent.IsBinary())
	fmt.Println()
	fmt.Println("Note: Binary clipboard support depends on the platform.")
	fmt.Println("Currently, text content is most widely supported.")
	fmt.Println()

	// Example 6: Content transformation
	fmt.Println("6. Content Transformation")
	fmt.Println("-------------------------")

	originalContent, _ := model.NewTextContent("Original")
	fmt.Printf("Original MIME type: %s\n", originalContent.MIMEType())

	// Transform to HTML
	htmlContent := originalContent.WithMIMEType(value.MIMETypeHTML)
	fmt.Printf("Transformed MIME type: %s\n", htmlContent.MIMEType())

	// Transform encoding
	base64Content := htmlContent.WithEncoding(value.EncodingBase64)
	fmt.Printf("Transformed encoding: %s\n", base64Content.Encoding())

	// Original is unchanged (immutability)
	fmt.Printf("Original still: %s\n", originalContent.MIMEType())
	fmt.Println()

	// Example 7: Large content
	fmt.Println("7. Large Content")
	fmt.Println("----------------")

	largeText := ""
	for i := 0; i < 1000; i++ {
		largeText += fmt.Sprintf("Line %d\n", i+1)
	}

	err = clipboard.Write(largeText)
	if err != nil {
		log.Fatalf("Failed to write large content: %v", err)
	}
	fmt.Printf("✓ Written large content: %d bytes\n", len(largeText))

	readLarge, err := clipboard.Read()
	if err != nil {
		log.Fatalf("Failed to read large content: %v", err)
	}
	fmt.Printf("✓ Read large content: %d bytes\n", len(readLarge))

	if len(readLarge) == len(largeText) {
		fmt.Println("✓ Large content preserved correctly")
	}

	fmt.Println()
	fmt.Println("All format examples completed successfully!")
}
