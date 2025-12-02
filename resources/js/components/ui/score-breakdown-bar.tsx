import React from 'react';
import { cn } from '@/lib/utils';
import { CheckCircle2, Users, DollarSign } from 'lucide-react';

interface ScoreBreakdown {
  formalisation: number; // Actual score from 7 checkboxes (0-70)
  teamStructure: number; // Actual score from team data (0-20)
  financialData: number; // Actual score from financial info (0-10)
}

interface ScoreBreakdownBarProps {
  breakdown: ScoreBreakdown;
  className?: string;
}

export function ScoreBreakdownBar({ breakdown, className }: ScoreBreakdownBarProps) {
  const totalScore = breakdown.formalisation + breakdown.teamStructure + breakdown.financialData;

  // Calculate percentages for the bar (out of 100 total)
  const formalisationPercent = breakdown.formalisation;
  const teamPercent = breakdown.teamStructure;
  const financialPercent = breakdown.financialData;

  // Get color based on total score
  const getOverallColor = () => {
    if (totalScore >= 67) return 'text-emerald-600';
    if (totalScore >= 34) return 'text-amber-500';
    return 'text-red-500';
  };

  const getOverallBadgeColor = () => {
    if (totalScore >= 67) return 'bg-emerald-600 text-white';
    if (totalScore >= 34) return 'bg-amber-500 text-white';
    return 'bg-red-500 text-white';
  };

  const sections = [
    {
      label: 'Formalisation Checkboxes',
      score: breakdown.formalisation,
      maxScore: 70,
      color: 'bg-blue-500',
      lightColor: 'bg-blue-100 dark:bg-blue-950',
      textColor: 'text-blue-700 dark:text-blue-300',
      icon: CheckCircle2,
      description: '7 compliance indicators, 10 points each'
    },
    {
      label: 'Team Structure',
      score: breakdown.teamStructure,
      maxScore: 20,
      color: 'bg-purple-500',
      lightColor: 'bg-purple-100 dark:bg-purple-950',
      textColor: 'text-purple-700 dark:text-purple-300',
      icon: Users,
      description: 'Owner, members, employees, team size'
    },
    {
      label: 'Financial Data',
      score: breakdown.financialData,
      maxScore: 10,
      color: 'bg-green-500',
      lightColor: 'bg-green-100 dark:bg-green-950',
      textColor: 'text-green-700 dark:text-green-300',
      icon: DollarSign,
      description: 'Annual turnover and estimated assets'
    }
  ];

  return (
    <div className={cn('space-y-6', className)}>
      {/* Overall Score Header */}
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-2xl font-bold text-foreground">Formalisation Score</h3>
          <p className="text-sm text-muted-foreground mt-1">
            Overall business formalisation assessment
          </p>
        </div>
        <div className="text-right">
          <div className={cn('text-5xl font-bold tabular-nums', getOverallColor())}>
            {totalScore}
          </div>
          <div className="text-sm text-muted-foreground mt-1">out of 100</div>
        </div>
      </div>

      {/* Main Progress Bar */}
      <div className="space-y-3">
        <div className="relative h-12 w-full overflow-hidden rounded-lg bg-slate-100 dark:bg-slate-800">
          {/* Formalisation segment */}
          {breakdown.formalisation > 0 && (
            <div
              className="absolute left-0 top-0 h-full bg-blue-500 transition-all duration-500 flex items-center justify-center"
              style={{ width: `${formalisationPercent}%` }}
            >
              {formalisationPercent >= 8 && (
                <span className="text-xs font-semibold text-white">
                  {breakdown.formalisation}
                </span>
              )}
            </div>
          )}

          {/* Team Structure segment */}
          {breakdown.teamStructure > 0 && (
            <div
              className="absolute top-0 h-full bg-purple-500 transition-all duration-500 flex items-center justify-center"
              style={{
                left: `${formalisationPercent}%`,
                width: `${teamPercent}%`
              }}
            >
              {teamPercent >= 5 && (
                <span className="text-xs font-semibold text-white">
                  {breakdown.teamStructure}
                </span>
              )}
            </div>
          )}

          {/* Financial Data segment */}
          {breakdown.financialData > 0 && (
            <div
              className="absolute top-0 h-full bg-green-500 transition-all duration-500 flex items-center justify-center"
              style={{
                left: `${formalisationPercent + teamPercent}%`,
                width: `${financialPercent}%`
              }}
            >
              {financialPercent >= 3 && (
                <span className="text-xs font-semibold text-white">
                  {breakdown.financialData}
                </span>
              )}
            </div>
          )}

          {/* Score marker overlay */}
          {totalScore === 0 && (
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-sm text-muted-foreground">No score data</span>
            </div>
          )}
        </div>

        {/* Score range labels */}
        <div className="flex justify-between text-xs text-muted-foreground px-1">
          <span>0</span>
          <span className="text-red-500">0-33 Low</span>
          <span className="text-amber-500">34-66 Medium</span>
          <span className="text-emerald-600">67-100 High</span>
          <span>100</span>
        </div>
      </div>

      {/* Breakdown Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {sections.map((section) => {
          const Icon = section.icon;
          const percentage = (section.score / section.maxScore) * 100;

          return (
            <div
              key={section.label}
              className={cn(
                'rounded-lg border p-4 transition-all hover:shadow-md',
                section.lightColor,
                'border-slate-200 dark:border-slate-700'
              )}
            >
              <div className="flex items-start justify-between mb-3">
                <div className={cn('p-2 rounded-lg', section.color)}>
                  <Icon className="h-4 w-4 text-white" />
                </div>
                <div className="text-right">
                  <div className={cn('text-2xl font-bold tabular-nums', section.textColor)}>
                    {section.score}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    / {section.maxScore}
                  </div>
                </div>
              </div>

              <h4 className="font-semibold text-sm text-foreground mb-1">
                {section.label}
              </h4>
              <p className="text-xs text-muted-foreground mb-3">
                {section.description}
              </p>

              {/* Mini progress bar */}
              <div className="h-2 w-full overflow-hidden rounded-full bg-slate-200 dark:bg-slate-700">
                <div
                  className={cn('h-full transition-all duration-500', section.color)}
                  style={{ width: `${percentage}%` }}
                />
              </div>
              <div className="mt-1 text-xs text-muted-foreground text-right">
                {Math.round(percentage)}%
              </div>
            </div>
          );
        })}
      </div>

      {/* Status Badge */}
      <div className="flex items-center justify-center">
        <div className={cn(
          'inline-flex items-center gap-2 px-6 py-3 rounded-full font-semibold',
          getOverallBadgeColor()
        )}>
          {totalScore >= 67 ? (
            <>
              <CheckCircle2 className="h-5 w-5" />
              <span>Highly Formalised</span>
            </>
          ) : totalScore >= 34 ? (
            <>
              <CheckCircle2 className="h-5 w-5" />
              <span>Moderately Formalised</span>
            </>
          ) : (
            <>
              <CheckCircle2 className="h-5 w-5" />
              <span>Low Formalisation</span>
            </>
          )}
        </div>
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
  if (formalisation) {
    if (formalisation.hasBankAccount) formalisationScore += 10;
    if (formalisation.hasTaxClarification) formalisationScore += 10;
    if (formalisation.isRegisteredForVat) formalisationScore += 10;
    if (formalisation.isMemberOfAssociation) formalisationScore += 10;
    if (formalisation.isAffiliated) formalisationScore += 10;
    if (formalisation.hasExportLicense) formalisationScore += 10;
    if (formalisation.hasAccessedBds) formalisationScore += 10;
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
  if (employeeSummary) {
    const fullTimeEmployees = (employeeSummary.fullTimeMales || 0) + (employeeSummary.fullTimeFemales || 0);
    if (fullTimeEmployees > 0) {
      teamScore += 5;
    }

    // Team size scoring: 5 points max
    let totalTeamSize = 0;
    if (hasPrimaryOwner) totalTeamSize += 1;
    totalTeamSize += additionalMembersCount;
    totalTeamSize += (employeeSummary.fullTimeMales || 0) + (employeeSummary.fullTimeFemales || 0);
    totalTeamSize += (employeeSummary.partTimeMales || 0) + (employeeSummary.partTimeFemales || 0);
    totalTeamSize += (employeeSummary.internMales || 0) + (employeeSummary.internFemales || 0);

    if (totalTeamSize >= 10) {
      teamScore += 5;
    } else if (totalTeamSize >= 5) {
      teamScore += 3;
    } else if (totalTeamSize >= 2) {
      teamScore += 1;
    }
  }

  // Financial Data (10 points max)
  if (formalisation) {
    // Annual turnover: 5 points
    if (formalisation.annualTurnover && formalisation.annualTurnover > 0) {
      financialScore += 5;
    }

    // Estimated assets: 5 points
    if (formalisation.estimatedValueOfAssets && formalisation.estimatedValueOfAssets > 0) {
      financialScore += 5;
    }
  }

  return {
    formalisation: formalisationScore,
    teamStructure: teamScore,
    financialData: financialScore
  };
}
