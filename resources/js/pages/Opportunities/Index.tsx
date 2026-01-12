import React, { useState } from 'react';
import { Head } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
  Dialog,
  DialogContent,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  CalendarDays,
  Building2,
  Clock,
  MapPin,
  Target,
  Sparkles,
  CheckCircle2,
  XCircle,
  AlertCircle,
  ChevronRight,
  FileText,
  Calendar,
  ArrowRight,
  X,
  Heart,
  Loader2,
  Users,
} from 'lucide-react';
import axios from 'axios';
import { cn } from '@/lib/utils';
import MarkdownEditor from '@uiw/react-markdown-editor';
import { EventsCalendarWidget, CalendarEvent } from '@/components/widgets';

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

// Helper to format date nicely
const formatDate = (dateStr: string): string => {
  if (!dateStr) return '';
  try {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-GB', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  } catch {
    return dateStr;
  }
};

// Opportunity Card Component
const OpportunityCard: React.FC<{
  opportunity: Opportunity;
  isPast?: boolean;
  onClick?: () => void;
}> = ({
  opportunity,
  isPast = false,
  onClick,
}) => {
  const statusColor = isPast
    ? 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'
    : opportunity.is_open
    ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
    : 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-300';

  const statusText = isPast
    ? 'Closed'
    : opportunity.is_open
    ? 'Open'
    : `Opens ${formatDate(opportunity.open_date)}`;

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
              <Badge variant="outline" className="text-xs font-mono">
                {opportunity.ref_no}
              </Badge>
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
        {/* Quick Info */}
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

        {/* Classification Badges */}
        {opportunity.classification && opportunity.classification.length > 0 && (
          <div className="flex flex-wrap gap-1 mt-3">
            {opportunity.classification.map((cls) => (
              <Badge key={cls} variant="secondary" className="text-xs">
                {cls}
              </Badge>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};

// Empty State Component
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

// Opportunity Detail Modal Component - Full screen modal with sidebar metadata and A4 document viewer
const OpportunityDetailModal: React.FC<{
  opportunity: Opportunity | null;
  isOpen: boolean;
  onClose: () => void;
  onInterestChange?: (opportunityId: number, hasInterest: boolean, count: number) => void;
}> = ({ opportunity, isOpen, onClose, onInterestChange }) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [localHasInterest, setLocalHasInterest] = useState(false);
  const [localInterestCount, setLocalInterestCount] = useState(0);

  // Sync local state with opportunity prop
  React.useEffect(() => {
    if (opportunity) {
      setLocalHasInterest(opportunity.has_shown_interest);
      setLocalInterestCount(opportunity.interested_count);
    }
  }, [opportunity]);

  if (!opportunity) return null;

  const statusColor = opportunity.is_open
    ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
    : 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-300';

  const statusText = opportunity.is_open
    ? 'Open'
    : `Opens ${formatDate(opportunity.open_date)}`;

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

  // Metadata item component for sidebar
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

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent
        showCloseButton={false}
        className="!max-w-[98vw] !w-[98vw] !max-h-[96vh] !h-[96vh] p-0 flex flex-col gap-0"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-border shrink-0 bg-background">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-8 w-8"
            >
              <X className="h-4 w-4" />
            </Button>
            <div>
              <DialogTitle className="text-lg font-semibold">
                {opportunity.invitation}
              </DialogTitle>
              <p className="text-xs text-muted-foreground">{opportunity.organization}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="text-xs font-mono">
              {opportunity.ref_no}
            </Badge>
            <Badge className={statusColor}>{statusText}</Badge>
            {opportunity.days_remaining > 0 && opportunity.days_remaining <= 7 && (
              <Badge variant="destructive" className="text-xs">
                {opportunity.days_remaining} days left
              </Badge>
            )}
          </div>
        </div>

        {/* Main Content: Sidebar + Document Viewer */}
        <div className="flex-1 flex overflow-hidden">
          {/* Left Sidebar - Metadata */}
          <aside className="w-72 border-r border-border bg-muted/30 flex flex-col shrink-0">
            <ScrollArea className="flex-1">
              <div className="p-4 space-y-6">
                {/* Organization Info */}
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
                    <MetadataItem
                      icon={<Building2 className="h-3.5 w-3.5 text-muted-foreground" />}
                      label="Procured By"
                      value={opportunity.procured_by}
                    />
                    <MetadataItem
                      icon={<FileText className="h-3.5 w-3.5 text-muted-foreground" />}
                      label="Type"
                      value={opportunity.procurement_type}
                    />
                    <MetadataItem
                      icon={<Sparkles className="h-3.5 w-3.5 text-muted-foreground" />}
                      label="Market Approach"
                      value={opportunity.market_approach}
                    />
                  </div>
                </div>

                <Separator />

                {/* Dates */}
                <div>
                  <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
                    Timeline
                  </h4>
                  <div className="space-y-3">
                    <MetadataItem
                      icon={<CalendarDays className="h-3.5 w-3.5 text-muted-foreground" />}
                      label="Opens"
                      value={formatDate(opportunity.open_date)}
                    />
                    <MetadataItem
                      icon={<Clock className="h-3.5 w-3.5 text-muted-foreground" />}
                      label="Closes"
                      value={formatDate(opportunity.close_date)}
                    />
                  </div>
                </div>

                <Separator />

                {/* Requirements */}
                <div>
                  <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
                    Requirements
                  </h4>
                  <div className="space-y-3">
                    <MetadataItem
                      icon={<Target className="h-3.5 w-3.5 text-muted-foreground" />}
                      label="Min. Score"
                      value={`${opportunity.minimum_qualifying_score}%`}
                    />
                    {opportunity.classification && opportunity.classification.length > 0 && (
                      <div>
                        <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-medium mb-1.5">
                          Classification
                        </p>
                        <div className="flex flex-wrap gap-1">
                          {opportunity.classification.map((cls) => (
                            <Badge key={cls} variant="secondary" className="text-[10px]">
                              {cls}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                    {opportunity.qualifying_districts && opportunity.qualifying_districts.length > 0 && (
                      <div>
                        <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-medium mb-1.5">
                          Districts
                        </p>
                        <div className="flex flex-wrap gap-1">
                          {opportunity.qualifying_districts.map((d) => (
                            <Badge key={d} variant="outline" className="text-[10px]">
                              {d}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                </div>

                <Separator />

                {/* Interest Section */}
                <div>
                  <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
                    Your Interest
                  </h4>
                  <div className="space-y-3">
                    <div className="flex items-center gap-2 text-sm">
                      <Users className="h-4 w-4 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {localInterestCount > 0
                          ? `${localInterestCount} MSME${localInterestCount > 1 ? 's' : ''} interested`
                          : 'No interest yet'}
                      </span>
                    </div>
                    <Button
                      onClick={handleInterestToggle}
                      disabled={isSubmitting || !opportunity.is_open}
                      variant={localHasInterest ? "default" : "outline"}
                      size="sm"
                      className="w-full"
                    >
                      {isSubmitting ? (
                        <>
                          <Loader2 className="mr-2 h-3.5 w-3.5 animate-spin" />
                          Processing...
                        </>
                      ) : localHasInterest ? (
                        <>
                          <Heart className="mr-2 h-3.5 w-3.5 fill-current" />
                          Interested
                        </>
                      ) : (
                        <>
                          <Heart className="mr-2 h-3.5 w-3.5" />
                          Show Interest
                        </>
                      )}
                    </Button>
                    {!opportunity.is_open && (
                      <p className="text-[10px] text-muted-foreground">
                        Available once opportunity opens
                      </p>
                    )}
                  </div>
                </div>
              </div>
            </ScrollArea>
          </aside>

          {/* Main Document Viewer - A4 aspect ratio */}
          <main className="flex-1 bg-muted/50 overflow-hidden flex items-start justify-center p-6">
            <div
              className="bg-background rounded-lg shadow-lg border overflow-hidden flex flex-col"
              style={{
                width: 'min(100%, 794px)', // A4 width at 96 DPI
                height: 'calc(96vh - 80px)',
                maxHeight: '1123px', // A4 height at 96 DPI
              }}
            >
              {/* Document Header */}
              <div className="px-8 py-6 border-b bg-gradient-to-r from-primary/5 to-transparent">
                <h1 className="text-2xl font-bold text-foreground mb-2">
                  {opportunity.invitation}
                </h1>
                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Building2 className="h-4 w-4" />
                    {opportunity.organization}
                  </span>
                  <span className="flex items-center gap-1">
                    <Calendar className="h-4 w-4" />
                    {formatDate(opportunity.open_date)} - {formatDate(opportunity.close_date)}
                  </span>
                </div>
              </div>

              {/* Document Content */}
              <ScrollArea className="flex-1">
                <div className="px-8 py-6 space-y-8">
                  {/* Details Section */}
                  <section>
                    <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
                      <FileText className="h-5 w-5 text-primary" />
                      Details
                    </h2>
                    <div className="prose prose-sm max-w-none dark:prose-invert prose-headings:text-foreground prose-p:text-foreground/90 prose-li:text-foreground/90">
                      <MarkdownEditor.Markdown
                        source={opportunity.details || 'No details provided.'}
                        style={{ backgroundColor: 'transparent' }}
                      />
                    </div>
                  </section>

                  <Separator />

                  {/* How to Apply Section */}
                  <section>
                    <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
                      <ArrowRight className="h-5 w-5 text-primary" />
                      How to Apply
                    </h2>
                    <div className="prose prose-sm max-w-none dark:prose-invert prose-headings:text-foreground prose-p:text-foreground/90 prose-li:text-foreground/90">
                      <MarkdownEditor.Markdown
                        source={opportunity.application_details || 'No application details provided.'}
                        style={{ backgroundColor: 'transparent' }}
                      />
                    </div>
                  </section>
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
  );
};

export default function OpportunitiesIndex({
  upcomingOpportunities = [],
  pastOpportunities = [],
  formalisationScore = 0,
  userName,
  smeName,
  usmeNumber,
  classification,
  district,
  error,
  calendarEvents = [],
  smeId,
}: OpportunitiesIndexProps) {
  const [selectedOpportunity, setSelectedOpportunity] = useState<Opportunity | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  // Local state for opportunities to handle interest updates without page reload
  const [localUpcoming, setLocalUpcoming] = useState(upcomingOpportunities);
  const [localPast, setLocalPast] = useState(pastOpportunities);

  // Sync with props when they change
  React.useEffect(() => {
    setLocalUpcoming(upcomingOpportunities);
    setLocalPast(pastOpportunities);
  }, [upcomingOpportunities, pastOpportunities]);

  const totalOpportunities = localUpcoming.length + localPast.length;
  const openOpportunities = localUpcoming.filter((o) => o.is_open).length;

  const handleOpportunityClick = (opportunity: Opportunity) => {
    setSelectedOpportunity(opportunity);
    setIsModalOpen(true);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
    setSelectedOpportunity(null);
  };

  const handleInterestChange = (opportunityId: number, hasInterest: boolean, count: number) => {
    // Update local state for upcoming opportunities
    setLocalUpcoming(prev => prev.map(opp =>
      opp.id === opportunityId
        ? { ...opp, has_shown_interest: hasInterest, interested_count: count }
        : opp
    ));
    // Update local state for past opportunities
    setLocalPast(prev => prev.map(opp =>
      opp.id === opportunityId
        ? { ...opp, has_shown_interest: hasInterest, interested_count: count }
        : opp
    ));
    // Update selected opportunity
    if (selectedOpportunity?.id === opportunityId) {
      setSelectedOpportunity(prev => prev ? {
        ...prev,
        has_shown_interest: hasInterest,
        interested_count: count
      } : null);
    }
  };

  return (
    <Admin title="Opportunities">
      <Head title="Opportunities" />

      <div className="p-4 md:p-6 flex flex-col gap-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight flex items-center gap-2">
              <Sparkles className="h-6 w-6 text-amber-500" />
              Opportunities
            </h1>
            <p className="text-muted-foreground">
              Procurement opportunities matching your business profile and formalisation score
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

        {/* Error State */}
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Main Content + Sidebar Layout */}
        <div className="flex flex-col xl:flex-row gap-6 items-start">
          {/* Main Content Area */}
          <div className="flex-1 flex flex-col gap-6 min-w-0">
            {/* Stats Summary */}
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
              <Card>
                <CardContent className="pt-6">
                  <div className="flex items-center gap-2">
                    <CheckCircle2 className="h-5 w-5 text-green-500" />
                    <div>
                      <p className="text-2xl font-bold">{openOpportunities}</p>
                      <p className="text-xs text-muted-foreground">Open Now</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="pt-6">
                  <div className="flex items-center gap-2">
                    <MapPin className="h-5 w-5 text-blue-500" />
                    <div>
                      <p className="text-2xl font-bold truncate">{district || 'All'}</p>
                      <p className="text-xs text-muted-foreground">Your District</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Opportunities Tabs */}
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

              <TabsContent value="past" className="mt-6">
                {localPast.length === 0 ? (
                  <EmptyState type="past" />
                ) : (
                  <div className="grid gap-4">
                    {localPast.map((opportunity) => (
                      <OpportunityCard
                        key={opportunity.id}
                        opportunity={opportunity}
                        isPast
                        onClick={() => handleOpportunityClick(opportunity)}
                      />
                    ))}
                  </div>
                )}
              </TabsContent>
            </Tabs>

            {/* Tip Card */}
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
          </div>

          {/* Right Sidebar - Calendar Widget */}
          <aside className="w-full xl:w-[380px] 2xl:w-[420px] shrink-0">
            <EventsCalendarWidget events={calendarEvents} district={district} smeId={smeId} />
          </aside>
        </div>
      </div>

      {/* Opportunity Detail Modal */}
      <OpportunityDetailModal
        opportunity={selectedOpportunity}
        isOpen={isModalOpen}
        onClose={handleModalClose}
        onInterestChange={handleInterestChange}
      />
    </Admin>
  );
}
