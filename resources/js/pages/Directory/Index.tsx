import { useState, useEffect, useCallback, useRef } from "react";
import { Head } from "@inertiajs/react";
import { DirectoryPageProps, DirectorySME } from "@/types/directory";
import { DirectoryFilterSidebar } from "./components/DirectoryFilterSidebar";
import { DirectorySMECard } from "./components/DirectorySMECard";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { useDebounce } from "@/hooks/useDebounce";
import { useIsMobile } from "@/hooks/use-mobile";
import { Search, Filter, X, ChevronLeft, ChevronRight } from "lucide-react";
import axios from "@/lib/axios";
import Admin from "@/layouts/Admin";

export default function DirectoryIndex({
  smes: initialSmes = [],
  filters: initialFilters,
  pagination: initialPagination,
  userFormalisationScore = 0,
  minScoreForContactView = 70,
}: DirectoryPageProps) {
  const isMobile = useIsMobile();

  // Default values for filters
  const filters = initialFilters || { sectors: [], districts: [] };

  // State
  const [smes, setSmes] = useState<DirectorySME[]>(initialSmes || []);
  const [pagination, setPagination] = useState(initialPagination || { total: 0, page: 1, pageSize: 12, totalPages: 0 });
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedSectors, setSelectedSectors] = useState<string[]>([]);
  const [selectedDistricts, setSelectedDistricts] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);

  const debouncedSearch = useDebounce(searchQuery, 300);
  const isFirstRender = useRef(true);

  // Fetch results from API
  const fetchResults = useCallback(
    async (page: number = 1) => {
      setIsLoading(true);

      try {
        const params: Record<string, any> = {
          page,
          pageSize: 12,
        };

        if (debouncedSearch) {
          params.search = debouncedSearch;
        }

        if (selectedSectors.length > 0) {
          params["sectors[]"] = selectedSectors;
        }

        if (selectedDistricts.length > 0) {
          params["districts[]"] = selectedDistricts;
        }

        const response = await axios.get("/api/directory/search", { params });

        setSmes(response.data.smes);
        setPagination(response.data.pagination);
      } catch (error) {
        console.error("Failed to fetch directory results:", error);
      } finally {
        setIsLoading(false);
      }
    },
    [debouncedSearch, selectedSectors, selectedDistricts]
  );

  // Trigger search when filters change (but not on initial render)
  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }

    // Only search if there's a search term or filters selected
    if (debouncedSearch || selectedSectors.length > 0 || selectedDistricts.length > 0) {
      setHasSearched(true);
      fetchResults(1);
    } else {
      // Clear results if all filters are removed
      setHasSearched(false);
      setSmes([]);
      setPagination({ total: 0, page: 1, pageSize: 12, totalPages: 0 });
    }
  }, [debouncedSearch, selectedSectors, selectedDistricts]);

  // Filter handlers
  const handleSectorChange = (sector: string, checked: boolean) => {
    setSelectedSectors((prev) =>
      checked ? [...prev, sector] : prev.filter((s) => s !== sector)
    );
  };

  const handleDistrictChange = (district: string, checked: boolean) => {
    setSelectedDistricts((prev) =>
      checked ? [...prev, district] : prev.filter((d) => d !== district)
    );
  };

  const handleClearAllFilters = () => {
    setSelectedSectors([]);
    setSelectedDistricts([]);
    setSearchQuery("");
  };

  const handleRemoveFilter = (type: "sector" | "district", value: string) => {
    if (type === "sector") {
      setSelectedSectors((prev) => prev.filter((s) => s !== value));
    } else {
      setSelectedDistricts((prev) => prev.filter((d) => d !== value));
    }
  };

  const handlePageChange = (newPage: number) => {
    fetchResults(newPage);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const activeFilterCount = selectedSectors.length + selectedDistricts.length;

  return (
    <Admin title="MSME Directory">
      <Head title="MSME Directory" />

      <div className="p-4 md:p-6 flex flex-col gap-6">
        {/* Header */}
        <div className="flex flex-col gap-1">
          <h1 className="text-2xl font-semibold tracking-tight">
            MSME Directory
          </h1>
          <p className="text-muted-foreground">
            <span className="hidden sm:inline">Search and explore registered Micro, Small and Medium Enterprises</span>
            <span className="sm:hidden">Explore registered MSMEs</span>
          </p>
        </div>

        {/* Search Bar */}
        <div className="flex items-center gap-4">
          <div className="relative flex-1 max-w-xl">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search MSMEs by name, sector, location..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 h-10"
            />
          </div>
          {isMobile && (
            <Sheet open={mobileFiltersOpen} onOpenChange={setMobileFiltersOpen}>
              <SheetTrigger asChild>
                <Button variant="outline" size="sm" className="gap-2">
                  <Filter className="h-4 w-4" />
                  Filters
                  {activeFilterCount > 0 && (
                    <Badge variant="secondary" className="ml-1 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs">
                      {activeFilterCount}
                    </Badge>
                  )}
                </Button>
              </SheetTrigger>
              <SheetContent side="left" className="w-[85vw] sm:w-[350px] p-0">
                <SheetHeader className="px-4 py-3 border-b">
                  <SheetTitle>Filters</SheetTitle>
                </SheetHeader>
                <DirectoryFilterSidebar
                  filters={filters}
                  selectedSectors={selectedSectors}
                  selectedDistricts={selectedDistricts}
                  onSectorChange={handleSectorChange}
                  onDistrictChange={handleDistrictChange}
                  onClearAll={handleClearAllFilters}
                  className="h-[calc(100vh-4rem)]"
                />
              </SheetContent>
            </Sheet>
          )}
        </div>

        {/* Active Filter Badges */}
        {activeFilterCount > 0 && (
          <div className="flex flex-wrap gap-2">
            {selectedSectors.map((sector) => (
              <Badge key={sector} variant="secondary" className="gap-1 pl-3 pr-2">
                {filters.sectors.find((s) => s.value === sector)?.label || sector}
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-4 w-4 p-0 hover:bg-transparent"
                  onClick={() => handleRemoveFilter("sector", sector)}
                >
                  <X className="h-3 w-3" />
                </Button>
              </Badge>
            ))}
            {selectedDistricts.map((district) => {
              const districtLabel = filters.districts
                .flatMap((g) => g.districts)
                .find((d) => d.value === district)?.label;
              return (
                <Badge key={district} variant="secondary" className="gap-1 pl-3 pr-2">
                  {districtLabel || district}
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-4 w-4 p-0 hover:bg-transparent"
                    onClick={() => handleRemoveFilter("district", district)}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </Badge>
              );
            })}
          </div>
        )}

        {/* Main Content */}
        <div className="flex gap-6">
          {/* Desktop Sidebar */}
          {!isMobile && (
            <aside className="w-64 flex-shrink-0">
              <div className="sticky top-20">
                <div className="border rounded-lg bg-card">
                  <DirectoryFilterSidebar
                    filters={filters}
                    selectedSectors={selectedSectors}
                    selectedDistricts={selectedDistricts}
                    onSectorChange={handleSectorChange}
                    onDistrictChange={handleDistrictChange}
                    onClearAll={handleClearAllFilters}
                  />
                </div>
              </div>
            </aside>
          )}

          {/* Results */}
          <div className="flex-1 min-w-0">
            {/* Results Header */}
            {hasSearched && (
              <div className="mb-6">
                <h2 className="text-lg font-semibold">
                  {isLoading ? "Searching..." : `${pagination.total} ${pagination.total === 1 ? "MSME" : "MSMEs"} found`}
                </h2>
              </div>
            )}

            {/* SME Grid */}
            {isLoading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {[...Array(6)].map((_, i) => (
                  <div key={i} className="h-64 rounded-lg border bg-card animate-pulse" />
                ))}
              </div>
            ) : !hasSearched ? (
              <div className="text-center py-16 border rounded-lg bg-muted/30">
                <Search className="mx-auto h-12 w-12 text-muted-foreground/50" />
                <h3 className="mt-4 text-lg font-medium">Search the Directory</h3>
                <p className="mt-2 text-muted-foreground max-w-sm mx-auto">
                  Enter a search term or select filters to find MSMEs
                </p>
              </div>
            ) : smes.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {smes.map((sme) => (
                  <DirectorySMECard
                    key={sme.usme_number}
                    sme={sme}
                    canViewContactDetails={userFormalisationScore >= minScoreForContactView}
                    minScoreRequired={minScoreForContactView}
                  />
                ))}
              </div>
            ) : (
              <div className="text-center py-16 border rounded-lg bg-muted/30">
                <p className="text-muted-foreground text-lg">
                  No MSMEs found matching your criteria.
                </p>
                <Button variant="outline" onClick={handleClearAllFilters} className="mt-4">
                  Clear all filters
                </Button>
              </div>
            )}

            {/* Pagination */}
            {pagination.totalPages > 1 && !isLoading && (
              <div className="flex items-center justify-center gap-2 mt-8">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handlePageChange(pagination.page - 1)}
                  disabled={pagination.page === 1}
                >
                  <ChevronLeft className="h-4 w-4 mr-1" />
                  Previous
                </Button>

                <div className="flex items-center gap-1">
                  {[...Array(pagination.totalPages)].map((_, i) => {
                    const page = i + 1;
                    const shouldShow =
                      page === 1 ||
                      page === pagination.totalPages ||
                      Math.abs(page - pagination.page) <= 1;

                    if (!shouldShow) {
                      if (page === pagination.page - 2 || page === pagination.page + 2) {
                        return <span key={page} className="px-2 text-muted-foreground">...</span>;
                      }
                      return null;
                    }

                    return (
                      <Button
                        key={page}
                        variant={page === pagination.page ? "default" : "outline"}
                        size="sm"
                        onClick={() => handlePageChange(page)}
                        className="w-9"
                      >
                        {page}
                      </Button>
                    );
                  })}
                </div>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handlePageChange(pagination.page + 1)}
                  disabled={pagination.page === pagination.totalPages}
                >
                  Next
                  <ChevronRight className="h-4 w-4 ml-1" />
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>
    </Admin>
  );
}
