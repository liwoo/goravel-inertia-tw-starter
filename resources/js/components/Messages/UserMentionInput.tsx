"use client";

import * as React from "react";
import { Search, AtSign, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";
import { useMessages } from "@/contexts/MessageContext";
import { useDebounce } from "@/hooks/useDebounce";

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

interface UserMentionInputProps {
  value: string;
  onChange: (value: string) => void;
  onUserSelect?: (user: User) => void;
  onKeyDown?: (e: React.KeyboardEvent<HTMLInputElement>) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

export function UserMentionInput({
  value,
  onChange,
  onUserSelect,
  onKeyDown,
  placeholder = "Type @ to mention users...",
  className,
  disabled
}: UserMentionInputProps) {
  const [isOpen, setIsOpen] = React.useState(false);
  const [searchTerm, setSearchTerm] = React.useState("");
  const [mentionUsers, setMentionUsers] = React.useState<User[]>([]);
  const [caretPosition, setCaretPosition] = React.useState(0);
  const [mentionStart, setMentionStart] = React.useState(-1);
  const [loading, setLoading] = React.useState(false);
  
  const inputRef = React.useRef<HTMLInputElement>(null);
  const { searchUsers } = useMessages();
  
  const debouncedSearchTerm = useDebounce(searchTerm, 300);

  // Search for users when term changes
  React.useEffect(() => {
    const performSearch = async () => {
      if (debouncedSearchTerm.length >= 2) {
        setLoading(true);
        try {
          const users = await searchUsers(debouncedSearchTerm);
          setMentionUsers(users);
        } catch (error) {
          console.error("Failed to search users:", error);
          setMentionUsers([]);
        } finally {
          setLoading(false);
        }
      } else {
        setMentionUsers([]);
      }
    };

    if (isOpen && debouncedSearchTerm) {
      performSearch();
    }
  }, [debouncedSearchTerm, isOpen, searchUsers]);

  // Handle input changes and detect @ mentions
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    const cursorPos = e.target.selectionStart || 0;
    
    onChange(newValue);
    setCaretPosition(cursorPos);

    // Check for @ mention
    const textBeforeCursor = newValue.substring(0, cursorPos);
    const atIndex = textBeforeCursor.lastIndexOf('@');
    
    if (atIndex !== -1) {
      const afterAt = textBeforeCursor.substring(atIndex + 1);
      
      // Check if there's a space after @, if so, close mentions
      if (afterAt.includes(' ')) {
        setIsOpen(false);
        setMentionStart(-1);
        return;
      }
      
      // Start mention search
      setMentionStart(atIndex);
      setSearchTerm(afterAt);
      setIsOpen(true);
    } else {
      setIsOpen(false);
      setMentionStart(-1);
    }
  };

  // Handle user selection
  const handleUserSelect = (user: User) => {
    if (mentionStart === -1) return;

    const beforeMention = value.substring(0, mentionStart);
    const afterCursor = value.substring(caretPosition);
    const newValue = `${beforeMention}@${user.name} ${afterCursor}`;
    
    onChange(newValue);
    setIsOpen(false);
    setMentionStart(-1);
    setSearchTerm("");
    
    // Focus back to input and position cursor after the mention
    setTimeout(() => {
      if (inputRef.current) {
        const newCursorPos = beforeMention.length + user.name.length + 2; // +2 for @ and space
        inputRef.current.setSelectionRange(newCursorPos, newCursorPos);
        inputRef.current.focus();
      }
    }, 0);

    onUserSelect?.(user);
  };

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (isOpen) {
      if (e.key === 'Escape') {
        setIsOpen(false);
        setMentionStart(-1);
        return;
      }
    }
    
    // Call the external onKeyDown handler if provided
    onKeyDown?.(e);
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map(word => word[0])
      .join("")
      .toUpperCase()
      .substring(0, 2);
  };

  // Extract mentions from current value
  const extractMentions = (text: string): string[] => {
    const mentionRegex = /@(\w+)/g;
    const mentions = [];
    let match;
    
    while ((match = mentionRegex.exec(text)) !== null) {
      mentions.push(match[1]);
    }
    
    return mentions;
  };

  const currentMentions = extractMentions(value);

  return (
    <div className={cn("relative", className)}>
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <div className="relative">
            <Input
              ref={inputRef}
              value={value}
              onChange={handleInputChange}
              onKeyDown={handleKeyDown}
              placeholder={placeholder}
              disabled={disabled}
              className="pr-10"
            />
            <div className="absolute right-3 top-1/2 -translate-y-1/2">
              <AtSign className="h-4 w-4 text-muted-foreground" />
            </div>
          </div>
        </PopoverTrigger>
        
        <PopoverContent 
          className="p-0 w-80" 
          align="start"
          side="bottom"
          sideOffset={4}
        >
          <Command>
            <CommandInput 
              placeholder="Search users..." 
              value={searchTerm}
              onValueChange={setSearchTerm}
            />
            <CommandList>
              {loading ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  Searching users...
                </div>
              ) : mentionUsers.length === 0 ? (
                <CommandEmpty>
                  {searchTerm.length < 2 
                    ? "Type at least 2 characters to search" 
                    : "No users found"}
                </CommandEmpty>
              ) : (
                <CommandGroup heading="Users">
                  {mentionUsers.map((user) => (
                    <CommandItem
                      key={user.id}
                      value={user.name}
                      onSelect={() => handleUserSelect(user)}
                      className="flex items-center gap-3 p-3"
                    >
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={`/avatars/${user.id}.jpg`} />
                        <AvatarFallback className="text-xs">
                          {getInitials(user.name)}
                        </AvatarFallback>
                      </Avatar>
                      
                      <div className="flex-1 min-w-0">
                        <div className="font-medium truncate">{user.name}</div>
                        <div className="text-xs text-muted-foreground truncate">
                          {user.email}
                        </div>
                        {user.roles && user.roles.length > 0 && (
                          <div className="flex gap-1 mt-1">
                            {user.roles.slice(0, 2).map((role) => (
                              <Badge key={role.id} variant="outline" className="text-xs h-4 px-1">
                                {role.name}
                              </Badge>
                            ))}
                            {user.roles.length > 2 && (
                              <Badge variant="outline" className="text-xs h-4 px-1">
                                +{user.roles.length - 2}
                              </Badge>
                            )}
                          </div>
                        )}
                      </div>
                      
                      <div className="flex items-center gap-1">
                        {user.is_active && (
                          <div className="h-2 w-2 rounded-full bg-green-500" />
                        )}
                        <AtSign className="h-4 w-4 text-muted-foreground" />
                      </div>
                    </CommandItem>
                  ))}
                </CommandGroup>
              )}
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {/* Show current mentions */}
      {currentMentions.length > 0 && (
        <div className="flex flex-wrap gap-1 mt-2">
          <span className="text-xs text-muted-foreground">Mentioning:</span>
          {currentMentions.map((mention, index) => (
            <Badge key={index} variant="secondary" className="text-xs">
              @{mention}
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
}

