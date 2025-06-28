import React, { useState, useEffect } from 'react';
import { Search, X, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { SearchBarProps } from '@/types/crud';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

export function SearchBar({ 
  value, 
  onChange, 
  onSearch, 
  placeholder = 'Search...', 
  debounceMs = 300,
  loading = false,
  className 
}: SearchBarProps) {
  const [localValue, setLocalValue] = useState(value);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      if (localValue !== value) {
        onSearch(localValue);
      }
    }, debounceMs);

    return () => clearTimeout(timer);
  }, [localValue, onSearch, debounceMs, value]);

  // Sync with external value changes
  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(localValue);
  };

  const handleClear = () => {
    setLocalValue('');
    onChange('');
    onSearch('');
  };

  return (
    <form onSubmit={handleSubmit} className={cn('relative', className)}>
      <div className="relative">
        <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
          {loading ? (
            <Loader2 className="h-4 w-4 text-gray-400 animate-spin" />
          ) : (
            <Search className="h-4 w-4 text-gray-400" />
          )}
        </div>
        
        <Input
          type="text"
          value={localValue}
          onChange={(e) => {
            setLocalValue(e.target.value);
            onChange(e.target.value);
          }}
          className={cn(
            'pl-10',
            localValue && 'pr-10'
          )}
          placeholder={placeholder}
          disabled={loading}
        />
        
        {localValue && !loading && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={handleClear}
            className="absolute inset-y-0 right-0 h-full px-3 py-0 hover:bg-transparent"
          >
            <X className="h-4 w-4 text-gray-400 hover:text-gray-600" />
            <span className="sr-only">Clear search</span>
          </Button>
        )}
      </div>
      
      {/* Submit button for accessibility - hidden but functional */}
      <button type="submit" className="sr-only">
        Search
      </button>
    </form>
  );
}