package helpers

import (
	"regexp"
	"strings"
)

// MalawiValidationHelper provides validation for Malawi-specific formats
type MalawiValidationHelper struct{}

// NewMalawiValidator creates a new MalawiValidationHelper instance
func NewMalawiValidator() *MalawiValidationHelper {
	return &MalawiValidationHelper{}
}

// ValidatePhoneNumber validates Malawian phone numbers
// Accepts formats: +265XXXXXXXXX (12 digits total including country code)
// Valid prefixes after +265: 88, 99, 98, 31, 21, 1 (for landlines)
func (v *MalawiValidationHelper) ValidatePhoneNumber(phone string) bool {
	// Remove spaces, hyphens, and parentheses
	phone = regexp.MustCompile(`[\s\-\(\)]`).ReplaceAllString(phone, "")

	// Check if starts with +265
	if strings.HasPrefix(phone, "+265") {
		phone = phone[4:] // Remove +265

		// Should have exactly 9 digits after country code
		if len(phone) != 9 {
			return false
		}

		// Check for valid mobile/landline prefixes
		// Mobile: 88X XXX XXX (TNM), 99X XXX XXX (Airtel), 98X XXX XXX (Access)
		// VOIP: 31X XXX XXX (TNM VOIP)
		// Landline: 1XX XXX XXX (MTL), 21X XXX XXX (ACL)
		validPrefixes := []string{"88", "99", "98", "31", "21", "1"}
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(phone, prefix) {
				// Ensure all characters are digits
				matched, _ := regexp.MatchString(`^\d{9}$`, phone)
				return matched
			}
		}
	} else if strings.HasPrefix(phone, "0") {
		// Local format without country code
		phone = phone[1:] // Remove leading 0

		// Should have exactly 9 digits after removing 0
		if len(phone) != 9 {
			return false
		}

		// Check for valid prefixes
		validPrefixes := []string{"88", "99", "98", "31", "21", "1"}
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(phone, prefix) {
				matched, _ := regexp.MatchString(`^\d{9}$`, phone)
				return matched
			}
		}
	}

	return false
}

// ValidateBusinessRegistration validates Malawian business registration numbers
// Format: BRNR-XXXXXX where X is alphanumeric (e.g., BRNR-EP5CWE3)
func (v *MalawiValidationHelper) ValidateBusinessRegistration(regNumber string) bool {
	// Remove spaces and convert to uppercase
	regNumber = strings.ToUpper(strings.TrimSpace(regNumber))

	// Pattern: BRNR- followed by 6-7 alphanumeric characters
	pattern := `^BRNR-[A-Z0-9]{6,7}$`
	matched, _ := regexp.MatchString(pattern, regNumber)
	return matched
}

// ValidateTIN validates Malawian Tax Identification Numbers
// Format: 8 digits (e.g., 70543634)
func (v *MalawiValidationHelper) ValidateTIN(tin string) bool {
	// Remove spaces and hyphens
	tin = regexp.MustCompile(`[\s\-]`).ReplaceAllString(tin, "")

	// Should be exactly 8 digits
	pattern := `^\d{8}$`
	matched, _ := regexp.MatchString(pattern, tin)
	return matched
}

// ValidateNationalID validates Malawian National ID numbers
// Format: 8 alphanumeric characters (e.g., T6N8SARR)
func (v *MalawiValidationHelper) ValidateNationalID(nationalID string) bool {
	// Remove spaces and convert to uppercase
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))

	// Should be exactly 8 alphanumeric characters
	pattern := `^[A-Z0-9]{8}$`
	matched, _ := regexp.MatchString(pattern, nationalID)
	return matched
}

// FormatPhoneNumber formats a phone number to the standard +265 format
func (v *MalawiValidationHelper) FormatPhoneNumber(phone string) string {
	// Remove all non-digit characters except +
	phone = regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// If starts with 0, replace with +265
	if strings.HasPrefix(phone, "0") {
		phone = "+265" + phone[1:]
	}

	// If doesn't start with +265, add it
	if !strings.HasPrefix(phone, "+265") && len(phone) == 9 {
		phone = "+265" + phone
	}

	return phone
}

// FormatBusinessRegistration formats a business registration number to standard format
func (v *MalawiValidationHelper) FormatBusinessRegistration(regNumber string) string {
	// Remove spaces and convert to uppercase
	regNumber = strings.ToUpper(strings.TrimSpace(regNumber))

	// If doesn't start with BRNR-, add it
	if !strings.HasPrefix(regNumber, "BRNR-") && len(regNumber) >= 6 {
		regNumber = "BRNR-" + regNumber
	}

	return regNumber
}