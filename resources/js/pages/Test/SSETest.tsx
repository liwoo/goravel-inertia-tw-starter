import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useSSE, useSSEConnection, useSSEEvent } from '@/hooks/useSSE';
import sseManager from '@/services/sseManager';
import axios from '@/lib/axios';

interface SSELogEntry {
  id: string;
  timestamp: Date;
  event: string;
  data: any;
}

export function SSETest() {
  const { isConnected, connectionState, connect, disconnect } = useSSEConnection();
  const [logs, setLogs] = useState<SSELogEntry[]>([]);
  const [testUserId, setTestUserId] = useState<string>('2');

  // Listen to all events and log them
  const logEvent = (eventName: string) => (event: any) => {
    const logEntry: SSELogEntry = {
      id: `${Date.now()}-${Math.random()}`,
      timestamp: new Date(),
      event: eventName,
      data: event.data
    };
    setLogs(prev => [logEntry, ...prev].slice(0, 100)); // Keep last 100 entries
  };

  // Subscribe to various events
  useSSE('connected', logEvent('connected'));
  useSSE('heartbeat', logEvent('heartbeat'));
  useSSE('message:new', logEvent('message:new'));
  useSSE('message:read', logEvent('message:read'));
  useSSE('message:updated', logEvent('message:updated'));
  useSSE('message:deleted', logEvent('message:deleted'));
  useSSE('message:unread_count', logEvent('message:unread_count'));
  useSSE('notification:new', logEvent('notification:new'));
  useSSE('notification:read', logEvent('notification:read'));
  useSSE('notification:dismissed', logEvent('notification:dismissed'));
  useSSE('notification:counts', logEvent('notification:counts'));
  useSSE('notification:system', logEvent('notification:system'));

  // Test actions
  const sendTestMessage = async () => {
    try {
      await axios.post('/api/messages', {
        recipient_id: parseInt(testUserId),
        content: `Test message sent at ${new Date().toLocaleTimeString()}`,
        type: 'direct'
      });
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  };

  const createTestNotification = async () => {
    try {
      await axios.post('/api/notifications', {
        user_id: parseInt(testUserId),
        title: 'Test Notification',
        message: `This is a test notification created at ${new Date().toLocaleTimeString()}`,
        type: 'system',
        priority: 'normal'
      });
    } catch (error) {
      console.error('Failed to create notification:', error);
    }
  };

  const clearLogs = () => {
    setLogs([]);
  };

  const connectionColor = {
    connected: 'bg-green-500',
    connecting: 'bg-yellow-500',
    closed: 'bg-red-500',
    disconnected: 'bg-gray-500',
    error: 'bg-red-500',
    unknown: 'bg-gray-500'
  }[connectionState] || 'bg-gray-500';

  return (
    <div className="container mx-auto p-6 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            SSE Test Dashboard
            <div className="flex items-center gap-2">
              <div className={`w-3 h-3 rounded-full ${connectionColor}`} />
              <span className="text-sm font-normal">{connectionState}</span>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2">
            <Button
              onClick={isConnected ? disconnect : connect}
              variant={isConnected ? 'destructive' : 'default'}
            >
              {isConnected ? 'Disconnect' : 'Connect'}
            </Button>
            <Button onClick={clearLogs} variant="outline">
              Clear Logs
            </Button>
          </div>

          <div className="space-y-2">
            <h3 className="text-sm font-medium">Test Actions</h3>
            <div className="flex gap-2 items-end">
              <div className="flex-1">
                <label className="text-xs text-muted-foreground">Target User ID</label>
                <input
                  type="number"
                  value={testUserId}
                  onChange={(e) => setTestUserId(e.target.value)}
                  className="w-full px-3 py-1 border rounded-md"
                  placeholder="Enter user ID"
                />
              </div>
              <Button onClick={sendTestMessage} size="sm">
                Send Test Message
              </Button>
              <Button onClick={createTestNotification} size="sm" variant="secondary">
                Create Test Notification
              </Button>
            </div>
          </div>

          <div className="space-y-2">
            <h3 className="text-sm font-medium">
              Event Log ({logs.length} events)
            </h3>
            <ScrollArea className="h-96 w-full border rounded-md">
              <div className="p-4 space-y-2">
                {logs.length === 0 ? (
                  <p className="text-center text-muted-foreground py-8">
                    No events logged yet. Connect to SSE to start receiving events.
                  </p>
                ) : (
                  logs.map((log) => (
                    <div
                      key={log.id}
                      className="border rounded-md p-3 space-y-1 text-sm"
                    >
                      <div className="flex items-center justify-between">
                        <Badge variant="outline">{log.event}</Badge>
                        <span className="text-xs text-muted-foreground">
                          {log.timestamp.toLocaleTimeString()}
                        </span>
                      </div>
                      <pre className="text-xs bg-muted p-2 rounded overflow-x-auto">
                        {JSON.stringify(log.data, null, 2)}
                      </pre>
                    </div>
                  ))
                )}
              </div>
            </ScrollArea>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}