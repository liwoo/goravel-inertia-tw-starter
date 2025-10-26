# 4. Roles and Permissions

## Overview

The SMEDI SME Database uses a comprehensive Role-Based Access Control (RBAC) system to ensure users only access data and features appropriate to their responsibilities. This section explains the different user roles, their permissions, and how the system manages access control.

## 👥 User Roles Hierarchy

### System Roles (Highest to Lowest Authority)

```
Super Administrator
    ├── Administrator
    ├── Librarian
    ├── Moderator
    ├── Member
    └── Guest
```

## Role Descriptions

### 1. Super Administrator
**Level**: 100 (Highest Authority)
**Description**: Complete system access with all permissions

**Capabilities**:
- ✅ **Full System Control**: Configure all system settings
- ✅ **User Management**: Create, modify, and delete all user accounts
- ✅ **Role Assignment**: Assign and modify any user role
- ✅ **Data Management**: Full CRUD access to all data
- ✅ **System Administration**: Backup, restore, and system maintenance
- ✅ **Audit Access**: View all system logs and audit trails
- ✅ **Security Management**: Configure security settings and policies

**Typical Users**: System developers, IT administrators, senior government officials

### 2. Administrator
**Level**: 80 (High Authority)
**Description**: Administrative access to most system features

**Capabilities**:
- ✅ **User Management**: Create and manage user accounts (except Super Admins)
- ✅ **Data Oversight**: Full access to SME and lender data
- ✅ **Report Generation**: Create and export all types of reports
- ✅ **Configuration**: Modify system configurations (limited)
- ✅ **Monitoring**: View system performance and user activities
- ❌ **System Settings**: Cannot modify core system settings
- ❌ **Super Admin Management**: Cannot manage Super Administrator accounts

**Typical Users**: Ministry department heads, senior data managers, IT supervisors

### 3. Librarian
**Level**: 60 (Medium-High Authority)
**Description**: Full access to data management with book/record focus

**Capabilities**:
- ✅ **SME Data Management**: Full CRUD operations on SME records
- ✅ **Lender Management**: Full CRUD operations on lender data
- ✅ **Data Import/Export**: Bulk data operations and file management
- ✅ **Standard Reports**: Generate most reports and analytics
- ✅ **Data Validation**: Review and approve data entries
- ❌ **User Management**: Cannot create or modify user accounts
- ❌ **System Configuration**: Cannot modify system settings

**Typical Users**: Data managers, records officers, SME specialists

### 4. Moderator
**Level**: 40 (Medium Authority)
**Description**: Limited administrative access with oversight capabilities

**Capabilities**:
- ✅ **Data Review**: Review and approve pending data entries
- ✅ **Basic Reports**: Generate standard reports within their scope
- ✅ **Data Editing**: Edit SME and lender records (with restrictions)
- ✅ **Quality Control**: Flag data inconsistencies and errors
- ❌ **Data Deletion**: Cannot delete records permanently
- ❌ **Bulk Operations**: Limited bulk data operations
- ❌ **User Management**: Cannot manage user accounts

**Typical Users**: Supervisors, quality control officers, regional coordinators

### 5. Member
**Level**: 20 (Basic Authority)
**Description**: Standard user with data entry and viewing privileges

**Capabilities**:
- ✅ **Data Entry**: Create new SME and lender records
- ✅ **Data Viewing**: View assigned data within their scope
- ✅ **Basic Editing**: Edit records they created (with approval workflow)
- ✅ **Personal Reports**: Generate reports for their own data
- ❌ **Data Deletion**: Cannot delete any records
- ❌ **Bulk Operations**: Cannot perform bulk operations
- ❌ **Administrative Functions**: No administrative access

**Typical Users**: Data entry clerks, field officers, junior staff

### 6. Guest
**Level**: 10 (Lowest Authority)
**Description**: Read-only access to limited public information

**Capabilities**:
- ✅ **Public Data Viewing**: View non-sensitive public statistics
- ✅ **Basic Reports**: Access to pre-approved public reports
- ✅ **Search**: Search public SME directory (limited fields)
- ❌ **Data Modification**: Cannot create, edit, or delete data
- ❌ **Sensitive Data**: Cannot access financial or personal information
- ❌ **Detailed Reports**: Cannot generate detailed analytics

**Typical Users**: External researchers, public users, temporary visitors

## 🔐 Permission Categories

### Data Access Permissions

#### SME Management
| Permission | Super Admin | Admin | Librarian | Moderator | Member | Guest |
|------------|-------------|--------|-----------|-----------|---------|-------|
| Create SME Records | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| View SME Records | ✅ | ✅ | ✅ | ✅ | Limited | Public Only |
| Edit SME Records | ✅ | ✅ | ✅ | Limited | Own Only | ❌ |
| Delete SME Records | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| Bulk SME Operations | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| Export SME Data | ✅ | ✅ | ✅ | Limited | Limited | ❌ |

#### Lender Management
| Permission | Super Admin | Admin | Librarian | Moderator | Member | Guest |
|------------|-------------|--------|-----------|-----------|---------|-------|
| Create Lenders | ✅ | ✅ | ✅ | ✅ | Limited | ❌ |
| View Lenders | ✅ | ✅ | ✅ | ✅ | Limited | ❌ |
| Edit Lenders | ✅ | ✅ | ✅ | Limited | ❌ | ❌ |
| Delete Lenders | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |

#### User Management
| Permission | Super Admin | Admin | Librarian | Moderator | Member | Guest |
|------------|-------------|--------|-----------|-----------|---------|-------|
| Create Users | ✅ | Limited | ❌ | ❌ | ❌ | ❌ |
| View Users | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Edit Users | ✅ | Limited | ❌ | ❌ | ❌ | ❌ |
| Delete Users | ✅ | Limited | ❌ | ❌ | ❌ | ❌ |
| Assign Roles | ✅ | Limited | ❌ | ❌ | ❌ | ❌ |

## 🎯 Data Scope and Filtering

### Geographic Scope
Users may be restricted to specific geographic regions:
- **National Level**: Full country access
- **Regional Level**: Specific regions only
- **District Level**: Specific districts only
- **Local Level**: Specific areas or communities

### Organizational Scope
Access may be limited by organization type:
- **Government Entities**: Ministry departments and agencies
- **Financial Institutions**: Banks and microfinance organizations
- **Development Partners**: NGOs and international organizations
- **Private Sector**: Consulting firms and service providers

### Data Sensitivity Levels
Information is classified by sensitivity:
- **Public**: General business information, public statistics
- **Internal**: Operational data, performance metrics
- **Confidential**: Financial details, personal information
- **Restricted**: Highly sensitive commercial or personal data

## 📊 Permission Management

### How Permissions Work

#### Automatic Assignment
- **Role-Based**: Permissions are automatically assigned based on user role
- **Inheritance**: Lower roles inherit appropriate permissions from higher roles
- **Scope-Based**: Geographic and organizational restrictions are applied automatically

#### Permission Checking
The system checks permissions at multiple levels:
1. **Page Level**: Determines if user can access specific pages
2. **Feature Level**: Controls access to specific features within pages
3. **Data Level**: Filters data based on user's scope and role
4. **Action Level**: Validates permissions before allowing specific actions

#### Real-Time Validation
- **Dynamic Menus**: Only show menu items user has permission to access
- **Button States**: Disable buttons for actions user cannot perform
- **Data Filtering**: Automatically filter displayed data based on user scope
- **Form Fields**: Hide or disable fields user cannot modify

### Permission Inheritance

```
Super Administrator (All Permissions)
    ├── Administrator (Most Permissions)
    │   ├── Librarian (Data Management Permissions)
    │   │   ├── Moderator (Review Permissions)
    │   │   │   ├── Member (Basic Permissions)
    │   │   │   │   └── Guest (View-Only Permissions)
```

## 🔧 Managing User Roles

### Assigning Roles (Admin Only)

#### Step 1: Access User Management
1. Navigate to **Administration** → **User Management**
2. Click on the user you want to modify
3. Select the **Roles** tab

#### Step 2: Role Selection
1. Choose the appropriate role from the dropdown
2. Consider the user's responsibilities and data access needs
3. Review the permission summary displayed
4. Click **"Assign Role"**

#### Step 3: Scope Configuration
1. If applicable, set geographic scope restrictions
2. Configure organizational access limitations
3. Set data sensitivity level access
4. Save the configuration

### Modifying Permissions (Super Admin Only)

#### Custom Permission Sets
1. Navigate to **Administration** → **Permissions**
2. Create custom permission groups if needed
3. Assign specific permissions to roles
4. Test permissions with test accounts

#### Permission Auditing
1. Regular review of assigned permissions
2. Remove unnecessary access rights
3. Ensure compliance with data protection policies
4. Document permission changes

## ⚠️ Security Considerations

### Principle of Least Privilege
- Users receive **minimum necessary permissions** for their role
- **Regular reviews** ensure permissions remain appropriate
- **Temporary access** for specific projects with automatic expiration
- **Escalation procedures** for requesting additional permissions

### Access Monitoring
- **Audit Logs**: All permission changes are logged
- **Usage Tracking**: Monitor which permissions are actually used
- **Anomaly Detection**: Alert on unusual access patterns
- **Regular Reviews**: Quarterly permission audits

### Best Practices
1. **Role Assignment**: Assign the lowest appropriate role first
2. **Scope Limitation**: Restrict geographic and organizational scope when possible
3. **Regular Reviews**: Conduct quarterly access reviews
4. **Documentation**: Document reasons for special permissions
5. **Training**: Ensure users understand their permission limitations

## 🚨 Common Permission Issues

### Access Denied Errors
**Problem**: User cannot access expected features
**Causes**:
- Role assignment is too restrictive
- Geographic scope limitations
- Data sensitivity restrictions
- Expired account or permissions

**Solutions**:
1. Check user's current role assignment
2. Verify geographic and organizational scope
3. Review data sensitivity access levels
4. Contact administrator for role elevation if justified

### Missing Data or Features
**Problem**: Expected data or menu items are not visible
**Causes**:
- Insufficient permissions for data sensitivity level
- Geographic scope restrictions
- Organizational access limitations
- Feature requires higher role level

**Solutions**:
1. Request appropriate role assignment from administrator
2. Verify data falls within user's scope
3. Check if feature requires specific permissions
4. Contact administrator for scope expansion if needed

## 📋 Permission Quick Reference

### Common Tasks and Required Roles

| Task | Minimum Role Required | Notes |
|------|----------------------|-------|
| View SME basic info | Guest | Public information only |
| Add new SME record | Member | Within assigned scope |
| Edit SME financial data | Librarian | Sensitive financial information |
| Delete SME record | Librarian | With appropriate justification |
| Generate reports | Member | Limited to own data scope |
| Export data | Moderator | Bulk operations restricted |
| Manage users | Administrator | Cannot manage Super Admins |
| Configure system | Super Administrator | System-wide settings |
| View audit logs | Administrator | Security monitoring |
| Assign roles | Administrator | Limited to lower roles |

### Emergency Access Procedures
1. **Contact Super Administrator** for urgent access needs
2. **Temporary role elevation** available for critical situations
3. **Emergency accounts** maintained for system recovery
4. **24/7 support** available for critical access issues

---

**Previous**: [Login and Authentication](03-login-authentication.md) ← | **Next**: [Navigation and Interface](05-navigation-interface.md) →
