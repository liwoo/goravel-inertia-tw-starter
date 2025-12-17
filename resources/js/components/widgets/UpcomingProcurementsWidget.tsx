"use client";

import * as React from "react";
import { FileText, Clock, Building2, ArrowRight } from "lucide-react";
import { Link } from "@inertiajs/react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

// Interface for upcoming procurement notice data
export interface UpcomingProcurement {
  id: number;
  organization: string;
  refNo: string;
  closeDate: string;
  procurementType?: string;
  invitation?: string;
}

interface UpcomingProcurementsWidgetProps {
  procurements: UpcomingProcurement[];
  isLoading?: boolean;
}

// Parse date string safely - handles various formats including ISO and carbon.DateTime
function parseDate(dateString: string | null | undefined): Date | null {
  if (!dateString) return null;

  // Try parsing as-is first
  let date = new Date(dateString);

  // If invalid, try other formats
  if (isNaN(date.getTime())) {
    // Try parsing "YYYY-MM-DD HH:mm:ss" format (common from Go/Carbon)
    const match = dateString.match(/^(\d{4})-(\d{2})-(\d{2})/);
    if (match) {
      date = new Date(parseInt(match[1]), parseInt(match[2]) - 1, parseInt(match[3]));
    }
  }

  return isNaN(date.getTime()) ? null : date;
}

// Format date for display
function formatCloseDate(dateString: string): string {
  const date = parseDate(dateString);
  if (!date) return "Date TBD";

  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

// Calculate days until close
function getDaysUntilClose(dateString: string): number | null {
  const closeDate = parseDate(dateString);
  if (!closeDate) return null;

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  closeDate.setHours(0, 0, 0, 0);
  const diffTime = closeDate.getTime() - today.getTime();
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}

// Get badge variant based on days until close
function getUrgencyBadge(daysUntil: number | null): { variant: "default" | "secondary" | "destructive" | "outline"; label: string } | null {
  if (daysUntil === null) return null;

  if (daysUntil <= 0) {
    return { variant: "destructive", label: "Closing Today" };
  } else if (daysUntil === 1) {
    return { variant: "destructive", label: "Closing Tomorrow" };
  } else if (daysUntil <= 3) {
    return { variant: "destructive", label: `${daysUntil} days left` };
  } else if (daysUntil <= 7) {
    return { variant: "default", label: `${daysUntil} days left` };
  } else {
    return { variant: "secondary", label: `${daysUntil} days left` };
  }
}

export function UpcomingProcurementsWidget({
  procurements,
  isLoading = false,
}: UpcomingProcurementsWidgetProps) {
  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5" />
            Active Procurement Notices
          </CardTitle>
          <CardDescription>Open tenders and procurement opportunities</CardDescription>
        </CardHeader>
        <CardContent className="flex-1">
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-16 animate-pulse rounded-md bg-muted" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!procurements || procurements.length === 0) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5" />
            Active Procurement Notices
          </CardTitle>
          <CardDescription>Open tenders and procurement opportunities</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center py-8">
          <div className="text-center">
            <FileText className="mx-auto h-12 w-12 text-muted-foreground/50" />
            <p className="mt-2 text-sm text-muted-foreground">
              No active procurement notices
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="flex flex-col">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileText className="h-5 w-5" />
          Active Procurement Notices
        </CardTitle>
        <CardDescription>Open tenders and procurement opportunities</CardDescription>
      </CardHeader>
      <CardContent className="flex-1">
        <div className="space-y-3">
          {procurements.slice(0, 5).map((procurement) => {
            const daysUntil = getDaysUntilClose(procurement.closeDate);
            const urgency = getUrgencyBadge(daysUntil);

            return (
              <div
                key={procurement.id}
                className="flex items-start gap-3 rounded-lg border p-3 transition-colors hover:bg-muted/50"
              >
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h4 className="font-medium text-sm truncate flex items-center gap-1">
                      <Building2 className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                      {procurement.organization}
                    </h4>
                  </div>
                  <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                    <span className="font-mono">{procurement.refNo}</span>
                    {procurement.procurementType && (
                      <Badge variant="outline" className="text-xs">
                        {procurement.procurementType}
                      </Badge>
                    )}
                  </div>
                  <div className="mt-1.5 flex items-center gap-2">
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Clock className="h-3 w-3" />
                      Closes: {formatCloseDate(procurement.closeDate)}
                    </span>
                    {urgency && (
                      <Badge variant={urgency.variant} className="text-xs">
                        {urgency.label}
                      </Badge>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
      <CardFooter className="border-t pt-4">
        <Button variant="ghost" className="w-full" asChild>
          <Link href="/admin/procurement-notices">
            View All Procurements
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

export default UpcomingProcurementsWidget;
