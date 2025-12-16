"use client";

import * as React from "react";
import { Shield } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ScoreBreakdownBar, calculateScoreBreakdown } from "@/components/ui/score-breakdown-bar";

// Interface for formalisation data
export interface FormalisationData {
  formalisation_score?: number;
  formalisationScore?: number;
  hasBankAccount?: boolean;
  has_bank_account?: boolean;
  hasTaxClarification?: boolean;
  has_tax_clarification?: boolean;
  isRegisteredForVat?: boolean;
  is_registered_for_vat?: boolean;
  isMemberOfAssociation?: boolean;
  is_member_of_association?: boolean;
  isAffiliated?: boolean;
  is_affiliated?: boolean;
  hasExportLicense?: boolean;
  has_export_license?: boolean;
  hasAccessedBds?: boolean;
  has_accessed_bds?: boolean;
  annualTurnover?: number;
  annual_turnover?: number;
  estimatedValueOfAssets?: number;
  estimated_value_of_assets?: number;
}

interface FormalisationScoreWidgetProps {
  formalisation?: FormalisationData | null;
  primaryOwner?: any;
  additionalMembers?: any[];
  employeeSummary?: any;
  isLoading?: boolean;
}

export function FormalisationScoreWidget({
  formalisation,
  primaryOwner,
  additionalMembers = [],
  employeeSummary,
  isLoading = false,
}: FormalisationScoreWidgetProps) {
  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Business Formalisation Progress
          </CardTitle>
          <CardDescription>Your business compliance and development score</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-48 animate-pulse rounded-md bg-muted" />
        </CardContent>
      </Card>
    );
  }

  if (!formalisation) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Business Formalisation Progress
          </CardTitle>
          <CardDescription>Your business compliance and development score</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center py-8">
          <div className="text-center">
            <Shield className="mx-auto h-12 w-12 text-muted-foreground/50" />
            <p className="mt-2 text-sm text-muted-foreground">
              No formalisation data available
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const breakdown = calculateScoreBreakdown(
    formalisation,
    !!primaryOwner,
    additionalMembers.length,
    employeeSummary
  );

  const score = formalisation.formalisation_score ?? formalisation.formalisationScore ?? 0;

  return (
    <Card className="flex flex-col">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Shield className="h-5 w-5" />
          Business Formalisation Progress
        </CardTitle>
        <CardDescription>Your business compliance and development score</CardDescription>
      </CardHeader>
      <CardContent>
        <ScoreBreakdownBar breakdown={breakdown} score={score} />
      </CardContent>
    </Card>
  );
}

export default FormalisationScoreWidget;
