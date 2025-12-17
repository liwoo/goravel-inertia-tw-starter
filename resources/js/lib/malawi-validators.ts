/**
 * Malawi-specific validation utilities for forms
 * Based on Malawi telecommunications and government standards
 */

/**
 * Validates Malawian phone numbers
 * Accepts formats: +265XXXXXXXXX (12 digits total including country code)
 * Valid prefixes after +265: 88, 99, 98, 31, 21, 1 (for landlines)
 * @param phone Phone number to validate
 * @returns boolean indicating if the phone number is valid
 */
export function validateMalawiPhone(phone: string): boolean {
  // Remove spaces, hyphens, and parentheses
  const cleaned = phone.replace(/[\s\-\(\)]/g, '');

  // Check if starts with +265
  if (cleaned.startsWith('+265')) {
    const withoutCode = cleaned.substring(4);

    // Should have exactly 9 digits after country code
    if (withoutCode.length !== 9) {
      return false;
    }

    // Check for valid mobile/landline prefixes
    // Mobile: 88X XXX XXX (TNM), 99X XXX XXX (Airtel), 98X XXX XXX (Access)
    // VOIP: 31X XXX XXX (TNM VOIP)
    // Landline: 1XX XXX XXX (MTL), 21X XXX XXX (ACL)
    const validPrefixes = ['88', '99', '98', '31', '21', '1'];
    const hasValidPrefix = validPrefixes.some(prefix => withoutCode.startsWith(prefix));

    if (!hasValidPrefix) {
      return false;
    }

    // Ensure all characters are digits
    return /^\d{9}$/.test(withoutCode);
  } else if (cleaned.startsWith('0')) {
    // Local format without country code
    const withoutZero = cleaned.substring(1);

    // Should have exactly 9 digits after removing 0
    if (withoutZero.length !== 9) {
      return false;
    }

    // Check for valid prefixes
    const validPrefixes = ['88', '99', '98', '31', '21', '1'];
    const hasValidPrefix = validPrefixes.some(prefix => withoutZero.startsWith(prefix));

    if (!hasValidPrefix) {
      return false;
    }

    // Ensure all characters are digits
    return /^\d{9}$/.test(withoutZero);
  }

  return false;
}

/**
 * Validates Malawian business registration numbers
 * Format: BRNR-XXXXXX where X is alphanumeric (e.g., BRNR-EP5CWE3)
 * @param regNumber Registration number to validate
 * @returns boolean indicating if the registration number is valid
 */
export function validateBusinessRegistration(regNumber: string): boolean {
  // Remove spaces and convert to uppercase
  const cleaned = regNumber.trim().toUpperCase();

  // Pattern: BRNR- followed by 6-7 alphanumeric characters
  const pattern = /^BRNR-[A-Z0-9]{6,7}$/;
  return pattern.test(cleaned);
}

/**
 * Validates Malawian Tax Identification Numbers (TIN)
 * Format: 8 digits (e.g., 70543634)
 * @param tin Tax Identification Number to validate
 * @returns boolean indicating if the TIN is valid
 */
export function validateTIN(tin: string): boolean {
  // Remove spaces and hyphens
  const cleaned = tin.replace(/[\s\-]/g, '');

  // Should be exactly 8 digits
  const pattern = /^\d{8}$/;
  return pattern.test(cleaned);
}

/**
 * Validates Malawian National ID numbers
 * Format: 8 alphanumeric characters (e.g., T6N8SARR)
 * @param nationalID National ID to validate
 * @returns boolean indicating if the National ID is valid
 */
export function validateNationalID(nationalID: string): boolean {
  // Remove spaces and convert to uppercase
  const cleaned = nationalID.trim().toUpperCase();

  // Should be exactly 8 alphanumeric characters
  const pattern = /^[A-Z0-9]{8}$/;
  return pattern.test(cleaned);
}

/**
 * Formats a phone number to the standard +265 format
 * @param phone Phone number to format
 * @returns Formatted phone number
 */
export function formatMalawiPhone(phone: string): string {
  // Remove all non-digit characters except +
  let cleaned = phone.replace(/[^\d+]/g, '');

  // If starts with 0, replace with +265
  if (cleaned.startsWith('0')) {
    cleaned = '+265' + cleaned.substring(1);
  }

  // If doesn't start with +265 and has 9 digits, add +265
  if (!cleaned.startsWith('+265') && cleaned.length === 9) {
    cleaned = '+265' + cleaned;
  }

  return cleaned;
}

/**
 * Formats a business registration number to standard format
 * @param regNumber Registration number to format
 * @returns Formatted registration number
 */
export function formatBusinessRegistration(regNumber: string): string {
  // Remove spaces and convert to uppercase
  let cleaned = regNumber.trim().toUpperCase();

  // If doesn't start with BRNR- and has at least 6 characters, add BRNR-
  if (!cleaned.startsWith('BRNR-') && cleaned.length >= 6) {
    cleaned = 'BRNR-' + cleaned;
  }

  return cleaned;
}

// Zod refinement functions for use in schemas
export const malawiPhoneRefinement = (val: string) => validateMalawiPhone(val);
export const businessRegistrationRefinement = (val: string) => !val || val === '' || validateBusinessRegistration(val);
export const tinRefinement = (val: string) => !val || val === '' || validateTIN(val);
export const nationalIDRefinement = (val: string) => validateNationalID(val);

// Error messages for validation
export const VALIDATION_MESSAGES = {
  PHONE: 'Must be a valid Malawian phone number (e.g., +265881234567 or 0881234567)',
  BUSINESS_REG: 'Must be in format BRNR-XXXXXX (e.g., BRNR-EP5CWE3)',
  TIN: 'Must be exactly 8 digits (e.g., 70543634)',
  NATIONAL_ID: 'Must be exactly 8 alphanumeric characters (e.g., T6N8SARR)',
} as const;

// Example valid formats for UI hints
export const EXAMPLE_FORMATS = {
  PHONE: '+265881234567',
  BUSINESS_REG: 'BRNR-EP5CWE3',
  TIN: '70543634',
  NATIONAL_ID: 'T6N8SARR',
} as const;