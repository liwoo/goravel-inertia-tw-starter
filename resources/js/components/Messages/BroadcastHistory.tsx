"use client";

import * as React from "react";
import { Megaphone, Users, Check, Clock, Loader2 } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import axios from "@/lib/axios";

interface Broadcast {
  id: number;
  content: string;
  created_at: string;
  recipient_count: number;
  read_count: number;
}

interface BroadcastHistoryProps {
  className?: string;
  refreshTrigger?: number;
}

export function BroadcastHistory({ className, refreshTrigger }: BroadcastHistoryProps) {
  const [broadcasts, setBroadcasts] = React.useState<Broadcast[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    loadBroadcasts();
  }, [refreshTrigger]);

  const loadBroadcasts = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await axios.get("/api/messages/broadcast-history?pageSize=50");
      if (response.data.success && response.data.data?.data) {
        setBroadcasts(response.data.data.data);
      }
    } catch (err: any) {
      console.error("Failed to load broadcast history:", err);
      setError(err.response?.data?.message || "Failed to load broadcast history");
    } finally {
      setLoading(false);
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60);

    if (diffInHours < 1) {
      return "Just now";
    } else if (diffInHours < 24) {
      return `${Math.floor(diffInHours)}h ago`;
    } else if (diffInHours < 168) {
      return `${Math.floor(diffInHours / 24)}d ago`;
    } else {
      return date.toLocaleDateString() + " " + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
  };

  if (loading) {
    return (
      <div className={cn("flex flex-col h-full items-center justify-center", className)}>
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        <p className="mt-2 text-sm text-muted-foreground">Loading broadcast history...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={cn("flex flex-col h-full items-center justify-center p-4", className)}>
        <p className="text-sm text-destructive">{error}</p>
        <Button variant="outline" size="sm" className="mt-2" onClick={loadBroadcasts}>
          Retry
        </Button>
      </div>
    );
  }

  if (broadcasts.length === 0) {
    return (
      <div className={cn("flex flex-col h-full items-center justify-center p-4", className)}>
        <Megaphone className="h-12 w-12 text-muted-foreground/50 mb-4" />
        <h3 className="font-medium text-lg mb-1">No Broadcasts Yet</h3>
        <p className="text-sm text-muted-foreground text-center max-w-sm">
          Broadcasts you send to roles will appear here. Use the Broadcast tab on the left to send your first message.
        </p>
      </div>
    );
  }

  return (
    <div className={cn("flex flex-col h-full", className)}>
      <div className="p-4 border-b">
        <div className="flex items-center gap-2">
          <Megaphone className="h-5 w-5" />
          <h2 className="font-semibold">Broadcast History</h2>
          <Badge variant="secondary" className="ml-auto">
            {broadcasts.length} broadcast{broadcasts.length !== 1 ? "s" : ""}
          </Badge>
        </div>
      </div>

      <ScrollArea className="flex-1">
        <div className="p-4 space-y-3">
          {broadcasts.map((broadcast) => (
            <div
              key={broadcast.id}
              className="border rounded-lg p-4"
            >
              <p className="text-sm whitespace-pre-wrap break-words">
                {broadcast.content}
              </p>

              <div className="flex items-center gap-4 mt-3 text-xs text-muted-foreground">
                <span className="flex items-center gap-1">
                  <Clock className="h-3 w-3" />
                  {formatTime(broadcast.created_at)}
                </span>
                <span className="flex items-center gap-1">
                  <Users className="h-3 w-3" />
                  {broadcast.recipient_count} recipient{broadcast.recipient_count !== 1 ? "s" : ""}
                </span>
                <span className="flex items-center gap-1">
                  <Check className="h-3 w-3" />
                  {broadcast.read_count} read
                </span>
              </div>
            </div>
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}
