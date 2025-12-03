"use client";

import * as React from "react";
import sseManager, { SSEEvent } from "@/services/sseManager";

interface PresenceContextType {
  onlineUsers: Set<number>;
  isUserOnline: (userId: number) => boolean;
  onlineCount: number;
  refreshPresence: () => Promise<void>;
  isRefreshing: boolean;
}

const PresenceContext = React.createContext<PresenceContextType | undefined>(undefined);

export function PresenceProvider({ children }: { children: React.ReactNode }) {
  const [onlineUsers, setOnlineUsers] = React.useState<Set<number>>(new Set());
  const [isRefreshing, setIsRefreshing] = React.useState(false);

  // Function to fetch presence from REST API
  const fetchPresence = React.useCallback(async () => {
    try {
      const res = await fetch('/api/presence/online', {
        credentials: 'include',
      });
      const data = await res.json();
      if (data.success && Array.isArray(data.data?.online_users)) {
        setOnlineUsers(new Set(data.data.online_users));
        return data.data.online_users;
      }
    } catch (err) {
      console.warn('Presence: failed to fetch state', err);
    }
    return [];
  }, []);

  // Manual refresh function exposed to consumers
  const refreshPresence = React.useCallback(async () => {
    setIsRefreshing(true);
    try {
      await fetchPresence();
    } finally {
      setIsRefreshing(false);
    }
  }, [fetchPresence]);

  React.useEffect(() => {
    // Handle initial presence state from server
    const handlePresenceInitial = (event: SSEEvent) => {
      console.log('Presence: received initial state', event.data);
      const onlineUsersList = event.data?.online_users;
      if (Array.isArray(onlineUsersList)) {
        setOnlineUsers(new Set(onlineUsersList));
        console.log('Presence: set online users', onlineUsersList);
      }
    };

    // Handle presence changes (user came online or went offline)
    const handlePresenceChange = (event: SSEEvent) => {
      console.log('Presence: received change', event.data);
      const userId = event.data?.user_id;
      const isOnline = event.data?.is_online;
      if (typeof userId === 'number') {
        setOnlineUsers(prev => {
          const newSet = new Set(prev);
          if (isOnline) {
            newSet.add(userId);
          } else {
            newSet.delete(userId);
          }
          console.log('Presence: updated online users', Array.from(newSet));
          return newSet;
        });
      }
    };

    // Subscribe to presence events
    const unsubscribeInitial = sseManager.on('presence:initial', handlePresenceInitial);
    const unsubscribeChange = sseManager.on('presence:change', handlePresenceChange);

    // Fetch initial presence state via REST API in case we missed the SSE event
    fetchPresence();

    return () => {
      unsubscribeInitial();
      unsubscribeChange();
    };
  }, []);

  const isUserOnline = React.useCallback((userId: number): boolean => {
    return onlineUsers.has(userId);
  }, [onlineUsers]);

  const contextValue: PresenceContextType = {
    onlineUsers,
    isUserOnline,
    onlineCount: onlineUsers.size,
    refreshPresence,
    isRefreshing,
  };

  return (
    <PresenceContext.Provider value={contextValue}>
      {children}
    </PresenceContext.Provider>
  );
}

export function usePresence() {
  const context = React.useContext(PresenceContext);
  if (context === undefined) {
    throw new Error("usePresence must be used within a PresenceProvider");
  }
  return context;
}

// Hook for checking a single user's status
export function useUserOnlineStatus(userId: number): boolean {
  const { isUserOnline } = usePresence();
  return isUserOnline(userId);
}
