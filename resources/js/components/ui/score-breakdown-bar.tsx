import React from 'react';
import { cn } from '@/lib/utils';

interface ScoreBreakdown {
  formalisation: number; // Actual score from 7 checkboxes (0-70)
  teamStructure: number; // Actual score from team data (0-20)
  financialData: number; // Actual score from financial info (0-10)
}

interface ScoreBreakdownBarProps {
  breakdown: ScoreBreakdown;
  className?: string;
  /** When provided, uses this score from API instead of calculating from breakdown */
  score?: number;
}

export function ScoreBreakdownBar({ breakdown, className, score }: ScoreBreakdownBarProps) {
  // Use API score when provided (source of truth), otherwise fall back to calculated
  const totalScore = score ?? (breakdown.formalisation + breakdown.teamStructure + breakdown.financialData);

  const categories = [
    {
      label: 'Compliance',
      score: breakdown.formalisation,
      maxScore: 70,
      color: 'bg-blue-500',
      dotColor: 'bg-blue-500',
    },
    {
      label: 'Team Structure',
      score: breakdown.teamStructure,
      maxScore: 20,
      color: 'bg-amber-500',
      dotColor: 'bg-amber-500',
    },
    {
      label: 'Financial',
      score: breakdown.financialData,
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

// Helper function to calculate the breakdown from formalisation data
export function calculateScoreBreakdown(
  formalisation: any,
  hasPrimaryOwner: boolean,
  additionalMembersCount: number,
  employeeSummary: any
): ScoreBreakdown {
  let formalisationScore = 0;
  let teamScore = 0;
  let financialScore = 0;

  // Formalisation Checkboxes (70 points max)
  // Support both camelCase and snake_case field names
  if (formalisation) {
    if (formalisation.hasBankAccount ?? formalisation.has_bank_account) formalisationScore += 10;
    if (formalisation.hasTaxClarification ?? formalisation.has_tax_clarification) formalisationScore += 10;
    if (formalisation.isRegisteredForVat ?? formalisation.is_registered_for_vat) formalisationScore += 10;
    if (formalisation.isMemberOfAssociation ?? formalisation.is_member_of_association) formalisationScore += 10;
    if (formalisation.isAffiliated ?? formalisation.is_affiliated) formalisationScore += 10;
    if (formalisation.hasExportLicense ?? formalisation.has_export_license) formalisationScore += 10;
    if (formalisation.hasAccessedBds ?? formalisation.has_accessed_bds) formalisationScore += 10;
  }

  // Team Structure (20 points max)
  // Has primary owner: 5 points
  if (hasPrimaryOwner) {
    teamScore += 5;
  }

  // Has additional members: 5 points
  if (additionalMembersCount > 0) {
    teamScore += 5;
  }

  // Has full-time employees: 5 points
  // Support both camelCase and snake_case field names
  if (employeeSummary) {
    const fullTimeMales = employeeSummary.fullTimeMales ?? employeeSummary.full_time_males ?? 0;
    const fullTimeFemales = employeeSummary.fullTimeFemales ?? employeeSummary.full_time_females ?? 0;
    const partTimeMales = employeeSummary.partTimeMales ?? employeeSummary.part_time_males ?? 0;
    const partTimeFemales = employeeSummary.partTimeFemales ?? employeeSummary.part_time_females ?? 0;
    const internMales = employeeSummary.internMales ?? employeeSummary.intern_males ?? 0;
    const internFemales = employeeSummary.internFemales ?? employeeSummary.intern_females ?? 0;

    const fullTimeEmployees = fullTimeMales + fullTimeFemales;
    if (fullTimeEmployees > 0) {
      teamScore += 5;
    }

    // Team size scoring: 5 points max
    let totalTeamSize = 0;
    if (hasPrimaryOwner) totalTeamSize += 1;
    totalTeamSize += additionalMembersCount;
    totalTeamSize += fullTimeMales + fullTimeFemales;
    totalTeamSize += partTimeMales + partTimeFemales;
    totalTeamSize += internMales + internFemales;

    if (totalTeamSize >= 10) {
      teamScore += 5;
    } else if (totalTeamSize >= 5) {
      teamScore += 3;
    } else if (totalTeamSize >= 2) {
      teamScore += 1;
    }
  }

  // Financial Data (10 points max)
  // Support both camelCase and snake_case field names
  if (formalisation) {
    const annualTurnover = formalisation.annualTurnover ?? formalisation.annual_turnover ?? 0;
    const estimatedAssets = formalisation.estimatedValueOfAssets ?? formalisation.estimated_value_of_assets ?? 0;

    // Annual turnover: 5 points
    if (annualTurnover > 0) {
      financialScore += 5;
    }

    // Estimated assets: 5 points
    if (estimatedAssets > 0) {
      financialScore += 5;
    }
  }

  return {
    formalisation: formalisationScore,
    teamStructure: teamScore,
    financialData: financialScore
  };
}
