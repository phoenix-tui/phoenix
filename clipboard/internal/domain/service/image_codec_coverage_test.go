package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

func TestImageCodec_EncodeJPEG_QualityVariations(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()

	qualities := []int{1, 50, 90, 100, 0, -10, 200}
	for _, q := range qualities {
		data, err := codec.EncodeJPEG(img, q)
		if err != nil {
			t.Errorf("EncodeJPEG(quality=%d) error = %v", q, err)
		}
		if len(data) == 0 {
			t.Errorf("EncodeJPEG(quality=%d) returned empty data", q)
		}
	}
}

func TestImageCodec_DetectFormat_AllMagicBytes(t *testing.T) {
	codec := NewImageCodec()

	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	format, err := codec.DetectFormat(pngMagic)
	if err != nil || format != value.MIMETypeImagePNG {
		t.Errorf("DetectFormat(PNG) = %v, %v", format, err)
	}

	jpegMagic := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	format, err = codec.DetectFormat(jpegMagic)
	if err != nil || format != value.MIMETypeImageJPEG {
		t.Errorf("DetectFormat(JPEG) = %v, %v", format, err)
	}

	gifMagic := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}
	format, err = codec.DetectFormat(gifMagic)
	if err != nil || format != value.MIMETypeImageGIF {
		t.Errorf("DetectFormat(GIF) = %v, %v", format, err)
	}

	bmpMagic := []byte{0x42, 0x4D, 0x00, 0x00, 0x00, 0x00}
	format, err = codec.DetectFormat(bmpMagic)
	if err != nil || format != value.MIMETypeImageBMP {
		t.Errorf("DetectFormat(BMP) = %v, %v", format, err)
	}
}

func TestImageCodec_ConvertFormat_AllCombinations(t *testing.T) {
	codec := NewImageCodec()
	img := createTestImage()

	pngData, _ := codec.EncodePNG(img)
	jpegData, _ := codec.EncodeJPEG(img, 90)
	gifData, _ := codec.EncodeGIF(img)

	sources := []struct {
		name string
		data []byte
	}{
		{"PNG", pngData},
		{"JPEG", jpegData},
		{"GIF", gifData},
	}

	targets := []value.MIMEType{
		value.MIMETypeImagePNG,
		value.MIMETypeImageJPEG,
		value.MIMETypeImageGIF,
	}

	for _, src := range sources {
		for _, tgt := range targets {
			_, err := codec.ConvertFormat(src.data, tgt)
			if err != nil {
				t.Errorf("ConvertFormat(%s to %s) error = %v", src.name, tgt, err)
			}
		}
	}

	_, err := codec.ConvertFormat(pngData, value.MIMEType("image/webp"))
	if err == nil {
		t.Error("ConvertFormat should fail with unsupported format")
	}
}

func TestImageCodec_formatToMIMEType_AllFormats(t *testing.T) {
	codec := NewImageCodec()

	formats := map[string]value.MIMEType{
		"png":  value.MIMETypeImagePNG,
		"jpeg": value.MIMETypeImageJPEG,
		"gif":  value.MIMETypeImageGIF,
		"bmp":  value.MIMETypeImageBMP,
	}

	for format, expected := range formats {
		got, err := codec.formatToMIMEType(format)
		if err != nil || got != expected {
			t.Errorf("formatToMIMEType(%s) = %v, %v; want %v", format, got, err, expected)
		}
	}

	unsupported := []string{"webp", "tiff", "unknown", ""}
	for _, format := range unsupported {
		_, err := codec.formatToMIMEType(format)
		if err == nil {
			t.Errorf("formatToMIMEType(%s) should fail", format)
		}
	}
}

func TestImageCodec_AllErrorPaths(t *testing.T) {
	codec := NewImageCodec()

	_, err := codec.EncodePNG(nil)
	if err == nil {
		t.Error("EncodePNG(nil) should fail")
	}

	_, err = codec.EncodeJPEG(nil, 90)
	if err == nil {
		t.Error("EncodeJPEG(nil) should fail")
	}

	_, err = codec.EncodeGIF(nil)
	if err == nil {
		t.Error("EncodeGIF(nil) should fail")
	}

	_, err = codec.DecodePNG([]byte{})
	if err == nil {
		t.Error("DecodePNG(empty) should fail")
	}

	_, err = codec.DecodePNG([]byte{0x00, 0x01})
	if err == nil {
		t.Error("DecodePNG(invalid) should fail")
	}

	_, err = codec.DecodeJPEG([]byte{})
	if err == nil {
		t.Error("DecodeJPEG(empty) should fail")
	}

	_, err = codec.DecodeJPEG([]byte{0x00, 0x01})
	if err == nil {
		t.Error("DecodeJPEG(invalid) should fail")
	}

	_, err = codec.DecodeGIF([]byte{})
	if err == nil {
		t.Error("DecodeGIF(empty) should fail")
	}

	_, err = codec.DecodeGIF([]byte{0x00, 0x01})
	if err == nil {
		t.Error("DecodeGIF(invalid) should fail")
	}

	_, _, err = codec.Decode([]byte{})
	if err == nil {
		t.Error("Decode(empty) should fail")
	}

	_, _, err = codec.Decode([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Error("Decode(invalid) should fail")
	}

	_, err = codec.DetectFormat([]byte{})
	if err == nil {
		t.Error("DetectFormat(empty) should fail")
	}

	_, err = codec.DetectFormat([]byte{0x00, 0x01})
	if err == nil {
		t.Error("DetectFormat(invalid) should fail")
	}

	_, err = codec.ConvertFormat([]byte{0x00}, value.MIMETypeImagePNG)
	if err == nil {
		t.Error("ConvertFormat(invalid source) should fail")
	}
}

func TestImageCodec_DetectFormat_ShortData(t *testing.T) {
	codec := NewImageCodec()

	shortData := []byte{0x00}
	_, err := codec.DetectFormat(shortData)
	if err == nil {
		t.Error("DetectFormat should fail with very short data")
	}

	twoBytes := []byte{0xFF, 0xD8}
	_, err = codec.DetectFormat(twoBytes)
	if err == nil {
		t.Error("DetectFormat should fail with incomplete JPEG header")
	}
}
