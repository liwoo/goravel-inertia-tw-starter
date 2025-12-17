import React, { useEffect, useState } from 'react';
import {
  Calendar, FileText, User, FolderOpen, Phone, Mail, MapPin,
  Building2, Globe, Users, Briefcase, CheckCircle2, XCircle, DollarSign,
  TrendingUp, Shield, CreditCard, Award, ChevronLeft, ChevronRight
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { CrudDetailViewProps } from '@/types/crud';
import { Sme, SmeClassification } from '@/types/sme';
import { BusinessFormalisation } from '@/types/business_formalisation';
import axios from 'axios';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { EmptyState } from '@/components/EmptyState';
import { ScoreBreakdownBar } from '@/components/ui/score-breakdown-bar';
import { CopyableText } from '@/components/ui/copyable-text';

export function SmeDetailView({
  item: sme,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Sme>) {
  const [primaryOwner, setPrimaryOwner] = useState<any>(null);
  const [additionalMembers, setAdditionalMembers] = useState<any[]>([]);
  const [employeeSummary, setEmployeeSummary] = useState<any>(null);
  const [formalisation, setFormalisation] = useState<BusinessFormalisation | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentMemberIndex, setCurrentMemberIndex] = useState(0);

  useEffect(() => {
    // Reset state before fetching to ensure fresh data display
    setLoading(true);
    setPrimaryOwner(null);
    setAdditionalMembers([]);
    setEmployeeSummary(null);
    setFormalisation(null);

    // Fetch related data
    const fetchRelatedData = async () => {
      try {
        // Fetch all related data in parallel
        const [ownerRes, membersRes, summaryRes, formalisationRes] = await Promise.all([
          axios.get(`/api/smes/${sme.id}/primary_business_owner`).catch(() => ({ data: { data: null } })),
          axios.get(`/api/smes/${sme.id}/additional_business_members`).catch(() => ({ data: { data: [] } })),
          axios.get(`/api/smes/${sme.id}/business_employee_summary`).catch(() => ({ data: { data: null } })),
          axios.get(`/api/smes/${sme.id}/business_formalisation`).catch(() => ({ data: { data: null } })),
        ]);

        setPrimaryOwner(ownerRes.data.data);
        setAdditionalMembers(membersRes.data.data || []);
        setEmployeeSummary(summaryRes.data.data);
        setFormalisation(formalisationRes.data.data);
      } catch (error) {
        console.error('Error fetching related data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRelatedData();
  }, [sme.id]);

  const formatDate = (date: string | Date | null | undefined) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const DetailRow = ({ label, value }: { label: string; value: any }) => (
    <div className="space-y-1">
      <p className="text-sm text-muted-foreground">{label}</p>
      <p className="font-medium text-foreground">{value || '-'}</p>
    </div>
  );

  const BooleanBadge = ({ value, trueLabel = 'Yes', falseLabel = 'No' }: { value: boolean; trueLabel?: string; falseLabel?: string }) => (
    <Badge variant={value ? 'default' : 'secondary'} className="gap-1">
      {value ? <CheckCircle2 className="h-3 w-3" /> : <XCircle className="h-3 w-3" />}
      {value ? trueLabel : falseLabel}
    </Badge>
  );

  const getClassificationBadge = (classification: SmeClassification | undefined) => {
    switch (classification) {
      case 'Micro':
        return { variant: 'outline' as const, className: 'border-blue-500 text-blue-600 bg-blue-50' };
      case 'Small':
        return { variant: 'outline' as const, className: 'border-emerald-500 text-emerald-600 bg-emerald-50' };
      case 'Medium':
        return { variant: 'outline' as const, className: 'border-purple-500 text-purple-600 bg-purple-50' };
      default:
        return { variant: 'secondary' as const, className: 'text-muted-foreground' };
    }
  };

  const classificationBadgeStyle = getClassificationBadge(sme.classification as SmeClassification);

  return (
    <div className="space-y-6">
      <Tabs defaultValue="business" className="w-full">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="business" title="Business">
            <Building2 className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Business</span>
          </TabsTrigger>
          <TabsTrigger value="owner" title="Primary Owner">
            <User className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Owner</span>
          </TabsTrigger>
          <TabsTrigger value="members" title="Members">
            <Users className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Members</span>
            {additionalMembers.length > 0 && (
              <Badge variant="secondary" className="ml-1 sm:ml-2 h-5 px-1 sm:px-1.5 text-xs">
                {additionalMembers.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="employees" title="Employees">
            <Briefcase className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Employees</span>
          </TabsTrigger>
          <TabsTrigger value="formalisation" title="Formalisation">
            <Shield className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Formal</span>
          </TabsTrigger>
        </TabsList>

        {/* Business Details Tab */}
        <TabsContent value="business" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Building2 className="h-5 w-5" />
                Business Information
              </CardTitle>
              <CardDescription>Core business details and registration information</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-1">
                  <p className="text-sm text-muted-foreground">USME Number</p>
                  <CopyableText value={sme.usmeNumber || sme.usme_number} className="font-mono font-medium text-foreground" iconSize="md" />
                </div>
                <div className="space-y-1">
                  <p className="text-sm text-muted-foreground">Classification</p>
                  <Badge variant={classificationBadgeStyle.variant} className={`text-sm ${classificationBadgeStyle.className}`}>
                    {sme.classification || 'Unclassified'}
                  </Badge>
                  <p className="text-xs text-muted-foreground mt-1">Based on employees, turnover & assets</p>
                </div>
                {/* Formalisation Score - inline after USME */}
                {!loading && formalisation && (
                  <div className="md:col-span-2">
                    <ScoreBreakdownBar
                      complianceScore={(formalisation as any).compliance_score ?? formalisation.complianceScore ?? 0}
                      teamStructureScore={(formalisation as any).team_structure_score ?? formalisation.teamStructureScore ?? 0}
                      financialScore={(formalisation as any).financial_score ?? formalisation.financialScore ?? 0}
                      totalScore={(formalisation as any).formalisation_score ?? formalisation.formalisationScore ?? 0}
                    />
                  </div>
                )}
                <DetailRow label="Business Name" value={sme.name} />
                <div className="space-y-1">
                  <p className="text-sm text-muted-foreground">Registration Number</p>
                  <CopyableText value={sme.registrationNumber || sme.registration_number} className="font-medium text-foreground" iconSize="md" />
                </div>
                <div className="space-y-1">
                  <p className="text-sm text-muted-foreground">Tax Identification Number</p>
                  <CopyableText value={sme.taxIdentificationNumber || sme.tax_identification_number} className="font-mono font-medium text-foreground" iconSize="md" />
                </div>
                <DetailRow label="Operational Since" value={formatDate(sme.operationalStartDate || sme.operational_start_date)} />
                <DetailRow label="Business Category" value={sme.businessCategory || sme.business_category} />
                <DetailRow label="Sector" value={sme.sector} />
                <DetailRow label="Sub Sector" value={sme.subSector || sme.sub_sector} />
              </div>

              <Separator />

              <div>
                <p className="text-sm text-muted-foreground mb-2">Business Description</p>
                <p className="text-sm text-foreground">{sme.businessDescription || sme.business_description || 'No description provided'}</p>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Contact Information</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <p className="text-sm text-muted-foreground">Phone</p>
                    <CopyableText value={sme.contactPhone || sme.contact_phone} className="font-mono font-medium text-foreground" iconSize="md" />
                  </div>
                  <div className="space-y-1">
                    <p className="text-sm text-muted-foreground">Email</p>
                    <CopyableText value={sme.contactEmail || sme.contact_email} className="font-medium text-foreground" iconSize="md" />
                  </div>
                  <DetailRow label="Website" value={sme.website} />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Location</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <DetailRow label="Region" value={(sme as any).region} />
                  <DetailRow label="District" value={sme.district} />
                  <DetailRow label="Traditional Authority" value={sme.traditionalAuthority || sme.traditional_authority} />
                  <DetailRow label="Physical Address" value={sme.physicalAddress || sme.physical_address} />
                  <DetailRow label="Postal Address" value={sme.postalAddress || sme.postal_address} />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Business Improvement & Financing</h4>
                <div>
                  <p className="text-sm text-muted-foreground mb-2">Improvement Aspects</p>
                  <div className="flex flex-wrap gap-2">
                    {((sme.businessImprovementAspects || sme.business_improvement_aspects) && (sme.businessImprovementAspects || sme.business_improvement_aspects || []).length > 0) ? (
                      (sme.businessImprovementAspects || sme.business_improvement_aspects || []).map((aspect, idx) => (
                        <Badge key={idx} variant="outline">{aspect}</Badge>
                      ))
                    ) : (
                      <span className="text-sm text-muted-foreground">No aspects specified</span>
                    )}
                  </div>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground mb-2">Accessed Financing</p>
                  <div className="flex flex-wrap gap-2">
                    {((sme.businessAccessedFinancing || sme.business_accessed_financing) && (sme.businessAccessedFinancing || sme.business_accessed_financing || []).length > 0) ? (
                      (sme.businessAccessedFinancing || sme.business_accessed_financing || []).map((financing, idx) => (
                        <Badge key={idx} variant="secondary">{financing}</Badge>
                      ))
                    ) : (
                      <span className="text-sm text-muted-foreground">No financing accessed</span>
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Primary Owner Tab */}
        <TabsContent value="owner" className="space-y-6">
          {loading ? (
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">Loading...</p>
              </CardContent>
            </Card>
          ) : primaryOwner ? (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  Primary Business Owner
                </CardTitle>
                <CardDescription>Details of the primary business owner</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <DetailRow label="First Name" value={primaryOwner.firstName || primaryOwner.first_name} />
                  <DetailRow label="Last Name" value={primaryOwner.lastName || primaryOwner.last_name} />
                  <DetailRow label="Other Names" value={primaryOwner.otherNames || primaryOwner.other_names} />
                  <div className="space-y-1">
                    <p className="text-sm text-muted-foreground">National ID</p>
                    <CopyableText value={primaryOwner.nationalIdNumber || primaryOwner.national_id_number} className="font-mono font-medium text-foreground" iconSize="md" />
                  </div>
                  <DetailRow label="Nationality" value={primaryOwner.nationality} />
                  <DetailRow label="Date of Birth" value={formatDate(primaryOwner.dateOfBirth || primaryOwner.date_of_birth)} />
                  <DetailRow label="Gender" value={primaryOwner.gender} />
                  <DetailRow label="Education Level" value={primaryOwner.educationLevel || primaryOwner.education_level} />
                  <DetailRow label="Malawian Status" value={primaryOwner.malawianStatus || primaryOwner.malawian_status} />
                  <div className="space-y-1">
                    <p className="text-sm text-muted-foreground">Special Needs</p>
                    <BooleanBadge value={primaryOwner.hasSpecialNeeds ?? primaryOwner.has_special_needs ?? false} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Contact Information</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-1">
                      <p className="text-sm text-muted-foreground">Phone Number</p>
                      <CopyableText value={primaryOwner.phoneNumber || primaryOwner.phone_number} className="font-mono font-medium text-foreground" iconSize="md" />
                    </div>
                    <div className="space-y-1">
                      <p className="text-sm text-muted-foreground">Landline</p>
                      <CopyableText value={primaryOwner.landlineNumber || primaryOwner.landline_number} className="font-mono font-medium text-foreground" iconSize="md" />
                    </div>
                    <div className="space-y-1">
                      <p className="text-sm text-muted-foreground">Email</p>
                      <CopyableText value={primaryOwner.email} className="font-medium text-foreground" iconSize="md" />
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Location</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="Region" value={primaryOwner.region} />
                    <DetailRow label="District" value={primaryOwner.district} />
                    <DetailRow label="Traditional Authority" value={primaryOwner.traditionalAuthority || primaryOwner.traditional_authority} />
                    <DetailRow label="Physical Address" value={primaryOwner.physicalAddress || primaryOwner.physical_address} />
                    <DetailRow label="Postal Address" value={primaryOwner.postalAddress || primaryOwner.postal_address} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Alternative Contact</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="Contact Name" value={primaryOwner.altContactName || primaryOwner.alt_contact_name} />
                    <DetailRow label="Relationship" value={primaryOwner.altContactRelationship || primaryOwner.alt_contact_relationship} />
                    <DetailRow label="Contact Phone" value={primaryOwner.altContactPhone || primaryOwner.alt_contact_phone} />
                  </div>
                </div>
              </CardContent>
            </Card>
          ) : (
            <EmptyState
              icon={User}
              title="No Primary Owner"
              description="Primary business owner information has not been added yet."
            />
          )}
        </TabsContent>

        {/* Additional Members Tab */}
        <TabsContent value="members" className="space-y-6">
          {loading ? (
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">Loading...</p>
              </CardContent>
            </Card>
          ) : additionalMembers.length > 0 ? (
            <div className="space-y-4">
              {/* Carousel Navigation */}
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => setCurrentMemberIndex(Math.max(0, currentMemberIndex - 1))}
                    disabled={currentMemberIndex === 0}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => setCurrentMemberIndex(Math.min(additionalMembers.length - 1, currentMemberIndex + 1))}
                    disabled={currentMemberIndex === additionalMembers.length - 1}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">
                    {currentMemberIndex + 1} of {additionalMembers.length}
                  </span>
                  {/* Dot indicators */}
                  <div className="flex gap-1">
                    {additionalMembers.map((_, idx) => (
                      <button
                        key={idx}
                        onClick={() => setCurrentMemberIndex(idx)}
                        className={`h-2 w-2 rounded-full transition-colors ${
                          idx === currentMemberIndex
                            ? 'bg-primary'
                            : 'bg-muted-foreground/30 hover:bg-muted-foreground/50'
                        }`}
                      />
                    ))}
                  </div>
                </div>
              </div>

              {/* Current Member Card */}
              {(() => {
                const member = additionalMembers[currentMemberIndex];
                return (
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <Users className="h-5 w-5" />
                        {member.firstName || member.first_name} {member.lastName || member.last_name}
                      </CardTitle>
                      <CardDescription>Additional business member</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <DetailRow label="First Name" value={member.firstName || member.first_name} />
                        <DetailRow label="Last Name" value={member.lastName || member.last_name} />
                        <DetailRow label="Other Names" value={member.otherNames || member.other_names} />
                        <DetailRow label="Gender" value={member.gender} />
                        <DetailRow label="National ID" value={member.nationalIdNumber || member.national_id_number} />
                        <DetailRow label="Nationality" value={member.nationality} />
                        <DetailRow label="Date of Birth" value={formatDate(member.dateOfBirth || member.date_of_birth)} />
                        <DetailRow label="Phone Number" value={member.phoneNumber || member.phone_number} />
                        <DetailRow label="Email" value={member.email} />
                      </div>

                      <Separator />

                      <div className="flex gap-4">
                        <div className="flex items-center gap-2">
                          <span className="text-sm text-muted-foreground">Intern:</span>
                          <BooleanBadge value={member.isIntern ?? member.is_intern ?? false} />
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-sm text-muted-foreground">Part Time:</span>
                          <BooleanBadge value={member.isPartTime ?? member.is_part_time ?? false} />
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })()}
            </div>
          ) : (
            <EmptyState
              icon={Users}
              title="No Team Members"
              description="No additional business members have been added yet."
            />
          )}
        </TabsContent>

        {/* Employee Summary Tab */}
        <TabsContent value="employees" className="space-y-6">
          {loading ? (
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">Loading...</p>
              </CardContent>
            </Card>
          ) : employeeSummary ? (
            <>
              {/* By Employee Type */}
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">By Employee Type</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 gap-6">
                    {/* Full-Time Employees */}
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-full bg-blue-500" />
                        <span className="text-sm text-muted-foreground">Full-Time Employees</span>
                      </div>
                      <div className="flex items-baseline gap-2">
                        <span className="text-4xl font-bold tabular-nums text-foreground">
                          {(employeeSummary.fullTimeMales ?? employeeSummary.full_time_males ?? 0) +
                           (employeeSummary.fullTimeFemales ?? employeeSummary.full_time_females ?? 0)}
                        </span>
                        <span className="text-sm text-muted-foreground">Employees</span>
                      </div>
                    </div>

                    {/* Part-Time Employees */}
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-full bg-emerald-500" />
                        <span className="text-sm text-muted-foreground">Part-Time Employees</span>
                      </div>
                      <div className="flex items-baseline gap-2">
                        <span className="text-4xl font-bold tabular-nums text-foreground">
                          {(employeeSummary.partTimeMales ?? employeeSummary.part_time_males ?? 0) +
                           (employeeSummary.partTimeFemales ?? employeeSummary.part_time_females ?? 0)}
                        </span>
                        <span className="text-sm text-muted-foreground">Employees</span>
                      </div>
                    </div>

                    {/* Internship Employees */}
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-full bg-amber-500" />
                        <span className="text-sm text-muted-foreground">Internship Employees</span>
                      </div>
                      <div className="flex items-baseline gap-2">
                        <span className="text-4xl font-bold tabular-nums text-foreground">
                          {(employeeSummary.internMales ?? employeeSummary.intern_males ?? 0) +
                           (employeeSummary.internFemales ?? employeeSummary.intern_females ?? 0)}
                        </span>
                        <span className="text-sm text-muted-foreground">Employees</span>
                      </div>
                    </div>

                    {/* Total */}
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-full bg-slate-500" />
                        <span className="text-sm text-muted-foreground">Total Employees</span>
                      </div>
                      <div className="flex items-baseline gap-2">
                        <span className="text-4xl font-bold tabular-nums text-foreground">
                          {(employeeSummary.fullTimeMales ?? employeeSummary.full_time_males ?? 0) +
                           (employeeSummary.fullTimeFemales ?? employeeSummary.full_time_females ?? 0) +
                           (employeeSummary.partTimeMales ?? employeeSummary.part_time_males ?? 0) +
                           (employeeSummary.partTimeFemales ?? employeeSummary.part_time_females ?? 0) +
                           (employeeSummary.internMales ?? employeeSummary.intern_males ?? 0) +
                           (employeeSummary.internFemales ?? employeeSummary.intern_females ?? 0)}
                        </span>
                        <span className="text-sm text-muted-foreground">Employees</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* By Gender */}
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">By Gender</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 gap-6">
                    {/* Male Employees */}
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-full bg-sky-500" />
                        <span className="text-sm text-muted-foreground">Male Employees</span>
                      </div>
                      <div className="flex items-baseline gap-2">
                        <span className="text-4xl font-bold tabular-nums text-foreground">
                          {(employeeSummary.fullTimeMales ?? employeeSummary.full_time_males ?? 0) +
                           (employeeSummary.partTimeMales ?? employeeSummary.part_time_males ?? 0) +
                           (employeeSummary.internMales ?? employeeSummary.intern_males ?? 0)}
                        </span>
                        <span className="text-sm text-muted-foreground">Employees</span>
                      </div>
                    </div>

                    {/* Female Employees */}
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-full bg-pink-500" />
                        <span className="text-sm text-muted-foreground">Female Employees</span>
                      </div>
                      <div className="flex items-baseline gap-2">
                        <span className="text-4xl font-bold tabular-nums text-foreground">
                          {(employeeSummary.fullTimeFemales ?? employeeSummary.full_time_females ?? 0) +
                           (employeeSummary.partTimeFemales ?? employeeSummary.part_time_females ?? 0) +
                           (employeeSummary.internFemales ?? employeeSummary.intern_females ?? 0)}
                        </span>
                        <span className="text-sm text-muted-foreground">Employees</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </>
          ) : (
            <EmptyState
              icon={Briefcase}
              title="No Employee Summary"
              description="Employee summary data has not been recorded yet."
            />
          )}
        </TabsContent>

        {/* Formalisation Tab */}
        <TabsContent value="formalisation" className="space-y-6">
          {loading ? (
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">Loading...</p>
              </CardContent>
            </Card>
          ) : formalisation ? (
            <div className="space-y-6">
              {/* Registration & Compliance Status */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Shield className="h-5 w-5" />
                    Registration & Compliance Status
                  </CardTitle>
                  <CardDescription>Legal registration and tax compliance information</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <CreditCard className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Bank Account</p>
                        <BooleanBadge
                          value={(formalisation as any).hasBankAccount ?? (formalisation as any).has_bank_account ?? false}
                          trueLabel="Has Account"
                          falseLabel="No Account"
                        />
                      </div>
                    </div>

                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <FileText className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Tax Clarification Certificate</p>
                        <BooleanBadge
                          value={(formalisation as any).hasTaxClarification ?? (formalisation as any).has_tax_clarification ?? false}
                          trueLabel="Has Certificate"
                          falseLabel="No Certificate"
                        />
                      </div>
                    </div>

                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <Shield className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">VAT Registration</p>
                        <BooleanBadge
                          value={(formalisation as any).isRegisteredForVat ?? (formalisation as any).is_registered_for_vat ?? false}
                          trueLabel="Registered"
                          falseLabel="Not Registered"
                        />
                      </div>
                    </div>

                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <Globe className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Export License</p>
                        <BooleanBadge
                          value={(formalisation as any).hasExportLicense ?? (formalisation as any).has_export_license ?? false}
                          trueLabel="Has License"
                          falseLabel="No License"
                        />
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Business Associations & Support */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Users className="h-5 w-5" />
                    Business Associations & Support
                  </CardTitle>
                  <CardDescription>Membership and business development services</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <Users className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Business Association</p>
                        <BooleanBadge
                          value={(formalisation as any).isMemberOfAssociation ?? (formalisation as any).is_member_of_association ?? false}
                          trueLabel="Member"
                          falseLabel="Non-Member"
                        />
                      </div>
                    </div>

                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <Building2 className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Business Affiliation</p>
                        <BooleanBadge
                          value={(formalisation as any).isAffiliated ?? (formalisation as any).is_affiliated ?? false}
                          trueLabel="Affiliated"
                          falseLabel="Not Affiliated"
                        />
                      </div>
                    </div>

                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <Briefcase className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Business Development Services</p>
                        <BooleanBadge
                          value={(formalisation as any).hasAccessedBds ?? (formalisation as any).has_accessed_bds ?? false}
                          trueLabel="Accessed"
                          falseLabel="Not Accessed"
                        />
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Financial Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <TrendingUp className="h-5 w-5" />
                    Financial Information
                  </CardTitle>
                  <CardDescription>Business financial metrics and valuation</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                      <div className="flex items-center gap-2">
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                        <p className="text-sm text-muted-foreground">Annual Turnover</p>
                      </div>
                      <p className="text-2xl font-bold text-foreground">
                        {((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) ? `MWK ${((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover).toLocaleString()}` : 'Not Specified'}
                      </p>
                      {((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) > 0 && (
                        <p className="text-xs text-muted-foreground">
                          {((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) >= 1000000 ? 'Above MWK 1M' : 'Below MWK 1M'}
                        </p>
                      )}
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2">
                        <Building2 className="h-4 w-4 text-muted-foreground" />
                        <p className="text-sm text-muted-foreground">Estimated Value of Assets</p>
                      </div>
                      <p className="text-2xl font-bold text-foreground">
                        {((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets) ? `MWK ${((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets).toLocaleString()}` : 'Not Specified'}
                      </p>
                      {((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets) > 0 && (
                        <p className="text-xs text-muted-foreground">
                          {((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets) >= 5000000 ? 'High Value Assets' : ((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets) >= 1000000 ? 'Medium Value Assets' : 'Low Value Assets'}
                        </p>
                      )}
                    </div>
                  </div>

                  {/* Financial Metrics Summary */}
                  {(((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) || ((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets)) && (
                    <>
                      <Separator />
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                        <div>
                          <p className="text-xs text-muted-foreground">Asset to Turnover Ratio</p>
                          <p className="text-lg font-semibold">
                            {((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) && ((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets)
                              ? `${(((formalisation as any).estimatedValueOfAssets ?? (formalisation as any).estimated_value_of_assets) / ((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover)).toFixed(2)}x`
                              : '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Monthly Turnover</p>
                          <p className="text-lg font-semibold">
                            {((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover)
                              ? `MWK ${Math.round(((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) / 12).toLocaleString()}`
                              : '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Daily Turnover</p>
                          <p className="text-lg font-semibold">
                            {((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover)
                              ? `MWK ${Math.round(((formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover) / 365).toLocaleString()}`
                              : '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Business Scale</p>
                          <p className="text-lg font-semibold">
                            {(() => {
                              const turnover = (formalisation as any).annualTurnover ?? (formalisation as any).annual_turnover ?? 0;
                              if (turnover >= 50000000) return 'Large';
                              if (turnover >= 10000000) return 'Medium';
                              if (turnover >= 1000000) return 'Small';
                              return 'Micro';
                            })()}
                          </p>
                        </div>
                      </div>
                    </>
                  )}
                </CardContent>
              </Card>

              {/* Compliance Summary */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-sm font-medium text-muted-foreground">Compliance Summary</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2">
                    {((formalisation as any).hasBankAccount ?? (formalisation as any).has_bank_account) && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Bank Account
                      </Badge>
                    )}
                    {((formalisation as any).hasTaxClarification ?? (formalisation as any).has_tax_clarification) && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Tax Compliant
                      </Badge>
                    )}
                    {((formalisation as any).isRegisteredForVat ?? (formalisation as any).is_registered_for_vat) && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        VAT Registered
                      </Badge>
                    )}
                    {((formalisation as any).hasExportLicense ?? (formalisation as any).has_export_license) && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Export Ready
                      </Badge>
                    )}
                    {((formalisation as any).isMemberOfAssociation ?? (formalisation as any).is_member_of_association) && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Association Member
                      </Badge>
                    )}
                    {((formalisation as any).hasAccessedBds ?? (formalisation as any).has_accessed_bds) && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        BDS Beneficiary
                      </Badge>
                    )}
                    {!((formalisation as any).hasBankAccount ?? (formalisation as any).has_bank_account) &&
                     !((formalisation as any).hasTaxClarification ?? (formalisation as any).has_tax_clarification) &&
                     !((formalisation as any).isRegisteredForVat ?? (formalisation as any).is_registered_for_vat) &&
                     !((formalisation as any).hasExportLicense ?? (formalisation as any).has_export_license) &&
                     !((formalisation as any).isMemberOfAssociation ?? (formalisation as any).is_member_of_association) &&
                     !((formalisation as any).hasAccessedBds ?? (formalisation as any).has_accessed_bds) && (
                      <span className="text-sm text-muted-foreground">No compliance items completed</span>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : (
            <EmptyState
              icon={Shield}
              title="No Formalization Data"
              description="Business formalization information has not been recorded yet."
            />
          )}
        </TabsContent>
      </Tabs>

      {/* Metadata Section */}
      <Separator />
      <div>
        <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Created</p>
            <p className="font-medium text-sm text-foreground">{formatDate(sme.createdAt || sme.created_at || null)}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Last Updated</p>
            <p className="font-medium text-sm text-foreground">{formatDate(sme.updatedAt || sme.updated_at || null)}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
