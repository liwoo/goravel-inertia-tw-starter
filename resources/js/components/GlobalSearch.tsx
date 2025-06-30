import React, { useState, useEffect, useCallback } from 'react';
import { router } from '@inertiajs/react';
import { 
  Search, 
  BookOpen, 
  Users, 
  Shield,
  FileText,
  ChevronRight,
  Command,
  Loader2
} from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { usePermissions } from '@/contexts/PermissionsContext';
import { useDebounce } from '@/hooks/useDebounce';
import axios from '@/lib/axios';

interface SearchResult {
  id: number;
  title: string;
  subtitle?: string;
  type: 'book' | 'user' | 'role' | 'permission';
  url: string;
}

interface GlobalSearchProps {
  isOpen: boolean;
  onClose: () => void;
}

export function GlobalSearch({ isOpen, onClose }: GlobalSearchProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const { canPerformAction } = usePermissions();
  
  const debouncedSearchTerm = useDebounce(searchTerm, 300);

  // Define available search categories based on permissions
  const searchableEntities = React.useMemo(() => {
    const entities = [];
    
    if (canPerformAction('books', 'read')) {
      entities.push({ type: 'book', label: 'Books', icon: <BookOpen className="h-4 w-4" /> });
    }
    
    if (canPerformAction('users', 'read')) {
      entities.push({ type: 'user', label: 'Users', icon: <Users className="h-4 w-4" /> });
    }
    
    if (canPerformAction('roles', 'read') || canPerformAction('permissions', 'read')) {
      entities.push({ type: 'role', label: 'Roles', icon: <Shield className="h-4 w-4" /> });
      entities.push({ type: 'permission', label: 'Permissions', icon: <Shield className="h-4 w-4" /> });
    }
    
    return entities;
  }, [canPerformAction]);

  // Real search function using API
  const performSearch = useCallback(async (term: string) => {
    if (!term.trim()) {
      setResults([]);
      return;
    }

    setIsLoading(true);
    
    try {
      const response = await axios.get('/api/search', {
        params: { q: term },
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        }
      });
      
      if (response.data && response.data.results) {
        setResults(response.data.results);
      } else {
        setResults([]);
      }
    } catch (error) {
      console.error('Search error:', error);
      setResults([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Perform search when debounced term changes
  useEffect(() => {
    performSearch(debouncedSearchTerm);
  }, [debouncedSearchTerm, performSearch]);

  // Reset state when dialog closes
  useEffect(() => {
    if (!isOpen) {
      setSearchTerm('');
      setResults([]);
      setSelectedIndex(0);
    }
  }, [isOpen]);

  // Handle keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isOpen) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex(prev => (prev + 1) % results.length);
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex(prev => (prev - 1 + results.length) % results.length);
          break;
        case 'Enter':
          e.preventDefault();
          if (results[selectedIndex]) {
            router.visit(results[selectedIndex].url);
            onClose();
          }
          break;
        case 'Escape':
          e.preventDefault();
          onClose();
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, results, selectedIndex, onClose]);

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'book':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
      case 'user':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'role':
      case 'permission':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400';
    }
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'book':
        return <BookOpen className="h-4 w-4" />;
      case 'user':
        return <Users className="h-4 w-4" />;
      case 'role':
      case 'permission':
        return <Shield className="h-4 w-4" />;
      default:
        return <FileText className="h-4 w-4" />;
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[600px] p-0">
        <div className="flex items-center border-b px-4 py-3">
          <Search className="mr-3 h-5 w-5 text-muted-foreground" />
          <Input
            placeholder="Search for books, users, roles..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="flex-1 border-0 bg-transparent p-0 text-base placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-0"
            autoFocus
          />
          <kbd className="pointer-events-none ml-2 hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:inline-flex">
            <span className="text-xs">ESC</span>
          </kbd>
        </div>
        
        <div className="max-h-[400px] overflow-y-auto">
          {isLoading && (
            <div className="flex items-center justify-center p-8">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
              <span className="ml-2 text-sm text-muted-foreground">Searching...</span>
            </div>
          )}
          
          {!isLoading && searchTerm && results.length === 0 && (
            <div className="p-8 text-center">
              <p className="text-sm text-muted-foreground">
                No results found for "{searchTerm}"
              </p>
              <p className="mt-2 text-xs text-muted-foreground">
                Try searching in specific sections:
              </p>
              <div className="mt-4 flex flex-wrap justify-center gap-2">
                {searchableEntities.map((entity) => (
                  <Badge
                    key={entity.type}
                    variant="secondary"
                    className="cursor-pointer"
                    onClick={() => {
                      router.visit(`/admin/${entity.type}s?search=${encodeURIComponent(searchTerm)}`);
                      onClose();
                    }}
                  >
                    {entity.icon}
                    <span className="ml-1">{entity.label}</span>
                  </Badge>
                ))}
              </div>
            </div>
          )}
          
          {!isLoading && results.length > 0 && (
            <div className="py-2">
              {results.map((result, index) => (
                <button
                  key={result.id}
                  className={cn(
                    "flex w-full items-center justify-between px-4 py-2 text-left text-sm hover:bg-accent",
                    selectedIndex === index && "bg-accent"
                  )}
                  onClick={() => {
                    router.visit(result.url);
                    onClose();
                  }}
                  onMouseEnter={() => setSelectedIndex(index)}
                >
                  <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
                      {getTypeIcon(result.type)}
                    </div>
                    <div className="flex flex-col">
                      <span className="font-medium">{result.title}</span>
                      {result.subtitle && (
                        <span className="text-xs text-muted-foreground">
                          {result.subtitle}
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className={cn("text-xs", getTypeColor(result.type))}>
                      {result.type}
                    </Badge>
                    <ChevronRight className="h-4 w-4 text-muted-foreground" />
                  </div>
                </button>
              ))}
            </div>
          )}
          
          {!searchTerm && (
            <div className="p-4">
              <p className="mb-3 text-xs font-medium text-muted-foreground">
                QUICK ACCESS
              </p>
              <div className="space-y-1">
                {searchableEntities.map((entity) => (
                  <button
                    key={entity.type}
                    className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm hover:bg-accent"
                    onClick={() => {
                      router.visit(`/admin/${entity.type}s`);
                      onClose();
                    }}
                  >
                    <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
                      {entity.icon}
                    </div>
                    <span>Browse {entity.label}</span>
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
        
        <div className="border-t p-3">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <div className="flex items-center gap-4">
              <span className="flex items-center gap-1">
                <kbd className="rounded border bg-muted px-1 font-mono text-[10px]">↑↓</kbd>
                Navigate
              </span>
              <span className="flex items-center gap-1">
                <kbd className="rounded border bg-muted px-1 font-mono text-[10px]">↵</kbd>
                Open
              </span>
            </div>
            <span className="flex items-center gap-1">
              Press
              <kbd className="rounded border bg-muted px-1 font-mono text-[10px]">
                <Command className="h-3 w-3" />
              </kbd>
              <kbd className="rounded border bg-muted px-1 font-mono text-[10px]">K</kbd>
              to open
            </span>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}