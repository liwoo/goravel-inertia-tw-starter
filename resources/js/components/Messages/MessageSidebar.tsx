"use client";

import * as React from "react";
import { 
  MessageCircle, 
  Search, 
  Send, 
  Archive, 
  Users,
  Plus,
  MoreHorizontal,
  Pin,
  Clock
} from "lucide-react";

import { NavUser } from "@/components/nav-user";
import { Label } from "@/components/ui/label";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { useMessages } from "@/contexts/MessageContext";

interface User {
  id: number;
  name: string;
  email: string;
  is_active: boolean;
  roles?: Array<{
    id: number;
    name: string;
    slug: string;
  }>;
}

interface Conversation {
  user: User;
  latest_message: {
    id: number;
    content: string;
    created_at: string;
    sender_id: number;
    is_edited: boolean;
  };
  unread_count: number;
  last_activity: string;
}

interface MessageSidebarProps extends React.ComponentProps<typeof Sidebar> {
  user: User;
  onNewMessage?: () => void;
}

const navigationItems = [
  {
    title: "Inbox",
    icon: MessageCircle,
    isActive: true,
    category: "inbox"
  },
  {
    title: "Sent",
    icon: Send,
    isActive: false,
    category: "sent"
  },
  {
    title: "Archive",
    icon: Archive,
    isActive: false,
    category: "archive"
  },
  {
    title: "All Users",
    icon: Users,
    isActive: false,
    category: "users"
  },
];

export function MessageSidebar({ user, onNewMessage, ...props }: MessageSidebarProps) {
  const [activeItem, setActiveItem] = React.useState(navigationItems[0]);
  const [searchTerm, setSearchTerm] = React.useState("");
  const [showUnreadOnly, setShowUnreadOnly] = React.useState(false);
  const { setOpen } = useSidebar();

  const {
    conversations,
    messagableUsers,
    selectedConversation,
    setSelectedConversation,
    loading,
    unreadCount,
    searchUsers,
    loadConversations,
    loadMessagableUsers
  } = useMessages();

  // Load data based on active item
  React.useEffect(() => {
    if (activeItem.category === "inbox" || activeItem.category === "sent") {
      loadConversations({
        page: 1,
        pageSize: 20,
        search: searchTerm,
        filters: {
          unread_only: showUnreadOnly
        }
      });
    } else if (activeItem.category === "users") {
      loadMessagableUsers({
        page: 1,
        pageSize: 50,
        search: searchTerm
      });
    }
  }, [activeItem, searchTerm, showUnreadOnly]);

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60);

    if (diffInHours < 1) {
      return "Just now";
    } else if (diffInHours < 24) {
      return `${Math.floor(diffInHours)}h ago`;
    } else if (diffInHours < 168) { // 7 days
      return `${Math.floor(diffInHours / 24)}d ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  const truncateMessage = (content: string, maxLength: number = 60) => {
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
    setOpen(true);
  };

  const handleUserSelect = (selectedUser: User) => {
    // Create a mock conversation for new message
    const newConversation: Conversation = {
      user: selectedUser,
      latest_message: {
        id: 0,
        content: "",
        created_at: new Date().toISOString(),
        sender_id: user.id,
        is_edited: false
      },
      unread_count: 0,
      last_activity: new Date().toISOString()
    };
    setSelectedConversation(newConversation);
    setOpen(true);
  };

  const renderConversations = () => {
    if (loading) {
      return (
        <div className="p-4 text-center text-sm text-muted-foreground">
          Loading conversations...
        </div>
      );
    }

    if (!conversations || conversations.length === 0) {
      return (
        <div className="p-4 text-center text-sm text-muted-foreground">
          No conversations found
        </div>
      );
    }

    return conversations.map((conversation) => (
      <button
        key={conversation.user.id}
        className={cn(
          "flex w-full items-start gap-3 border-b p-4 text-left text-sm transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground last:border-b-0",
          selectedConversation?.user.id === conversation.user.id && "bg-sidebar-accent"
        )}
        onClick={() => handleConversationSelect(conversation)}
      >
        <Avatar className="h-10 w-10 shrink-0">
          <AvatarImage src={`/avatars/${conversation.user.id}.jpg`} />
          <AvatarFallback className="text-xs">
            {getInitials(conversation.user.name)}
          </AvatarFallback>
        </Avatar>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-1">
            <span className="font-medium truncate">{conversation.user.name}</span>
            <div className="flex items-center gap-1 shrink-0">
              {conversation.unread_count > 0 && (
                <Badge variant="destructive" className="h-5 min-w-5 text-xs px-1">
                  {conversation.unread_count > 99 ? "99+" : conversation.unread_count}
                </Badge>
              )}
              <span className="text-xs text-muted-foreground">
                {formatTime(conversation.last_activity)}
              </span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <p className={cn(
              "text-xs text-muted-foreground truncate flex-1",
              conversation.unread_count > 0 && "font-medium text-foreground"
            )}>
              {conversation.latest_message.sender_id === user.id && "You: "}
              {truncateMessage(conversation.latest_message.content)}
              {conversation.latest_message.is_edited && (
                <span className="ml-1 text-xs opacity-70">(edited)</span>
              )}
            </p>
          </div>
          {conversation.user.roles && conversation.user.roles.length > 0 && (
            <div className="flex gap-1 mt-1">
              {conversation.user.roles.slice(0, 2).map((role) => (
                <Badge key={role.id} variant="outline" className="text-xs px-1 h-4">
                  {role.name}
                </Badge>
              ))}
              {conversation.user.roles.length > 2 && (
                <Badge variant="outline" className="text-xs px-1 h-4">
                  +{conversation.user.roles.length - 2}
                </Badge>
              )}
            </div>
          )}
        </div>
      </button>
    ));
  };

  const renderUsers = () => {
    if (loading) {
      return (
        <div className="p-4 text-center text-sm text-muted-foreground">
          Loading users...
        </div>
      );
    }

    if (!messagableUsers || messagableUsers.length === 0) {
      return (
        <div className="p-4 text-center text-sm text-muted-foreground">
          No users found
        </div>
      );
    }

    return messagableUsers.map((msgUser) => (
      <button
        key={msgUser.id}
        className="flex w-full items-center gap-3 border-b p-4 text-left text-sm transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground last:border-b-0"
        onClick={() => handleUserSelect(msgUser)}
      >
        <Avatar className="h-10 w-10 shrink-0">
          <AvatarImage src={`/avatars/${msgUser.id}.jpg`} />
          <AvatarFallback className="text-xs">
            {getInitials(msgUser.name)}
          </AvatarFallback>
        </Avatar>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-1">
            <span className="font-medium truncate">{msgUser.name}</span>
            {msgUser.is_active && (
              <div className="h-2 w-2 rounded-full bg-green-500 shrink-0" />
            )}
          </div>
          <p className="text-xs text-muted-foreground truncate">
            {msgUser.email}
          </p>
          {msgUser.roles && msgUser.roles.length > 0 && (
            <div className="flex gap-1 mt-1">
              {msgUser.roles.slice(0, 2).map((role) => (
                <Badge key={role.id} variant="outline" className="text-xs px-1 h-4">
                  {role.name}
                </Badge>
              ))}
              {msgUser.roles.length > 2 && (
                <Badge variant="outline" className="text-xs px-1 h-4">
                  +{msgUser.roles.length - 2}
                </Badge>
              )}
            </div>
          )}
        </div>
      </button>
    ));
  };

  return (
    <Sidebar
      collapsible="icon"
      className="overflow-hidden *:data-[sidebar=sidebar]:flex-row"
      {...props}
    >
      {/* Navigation sidebar */}
      <Sidebar
        collapsible="none"
        className="w-[calc(var(--sidebar-width-icon)+1px)]! border-r"
      >
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton size="lg" asChild className="md:h-8 md:p-0">
                <div className="cursor-pointer">
                  <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                    <MessageCircle className="size-4" />
                  </div>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-medium">Messages</span>
                    <span className="truncate text-xs">
                      {unreadCount > 0 ? `${unreadCount} unread` : "All caught up"}
                    </span>
                  </div>
                </div>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupContent className="px-1.5 md:px-0">
              <SidebarMenu>
                {navigationItems.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      tooltip={{
                        children: item.title,
                        hidden: false,
                      }}
                      onClick={() => {
                        setActiveItem(item);
                        setOpen(true);
                      }}
                      isActive={activeItem?.title === item.title}
                      className="px-2.5 md:px-2"
                    >
                      <item.icon />
                      <span>{item.title}</span>
                      {item.title === "Inbox" && unreadCount > 0 && (
                        <Badge variant="destructive" className="ml-auto h-5 min-w-5 text-xs px-1">
                          {unreadCount > 99 ? "99+" : unreadCount}
                        </Badge>
                      )}
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <NavUser user={user} />
        </SidebarFooter>
      </Sidebar>

      {/* Content sidebar */}
      <Sidebar collapsible="none" className="hidden flex-1 md:flex">
        <SidebarHeader className="gap-3.5 border-b p-4">
          <div className="flex w-full items-center justify-between">
            <div className="text-foreground text-base font-medium">
              {activeItem?.title}
            </div>
            <div className="flex items-center gap-2">
              {(activeItem?.category === "inbox" || activeItem?.category === "sent") && (
                <Label className="flex items-center gap-2 text-sm">
                  <span>Unread</span>
                  <Switch 
                    checked={showUnreadOnly}
                    onCheckedChange={setShowUnreadOnly}
                    className="shadow-none" 
                  />
                </Label>
              )}
              {activeItem?.category === "users" && onNewMessage && (
                <Button onClick={onNewMessage} size="sm" variant="outline">
                  <Plus className="h-4 w-4" />
                </Button>
              )}
            </div>
          </div>
          <SidebarInput 
            placeholder={
              activeItem?.category === "users" 
                ? "Search users..." 
                : "Search conversations..."
            }
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup className="px-0">
            <SidebarGroupContent>
              {activeItem?.category === "users" ? renderUsers() : renderConversations()}
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
      </Sidebar>
    </Sidebar>
  );
}