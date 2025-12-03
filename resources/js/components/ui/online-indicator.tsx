"use client";

import * as React from "react";
import { cn } from "@/lib/utils";
import { usePresence } from "@/contexts/PresenceContext";

interface OnlineIndicatorProps {
  userId: number;
  className?: string;
  size?: "sm" | "md" | "lg";
  /** @deprecated Use alwaysShow instead */
  showOffline?: boolean;
  /** Always show indicator (default: true). Set to false to hide when offline */
  alwaysShow?: boolean;
}

const sizeClasses = {
  sm: "h-2 w-2",
  md: "h-2.5 w-2.5",
  lg: "h-3 w-3",
};

export function OnlineIndicator({
  userId,
  className,
  size = "md",
  showOffline,
  alwaysShow = true,
}: OnlineIndicatorProps) {
  const { isUserOnline } = usePresence();
  const isOnline = isUserOnline(userId);

  // Support legacy showOffline prop
  const shouldShow = alwaysShow || showOffline || isOnline;

  if (!shouldShow) {
    return null;
  }

  return (
    <span
      className={cn(
        "rounded-full shrink-0",
        sizeClasses[size],
        isOnline ? "bg-green-500" : "bg-red-400",
        className
      )}
      title={isOnline ? "Online" : "Offline"}
    />
  );
}

// For use with Avatar components - positioned in bottom-right corner
interface AvatarOnlineIndicatorProps extends OnlineIndicatorProps {
  avatarSize?: "sm" | "md" | "lg";
}

export function AvatarOnlineIndicator({
  userId,
  avatarSize = "md",
  showOffline,
  alwaysShow = true,
  ...props
}: AvatarOnlineIndicatorProps) {
  const { isUserOnline } = usePresence();
  const isOnline = isUserOnline(userId);

  // Support legacy showOffline prop
  const shouldShow = alwaysShow || showOffline || isOnline;

  if (!shouldShow) {
    return null;
  }

  const positionClasses = {
    sm: "bottom-0 right-0",
    md: "bottom-0 right-0",
    lg: "-bottom-0.5 -right-0.5",
  };

  return (
    <span
      className={cn(
        "absolute rounded-full border-2 border-background",
        sizeClasses[avatarSize === "sm" ? "sm" : "md"],
        positionClasses[avatarSize],
        isOnline ? "bg-green-500" : "bg-red-400"
      )}
      title={isOnline ? "Online" : "Offline"}
    />
  );
}
