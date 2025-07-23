# Scoped Permissions UI Design

## Overview
The UI needs to clearly show the three-level permission hierarchy (by_me, by_my_role, by_all) in an intuitive way that's easy to understand and manage.

## Design Approach

### 1. Permission Matrix View (Current Enhancement)

Instead of simple checkboxes, each permission cell becomes a **dropdown selector** with scope options:

```
Service: Books
┌─────────────┬────────────┬────────────┬────────────┬────────────┐
│             │   Create   │    Read    │   Update   │   Delete   │
├─────────────┼────────────┼────────────┼────────────┼────────────┤
│ Admin Role  │ [By All ▼] │ [By All ▼] │ [By Role▼] │ [By Me  ▼] │
│ Editor Role │ [By Role▼] │ [By All ▼] │ [By Me  ▼] │ [By Me  ▼] │
│ Author Role │ [By Me  ▼] │ [By Role▼] │ [By Me  ▼] │ [None   ▼] │
└─────────────┴────────────┴────────────┴────────────┴────────────┘
```

Each dropdown contains:
- 🌍 By All (Can access all resources)
- 👥 By My Role (Can access resources created by same/lower roles)
- 👤 By Me (Can only access own resources)
- ❌ None (No access)

### 2. Visual Indicators

Use colors and icons to make scopes immediately recognizable:

```tsx
const scopeConfig = {
  by_all: {
    label: "By All",
    icon: "🌍",
    color: "green",
    description: "Can access all resources"
  },
  by_my_role: {
    label: "By My Role", 
    icon: "👥",
    color: "blue",
    description: "Can access resources created by same or lower role levels"
  },
  by_me: {
    label: "By Me",
    icon: "👤", 
    color: "orange",
    description: "Can only access own resources"
  },
  none: {
    label: "None",
    icon: "❌",
    color: "gray",
    description: "No access"
  }
}
```

### 3. Enhanced Permission Component

```tsx
// PermissionScopeSelector.tsx
interface PermissionScopeSelectorProps {
  service: string;
  action: string;
  currentScope?: string;
  onChange: (scope: string) => void;
}

export function PermissionScopeSelector({ 
  service, 
  action, 
  currentScope = 'none',
  onChange 
}: PermissionScopeSelectorProps) {
  return (
    <Select value={currentScope} onValueChange={onChange}>
      <SelectTrigger className="w-[140px]">
        <SelectValue>
          <div className="flex items-center gap-2">
            <span>{scopeConfig[currentScope].icon}</span>
            <span>{scopeConfig[currentScope].label}</span>
          </div>
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {Object.entries(scopeConfig).map(([scope, config]) => (
          <SelectItem key={scope} value={scope}>
            <div className="flex items-center gap-2">
              <span>{config.icon}</span>
              <span>{config.label}</span>
            </div>
            <span className="text-xs text-muted-foreground block ml-6">
              {config.description}
            </span>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
```

### 4. Quick Presets

Add preset buttons for common permission patterns:

```tsx
<div className="flex gap-2 mb-4">
  <Button variant="outline" size="sm" onClick={() => applyPreset('admin')}>
    <Shield className="w-4 h-4 mr-2" />
    Admin Preset
  </Button>
  <Button variant="outline" size="sm" onClick={() => applyPreset('editor')}>
    <Edit className="w-4 h-4 mr-2" />
    Editor Preset
  </Button>
  <Button variant="outline" size="sm" onClick={() => applyPreset('viewer')}>
    <Eye className="w-4 h-4 mr-2" />
    Viewer Preset
  </Button>
</div>
```

Presets would set:
- **Admin**: All permissions with "by_all" scope
- **Editor**: Create/Update "by_my_role", Read "by_all", Delete "by_me"
- **Viewer**: Read "by_all", everything else "none"

### 5. Bulk Actions

For managing multiple permissions at once:

```tsx
<div className="border rounded p-4 mb-4">
  <h4 className="font-medium mb-2">Bulk Actions</h4>
  <div className="flex gap-4">
    <Select>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select services..." />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Services</SelectItem>
        <SelectItem value="books">Books</SelectItem>
        <SelectItem value="users">Users</SelectItem>
        <SelectItem value="roles">Roles</SelectItem>
      </SelectContent>
    </Select>
    
    <Select>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select actions..." />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Actions</SelectItem>
        <SelectItem value="read">Read</SelectItem>
        <SelectItem value="create">Create</SelectItem>
        <SelectItem value="update">Update</SelectItem>
        <SelectItem value="delete">Delete</SelectItem>
      </SelectContent>
    </Select>
    
    <Select>
      <SelectTrigger className="w-[140px]">
        <SelectValue placeholder="Set scope..." />
      </SelectTrigger>
      <SelectContent>
        {/* Scope options */}
      </SelectContent>
    </Select>
    
    <Button>Apply</Button>
  </div>
</div>
```

### 6. Permission Summary View

Show a readable summary of what each role can do:

```tsx
<Card>
  <CardHeader>
    <CardTitle>Editor Role Permissions</CardTitle>
  </CardHeader>
  <CardContent>
    <div className="space-y-4">
      <div>
        <h4 className="font-medium flex items-center gap-2">
          <Book className="w-4 h-4" />
          Books
        </h4>
        <ul className="ml-6 mt-2 space-y-1 text-sm">
          <li className="flex items-center gap-2">
            <span className="text-green-600">✓</span>
            Can create books (owned by them)
          </li>
          <li className="flex items-center gap-2">
            <span className="text-green-600">✓</span>
            Can read all books
          </li>
          <li className="flex items-center gap-2">
            <span className="text-blue-600">✓</span>
            Can update books created by editors or below
          </li>
          <li className="flex items-center gap-2">
            <span className="text-orange-600">✓</span>
            Can only delete their own books
          </li>
        </ul>
      </div>
    </div>
  </CardContent>
</Card>
```

### 7. Testing Interface

Add a permission tester to verify configurations:

```tsx
<Card>
  <CardHeader>
    <CardTitle>Test Permissions</CardTitle>
  </CardHeader>
  <CardContent>
    <div className="space-y-4">
      <Select>
        <SelectTrigger>
          <SelectValue placeholder="Select a user..." />
        </SelectTrigger>
        {/* User options */}
      </Select>
      
      <Select>
        <SelectTrigger>
          <SelectValue placeholder="Select a resource..." />
        </SelectTrigger>
        {/* Resource options with creator info */}
      </Select>
      
      <Select>
        <SelectTrigger>
          <SelectValue placeholder="Select an action..." />
        </SelectTrigger>
        {/* Action options */}
      </Select>
      
      <Button onClick={testPermission}>Test Access</Button>
      
      {result && (
        <Alert>
          <AlertTitle>{result.allowed ? "✅ Allowed" : "❌ Denied"}</AlertTitle>
          <AlertDescription>
            {result.reason}
          </AlertDescription>
        </Alert>
      )}
    </div>
  </CardContent>
</Card>
```

## Implementation Components

### 1. Enhanced RolePermissions Page

Update the existing RolePermissions component to use the scope selector:

```tsx
// pages/Permissions/RolePermissions.tsx
const [permissions, setPermissions] = useState<Record<string, string>>({});

// Initialize with current permissions
useEffect(() => {
  const perms: Record<string, string> = {};
  
  // For each service/action combination
  services.forEach(service => {
    actions.forEach(action => {
      // Check which scope the role has
      const byAllPerm = `${service.slug}_${action.slug}_by_all`;
      const byRolePerm = `${service.slug}_${action.slug}_by_my_role`;
      const byMePerm = `${service.slug}_${action.slug}_by_me`;
      
      if (currentPermissions[byAllPerm]) {
        perms[`${service.slug}_${action.slug}`] = 'by_all';
      } else if (currentPermissions[byRolePerm]) {
        perms[`${service.slug}_${action.slug}`] = 'by_my_role';
      } else if (currentPermissions[byMePerm]) {
        perms[`${service.slug}_${action.slug}`] = 'by_me';
      } else {
        perms[`${service.slug}_${action.slug}`] = 'none';
      }
    });
  });
  
  setPermissions(perms);
}, [services, actions, currentPermissions]);
```

### 2. Save Mechanism

When saving, convert the scope selections back to individual permissions:

```tsx
const handleSave = async () => {
  const permissionsToGrant: string[] = [];
  const permissionsToRevoke: string[] = [];
  
  Object.entries(permissions).forEach(([key, scope]) => {
    const [service, action] = key.split('_');
    
    // Remove all scoped versions first
    ['by_all', 'by_my_role', 'by_me'].forEach(s => {
      const perm = `${service}_${action}_${s}`;
      if (currentPermissions[perm] && s !== scope) {
        permissionsToRevoke.push(perm);
      }
    });
    
    // Add the selected scope
    if (scope !== 'none') {
      const perm = `${service}_${action}_${scope}`;
      if (!currentPermissions[perm]) {
        permissionsToGrant.push(perm);
      }
    }
  });
  
  // Send to API
  await updateRolePermissions(role.id, {
    grant: permissionsToGrant,
    revoke: permissionsToRevoke
  });
};
```

## Benefits

1. **Intuitive**: Users can immediately understand the hierarchy
2. **Visual**: Colors and icons make different scopes easy to distinguish
3. **Flexible**: Allows fine-grained control while keeping the UI simple
4. **Efficient**: Bulk actions and presets speed up configuration
5. **Testable**: Built-in testing helps verify configurations

This design maintains the familiar permission matrix layout while adding the power of scoped permissions in a user-friendly way.