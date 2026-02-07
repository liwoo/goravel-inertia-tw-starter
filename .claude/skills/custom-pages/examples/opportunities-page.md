# Opportunities Page - Reference Example

The full custom page implementation from `resources/js/pages/Opportunities/Index.tsx` and `app/http/controllers/opportunities/opportunities_page_controller.go`.

## TypeScript Interfaces

```typescript
interface Opportunity {
  id: number;
  ref_no: string;
  organization: string;
  procured_by: string;
  procurement_type: string;
  market_approach: string;
  invitation: string;
  details: string;
  application_details: string;
  open_date: string;
  close_date: string;
  minimum_qualifying_score: number;
  classification: string[];
  qualifying_districts: string[];
  is_open: boolean;
  days_remaining: number;
  has_shown_interest: boolean;
  interested_count: number;
}

interface OpportunitiesIndexProps {
  upcomingOpportunities: Opportunity[];
  pastOpportunities: Opportunity[];
  formalisationScore: number;
  userName?: string;
  smeName?: string;
  usmeNumber?: string;
  classification?: string;
  district?: string;
  error?: string;
  calendarEvents?: CalendarEvent[];
  smeId?: number;
}
```

## Page Header with Context Badges

Responsive description (long on desktop, short on mobile) plus SME identity badges:

```tsx
<div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
  <div className="flex flex-col gap-1">
    <h1 className="text-2xl font-semibold tracking-tight flex items-center gap-2">
      <Sparkles className="h-6 w-6 text-amber-500" />
      Opportunities
    </h1>
    <p className="text-muted-foreground">
      <span className="hidden md:inline">Procurement opportunities matching your business profile and formalisation score</span>
      <span className="md:hidden">Opportunities for your business</span>
    </p>
  </div>
  {smeName && (
    <div className="flex flex-col items-end gap-2 shrink-0">
      <div className="text-lg font-semibold text-right">{smeName}</div>
      <div className="flex items-center gap-2">
        {usmeNumber && (
          <Badge variant="outline" className="text-xs font-mono">
            {usmeNumber}
          </Badge>
        )}
        <Badge variant="default" className="bg-primary">
          Score: {formalisationScore}%
        </Badge>
      </div>
    </div>
  )}
</div>
```

## Stats Summary Row

Compact card grid with icon + value + label:

```tsx
<div className="grid grid-cols-2 md:grid-cols-4 gap-4">
  <Card>
    <CardContent className="pt-6">
      <div className="flex items-center gap-2">
        <Target className="h-5 w-5 text-primary" />
        <div>
          <p className="text-2xl font-bold">{formalisationScore}%</p>
          <p className="text-xs text-muted-foreground">Your Score</p>
        </div>
      </div>
    </CardContent>
  </Card>
  <Card>
    <CardContent className="pt-6">
      <div className="flex items-center gap-2">
        <Sparkles className="h-5 w-5 text-amber-500" />
        <div>
          <p className="text-2xl font-bold">{totalOpportunities}</p>
          <p className="text-xs text-muted-foreground">Total Matches</p>
        </div>
      </div>
    </CardContent>
  </Card>
  {/* ... more stat cards */}
</div>
```

## Tabbed Content with Counts

```tsx
<Tabs defaultValue="upcoming" className="w-full">
  <TabsList className="grid w-full max-w-md grid-cols-2">
    <TabsTrigger value="upcoming" className="flex items-center gap-2">
      <CalendarDays className="h-4 w-4" />
      Upcoming ({localUpcoming.length})
    </TabsTrigger>
    <TabsTrigger value="past" className="flex items-center gap-2">
      <Clock className="h-4 w-4" />
      Past ({localPast.length})
    </TabsTrigger>
  </TabsList>

  <TabsContent value="upcoming" className="mt-6">
    {localUpcoming.length === 0 ? (
      <EmptyState type="upcoming" />
    ) : (
      <div className="grid gap-4">
        {localUpcoming.map((opportunity) => (
          <OpportunityCard
            key={opportunity.id}
            opportunity={opportunity}
            onClick={() => handleOpportunityClick(opportunity)}
          />
        ))}
      </div>
    )}
  </TabsContent>
  {/* ... past tab */}
</Tabs>
```

## Opportunity Card (List Item)

Clickable card with badges, icon metadata, and hover effect:

```tsx
const OpportunityCard: React.FC<{
  opportunity: Opportunity;
  isPast?: boolean;
  onClick?: () => void;
}> = ({ opportunity, isPast = false, onClick }) => {
  const statusColor = isPast
    ? 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'
    : opportunity.is_open
    ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
    : 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-300';

  return (
    <Card
      className={cn(
        'transition-all duration-200 hover:shadow-md cursor-pointer',
        isPast && 'opacity-75'
      )}
      onClick={onClick}
    >
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <Badge variant="outline" className="text-xs font-mono">{opportunity.ref_no}</Badge>
              <Badge className={statusColor}>{statusText}</Badge>
              {!isPast && opportunity.days_remaining > 0 && opportunity.days_remaining <= 7 && (
                <Badge variant="destructive" className="text-xs">
                  {opportunity.days_remaining} days left
                </Badge>
              )}
            </div>
            <CardTitle className="text-lg line-clamp-2">{opportunity.organization}</CardTitle>
            <CardDescription className="mt-1 flex items-center gap-1">
              <Building2 className="h-3 w-3" />
              {opportunity.procured_by}
            </CardDescription>
          </div>
          <ChevronRight className="h-5 w-5 text-muted-foreground shrink-0" />
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <FileText className="h-4 w-4" />
            <span>{opportunity.procurement_type}</span>
          </div>
          <div className="flex items-center gap-1">
            <Target className="h-4 w-4" />
            <span>Min Score: {opportunity.minimum_qualifying_score}%</span>
          </div>
          <div className="flex items-center gap-1">
            <Calendar className="h-4 w-4" />
            <span>
              {formatDate(opportunity.open_date)} <ArrowRight className="h-3 w-3 inline" />{' '}
              {formatDate(opportunity.close_date)}
            </span>
          </div>
        </div>
        {opportunity.classification?.length > 0 && (
          <div className="flex flex-wrap gap-1 mt-3">
            {opportunity.classification.map((cls) => (
              <Badge key={cls} variant="secondary" className="text-xs">{cls}</Badge>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
```

## Fullscreen Detail Modal

The core pattern: `98vw x 96vh` Dialog with sidebar + A4 document viewer:

```tsx
<Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
  <DialogContent
    showCloseButton={false}
    className="!max-w-[98vw] !w-[98vw] !max-h-[96vh] !h-[96vh] p-0 flex flex-col gap-0"
  >
    {/* Header Bar */}
    <div className="flex items-center justify-between px-4 py-3 border-b border-border shrink-0 bg-background">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8">
          <X className="h-4 w-4" />
        </Button>
        <div>
          <DialogTitle className="text-lg font-semibold">{opportunity.invitation}</DialogTitle>
          <p className="text-xs text-muted-foreground">{opportunity.organization}</p>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Badge variant="outline" className="text-xs font-mono">{opportunity.ref_no}</Badge>
        <Badge className={statusColor}>{statusText}</Badge>
      </div>
    </div>

    {/* Two-Column Layout */}
    <div className="flex-1 flex overflow-hidden">
      {/* Left Sidebar - Metadata */}
      <aside className="w-72 border-r border-border bg-muted/30 flex flex-col shrink-0">
        <ScrollArea className="flex-1">
          <div className="p-4 space-y-6">
            <div>
              <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
                Organization
              </h4>
              <div className="space-y-3">
                <MetadataItem
                  icon={<Building2 className="h-3.5 w-3.5 text-muted-foreground" />}
                  label="Organization"
                  value={opportunity.organization}
                />
                {/* ... more metadata items */}
              </div>
            </div>
            <Separator />
            {/* More sections: Timeline, Requirements, Interest */}
          </div>
        </ScrollArea>
      </aside>

      {/* Main Document Viewer - A4 */}
      <main className="flex-1 bg-muted/50 overflow-hidden flex items-start justify-center p-6">
        <div
          className="bg-background rounded-lg shadow-lg border overflow-hidden flex flex-col"
          style={{
            width: 'min(100%, 794px)',      // A4 width at 96 DPI
            height: 'calc(96vh - 80px)',
            maxHeight: '1123px',            // A4 height at 96 DPI
          }}
        >
          {/* Document Header */}
          <div className="px-8 py-6 border-b bg-gradient-to-r from-primary/5 to-transparent">
            <h1 className="text-2xl font-bold text-foreground mb-2">{opportunity.invitation}</h1>
            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <span className="flex items-center gap-1">
                <Building2 className="h-4 w-4" />{opportunity.organization}
              </span>
              <span className="flex items-center gap-1">
                <Calendar className="h-4 w-4" />{formatDate(opportunity.open_date)} - {formatDate(opportunity.close_date)}
              </span>
            </div>
          </div>

          {/* Document Content */}
          <ScrollArea className="flex-1">
            <div className="px-8 py-6 space-y-8">
              <section>
                <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
                  <FileText className="h-5 w-5 text-primary" />Details
                </h2>
                <div className="prose prose-sm max-w-none dark:prose-invert">
                  <MarkdownEditor.Markdown source={opportunity.details || 'No details provided.'} />
                </div>
              </section>
              <Separator />
              {/* More sections... */}
            </div>
          </ScrollArea>

          {/* Document Footer */}
          <div className="px-8 py-3 border-t bg-muted/30 text-xs text-muted-foreground flex items-center justify-between">
            <span>Reference: {opportunity.ref_no}</span>
            <span>Minimum Score Required: {opportunity.minimum_qualifying_score}%</span>
          </div>
        </div>
      </main>
    </div>
  </DialogContent>
</Dialog>
```

## MetadataItem Component

Reusable sidebar metadata row:

```tsx
const MetadataItem: React.FC<{
  icon: React.ReactNode;
  label: string;
  value: React.ReactNode;
}> = ({ icon, label, value }) => (
  <div className="flex items-start gap-2">
    <div className="p-1.5 rounded bg-muted shrink-0">
      {icon}
    </div>
    <div className="min-w-0 flex-1">
      <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-medium">{label}</p>
      <p className="text-sm font-medium text-foreground truncate">{value}</p>
    </div>
  </div>
);
```

## Interest Toggle (State Sync Pattern)

Local state synced with parent via callback:

```tsx
const [localHasInterest, setLocalHasInterest] = useState(false);
const [localInterestCount, setLocalInterestCount] = useState(0);

// Sync local state with opportunity prop
React.useEffect(() => {
  if (opportunity) {
    setLocalHasInterest(opportunity.has_shown_interest);
    setLocalInterestCount(opportunity.interested_count);
  }
}, [opportunity]);

const handleInterestToggle = async () => {
  if (isSubmitting) return;
  setIsSubmitting(true);
  try {
    if (localHasInterest) {
      const response = await axios.delete(`/api/portal/opportunities/${opportunity.id}/interest`);
      setLocalHasInterest(false);
      setLocalInterestCount(response.data.interested_count);
      onInterestChange?.(opportunity.id, false, response.data.interested_count);
    } else {
      const response = await axios.post(`/api/portal/opportunities/${opportunity.id}/interest`);
      setLocalHasInterest(true);
      setLocalInterestCount(response.data.interested_count);
      onInterestChange?.(opportunity.id, true, response.data.interested_count);
    }
  } catch (error) {
    console.error('Error toggling interest:', error);
  } finally {
    setIsSubmitting(false);
  }
};
```

Parent updates local list state to avoid full page reload:

```tsx
const handleInterestChange = (opportunityId: number, hasInterest: boolean, count: number) => {
  setLocalUpcoming(prev => prev.map(opp =>
    opp.id === opportunityId
      ? { ...opp, has_shown_interest: hasInterest, interested_count: count }
      : opp
  ));
  setLocalPast(prev => prev.map(opp =>
    opp.id === opportunityId
      ? { ...opp, has_shown_interest: hasInterest, interested_count: count }
      : opp
  ));
};
```

## Empty State

Centered icon + heading + helpful message:

```tsx
const EmptyState: React.FC<{ type: 'upcoming' | 'past' }> = ({ type }) => (
  <div className="flex flex-col items-center justify-center py-12 text-center">
    {type === 'upcoming' ? (
      <>
        <Sparkles className="h-12 w-12 text-muted-foreground/50 mb-4" />
        <h3 className="text-lg font-medium">No Upcoming Opportunities</h3>
        <p className="text-sm text-muted-foreground mt-1 max-w-sm">
          There are no procurement opportunities matching your formalisation score at the moment.
          Keep improving your score to unlock more opportunities!
        </p>
      </>
    ) : (
      <>
        <Clock className="h-12 w-12 text-muted-foreground/50 mb-4" />
        <h3 className="text-lg font-medium">No Past Opportunities</h3>
        <p className="text-sm text-muted-foreground mt-1 max-w-sm">
          You haven't had any past procurement opportunities yet.
        </p>
      </>
    )}
  </div>
);
```

## Tip/CTA Card

Gradient background card with icon and actionable text:

```tsx
<Card className="bg-gradient-to-r from-primary/10 to-primary/5 border-primary/20">
  <CardContent className="pt-6">
    <div className="flex gap-4">
      <div className="shrink-0">
        <Target className="h-8 w-8 text-primary" />
      </div>
      <div>
        <h3 className="font-semibold mb-1">Improve Your Score</h3>
        <p className="text-sm text-muted-foreground">
          Increase your formalisation score to unlock more procurement opportunities.
          Update your business details in the MSME Portal to improve your score.
        </p>
      </div>
    </div>
  </CardContent>
</Card>
```

## Main Content + Sidebar Layout

```tsx
<div className="flex flex-col xl:flex-row gap-6 items-start">
  {/* Main Content Area */}
  <div className="flex-1 flex flex-col gap-6 min-w-0">
    {/* Stats, Tabs, Cards, Tip Card */}
  </div>

  {/* Right Sidebar - Calendar Widget */}
  <aside className="w-full xl:w-[380px] 2xl:w-[420px] shrink-0">
    <EventsCalendarWidget events={calendarEvents} district={district} smeId={smeId} />
  </aside>
</div>
```

## Go Backend: Custom Page Controller

```go
type OpportunitiesPageController struct {
    smeService       *services.SmeService
    dashboardService *services.DashboardService
}

func (c *OpportunitiesPageController) Index(ctx http.Context) http.Response {
    // 1. Auth check
    user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
    if user == nil {
        return ctx.Response().Redirect(http.StatusFound, "/login")
    }

    // 2. Get user's linked entity
    sme, err := c.smeService.GetSmeByUserEmail(user.Email)
    if err != nil || sme == nil {
        // Render with error state - still show the page
        return inertiaHelper.Render(ctx, "Opportunities/Index", map[string]interface{}{
            "upcomingOpportunities": []OpportunityItem{},
            "pastOpportunities":     []OpportunityItem{},
            "formalisationScore":    0,
            "error":                 "No MSME linked to your account",
            "calendarEvents":        c.dashboardService.GetEventsForDistrict(""),
        })
    }

    // 3. Fetch and filter data based on user context
    var allOpportunities []models.ProcurementNotice
    query := facades.Orm().Query().
        Where("is_published = ?", true).
        Where("minimum_qualifying_score <= ?", formalisationScore).
        Order("close_date DESC")
    query.Find(&allOpportunities)

    // 4. Transform and partition into categories
    now := time.Now()
    upcomingOpportunities := make([]OpportunityItem, 0)
    pastOpportunities := make([]OpportunityItem, 0)

    for _, opp := range allOpportunities {
        // District qualification check, interest check, date partitioning...
        if closeDate.After(now) {
            upcomingOpportunities = append(upcomingOpportunities, item)
        } else {
            pastOpportunities = append(pastOpportunities, item)
        }
    }

    // 5. Get sidebar widget data
    calendarEvents := c.dashboardService.GetEventsForDistrictWithAttendance(district, sme.ID)

    // 6. Render with full context
    return inertiaHelper.Render(ctx, "Opportunities/Index", map[string]interface{}{
        "upcomingOpportunities": upcomingOpportunities,
        "pastOpportunities":     pastOpportunities,
        "formalisationScore":    formalisationScore,
        "userName":              user.Name,
        "smeName":               sme.Name,
        "usmeNumber":            sme.UsmeNumber,
        "classification":        sme.Classification,
        "district":              district,
        "calendarEvents":        calendarEvents,
        "smeId":                 sme.ID,
    })
}
```

Key patterns:
- **Auth + entity lookup** at the top of every custom page
- **Graceful error state**: Render the page with empty data + error message instead of redirecting
- **Server-side filtering**: Filter by user's formalisation score and district
- **DTO transformation**: Convert models to simplified frontend-friendly structs
- **Sidebar data**: Include additional widget data (calendar events) in the same render call

## API Endpoints for Actions

Interest toggle uses separate POST/DELETE endpoints:

```go
// POST /api/portal/opportunities/:id/interest
func (c *OpportunitiesPageController) ShowInterest(ctx http.Context) http.Response {
    // Auth → Get SME → Get opportunity → Append to interested_smes array → Return updated count
}

// DELETE /api/portal/opportunities/:id/interest
func (c *OpportunitiesPageController) WithdrawInterest(ctx http.Context) http.Response {
    // Auth → Get SME → Get opportunity → Remove from interested_smes array → Return updated count
}
```

Response shape for both:
```json
{
  "message": "Interest shown successfully",
  "has_shown_interest": true,
  "interested_count": 5
}
```

## Source Files

- Page component: `resources/js/pages/Opportunities/Index.tsx`
- Page controller: `app/http/controllers/opportunities/opportunities_page_controller.go`
- Calendar widget: `resources/js/components/widgets/EventsCalendarWidget.tsx`
