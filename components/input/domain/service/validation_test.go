package service

import (
	"errors"
	"testing"
)

func TestValidationService_Validate(t *testing.T) {
	svc := NewValidationService()

	alwaysValid := func(s string) error { return nil }
	alwaysInvalid := func(s string) error { return errors.New("invalid") }

	tests := []struct {
		name      string
		content   string
		validator ValidationFunc
		wantErr   bool
	}{
		{"nil validator passes", "anything", nil, false},
		{"valid content passes", "test", alwaysValid, false},
		{"invalid content fails", "test", alwaysInvalid, true},
		{"empty with nil validator", "", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Validate(tt.content, tt.validator)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationService_IsValid(t *testing.T) {
	svc := NewValidationService()

	alwaysValid := func(s string) error { return nil }
	alwaysInvalid := func(s string) error { return errors.New("invalid") }

	tests := []struct {
		name      string
		content   string
		validator ValidationFunc
		want      bool
	}{
		{"nil validator is valid", "anything", nil, true},
		{"valid content", "test", alwaysValid, true},
		{"invalid content", "test", alwaysInvalid, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.IsValid(tt.content, tt.validator)
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotEmpty(t *testing.T) {
	validator := NotEmpty()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"non-empty passes", "hello", false},
		{"empty fails", "", true},
		{"whitespace passes", " ", false},
		{"single char passes", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotEmpty()(%q) error = %v, wantErr %v", tt.content, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrEmpty) {
				t.Errorf("expected ErrEmpty, got %v", err)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	validator := MinLength(5)

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"equal to min passes", "12345", false},
		{"above min passes", "123456", false},
		{"below min fails", "1234", true},
		{"empty fails", "", true},
		{"way above passes", "1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinLength(5)(%q) error = %v, wantErr %v", tt.content, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrTooShort) {
				t.Errorf("expected ErrTooShort, got %v", err)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	validator := MaxLength(5)

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"equal to max passes", "12345", false},
		{"below max passes", "1234", false},
		{"above max fails", "123456", true},
		{"empty passes", "", false},
		{"way above fails", "1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxLength(5)(%q) error = %v, wantErr %v", tt.content, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrTooLong) {
				t.Errorf("expected ErrTooLong, got %v", err)
			}
		})
	}
}

func TestRange(t *testing.T) {
	validator := Range(3, 7)

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"at min passes", "123", false},
		{"at max passes", "1234567", false},
		{"in range passes", "12345", false},
		{"below min fails", "12", true},
		{"above max fails", "12345678", true},
		{"empty fails", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Range(3,7)(%q) error = %v, wantErr %v", tt.content, err, tt.wantErr)
			}
		})
	}
}

func TestChain(t *testing.T) {
	// Create a chain of validators
	validator := Chain(
		NotEmpty(),
		MinLength(3),
		MaxLength(10),
	)

	tests := []struct {
		name    string
		content string
		wantErr bool
		errType error
	}{
		{"valid content", "hello", false, nil},
		{"empty fails on first", "", true, ErrEmpty},
		{"too short fails on second", "ab", true, ErrTooShort},
		{"too long fails on third", "12345678901", true, ErrTooLong},
		{"at boundaries passes", "abc", false, nil},
		{"at max boundary passes", "1234567890", false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chain()(%q) error = %v, wantErr %v", tt.content, err, tt.wantErr)
			}
			if tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("expected error type %v, got %v", tt.errType, err)
			}
		})
	}
}

func TestChain_StopsAtFirstError(t *testing.T) {
	callCount := 0
	countingValidator := func(s string) error {
		callCount++
		return errors.New("error")
	}

	validator := Chain(
		countingValidator,
		countingValidator,
		countingValidator,
	)

	_ = validator("test")

	// Should only call first validator
	if callCount != 1 {
		t.Errorf("Chain called %d validators, expected 1 (should stop at first error)", callCount)
	}
}

func TestChain_Empty(t *testing.T) {
	validator := Chain()
	err := validator("anything")
	if err != nil {
		t.Errorf("Chain() with no validators should always pass, got error: %v", err)
	}
}
