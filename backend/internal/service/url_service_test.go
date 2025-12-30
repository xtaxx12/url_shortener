package service

import (
	"testing"
)

// ============================================
// Tests para validateURL - Edge Cases y Security
// ============================================

func TestValidateURL_ShouldAcceptValidHTTPSURL(t *testing.T) {
	// Arrange
	validURL := "https://www.google.com/search?q=test"

	// Act
	err := validateURL(validURL)

	// Assert
	if err != nil {
		t.Errorf("Expected valid HTTPS URL to pass, got error: %v", err)
	}
}

func TestValidateURL_ShouldAcceptValidHTTPURL(t *testing.T) {
	// Arrange
	validURL := "http://example.com/page"

	// Act
	err := validateURL(validURL)

	// Assert
	if err != nil {
		t.Errorf("Expected valid HTTP URL to pass, got error: %v", err)
	}
}

func TestValidateURL_ShouldRejectJavaScriptScheme(t *testing.T) {
	// Arrange - Security: Prevent XSS via javascript: URLs
	maliciousURL := "javascript:alert('XSS')"

	// Act
	err := validateURL(maliciousURL)

	// Assert
	if err == nil {
		t.Error("Expected javascript: URL to be rejected")
	}
}

func TestValidateURL_ShouldRejectDataScheme(t *testing.T) {
	// Arrange - Security: Prevent data: URL attacks
	maliciousURL := "data:text/html,<script>alert('XSS')</script>"

	// Act
	err := validateURL(maliciousURL)

	// Assert
	if err == nil {
		t.Error("Expected data: URL to be rejected")
	}
}

func TestValidateURL_ShouldRejectFileScheme(t *testing.T) {
	// Arrange - Security: Prevent local file access
	maliciousURL := "file:///etc/passwd"

	// Act
	err := validateURL(maliciousURL)

	// Assert
	if err == nil {
		t.Error("Expected file: URL to be rejected")
	}
}

func TestValidateURL_ShouldRejectVBScriptScheme(t *testing.T) {
	// Arrange - Security: Prevent VBScript attacks
	maliciousURL := "vbscript:msgbox('attack')"

	// Act
	err := validateURL(maliciousURL)

	// Assert
	if err == nil {
		t.Error("Expected vbscript: URL to be rejected")
	}
}

func TestValidateURL_ShouldRejectURLWithoutHost(t *testing.T) {
	// Arrange - Edge case: URL sin host
	invalidURL := "https://"

	// Act
	err := validateURL(invalidURL)

	// Assert
	if err == nil {
		t.Error("Expected URL without host to be rejected")
	}
}

func TestValidateURL_ShouldRejectFTPScheme(t *testing.T) {
	// Arrange - Only HTTP/HTTPS allowed
	ftpURL := "ftp://files.example.com/file.txt"

	// Act
	err := validateURL(ftpURL)

	// Assert
	if err == nil {
		t.Error("Expected FTP URL to be rejected, only HTTP/HTTPS allowed")
	}
}

func TestValidateURL_ShouldAcceptURLWithPort(t *testing.T) {
	// Arrange
	urlWithPort := "https://localhost:8080/api/test"

	// Act
	err := validateURL(urlWithPort)

	// Assert
	if err != nil {
		t.Errorf("Expected URL with port to pass, got error: %v", err)
	}
}

func TestValidateURL_ShouldAcceptURLWithQueryParams(t *testing.T) {
	// Arrange
	urlWithQuery := "https://example.com/search?q=hello&lang=es&page=1"

	// Act
	err := validateURL(urlWithQuery)

	// Assert
	if err != nil {
		t.Errorf("Expected URL with query params to pass, got error: %v", err)
	}
}

func TestValidateURL_ShouldAcceptURLWithFragment(t *testing.T) {
	// Arrange
	urlWithFragment := "https://example.com/page#section-1"

	// Act
	err := validateURL(urlWithFragment)

	// Assert
	if err != nil {
		t.Errorf("Expected URL with fragment to pass, got error: %v", err)
	}
}

func TestValidateURL_ShouldAcceptInternationalDomain(t *testing.T) {
	// Arrange - Edge case: Unicode domain
	internationalURL := "https://例え.jp/test"

	// Act
	err := validateURL(internationalURL)

	// Assert
	if err != nil {
		t.Errorf("Expected international domain to pass, got error: %v", err)
	}
}

func TestValidateURL_ShouldAcceptIPAddress(t *testing.T) {
	// Arrange
	ipURL := "http://192.168.1.1:8080/api"

	// Act
	err := validateURL(ipURL)

	// Assert
	if err != nil {
		t.Errorf("Expected IP address URL to pass, got error: %v", err)
	}
}

// ============================================
// Tabla de tests para cobertura exhaustiva
// ============================================

func TestValidateURL_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantError bool
	}{
		// Valid URLs
		{"valid https", "https://google.com", false},
		{"valid http", "http://example.com", false},
		{"url with path", "https://example.com/path/to/page", false},
		{"url with query", "https://example.com?q=test", false},
		
		// Invalid URLs - Security
		{"javascript scheme", "javascript:alert(1)", true},
		{"data scheme", "data:text/html,test", true},
		{"file scheme", "file:///etc/passwd", true},
		{"vbscript scheme", "vbscript:code", true},
		
		// Invalid URLs - Format
		{"ftp scheme", "ftp://files.com", true},
		{"mailto scheme", "mailto:test@example.com", true},
		{"no host", "https://", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			gotError := err != nil
			if gotError != tt.wantError {
				t.Errorf("validateURL(%q) error = %v, wantError = %v", tt.url, err, tt.wantError)
			}
		})
	}
}
