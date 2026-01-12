import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import { CheckCircle, XCircle, Plus, Copy, Mail, FileEdit } from 'lucide-react';
import {
  Application,
  ApplicationListResponse,
  ApplicationListRequest,
  ApplicationType
} from '@/types/application';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  ApplicationDetailView,
  applicationColumns,
  applicationColumnsMobile,
  applicationFilters,
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import { CrudAction } from '@/types/crud';
import { toast } from 'sonner';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { SmeCreateForm, type InitialSmeData } from '@/pages/Sme/sections';
import { Textarea } from '@/components/ui/textarea';

// Props interface for the Application Index page
interface ApplicationIndexProps {
  data: ApplicationListResponse;
  filters: ApplicationListRequest;
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
  };
  meta?: {
    pagination: {
      defaultPageSize: number;
      maxPageSize: number;
      allowedSizes: number[];
    };
  };
}

interface Sme {
  id: number;
  name: string;
  usme_number: string;
}

export default function ApplicationIndex({
  data,
  filters,
  permissions,
  meta
}: ApplicationIndexProps) {
  const isMobile = useIsMobile();
  const [showApprovalDialog, setShowApprovalDialog] = useState(false);
  const [showRejectionDialog, setShowRejectionDialog] = useState(false);
  const [showCreateSmeDialog, setShowCreateSmeDialog] = useState(false);
  const [showWelcomeEmailDialog, setShowWelcomeEmailDialog] = useState(false);
  const [selectedApplication, setSelectedApplication] = useState<Application | null>(null);
  const [selectedSmeId, setSelectedSmeId] = useState<string>('');
  const [smes, setSmes] = useState<Sme[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingSmes, setIsLoadingSmes] = useState(false);
  const [needsReload, setNeedsReload] = useState(false);
  const [rejectionReason, setRejectionReason] = useState<string>('');
  const [showAmendmentConfirmDialog, setShowAmendmentConfirmDialog] = useState(false);
  const [showSignupRejectConfirmDialog, setShowSignupRejectConfirmDialog] = useState(false);
  const [welcomeEmailData, setWelcomeEmailData] = useState<{
    username: string;
    password: string;
    loginUrl: string;
    smeName: string;
    applicantName: string;
  } | null>(null);

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  // Fetch SMEs that match the applicant's email (primary owner email match)
  const fetchSmes = async (applicantEmail?: string) => {
    setIsLoadingSmes(true);
    try {
      // If we have an applicant email, filter by primary owner email
      const url = applicantEmail
        ? `/api/smes/by-owner-email?email=${encodeURIComponent(applicantEmail)}`
        : '/api/smes';

      const response = await fetch(url, {
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });

      if (response.ok) {
        const data = await response.json();
        // Handle both filtered response (array) and paginated response (nested data)
        const smesData = Array.isArray(data.data) ? data.data : (data.data?.data || []);
        setSmes(smesData);

        if (applicantEmail && smesData.length === 0) {
          toast.info('No MSMEs found matching the applicant\'s email. You can create a new MSME.');
        }
      } else {
        toast.info('Failed to load MSMEs');
      }
    } catch (error) {
      console.error('Error fetching SMEs:', error);
      toast.info('Failed to load SMEs');
    } finally {
      setIsLoadingSmes(false);
    }
  };

  const handleApproveClick = (application: Application) => {
    if (application.status !== 'Pending') {
      toast.info('Application has already been ' + application.status);
      return;
    }
    setSelectedApplication(application);

    // Amendment applications don't need SME selection - show confirmation first
    if (application.type === 'amend_formalisation') {
      setShowAmendmentConfirmDialog(true);
      return;
    }

    // Signup applications need SME selection
    // Filter SMEs by applicant's email to only show SMEs they own
    setShowApprovalDialog(true);
    fetchSmes(application.email);
  };

  const handleApproveAmendment = async (application: Application) => {
    setIsLoading(true);
    try {
      const response = await fetch(`/api/applications/${application.id}/approve`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({}),
      });

      const data = await response.json();

      if (response.ok) {
        toast.success(data.message || 'Amendment application approved successfully');
        setSelectedApplication(null);
        router.reload({ only: ['data'] });
      } else {
        toast.error(data.message || 'Failed to approve amendment application');
      }
    } catch (error) {
      console.error('Error approving amendment application:', error);
      toast.error('Failed to approve amendment application');
    } finally {
      setIsLoading(false);
    }
  };

  const handleApprove = async () => {
    if (!selectedApplication) return;

    if (!selectedSmeId) {
      toast.info('Please select an MSME');
      return;
    }

    setIsLoading(true);
    try {
      const response = await fetch(`/api/applications/${selectedApplication.id}/approve`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({ sme_id: parseInt(selectedSmeId) }),
      });

      const data = await response.json();

      if (response.ok) {
        toast.success(data.message || 'Application approved successfully');
        setShowApprovalDialog(false);

        // Show welcome email dialog with user credentials
        const selectedSme = smes.find(s => s.id === parseInt(selectedSmeId));
        setWelcomeEmailData({
          username: data.data?.user?.email,
          password: data.data?.user?.password || 'Generated by system',
          loginUrl: window.location.origin + '/login',
          smeName: selectedSme?.name || '',
          applicantName: selectedApplication.registrant_name || `${selectedApplication.first_name} ${selectedApplication.last_name}`,
        });
        setShowWelcomeEmailDialog(true);
        setNeedsReload(true);

        setSelectedApplication(null);
        setSelectedSmeId('');
      } else {
        toast.error(data.message || 'Failed to approve application');
      }
    } catch (error) {
      console.error('Error approving application:', error);
      toast.error('Failed to approve application');
    } finally {
      setIsLoading(false);
    }
  };

  // Handle reject click
  const handleRejectClick = (application: Application) => {
    if (application.status !== 'Pending') {
      toast.info('Application has already been ' + application.status);
      return;
    }
    setSelectedApplication(application);
    setRejectionReason('');

    // Signup applications need confirmation before showing rejection dialog
    if (application.type === 'signup') {
      setShowSignupRejectConfirmDialog(true);
      return;
    }

    // Amendment applications go directly to rejection reason dialog
    setShowRejectionDialog(true);
  };

  // Handle reject submission
  const handleReject = async () => {
    if (!selectedApplication) return;

    if (!rejectionReason.trim()) {
      toast.info('Please provide a reason for rejection');
      return;
    }

    setIsLoading(true);
    try {
      const response = await fetch(`/api/applications/${selectedApplication.id}/reject`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({ reason: rejectionReason.trim() }),
      });

      const data = await response.json();

      if (response.ok) {
        toast.success(data.message || 'Application rejected successfully');
        setShowRejectionDialog(false);
        setSelectedApplication(null);
        setRejectionReason('');
        router.reload({ only: ['data'] });
      } else {
        toast.error(data.message || 'Failed to reject application');
      }
    } catch (error) {
      console.error('Error rejecting application:', error);
      toast.error('Failed to reject application');
    } finally {
      setIsLoading(false);
    }
  };

  // Handle SME creation success - auto-approve the application with the new SME
  const handleSmeCreated = async (newSme: any) => {
    toast.success('MSME created successfully');
    setShowCreateSmeDialog(false);

    if (!selectedApplication || !newSme?.id) {
      await fetchSmes();
      return;
    }

    // Auto-approve the application with the newly created SME
    setIsLoading(true);
    try {
      const response = await fetch(`/api/applications/${selectedApplication.id}/approve`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({ sme_id: newSme.id }),
      });

      const data = await response.json();

      if (response.ok) {
        toast.success(data.message || 'Application approved successfully');
        setShowApprovalDialog(false);

        // Show welcome email dialog with user credentials
        setWelcomeEmailData({
          username: data.data?.user?.email,
          password: data.data?.user?.password || 'Generated by system',
          loginUrl: window.location.origin + '/login',
          smeName: newSme.name || '',
          applicantName: selectedApplication.registrant_name || `${selectedApplication.first_name} ${selectedApplication.last_name}`,
        });
        setShowWelcomeEmailDialog(true);
        setNeedsReload(true);

        setSelectedApplication(null);
        setSelectedSmeId('');
      } else {
        toast.error(data.message || 'Failed to approve application');
        // Still refresh SMEs list in case of error
        await fetchSmes();
        if (newSme?.id) {
          setSelectedSmeId(newSme.id.toString());
        }
      }
    } catch (error) {
      console.error('Error approving application:', error);
      toast.error('Failed to approve application');
      // Still refresh SMEs list in case of error
      await fetchSmes();
      if (newSme?.id) {
        setSelectedSmeId(newSme.id.toString());
      }
    } finally {
      setIsLoading(false);
    }
  };

  // Generate welcome email text
  const generateWelcomeEmail = () => {
    if (!welcomeEmailData) return '';

    return `Subject: Welcome to the SMEDI Database - Your Account Has Been Approved

Dear ${welcomeEmailData.applicantName},

Congratulations! Your application has been approved and your account has been successfully created.

Your account details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SME: ${welcomeEmailData.smeName}
Username: ${welcomeEmailData.username}
Password: ${welcomeEmailData.password}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

You can now log in to your account using the following link:
${welcomeEmailData.loginUrl}

For security reasons, we recommend changing your password after your first login.

If you have any questions or need assistance, please don't hesitate to contact our support team.

Best regards,
SMEDI Team`;
  };

  const handleCopyWelcomeEmail = async () => {
    const emailText = generateWelcomeEmail();
    try {
      await navigator.clipboard.writeText(emailText);
      toast.success('Welcome email copied to clipboard');
    } catch (error) {
      toast.error('Failed to copy to clipboard');
    }
  };

  // Handle welcome email dialog close
  const handleWelcomeEmailDialogClose = (open: boolean) => {
    setShowWelcomeEmailDialog(open);
    if (!open && needsReload) {
      setNeedsReload(false);
      router.reload({ only: ['data'] });
    }
  };

  // Map application data to SME initial data for prefilling the form
  const getInitialSmeDataFromApplication = (application: Application | null): InitialSmeData | undefined => {
    if (!application) return undefined;

    return {
      // Business Information (Step 1)
      name: application.sme || '',
      registrationNumber: application.sme_registration_number || '',
      taxIdentificationNumber: application.sme_tax_identification_number || '',
      contactPhone: application.phone || '',
      contactEmail: application.email || '',
      physicalAddress: application.physical_address || '',
      postalAddress: application.postal_address || '',
      district: application.district || '',
      traditionalAuthority: application.traditional_authority || '',
      // Primary Owner (Step 2)
      ownerFirstName: application.first_name || '',
      ownerLastName: application.last_name || '',
      ownerOtherNames: application.other_names || '',
      ownerNationality: application.nationality || '',
      ownerNationalIdNumber: application.national_id_number || '',
      ownerDateOfBirth: application.date_of_birth || '',
      ownerGender: application.gender || '',
      ownerEducationLevel: application.education_level || '',
      ownerMalawianStatus: application.malawian_status || '',
      ownerHasSpecialNeeds: application.has_special_needs || false,
      ownerPhoneNumber: application.phone || '',
      ownerLandlineNumber: application.landline_number || '',
      ownerEmail: application.email || '',
      ownerPhysicalAddress: application.physical_address || '',
      ownerPostalAddress: application.postal_address || '',
      ownerDistrict: application.district || '',
      ownerTraditionalAuthority: application.traditional_authority || '',
      ownerAltContactName: application.alt_contact_name || '',
      ownerAltContactRelationship: application.alt_contact_relationship || '',
      ownerAltContactPhone: application.alt_contact_phone || '',
    };
  };

  // Custom row actions for approve/reject - only show for pending applications
  const customActions: CrudAction<Application>[] = permissions.canEdit ? [
    {
      key: 'approve',
      label: 'Approve',
      icon: <CheckCircle className="w-4 h-4" />,
      onClick: handleApproveClick,
      hidden: (app) => app.status !== 'Pending',
    },
    {
      key: 'reject',
      label: 'Reject',
      icon: <XCircle className="w-4 h-4" />,
      onClick: handleRejectClick,
      hidden: (app) => app.status !== 'Pending',
    },
  ] : [];

  return (
    <Admin title={"Application"}>
      <Head title="Application - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Application>
            data={data}
            filters={filters}
            title="Applications"
            resourceName="applications"
            columns={isMobile ? applicationColumnsMobile : applicationColumns}
            customFilters={applicationFilters}
            paginationConfig={meta?.pagination}
            detailView={ApplicationDetailView}
            onRefresh={handleRefresh}
            canView={true}
            readOnly={true}
            actions={customActions}
          />
        </div>
      </div>

      {/* Approval Dialog */}
      <Dialog open={showApprovalDialog} onOpenChange={setShowApprovalDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Approve Application</DialogTitle>
            <DialogDescription>
              Select an MSME to link with this applicant. Only MSMEs where the primary owner's email matches
              the applicant's email ({selectedApplication?.email}) are shown. A user account will be created.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="sme">Select MSME *</Label>
              <Select value={selectedSmeId} onValueChange={setSelectedSmeId} disabled={isLoadingSmes}>
                <SelectTrigger>
                  <SelectValue placeholder={isLoadingSmes ? "Loading MSMEs..." : "Select an MSME"} />
                </SelectTrigger>
                <SelectContent>
                  {smes.map((sme) => (
                    <SelectItem key={sme.id} value={sme.id.toString()}>
                      {sme.name} ({sme.usme_number})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <Button
                variant="outline"
                className="w-full"
                onClick={() => setShowCreateSmeDialog(true)}
              >
                <Plus className="h-4 w-4 mr-2" />
                Create New MSME
              </Button>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowApprovalDialog(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button onClick={handleApprove} disabled={isLoading || !selectedSmeId}>
              {isLoading ? 'Approving...' : 'Approve Application'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Rejection Dialog */}
      <Dialog open={showRejectionDialog} onOpenChange={setShowRejectionDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Application</DialogTitle>
            <DialogDescription>
              {selectedApplication?.type === 'amend_formalisation'
                ? 'Reject this formalisation amendment request. The applicant will be notified.'
                : 'Reject this signup application. This action cannot be undone.'}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="rejection-reason">Reason for Rejection *</Label>
              <Textarea
                id="rejection-reason"
                placeholder="Please provide a reason for rejecting this application..."
                value={rejectionReason}
                onChange={(e) => setRejectionReason(e.target.value)}
                className="min-h-[100px]"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowRejectionDialog(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleReject} disabled={isLoading || !rejectionReason.trim()}>
              {isLoading ? 'Rejecting...' : 'Reject Application'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Welcome Email Dialog */}
      <Dialog open={showWelcomeEmailDialog} onOpenChange={handleWelcomeEmailDialogClose}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Mail className="h-5 w-5" />
              Welcome Email Preview
            </DialogTitle>
            <DialogDescription>
              Copy this welcome email to send to the applicant with their login credentials.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <Textarea
              value={generateWelcomeEmail()}
              readOnly
              className="min-h-[400px] font-mono text-sm"
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => handleWelcomeEmailDialogClose(false)}>
              Close
            </Button>
            <Button onClick={handleCopyWelcomeEmail}>
              <Copy className="h-4 w-4 mr-2" />
              Copy to Clipboard
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Create MSME Dialog - Full Screen */}
      <Dialog open={showCreateSmeDialog} onOpenChange={setShowCreateSmeDialog}>
        <DialogContent className="!max-w-[95vw] !w-[95vw] !h-[90vh] !max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create New MSME</DialogTitle>
            <DialogDescription>
              Create a new MSME to associate with this application. Form has been prefilled with application data.
            </DialogDescription>
          </DialogHeader>
          <SmeCreateForm
            onSuccess={() => {}} // No-op since we handle the SME data via onSmeCreated
            onSmeCreated={handleSmeCreated}
            onError={(error) => toast.error(error?.message || 'Failed to create MSME')}
            onCancel={() => setShowCreateSmeDialog(false)}
            initialData={getInitialSmeDataFromApplication(selectedApplication)}
          />
        </DialogContent>
      </Dialog>

      {/* Amendment Approval Confirmation Dialog */}
      <Dialog open={showAmendmentConfirmDialog} onOpenChange={setShowAmendmentConfirmDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Approve Amendment Request</DialogTitle>
            <DialogDescription>
              Are you sure you want to approve this formalisation amendment request?
              This will update the MSME's formalisation data with the proposed changes.
            </DialogDescription>
          </DialogHeader>
          <div className="py-4">
            {selectedApplication && (
              <div className="bg-muted p-3 rounded-md text-sm space-y-1">
                <p><strong>MSME:</strong> {selectedApplication.sme}</p>
                <p><strong>Submitted by:</strong> {selectedApplication.email}</p>
              </div>
            )}
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setShowAmendmentConfirmDialog(false);
                setSelectedApplication(null);
              }}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button
              onClick={() => {
                setShowAmendmentConfirmDialog(false);
                if (selectedApplication) {
                  handleApproveAmendment(selectedApplication);
                }
              }}
              disabled={isLoading}
            >
              {isLoading ? 'Approving...' : 'Yes, Approve Amendment'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Signup Rejection Confirmation Dialog */}
      <Dialog open={showSignupRejectConfirmDialog} onOpenChange={setShowSignupRejectConfirmDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Signup Application</DialogTitle>
            <DialogDescription>
              Are you sure you want to reject this signup application?
              The applicant will not be able to access the system.
            </DialogDescription>
          </DialogHeader>
          <div className="py-4">
            {selectedApplication && (
              <div className="bg-muted p-3 rounded-md text-sm space-y-1">
                <p><strong>MSME:</strong> {selectedApplication.sme}</p>
                <p><strong>Applicant:</strong> {selectedApplication.registrant_name || `${selectedApplication.first_name} ${selectedApplication.last_name}`}</p>
                <p><strong>Email:</strong> {selectedApplication.email}</p>
              </div>
            )}
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setShowSignupRejectConfirmDialog(false);
                setSelectedApplication(null);
              }}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                setShowSignupRejectConfirmDialog(false);
                setShowRejectionDialog(true);
              }}
            >
              Yes, Proceed to Rejection
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Admin>
  );
}
