import { useEffect, useState, useCallback, useRef } from 'react';
import sseManager, { SSEEvent, SSEEventHandler } from '../services/sseManager';

export function useSSE(event: string | string[], handler: SSEEventHandler) {
  const events = Array.isArray(event) ? event : [event];
  const unsubscribers = useRef<Array<() => void>>([]);

  useEffect(() => {
    // Subscribe to all events
    events.forEach(eventName => {
      const unsubscribe = sseManager.on(eventName, handler);
      unsubscribers.current.push(unsubscribe);
    });

    // Cleanup on unmount
    return () => {
      unsubscribers.current.forEach(unsubscribe => unsubscribe());
      unsubscribers.current = [];
    };
  }, [events.join(','), handler]);
}

export function useSSEConnection() {
  const [isConnected, setIsConnected] = useState(sseManager.isConnected());
  const [connectionState, setConnectionState] = useState(sseManager.getConnectionState());

  useEffect(() => {
    const handleConnectionOpen = () => {
      setIsConnected(true);
      setConnectionState('connected');
    };

    const handleConnectionError = () => {
      setIsConnected(false);
      setConnectionState('error');
    };

    const handleConnectionClosed = () => {
      setIsConnected(false);
      setConnectionState('closed');
    };

    const unsubscribeOpen = sseManager.on('connection:open', handleConnectionOpen);
    const unsubscribeError = sseManager.on('connection:error', handleConnectionError);
    const unsubscribeClosed = sseManager.on('connection:closed', handleConnectionClosed);

    // Check initial state
    setIsConnected(sseManager.isConnected());
    setConnectionState(sseManager.getConnectionState());

    return () => {
      unsubscribeOpen();
      unsubscribeError();
      unsubscribeClosed();
    };
  }, []);

  const connect = useCallback(() => {
    sseManager.connect();
  }, []);

  const disconnect = useCallback(() => {
    sseManager.disconnect();
  }, []);

  return {
    isConnected,
    connectionState,
    connect,
    disconnect,
  };
}

export function useSSEEvent<T = any>(event: string, initialValue?: T) {
  const [data, setData] = useState<T | undefined>(initialValue);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);

  useSSE(event, (sseEvent: SSEEvent) => {
    setData(sseEvent.data as T);
    setLastUpdate(new Date(sseEvent.time));
  });

  return { data, lastUpdate };
}

export function useMessageUnreadCount() {
  const { data: count = 0 } = useSSEEvent<{ count: number }>('message:unread_count');
  return count;
}

export function useNotificationCounts() {
  const { data: counts = {} } = useSSEEvent<Record<string, number>>('notification:counts');
  return counts;
}