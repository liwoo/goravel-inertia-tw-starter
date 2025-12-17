"use client";

import * as React from "react";
import { Calendar as CalendarIcon, MapPin, Clock } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Calendar } from "@/components/ui/calendar";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";

// Interface for calendar event data
export interface CalendarEvent {
  id: number;
  title: string;
  date: string;
  endDate?: string;
  venue: string;
  district?: string;
}

interface EventsCalendarWidgetProps {
  events: CalendarEvent[];
  district?: string;
  isLoading?: boolean;
}

// Format date for display
function formatEventDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
  });
}

// Format time for display
function formatEventTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString("en-US", {
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  });
}

// Check if two dates are the same day
function isSameDay(date1: Date, date2: Date): boolean {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
}

// Check if a date falls within a range
function isDateInRange(date: Date, startDate: Date, endDate?: Date): boolean {
  const d = new Date(date.getFullYear(), date.getMonth(), date.getDate());
  const start = new Date(startDate.getFullYear(), startDate.getMonth(), startDate.getDate());

  if (!endDate) {
    return isSameDay(d, start);
  }

  const end = new Date(endDate.getFullYear(), endDate.getMonth(), endDate.getDate());
  return d >= start && d <= end;
}

// Get days until event
function getDaysUntil(dateString: string): number {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const eventDate = new Date(dateString);
  eventDate.setHours(0, 0, 0, 0);
  const diffTime = eventDate.getTime() - today.getTime();
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}

// Get urgency badge
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

export function EventsCalendarWidget({
  events,
  district,
  isLoading = false,
}: EventsCalendarWidgetProps) {
  const [selectedDate, setSelectedDate] = React.useState<Date | undefined>(undefined);

  // Get dates that have events
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

  // Get events for selected date
  const eventsForSelectedDate = React.useMemo(() => {
    if (!selectedDate) return [];

    return events.filter((event) => {
      const eventDate = new Date(event.date);
      const endDate = event.endDate ? new Date(event.endDate) : undefined;
      return isDateInRange(selectedDate, eventDate, endDate);
    });
  }, [selectedDate, events]);

  // Get upcoming events (next 5)
  const upcomingEvents = React.useMemo(() => {
    const now = new Date();
    return events
      .filter((event) => new Date(event.date) >= now)
      .slice(0, 5);
  }, [events]);

  // Display events - selected date or upcoming
  const displayEvents = selectedDate ? eventsForSelectedDate : upcomingEvents;

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="pb-2 px-4">
          <CardTitle className="flex items-center gap-2 text-base">
            <CalendarIcon className="h-4 w-4" />
            Events
          </CardTitle>
          <CardDescription className="text-xs">Loading...</CardDescription>
        </CardHeader>
        <CardContent className="px-4 pb-4">
          <div className="h-48 animate-pulse rounded-md bg-muted" />
        </CardContent>
      </Card>
    );
  }

  // Modifiers for highlighting event days
  const modifiers = {
    hasEvent: eventDates,
  };

  const modifiersClassNames = {
    hasEvent: "bg-primary/20 font-semibold text-primary hover:bg-primary/30",
  };

  return (
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

        {/* Compact Events List */}
        <div className="flex flex-col gap-1.5">
          <div className="flex items-center justify-between">
            <h4 className="text-xs font-medium text-muted-foreground">
              {selectedDate
                ? formatEventDate(selectedDate.toISOString())
                : "Upcoming"
              }
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
                {displayEvents.map((event) => {
                  const daysUntil = getDaysUntil(event.date);
                  const urgency = getUrgencyBadge(daysUntil);

                  return (
                    <div
                      key={event.id}
                      className="flex flex-col gap-0.5 rounded-md border p-2 transition-colors hover:bg-muted/50"
                    >
                      <div className="flex items-start justify-between gap-1">
                        <h5 className="font-medium text-xs leading-tight line-clamp-1">
                          {event.title}
                        </h5>
                        <Badge variant={urgency.variant} className="shrink-0 text-[10px] px-1.5 py-0">
                          {urgency.label}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-2 text-[10px] text-muted-foreground">
                        <span className="flex items-center gap-0.5">
                          <Clock className="h-2.5 w-2.5" />
                          {formatEventTime(event.date)}
                        </span>
                        {event.venue && (
                          <span className="flex items-center gap-0.5 truncate">
                            <MapPin className="h-2.5 w-2.5 shrink-0" />
                            <span className="truncate">{event.venue}</span>
                          </span>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </ScrollArea>
        </div>
      </CardContent>
    </Card>
  );
}

export default EventsCalendarWidget;
