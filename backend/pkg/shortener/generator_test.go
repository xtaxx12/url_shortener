package shortener

import (
	"testing"
)

func TestGenerateCode_DefaultLength(t *testing.T) {
	code := GenerateCode(0)
	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}
}

func TestGenerateCode_CustomLength(t *testing.T) {
	code := GenerateCode(10)
	if len(code) != 10 {
		t.Errorf("Expected length 10, got %d", len(code))
	}
}

func TestGenerateCode_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		code := GenerateCode(7)
		if codes[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		codes[code] = true
	}
}

func TestGenerateCode_ValidCharacters(t *testing.T) {
	code := GenerateCode(100)
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9')) {
			t.Errorf("Invalid character in code: %c", char)
		}
	}
}

func TestIsValidCode_ValidCodes(t *testing.T) {
	validCodes := []string{"abc123", "ABC123", "abcDEF", "a1b2c3d"}
	for _, code := range validCodes {
		if !IsValidCode(code) {
			t.Errorf("Expected %s to be valid", code)
		}
	}
}

func TestIsValidCode_InvalidCodes(t *testing.T) {
	invalidCodes := []string{"ab", "abc-123", "abc_123", "", "verylongcodethatexceedslimit123"}
	for _, code := range invalidCodes {
		if IsValidCode(code) {
			t.Errorf("Expected %s to be invalid", code)
		}
	}
}

func BenchmarkGenerateCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateCode(7)
	}
}
