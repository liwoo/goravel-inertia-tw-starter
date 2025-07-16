"use client";

import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import axios from '@/lib/axios';

interface User {
  id: number;
  name: string;
  email: string;
}

interface Notification {
  id: number;
  title: string;
  message: string;
  type: string;
  is_read: boolean;
  is_dismissed: boolean;
  priority: string;
  created_at: string;
  read_at?: string;
  dismissed_at?: string;
  trigger_user?: User;
  related_type?: string;
  related_id?: number;
  data?: string;
  expires_at?: string;
}

interface ListRequest {
  page: number;
  pageSize: number;
  search?: string;
  sort?: string;
  direction?: string;
  filters?: Record<string, any>;
}

interface PaginatedResult<T> {
  data: T[];
  total: number;
  currentPage: number;
  lastPage: number;
  perPage: number;
  from: number;
  to: number;
  hasNext: boolean;
  hasPrev: boolean;
}

interface NotificationCounts {
  unread: number;
  unread_high: number;
  unread_messages: number;
  unread_mentions: number;
  total: number;
}

interface NotificationContextType {
  // State
  notifications: Notification[];
  unreadCount: number;
  counts: NotificationCounts | null;
  loading: boolean;
  error: string | null;

  // Actions
  loadNotifications: (request: ListRequest) => Promise<void>;
  loadNotificationCounts: () => Promise<void>;
  loadNotificationsByType: (type: string, limit?: number) => Promise<Notification[]>;
  markAsRead: (notificationId: number) => Promise<void>;
  markAllAsRead: () => Promise<void>;
  dismissNotification: (notificationId: number) => Promise<void>;
  dismissAllNotifications: () => Promise<void>;
  batchMarkAsRead: (notificationIds: number[]) => Promise<void>;
  batchDismiss: (notificationIds: number[]) => Promise<void>;
  createNotification: (data: {
    user_id: number;
    title: string;
    message?: string;
    type: string;
    priority?: string;
    expires_in_days?: number;
    data?: string;
  }) => Promise<Notification>;
  createSystemNotification: (data: {
    title: string;
    message: string;
    priority?: string;
    expires_in_days?: number;
    role_filter?: string[];
  }) => Promise<void>;
  refreshCounts: () => Promise<void>;
  clearError: () => void;
}

const NotificationContext = createContext<NotificationContextType | null>(null);

export function NotificationProvider({ children }: { children: React.ReactNode }) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [counts, setCounts] = useState<NotificationCounts | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load notifications
  const loadNotifications = useCallback(async (request: ListRequest) => {
    try {
      setLoading(true);
      setError(null);
      
      const params = new URLSearchParams({
        page: request.page.toString(),
        pageSize: request.pageSize.toString(),
      });

      if (request.search) params.append('search', request.search);
      if (request.sort) params.append('sort', request.sort);
      if (request.direction) params.append('direction', request.direction);
      
      // Add filters
      if (request.filters) {
        Object.entries(request.filters).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== '') {
            params.append(key, value.toString());
          }
        });
      }

      const response = await axios.get(`/api/notifications?${params}`);
      
      if (response.data.success && response.data.data) {
        const paginatedData = response.data.data as PaginatedResult<Notification>;
        setNotifications(paginatedData.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load notifications');
      console.error('Error loading notifications:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load notification counts
  const loadNotificationCounts = useCallback(async () => {
    try {
      setError(null);
      
      const response = await axios.get('/api/notifications/counts');
      
      if (response.data.success && response.data.data) {
        const countsData = response.data.data as NotificationCounts;
        setCounts(countsData);
        setUnreadCount(countsData.unread);
      }
    } catch (err: any) {
      console.error('Error loading notification counts:', err);
    }
  }, []);

  // Load notifications by type
  const loadNotificationsByType = useCallback(async (type: string, limit: number = 20): Promise<Notification[]> => {
    try {
      setError(null);
      
      const params = new URLSearchParams();
      if (limit > 0) params.append('limit', limit.toString());
      
      const response = await axios.get(`/api/notifications/type/${type}?${params}`);
      
      if (response.data.success && response.data.data) {
        return response.data.data as Notification[];
      }
      
      return [];
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to load notifications';
      setError(errorMessage);
      console.error('Error loading notifications by type:', err);
      return [];
    }
  }, []);

  // Mark notification as read
  const markAsRead = useCallback(async (notificationId: number) => {
    try {
      setError(null);
      
      const response = await axios.put(`/api/notifications/${notificationId}/read`);
      
      if (response.data.success) {
        // Update notification in state
        setNotifications(prev =>
          prev.map(notification =>
            notification.id === notificationId
              ? { ...notification, is_read: true, read_at: new Date().toISOString() }
              : notification
          )
        );
        
        // Refresh counts
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to mark notification as read';
      setError(errorMessage);
      console.error('Error marking notification as read:', err);
    }
  }, [loadNotificationCounts]);

  // Mark all notifications as read
  const markAllAsRead = useCallback(async () => {
    try {
      setError(null);
      
      const response = await axios.put('/api/notifications/read-all');
      
      if (response.data.success) {
        // Update all notifications in state
        setNotifications(prev =>
          prev.map(notification => ({
            ...notification,
            is_read: true,
            read_at: new Date().toISOString()
          }))
        );
        
        // Refresh counts
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to mark all notifications as read';
      setError(errorMessage);
      console.error('Error marking all notifications as read:', err);
    }
  }, [loadNotificationCounts]);

  // Dismiss notification
  const dismissNotification = useCallback(async (notificationId: number) => {
    try {
      setError(null);
      
      const response = await axios.delete(`/api/notifications/${notificationId}`);
      
      if (response.data.success) {
        // Remove notification from state
        setNotifications(prev =>
          prev.filter(notification => notification.id !== notificationId)
        );
        
        // Refresh counts
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to dismiss notification';
      setError(errorMessage);
      console.error('Error dismissing notification:', err);
    }
  }, [loadNotificationCounts]);

  // Dismiss all notifications
  const dismissAllNotifications = useCallback(async () => {
    try {
      setError(null);
      
      const response = await axios.delete('/api/notifications');
      
      if (response.data.success) {
        // Clear all notifications from state
        setNotifications([]);
        
        // Refresh counts
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to dismiss all notifications';
      setError(errorMessage);
      console.error('Error dismissing all notifications:', err);
    }
  }, [loadNotificationCounts]);

  // Batch mark as read
  const batchMarkAsRead = useCallback(async (notificationIds: number[]) => {
    try {
      setError(null);
      
      const response = await axios.put('/api/notifications/batch/read', {
        notification_ids: notificationIds
      });
      
      if (response.data.success) {
        // Update notifications in state
        setNotifications(prev =>
          prev.map(notification =>
            notificationIds.includes(notification.id)
              ? { ...notification, is_read: true, read_at: new Date().toISOString() }
              : notification
          )
        );
        
        // Refresh counts
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to mark notifications as read';
      setError(errorMessage);
      console.error('Error batch marking notifications as read:', err);
    }
  }, [loadNotificationCounts]);

  // Batch dismiss
  const batchDismiss = useCallback(async (notificationIds: number[]) => {
    try {
      setError(null);
      
      const response = await axios.delete('/api/notifications/batch', {
        data: { notification_ids: notificationIds }
      });
      
      if (response.data.success) {
        // Remove notifications from state
        setNotifications(prev =>
          prev.filter(notification => !notificationIds.includes(notification.id))
        );
        
        // Refresh counts
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to dismiss notifications';
      setError(errorMessage);
      console.error('Error batch dismissing notifications:', err);
    }
  }, [loadNotificationCounts]);

  // Create notification (admin only)
  const createNotification = useCallback(async (data: {
    user_id: number;
    title: string;
    message?: string;
    type: string;
    priority?: string;
    expires_in_days?: number;
    data?: string;
  }): Promise<Notification> => {
    try {
      setError(null);
      
      const response = await axios.post('/api/notifications', data);
      
      if (response.data.success && response.data.data) {
        const newNotification = response.data.data as Notification;
        
        // Add to notifications if it's for current user
        // (This would need user context to determine)
        
        return newNotification;
      }
      
      throw new Error('Failed to create notification');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to create notification';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, []);

  // Create system notification (super admin only)
  const createSystemNotification = useCallback(async (data: {
    title: string;
    message: string;
    priority?: string;
    expires_in_days?: number;
    role_filter?: string[];
  }) => {
    try {
      setError(null);
      
      const response = await axios.post('/api/notifications/system', data);
      
      if (response.data.success) {
        // Refresh notifications and counts
        await loadNotifications({ page: 1, pageSize: 20 });
        await loadNotificationCounts();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to create system notification';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [loadNotifications, loadNotificationCounts]);

  // Refresh counts
  const refreshCounts = useCallback(async () => {
    await loadNotificationCounts();
  }, [loadNotificationCounts]);

  // Clear error
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Load initial data
  useEffect(() => {
    loadNotificationCounts();
  }, [loadNotificationCounts]);

  // Refresh counts periodically
  useEffect(() => {
    const interval = setInterval(loadNotificationCounts, 30000); // Every 30 seconds
    return () => clearInterval(interval);
  }, [loadNotificationCounts]);

  const value: NotificationContextType = {
    // State
    notifications,
    unreadCount,
    counts,
    loading,
    error,

    // Actions
    loadNotifications,
    loadNotificationCounts,
    loadNotificationsByType,
    markAsRead,
    markAllAsRead,
    dismissNotification,
    dismissAllNotifications,
    batchMarkAsRead,
    batchDismiss,
    createNotification,
    createSystemNotification,
    refreshCounts,
    clearError,
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotifications() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider');
  }
  return context;
}