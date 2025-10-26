# SME Management Module - User Manual

## 📊 Overview

The SME (Small and Medium Enterprise) Management Module is the core component of the SMEDI Database system. This module allows authorized users to create, update, view, delete, and manage all SME records in the system. It provides comprehensive tools for data entry, validation, filtering, searching, and exporting SME information.

## 🎯 Module Objectives

- **Centralized SME Data Management**: Single source of truth for all SME information
- **Data Quality Assurance**: Built-in validation and verification processes
- **Efficient Data Entry**: Streamlined forms and bulk import capabilities
- **Advanced Search & Filtering**: Powerful tools to find and analyze SME data
- **Comprehensive Reporting**: Export and generate reports in multiple formats
- **Role-Based Access Control**: Secure access based on user permissions

## 🚀 Getting Started

### Accessing the SME Module

1. **Login** to the SMEDI Database system
2. Navigate to **SME Management** from the main menu
3. Select your desired action from the submenu:
   - **All SMEs**: View and manage existing records
   - **Add New SME**: Create new SME registration
   - **Import SMEs**: Bulk import SME data
   - **SME Categories**: Manage business categories
   - **Export Data**: Download SME information

### Required Permissions

To use the SME Management module, you need appropriate role permissions:
- **Member**: Basic viewing and limited editing
- **Moderator**: Full viewing and editing with approval workflow
- **Librarian**: Complete data management capabilities
- **Administrator**: Full system access including user management
- **Super Administrator**: Complete system control

## 📝 Creating New SME Records

### Step-by-Step Process

#### 1. Access the Creation Form
- Click **SME Management** → **Add New SME**
- The system will open a multi-step form

#### 2. Basic Business Information
**Required Fields (marked with *):**
- **Business Name***: Full legal name of the enterprise
- **Registration Number***: Official business registration number
- **Business Type***: Select from dropdown (Manufacturing, Agriculture, Services, etc.)
- **Founded Date**: Date when business was established
- **Business Status**: Active, Inactive, Pending, Suspended

**Example:**
```
Business Name: Malawi Agro Processing Ltd
Registration Number: BN2024001234
Business Type: Agriculture
Founded Date: 15/03/2020
Business Status: Active
```

#### 3. Contact Information
**Required Fields:**
- **Physical Address***: Complete business address
- **District***: Select from dropdown
- **Region***: Northern, Central, Southern
- **Phone Number***: Primary contact number
- **Email Address**: Business email (if available)
- **Website**: Business website (optional)

#### 4. Business Details
- **Industry Category**: Specific business sector
- **Number of Employees**: Current workforce size
- **Annual Revenue**: Estimated yearly revenue (optional)
- **Business Description**: Brief description of business activities
- **Products/Services**: Main offerings

#### 5. Owner/Manager Information
- **Owner Name***: Full name of business owner
- **Manager Name**: If different from owner
- **Owner Contact**: Personal contact information
- **ID Number**: National ID or passport number
- **Gender**: Male, Female, Other
- **Age Group**: Select appropriate range

#### 6. Financial Information (Optional)
- **Banking Institution**: Primary bank
- **Account Status**: Banked, Unbanked
- **Credit History**: Previous loan experience
- **Collateral**: Available security for loans

#### 7. Document Upload
**Supported Documents:**
- Business Registration Certificate (PDF, JPG, PNG)
- Tax Clearance Certificate
- Business License
- Bank Statements
- ID Copy of Owner/Manager

**File Requirements:**
- Maximum file size: 10MB per file
- Supported formats: PDF, DOC, DOCX, JPG, PNG
- Clear, readable documents only

#### 8. Review and Submit
- **Preview**: Review all entered information
- **Validation**: System checks for errors and missing required fields
- **Save Draft**: Save incomplete forms for later completion
- **Submit**: Final submission for approval (if workflow enabled)

### Form Validation Rules

#### Real-time Validation
- **Email Format**: Must contain @ and valid domain
- **Phone Number**: Must follow format +265-XXX-XXX-XXX
- **Registration Number**: Must be unique in the system
- **Required Fields**: Cannot be empty
- **Date Fields**: Must be valid dates in DD/MM/YYYY format

#### Business Rules
- **Duplicate Prevention**: System checks for existing businesses with same name/registration
- **Data Consistency**: Cross-field validation ensures logical consistency
- **Reference Data**: Dropdowns are populated from approved reference lists

### Error Handling
If validation fails:
1. **Red indicators** appear next to invalid fields
2. **Error messages** explain what needs to be corrected
3. **Form cannot be submitted** until all errors are resolved
4. **Draft saving** allows work to be saved while fixing errors

## 📋 Viewing and Managing SME Records

### SME List View

#### Accessing the List
- Navigate to **SME Management** → **All SMEs**
- Default view shows paginated list of all SMEs you have permission to view

#### List Display Options

**Table View (Default):**
```
┌──────────────────────────────────────────────────────────────────┐
│ SME Records                                            [Export]   │
├──────────────────────────────────────────────────────────────────┤
│ Name              │ Type        │ Location    │ Status  │ Actions  │
├──────────────────────────────────────────────────────────────────┤
│ ABC Trading Ltd   │ Commerce    │ Lilongwe    │ Active  │ [V][E][D]│
│ Sunrise Farms     │ Agriculture │ Blantyre    │ Active  │ [V][E][D]│
│ Tech Solutions    │ Services    │ Mzuzu       │ Pending │ [V][A][D]│
│ Craft Makers Co   │ Manufacturing│ Zomba      │ Active  │ [V][E][D]│
└──────────────────────────────────────────────────────────────────┘
```

**Card View:**
```
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ ABC Trading Ltd │ │ Sunrise Farms   │ │ Tech Solutions  │
│ Commerce        │ │ Agriculture     │ │ Services        │
│ Lilongwe        │ │ Blantyre        │ │ Mzuzu          │
│ Status: Active  │ │ Status: Active  │ │ Status: Pending │
│ 📅 2020        │ │ 📅 2019        │ │ 📅 2024        │
│ [View] [Edit]   │ │ [View] [Edit]   │ │ [View] [Approve]│
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

#### Action Buttons
- **[V] View**: Open detailed view (read-only)
- **[E] Edit**: Open edit form (if you have permission)
- **[D] Delete**: Remove record (with confirmation)
- **[A] Approve**: Approve pending records (Moderator+ only)

### Detailed SME View

#### Accessing Details
- Click **View** button or SME name in the list
- Opens comprehensive tabbed interface

#### Information Tabs

**Basic Information Tab:**
- Business name and registration details
- Contact information
- Business type and category
- Establishment date and status

**Financial Information Tab:**
- Banking details
- Revenue information
- Credit history
- Loan applications and status

**Location Tab:**
- Physical address
- Geographic coordinates (if available)
- Service areas
- Branch locations

**Documents Tab:**
- All uploaded documents
- Document verification status
- Download options
- Upload additional documents

**History Tab:**
- Creation and modification history
- Approval workflow status
- User activity log
- System-generated events

#### Quick Actions from Detail View
- **Edit**: Switch to edit mode
- **Print**: Generate printable version
- **Export**: Download as PDF or Excel
- **Share**: Email record details
- **Flag**: Report issues or inconsistencies

## ✏️ Updating SME Records

### Edit Process

#### 1. Access Edit Mode
- From SME list: Click **Edit** button
- From detail view: Click **Edit** button
- Ensure you have appropriate permissions

#### 2. Edit Form Interface
- **Same structure** as creation form
- **Pre-populated** with existing data
- **Track changes** highlights modified fields
- **Version control** maintains edit history

#### 3. Field-Level Permissions
Depending on your role:
- **Basic fields**: Name, contact info (Most users)
- **Financial data**: Revenue, banking (Librarian+)
- **Status changes**: Active/Inactive (Moderator+)
- **Verification flags**: Approved/Verified (Admin only)

#### 4. Change Tracking
System automatically tracks:
- **What changed**: Field-level change detection
- **Who changed it**: User identification
- **When changed**: Timestamp
- **Previous values**: Audit trail maintenance

### Bulk Edit Operations

#### Selecting Multiple Records
1. **Checkbox selection**: Check boxes next to SME names
2. **Select all**: Use "Select All" for current page
3. **Select filtered**: Select all records matching current filter

#### Available Bulk Operations
- **Status Update**: Change status for multiple SMEs
- **Category Assignment**: Update business categories
- **Tag Assignment**: Add or remove tags
- **Export Selected**: Export only selected records
- **Delete Selected**: Remove multiple records (with confirmation)

#### Bulk Edit Process
1. **Select records** using checkboxes
2. **Choose action** from bulk actions dropdown
3. **Configure changes** in popup dialog
4. **Preview changes** before applying
5. **Confirm and execute** bulk operation

### Approval Workflow

#### For Pending Records
If your organization uses approval workflows:

**Moderator Actions:**
- **Approve**: Mark record as approved and active
- **Reject**: Send back with comments for revision
- **Request Changes**: Specify required modifications

**Approval Process:**
1. **Review submitted data** for accuracy and completeness
2. **Verify supporting documents** are authentic and current
3. **Check for duplicates** against existing records
4. **Apply business rules** and policy compliance
5. **Make approval decision** with appropriate comments

## 🗑️ Deleting SME Records

### Soft Delete vs Hard Delete

#### Soft Delete (Recommended)
- **Marks record as inactive** rather than removing
- **Preserves data integrity** and audit trails
- **Allows recovery** if deletion was mistaken
- **Maintains relationships** with other system data

#### Hard Delete (Admin Only)
- **Permanently removes** record from database
- **Cannot be recovered** once executed
- **Breaks relationships** with related data
- **Used only for** data cleanup or compliance

### Deletion Process

#### Single Record Deletion
1. **Navigate to SME record** (list or detail view)
2. **Click Delete button** (trash icon)
3. **Confirm deletion** in popup dialog
4. **Provide reason** for deletion (optional)
5. **Execute deletion** with final confirmation

#### Bulk Deletion
1. **Select multiple records** using checkboxes
2. **Choose "Delete Selected"** from bulk actions
3. **Review selected records** in confirmation dialog
4. **Provide deletion reason** (required for bulk operations)
5. **Confirm bulk deletion** with password verification

### Safety Measures
- **Double confirmation**: Two-step confirmation process
- **Reason logging**: All deletions must include reason
- **Permission checks**: Only authorized users can delete
- **Backup notifications**: System creates backup before deletion
- **Audit trail**: Complete log of deletion activities

### Recovery Options
For accidentally deleted records:
- **Contact administrator** immediately
- **Provide deletion details** (time, record name, reason)
- **Recovery possible** within backup retention period
- **Data restoration** may take 24-48 hours

## 🔍 Search and Filtering

### Global Search

#### Quick Search
Located in the top header of SME Management:
- **Search box**: Type any text to search across all SME fields
- **Auto-complete**: Suggestions appear as you type
- **Recent searches**: Quick access to previous searches
- **Search scope**: Searches name, registration number, location, owner

#### Advanced Search
Click the gear icon (⚙️) next to search box:
- **Field-specific search**: Search in specific fields only
- **Operator selection**: Equals, contains, starts with, etc.
- **Multiple criteria**: Combine multiple search conditions
- **Date range search**: Search by date ranges
- **Saved searches**: Save frequently used search criteria

### Filter Panel

#### Accessing Filters
- **Filter sidebar**: Click "Filters" button to open/close
- **Quick filters**: Pre-defined filter buttons above the list
- **Active filters**: Display current filter status
- **Filter count**: Shows number of applied filters

#### Available Filters

**Business Type Filter:**
```
☐ Business Type
  ☑ Agriculture
  ☐ Manufacturing
  ☐ Services
  ☐ Commerce
  ☐ Technology
```

**Location Filter:**
```
☐ Region
  ☑ Northern
  ☐ Central
  ☐ Southern

☐ District
  ☑ Mzuzu
  ☐ Lilongwe
  ☐ Blantyre
```

**Status Filter:**
```
☐ Status
  ☑ Active
  ☐ Inactive
  ☐ Pending
  ☐ Suspended
```

**Size Filter:**
```
☐ Enterprise Size
  ☑ Micro (1-4 employees)
  ☐ Small (5-29 employees)
  ☐ Medium (30-99 employees)
```

**Date Filters:**
```
☐ Registration Date
  From: [DD/MM/YYYY]
  To:   [DD/MM/YYYY]

☐ Last Updated
  From: [DD/MM/YYYY]
  To:   [DD/MM/YYYY]
```

#### Filter Operations
- **Apply Filters**: Click "Apply" to execute filters
- **Clear Filters**: Remove all active filters
- **Reset to Default**: Return to system default filters
- **Save Filter Set**: Save current filters for future use
- **Load Saved Filters**: Apply previously saved filter combinations

### Search Results

#### Result Display
- **Relevance ranking**: Most relevant results appear first
- **Highlight matches**: Search terms highlighted in results
- **Result count**: Total number of matching records
- **Pagination**: Navigate through multiple pages of results

#### Refining Results
- **Additional filters**: Apply filters to search results
- **Sort options**: Sort results by various criteria
- **Export results**: Export filtered/searched results
- **Save search**: Save search criteria for future use

### Advanced Search Techniques

#### Wildcard Searches
- **Asterisk (*)**: Matches any number of characters
  - Example: `Agro*` finds "Agro", "Agriculture", "Agribusiness"
- **Question mark (?)**: Matches single character
  - Example: `B?nk` finds "Bank", "Bunk"

#### Boolean Searches
- **AND**: All terms must be present
  - Example: `agriculture AND northern`
- **OR**: Any term can be present
  - Example: `commerce OR trading`
- **NOT**: Exclude terms
  - Example: `farming NOT poultry`

#### Phrase Searches
- **Exact phrases**: Use quotes for exact matches
  - Example: `"Malawi Agro Processing"`
- **Field searches**: Search in specific fields
  - Example: `name:"ABC Trading"`

## 📤 Export and Reporting

### Export Options

#### Quick Export
From the SME list view:
- **Export All**: Download all visible records
- **Export Selected**: Download only selected records
- **Export Filtered**: Download records matching current filters

#### Export Formats

**Excel (.xlsx):**
- **Best for**: Data analysis and manipulation
- **Includes**: All fields with proper formatting
- **Features**: Multiple worksheets, charts, formulas
- **File size**: Larger but comprehensive

**CSV (.csv):**
- **Best for**: Data import into other systems
- **Includes**: Raw data in comma-separated format
- **Features**: Universal compatibility
- **File size**: Smallest, fastest download

**PDF (.pdf):**
- **Best for**: Printing and sharing reports
- **Includes**: Formatted tables and summaries
- **Features**: Professional layout, charts, headers
- **File size**: Medium, print-ready

**JSON (.json):**
- **Best for**: Technical integrations and APIs
- **Includes**: Structured data format
- **Features**: Machine-readable, nested data
- **File size**: Small, efficient

### Custom Export Configuration

#### Field Selection
Choose which fields to include in export:
```
☑ Basic Information
  ☑ Business Name
  ☑ Registration Number
  ☑ Business Type
  ☐ Founded Date

☑ Contact Information
  ☑ Address
  ☑ Phone
  ☐ Email
  ☐ Website

☐ Financial Information
  ☐ Revenue
  ☐ Banking Details
  ☐ Credit History
```

#### Export Settings
- **Date format**: DD/MM/YYYY, MM/DD/YYYY, YYYY-MM-DD
- **Number format**: Include currency symbols, decimal places
- **Text encoding**: UTF-8, ASCII, local encoding
- **Column headers**: Include/exclude, custom names
- **Empty fields**: Show as blank, "N/A", or custom text

### Report Generation

#### Standard Reports

**SME Summary Report:**
- Overview of all SMEs by region, type, status
- Statistical summaries and charts
- Growth trends and analytics

**Regional Distribution Report:**
- SME distribution across regions and districts
- Geographic analysis and mapping
- Concentration patterns

**Industry Analysis Report:**
- SME breakdown by business type and industry
- Sector performance analysis
- Market share distribution

**Status Summary Report:**
- Active vs inactive SMEs
- Approval status tracking
- Workflow performance metrics

#### Custom Report Builder

**Step 1: Report Type Selection**
- Summary report with charts
- Detailed list with filters
- Comparative analysis
- Trend analysis over time

**Step 2: Data Selection**
- Choose SME fields to include
- Select aggregation methods (count, sum, average)
- Define grouping criteria
- Set date ranges

**Step 3: Visualization**
- Chart types: Bar, pie, line, scatter
- Color schemes and themes
- Layout options
- Interactive features

**Step 4: Output Format**
- PDF with charts and tables
- Excel with multiple worksheets
- HTML for web viewing
- PowerPoint for presentations

### Scheduled Exports

#### Setting Up Automated Exports
1. **Configure export parameters** (fields, format, filters)
2. **Set schedule** (daily, weekly, monthly)
3. **Choose delivery method** (email, FTP, cloud storage)
4. **Test configuration** with sample export
5. **Activate schedule** and monitor

#### Schedule Options
- **Daily**: Every day at specified time
- **Weekly**: Specific day of week
- **Monthly**: Specific date of month
- **Quarterly**: End of each quarter
- **Custom**: Complex scheduling rules

## 🔐 Security and Access Control

### User Permissions

#### View Permissions
- **Public data**: Basic business information
- **Restricted data**: Financial and sensitive information
- **Geographic scope**: Limited to specific regions/districts
- **Status-based**: Access to active vs all records

#### Edit Permissions
- **Own records**: Edit records you created
- **Department records**: Edit within your organization
- **All records**: Edit any SME record (Librarian+)
- **Approval rights**: Approve pending submissions

#### Administrative Permissions
- **User management**: Create and manage user accounts
- **System settings**: Configure system parameters
- **Audit access**: View system logs and user activity
- **Data backup**: Access to backup and recovery functions

### Data Security

#### Audit Trail
System automatically logs:
- **User actions**: Create, read, update, delete operations
- **Data changes**: Before and after values for all modifications
- **Access attempts**: Successful and failed login attempts
- **Export activities**: What data was exported and by whom
- **System events**: Backups, maintenance, configuration changes

#### Access Logging
- **Session tracking**: Monitor user sessions and activity
- **IP address logging**: Track access locations
- **Device fingerprinting**: Identify access devices
- **Unusual activity detection**: Alert on suspicious patterns

### Data Privacy

#### Personal Information Protection
- **Consent tracking**: Record consent for data processing
- **Data minimization**: Collect only necessary information
- **Access controls**: Restrict access to sensitive personal data
- **Retention policies**: Automatic deletion of old data
- **Anonymization**: Remove personal identifiers for analytics

#### Compliance Features
- **GDPR compliance**: Data subject rights and consent management
- **Local regulations**: Compliance with Malawi data protection laws
- **Audit reports**: Generate compliance reports for regulators
- **Data breach protocols**: Incident response and notification procedures

## 🛠️ Troubleshooting

### Common Issues

#### Form Submission Problems

**Error: "Required fields missing"**
- **Cause**: One or more mandatory fields not completed
- **Solution**: Check for red error indicators, complete all required fields
- **Prevention**: Use "Save Draft" feature to save progress

**Error: "Invalid file format"**
- **Cause**: Uploaded file not in supported format
- **Solution**: Convert file to PDF, DOC, DOCX, JPG, or PNG
- **Prevention**: Check file format before uploading

**Error: "File too large"**
- **Cause**: Uploaded file exceeds 10MB limit
- **Solution**: Compress file or split into smaller files
- **Prevention**: Check file size before uploading

#### Search and Filter Issues

**No search results found**
- **Check spelling**: Verify search terms are correct
- **Clear filters**: Remove all filters and try again
- **Broaden search**: Use partial terms or wildcards
- **Check permissions**: Ensure you have access to requested data

**Filters not working**
- **Clear browser cache**: Refresh page and clear cache
- **Reset filters**: Click "Clear All Filters" and start over
- **Check filter logic**: Ensure filter combinations are logical
- **Try one filter**: Apply filters one at a time to identify issues

#### Performance Issues

**Slow page loading**
- **Check internet connection**: Ensure stable connectivity
- **Clear browser cache**: Remove cached files
- **Close other applications**: Free up system resources
- **Try different browser**: Test with Chrome, Firefox, or Edge

**Export timeouts**
- **Reduce data size**: Apply filters to limit records
- **Choose different format**: CSV loads faster than Excel
- **Try smaller batches**: Export data in smaller chunks
- **Contact administrator**: For large dataset exports

### Error Messages

#### Common Error Codes

**ERR_001: Access Denied**
- **Meaning**: Insufficient permissions for requested action
- **Solution**: Contact administrator to request appropriate permissions
- **Check**: Verify your user role and assigned permissions

**ERR_002: Data Validation Failed**
- **Meaning**: Submitted data doesn't meet validation rules
- **Solution**: Review error messages and correct invalid fields
- **Check**: Ensure all required fields are completed correctly

**ERR_003: Duplicate Record**
- **Meaning**: SME with same name/registration already exists
- **Solution**: Search for existing record or use different registration number
- **Check**: Verify business registration number is unique

**ERR_004: Session Expired**
- **Meaning**: User session has timed out due to inactivity
- **Solution**: Log in again and resume work
- **Prevention**: Save work frequently, use "Keep me logged in" option

**ERR_005: File Upload Failed**
- **Meaning**: Document upload was unsuccessful
- **Solution**: Check file format, size, and internet connection
- **Retry**: Try uploading file again

### Getting Help

#### Self-Service Options
- **Contextual help**: Click "?" icons for field-specific help
- **Tooltips**: Hover over interface elements for quick tips
- **User manual**: Reference this comprehensive guide
- **FAQ section**: Check frequently asked questions

#### Contact Support
- **Local administrator**: First point of contact for user issues
- **IT support**: Technical problems and system errors
- **System administrator**: Permission and access issues
- **Emergency support**: Critical system problems (24/7)

#### Reporting Bugs
When reporting issues, include:
- **Error message**: Exact text of any error messages
- **Steps taken**: What you were doing when error occurred
- **Browser info**: Browser type and version
- **Screenshot**: Visual documentation of the problem
- **Expected behavior**: What should have happened

---

## 📞 Support and Resources

### Quick Reference

#### Essential Keyboard Shortcuts
- **Ctrl + N**: Create new SME record
- **Ctrl + S**: Save current form
- **Ctrl + F**: Search in current page
- **Ctrl + E**: Edit current record
- **Ctrl + P**: Print current page
- **Esc**: Cancel current operation

#### Contact Information
- **System Administrator**: admin@sme-database.gov.mw
- **Technical Support**: support@sme-database.gov.mw
- **Training Support**: training@sme-database.gov.mw
- **Emergency Support**: +265-XXX-XXXX (24/7)

### Additional Resources
- **Video tutorials**: Available in system help section
- **Training workshops**: Regular sessions for new users
- **User community**: Forum for sharing tips and best practices
- **System updates**: Regular feature enhancements and improvements

---

**Document Version**: 1.0  
**Last Updated**: October 2024  
**Next Review**: January 2025

**Previous**: [Main Manual](../README.md) ← | **Next**: [Lender Management](lender-management.md) →