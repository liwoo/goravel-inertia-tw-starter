package unit

import (
	"smedi-sme-db/app/helpers"
	"testing"
)

func TestValidateBusinessRegistration(t *testing.T) {
	validator := helpers.NewMalawiValidator()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// BRN format - sole proprietorships/partnerships (7 alphanumeric)
		{"BRN valid uppercase", "BRN-ABC1234", true},
		{"BRN valid mixed case", "brn-abc1234", true},
		{"BRN valid with spaces", " BRN-ABC1234 ", true},
		{"BRN invalid - 6 chars", "BRN-ABC123", false},
		{"BRN invalid - 8 chars", "BRN-ABC12345", false},

		// COY format - incorporated companies (6 alphanumeric)
		{"COY valid uppercase", "COY-ABC123", true},
		{"COY valid lowercase", "coy-abc123", true},
		{"COY valid with spaces", " COY-ABC123 ", true},
		{"COY invalid - 5 chars", "COY-ABC12", false},
		{"COY invalid - 7 chars", "COY-ABC1234", false},

		// PVT format - private limited companies (8 alphanumeric)
		{"PVT valid uppercase", "PVT-ABCD1234", true},
		{"PVT valid lowercase", "pvt-abcd1234", true},
		{"PVT valid with spaces", " PVT-ABCD1234 ", true},
		{"PVT invalid - 7 chars", "PVT-ABC1234", false},
		{"PVT invalid - 9 chars", "PVT-ABCD12345", false},

		// BRNR legacy format (6-7 alphanumeric)
		{"BRNR valid 6 chars", "BRNR-EP5CWE", true},
		{"BRNR valid 7 chars", "BRNR-EP5CWE3", true},
		{"BRNR valid lowercase", "brnr-ep5cwe3", true},
		{"BRNR valid with spaces", " BRNR-EP5CWE3 ", true},
		{"BRNR invalid - 5 chars", "BRNR-ABC12", false},
		{"BRNR invalid - 8 chars", "BRNR-ABC12345", false},

		// Invalid formats
		{"Invalid prefix", "XYZ-ABC123", false},
		{"No prefix", "ABC1234", false},
		{"Empty string", "", false},
		{"Only prefix", "BRN-", false},
		{"Special characters", "BRN-ABC@123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateBusinessRegistration(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateBusinessRegistration(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatBusinessRegistration(t *testing.T) {
	validator := helpers.NewMalawiValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Already formatted - should return as-is (uppercase)
		{"BRN already formatted", "BRN-ABC1234", "BRN-ABC1234"},
		{"COY already formatted", "COY-ABC123", "COY-ABC123"},
		{"PVT already formatted", "PVT-ABCD1234", "PVT-ABCD1234"},
		{"BRNR already formatted", "BRNR-EP5CWE3", "BRNR-EP5CWE3"},

		// Lowercase - should be converted to uppercase
		{"BRN lowercase", "brn-abc1234", "BRN-ABC1234"},
		{"COY lowercase", "coy-abc123", "COY-ABC123"},
		{"PVT lowercase", "pvt-abcd1234", "PVT-ABCD1234"},
		{"BRNR lowercase", "brnr-ep5cwe3", "BRNR-EP5CWE3"},

		// With spaces - should be trimmed and uppercased
		{"BRN with spaces", " brn-abc1234 ", "BRN-ABC1234"},
		{"COY with spaces", " coy-abc123 ", "COY-ABC123"},

		// No valid prefix - should return uppercase without adding prefix
		{"No prefix", "abc1234", "ABC1234"},
		{"Invalid prefix", "xyz-abc123", "XYZ-ABC123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.FormatBusinessRegistration(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBusinessRegistration(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
