# Widgets & Overlays - Reference Examples

Reusable widgets and overlay patterns used across custom pages.

## Events Calendar Widget

Source: `resources/js/components/widgets/EventsCalendarWidget.tsx`

### Interface

```typescript
export interface CalendarEvent {
  id: number;
  title: string;
  date: string;
  endDate?: string;
  venue: string;
  district?: string;
  isAttending?: boolean;
}

interface EventsCalendarWidgetProps {
  events: CalendarEvent[];
  district?: string;
  isLoading?: boolean;
  smeId?: number;
}
```

### Event Date Highlighting

Uses the shadcn Calendar `modifiers` API to highlight dates that have events:

```typescript
const eventDates = React.useMemo(() => {
  const dates: Date[] = [];
  events.forEach((event) => {
    const startDate = new Date(event.date);
    const endDate = event.endDate ? new Date(event.endDate) : undefined;

    if (endDate) {
      // Add all dates in range
      const current = new Date(startDate);
      while (current <= endDate) {
        dates.push(new Date(current));
        current.setDate(current.getDate() + 1);
      }
    } else {
      dates.push(startDate);
    }
  });
  return dates;
}, [events]);

// Calendar modifiers
const modifiers = { hasEvent: eventDates };
const modifiersClassNames = {
  hasEvent: "bg-primary/20 font-semibold text-primary hover:bg-primary/30",
};
```

### Calendar with Event List

```tsx
<Card className="flex flex-col">
  <CardHeader className="pb-2 px-4">
    <CardTitle className="flex items-center gap-2 text-base">
      <CalendarIcon className="h-4 w-4" />
      Events
    </CardTitle>
    <CardDescription className="text-xs">
      {district ? `${district} & nationwide` : "Nationwide"}
    </CardDescription>
  </CardHeader>
  <CardContent className="flex flex-col gap-3 px-4 pb-4">
    {/* Compact Calendar */}
    <div className="flex justify-center">
      <Calendar
        mode="single"
        selected={selectedDate}
        onSelect={setSelectedDate}
        modifiers={modifiers}
        modifiersClassNames={modifiersClassNames}
        className="rounded-md border p-2 [--cell-size:1.75rem] text-xs"
      />
    </div>

    {/* Events List Header */}
    <div className="flex items-center justify-between">
      <h4 className="text-xs font-medium text-muted-foreground">
        {selectedDate ? formatEventDate(selectedDate.toISOString()) : "Upcoming"}
      </h4>
      {selectedDate && (
        <button
          onClick={() => setSelectedDate(undefined)}
          className="text-[10px] text-muted-foreground hover:text-foreground"
        >
          Clear
        </button>
      )}
    </div>

    {/* Scrollable Events */}
    <ScrollArea className="h-[140px]">
      {displayEvents.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-4 text-center">
          <CalendarIcon className="h-6 w-6 text-muted-foreground/50" />
          <p className="mt-1 text-xs text-muted-foreground">
            {selectedDate ? "No events" : "No upcoming events"}
          </p>
        </div>
      ) : (
        <div className="space-y-1.5 pr-2">
          {displayEvents.map((event) => (
            <EventItem key={event.id} event={event} />
          ))}
        </div>
      )}
    </ScrollArea>
  </CardContent>
</Card>
```

### Event Item with Urgency Badge

```tsx
const daysUntil = getDaysUntil(event.date);
const urgency = getUrgencyBadge(daysUntil);

<div className="flex flex-col gap-1 rounded-md border p-2 transition-colors hover:bg-muted/50">
  <div className="flex items-start justify-between gap-1">
    <h5 className="font-medium text-xs leading-tight line-clamp-1">{event.title}</h5>
    <Badge variant={urgency.variant} className="shrink-0 text-[10px] px-1.5 py-0">
      {urgency.label}
    </Badge>
  </div>
  <div className="flex items-center gap-2 text-[10px] text-muted-foreground">
    <span className="flex items-center gap-0.5">
      <Clock className="h-2.5 w-2.5" />{formatEventTime(event.date)}
    </span>
    {event.venue && (
      <span className="flex items-center gap-0.5 truncate">
        <MapPin className="h-2.5 w-2.5 shrink-0" />
        <span className="truncate">{event.venue}</span>
      </span>
    )}
  </div>
  {smeId && (
    <Button
      size="sm"
      variant={isAttending ? "default" : "outline"}
      className="w-full mt-1 h-7 text-[10px]"
      onClick={(e) => { e.stopPropagation(); handleToggleAttendance(event.id); }}
    >
      {isAttending ? "Attending ✓" : "Attend"}
    </Button>
  )}
</div>
```

### Urgency Badge Logic

```typescript
function getUrgencyBadge(daysUntil: number) {
  if (daysUntil <= 0) return { variant: "destructive", label: "Today" };
  if (daysUntil === 1) return { variant: "destructive", label: "Tomorrow" };
  if (daysUntil <= 7) return { variant: "default", label: `${daysUntil} days` };
  return { variant: "secondary", label: `${daysUntil} days` };
}
```

### Optimistic Attendance Toggle

```typescript
const handleToggleAttendance = async (eventId: number) => {
  if (!smeId) return;
  const isCurrentlyAttending = attendingEvents.has(eventId);

  // Optimistic update
  setAttendingEvents((prev) => {
    const next = new Set(prev);
    if (isCurrentlyAttending) { next.delete(eventId); }
    else { next.add(eventId); }
    return next;
  });

  try {
    if (isCurrentlyAttending) {
      await axios.delete(`/api/portal/events/${eventId}/attend`);
    } else {
      await axios.post(`/api/portal/events/${eventId}/attend`);
    }
  } catch (error) {
    // Revert on error
    setAttendingEvents((prev) => {
      const next = new Set(prev);
      if (isCurrentlyAttending) { next.add(eventId); }
      else { next.delete(eventId); }
      return next;
    });
    console.error("Error toggling attendance:", error);
  }
};
```

---

## Notification Drawer (Doorbell)

Source: `resources/js/components/Notifications/NotificationDrawer.tsx`

### Interface

```typescript
interface Notification {
  id: number;
  title: string;
  message: string;
  type: string;          // "message", "mention", "system", "warning", "success"
  is_read: boolean;
  is_dismissed: boolean;
  priority: string;      // "high", "medium", "low"
  created_at: string;
  read_at?: string;
  trigger_user?: { id: number; name: string; email: string };
  related_type?: string; // "application", "event", "procurement"
  related_id?: number;
}
```

### Drawer Trigger with Badge Counter

```tsx
<Drawer open={isOpen} onOpenChange={setIsOpen}>
  <DrawerTrigger asChild>
    <Button variant="ghost" size="icon" className="relative">
      <Bell className="h-5 w-5" />
      {(() => {
        const pendingApps = canManageApplications ? (counts?.pending_applications ?? 0) : 0;
        const total = unreadCount + pendingApps;
        return total > 0 ? (
          <Badge
            variant="destructive"
            className="absolute -top-1 -right-1 h-5 min-w-5 text-xs px-1"
          >
            {total > 99 ? "99+" : total}
          </Badge>
        ) : null;
      })()}
    </Button>
  </DrawerTrigger>
  <DrawerContent className="max-w-md mx-auto">
    {/* ... */}
  </DrawerContent>
</Drawer>
```

### Drawer Header with Actions

```tsx
<DrawerHeader className="pb-4">
  <div className="flex items-center justify-between">
    <div>
      <DrawerTitle className="flex items-center gap-2">
        <BellRing className="h-5 w-5" />Notifications
      </DrawerTitle>
      <DrawerDescription>
        {unreadCount > 0 ? `${unreadCount} unread notifications` : "All caught up!"}
      </DrawerDescription>
    </div>
    <div className="flex items-center gap-2">
      {unreadCount > 0 && (
        <Button variant="outline" size="sm" onClick={handleMarkAllAsRead} disabled={loading}>
          <Check className="h-4 w-4 mr-1" />Mark all read
        </Button>
      )}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon"><MoreHorizontal className="h-4 w-4" /></Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={handleMarkAllAsRead}>Mark all as read</DropdownMenuItem>
          <DropdownMenuItem onClick={handleDismissAll}>Dismiss all</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </div>
</DrawerHeader>
```

### Action Alert Banner (Permission-Gated)

```tsx
{canManageApplications && counts?.pending_applications > 0 && (
  <div
    className="mb-4 p-3 rounded-lg border border-l-4 border-l-amber-500 bg-amber-50 dark:bg-amber-950/20 cursor-pointer hover:bg-amber-100 dark:hover:bg-amber-950/30 transition-colors"
    onClick={() => { setIsOpen(false); router.visit('/admin/applications'); }}
  >
    <div className="flex items-center gap-3">
      <FileText className="h-5 w-5 text-amber-600" />
      <div className="flex-1">
        <p className="text-sm font-medium text-amber-800 dark:text-amber-200">
          {counts.pending_applications} Pending Application{counts.pending_applications !== 1 ? 's' : ''}
        </p>
        <p className="text-xs text-amber-600 dark:text-amber-400">Click to review and process</p>
      </div>
      <Badge variant="outline" className="bg-amber-100 text-amber-800 border-amber-300">
        Action Required
      </Badge>
    </div>
  </div>
)}
```

### Tabbed Filtering

```tsx
<Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
  <TabsList className="grid w-full grid-cols-4">
    <TabsTrigger value="all" className="text-xs">
      All
      {counts?.total > 0 && (
        <Badge variant="secondary" className="ml-1 h-4 min-w-4 text-xs px-1">{counts.total}</Badge>
      )}
    </TabsTrigger>
    <TabsTrigger value="unread" className="text-xs">
      Unread
      {counts?.unread > 0 && (
        <Badge variant="destructive" className="ml-1 h-4 min-w-4 text-xs px-1">{counts.unread}</Badge>
      )}
    </TabsTrigger>
    <TabsTrigger value="messages" className="text-xs">Messages</TabsTrigger>
    <TabsTrigger value="system" className="text-xs">System</TabsTrigger>
  </TabsList>
  <TabsContent value={activeTab} className="mt-4">
    <ScrollArea className="h-[400px] pr-4">
      {/* Notification items */}
    </ScrollArea>
  </TabsContent>
</Tabs>
```

### Priority-Colored Notification Items

```typescript
const getPriorityColor = (priority: string) => {
  switch (priority) {
    case "high":   return "border-l-red-500 bg-red-50 dark:bg-red-950/20";
    case "medium": return "border-l-yellow-500 bg-yellow-50 dark:bg-yellow-950/20";
    case "low":    return "border-l-blue-500 bg-blue-50 dark:bg-blue-950/20";
    default:       return "border-l-gray-500 bg-gray-50 dark:bg-gray-950/20";
  }
};
```

```tsx
<div
  className={cn(
    "relative p-4 rounded-lg border border-l-4 cursor-pointer transition-colors hover:bg-muted/50",
    getPriorityColor(notification.priority),
    notification.is_read ? "opacity-70" : ""
  )}
  onClick={() => handleNotificationClick(notification)}
>
  {/* Icon + content + actions */}
</div>
```

### Smart Routing by Type

```typescript
const getNotificationRoute = (notification: Notification): string | null => {
  if (notification.related_type) {
    switch (notification.related_type) {
      case "application": return "/admin/applications";
      case "event":
      case "procurement": return "/opportunities";
    }
  }
  switch (notification.type) {
    case "application_approved":
    case "application_rejected": return "/admin/applications";
    case "event":
    case "procurement": return "/opportunities";
    case "message":
    case "mention": return "/portal";
    default: return null;
  }
};
```

---

## CrudPage Read-Only Integration

Source: `resources/js/pages/MyApplications/Index.tsx`

When a custom page needs a data table with detail views but no create/edit/delete:

```tsx
export default function MyApplicationsIndex({
  data, filters, permissions, meta,
  userName, smeName, usmeNumber, classification, error,
}: MyApplicationsIndexProps) {
  const isMobile = useIsMobile();
  const handleRefresh = () => router.reload({ only: ['data'] });

  // Filter out admin-specific columns
  const myApplicationColumns = applicationColumns.filter(col =>
    !['registrant_name', 'email', 'phone'].includes(col.key as string)
  );

  return (
    <Admin title="My Applications">
      <Head title="My Applications" />
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Custom header with SME context badges */}
        <div className="px-4 lg:px-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight">My Applications</h1>
            <p className="text-muted-foreground">
              <span className="hidden sm:inline">Track the status of your applications and formalisation requests</span>
              <span className="sm:hidden">Track your applications</span>
            </p>
          </div>
          {smeName && (
            <div className="flex flex-col items-end gap-2 shrink-0">
              <div className="text-lg font-semibold text-right">{smeName}</div>
              <div className="flex items-center gap-2">
                {usmeNumber && <Badge variant="outline" className="text-xs font-mono">{usmeNumber}</Badge>}
                {classification && (
                  <Badge variant="default" className={`text-sm px-3 py-1 font-semibold ${
                    classification === 'Micro' ? 'bg-blue-600 hover:bg-blue-700'
                    : classification === 'Small' ? 'bg-emerald-600 hover:bg-emerald-700'
                    : classification === 'Medium' ? 'bg-amber-600 hover:bg-amber-700'
                    : ''
                  }`}>
                    {classification} Enterprise
                  </Badge>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Error state */}
        {error && (
          <div className="px-4 lg:px-6">
            <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-md">{error}</div>
          </div>
        )}

        {/* CrudPage in read-only mode */}
        <div className="px-0">
          <CrudPage<Application>
            data={data}
            filters={filters}
            title="Applications"
            resourceName="my-applications"
            columns={isMobile ? myApplicationColumnsMobile : myApplicationColumns}
            paginationConfig={meta?.pagination}
            detailView={ApplicationDetailView}
            onRefresh={handleRefresh}
            canView={true}
            readOnly={true}
          />
        </div>
      </div>
    </Admin>
  );
}
```

Key points:
- **`readOnly={true}`** hides create/edit/delete buttons
- **`canView={true}`** enables detail view on row click
- **Custom header** above the CrudPage for page-specific context
- **Column filtering** removes admin-specific columns for portal users
- **`router.reload({ only: ['data'] })`** for efficient refresh without full page reload

### Go Backend: GenericPageController with Mandatory Filter

For scoping CrudPage data to the current user's entity, use `SetMandatoryFilterProvider`:

```go
func NewMyApplicationsPageController() *GenericPageController[models.Application] {
    controller := NewGenericPageController[models.Application](/* ... */)

    controller.SetMandatoryFilterProvider(func(ctx http.Context) func(query contractsorm.Query) contractsorm.Query {
        user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
        if user == nil {
            return func(q contractsorm.Query) contractsorm.Query {
                return q.Where("1 = 0") // Return nothing
            }
        }
        sme, _ := services.NewSmeService().GetSmeByUserEmail(user.Email)
        if sme == nil {
            return func(q contractsorm.Query) contractsorm.Query {
                return q.Where("1 = 0")
            }
        }
        return func(q contractsorm.Query) contractsorm.Query {
            return q.Where("sme_id = ?", sme.ID)
        }
    })

    return controller
}
```

## Source Files

- Calendar widget: `resources/js/components/widgets/EventsCalendarWidget.tsx`
- Notification drawer: `resources/js/components/Notifications/NotificationDrawer.tsx`
- My Applications page: `resources/js/pages/MyApplications/Index.tsx`
- Notification context: `resources/js/contexts/NotificationContext.tsx`
- Permissions context: `resources/js/contexts/PermissionsContext.tsx`
