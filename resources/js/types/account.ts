/**
 * Account page TypeScript types
 * These types define the data structures for the account/profile management features
 */

// Profile data returned from GET /api/account/profile
export interface ProfileData {
  id: number;
  name: string;
  email: string;
  role: string;
  is_active: boolean;
  email_verified_at?: string;
  created_at: string;
  updated_at?: string;
  last_login_at?: string;
  last_login_ip?: string;
  roles: Array<{
    id: number;
    name: string;
    slug?: string;
  }>;
}

// Request payload for PUT /api/account/profile
export interface UpdateProfileRequest {
  name: string;
  email: string;
}

// Request payload for PUT /api/account/password
export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
  confirm_password: string;
}

// Activity item from GET /api/account/activities
export interface ActivityItem {
  id: number;
  type: string;
  description: string;
  ip_address?: string;
  user_agent?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
}

// Activity type from GET /api/account/activities/types
export interface ActivityType {
  value: string;
  label: string;
}

// Activity list response with pagination
export interface ActivityListResponse {
  data: ActivityItem[];
  meta: {
    current_page: number;
    per_page: number;
    total: number;
    total_pages: number;
    from: number;
    to: number;
  };
}

// Activity summary from GET /api/account/activities/summary
export interface ActivitySummary {
  total_activities: number;
  activities_by_type: Record<string, number>;
  last_activity_at?: string;
}

// Activity filter/request params for GET /api/account/activities
export interface ActivityListRequest {
  page?: number;
  pageSize?: number;
  type?: string;
}

// Form error types for API responses
export interface ProfileFormErrors {
  name?: string;
  email?: string;
}

export interface PasswordFormErrors {
  current_password?: string;
  new_password?: string;
  confirm_password?: string;
}

// API error response structure
export interface ApiErrorResponse {
  message: string;
  errors?: Record<string, string[]>;
}

// =============================================================================
// TOTP Two-Factor Authentication Types
// =============================================================================

// Generic API response wrapper
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

// Status data from GET /api/2fa/status
export interface TOTPStatusData {
  enabled: boolean;
  verified_at: string | null;
  backup_codes_remaining: number;
}

// Status response (wrapped)
export type TOTPStatus = ApiResponse<TOTPStatusData>;

// Setup data from POST /api/2fa/setup
export interface TOTPSetupData {
  qr_code: string;  // base64 PNG data URI
  secret: string;   // Manual entry key
}

// Setup response (wrapped)
export type TOTPSetupResponse = ApiResponse<TOTPSetupData>;

// Verify data from POST /api/2fa/verify
export interface TOTPVerifyData {
  backup_codes: string[];
  enabled_at: string;
  redirect?: string;
}

// Verify response (wrapped)
export type TOTPVerifyResponse = ApiResponse<TOTPVerifyData>;

// Backup codes data from POST /api/2fa/backup-codes
export interface TOTPBackupCodesData {
  backup_codes: string[];
  generated_at: string;
}

// Backup codes response (wrapped)
export type TOTPBackupCodesResponse = ApiResponse<TOTPBackupCodesData>;

// Request payloads
export interface TOTPSetupRequest {
  password: string;
}

export interface TOTPVerifyRequest {
  code: string;
}

export interface TOTPDisableRequest {
  password: string;
  code: string;
}

export interface TOTPBackupCodesRequest {
  password: string;
}

// Login 2FA verification
export interface TwoFactorLoginRequest {
  temp_token: string;
  code: string;
}

// Login response when 2FA is required
export interface TwoFactorRequiredResponse {
  requires_2fa: boolean;
  temp_token: string;
}

// Form error types for 2FA
export interface TOTPFormErrors {
  password?: string;
  code?: string;
  general?: string;
}
