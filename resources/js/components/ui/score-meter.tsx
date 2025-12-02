import React from 'react';
import { cn } from '@/lib/utils';

interface ScoreMeterProps {
  score: number;
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
  className?: string;
}

export function ScoreMeter({
  score,
  size = 'sm',
  showLabel = true,
  className
}: ScoreMeterProps) {
  // Clamp score between 0-100
  const clampedScore = Math.min(100, Math.max(0, score));

  // Determine color based on score ranges
  const getBarColor = (value: number) => {
    if (value >= 67) return 'bg-emerald-600';
    if (value >= 34) return 'bg-amber-500';
    return 'bg-red-500';
  };

  const getTextColor = (value: number) => {
    if (value >= 67) return 'text-emerald-600';
    if (value >= 34) return 'text-amber-500';
    return 'text-red-500';
  };

  // Size configurations
  const sizeConfig = {
    sm: { width: 'w-16', height: 'h-2', fontSize: 'text-xs', gap: 'gap-1.5' },
    md: { width: 'w-20', height: 'h-2.5', fontSize: 'text-sm', gap: 'gap-2' },
    lg: { width: 'w-24', height: 'h-3', fontSize: 'text-base', gap: 'gap-2' },
  };

  const config = sizeConfig[size];

  return (
    <div className={cn('flex items-center', config.gap, className)}>
      {/* Progress bar */}
      <div className={cn(
        'overflow-hidden rounded-full bg-slate-200 dark:bg-slate-700',
        config.width,
        config.height
      )}>
        <div
          className={cn(
            'h-full rounded-full transition-all duration-500',
            getBarColor(clampedScore)
          )}
          style={{ width: `${clampedScore}%` }}
        />
      </div>

      {/* Score label */}
      {showLabel && (
        <span className={cn(
          'font-semibold tabular-nums',
          config.fontSize,
          getTextColor(clampedScore)
        )}>
          {Math.round(clampedScore)}
        </span>
      )}
    </div>
  );
}
