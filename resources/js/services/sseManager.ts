import { router } from "@inertiajs/react";

export interface SSEEvent {
  type: string;
  data: any;
  time: string;
}

export type SSEEventHandler = (event: SSEEvent) => void;

class SSEManager {
  private eventSource: EventSource | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private handlers: Map<string, Set<SSEEventHandler>> = new Map();
  private isConnecting = false;
  private authToken: string | null = null;

  constructor() {
    // Listen for auth changes
    if (typeof window !== 'undefined') {
      window.addEventListener('auth:login', this.handleAuthChange.bind(this));
      window.addEventListener('auth:logout', this.disconnect.bind(this));
      
      // Auto-connect if we have a token
      this.authToken = this.getAuthToken();
      if (this.authToken) {
        this.connect();
      }
    }
  }

  private getAuthToken(): string | null {
    // Try to get token from cookie
    const cookies = document.cookie.split(';');
    for (const cookie of cookies) {
      const [name, value] = cookie.trim().split('=');
      if (name === 'jwt_token') {
        return decodeURIComponent(value);
      }
    }
    
    // Try localStorage as fallback
    return localStorage.getItem('jwt_token');
  }

  private handleAuthChange(event: CustomEvent) {
    this.authToken = event.detail?.token || this.getAuthToken();
    if (this.authToken) {
      this.connect();
    }
  }

  connect() {
    if (this.isConnecting || this.eventSource?.readyState === EventSource.OPEN) {
      return;
    }

    if (!this.authToken) {
      console.warn('SSE: No auth token available');
      return;
    }

    this.isConnecting = true;

    try {
      // Create SSE connection with auth token as query parameter
      const url = new URL('/api/sse/stream', window.location.origin);
      url.searchParams.append('token', this.authToken);
      
      this.eventSource = new EventSource(url.toString());

      this.eventSource.onopen = () => {
        console.log('SSE: Connected');
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.emit('connection:open', {});
      };

      this.eventSource.onerror = (error) => {
        console.error('SSE: Connection error', error);
        this.isConnecting = false;
        this.emit('connection:error', { error });

        if (this.eventSource?.readyState === EventSource.CLOSED) {
          this.handleReconnect();
        }
      };

      // Handle standard SSE events
      this.eventSource.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.emit('message', data);
        } catch (error) {
          console.error('SSE: Failed to parse message', error);
        }
      };

      // Register specific event handlers
      this.registerEventHandlers();

    } catch (error) {
      console.error('SSE: Failed to create connection', error);
      this.isConnecting = false;
      this.handleReconnect();
    }
  }

  private registerEventHandlers() {
    if (!this.eventSource) return;

    // Connection events
    this.eventSource.addEventListener('connected', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('connected', data);
    });

    this.eventSource.addEventListener('heartbeat', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('heartbeat', data);
    });

    // Message events
    this.eventSource.addEventListener('message:new', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('message:new', data);
    });

    this.eventSource.addEventListener('message:read', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('message:read', data);
    });

    this.eventSource.addEventListener('message:updated', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('message:updated', data);
    });

    this.eventSource.addEventListener('message:deleted', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('message:deleted', data);
    });

    this.eventSource.addEventListener('message:unread_count', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('message:unread_count', data);
    });

    // Notification events
    this.eventSource.addEventListener('notification:new', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('notification:new', data);
    });

    this.eventSource.addEventListener('notification:read', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('notification:read', data);
    });

    this.eventSource.addEventListener('notification:dismissed', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('notification:dismissed', data);
    });

    this.eventSource.addEventListener('notification:counts', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('notification:counts', data);
    });

    this.eventSource.addEventListener('notification:system', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('notification:system', data);
    });

    this.eventSource.addEventListener('notification:initial', (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      this.emit('notification:initial', data);
    });
  }

  private handleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('SSE: Max reconnection attempts reached');
      this.emit('connection:failed', { attempts: this.reconnectAttempts });
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

    console.log(`SSE: Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
    
    setTimeout(() => {
      this.connect();
    }, delay);
  }

  disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
      this.isConnecting = false;
      this.reconnectAttempts = 0;
      this.emit('connection:closed', {});
      console.log('SSE: Disconnected');
    }
  }

  on(event: string, handler: SSEEventHandler) {
    if (!this.handlers.has(event)) {
      this.handlers.set(event, new Set());
    }
    this.handlers.get(event)!.add(handler);

    // Return unsubscribe function
    return () => {
      this.off(event, handler);
    };
  }

  off(event: string, handler: SSEEventHandler) {
    const eventHandlers = this.handlers.get(event);
    if (eventHandlers) {
      eventHandlers.delete(handler);
      if (eventHandlers.size === 0) {
        this.handlers.delete(event);
      }
    }
  }

  once(event: string, handler: SSEEventHandler) {
    const wrappedHandler: SSEEventHandler = (data) => {
      handler(data);
      this.off(event, wrappedHandler);
    };
    this.on(event, wrappedHandler);
  }

  private emit(event: string, data: any) {
    const eventHandlers = this.handlers.get(event);
    if (eventHandlers) {
      const sseEvent: SSEEvent = {
        type: event,
        data,
        time: new Date().toISOString(),
      };
      
      eventHandlers.forEach(handler => {
        try {
          handler(sseEvent);
        } catch (error) {
          console.error(`SSE: Error in event handler for '${event}'`, error);
        }
      });
    }
  }

  isConnected(): boolean {
    return this.eventSource?.readyState === EventSource.OPEN;
  }

  getConnectionState(): string {
    if (!this.eventSource) return 'disconnected';
    
    switch (this.eventSource.readyState) {
      case EventSource.CONNECTING:
        return 'connecting';
      case EventSource.OPEN:
        return 'connected';
      case EventSource.CLOSED:
        return 'closed';
      default:
        return 'unknown';
    }
  }
}

// Create singleton instance
const sseManager = new SSEManager();

// Export for use in other modules
export default sseManager;