package service

import "errors"

// ValidationFunc is a function that validates input content.
// It returns nil if the content is valid, or an error describing the problem.
type ValidationFunc func(string) error

// ValidationService provides common validation functions and helpers.
type ValidationService struct{}

// NewValidationService creates a new validation service.
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// Validate runs the validator on the content and returns any error.
// If validator is nil, returns nil (no validation).
func (s *ValidationService) Validate(content string, validator ValidationFunc) error {
	if validator == nil {
		return nil
	}
	return validator(content)
}

// IsValid returns true if the content passes validation.
func (s *ValidationService) IsValid(content string, validator ValidationFunc) bool {
	return s.Validate(content, validator) == nil
}

// Common validation errors.
var (
	ErrEmpty         = errors.New("content cannot be empty")
	ErrTooShort      = errors.New("content is too short")
	ErrTooLong       = errors.New("content is too long")
	ErrInvalidFormat = errors.New("content has invalid format")
)

// NotEmpty returns a validator that ensures content is not empty.
func NotEmpty() ValidationFunc {
	return func(content string) error {
		if content == "" {
			return ErrEmpty
		}
		return nil
	}
}

// MinLength returns a validator that ensures content has minimum length.
func MinLength(minLen int) ValidationFunc {
	return func(content string) error {
		if len(content) < minLen {
			return ErrTooShort
		}
		return nil
	}
}

// MaxLength returns a validator that ensures content doesn't exceed maximum length.
func MaxLength(maxLen int) ValidationFunc {
	return func(content string) error {
		if len(content) > maxLen {
			return ErrTooLong
		}
		return nil
	}
}

// Range returns a validator that ensures content length is within range.
func Range(minLen, maxLen int) ValidationFunc {
	return func(content string) error {
		length := len(content)
		if length < minLen {
			return ErrTooShort
		}
		if length > maxLen {
			return ErrTooLong
		}
		return nil
	}
}

// Chain combines multiple validators into one.
// Validators are executed in order, stopping at the first error.
func Chain(validators ...ValidationFunc) ValidationFunc {
	return func(content string) error {
		for _, validator := range validators {
			if err := validator(content); err != nil {
				return err
			}
		}
		return nil
	}
}
