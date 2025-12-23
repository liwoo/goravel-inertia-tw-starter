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
} from 'lucide-react';
import { cn } from '@/lib/utils';
import MarkdownEditor from '@uiw/react-markdown-editor';

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

// Opportunity Detail Modal Component
const OpportunityDetailModal: React.FC<{
  opportunity: Opportunity | null;
  isOpen: boolean;
  onClose: () => void;
}> = ({ opportunity, isOpen, onClose }) => {
  if (!opportunity) return null;

  const statusColor = opportunity.is_open
    ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
    : 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-300';

  const statusText = opportunity.is_open
    ? 'Open'
    : `Opens ${formatDate(opportunity.open_date)}`;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent
        showCloseButton={false}
        className="fixed inset-4 w-auto max-w-none h-auto translate-x-0 translate-y-0 top-0 left-0 p-0 flex flex-col gap-0"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-border shrink-0">
          <div className="flex items-center gap-4">
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
                {opportunity.organization}
              </DialogTitle>
              <div className="flex items-center gap-2 mt-1">
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
          </div>
        </div>

        {/* Scrollable Content */}
        <ScrollArea className="flex-1">
          <div className="p-6 space-y-6">
            {/* Quick Info */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="flex items-center gap-2 text-sm">
                <Building2 className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-muted-foreground">Procured By</p>
                  <p className="font-medium">{opportunity.procured_by}</p>
                </div>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <FileText className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-muted-foreground">Type</p>
                  <p className="font-medium">{opportunity.procurement_type}</p>
                </div>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <Target className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-muted-foreground">Min Score</p>
                  <p className="font-medium">{opportunity.minimum_qualifying_score}%</p>
                </div>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <Calendar className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-muted-foreground">Dates</p>
                  <p className="font-medium">
                    {formatDate(opportunity.open_date)} - {formatDate(opportunity.close_date)}
                  </p>
                </div>
              </div>
            </div>

            {/* Classification Badges */}
            {opportunity.classification && opportunity.classification.length > 0 && (
              <div>
                <h4 className="font-medium text-sm mb-2">Classification</h4>
                <div className="flex flex-wrap gap-2">
                  {opportunity.classification.map((cls) => (
                    <Badge key={cls} variant="secondary">
                      {cls}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {/* Qualifying Districts */}
            {opportunity.qualifying_districts && opportunity.qualifying_districts.length > 0 && (
              <div>
                <h4 className="font-medium text-sm mb-2 flex items-center gap-1">
                  <MapPin className="h-4 w-4" />
                  Qualifying Districts
                </h4>
                <div className="flex flex-wrap gap-2">
                  {opportunity.qualifying_districts.map((district) => (
                    <Badge key={district} variant="outline">
                      {district}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            <Separator />

            {/* Details Section */}
            <div>
              <h3 className="text-lg font-semibold mb-4">Details</h3>
              <div className="prose prose-sm max-w-none dark:prose-invert border rounded-md p-4 bg-card">
                <MarkdownEditor.Markdown
                  source={opportunity.details || 'No details provided.'}
                  style={{ backgroundColor: 'transparent' }}
                />
              </div>
            </div>

            {/* How to Apply Section */}
            <div>
              <h3 className="text-lg font-semibold mb-4">How to Apply</h3>
              <div className="prose prose-sm max-w-none dark:prose-invert border rounded-md p-4 bg-card">
                <MarkdownEditor.Markdown
                  source={opportunity.application_details || 'No application details provided.'}
                  style={{ backgroundColor: 'transparent' }}
                />
              </div>
            </div>

            {/* Additional Info */}
            <div className="flex items-center gap-4 text-sm text-muted-foreground pt-4 border-t">
              <span>
                <strong>Market Approach:</strong> {opportunity.market_approach}
              </span>
              <span>
                <strong>Invitation:</strong> {opportunity.invitation}
              </span>
            </div>
          </div>
        </ScrollArea>
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
}: OpportunitiesIndexProps) {
  const [selectedOpportunity, setSelectedOpportunity] = useState<Opportunity | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const totalOpportunities = upcomingOpportunities.length + pastOpportunities.length;
  const openOpportunities = upcomingOpportunities.filter((o) => o.is_open).length;

  const handleOpportunityClick = (opportunity: Opportunity) => {
    setSelectedOpportunity(opportunity);
    setIsModalOpen(true);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
    setSelectedOpportunity(null);
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
              Upcoming ({upcomingOpportunities.length})
            </TabsTrigger>
            <TabsTrigger value="past" className="flex items-center gap-2">
              <Clock className="h-4 w-4" />
              Past ({pastOpportunities.length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="upcoming" className="mt-6">
            {upcomingOpportunities.length === 0 ? (
              <EmptyState type="upcoming" />
            ) : (
              <div className="grid gap-4">
                {upcomingOpportunities.map((opportunity) => (
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
            {pastOpportunities.length === 0 ? (
              <EmptyState type="past" />
            ) : (
              <div className="grid gap-4">
                {pastOpportunities.map((opportunity) => (
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
                  Update your business details in the SME Portal to improve your score.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Opportunity Detail Modal */}
      <OpportunityDetailModal
        opportunity={selectedOpportunity}
        isOpen={isModalOpen}
        onClose={handleModalClose}
      />
    </Admin>
  );
}
