package service

import (
	"image"
	"image/color"
	"testing"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

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

func TestNewImageCodec(t *testing.T) {
	codec := NewImageCodec()
	if codec == nil {
		t.Fatal("NewImageCodec() returned nil")
	}
}

func TestImageCodec_EncodePNG(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	data, err := codec.EncodePNG(img)
	if err != nil {
		t.Fatalf("EncodePNG() error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("EncodePNG() returned empty data")
	}
	if len(data) < 8 || data[0] != 0x89 || data[1] != 0x50 {
		t.Error("EncodePNG() did not produce valid PNG magic bytes")
	}
}

func TestImageCodec_EncodeJPEG(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	data, err := codec.EncodeJPEG(img, 90)
	if err != nil {
		t.Fatalf("EncodeJPEG() error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("EncodeJPEG() returned empty data")
	}
	if len(data) < 3 || data[0] != 0xFF || data[1] != 0xD8 {
		t.Error("EncodeJPEG() did not produce valid JPEG magic bytes")
	}
}

func TestImageCodec_EncodeGIF(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	data, err := codec.EncodeGIF(img)
	if err != nil {
		t.Fatalf("EncodeGIF() error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("EncodeGIF() returned empty data")
	}
}

func TestImageCodec_DecodePNG(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	data, err := codec.EncodePNG(img)
	if err != nil {
		t.Fatalf("Setup error = %v", err)
	}
	decoded, err := codec.DecodePNG(data)
	if err != nil {
		t.Fatalf("DecodePNG() error = %v", err)
	}
	if decoded == nil {
		t.Fatal("DecodePNG() returned nil image")
	}
	bounds := decoded.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("DecodePNG() dimensions = %dx%d, want 10x10", bounds.Dx(), bounds.Dy())
	}
}

func TestImageCodec_DecodeJPEG(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	data, err := codec.EncodeJPEG(img, 90)
	if err != nil {
		t.Fatalf("Setup error = %v", err)
	}
	decoded, err := codec.DecodeJPEG(data)
	if err != nil {
		t.Fatalf("DecodeJPEG() error = %v", err)
	}
	if decoded == nil {
		t.Fatal("DecodeJPEG() returned nil image")
	}
}

func TestImageCodec_DecodeGIF(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	data, err := codec.EncodeGIF(img)
	if err != nil {
		t.Fatalf("Setup error = %v", err)
	}
	decoded, err := codec.DecodeGIF(data)
	if err != nil {
		t.Fatalf("DecodeGIF() error = %v", err)
	}
	if decoded == nil {
		t.Fatal("DecodeGIF() returned nil image")
	}
}

func TestImageCodec_Decode(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	pngData, _ := codec.EncodePNG(img)
	decoded, format, err := codec.Decode(pngData)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if decoded == nil {
		t.Fatal("Decode() returned nil image")
	}
	if format != value.MIMETypeImagePNG {
		t.Errorf("Decode() format = %v, want %v", format, value.MIMETypeImagePNG)
	}
}

func TestImageCodec_DetectFormat(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	pngData, _ := codec.EncodePNG(img)
	format, err := codec.DetectFormat(pngData)
	if err != nil {
		t.Fatalf("DetectFormat() error = %v", err)
	}
	if format != value.MIMETypeImagePNG {
		t.Errorf("DetectFormat() = %v, want %v", format, value.MIMETypeImagePNG)
	}
}

func TestImageCodec_ConvertFormat(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()
	pngData, _ := codec.EncodePNG(img)
	jpegData, err := codec.ConvertFormat(pngData, value.MIMETypeImageJPEG)
	if err != nil {
		t.Fatalf("ConvertFormat() error = %v", err)
	}
	if len(jpegData) == 0 {
		t.Fatal("ConvertFormat() returned empty data")
	}
	format, _ := codec.DetectFormat(jpegData)
	if format != value.MIMETypeImageJPEG {
		t.Errorf("ConvertFormat() produced format %v, want %v", format, value.MIMETypeImageJPEG)
	}
}

func TestImageCodec_RoundTrip(t *testing.T) {
	codec := NewImageCodec()
	original := createTestImage()

	encoded, _ := codec.EncodePNG(original)
	decoded, format, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode error = %v", err)
	}
	if format != value.MIMETypeImagePNG {
		t.Errorf("Format mismatch: got %v, want %v", format, value.MIMETypeImagePNG)
	}
	origBounds := original.Bounds()
	decodedBounds := decoded.Bounds()
	if origBounds.Dx() != decodedBounds.Dx() || origBounds.Dy() != decodedBounds.Dy() {
		t.Errorf("Dimensions mismatch")
	}
}

func TestImageCodec_ErrorCases(t *testing.T) {
	codec := NewImageCodec()

	_, err := codec.EncodePNG(nil)
	if err == nil {
		t.Error("EncodePNG should fail with nil image")
	}

	_, err = codec.EncodeJPEG(nil, 90)
	if err == nil {
		t.Error("EncodeJPEG should fail with nil image")
	}

	_, err = codec.EncodeGIF(nil)
	if err == nil {
		t.Error("EncodeGIF should fail with nil image")
	}

	_, err = codec.DecodePNG([]byte{})
	if err == nil {
		t.Error("DecodePNG should fail with empty data")
	}

	_, err = codec.DecodeJPEG([]byte{})
	if err == nil {
		t.Error("DecodeJPEG should fail with empty data")
	}

	_, err = codec.DecodeGIF([]byte{})
	if err == nil {
		t.Error("DecodeGIF should fail with empty data")
	}

	_, _, err = codec.Decode([]byte{})
	if err == nil {
		t.Error("Decode should fail with empty data")
	}

	_, err = codec.DetectFormat([]byte{})
	if err == nil {
		t.Error("DetectFormat should fail with empty data")
	}

	_, err = codec.ConvertFormat([]byte{0x00}, value.MIMEType("image/webp"))
	if err == nil {
		t.Error("ConvertFormat should fail with unsupported format")
	}
}
