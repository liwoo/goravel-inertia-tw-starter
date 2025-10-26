# 12. Troubleshooting

## 🛠️ Overview

This section provides solutions to common problems users may encounter while using the SMEDI SME Database. Issues are organized by category with step-by-step solutions and escalation procedures.

## 🔐 Login and Authentication Issues

### Cannot Access Login Page

#### Problem: Browser shows "Site can't be reached" or similar error
**Possible Causes:**
- Internet connection problems
- Server maintenance
- Incorrect URL
- Firewall blocking access

**Solutions:**
1. **Check Internet Connection**
   - Try accessing other websites
   - Restart your router/modem
   - Contact your internet service provider

2. **Verify URL**
   - Ensure you're using: `https://sme-database.gov.mw`
   - Check for typos in the address
   - Try accessing from a bookmark or saved link

3. **Clear Browser Data**
   ```
   Chrome: Settings > Privacy and Security > Clear Browsing Data
   Firefox: Settings > Privacy & Security > Clear Data
   Edge: Settings > Privacy, Search, and Services > Clear Browsing Data
   ```

4. **Try Different Browser**
   - Test with Chrome, Firefox, Edge, or Safari
   - Ensure browser is updated to latest version

### Login Credentials Not Working

#### Problem: "Invalid credentials" error message
**Step-by-Step Solution:**

1. **Verify Credentials**
   - Double-check email address spelling
   - Ensure password is correct (use eye icon to verify)
   - Check if Caps Lock is enabled

2. **Reset Password**
   - Click "Forgot Password?" on login page
   - Enter your registered email address
   - Check email for reset link (including spam folder)
   - Follow reset instructions

3. **Account Status Check**
   - Contact administrator to verify account is active
   - Ensure account hasn't been suspended
   - Confirm email address is registered in system

#### Problem: Account locked due to failed attempts
**Solution:**
- **Automatic Unlock**: Wait 30 minutes for automatic unlock
- **Manual Unlock**: Contact system administrator
- **Prevention**: Use password manager to avoid future lockouts

### Two-Factor Authentication Issues

#### Problem: 2FA code not accepted
**Solutions:**

1. **Time Synchronization**
   - Ensure device time is correct
   - Sync authenticator app time
   - Account for time zone differences

2. **Code Timing**
   - Use fresh code (codes expire every 30 seconds)
   - Don't reuse previously entered codes
   - Wait for new code if current one is about to expire

3. **Backup Options**
   - Use backup codes if available
   - Contact administrator for 2FA reset
   - Verify authenticator app is working correctly

## 🌐 Browser and Display Issues

### Page Loading Problems

#### Problem: Pages load slowly or incompletely
**Diagnostic Steps:**

1. **Check Internet Speed**
   - Run speed test (minimum 1 Mbps recommended)
   - Close other bandwidth-intensive applications
   - Try different network if possible

2. **Browser Optimization**
   ```
   Steps to optimize browser:
   1. Close unnecessary tabs
   2. Clear cache and cookies
   3. Disable unnecessary extensions
   4. Update browser to latest version
   5. Restart browser
   ```

3. **System Resources**
   - Close other applications
   - Check available RAM (minimum 4GB recommended)
   - Restart computer if necessary

#### Problem: Page displays incorrectly or missing elements
**Solutions:**

1. **Browser Compatibility**
   - Use supported browsers: Chrome, Firefox, Edge, Safari
   - Update to latest browser version
   - Enable JavaScript and cookies

2. **Clear Browser Cache**
   ```
   Quick cache clear:
   - Chrome: Ctrl+Shift+Delete
   - Firefox: Ctrl+Shift+Delete
   - Edge: Ctrl+Shift+Delete
   - Safari: Cmd+Option+E
   ```

3. **Disable Extensions**
   - Temporarily disable ad blockers
   - Turn off privacy extensions
   - Test in incognito/private mode

### Mobile Display Issues

#### Problem: Interface not responsive on mobile devices
**Solutions:**

1. **Browser Settings**
   - Ensure mobile view is enabled
   - Clear mobile browser cache
   - Update mobile browser

2. **Device Orientation**
   - Try both portrait and landscape modes
   - Refresh page after orientation change
   - Check if specific features require landscape mode

3. **Mobile-Specific Troubleshooting**
   - Force close and reopen browser app
   - Restart mobile device
   - Clear app data for browser

## 📊 Data Entry and Form Issues

### Form Submission Problems

#### Problem: "Form submission failed" or similar errors
**Step-by-Step Resolution:**

1. **Check Required Fields**
   - Look for red error messages
   - Ensure all required fields (marked with *) are completed
   - Verify field format requirements (email, phone, etc.)

2. **Validate Data Formats**
   - **Email**: Must include @ and valid domain
   - **Phone**: Follow expected format (+265-XXX-XXXX)
   - **Dates**: Use date picker or DD/MM/YYYY format
   - **Numbers**: Ensure numeric fields contain only numbers

3. **File Upload Issues**
   - Check file size limits (typically 10MB maximum)
   - Verify file types are supported (PDF, DOC, JPG, PNG)
   - Ensure stable internet connection for uploads

#### Problem: Data not saving properly
**Troubleshooting Steps:**

1. **Session Timeout**
   - Save work frequently (Ctrl+S)
   - Check for session timeout warnings
   - Log out and log back in

2. **Permission Issues**
   - Verify you have permission to edit the record
   - Check if record is locked by another user
   - Contact administrator if permissions seem incorrect

3. **Data Validation**
   - Review all error messages carefully
   - Ensure data meets system requirements
   - Check for duplicate entries (registration numbers, emails)

### Search and Filter Problems

#### Problem: Search returns no results when data should exist
**Diagnostic Process:**

1. **Search Criteria**
   - Check spelling of search terms
   - Try partial search terms
   - Use wildcard characters if supported

2. **Filter Settings**
   - Clear all active filters
   - Check date range filters
   - Verify geographic or category filters

3. **Permission Scope**
   - Ensure you have permission to view the data
   - Check if data falls within your assigned scope
   - Contact administrator to verify access rights

#### Problem: Filters not working as expected
**Solutions:**

1. **Reset Filters**
   - Click "Clear All Filters" button
   - Refresh page to reset filter state
   - Apply filters one at a time to identify issues

2. **Filter Conflicts**
   - Avoid conflicting filter combinations
   - Check for logical inconsistencies
   - Use broader criteria initially, then narrow down

## 🔒 Permission and Access Issues

### Access Denied Errors

#### Problem: "Access Denied" or "Insufficient Permissions"
**Resolution Steps:**

1. **Verify Role Assignment**
   - Check your current role in Profile Settings
   - Compare required permissions with your role
   - Contact administrator for role review

2. **Data Scope Limitations**
   - Ensure data falls within your geographic scope
   - Check organizational access restrictions
   - Verify data sensitivity level permissions

3. **Temporary Access**
   - Request temporary elevated permissions
   - Use shared account if available and permitted
   - Contact supervisor for urgent access needs

### Missing Menu Items or Features

#### Problem: Expected menu items or buttons are not visible
**Troubleshooting:**

1. **Role-Based Visibility**
   - Confirm your role includes access to the feature
   - Check if feature requires specific permissions
   - Review role description and capabilities

2. **Feature Availability**
   - Verify feature is available in your system version
   - Check if feature is temporarily disabled
   - Contact administrator for feature status

3. **Browser Issues**
   - Clear browser cache and reload
   - Try different browser
   - Disable browser extensions temporarily

## 📈 Report and Export Issues

### Report Generation Problems

#### Problem: Reports fail to generate or show errors
**Step-by-Step Fix:**

1. **Report Parameters**
   - Verify all required parameters are selected
   - Check date ranges are valid
   - Ensure data exists for selected criteria

2. **System Resources**
   - Try generating smaller reports first
   - Avoid peak usage times
   - Break large reports into smaller segments

3. **Export Format Issues**
   - Try different export formats (PDF, Excel, CSV)
   - Check if specific format requires additional software
   - Verify browser can handle downloads

#### Problem: Exported data appears incorrect or incomplete
**Solutions:**

1. **Data Filters**
   - Review applied filters before export
   - Ensure all required data is selected
   - Check column visibility settings

2. **Format Settings**
   - Verify export format supports your data types
   - Check date and number format settings
   - Ensure special characters are handled correctly

3. **Permission Restrictions**
   - Confirm you have export permissions for all data
   - Check if some fields are restricted in exports
   - Contact administrator for export permission review

## 💾 Data Import Issues

### File Upload Problems

#### Problem: Import file rejected or upload fails
**Common Solutions:**

1. **File Format Verification**
   ```
   Supported formats:
   - CSV (Comma-separated values)
   - Excel (.xlsx, .xls)
   - Text files (.txt)
   ```

2. **File Size and Structure**
   - Maximum file size: 10MB
   - Ensure column headers match template
   - Check for empty rows or columns
   - Verify character encoding (UTF-8 recommended)

3. **Data Validation**
   - Remove special characters from critical fields
   - Ensure dates are in correct format
   - Check for required field completeness
   - Validate unique identifiers (no duplicates)

#### Problem: Import completes but data appears incorrect
**Troubleshooting:**

1. **Column Mapping**
   - Verify column mapping during import
   - Check for switched or misaligned columns
   - Review data preview before confirming import

2. **Data Transformation**
   - Check date format conversions
   - Verify number format handling
   - Review text encoding issues

## 🔧 System Performance Issues

### Slow System Response

#### Problem: System responds slowly to clicks and actions
**Performance Optimization:**

1. **Browser Optimization**
   ```
   Browser cleanup checklist:
   □ Clear cache and cookies
   □ Close unnecessary tabs
   □ Disable unused extensions
   □ Update to latest version
   □ Restart browser
   ```

2. **System Resources**
   - Close other applications
   - Check available memory (RAM)
   - Restart computer if necessary
   - Run disk cleanup

3. **Network Optimization**
   - Check internet speed and stability
   - Try different network if possible
   - Contact IT support for network issues

### Timeout Issues

#### Problem: Operations timeout or system logs you out frequently
**Solutions:**

1. **Session Management**
   - Save work frequently (every 10-15 minutes)
   - Be aware of 30-minute inactivity timeout
   - Use "Keep me logged in" option if available

2. **Network Stability**
   - Ensure stable internet connection
   - Avoid operations during network issues
   - Use wired connection if possible

## 📞 Getting Additional Help

### Self-Service Resources

1. **In-System Help**
   - Click ? icons for contextual help
   - Use F1 key for help in current context
   - Check tooltips on form fields

2. **Documentation**
   - Review relevant manual sections
   - Check FAQ for common issues
   - Search manual for specific topics

### Contact Support

#### Level 1: Local Support
- **Supervisor or Team Lead**: First point of contact
- **Local IT Support**: Basic technical issues
- **Department Administrator**: Permission and access issues

#### Level 2: System Administration
- **Email**: support@sme-database.gov.mw
- **Phone**: +265-XXX-XXXX (Business hours: 8:00 AM - 5:00 PM)
- **Response Time**: 4-8 hours for non-critical issues

#### Level 3: Technical Support
- **Critical Issues**: 24/7 emergency support
- **System Outages**: Immediate response
- **Data Loss**: Emergency recovery procedures

### Information to Provide When Seeking Help

#### Essential Information
1. **User Details**
   - Your name and email address
   - User role and permissions
   - Department/organization

2. **Problem Description**
   - Specific error messages (screenshot if possible)
   - Steps taken before the problem occurred
   - What you were trying to accomplish

3. **Technical Details**
   - Browser type and version
   - Operating system
   - Date and time of issue
   - Whether problem is recurring

#### Helpful Screenshots
- Error messages
- Unexpected display issues
- Form validation problems
- Missing interface elements

### Emergency Procedures

#### Critical System Issues
1. **System Outage**
   - Check system status page if available
   - Contact emergency support number
   - Document critical work that needs immediate attention

2. **Data Loss or Corruption**
   - Stop using the system immediately
   - Contact technical support urgently
   - Provide details of last known good state

3. **Security Incidents**
   - Report suspected security breaches immediately
   - Change passwords if compromise suspected
   - Document suspicious activities

#### Escalation Process
```
User Reports Issue
        ↓
Local Support (4 hours)
        ↓
System Administrator (8 hours)
        ↓
Technical Support (24 hours)
        ↓
Vendor Support (48 hours)
```

---

**Previous**: [Keyboard Shortcuts](11-shortcuts.md) ← | **Next**: [FAQ](13-faq.md) →