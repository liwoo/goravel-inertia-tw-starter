import React from 'react';
import { cn } from '@/lib/utils';

interface ScoreBreakdownBarProps {
  /** Compliance score from backend (0-70) */
  complianceScore: number;
  /** Team structure score from backend (0-20) */
  teamStructureScore: number;
  /** Financial score from backend (0-10) */
  financialScore: number;
  /** Total formalisation score from backend (0-100) */
  totalScore: number;
  className?: string;
}

export function ScoreBreakdownBar({
  complianceScore,
  teamStructureScore,
  financialScore,
  totalScore,
  className
}: ScoreBreakdownBarProps) {

  const categories = [
    {
      label: 'Compliance',
      score: complianceScore,
      maxScore: 70,
      color: 'bg-blue-500',
      dotColor: 'bg-blue-500',
    },
    {
      label: 'Team Structure',
      score: teamStructureScore,
      maxScore: 20,
      color: 'bg-amber-500',
      dotColor: 'bg-amber-500',
    },
    {
      label: 'Financial',
      score: financialScore,
      maxScore: 10,
      color: 'bg-emerald-500',
      dotColor: 'bg-emerald-500',
    },
  ];

  // Calculate cumulative positions for the segmented bar
  let cumulativePercent = 0;

  return (
    <div className={cn('space-y-4', className)}>
      {/* Header with title and total score */}
      <div className="flex items-baseline justify-between">
        <h3 className="text-sm font-medium text-muted-foreground">Formalisation Score</h3>
      </div>

      {/* Large Score Display */}
      <div className="flex items-baseline gap-2">
        <span className="text-5xl font-bold tabular-nums text-foreground">{totalScore}</span>
        <span className="text-lg text-muted-foreground">/ 100</span>
      </div>

      {/* Segmented Progress Bar */}
      <div className="relative h-3 w-full overflow-hidden rounded-full bg-slate-200 dark:bg-slate-700">
        {categories.map((category, index) => {
          const width = category.score;
          const left = cumulativePercent;
          cumulativePercent += category.score;

          if (width === 0) return null;

          return (
            <div
              key={category.label}
              className={cn('absolute top-0 h-full transition-all duration-500', category.color)}
              style={{
                left: `${left}%`,
                width: `${width}%`,
              }}
            />
          );
        })}
      </div>

      {/* Category Breakdown - 3 columns */}
      <div className="grid grid-cols-3 gap-4 pt-2">
        {categories.map((category) => (
          <div key={category.label} className="space-y-1">
            <div className="flex items-center gap-2">
              <div className={cn('h-2.5 w-2.5 rounded-full', category.dotColor)} />
              <span className="text-sm text-muted-foreground">{category.label}</span>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-2xl font-bold tabular-nums text-foreground">{category.score}</span>
              <span className="text-sm text-muted-foreground">/ {category.maxScore}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
