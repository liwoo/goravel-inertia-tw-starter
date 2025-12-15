/**
 * Two-Factor Authentication API Client
 *
 * This module provides typed API functions for TOTP 2FA operations.
 */

import axios from '@/lib/axios';
import type {
  TOTPStatus,
  TOTPSetupResponse,
  TOTPVerifyResponse,
  TOTPBackupCodesResponse,
} from '@/types/account';

/**
 * Get the current 2FA status for the authenticated user
 */
export const getTOTPStatus = () => {
  return axios.get<TOTPStatus>('/api/2fa/status');
};

/**
 * Initiate TOTP setup - returns QR code and secret
 * @param password - User's current password for verification
 */
export const setupTOTP = (password: string) => {
  return axios.post<TOTPSetupResponse>('/api/2fa/setup', { password });
};

/**
 * Verify TOTP code and complete 2FA setup
 * @param code - 6-digit TOTP code from authenticator app
 */
export const verifyTOTP = (code: string) => {
  return axios.post<TOTPVerifyResponse>('/api/2fa/verify', { code });
};

/**
 * Disable 2FA for the authenticated user
 * @param password - User's current password
 * @param code - Current TOTP code for verification
 */
export const disableTOTP = (password: string, code: string) => {
  return axios.post('/api/2fa/disable', { password, code });
};

/**
 * Regenerate backup codes
 * @param password - User's current password for verification
 */
export const regenerateBackupCodes = (password: string) => {
  return axios.post<TOTPBackupCodesResponse>('/api/2fa/backup-codes', { password });
};

/**
 * Verify 2FA code during login
 * @param tempToken - Temporary token from initial login response
 * @param code - 6-digit TOTP code or backup code
 */
export const verify2FALogin = (tempToken: string, code: string) => {
  return axios.post('/api/auth/verify-2fa', {
    temp_token: tempToken,
    code
  });
};
