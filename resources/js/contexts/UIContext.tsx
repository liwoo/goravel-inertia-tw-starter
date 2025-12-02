"use client";

import React, { createContext, useContext, useState, useCallback } from 'react';

interface UIContextType {
  // Messages drawer
  isMessagesOpen: boolean;
  openMessages: () => void;
  closeMessages: () => void;
  toggleMessages: () => void;

  // Notifications drawer
  isNotificationsOpen: boolean;
  openNotifications: () => void;
  closeNotifications: () => void;
  toggleNotifications: () => void;
}

const UIContext = createContext<UIContextType | null>(null);

export function UIProvider({ children }: { children: React.ReactNode }) {
  const [isMessagesOpen, setIsMessagesOpen] = useState(false);
  const [isNotificationsOpen, setIsNotificationsOpen] = useState(false);

  const openMessages = useCallback(() => setIsMessagesOpen(true), []);
  const closeMessages = useCallback(() => setIsMessagesOpen(false), []);
  const toggleMessages = useCallback(() => setIsMessagesOpen(prev => !prev), []);

  const openNotifications = useCallback(() => setIsNotificationsOpen(true), []);
  const closeNotifications = useCallback(() => setIsNotificationsOpen(false), []);
  const toggleNotifications = useCallback(() => setIsNotificationsOpen(prev => !prev), []);

  const value: UIContextType = {
    isMessagesOpen,
    openMessages,
    closeMessages,
    toggleMessages,
    isNotificationsOpen,
    openNotifications,
    closeNotifications,
    toggleNotifications,
  };

  return (
    <UIContext.Provider value={value}>
      {children}
    </UIContext.Provider>
  );
}

export function useUI() {
  const context = useContext(UIContext);
  if (!context) {
    throw new Error('useUI must be used within a UIProvider');
  }
  return context;
}
