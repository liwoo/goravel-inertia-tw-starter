"use client";

import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import axios from '@/lib/axios';
import { useSSE, useSSEEvent } from '@/hooks/useSSE';
import sseManager from '@/services/sseManager';

export interface MessageUser {
  id: number;
  name: string;
  email: string;
  is_active?: boolean;
  role?: string;
  roles?: Array<{
    id: number;
    name: string;
    slug?: string;
  }>;
}

export interface Message {
  id: number;
  content: string;
  type: string;
  status: string;
  sender_id: number;
  recipient_id?: number;
  is_edited: boolean;
  edited_at?: string;
  read_at?: string;
  created_at: string;
  updated_at: string;
  sender?: MessageUser;
  recipient?: MessageUser;
  parent_message_id?: number;
  parent_message?: Message;
  replies?: Message[];
}

export interface Conversation {
  user: MessageUser;
  latest_message: Message;
  unread_count: number;
  last_activity: string;
}

// Alias for internal use
type User = MessageUser;

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

interface MessageContextType {
  // State
  conversations: Conversation[];
  messages: Message[];
  messagableUsers: User[];
  selectedConversation: Conversation | null;
  unreadCount: number;
  loading: boolean;
  error: string | null;

  // Actions
  loadConversations: (request: ListRequest) => Promise<void>;
  loadConversation: (userId: number, request: ListRequest) => Promise<void>;
  loadMessagableUsers: (request: ListRequest) => Promise<void>;
  sendMessage: (recipientId: number, content: string, type?: string) => Promise<Message>;
  editMessage: (messageId: number, content: string) => Promise<Message>;
  deleteMessage: (messageId: number) => Promise<void>;
  markAsRead: (userId: number) => Promise<void>;
  replyToMessage: (messageId: number, content: string, recipientId?: number) => Promise<Message>;
  searchUsers: (query: string) => Promise<User[]>;
  setSelectedConversation: (conversation: Conversation | null) => void;
  refreshUnreadCount: () => Promise<void>;
  clearError: () => void;
}

const MessageContext = createContext<MessageContextType | null>(null);

export function MessageProvider({ children }: { children: React.ReactNode }) {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [messages, setMessages] = useState<Message[]>([]);
  const [messagableUsers, setMessagableUsers] = useState<User[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<Conversation | null>(null);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // SSE Event Handlers
  useSSE('message:unread_count', (event) => {
    setUnreadCount(event.data.count || 0);
  });

  useSSE('message:new', (event) => {
    const newMessage = event.data.message as Message;
    
    // Update messages if it's in the current conversation
    if (selectedConversation && 
        (selectedConversation.user.id === newMessage.sender_id || 
         selectedConversation.user.id === newMessage.recipient_id)) {
      setMessages(prev => [newMessage, ...prev]);
    }
    
    // Refresh conversations to update latest message
    loadConversations({ page: 1, pageSize: 20 });
  });

  useSSE('message:read', (event) => {
    const messageId = event.data.message_id;
    
    // Update message read status
    setMessages(prev => 
      prev.map(msg => 
        msg.id === messageId 
          ? { ...msg, status: 'read', read_at: new Date().toISOString() } 
          : msg
      )
    );
  });

  useSSE('message:updated', (event) => {
    const updatedMessage = event.data as Message;
    
    // Update message in current conversation
    setMessages(prev => 
      prev.map(msg => msg.id === updatedMessage.id ? updatedMessage : msg)
    );
  });

  useSSE('message:deleted', (event) => {
    const messageId = event.data.message_id;
    
    // Remove message from current conversation
    setMessages(prev => prev.filter(msg => msg.id !== messageId));
    
    // Refresh conversations if needed
    if (messages.some(msg => msg.id === messageId)) {
      loadConversations({ page: 1, pageSize: 20 });
    }
  });

  // Load conversations
  const loadConversations = useCallback(async (request: ListRequest) => {
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

      const response = await axios.get(`/api/messages/conversations?${params}`);
      
      if (response.data.success && response.data.data) {
        const paginatedData = response.data.data as PaginatedResult<Conversation>;
        setConversations(paginatedData.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load conversations');
      console.error('Error loading conversations:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load specific conversation messages
  const loadConversation = useCallback(async (userId: number, request: ListRequest) => {
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

      const response = await axios.get(`/api/messages/conversation/${userId}?${params}`);
      
      if (response.data.success && response.data.data) {
        const paginatedData = response.data.data as PaginatedResult<Message>;
        setMessages(paginatedData.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load conversation');
      console.error('Error loading conversation:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load messagable users
  const loadMessagableUsers = useCallback(async (request: ListRequest) => {
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

      const response = await axios.get(`/api/messages/users?${params}`);
      
      if (response.data.success && response.data.data) {
        const paginatedData = response.data.data as PaginatedResult<User>;
        setMessagableUsers(paginatedData.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load users');
      console.error('Error loading messagable users:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Send message
  const sendMessage = useCallback(async (recipientId: number, content: string, type: string = 'direct'): Promise<Message> => {
    try {
      setError(null);
      
      const response = await axios.post('/api/messages', {
        recipient_id: recipientId,
        content,
        type
      });
      
      if (response.data.success && response.data.data) {
        const newMessage = response.data.data as Message;
        
        // Add to messages if this is the current conversation
        if (selectedConversation && selectedConversation.user.id === recipientId) {
          setMessages(prev => [newMessage, ...prev]);
        }
        
        // Refresh conversations to update latest message
        await loadConversations({ page: 1, pageSize: 20 });
        await refreshUnreadCount();
        
        return newMessage;
      }
      
      throw new Error('Failed to send message');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to send message';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [selectedConversation, loadConversations]);

  // Edit message
  const editMessage = useCallback(async (messageId: number, content: string): Promise<Message> => {
    try {
      setError(null);
      
      const response = await axios.put(`/api/messages/${messageId}`, {
        content
      });
      
      if (response.data.success && response.data.data) {
        const updatedMessage = response.data.data as Message;
        
        // Update message in current conversation
        setMessages(prev => 
          prev.map(msg => msg.id === messageId ? updatedMessage : msg)
        );
        
        return updatedMessage;
      }
      
      throw new Error('Failed to edit message');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to edit message';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, []);

  // Delete message
  const deleteMessage = useCallback(async (messageId: number): Promise<void> => {
    try {
      setError(null);
      
      const response = await axios.delete(`/api/messages/${messageId}`);
      
      if (response.data.success) {
        // Remove message from current conversation
        setMessages(prev => prev.filter(msg => msg.id !== messageId));
        
        // Refresh conversations to update latest message
        await loadConversations({ page: 1, pageSize: 20 });
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to delete message';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [loadConversations]);

  // Mark messages as read
  const markAsRead = useCallback(async (userId: number): Promise<void> => {
    try {
      await axios.put(`/api/messages/conversation/${userId}/read`);

      // Update conversation to show messages as read
      setConversations(prev =>
        prev.map(conv =>
          conv.user.id === userId
            ? { ...conv, unread_count: 0 }
            : conv
        )
      );

      // Refresh unread count
      try {
        const response = await axios.get('/api/messages/unread-count');
        if (response.data.success && response.data.data) {
          setUnreadCount(response.data.data.unread_count || 0);
        }
      } catch {
        // Ignore unread count refresh errors
      }
    } catch (err: any) {
      // Don't set error for mark as read failures - it's not critical
      console.error('Error marking messages as read:', err);
    }
  }, []);

  // Reply to message
  const replyToMessage = useCallback(async (messageId: number, content: string, recipientId?: number): Promise<Message> => {
    try {
      setError(null);
      
      const payload: any = { content };
      if (recipientId) {
        payload.recipient_id = recipientId;
      }
      
      const response = await axios.post(`/api/messages/${messageId}/reply`, payload);
      
      if (response.data.success && response.data.data) {
        const newMessage = response.data.data as Message;
        
        // Add to messages
        setMessages(prev => [newMessage, ...prev]);
        
        // Refresh conversations
        await loadConversations({ page: 1, pageSize: 20 });
        await refreshUnreadCount();
        
        return newMessage;
      }
      
      throw new Error('Failed to send reply');
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to send reply';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [loadConversations]);

  // Search users
  const searchUsers = useCallback(async (query: string): Promise<User[]> => {
    try {
      setError(null);
      
      const response = await axios.get(`/api/messages/search-users?q=${encodeURIComponent(query)}`);
      
      if (response.data.success && response.data.data) {
        return response.data.data as User[];
      }
      
      return [];
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to search users';
      setError(errorMessage);
      console.error('Error searching users:', err);
      return [];
    }
  }, []);

  // Refresh unread count
  const refreshUnreadCount = useCallback(async () => {
    try {
      const response = await axios.get('/api/messages/unread-count');
      
      if (response.data.success && response.data.data) {
        setUnreadCount(response.data.data.unread_count || 0);
      }
    } catch (err: any) {
      console.error('Error refreshing unread count:', err);
    }
  }, []);

  // Clear error
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Load initial data and connect SSE
  useEffect(() => {
    refreshUnreadCount();
    
    // Ensure SSE is connected
    if (!sseManager.isConnected()) {
      sseManager.connect();
    }
  }, [refreshUnreadCount]);

  const value: MessageContextType = {
    // State
    conversations,
    messages,
    messagableUsers,
    selectedConversation,
    unreadCount,
    loading,
    error,

    // Actions
    loadConversations,
    loadConversation,
    loadMessagableUsers,
    sendMessage,
    editMessage,
    deleteMessage,
    markAsRead,
    replyToMessage,
    searchUsers,
    setSelectedConversation,
    refreshUnreadCount,
    clearError,
  };

  return (
    <MessageContext.Provider value={value}>
      {children}
    </MessageContext.Provider>
  );
}

export function useMessages() {
  const context = useContext(MessageContext);
  if (!context) {
    throw new Error('useMessages must be used within a MessageProvider');
  }
  return context;
}