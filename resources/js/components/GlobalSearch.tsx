import React, { useState, useEffect, useCallback } from 'react';
import { router } from '@inertiajs/react';
import {
  Search,
  ChevronRight,
  Command,
  Loader2,
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
import {
  SEARCH_ENTITIES,
  getEntityColors,
  getEntityIcon,
  type SearchEntityType,
} from '@/config/search_config';

interface SearchResult {
  id: number;
  title: string;
  subtitle?: string;
  type: SearchEntityType;
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
    return SEARCH_ENTITIES.filter(entity =>
      canPerformAction(entity.permissionService, entity.permissionAction)
    );
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

  // Generate search placeholder dynamically from available entities
  const searchPlaceholder = React.useMemo(() => {
    const labels = searchableEntities.map(e => e.label.toLowerCase()).slice(0, 3);
    return `Search for ${labels.join(', ')}...`;
  }, [searchableEntities]);

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[600px] p-0">
        <div className="flex items-center border-b px-4 py-3">
          <Search className="mr-3 h-5 w-5 text-muted-foreground" />
          <Input
            placeholder={searchPlaceholder}
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
                      router.visit(`${entity.urlPrefix}?search=${encodeURIComponent(searchTerm)}`);
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
                      {getEntityIcon(result.type)}
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
                    <Badge variant="secondary" className={cn("text-xs", getEntityColors(result.type))}>
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
                      router.visit(entity.urlPrefix);
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