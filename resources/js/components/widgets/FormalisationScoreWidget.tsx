"use client";

import * as React from "react";
import { Shield, Pencil, Info } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
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
  // Edit functionality props
  canEdit?: boolean;
  onEditClick?: () => void;
}

export function FormalisationScoreWidget({
  formalisation,
  primaryOwner,
  additionalMembers = [],
  employeeSummary,
  isLoading = false,
  canEdit = false,
  onEditClick,
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
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5" />
              Business Formalisation Progress
              <Popover>
                <PopoverTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-5 w-5 rounded-full p-0 hover:bg-muted">
                    <Info className="h-4 w-4 text-muted-foreground" />
                    <span className="sr-only">How is this score calculated?</span>
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-96" align="start">
                  <div className="space-y-4">
                    <div>
                      <h4 className="font-semibold text-sm">How is this score calculated?</h4>
                      <p className="text-xs text-muted-foreground mt-1">
                        Your formalisation score (0-100) reflects how well-established and compliant your business is.
                      </p>
                    </div>

                    <div className="space-y-3 text-xs">
                      <div>
                        <div className="flex items-center gap-2 font-medium">
                          <div className="h-2 w-2 rounded-full bg-blue-500" />
                          Compliance (up to 70 points)
                        </div>
                        <p className="text-muted-foreground ml-4 mt-1">
                          10 points each for: Bank Account, Tax Clarification, VAT Registration,
                          Association Membership, Business Affiliation, Export License, and BDS Access.
                        </p>
                      </div>

                      <div>
                        <div className="flex items-center gap-2 font-medium">
                          <div className="h-2 w-2 rounded-full bg-amber-500" />
                          Team Structure (up to 20 points)
                        </div>
                        <p className="text-muted-foreground ml-4 mt-1">
                          Points for having a registered owner, team members, employees, and overall team size.
                        </p>
                      </div>

                      <div>
                        <div className="flex items-center gap-2 font-medium">
                          <div className="h-2 w-2 rounded-full bg-emerald-500" />
                          Financial Data (up to 10 points)
                        </div>
                        <p className="text-muted-foreground ml-4 mt-1">
                          5 points each for reporting annual turnover and estimated asset value.
                        </p>
                      </div>
                    </div>

                    <div className="border-t pt-3">
                      <h4 className="font-semibold text-sm">Why does this matter?</h4>
                      <p className="text-xs text-muted-foreground mt-1">
                        A higher formalisation score improves your visibility and eligibility for:
                      </p>
                      <ul className="text-xs text-muted-foreground mt-2 ml-4 list-disc space-y-1">
                        <li>Government procurement opportunities</li>
                        <li>Business development support programs</li>
                        <li>Access to financing and credit facilities</li>
                        <li>Partnership opportunities with larger enterprises</li>
                        <li>Export and trade facilitation programs</li>
                      </ul>
                    </div>
                  </div>
                </PopoverContent>
              </Popover>
            </CardTitle>
            <CardDescription>Your business compliance and development score</CardDescription>
          </div>
          {canEdit && onEditClick && (
            <Button
              variant="outline"
              size="sm"
              onClick={onEditClick}
              className="flex items-center gap-2"
            >
              <Pencil className="h-4 w-4" />
              Edit Details
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <ScoreBreakdownBar breakdown={breakdown} score={score} />
      </CardContent>
    </Card>
  );
}

export default FormalisationScoreWidget;
