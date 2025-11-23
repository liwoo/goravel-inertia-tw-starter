import React, { useEffect, useState } from 'react';
import {
  Calendar, FileText, User, FolderOpen, Phone, Mail, MapPin,
  Building2, Globe, Users, Briefcase, CheckCircle2, XCircle, DollarSign,
  TrendingUp, Shield, CreditCard, Award
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { CrudDetailViewProps } from '@/types/crud';
import { Sme } from '@/types/sme';
import { BusinessFormalisation } from '@/types/business_formalisation';
import axios from 'axios';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { EmptyState } from '@/components/EmptyState';

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

  useEffect(() => {
    // Fetch related data
    const fetchRelatedData = async () => {
      try {
        setLoading(true);

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

  const formatDate = (date: string | Date | null) => {
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

  return (
    <div className="space-y-6">
      <Tabs defaultValue="business" className="w-full">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="business">Business</TabsTrigger>
          <TabsTrigger value="owner">Primary Owner</TabsTrigger>
          <TabsTrigger value="members">
            Members
            {additionalMembers.length > 0 && (
              <Badge variant="secondary" className="ml-2 h-5 px-1.5 text-xs">
                {additionalMembers.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="employees">Employees</TabsTrigger>
          <TabsTrigger value="formalisation">Formalisation</TabsTrigger>
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
                <DetailRow label="USME Number" value={sme.usmeNumber} />
                <DetailRow label="Business Name" value={sme.name} />
                <DetailRow label="Registration Number" value={sme.registrationNumber} />
                <DetailRow label="Tax Identification Number" value={sme.taxIdentificationNumber} />
                <DetailRow label="Operational Since" value={formatDate(sme.operationalStartDate)} />
                <DetailRow label="Business Category" value={sme.businessCategory} />
                <DetailRow label="Sector" value={sme.sector} />
                <DetailRow label="Sub Sector" value={sme.subSector} />
              </div>

              <Separator />

              <div>
                <p className="text-sm text-muted-foreground mb-2">Business Description</p>
                <p className="text-sm text-foreground">{sme.businessDescription || 'No description provided'}</p>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Contact Information</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <DetailRow label="Phone" value={sme.contactPhone} />
                  <DetailRow label="Email" value={sme.contactEmail} />
                  <DetailRow label="Website" value={sme.website} />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Location</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <DetailRow label="Region" value={sme.region} />
                  <DetailRow label="District" value={sme.district} />
                  <DetailRow label="Traditional Authority" value={sme.traditionalAuthority} />
                  <DetailRow label="Physical Address" value={sme.physicalAddress} />
                  <DetailRow label="Postal Address" value={sme.postalAddress} />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Business Improvement & Financing</h4>
                <div>
                  <p className="text-sm text-muted-foreground mb-2">Improvement Aspects</p>
                  <div className="flex flex-wrap gap-2">
                    {sme.businessImprovementAspects && sme.businessImprovementAspects.length > 0 ? (
                      sme.businessImprovementAspects.map((aspect, idx) => (
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
                    {sme.businessAccessedFinancing && sme.businessAccessedFinancing.length > 0 ? (
                      sme.businessAccessedFinancing.map((financing, idx) => (
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
                  <DetailRow label="First Name" value={primaryOwner.firstName} />
                  <DetailRow label="Last Name" value={primaryOwner.lastName} />
                  <DetailRow label="Other Names" value={primaryOwner.otherNames} />
                  <DetailRow label="National ID" value={primaryOwner.nationalIdNumber} />
                  <DetailRow label="Nationality" value={primaryOwner.nationality} />
                  <DetailRow label="Date of Birth" value={formatDate(primaryOwner.dateOfBirth)} />
                  <DetailRow label="Gender" value={primaryOwner.gender} />
                  <DetailRow label="Education Level" value={primaryOwner.educationLevel} />
                  <DetailRow label="Malawian Status" value={primaryOwner.malawianStatus} />
                  <div className="space-y-1">
                    <p className="text-sm text-muted-foreground">Special Needs</p>
                    <BooleanBadge value={primaryOwner.hasSpecialNeeds} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Contact Information</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="Phone Number" value={primaryOwner.phoneNumber} />
                    <DetailRow label="Landline" value={primaryOwner.landlineNumber} />
                    <DetailRow label="Email" value={primaryOwner.email} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Location</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="Region" value={primaryOwner.region} />
                    <DetailRow label="District" value={primaryOwner.district} />
                    <DetailRow label="Traditional Authority" value={primaryOwner.traditionalAuthority} />
                    <DetailRow label="Physical Address" value={primaryOwner.physicalAddress} />
                    <DetailRow label="Postal Address" value={primaryOwner.postalAddress} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Alternative Contact</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="Contact Name" value={primaryOwner.altContactName} />
                    <DetailRow label="Relationship" value={primaryOwner.altContactRelationship} />
                    <DetailRow label="Contact Phone" value={primaryOwner.altContactPhone} />
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
              {additionalMembers.map((member, idx) => (
                <Card key={member.id || idx}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Users className="h-5 w-5" />
                      {member.firstName} {member.lastName}
                    </CardTitle>
                    <CardDescription>Additional business member</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <DetailRow label="First Name" value={member.firstName} />
                      <DetailRow label="Last Name" value={member.lastName} />
                      <DetailRow label="Other Names" value={member.otherNames} />
                      <DetailRow label="National ID" value={member.nationalIdNumber} />
                      <DetailRow label="Nationality" value={member.nationality} />
                      <DetailRow label="Date of Birth" value={formatDate(member.dateOfBirth)} />
                      <DetailRow label="Phone Number" value={member.phoneNumber} />
                      <DetailRow label="Email" value={member.email} />
                    </div>

                    <Separator />

                    <div className="flex gap-4">
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-muted-foreground">Intern:</span>
                        <BooleanBadge value={member.isIntern} />
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-muted-foreground">Part Time:</span>
                        <BooleanBadge value={member.isPartTime} />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
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
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Users className="h-5 w-5" />
                  Employee Summary
                </CardTitle>
                <CardDescription>Breakdown of employees by type and gender</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  {/* Full Time */}
                  <div className="space-y-3">
                    <h4 className="font-semibold flex items-center gap-2">
                      <Badge variant="default">Full Time</Badge>
                    </h4>
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Males</span>
                        <Badge variant="outline">{employeeSummary.fullTimeMales || 0}</Badge>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Females</span>
                        <Badge variant="outline">{employeeSummary.fullTimeFemales || 0}</Badge>
                      </div>
                      <Separator />
                      <div className="flex justify-between items-center font-semibold">
                        <span className="text-sm">Total</span>
                        <Badge>{(employeeSummary.fullTimeMales || 0) + (employeeSummary.fullTimeFemales || 0)}</Badge>
                      </div>
                    </div>
                  </div>

                  {/* Part Time */}
                  <div className="space-y-3">
                    <h4 className="font-semibold flex items-center gap-2">
                      <Badge variant="secondary">Part Time</Badge>
                    </h4>
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Males</span>
                        <Badge variant="outline">{employeeSummary.partTimeMales || 0}</Badge>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Females</span>
                        <Badge variant="outline">{employeeSummary.partTimeFemales || 0}</Badge>
                      </div>
                      <Separator />
                      <div className="flex justify-between items-center font-semibold">
                        <span className="text-sm">Total</span>
                        <Badge>{(employeeSummary.partTimeMales || 0) + (employeeSummary.partTimeFemales || 0)}</Badge>
                      </div>
                    </div>
                  </div>

                  {/* Interns */}
                  <div className="space-y-3">
                    <h4 className="font-semibold flex items-center gap-2">
                      <Badge variant="outline">Interns</Badge>
                    </h4>
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Males</span>
                        <Badge variant="outline">{employeeSummary.internMales || 0}</Badge>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Females</span>
                        <Badge variant="outline">{employeeSummary.internFemales || 0}</Badge>
                      </div>
                      <Separator />
                      <div className="flex justify-between items-center font-semibold">
                        <span className="text-sm">Total</span>
                        <Badge>{(employeeSummary.internMales || 0) + (employeeSummary.internFemales || 0)}</Badge>
                      </div>
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="bg-muted p-4 rounded-lg">
                  <div className="flex justify-between items-center">
                    <span className="font-semibold">Grand Total Employees</span>
                    <Badge variant="default" className="text-lg px-4 py-1">
                      {(employeeSummary.fullTimeMales || 0) +
                       (employeeSummary.fullTimeFemales || 0) +
                       (employeeSummary.partTimeMales || 0) +
                       (employeeSummary.partTimeFemales || 0) +
                       (employeeSummary.internMales || 0) +
                       (employeeSummary.internFemales || 0)}
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>
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
              {/* Formalisation Score Card */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Award className="h-5 w-5" />
                      Formalisation Score
                    </div>
                    <Badge variant={formalisation.formalisationScore >= 70 ? 'default' : formalisation.formalisationScore >= 40 ? 'secondary' : 'destructive'} className="text-lg px-4 py-1">
                      {formalisation.formalisationScore || 0} / 100
                    </Badge>
                  </CardTitle>
                  <CardDescription>Overall business formalisation assessment</CardDescription>
                </CardHeader>
                <CardContent>
                  <Progress value={formalisation.formalisationScore || 0} className="h-3" />
                  <div className="mt-2 flex justify-between text-xs text-muted-foreground">
                    <span>Informal</span>
                    <span>Semi-Formal</span>
                    <span>Formal</span>
                  </div>
                </CardContent>
              </Card>

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
                          value={formalisation.hasBankAccount}
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
                          value={formalisation.hasTaxClarification}
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
                          value={formalisation.isRegisteredForVat}
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
                          value={formalisation.hasExportLicense}
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
                          value={formalisation.isMemberOfAssociation}
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
                          value={formalisation.isAffiliated}
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
                          value={formalisation.hasAccessedBds}
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
                        {formalisation.annualTurnover ? `MWK ${formalisation.annualTurnover.toLocaleString()}` : 'Not Specified'}
                      </p>
                      {formalisation.annualTurnover && (
                        <p className="text-xs text-muted-foreground">
                          {formalisation.annualTurnover >= 1000000 ? 'Above MWK 1M' : 'Below MWK 1M'}
                        </p>
                      )}
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2">
                        <Building2 className="h-4 w-4 text-muted-foreground" />
                        <p className="text-sm text-muted-foreground">Estimated Value of Assets</p>
                      </div>
                      <p className="text-2xl font-bold text-foreground">
                        {formalisation.estimatedValueOfAssets ? `MWK ${formalisation.estimatedValueOfAssets.toLocaleString()}` : 'Not Specified'}
                      </p>
                      {formalisation.estimatedValueOfAssets && (
                        <p className="text-xs text-muted-foreground">
                          {formalisation.estimatedValueOfAssets >= 5000000 ? 'High Value Assets' : formalisation.estimatedValueOfAssets >= 1000000 ? 'Medium Value Assets' : 'Low Value Assets'}
                        </p>
                      )}
                    </div>
                  </div>

                  {/* Financial Metrics Summary */}
                  {(formalisation.annualTurnover || formalisation.estimatedValueOfAssets) && (
                    <>
                      <Separator />
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                        <div>
                          <p className="text-xs text-muted-foreground">Asset to Turnover Ratio</p>
                          <p className="text-lg font-semibold">
                            {formalisation.annualTurnover && formalisation.estimatedValueOfAssets
                              ? `${(formalisation.estimatedValueOfAssets / formalisation.annualTurnover).toFixed(2)}x`
                              : '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Monthly Turnover</p>
                          <p className="text-lg font-semibold">
                            {formalisation.annualTurnover
                              ? `MWK ${Math.round(formalisation.annualTurnover / 12).toLocaleString()}`
                              : '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Daily Turnover</p>
                          <p className="text-lg font-semibold">
                            {formalisation.annualTurnover
                              ? `MWK ${Math.round(formalisation.annualTurnover / 365).toLocaleString()}`
                              : '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Business Scale</p>
                          <p className="text-lg font-semibold">
                            {formalisation.annualTurnover >= 50000000 ? 'Large' :
                             formalisation.annualTurnover >= 10000000 ? 'Medium' :
                             formalisation.annualTurnover >= 1000000 ? 'Small' : 'Micro'}
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
                    {formalisation.hasBankAccount && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Bank Account
                      </Badge>
                    )}
                    {formalisation.hasTaxClarification && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Tax Compliant
                      </Badge>
                    )}
                    {formalisation.isRegisteredForVat && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        VAT Registered
                      </Badge>
                    )}
                    {formalisation.hasExportLicense && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Export Ready
                      </Badge>
                    )}
                    {formalisation.isMemberOfAssociation && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        Association Member
                      </Badge>
                    )}
                    {formalisation.hasAccessedBds && (
                      <Badge variant="outline" className="gap-1">
                        <CheckCircle2 className="h-3 w-3" />
                        BDS Beneficiary
                      </Badge>
                    )}
                    {!formalisation.hasBankAccount && !formalisation.hasTaxClarification && !formalisation.isRegisteredForVat && !formalisation.hasExportLicense && !formalisation.isMemberOfAssociation && !formalisation.hasAccessedBds && (
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
