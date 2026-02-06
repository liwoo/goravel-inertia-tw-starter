"use client";

import * as React from "react";
import { Calendar, MapPin, ArrowRight } from "lucide-react";
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

// Interface for upcoming event data
export interface UpcomingEvent {
  id: number;
  title: string;
  date: string;
  venue: string;
  district?: string;
}

interface UpcomingEventsWidgetProps {
  events: UpcomingEvent[];
  isLoading?: boolean;
}

// Format date for display
function formatEventDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

// Calculate days until event
function getDaysUntil(dateString: string): number {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const eventDate = new Date(dateString);
  eventDate.setHours(0, 0, 0, 0);
  const diffTime = eventDate.getTime() - today.getTime();
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}

// Get badge variant based on days until event
function getUrgencyBadge(daysUntil: number): { variant: "default" | "secondary" | "destructive" | "outline"; label: string } {
  if (daysUntil <= 0) {
    return { variant: "destructive", label: "Today" };
  } else if (daysUntil === 1) {
    return { variant: "destructive", label: "Tomorrow" };
  } else if (daysUntil <= 7) {
    return { variant: "default", label: `${daysUntil} days` };
  } else {
    return { variant: "secondary", label: `${daysUntil} days` };
  }
}

export function UpcomingEventsWidget({
  events,
  isLoading = false,
}: UpcomingEventsWidgetProps) {
  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Upcoming Events
          </CardTitle>
          <CardDescription className="hidden xl:inline">Events scheduled in the coming days</CardDescription>
          <CardDescription className="xl:hidden">Upcoming events</CardDescription>
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

  if (!events || events.length === 0) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Upcoming Events
          </CardTitle>
          <CardDescription className="hidden xl:inline">Events scheduled in the coming days</CardDescription>
          <CardDescription className="xl:hidden">Upcoming events</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center py-8">
          <div className="text-center">
            <Calendar className="mx-auto h-12 w-12 text-muted-foreground/50" />
            <p className="mt-2 text-sm text-muted-foreground">
              No upcoming events scheduled
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
          <Calendar className="h-5 w-5" />
          Upcoming Events
        </CardTitle>
        <CardDescription>Events scheduled in the coming days</CardDescription>
      </CardHeader>
      <CardContent className="flex-1">
        <div className="space-y-3">
          {events.slice(0, 5).map((event) => {
            const daysUntil = getDaysUntil(event.date);
            const urgency = getUrgencyBadge(daysUntil);

            return (
              <div
                key={event.id}
                className="flex items-start gap-3 rounded-lg border p-3 transition-colors hover:bg-muted/50"
              >
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h4 className="font-medium text-sm truncate">
                      {event.title}
                    </h4>
                    <Badge variant={urgency.variant} className="shrink-0 text-xs">
                      {urgency.label}
                    </Badge>
                  </div>
                  <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {formatEventDate(event.date)}
                    </span>
                    {event.venue && (
                      <span className="flex items-center gap-1">
                        <MapPin className="h-3 w-3" />
                        {event.venue}
                        {event.district && `, ${event.district}`}
                      </span>
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
          <Link href="/admin/events">
            View All Events
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

export default UpcomingEventsWidget;
