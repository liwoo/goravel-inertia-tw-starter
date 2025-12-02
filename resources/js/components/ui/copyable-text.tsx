import React, { useState } from 'react';
import { cn } from '@/lib/utils';
import { Copy, Check } from 'lucide-react';

interface CopyableTextProps {
  value: string | null | undefined;
  className?: string;
  iconSize?: 'sm' | 'md';
}

export function CopyableText({ value, className, iconSize = 'sm' }: CopyableTextProps) {
  const [copied, setCopied] = useState(false);

  if (!value) {
    return <span className={cn('text-muted-foreground', className)}>-</span>;
  }

  const handleCopy = async (e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(value);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const iconClasses = iconSize === 'sm' ? 'h-3 w-3' : 'h-4 w-4';

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 group cursor-pointer hover:text-primary transition-colors',
        className
      )}
      onClick={handleCopy}
      title="Click to copy"
    >
      <span>{value}</span>
      {copied ? (
        <Check className={cn(iconClasses, 'text-green-600')} />
      ) : (
        <Copy className={cn(iconClasses, 'opacity-0 group-hover:opacity-50 transition-opacity')} />
      )}
    </span>
  );
}
