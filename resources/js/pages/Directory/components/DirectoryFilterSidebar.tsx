import { DirectoryFilters } from "@/types/directory";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { ChevronDown, X } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";

interface DirectoryFilterSidebarProps {
  filters: DirectoryFilters;
  selectedSectors: string[];
  selectedDistricts: string[];
  onSectorChange: (sector: string, checked: boolean) => void;
  onDistrictChange: (district: string, checked: boolean) => void;
  onClearAll: () => void;
  className?: string;
}

export function DirectoryFilterSidebar({
  filters,
  selectedSectors,
  selectedDistricts,
  onSectorChange,
  onDistrictChange,
  onClearAll,
  className,
}: DirectoryFilterSidebarProps) {
  const [sectorsOpen, setSectorsOpen] = useState(true);
  // Only one region can be open at a time, start with none expanded
  const [openRegion, setOpenRegion] = useState<string | null>(null);

  // Guard against undefined filters
  const sectors = filters?.sectors || [];
  const districts = filters?.districts || [];

  const activeFilterCount = selectedSectors.length + selectedDistricts.length;

  const toggleDistrictGroup = (region: string) => {
    // If clicking the currently open region, close it; otherwise open the new one
    setOpenRegion((prev) => (prev === region ? null : region));
  };

  return (
    <div className={cn("flex flex-col h-full", className)}>
      {/* Header */}
      <div className="px-4 py-3 border-b">
        <div className="flex items-center justify-between">
          <h2 className="font-semibold text-sm">
            Filters
            {activeFilterCount > 0 && (
              <span className="ml-2 text-xs font-normal text-muted-foreground">
                ({activeFilterCount})
              </span>
            )}
          </h2>
          {activeFilterCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onClearAll}
              className="h-7 px-2 text-xs"
            >
              Clear all
            </Button>
          )}
        </div>
      </div>

      {/* Filters */}
      <ScrollArea className="flex-1">
        <div className="p-4 space-y-4">
          {/* Sectors Filter */}
          <Collapsible open={sectorsOpen} onOpenChange={setSectorsOpen}>
            <CollapsibleTrigger className="flex w-full items-center justify-between hover:bg-accent hover:text-accent-foreground rounded-md px-2 py-1.5 -mx-2">
              <span className="text-sm font-medium">Sector</span>
              <ChevronDown
                className={cn(
                  "h-4 w-4 transition-transform",
                  sectorsOpen && "transform rotate-180"
                )}
              />
            </CollapsibleTrigger>
            <CollapsibleContent className="mt-2 space-y-2">
              {sectors.map((sector) => (
                <div
                  key={sector.value}
                  className="flex items-start space-x-2 py-1"
                >
                  <Checkbox
                    id={`sector-${sector.value}`}
                    checked={selectedSectors.includes(sector.value)}
                    onCheckedChange={(checked) =>
                      onSectorChange(sector.value, checked as boolean)
                    }
                  />
                  <Label
                    htmlFor={`sector-${sector.value}`}
                    className="flex-1 text-sm font-normal leading-none cursor-pointer peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  >
                    <div className="flex items-baseline justify-between gap-2">
                      <span className="flex-1">{sector.label}</span>
                      <span className="text-xs text-muted-foreground">
                        {sector.count}
                      </span>
                    </div>
                  </Label>
                </div>
              ))}
            </CollapsibleContent>
          </Collapsible>

          <Separator />

          {/* Districts Filter */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium px-2">District</h3>
            {districts.map((group) => (
              <Collapsible
                key={group.region}
                open={openRegion === group.region}
                onOpenChange={() => toggleDistrictGroup(group.region)}
              >
                <CollapsibleTrigger className="flex w-full items-center justify-between hover:bg-accent hover:text-accent-foreground rounded-md px-2 py-1.5 -mx-2">
                  <span className="text-sm font-medium text-muted-foreground">
                    {group.region}
                  </span>
                  <ChevronDown
                    className={cn(
                      "h-4 w-4 transition-transform",
                      openRegion === group.region && "transform rotate-180"
                    )}
                  />
                </CollapsibleTrigger>
                <CollapsibleContent className="mt-2 space-y-2 pl-2">
                  {group.districts.map((district) => (
                    <div
                      key={district.value}
                      className="flex items-start space-x-2 py-1"
                    >
                      <Checkbox
                        id={`district-${district.value}`}
                        checked={selectedDistricts.includes(district.value)}
                        onCheckedChange={(checked) =>
                          onDistrictChange(district.value, checked as boolean)
                        }
                      />
                      <Label
                        htmlFor={`district-${district.value}`}
                        className="flex-1 text-sm font-normal leading-none cursor-pointer peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        <div className="flex items-baseline justify-between gap-2">
                          <span className="flex-1">{district.label}</span>
                          <span className="text-xs text-muted-foreground">
                            {district.count}
                          </span>
                        </div>
                      </Label>
                    </div>
                  ))}
                </CollapsibleContent>
              </Collapsible>
            ))}
          </div>
        </div>
      </ScrollArea>
    </div>
  );
}
