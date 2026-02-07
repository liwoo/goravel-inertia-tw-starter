---
name: custom-pages
description: Guide for building non-CRUD pages with consistent design. Use when creating portal pages, custom views, detail modals, sidebar widgets, notification drawers, or any page that is not a standard CRUD table. Covers fullscreen modals, minimal text patterns, calendar widgets, doorbell notifications, optimistic updates, and CrudPage read-only integration.
---

# Custom Pages Pattern

This skill defines how to build pages that are not standard CRUD tables. Follow these patterns for portal pages, detail views, dashboard-style layouts, and interactive features to ensure visual consistency across the application.

## Design Principles

1. **Minimal text, maximum signal** - Use icons, badges, and compact layouts instead of paragraphs. Show data, not descriptions.
2. **Fullscreen modals for detail views** - Open items in near-fullscreen `Dialog` (98vw x 96vh) with sidebar metadata + main content. Never navigate away from the list.
3. **Progressive disclosure** - Summary cards on the page, full details in a modal/sheet. Tabs for complex data.
4. **Icon-driven metadata** - Small icons (`h-3 w-3` to `h-4 w-4`) next to compact key-value pairs. Never bare text.
5. **Responsive descriptions** - Long text on desktop, short on mobile: `hidden md:inline` / `md:hidden`.
6. **Sidebar widgets** - Calendar, events, activities as right-sidebar cards. Full width on mobile, fixed `xl:w-[380px]` on desktop.
7. **Optimistic updates** - Update UI state immediately, make API call, revert on error.
8. **Permission-gated** - Wrap sections in `PermissionGate` or conditionally render based on `canPerformAction()`.

## Page Types

| Type | Pattern | Example |
|------|---------|---------|
| **Fully custom** | Custom controller + custom React page | Opportunities (procurement for SME users) |
| **CrudPage read-only** | GenericPageController with mandatory filters + CrudPage `readOnly` | My Applications |
| **Role-specific dashboard** | Detect role in dashboard, render alternate component | SME Dashboard |
| **Sheet/Drawer overlay** | Slide-in panel for editing without leaving page | MySmeModal, FormalisationEditSheet |
| **Drawer notification** | Bottom-slide drawer for notifications | NotificationDrawer (doorbell) |

## File Structure

```
# Fully custom page
app/http/controllers/{feature}/
  {feature}_page_controller.go        # Custom controller (NOT GenericPageController)
resources/js/pages/{Feature}/
  Index.tsx                            # Main page component

# CrudPage read-only
app/http/controllers/{feature}/
  {feature}_page_controller.go         # GenericPageController with SetMandatoryFilterProvider
resources/js/pages/{Feature}/
  Index.tsx                            # Wraps CrudPage with readOnly={true}

# Widgets
resources/js/components/widgets/
  {Widget}Widget.tsx                   # Reusable widget card
  index.ts                             # Barrel exports

# Overlays
resources/js/components/
  {Feature}Modal.tsx                   # Sheet or Dialog overlay
```

## Page Layout Template

Every custom page follows this layout:

```tsx
<Admin title="Page Title">
  <Head title="Page Title" />
  <div className="p-4 md:p-6 flex flex-col gap-6">
    {/* 1. Header */}
    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-semibold tracking-tight flex items-center gap-2">
          <Icon className="h-6 w-6 text-amber-500" />
          Page Title
        </h1>
        <p className="text-muted-foreground">
          <span className="hidden md:inline">Full description for desktop</span>
          <span className="md:hidden">Short description for mobile</span>
        </p>
      </div>
      {/* Optional: Context badges (SME name, score, etc.) */}
      {smeName && (
        <div className="flex flex-col items-end gap-2 shrink-0">
          <div className="text-lg font-semibold text-right">{smeName}</div>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="text-xs font-mono">{usmeNumber}</Badge>
            <Badge variant="default" className="bg-primary">Score: {score}%</Badge>
          </div>
        </div>
      )}
    </div>

    {/* 2. Error state */}
    {error && (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    )}

    {/* 3. Main + Sidebar */}
    <div className="flex flex-col xl:flex-row gap-6 items-start">
      <div className="flex-1 flex flex-col gap-6 min-w-0">
        {/* Stats row */}
        {/* Tabbed content */}
        {/* Cards / lists */}
        {/* Tip card */}
      </div>
      <aside className="w-full xl:w-[380px] 2xl:w-[420px] shrink-0">
        {/* Widgets */}
      </aside>
    </div>
  </div>

  {/* 4. Detail modal (rendered outside flow) */}
  <DetailModal item={selected} isOpen={isOpen} onClose={handleClose} />
</Admin>
```

## Fullscreen Detail Modal

Use for viewing item details. Two-column layout: metadata sidebar + main content area.

```tsx
const DetailModal: React.FC<{
  item: Item | null;
  isOpen: boolean;
  onClose: () => void;
}> = ({ item, isOpen, onClose }) => {
  if (!item) return null;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent
        showCloseButton={false}
        className="!max-w-[98vw] !w-[98vw] !max-h-[96vh] !h-[96vh] p-0 flex flex-col gap-0"
      >
        {/* Header bar */}
        <div className="flex items-center justify-between px-4 py-3 border-b shrink-0 bg-background">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8">
              <X className="h-4 w-4" />
            </Button>
            <div>
              <DialogTitle className="text-lg font-semibold">{item.title}</DialogTitle>
              <p className="text-xs text-muted-foreground">{item.subtitle}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="text-xs font-mono">{item.refNo}</Badge>
            <Badge className={statusColor}>{statusText}</Badge>
          </div>
        </div>

        {/* Two-column body */}
        <div className="flex-1 flex overflow-hidden">
          {/* Left sidebar - metadata */}
          <aside className="w-72 border-r bg-muted/30 flex flex-col shrink-0">
            <ScrollArea className="flex-1">
              <div className="p-4 space-y-6">
                <MetadataSection title="Organization">
                  <MetadataItem icon={<Building2 />} label="Name" value={item.org} />
                  <MetadataItem icon={<FileText />} label="Type" value={item.type} />
                </MetadataSection>
                <Separator />
                <MetadataSection title="Timeline">
                  <MetadataItem icon={<CalendarDays />} label="Opens" value={formatDate(item.openDate)} />
                  <MetadataItem icon={<Clock />} label="Closes" value={formatDate(item.closeDate)} />
                </MetadataSection>
                <Separator />
                {/* Action buttons in sidebar */}
                <Button onClick={handleAction} className="w-full" size="sm">
                  <Heart className="mr-2 h-3.5 w-3.5" /> Take Action
                </Button>
              </div>
            </ScrollArea>
          </aside>

          {/* Main content - document viewer */}
          <main className="flex-1 bg-muted/50 overflow-hidden flex items-start justify-center p-6">
            <div className="bg-background rounded-lg shadow-lg border overflow-hidden flex flex-col"
              style={{ width: 'min(100%, 794px)', height: 'calc(96vh - 80px)', maxHeight: '1123px' }}>
              {/* Document header */}
              <div className="px-8 py-6 border-b bg-gradient-to-r from-primary/5 to-transparent">
                <h1 className="text-2xl font-bold">{item.title}</h1>
              </div>
              {/* Document body */}
              <ScrollArea className="flex-1">
                <div className="px-8 py-6 space-y-8">
                  <section>
                    <h2 className="text-lg font-semibold flex items-center gap-2">
                      <FileText className="h-5 w-5 text-primary" /> Details
                    </h2>
                    <div className="prose prose-sm max-w-none dark:prose-invert">
                      <MarkdownEditor.Markdown source={item.details} />
                    </div>
                  </section>
                </div>
              </ScrollArea>
              {/* Document footer */}
              <div className="px-8 py-3 border-t bg-muted/30 text-xs text-muted-foreground flex justify-between">
                <span>Reference: {item.refNo}</span>
              </div>
            </div>
          </main>
        </div>
      </DialogContent>
    </Dialog>
  );
};
```

## MetadataItem Component

Reusable sidebar metadata display:

```tsx
const MetadataItem: React.FC<{
  icon: React.ReactNode;
  label: string;
  value: React.ReactNode;
}> = ({ icon, label, value }) => (
  <div className="flex items-start gap-2">
    <div className="p-1.5 rounded bg-muted shrink-0">{icon}</div>
    <div className="min-w-0 flex-1">
      <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-medium">{label}</p>
      <p className="text-sm font-medium text-foreground truncate">{value}</p>
    </div>
  </div>
);
```

## Sheet (Side Panel) for Editing

Use `Sheet` for editing forms that should overlay without navigating away:

```tsx
<Sheet open={open} onOpenChange={setOpen}>
  <SheetContent className="sm:max-w-lg">
    <SheetHeader>
      <SheetTitle>Edit Details</SheetTitle>
    </SheetHeader>
    <Tabs value={activeTab} onValueChange={setActiveTab}>
      <TabsList className="grid w-full grid-cols-2">
        <TabsTrigger value="business">Business</TabsTrigger>
        <TabsTrigger value="owner">Owner</TabsTrigger>
      </TabsList>
      <TabsContent value="business">
        <Card><CardContent className="space-y-4 pt-4">
          {/* Form fields */}
        </CardContent></Card>
      </TabsContent>
    </Tabs>
    <SheetFooter className="mt-4">
      <Button variant="outline" onClick={() => setOpen(false)}>Cancel</Button>
      <Button onClick={handleSave} disabled={isSaving}>
        {isSaving ? <><Loader2 className="mr-2 h-4 w-4 animate-spin" />Saving...</> : 'Save'}
      </Button>
    </SheetFooter>
  </SheetContent>
</Sheet>
```

## Notification Drawer (Doorbell)

Bottom-slide `Drawer` for notifications. Triggered from the site header bell icon.

```tsx
<Drawer open={isOpen} onOpenChange={setIsOpen}>
  <DrawerTrigger asChild>
    <Button variant="ghost" size="icon" className="relative">
      <Bell className="h-5 w-5" />
      {totalCount > 0 && (
        <Badge variant="destructive" className="absolute -top-1 -right-1 h-5 min-w-5 text-xs px-1">
          {totalCount > 99 ? "99+" : totalCount}
        </Badge>
      )}
    </Button>
  </DrawerTrigger>
  <DrawerContent className="max-w-md mx-auto">
    <DrawerHeader className="pb-4">
      <DrawerTitle className="flex items-center gap-2">
        <BellRing className="h-5 w-5" /> Notifications
      </DrawerTitle>
      <DrawerDescription>
        {unreadCount > 0 ? `${unreadCount} unread` : "All caught up!"}
      </DrawerDescription>
    </DrawerHeader>

    {/* Optional: Action-required alert */}
    {pendingCount > 0 && (
      <div className="mx-4 mb-4 p-3 rounded-lg border border-l-4 border-l-amber-500 bg-amber-50 dark:bg-amber-950/20 cursor-pointer"
        onClick={() => router.visit('/admin/applications')}>
        <div className="flex items-center gap-3">
          <FileText className="h-5 w-5 text-amber-600" />
          <div className="flex-1">
            <p className="text-sm font-medium">{pendingCount} Pending</p>
            <p className="text-xs text-amber-600">Click to review</p>
          </div>
          <Badge variant="outline">Action Required</Badge>
        </div>
      </div>
    )}

    {/* Tabbed notification list */}
    <Tabs value={activeTab} onValueChange={setActiveTab} className="px-4">
      <TabsList className="grid w-full grid-cols-4">
        <TabsTrigger value="all" className="text-xs">All</TabsTrigger>
        <TabsTrigger value="unread" className="text-xs">Unread</TabsTrigger>
        <TabsTrigger value="messages" className="text-xs">Messages</TabsTrigger>
        <TabsTrigger value="system" className="text-xs">System</TabsTrigger>
      </TabsList>
      <TabsContent value={activeTab}>
        <ScrollArea className="h-[400px]">
          {/* Notification items with priority color coding */}
        </ScrollArea>
      </TabsContent>
    </Tabs>
  </DrawerContent>
</Drawer>
```

**Notification item pattern:**
- Left border color by priority: red (high), yellow (medium), blue (low)
- Icon by type: MessageCircle, Shield, AlertTriangle, CheckCircle, Info
- Unread indicator: blue dot `h-2 w-2 rounded-full bg-blue-500`
- Relative timestamp: "Just now", "5m ago", "3h ago", "2d ago"
- Click to navigate: `router.visit(route)` based on `related_type`

## Calendar Widget

Smart sidebar widget combining calendar picker with event list and RSVP:

```tsx
<Card className="flex flex-col">
  <CardHeader className="pb-2 px-4">
    <CardTitle className="flex items-center gap-2 text-base">
      <CalendarIcon className="h-4 w-4" /> Events
    </CardTitle>
    <CardDescription className="text-xs">
      {district ? `${district} & nationwide` : "Nationwide"}
    </CardDescription>
  </CardHeader>
  <CardContent className="flex flex-col gap-3 px-4 pb-4">
    {/* Calendar with event date highlighting */}
    <Calendar
      mode="single"
      selected={selectedDate}
      onSelect={setSelectedDate}
      modifiers={{ hasEvent: eventDates }}
      modifiersClassNames={{
        hasEvent: "bg-primary/20 font-semibold text-primary hover:bg-primary/30",
      }}
      className="rounded-md border p-2 text-xs"
    />

    {/* Event list: selected date or upcoming */}
    <ScrollArea className="h-[140px]">
      {displayEvents.map(event => (
        <div className="flex flex-col gap-1 rounded-md border p-2 hover:bg-muted/50">
          <div className="flex items-start justify-between gap-1">
            <h5 className="font-medium text-xs line-clamp-1">{event.title}</h5>
            <Badge variant={urgency.variant} className="text-[10px] px-1.5 py-0">
              {urgency.label}
            </Badge>
          </div>
          <div className="flex items-center gap-2 text-[10px] text-muted-foreground">
            <span className="flex items-center gap-0.5">
              <Clock className="h-2.5 w-2.5" /> {formatTime(event.date)}
            </span>
            <span className="flex items-center gap-0.5 truncate">
              <MapPin className="h-2.5 w-2.5" /> {event.venue}
            </span>
          </div>
          {/* RSVP button with optimistic toggle */}
          <Button size="sm" variant={isAttending ? "default" : "outline"}
            className="w-full mt-1 h-7 text-[10px]"
            onClick={() => handleToggleAttendance(event.id)}>
            {isAttending ? "Attending" : "Attend"}
          </Button>
        </div>
      ))}
    </ScrollArea>
  </CardContent>
</Card>
```

## Optimistic Updates

For any toggle/action that hits an API:

```typescript
const handleToggle = async (itemId: number) => {
  const wasActive = activeItems.has(itemId);

  // 1. Update UI immediately
  setActiveItems(prev => {
    const next = new Set(prev);
    wasActive ? next.delete(itemId) : next.add(itemId);
    return next;
  });

  try {
    // 2. API call
    if (wasActive) {
      await axios.delete(`/api/resource/${itemId}/action`);
    } else {
      await axios.post(`/api/resource/${itemId}/action`);
    }
  } catch (error) {
    // 3. Revert on failure
    setActiveItems(prev => {
      const next = new Set(prev);
      wasActive ? next.add(itemId) : next.delete(itemId);
      return next;
    });
  }
};
```

## CrudPage Read-Only Integration

Wrap `CrudPage` with a custom header for user-scoped views:

```tsx
export default function MyItemsIndex({ data, filters, permissions, smeName, error }) {
  const isMobile = useIsMobile();

  return (
    <Admin title="My Items">
      <Head title="My Items" />
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Custom header with context */}
        <div className="px-4 lg:px-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight">My Items</h1>
            <p className="text-muted-foreground">
              <span className="hidden sm:inline">Full description</span>
              <span className="sm:hidden">Short</span>
            </p>
          </div>
          {smeName && (
            <div className="flex flex-col items-end gap-2 shrink-0">
              <div className="text-lg font-semibold">{smeName}</div>
              <Badge variant="outline" className="text-xs font-mono">{usmeNumber}</Badge>
            </div>
          )}
        </div>

        {/* CrudPage in read-only mode */}
        <div className="px-0">
          <CrudPage<Item>
            data={data}
            filters={filters}
            title="Items"
            resourceName="my-items"
            columns={isMobile ? mobileColumns : desktopColumns}
            detailView={ItemDetailView}
            onRefresh={() => router.reload({ only: ['data'] })}
            canView={true}
            readOnly={true}
          />
        </div>
      </div>
    </Admin>
  );
}
```

**Backend**: Use `SetMandatoryFilterProvider` to scope data:

```go
controller.SetMandatoryFilterProvider(func(ctx http.Context) map[string]interface{} {
    user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
    sme, _ := smeService.GetSmeByUserEmail(user.Email)
    return map[string]interface{}{"sme_id": sme.ID}
})
```

## Custom Page Controller (Go)

For fully custom pages that don't use GenericPageController:

```go
type FeaturePageController struct {
    smeService       *services.SmeService
    dashboardService *services.DashboardService
}

func (c *FeaturePageController) Index(ctx http.Context) http.Response {
    // 1. Authenticate
    user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
    if user == nil { return ctx.Response().Redirect(http.StatusFound, "/login") }

    // 2. Get user context (linked SME)
    sme, err := c.smeService.GetSmeByUserEmail(user.Email)
    if err != nil {
        return inertiaHelper.Render(ctx, "Feature/Index", map[string]interface{}{
            "items": []ItemDTO{}, "error": "No MSME linked to your account",
        })
    }

    // 3. Fetch and transform data based on user context
    items := fetchItemsForSme(sme)

    // 4. Include sidebar widget data
    calendarEvents := c.dashboardService.GetEventsForDistrictWithAttendance(sme.District, sme.ID)

    // 5. Render with all props
    return inertiaHelper.Render(ctx, "Feature/Index", map[string]interface{}{
        "items":            items,
        "userName":         user.Name,
        "smeName":          sme.Name,
        "usmeNumber":       sme.UsmeNumber,
        "classification":   sme.Classification,
        "district":         sme.District,
        "calendarEvents":   calendarEvents,
        "smeId":            sme.ID,
    })
}
```

## Stats Summary Row

Compact KPI cards for custom pages:

```tsx
<div className="grid grid-cols-2 md:grid-cols-4 gap-4">
  <Card>
    <CardContent className="pt-6">
      <div className="flex items-center gap-2">
        <Target className="h-5 w-5 text-primary" />
        <div>
          <p className="text-2xl font-bold">{score}%</p>
          <p className="text-xs text-muted-foreground">Your Score</p>
        </div>
      </div>
    </CardContent>
  </Card>
  {/* More stat cards... */}
</div>
```

## Empty State

```tsx
const EmptyState: React.FC<{ type: string }> = ({ type }) => (
  <div className="flex flex-col items-center justify-center py-12 text-center">
    <Icon className="h-12 w-12 text-muted-foreground/50 mb-4" />
    <h3 className="text-lg font-medium">No {type}</h3>
    <p className="text-sm text-muted-foreground mt-1 max-w-sm">
      Helpful guidance text here.
    </p>
  </div>
);
```

## Tip / CTA Card

Gradient card for tips or calls to action:

```tsx
<Card className="bg-gradient-to-r from-primary/10 to-primary/5 border-primary/20">
  <CardContent className="pt-6">
    <div className="flex gap-4">
      <div className="shrink-0"><Target className="h-8 w-8 text-primary" /></div>
      <div>
        <h3 className="font-semibold mb-1">Improve Your Score</h3>
        <p className="text-sm text-muted-foreground">
          Helpful guidance about what the user can do next.
        </p>
      </div>
    </div>
  </CardContent>
</Card>
```

## Checklist

When building a new custom page:

- [ ] Use `Admin` layout wrapper with `<Head title="..." />`
- [ ] Header: icon + title + responsive descriptions (`hidden md:inline` / `md:hidden`)
- [ ] Context badges in header (SME name, USME number, score, classification)
- [ ] Error state with `Alert variant="destructive"`
- [ ] Main + sidebar layout: `flex flex-col xl:flex-row gap-6`
- [ ] Sidebar: `w-full xl:w-[380px] 2xl:w-[420px] shrink-0`
- [ ] Detail view: fullscreen `Dialog` (98vw x 96vh) with sidebar metadata + main content
- [ ] MetadataItem: icon + `text-[10px] uppercase` label + `text-sm font-medium` value
- [ ] Cards as click targets with `cursor-pointer hover:shadow-md transition-all`
- [ ] Tabs for content categories (e.g., Upcoming / Past)
- [ ] Empty states: centered icon + title + help text
- [ ] Loading states: `animate-pulse` skeletons
- [ ] Optimistic updates for toggles/actions
- [ ] Local state synced with Inertia props via `useEffect`
- [ ] Mobile columns via `useIsMobile()` hook
- [ ] Badge color coding for status/urgency/classification
- [ ] `toast()` from sonner for save feedback (not validation errors)
- [ ] Calendar widget in sidebar if events are relevant
- [ ] Permission gating for sections and actions

For detailed code examples, see [examples/opportunities-page.md](examples/opportunities-page.md) and [examples/widgets-and-overlays.md](examples/widgets-and-overlays.md).
