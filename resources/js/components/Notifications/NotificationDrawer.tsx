"use client";

import * as React from "react";
import {
  Bell,
  BellRing,
  Check,
  X,
  MessageCircle,
  User,
  Shield,
  AlertTriangle,
  Info,
  CheckCircle,
  Clock,
  MoreHorizontal,
  Archive,
  FileText
} from "lucide-react";
import { router } from "@inertiajs/react";

import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import { useNotifications } from "@/contexts/NotificationContext";
import { usePermissions } from "@/contexts/PermissionsContext";

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
  trigger_user?: {
    id: number;
    name: string;
    email: string;
  };
  related_type?: string;
  related_id?: number;
}

interface NotificationDrawerProps {
  children?: React.ReactNode;
  isOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function NotificationDrawer({ children, isOpen: controlledIsOpen, onOpenChange }: NotificationDrawerProps) {
  const [internalIsOpen, setInternalIsOpen] = React.useState(false);
  const [activeTab, setActiveTab] = React.useState("all");

  // Use controlled state if provided, otherwise use internal state
  const isOpen = controlledIsOpen !== undefined ? controlledIsOpen : internalIsOpen;
  const setIsOpen = onOpenChange || setInternalIsOpen;

  const {
    notifications,
    unreadCount,
    counts,
    loading,
    loadNotifications,
    markAsRead,
    markAllAsRead,
    dismissNotification,
    dismissAllNotifications,
    batchMarkAsRead,
    batchDismiss
  } = useNotifications();

  const { canPerformAction } = usePermissions();
  const canManageApplications = canPerformAction('applications', 'update');

  // Load notifications when drawer opens
  React.useEffect(() => {
    if (isOpen) {
      loadNotifications({
        page: 1,
        pageSize: 50,
        filters: activeTab === "unread" ? { unread_only: true } : {}
      });
    }
  }, [isOpen, activeTab]);

  const getNotificationIcon = (type: string, priority: string) => {
    const iconClass = cn(
      "h-4 w-4",
      priority === "high" ? "text-red-500" : 
      priority === "medium" ? "text-yellow-500" : "text-blue-500"
    );

    switch (type) {
      case "message":
        return <MessageCircle className={iconClass} />;
      case "mention":
        return <User className={iconClass} />;
      case "system":
        return <Shield className={iconClass} />;
      case "warning":
        return <AlertTriangle className={iconClass} />;
      case "success":
        return <CheckCircle className={iconClass} />;
      default:
        return <Info className={iconClass} />;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "high":
        return "border-l-red-500 bg-red-50 dark:bg-red-950/20";
      case "medium":
        return "border-l-yellow-500 bg-yellow-50 dark:bg-yellow-950/20";
      case "low":
        return "border-l-blue-500 bg-blue-50 dark:bg-blue-950/20";
      default:
        return "border-l-gray-500 bg-gray-50 dark:bg-gray-950/20";
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInMinutes = (now.getTime() - date.getTime()) / (1000 * 60);

    if (diffInMinutes < 1) {
      return "Just now";
    } else if (diffInMinutes < 60) {
      return `${Math.floor(diffInMinutes)}m ago`;
    } else if (diffInMinutes < 1440) { // 24 hours
      return `${Math.floor(diffInMinutes / 60)}h ago`;
    } else {
      return `${Math.floor(diffInMinutes / 1440)}d ago`;
    }
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map(word => word[0])
      .join("")
      .toUpperCase()
      .substring(0, 2);
  };

  const handleNotificationClick = async (notification: Notification) => {
    if (!notification.is_read) {
      await markAsRead(notification.id);
    }
    
    // Handle navigation based on notification type
    if (notification.related_type && notification.related_id) {
      // Navigate to related resource
      // This would be implemented based on your routing structure
      console.log(`Navigate to ${notification.related_type}:${notification.related_id}`);
    }
  };

  const handleMarkAllAsRead = async () => {
    await markAllAsRead();
  };

  const handleDismissAll = async () => {
    await dismissAllNotifications();
  };

  const getFilteredNotifications = () => {
    if (!notifications) return [];
    
    switch (activeTab) {
      case "unread":
        return notifications.filter(n => !n.is_read);
      case "messages":
        return notifications.filter(n => n.type === "message" || n.type === "mention");
      case "system":
        return notifications.filter(n => n.type === "system");
      default:
        return notifications;
    }
  };

  const filteredNotifications = getFilteredNotifications();

  return (
    <Drawer open={isOpen} onOpenChange={setIsOpen}>
      <DrawerTrigger asChild>
        {children || (
          <Button variant="ghost" size="icon" className="relative">
            <Bell className="h-5 w-5" />
            {/* Combined badge for notifications + pending applications (if user can manage applications) */}
            {(() => {
              const pendingApps = canManageApplications ? (counts?.pending_applications ?? 0) : 0;
              const total = unreadCount + pendingApps;
              return total > 0 ? (
                <Badge
                  variant="destructive"
                  className="absolute -top-1 -right-1 h-5 min-w-5 text-xs px-1"
                >
                  {total > 99 ? "99+" : total}
                </Badge>
              ) : null;
            })()}
          </Button>
        )}
      </DrawerTrigger>
      <DrawerContent className="max-w-md mx-auto">
        <DrawerHeader className="pb-4">
          <div className="flex items-center justify-between">
            <div>
              <DrawerTitle className="flex items-center gap-2">
                <BellRing className="h-5 w-5" />
                Notifications
              </DrawerTitle>
              <DrawerDescription>
                {unreadCount > 0 ? `${unreadCount} unread notifications` : "All caught up!"}
              </DrawerDescription>
            </div>
            <div className="flex items-center gap-2">
              {unreadCount > 0 && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleMarkAllAsRead}
                  disabled={loading}
                >
                  <Check className="h-4 w-4 mr-1" />
                  Mark all read
                </Button>
              )}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={handleMarkAllAsRead}>
                    <Check className="h-4 w-4 mr-2" />
                    Mark all as read
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={handleDismissAll}>
                    <Archive className="h-4 w-4 mr-2" />
                    Dismiss all
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </DrawerHeader>

        <div className="px-4">
          {/* Pending Applications Alert - only show for users with application management permissions */}
          {canManageApplications && counts?.pending_applications != null && counts.pending_applications > 0 && (
            <div
              className="mb-4 p-3 rounded-lg border border-l-4 border-l-amber-500 bg-amber-50 dark:bg-amber-950/20 cursor-pointer hover:bg-amber-100 dark:hover:bg-amber-950/30 transition-colors"
              onClick={() => {
                setIsOpen(false);
                router.visit('/applications');
              }}
            >
              <div className="flex items-center gap-3">
                <FileText className="h-5 w-5 text-amber-600" />
                <div className="flex-1">
                  <p className="text-sm font-medium text-amber-800 dark:text-amber-200">
                    {counts.pending_applications} Pending Application{counts.pending_applications !== 1 ? 's' : ''}
                  </p>
                  <p className="text-xs text-amber-600 dark:text-amber-400">
                    Click to review and process
                  </p>
                </div>
                <Badge variant="outline" className="bg-amber-100 text-amber-800 border-amber-300">
                  Action Required
                </Badge>
              </div>
            </div>
          )}

          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="all" className="text-xs">
                All
                {counts?.total != null && counts.total > 0 && (
                  <Badge variant="secondary" className="ml-1 h-4 min-w-4 text-xs px-1">
                    {counts.total}
                  </Badge>
                )}
              </TabsTrigger>
              <TabsTrigger value="unread" className="text-xs">
                Unread
                {counts?.unread != null && counts.unread > 0 && (
                  <Badge variant="destructive" className="ml-1 h-4 min-w-4 text-xs px-1">
                    {counts.unread}
                  </Badge>
                )}
              </TabsTrigger>
              <TabsTrigger value="messages" className="text-xs">
                Messages
                {counts?.unread_messages != null && counts.unread_messages > 0 && (
                  <Badge variant="secondary" className="ml-1 h-4 min-w-4 text-xs px-1">
                    {counts.unread_messages}
                  </Badge>
                )}
              </TabsTrigger>
              <TabsTrigger value="system" className="text-xs">
                System
              </TabsTrigger>
            </TabsList>

            <TabsContent value={activeTab} className="mt-4">
              <ScrollArea className="h-[400px] pr-4">
                {loading ? (
                  <div className="flex items-center justify-center py-8">
                    <div className="text-sm text-muted-foreground">Loading notifications...</div>
                  </div>
                ) : filteredNotifications.length === 0 ? (
                  <div className="flex flex-col items-center justify-center py-8 text-center">
                    <Bell className="h-8 w-8 text-muted-foreground mb-2" />
                    <div className="text-sm text-muted-foreground">
                      {activeTab === "unread" ? "No unread notifications" : "No notifications found"}
                    </div>
                  </div>
                ) : (
                  <div className="space-y-2">
                    {filteredNotifications.map((notification) => (
                      <div
                        key={notification.id}
                        className={cn(
                          "relative p-4 rounded-lg border border-l-4 cursor-pointer transition-colors hover:bg-muted/50",
                          getPriorityColor(notification.priority),
                          notification.is_read ? "opacity-70" : ""
                        )}
                        onClick={() => handleNotificationClick(notification)}
                      >
                        <div className="flex items-start gap-3">
                          <div className="flex-shrink-0 mt-0.5">
                            {getNotificationIcon(notification.type, notification.priority)}
                          </div>
                          
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between gap-2">
                              <div className="flex-1">
                                <p className="text-sm font-medium leading-tight">
                                  {notification.title}
                                </p>
                                {notification.message && (
                                  <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                                    {notification.message}
                                  </p>
                                )}
                              </div>
                              
                              <div className="flex items-center gap-2 flex-shrink-0">
                                {!notification.is_read && (
                                  <div className="h-2 w-2 rounded-full bg-blue-500" />
                                )}
                                <DropdownMenu>
                                  <DropdownMenuTrigger asChild>
                                    <Button 
                                      variant="ghost" 
                                      size="icon" 
                                      className="h-6 w-6"
                                      onClick={(e) => e.stopPropagation()}
                                    >
                                      <MoreHorizontal className="h-3 w-3" />
                                    </Button>
                                  </DropdownMenuTrigger>
                                  <DropdownMenuContent align="end">
                                    {!notification.is_read && (
                                      <DropdownMenuItem onClick={() => markAsRead(notification.id)}>
                                        <Check className="h-4 w-4 mr-2" />
                                        Mark as read
                                      </DropdownMenuItem>
                                    )}
                                    <DropdownMenuItem onClick={() => dismissNotification(notification.id)}>
                                      <X className="h-4 w-4 mr-2" />
                                      Dismiss
                                    </DropdownMenuItem>
                                  </DropdownMenuContent>
                                </DropdownMenu>
                              </div>
                            </div>
                            
                            <div className="flex items-center justify-between mt-2">
                              <div className="flex items-center gap-2">
                                {notification.trigger_user && (
                                  <div className="flex items-center gap-1">
                                    <Avatar className="h-4 w-4">
                                      <AvatarImage src={`/avatars/${notification.trigger_user.id}.jpg`} />
                                      <AvatarFallback className="text-xs">
                                        {getInitials(notification.trigger_user.name)}
                                      </AvatarFallback>
                                    </Avatar>
                                    <span className="text-xs text-muted-foreground">
                                      {notification.trigger_user.name}
                                    </span>
                                  </div>
                                )}
                                <Badge variant="outline" className="text-xs h-4 px-1">
                                  {notification.type}
                                </Badge>
                              </div>
                              
                              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                                <Clock className="h-3 w-3" />
                                {formatTime(notification.created_at)}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </ScrollArea>
            </TabsContent>
          </Tabs>
        </div>
      </DrawerContent>
    </Drawer>
  );
}