import React, { useState, useEffect, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Activity,
  Building2,
  Calendar,
  FileText,
  User,
  ChevronLeft,
  ChevronRight,
  RefreshCw,
  AlertCircle
} from 'lucide-react';
import { toast } from 'sonner';
import axios from '@/lib/axios';
import { cn } from '@/lib/utils';

interface ActivitySectionProps {
  userId?: number;
}

// Recent activity data structure (matches dashboard RecentActivityDTO)
interface RecentActivity {
  id: number;
  entityType: "sme" | "event" | "procurement";
  entityId: number;
  entityName: string;
  action: "created" | "updated";
  userId: number;
  userName: string;
  timestamp: string;
}

// Get icon for entity type
function getEntityIcon(entityType: string) {
  switch (entityType) {
    case "sme":
      return Building2;
    case "event":
      return Calendar;
    case "procurement":
      return FileText;
    default:
      return Activity;
  }
}

// Get label for entity type
function getEntityLabel(entityType: string): string {
  switch (entityType) {
    case "sme":
      return "MSME";
    case "event":
      return "Event";
    case "procurement":
      return "Procurement";
    default:
      return entityType;
  }
}

// Format timestamp to relative time
function formatRelativeTime(timestamp: string): string {
  const date = new Date(timestamp);
  if (isNaN(date.getTime())) return "Recently";

  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) return "Just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
}

const ITEMS_PER_PAGE = 5;

export const ActivitySection: React.FC<ActivitySectionProps> = () => {
  const [activities, setActivities] = useState<RecentActivity[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(0);

  const totalPages = Math.ceil((activities?.length || 0) / ITEMS_PER_PAGE);
  const startIndex = currentPage * ITEMS_PER_PAGE;
  const paginatedActivities = activities?.slice(startIndex, startIndex + ITEMS_PER_PAGE) || [];

  const fetchActivities = useCallback(async (showRefreshIndicator = false) => {
    if (showRefreshIndicator) {
      setIsRefreshing(true);
    } else {
      setIsLoading(true);
    }
    setError(null);

    try {
      const response = await axios.get<{
        success: boolean;
        data: RecentActivity[];
      }>('/api/account/recent-activities', {
        params: { limit: 50 }
      });

      setActivities(response.data.data || []);
    } catch (err) {
      console.error('Failed to fetch activities:', err);
      setError('Failed to load activity history. Please try again.');
      toast.error('Failed to load activities');
    } finally {
      setIsLoading(false);
      setIsRefreshing(false);
    }
  }, []);

  useEffect(() => {
    fetchActivities();
  }, [fetchActivities]);

  const handleRefresh = () => {
    fetchActivities(true);
  };

  const goToPreviousPage = () => {
    setCurrentPage((prev) => Math.max(0, prev - 1));
  };

  const goToNextPage = () => {
    setCurrentPage((prev) => Math.min(totalPages - 1, prev + 1));
  };

  // Loading skeleton
  if (isLoading) {
    return (
      <Card className="flex flex-col h-full">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Recent Activity
          </CardTitle>
          <CardDescription>Your recent changes across the system</CardDescription>
        </CardHeader>
        <CardContent className="flex-1">
          <div className="space-y-3">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="h-12 animate-pulse rounded-md bg-muted" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  // Error state
  if (error) {
    return (
      <Card className="flex flex-col h-full">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Recent Activity
          </CardTitle>
          <CardDescription>Your recent changes across the system</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center py-8">
          <div className="text-center">
            <AlertCircle className="mx-auto h-12 w-12 text-muted-foreground/50" />
            <p className="mt-2 text-sm text-muted-foreground mb-4">{error}</p>
            <Button variant="outline" size="sm" onClick={handleRefresh}>
              <RefreshCw className="h-4 w-4 mr-2" />
              Try Again
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Empty state
  if (!activities || activities.length === 0) {
    return (
      <Card className="flex flex-col h-full">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Activity className="h-5 w-5" />
                Recent Activity
              </CardTitle>
              <CardDescription>Your recent changes across the system</CardDescription>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={handleRefresh}
              disabled={isRefreshing}
              title="Refresh"
            >
              <RefreshCw className={cn("h-4 w-4", isRefreshing && "animate-spin")} />
            </Button>
          </div>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center py-8">
          <div className="text-center">
            <Activity className="mx-auto h-12 w-12 text-muted-foreground/50" />
            <p className="mt-2 text-sm text-muted-foreground">
              No recent activity
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="flex flex-col h-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Recent Activity
            </CardTitle>
            <CardDescription>Your recent changes across the system</CardDescription>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleRefresh}
            disabled={isRefreshing}
            title="Refresh"
          >
            <RefreshCw className={cn("h-4 w-4", isRefreshing && "animate-spin")} />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="flex-1">
        <div className="space-y-3">
          {paginatedActivities.map((activity) => {
            const Icon = getEntityIcon(activity.entityType);

            return (
              <div
                key={`${activity.entityType}-${activity.entityId}-${activity.timestamp}`}
                className="flex items-start gap-3 text-sm"
              >
                <div
                  className={cn(
                    "mt-0.5 flex h-7 w-7 shrink-0 items-center justify-center rounded-full",
                    activity.action === "created"
                      ? "bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400"
                      : "bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400"
                  )}
                >
                  <Icon className="h-3.5 w-3.5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium truncate">
                      {activity.entityName}
                    </span>
                    <Badge
                      variant={activity.action === "created" ? "default" : "secondary"}
                      className="text-xs shrink-0"
                    >
                      {activity.action}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-2 text-xs text-muted-foreground mt-0.5">
                    <Badge variant="outline" className="text-xs px-1 h-4">
                      {getEntityLabel(activity.entityType)}
                    </Badge>
                    <span>-</span>
                    <span>{formatRelativeTime(activity.timestamp)}</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
      {totalPages > 1 && (
        <CardFooter className="border-t pt-4">
          <div className="flex w-full items-center justify-between">
            <Button
              variant="ghost"
              size="sm"
              onClick={goToPreviousPage}
              disabled={currentPage === 0}
            >
              <ChevronLeft className="h-4 w-4 mr-1" />
              Prev
            </Button>
            <span className="text-xs text-muted-foreground">
              {currentPage + 1} / {totalPages}
            </span>
            <Button
              variant="ghost"
              size="sm"
              onClick={goToNextPage}
              disabled={currentPage === totalPages - 1}
            >
              Next
              <ChevronRight className="h-4 w-4 ml-1" />
            </Button>
          </div>
        </CardFooter>
      )}
    </Card>
  );
};

export default ActivitySection;
