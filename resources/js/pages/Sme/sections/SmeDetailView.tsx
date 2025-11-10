import React, { useEffect, useState } from 'react';
import {
  Calendar, FileText, User, FolderOpen, Phone, Mail, MapPin,
  Building2, Globe, Users, Briefcase, CheckCircle2, XCircle, DollarSign
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { CrudDetailViewProps } from '@/types/crud';
import { Sme } from '@/types/sme';
import axios from 'axios';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export function SmeDetailView({
  item: sme,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Sme>) {
  const [primaryOwner, setPrimaryOwner] = useState<any>(null);
  const [additionalMembers, setAdditionalMembers] = useState<any[]>([]);
  const [employeeSummary, setEmployeeSummary] = useState<any>(null);
  const [formalisation, setFormalisation] = useState<any>(null);
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

  const DetailRow = ({ icon: Icon, label, value }: { icon: any; label: string; value: any }) => (
    <div className="flex items-start gap-3">
      <div className="p-2 rounded-lg bg-muted">
        <Icon className="h-4 w-4 text-muted-foreground" />
      </div>
      <div className="flex-1 space-y-1">
        <p className="text-sm text-muted-foreground">{label}</p>
        <p className="font-medium text-foreground">{value || '-'}</p>
      </div>
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
                <DetailRow icon={FileText} label="USME Number" value={sme.usmeNumber} />
                <DetailRow icon={Building2} label="Business Name" value={sme.name} />
                <DetailRow icon={FileText} label="Registration Number" value={sme.registrationNumber} />
                <DetailRow icon={FileText} label="Tax Identification Number" value={sme.taxIdentificationNumber} />
                <DetailRow icon={Calendar} label="Operational Since" value={formatDate(sme.operationalStartDate)} />
                <DetailRow icon={FolderOpen} label="Business Category" value={sme.businessCategory} />
                <DetailRow icon={Briefcase} label="Sector" value={sme.sector} />
                <DetailRow icon={Briefcase} label="Sub Sector" value={sme.subSector} />
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
                  <DetailRow icon={Phone} label="Phone" value={sme.contactPhone} />
                  <DetailRow icon={Mail} label="Email" value={sme.contactEmail} />
                  <DetailRow icon={Globe} label="Website" value={sme.website} />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-semibold">Location</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <DetailRow icon={MapPin} label="Region" value={sme.region} />
                  <DetailRow icon={MapPin} label="District" value={sme.district} />
                  <DetailRow icon={MapPin} label="Traditional Authority" value={sme.traditionalAuthority} />
                  <DetailRow icon={MapPin} label="Physical Address" value={sme.physicalAddress} />
                  <DetailRow icon={MapPin} label="Postal Address" value={sme.postalAddress} />
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
                  <DetailRow icon={User} label="First Name" value={primaryOwner.firstName} />
                  <DetailRow icon={User} label="Last Name" value={primaryOwner.lastName} />
                  <DetailRow icon={User} label="Other Names" value={primaryOwner.otherNames} />
                  <DetailRow icon={FileText} label="National ID" value={primaryOwner.nationalIdNumber} />
                  <DetailRow icon={Globe} label="Nationality" value={primaryOwner.nationality} />
                  <DetailRow icon={Calendar} label="Date of Birth" value={formatDate(primaryOwner.dateOfBirth)} />
                  <DetailRow icon={User} label="Gender" value={primaryOwner.gender} />
                  <DetailRow icon={FileText} label="Education Level" value={primaryOwner.educationLevel} />
                  <DetailRow icon={FileText} label="Malawian Status" value={primaryOwner.malawianStatus} />
                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <User className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Special Needs</p>
                      <BooleanBadge value={primaryOwner.hasSpecialNeeds} />
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Contact Information</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow icon={Phone} label="Phone Number" value={primaryOwner.phoneNumber} />
                    <DetailRow icon={Phone} label="Landline" value={primaryOwner.landlineNumber} />
                    <DetailRow icon={Mail} label="Email" value={primaryOwner.email} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Location</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow icon={MapPin} label="Region" value={primaryOwner.region} />
                    <DetailRow icon={MapPin} label="District" value={primaryOwner.district} />
                    <DetailRow icon={MapPin} label="Traditional Authority" value={primaryOwner.traditionalAuthority} />
                    <DetailRow icon={MapPin} label="Physical Address" value={primaryOwner.physicalAddress} />
                    <DetailRow icon={MapPin} label="Postal Address" value={primaryOwner.postalAddress} />
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Alternative Contact</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow icon={User} label="Contact Name" value={primaryOwner.altContactName} />
                    <DetailRow icon={User} label="Relationship" value={primaryOwner.altContactRelationship} />
                    <DetailRow icon={Phone} label="Contact Phone" value={primaryOwner.altContactPhone} />
                  </div>
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">No primary business owner information available</p>
              </CardContent>
            </Card>
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
                      <DetailRow icon={User} label="First Name" value={member.firstName} />
                      <DetailRow icon={User} label="Last Name" value={member.lastName} />
                      <DetailRow icon={User} label="Other Names" value={member.otherNames} />
                      <DetailRow icon={FileText} label="National ID" value={member.nationalIdNumber} />
                      <DetailRow icon={Globe} label="Nationality" value={member.nationality} />
                      <DetailRow icon={Calendar} label="Date of Birth" value={formatDate(member.dateOfBirth)} />
                      <DetailRow icon={Phone} label="Phone Number" value={member.phoneNumber} />
                      <DetailRow icon={Mail} label="Email" value={member.email} />
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
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">No additional business members</p>
              </CardContent>
            </Card>
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
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">No employee summary available</p>
              </CardContent>
            </Card>
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
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Business Formalisation
                </CardTitle>
                <CardDescription>Formal business registration and compliance status</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <FileText className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Has Bank Account</p>
                      <BooleanBadge value={formalisation.hasBankAccount} />
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <FileText className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Has Tax Clarification</p>
                      <BooleanBadge value={formalisation.hasTaxClarification} />
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <FileText className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Registered for VAT</p>
                      <BooleanBadge value={formalisation.isRegisteredForVat} />
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <Users className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Member of Association</p>
                      <BooleanBadge value={formalisation.isMemberOfAssociation} />
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <Users className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Is Affiliated</p>
                      <BooleanBadge value={formalisation.isAffiliated} />
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <FileText className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Has Export License</p>
                      <BooleanBadge value={formalisation.hasExportLicense} />
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                      <Briefcase className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="text-sm text-muted-foreground">Has Accessed BDS</p>
                      <BooleanBadge value={formalisation.hasAccessedBds} />
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Financial Information</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Annual Turnover</p>
                        <p className="font-medium text-foreground">
                          {formalisation.annualTurnover ? `MWK ${formalisation.annualTurnover.toLocaleString()}` : '-'}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-start gap-3">
                      <div className="p-2 rounded-lg bg-muted">
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <p className="text-sm text-muted-foreground">Estimated Value of Assets</p>
                        <p className="font-medium text-foreground">
                          {formalisation.estimatedValueOfAssets ? `MWK ${formalisation.estimatedValueOfAssets.toLocaleString()}` : '-'}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="bg-muted p-4 rounded-lg">
                  <div className="flex justify-between items-center">
                    <span className="font-semibold">Formalisation Score</span>
                    <Badge variant="default" className="text-lg px-4 py-1">
                      {formalisation.formalisationScore || 0}
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardContent className="py-10">
                <p className="text-center text-muted-foreground">No formalisation data available</p>
              </CardContent>
            </Card>
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
