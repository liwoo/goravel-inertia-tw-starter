import { DirectorySME } from "@/types/directory";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Building2, MapPin, Phone, Mail, Globe, Calendar, Info, Lock } from "lucide-react";
import { cn } from "@/lib/utils";

interface DirectorySMECardProps {
  sme: DirectorySME;
  canViewContactDetails?: boolean;
  minScoreRequired?: number;
}

const classificationColors = {
  Micro: "bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-950 dark:text-emerald-400",
  Small: "bg-blue-100 text-blue-800 border-blue-200 dark:bg-blue-950 dark:text-blue-400",
  Medium: "bg-amber-100 text-amber-800 border-amber-200 dark:bg-amber-950 dark:text-amber-400",
};

export function DirectorySMECard({
  sme,
  canViewContactDetails = true,
  minScoreRequired = 70,
}: DirectorySMECardProps) {
  return (
    <TooltipProvider>
      <Card className="h-full transition-all hover:shadow-lg hover:border-primary/20">
        <CardHeader className="space-y-3">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1 min-w-0">
              <h3 className="text-lg font-semibold leading-tight line-clamp-2">
                {sme.name}
              </h3>
              <div className="flex items-center gap-2 mt-1.5 flex-wrap">
                <Badge
                  variant="outline"
                  className="text-xs font-mono"
                >
                  {sme.usme_number}
                </Badge>
                {canViewContactDetails ? (
                  <Badge
                    className={cn(
                      "text-xs border",
                      classificationColors[sme.classification as keyof typeof classificationColors] ||
                      "bg-slate-100 text-slate-800 border-slate-200"
                    )}
                  >
                    {sme.classification}
                  </Badge>
                ) : (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Badge
                        variant="outline"
                        className="text-xs border-dashed cursor-help"
                      >
                        <Lock className="h-3 w-3 mr-1" />
                        Classification Hidden
                      </Badge>
                    </TooltipTrigger>
                    <TooltipContent className="max-w-xs">
                      <p className="text-sm">
                        Complete your business formalisation to achieve a score of {minScoreRequired}+ to view classification details.
                      </p>
                    </TooltipContent>
                  </Tooltip>
                )}
              </div>
            </div>
          </div>
        </CardHeader>

      <CardContent className="space-y-4">
        {/* Sector */}
        <div className="flex items-start gap-2 text-sm">
          <Building2 className="h-4 w-4 mt-0.5 text-muted-foreground flex-shrink-0" />
          <div className="flex-1 min-w-0">
            <p className="font-medium">{sme.sector}</p>
            {sme.sub_sector && (
              <p className="text-xs text-muted-foreground mt-0.5">
                {sme.sub_sector}
              </p>
            )}
          </div>
        </div>

        {/* Location */}
        <div className="flex items-center gap-2 text-sm">
          <MapPin className="h-4 w-4 text-muted-foreground flex-shrink-0" />
          <span className="text-muted-foreground">
            {sme.district}, {sme.region}
          </span>
        </div>

        {/* Description */}
        {sme.business_description && (
          <p className="text-sm text-muted-foreground line-clamp-3 leading-relaxed">
            {sme.business_description}
          </p>
        )}

        {/* Operational Date */}
        {sme.operational_start_date && (
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <Calendar className="h-3.5 w-3.5 flex-shrink-0" />
            <span>Operational since {sme.operational_start_date}</span>
          </div>
        )}

        {/* Contact Information */}
        <div className="pt-3 border-t space-y-2">
          {canViewContactDetails ? (
            <>
              {sme.contact_phone && (
                <a
                  href={`tel:${sme.contact_phone}`}
                  className="flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  <Phone className="h-3.5 w-3.5 flex-shrink-0" />
                  <span className="truncate">{sme.contact_phone}</span>
                </a>
              )}

              {sme.contact_email && (
                <a
                  href={`mailto:${sme.contact_email}`}
                  className="flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  <Mail className="h-3.5 w-3.5 flex-shrink-0" />
                  <span className="truncate">{sme.contact_email}</span>
                </a>
              )}

              {sme.website && (
                <a
                  href={sme.website.startsWith("http") ? sme.website : `https://${sme.website}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  <Globe className="h-3.5 w-3.5 flex-shrink-0" />
                  <span className="truncate">{sme.website}</span>
                </a>
              )}
            </>
          ) : (
            <div className="relative">
              {/* Blurred contact details */}
              <div className="space-y-2 blur-sm select-none pointer-events-none">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Phone className="h-3.5 w-3.5 flex-shrink-0" />
                  <span className="truncate">+265 999 *** ***</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Mail className="h-3.5 w-3.5 flex-shrink-0" />
                  <span className="truncate">contact@***.com</span>
                </div>
              </div>

              {/* Overlay with info */}
              <div className="absolute inset-0 flex items-center justify-center">
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div className="flex items-center gap-2 bg-background/90 px-3 py-1.5 rounded-md border shadow-sm cursor-help">
                      <Lock className="h-4 w-4 text-muted-foreground" />
                      <span className="text-xs font-medium text-muted-foreground">
                        Contact Hidden
                      </span>
                      <Info className="h-3.5 w-3.5 text-muted-foreground" />
                    </div>
                  </TooltipTrigger>
                  <TooltipContent side="top" className="max-w-xs">
                    <div className="space-y-2">
                      <p className="text-sm font-medium">
                        Formalisation score of {minScoreRequired}+ required
                      </p>
                      <p className="text-xs text-muted-foreground">
                        To view contact details, complete your business formalisation profile:
                      </p>
                      <ul className="text-xs text-muted-foreground list-disc list-inside space-y-1">
                        <li>Register with tax authority</li>
                        <li>Open a business bank account</li>
                        <li>Join a business association</li>
                        <li>Complete your business profile</li>
                      </ul>
                    </div>
                  </TooltipContent>
                </Tooltip>
              </div>
            </div>
          )}
        </div>
      </CardContent>
      </Card>
    </TooltipProvider>
  );
}
