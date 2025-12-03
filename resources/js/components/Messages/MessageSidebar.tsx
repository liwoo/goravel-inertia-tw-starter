"use client";

import * as React from "react";
import {
  MessageCircle,
  Send,
  Plus,
  Loader2,
  Megaphone,
  RefreshCw
} from "lucide-react";

import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import { useMessages, Conversation, MessageUser } from "@/contexts/MessageContext";
import { usePresence } from "@/contexts/PresenceContext";
import { BroadcastToRole } from "./BroadcastToRole";
import { OnlineIndicator } from "@/components/ui/online-indicator";

interface User extends MessageUser {
  is_super_admin?: boolean;
  isSuperAdmin?: boolean;
}

interface MessageSidebarProps {
  user: User;
  showNewMessage?: boolean;
  onShowNewMessageChange?: (show: boolean) => void;
  onBroadcastModeChange?: (isBroadcast: boolean) => void;
  onBroadcastSent?: () => void;
}

type TabCategory = "inbox" | "sent" | "users" | "broadcast";

const tabs: { title: string; icon: React.ElementType; category: TabCategory }[] = [
  { title: "Inbox", icon: MessageCircle, category: "inbox" },
  { title: "Sent", icon: Send, category: "sent" },
];

export function MessageSidebar({ user, showNewMessage, onShowNewMessageChange, onBroadcastModeChange, onBroadcastSent }: MessageSidebarProps) {
  const [activeTab, setActiveTab] = React.useState<TabCategory>("inbox");
  const [searchTerm, setSearchTerm] = React.useState("");
  const [showUnreadOnly, setShowUnreadOnly] = React.useState(false);

  const {
    conversations,
    messagableUsers,
    selectedConversation,
    setSelectedConversation,
    loading,
    unreadCount,
    loadConversations,
    loadMessagableUsers
  } = useMessages();

  const { refreshPresence, isRefreshing, onlineCount } = usePresence();

  // Switch to users view when requested externally
  React.useEffect(() => {
    if (showNewMessage) {
      setActiveTab("users");
    }
  }, [showNewMessage]);

  // Notify parent when broadcast mode changes
  React.useEffect(() => {
    onBroadcastModeChange?.(activeTab === "broadcast");
  }, [activeTab, onBroadcastModeChange]);

  // Load data based on active tab
  React.useEffect(() => {
    if (activeTab === "inbox" || activeTab === "sent") {
      loadConversations({
        page: 1,
        pageSize: 20,
        search: searchTerm,
        filters: {
          unread_only: showUnreadOnly
        }
      });
    } else if (activeTab === "users") {
      loadMessagableUsers({
        page: 1,
        pageSize: 50,
        search: searchTerm
      });
    }
  }, [activeTab, searchTerm, showUnreadOnly, loadConversations, loadMessagableUsers]);

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
      return date.toLocaleDateString();
    }
  };

  const truncateMessage = (content: string, maxLength: number = 50) => {
    if (!content) return "";
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + "...";
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map(word => word[0])
      .join("")
      .toUpperCase()
      .substring(0, 2);
  };

  const handleConversationSelect = (conversation: Conversation) => {
    setSelectedConversation(conversation);
  };

  const handleUserSelect = (selectedUser: User) => {
    const now = new Date().toISOString();
    const newConversation: Conversation = {
      user: selectedUser,
      latest_message: {
        id: 0,
        content: "",
        type: "direct",
        status: "draft",
        created_at: now,
        updated_at: now,
        sender_id: user.id,
        is_edited: false
      },
      unread_count: 0,
      last_activity: now
    };
    setSelectedConversation(newConversation);
    onShowNewMessageChange?.(false);
    setActiveTab("inbox");
  };

  const handleNewMessageClick = () => {
    setActiveTab("users");
    onShowNewMessageChange?.(true);
  };

  const renderConversations = () => {
    if (loading) {
      return (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      );
    }

    if (!conversations || conversations.length === 0) {
      return (
        <div className="p-4 text-center">
          <p className="text-sm text-muted-foreground mb-3">No conversations yet</p>
          <Button variant="outline" size="sm" onClick={handleNewMessageClick}>
            <Plus className="h-4 w-4 mr-2" />
            Start a conversation
          </Button>
        </div>
      );
    }

    return conversations.map((conversation) => (
      <button
        key={conversation.user.id}
        className={cn(
          "flex w-full items-start gap-3 p-3 text-left text-sm transition-colors hover:bg-accent rounded-lg",
          selectedConversation?.user.id === conversation.user.id && "bg-accent"
        )}
        onClick={() => handleConversationSelect(conversation)}
      >
        <div className="relative shrink-0">
          <Avatar className="h-10 w-10">
            <AvatarImage src={`/avatars/${conversation.user.id}.jpg`} />
            <AvatarFallback className="text-xs">
              {getInitials(conversation.user.name)}
            </AvatarFallback>
          </Avatar>
          <OnlineIndicator
            userId={conversation.user.id}
            size="sm"
            className="absolute bottom-0 right-0 border-2 border-background"
          />
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-0.5">
            <span className="font-medium truncate">{conversation.user.name}</span>
            <span className="text-xs text-muted-foreground shrink-0">
              {formatTime(conversation.last_activity)}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <p className={cn(
              "text-xs text-muted-foreground truncate flex-1",
              conversation.unread_count > 0 && "font-medium text-foreground"
            )}>
              {conversation.latest_message.sender_id === user.id && "You: "}
              {truncateMessage(conversation.latest_message.content)}
            </p>
            {conversation.unread_count > 0 && (
              <Badge variant="destructive" className="h-5 min-w-5 text-xs px-1.5">
                {conversation.unread_count > 99 ? "99+" : conversation.unread_count}
              </Badge>
            )}
          </div>
        </div>
      </button>
    ));
  };

  const renderUsers = () => {
    if (loading) {
      return (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      );
    }

    if (!messagableUsers || messagableUsers.length === 0) {
      return (
        <div className="p-4 text-center text-sm text-muted-foreground">
          No users available to message
        </div>
      );
    }

    return messagableUsers.map((msgUser) => (
      <button
        key={msgUser.id}
        className="flex w-full items-center gap-3 p-3 text-left text-sm transition-colors hover:bg-accent rounded-lg"
        onClick={() => handleUserSelect(msgUser)}
      >
        <div className="relative shrink-0">
          <Avatar className="h-10 w-10">
            <AvatarImage src={`/avatars/${msgUser.id}.jpg`} />
            <AvatarFallback className="text-xs">
              {getInitials(msgUser.name)}
            </AvatarFallback>
          </Avatar>
          <OnlineIndicator
            userId={msgUser.id}
            size="sm"
            className="absolute bottom-0 right-0 border-2 border-background"
          />
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-0.5">
            <span className="font-medium truncate">{msgUser.name}</span>
          </div>
          <p className="text-xs text-muted-foreground truncate">
            {msgUser.email}
          </p>
          {msgUser.roles && msgUser.roles.length > 0 && (
            <div className="flex gap-1 mt-1">
              {msgUser.roles.slice(0, 2).map((role) => (
                <Badge key={role.id} variant="secondary" className="text-xs px-1.5 h-4">
                  {role.name}
                </Badge>
              ))}
            </div>
          )}
        </div>
      </button>
    ));
  };

  return (
    <div className="flex flex-col h-full">
      {/* Header with New Message button */}
      <div className="p-4 border-b">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <MessageCircle className="h-5 w-5" />
            <span className="font-semibold">Messages</span>
            {unreadCount > 0 && (
              <Badge variant="destructive" className="h-5 min-w-5 text-xs px-1.5">
                {unreadCount}
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-1">
            {/* Online count indicator */}
            <div className="flex items-center gap-1 text-xs text-muted-foreground mr-1">
              <span className="h-2 w-2 rounded-full bg-green-500" />
              <span>{onlineCount}</span>
            </div>
            {/* Refresh presence button */}
            <Button
              variant="ghost"
              size="icon"
              onClick={refreshPresence}
              disabled={isRefreshing}
              className="h-8 w-8"
              title="Refresh online status"
            >
              <RefreshCw className={cn("h-4 w-4", isRefreshing && "animate-spin")} />
            </Button>
            <Button
              variant="default"
              size="sm"
              onClick={handleNewMessageClick}
              className="h-8"
            >
              <Plus className="h-4 w-4 mr-1" />
              New
            </Button>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 flex-wrap">
          {tabs.map((tab) => (
            <Button
              key={tab.category}
              variant={activeTab === tab.category ? "secondary" : "ghost"}
              size="sm"
              className="h-8 px-3"
              onClick={() => {
                setActiveTab(tab.category);
                onShowNewMessageChange?.(false);
              }}
            >
              <tab.icon className="h-4 w-4 mr-1.5" />
              {tab.title}
            </Button>
          ))}
          {/* Broadcast tab - Super admin only */}
          {(user.is_super_admin || user.isSuperAdmin) && (
            <Button
              variant={activeTab === "broadcast" ? "secondary" : "ghost"}
              size="sm"
              className="h-8 px-3"
              onClick={() => {
                setActiveTab("broadcast");
                onShowNewMessageChange?.(false);
              }}
            >
              <Megaphone className="h-4 w-4 mr-1.5" />
              Broadcast
            </Button>
          )}
        </div>
      </div>

      {/* Broadcast view - Super admin only */}
      {activeTab === "broadcast" && (user.is_super_admin || user.isSuperAdmin) && (
        <BroadcastToRole
          onSuccess={() => {
            onBroadcastSent?.();
          }}
        />
      )}

      {/* Search and filters - only for non-broadcast views */}
      {activeTab !== "broadcast" && (
        <div className="p-3 border-b space-y-2">
          {activeTab === "users" ? (
            <>
              <p className="text-sm font-medium">Select Recipient</p>
              <Input
                placeholder="Search users..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="h-8"
              />
            </>
          ) : (
            <>
              <Input
                placeholder="Search conversations..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="h-8"
              />
              <Label className="flex items-center gap-2 text-xs">
                <Switch
                  checked={showUnreadOnly}
                  onCheckedChange={setShowUnreadOnly}
                  className="scale-75"
                />
                <span>Show unread only</span>
              </Label>
            </>
          )}
        </div>
      )}

      {/* Content - only for non-broadcast views */}
      {activeTab !== "broadcast" && (
        <ScrollArea className="flex-1">
          <div className="p-2">
            {activeTab === "users" ? renderUsers() : renderConversations()}
          </div>
        </ScrollArea>
      )}
    </div>
  );
}
