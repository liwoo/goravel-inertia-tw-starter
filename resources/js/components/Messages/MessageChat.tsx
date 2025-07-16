"use client";

import * as React from "react";
import { Send, Edit3, Trash2, Reply, MoreHorizontal, Smile } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { UserMentionInput } from "./UserMentionInput";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Textarea } from "@/components/ui/textarea";
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

interface Message {
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
  sender?: User;
  recipient?: User;
  parent_message_id?: number;
  parent_message?: Message;
  replies?: Message[];
}

interface Conversation {
  user: User;
  latest_message: Message;
  unread_count: number;
  last_activity: string;
}

interface MessageChatProps {
  currentUser: User;
  conversation: Conversation | null;
  className?: string;
}

export function MessageChat({ currentUser, conversation, className }: MessageChatProps) {
  const [newMessage, setNewMessage] = React.useState("");
  const [editingMessage, setEditingMessage] = React.useState<number | null>(null);
  const [editContent, setEditContent] = React.useState("");
  const [replyingTo, setReplyingTo] = React.useState<Message | null>(null);
  const messagesEndRef = React.useRef<HTMLDivElement>(null);

  const {
    messages,
    loading,
    error,
    loadConversation,
    sendMessage,
    editMessage,
    deleteMessage,
    markAsRead,
    replyToMessage,
    clearError
  } = useMessages();

  // Load conversation messages when conversation changes
  React.useEffect(() => {
    if (conversation?.user.id) {
      loadConversation(conversation.user.id, {
        page: 1,
        pageSize: 50,
        sort: "created_at",
        direction: "ASC"
      });
      
      // Mark messages as read
      markAsRead(conversation.user.id);
    }
  }, [conversation?.user.id, loadConversation, markAsRead]);

  // Scroll to bottom when new messages arrive
  React.useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === today.toDateString()) {
      return "Today";
    } else if (date.toDateString() === yesterday.toDateString()) {
      return "Yesterday";
    } else {
      return date.toLocaleDateString();
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

  const handleSendMessage = async () => {
    if (!newMessage.trim() || !conversation?.user.id) return;

    try {
      if (replyingTo) {
        await replyToMessage(replyingTo.id, newMessage.trim(), conversation.user.id);
        setReplyingTo(null);
      } else {
        await sendMessage(conversation.user.id, newMessage.trim());
      }
      setNewMessage("");
    } catch (err) {
      console.error("Failed to send message:", err);
    }
  };

  const handleEditMessage = async (messageId: number) => {
    if (!editContent.trim()) return;

    try {
      await editMessage(messageId, editContent.trim());
      setEditingMessage(null);
      setEditContent("");
    } catch (err) {
      console.error("Failed to edit message:", err);
    }
  };

  const handleDeleteMessage = async (messageId: number) => {
    try {
      await deleteMessage(messageId);
    } catch (err) {
      console.error("Failed to delete message:", err);
    }
  };

  const startEditing = (message: Message) => {
    setEditingMessage(message.id);
    setEditContent(message.content);
  };

  const cancelEditing = () => {
    setEditingMessage(null);
    setEditContent("");
  };

  const startReply = (message: Message) => {
    setReplyingTo(message);
    setNewMessage(`@${message.sender?.name} `);
  };

  const cancelReply = () => {
    setReplyingTo(null);
    setNewMessage("");
  };

  // Group messages by date
  const groupMessagesByDate = (messages: Message[]) => {
    const groups: { [date: string]: Message[] } = {};
    
    messages.forEach(message => {
      const date = new Date(message.created_at).toDateString();
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(message);
    });
    
    return groups;
  };

  const messageGroups = groupMessagesByDate(messages);

  if (!conversation) {
    return (
      <div className={cn("flex items-center justify-center h-full", className)}>
        <div className="text-center text-muted-foreground">
          <div className="text-lg font-medium mb-2">No conversation selected</div>
          <div className="text-sm">Choose a conversation from the sidebar to start messaging</div>
        </div>
      </div>
    );
  }

  return (
    <div className={cn("flex flex-col h-full", className)}>
      {/* Chat Header */}
      <div className="flex items-center gap-3 p-4 border-b">
        <Avatar className="h-10 w-10">
          <AvatarImage src={`/avatars/${conversation.user.id}.jpg`} />
          <AvatarFallback>
            {getInitials(conversation.user.name)}
          </AvatarFallback>
        </Avatar>
        <div className="flex-1">
          <div className="font-medium">{conversation.user.name}</div>
          <div className="text-sm text-muted-foreground">{conversation.user.email}</div>
        </div>
        {conversation.user.roles && conversation.user.roles.length > 0 && (
          <div className="flex gap-1">
            {conversation.user.roles.slice(0, 2).map((role) => (
              <Badge key={role.id} variant="outline" className="text-xs">
                {role.name}
              </Badge>
            ))}
          </div>
        )}
      </div>

      {/* Messages */}
      <ScrollArea className="flex-1 p-4">
        {loading && messages.length === 0 ? (
          <div className="flex items-center justify-center py-8">
            <div className="text-sm text-muted-foreground">Loading messages...</div>
          </div>
        ) : Object.keys(messageGroups).length === 0 ? (
          <div className="flex items-center justify-center py-8">
            <div className="text-center text-muted-foreground">
              <div className="text-sm mb-2">No messages yet</div>
              <div className="text-xs">Start a conversation with {conversation.user.name}</div>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {Object.entries(messageGroups).map(([date, dateMessages]) => (
              <div key={date}>
                {/* Date Separator */}
                <div className="flex items-center gap-3 my-4">
                  <Separator className="flex-1" />
                  <span className="text-xs text-muted-foreground bg-background px-2">
                    {formatDate(date)}
                  </span>
                  <Separator className="flex-1" />
                </div>

                {/* Messages for this date */}
                <div className="space-y-3">
                  {dateMessages.map((message) => {
                    const isCurrentUser = message.sender_id === currentUser.id;
                    const canEdit = isCurrentUser && !message.is_edited && 
                      (new Date().getTime() - new Date(message.created_at).getTime()) < 15 * 60 * 1000; // 15 minutes

                    return (
                      <div
                        key={message.id}
                        className={cn(
                          "flex gap-3",
                          isCurrentUser ? "flex-row-reverse" : "flex-row"
                        )}
                      >
                        <Avatar className="h-8 w-8 flex-shrink-0">
                          <AvatarImage src={`/avatars/${message.sender?.id}.jpg`} />
                          <AvatarFallback className="text-xs">
                            {getInitials(message.sender?.name || "")}
                          </AvatarFallback>
                        </Avatar>

                        <div className={cn("flex-1 max-w-[70%]", isCurrentUser && "text-right")}>
                          {/* Reply indicator */}
                          {message.parent_message && (
                            <div className="text-xs text-muted-foreground mb-1 opacity-70">
                              Replying to: {message.parent_message.content.substring(0, 50)}...
                            </div>
                          )}

                          <div
                            className={cn(
                              "inline-block rounded-lg px-3 py-2 text-sm",
                              isCurrentUser
                                ? "bg-primary text-primary-foreground"
                                : "bg-muted"
                            )}
                          >
                            {editingMessage === message.id ? (
                              <div className="space-y-2 min-w-48">
                                <Textarea
                                  value={editContent}
                                  onChange={(e) => setEditContent(e.target.value)}
                                  className="resize-none"
                                  rows={2}
                                />
                                <div className="flex gap-2">
                                  <Button
                                    size="sm"
                                    onClick={() => handleEditMessage(message.id)}
                                  >
                                    Save
                                  </Button>
                                  <Button
                                    size="sm"
                                    variant="outline"
                                    onClick={cancelEditing}
                                  >
                                    Cancel
                                  </Button>
                                </div>
                              </div>
                            ) : (
                              <>
                                <div>{message.content}</div>
                                {message.is_edited && (
                                  <div className="text-xs opacity-70 mt-1">(edited)</div>
                                )}
                              </>
                            )}
                          </div>

                          <div className={cn(
                            "flex items-center gap-2 mt-1 text-xs text-muted-foreground",
                            isCurrentUser ? "justify-end" : "justify-start"
                          )}>
                            <span>{formatTime(message.created_at)}</span>
                            {isCurrentUser && (
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="icon" className="h-4 w-4">
                                    <MoreHorizontal className="h-3 w-3" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem onClick={() => startReply(message)}>
                                    <Reply className="h-4 w-4 mr-2" />
                                    Reply
                                  </DropdownMenuItem>
                                  {canEdit && (
                                    <DropdownMenuItem onClick={() => startEditing(message)}>
                                      <Edit3 className="h-4 w-4 mr-2" />
                                      Edit
                                    </DropdownMenuItem>
                                  )}
                                  <DropdownMenuSeparator />
                                  <DropdownMenuItem
                                    onClick={() => handleDeleteMessage(message.id)}
                                    className="text-destructive"
                                  >
                                    <Trash2 className="h-4 w-4 mr-2" />
                                    Delete
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            )}
                            {!isCurrentUser && (
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-4 w-4"
                                onClick={() => startReply(message)}
                              >
                                <Reply className="h-3 w-3" />
                              </Button>
                            )}
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}
      </ScrollArea>

      {/* Message Input */}
      <div className="border-t p-4">
        {replyingTo && (
          <div className="flex items-center gap-2 mb-3 text-sm text-muted-foreground">
            <Reply className="h-4 w-4" />
            <span>Replying to {replyingTo.sender?.name}</span>
            <Button variant="ghost" size="icon" className="h-4 w-4 ml-auto" onClick={cancelReply}>
              ×
            </Button>
          </div>
        )}
        
        {error && (
          <div className="mb-3 text-sm text-destructive">
            {error}
            <Button variant="ghost" size="sm" className="ml-2" onClick={clearError}>
              Dismiss
            </Button>
          </div>
        )}

        <div className="flex gap-2">
          <div className="flex-1">
            <UserMentionInput
              placeholder={`Message ${conversation.user.name}...`}
              value={newMessage}
              onChange={setNewMessage}
              onUserSelect={(user) => {
                console.log(`Mentioned user: ${user.name}`);
              }}
              onKeyDown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  handleSendMessage();
                }
              }}
              disabled={loading}
            />
          </div>
          <Button
            onClick={handleSendMessage}
            disabled={!newMessage.trim() || loading}
            size="icon"
          >
            <Send className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}