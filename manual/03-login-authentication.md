# 3. Login and Authentication

## Overview

The SMEDI SME Database uses a secure authentication system to protect sensitive business and financial data. This section covers all aspects of logging in, password management, and account security.

## How to Login

### Step 1: Access the System
1. Open your web browser
2. Navigate to the system URL: `https://sme-database.gov.mw` (example for now since we've not moved to production yet)
3. You will be redirected to the login page automatically

### Step 2: Enter Your Credentials
1. **Email Address**: Enter your registered email address
2. **Password**: Enter your secure password
3. Click the **"Login"** button

### Step 3: Successful Login
- Upon successful authentication, you'll be redirected to your dashboard
- The system will display a welcome message with your name and role
- Your last login time will be shown for security verification
- You will only be able to access your account if you have a valid email address and password
## 🔑 Password Requirements

The system enforces strict password security to protect sensitive data:

### Required Elements
- ✅ **Minimum 8 characters** - Must be at least 8 characters long
- ✅ **1 Capital letter** - At least one uppercase letter (A-Z)
- ✅ **1 Number** - At least one digit (0-9)
- ✅ **1 Symbol** - At least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)
- ✅ **No spaces** - Spaces are not allowed in passwords

### Password Examples
| ❌ **Weak Passwords** | ✅ **Strong Passwords** |
|----------------------|------------------------|
| `password` | `MySecure123!` |
| `123456` | `Gov@Malawi2024` |
| `qwerty` | `SME#Database99` |
| `Password` | `Secure&Strong8` |

### Visual Password Validation
When typing your password, the system shows real-time validation:
- 🟢 **Green checkmarks** indicate requirements that are met
- 🔴 **Red X marks** indicate requirements that need attention
- The password field border changes color based on validation status

## Password Visibility Toggle

For your convenience, the password field includes a visibility toggle:
- **Eye icon (👁️)**: Click to show your password
- **Eye-off icon (👁️‍🗨️)**: Click to hide your password
- Useful for verifying your password entry without retyping

##  Security Features

### Account Protection
- **Failed Login Attempts**: Account is temporarily locked after 5 failed attempts
- **Session Timeout**: Automatic logout after 30 minutes of inactivity
- **Secure Transmission**: All login data is encrypted using HTTPS
- **Password Encryption**: Passwords are securely hashed and never stored in plain text

### Login Monitoring
The system tracks all login activities:
- **Login Time**: Records when you log in
- **IP Address**: Tracks the location of login attempts
- **Device Information**: Monitors which devices access your account
- **Failed Attempts**: Logs unsuccessful login attempts

## 🔄 Password Management

### Changing Your Password
1. Only the admin who assigned you can change your password.
2. Contact the admin to request a password change.

### Forgot Password Process
1. If you forget your password, also contact the admin to request a new password.
2. The admin will trigger a reset password link to your email address.
3. Use the link to reset your password.

### Password Reset Security
- Reset links expire after 1 hour for security
- Links can only be used once
- Email notifications are sent for all password changes
- System administrator can assist with password resets if needed

## Two-Factor Authentication (2FA)

### Enabling 2FA (For High-Privilege Accounts)
1. Go to **Profile Settings** → **Security**
2. Click **"Enable Two-Factor Authentication"**
3. Scan the QR code with your authenticator app
4. Enter the verification code from your app
5. Save your backup codes in a secure location

### Using 2FA During Login
1. Enter your email and password normally
2. When prompted, open your authenticator app
3. Enter the 6-digit code displayed
4. Complete login process

### Supported Authenticator Apps
- Google Authenticator
- Microsoft Authenticator
- Authy
- Any TOTP-compatible app

## Login from Different Devices

### Desktop/Laptop Access
- **Recommended browsers**: Chrome, Firefox, Edge, Safari
- **Screen resolution**: Minimum 1024x768 for optimal experience
- **Hardware**: Modern computer with stable internet connection

### Mobile Access
- **Responsive design** adapts to mobile screens
- **Touch-friendly** interface for mobile devices
- **Offline capability** for viewing cached data

### Security Considerations
- Always log out from shared or public computers
- Use private/incognito mode on public devices
- Verify the URL is correct before entering credentials
- Report suspicious login notifications immediately

## ⚠️ Security Best Practices

### For Users
1. **Never share passwords** with colleagues or others
2. **Use unique passwords** - don't reuse passwords from other systems
3. **Log out properly** when finished using the system
4. **Keep credentials confidential** - don't write down passwords
5. **Report security concerns** to your administrator immediately

### For Administrators
1. **Regular password audits** to ensure compliance
2. **Monitor failed login attempts** for suspicious activity
3. **Regular security training** for all users
4. **Update security policies** as needed
5. **Maintain backup access** methods for account recovery

## 🚨 Troubleshooting Login Issues

### Common Problems and Solutions

#### Cannot Remember Password
- **Solution**: Use the "Forgot Password" feature
- **Timeline**: Reset emails are typically delivered within 5 minutes
- **Backup**: Contact your system administrator if email doesn't arrive

#### Account Locked Due to Failed Attempts
- **Automatic unlock**: Account unlocks after 30 minutes
- **Manual unlock**: Contact system administrator for immediate unlock
- **Prevention**: Ensure you're using the correct password

#### Browser Compatibility Issues
- **Clear browser cache** and cookies
- **Update browser** to the latest version
- **Try different browser** as a temporary solution
- **Disable browser extensions** that might interfere

#### Network Connection Problems
- **Check internet connection** stability
- **Try different network** (mobile hotspot, different WiFi)
- **Contact IT support** if persistent connectivity issues
- **Use different device** to isolate the problem

### Error Messages and Meanings

| Error Message | Meaning | Solution |
|---------------|---------|----------|
| "Invalid credentials" | Wrong email or password | Verify spelling and try again |
| "Account locked" | Too many failed attempts | Wait 30 minutes or contact admin |
| "Session expired" | Session timeout reached | Log in again |
| "Access denied" | Account permissions issue | Contact administrator |
| "Connection error" | Network connectivity problem | Check internet connection |

## 📞 Getting Help

### Self-Service Options
- **Password reset** feature on login page
- **Browser troubleshooting** steps in this manual
- **FAQ section** for common issues

### Contact Support
- **System Administrator**: [admin@sme-database.gov.mw]
- **Technical Support**: [support@sme-database.gov.mw]
- **Phone Support**: +265-XXX-XXXX (Business hours: 8:00 AM - 5:00 PM)
- **Emergency Access**: 24/7 support for critical issues

---

**Previous**: [Getting Started](02-getting-started.md) ← | **Next**: [Roles and Permissions](04-roles-permissions.md) →
