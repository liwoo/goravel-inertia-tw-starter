import React, { useState } from 'react';
import { Head } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
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
} from 'lucide-react';
import { cn } from '@/lib/utils';

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
const OpportunityCard: React.FC<{ opportunity: Opportunity; isPast?: boolean }> = ({
  opportunity,
  isPast = false,
}) => {
  const [isExpanded, setIsExpanded] = useState(false);

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
      onClick={() => setIsExpanded(!isExpanded)}
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
          <ChevronRight
            className={cn(
              'h-5 w-5 text-muted-foreground transition-transform shrink-0',
              isExpanded && 'rotate-90'
            )}
          />
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

        {/* Expanded Details */}
        {isExpanded && (
          <div className="mt-4 pt-4 border-t space-y-4">
            <div>
              <h4 className="font-medium text-sm mb-2">Details</h4>
              <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                {opportunity.details}
              </p>
            </div>

            <div>
              <h4 className="font-medium text-sm mb-2">How to Apply</h4>
              <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                {opportunity.application_details}
              </p>
            </div>

            {opportunity.qualifying_districts && opportunity.qualifying_districts.length > 0 && (
              <div>
                <h4 className="font-medium text-sm mb-2 flex items-center gap-1">
                  <MapPin className="h-4 w-4" />
                  Qualifying Districts
                </h4>
                <div className="flex flex-wrap gap-1">
                  {opportunity.qualifying_districts.map((district) => (
                    <Badge key={district} variant="outline" className="text-xs">
                      {district}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            <div className="flex items-center gap-4 text-xs text-muted-foreground">
              <span>
                <strong>Market Approach:</strong> {opportunity.market_approach}
              </span>
              <span>
                <strong>Invitation:</strong> {opportunity.invitation}
              </span>
            </div>
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
  const totalOpportunities = upcomingOpportunities.length + pastOpportunities.length;
  const openOpportunities = upcomingOpportunities.filter((o) => o.is_open).length;

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
                  <OpportunityCard key={opportunity.id} opportunity={opportunity} />
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
                  <OpportunityCard key={opportunity.id} opportunity={opportunity} isPast />
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
    </Admin>
  );
}
