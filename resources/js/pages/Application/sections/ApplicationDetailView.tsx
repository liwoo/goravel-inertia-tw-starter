import React, { useState } from 'react';
import { Calendar, Tag, FileText, User, Mail, Phone, CheckCircle, XCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Application } from '@/types/application';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
import { router } from '@inertiajs/react';
import { toast } from 'sonner';

interface Sme {
  id: number;
  name: string;
  usme_number: string;
}

export function ApplicationDetailView({
  item: application,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Application>) {
  const [showApprovalDialog, setShowApprovalDialog] = useState(false);
  const [showRejectionDialog, setShowRejectionDialog] = useState(false);
  const [selectedSmeId, setSelectedSmeId] = useState<string>('');
  const [smes, setSmes] = useState<Sme[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingSmes, setIsLoadingSmes] = useState(false);

  const formatDate = (date: string | Date | null) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const fetchSmes = async () => {
    setIsLoadingSmes(true);
    try {
      const response = await fetch('/api/smes', {
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });

      if (response.ok) {
        const data = await response.json();
        setSmes(data.data?.data || []);
      } else {
        toast.info('Failed to load SMEs',)
      }
    } catch (error) {
      console.error('Error fetching SMEs:', error);
      toast.info('Failed to load SMEs',)
    }
    finally {
      setIsLoadingSmes(false);
    }
  };

  const handleApproveClick = () => {
    setShowApprovalDialog(true);
    fetchSmes();
  };

  const handleApprove = async () => {
    if (!selectedSmeId) {
      toast.info('Please select an SME',)
      return;
    }

    setIsLoading(true);
    try {
      const response = await fetch(`/api/applications/${application.id}/approve`, {
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
        toast.success(data.message || 'Application approved successfully',)
        setShowApprovalDialog(false);
        // Refresh the page to show updated status
        router.reload({ only: ['data'] });
      } else {
        toast.error(data.message || 'Failed to approve application',)
      }
    } catch (error) {
      console.error('Error approving application:', error);
      toast.error('Failed to approve application',)
    } finally {
      setIsLoading(false);
    }
  };

  const handleReject = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`/api/applications/${application.id}/reject`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });

      const data = await response.json();

      if (response.ok) {
        toast.success(data.message || 'Application rejected successfully',)
        setShowRejectionDialog(false);
        // Refresh the page to show updated status
        router.reload({ only: ['data'] });
      } else {
        toast.error(data.message || 'Failed to reject application',)
      }
    } catch (error) {
      console.error('Error rejecting application:', error);
      toast.error('Failed to reject application',)
    } finally {
      setIsLoading(false);
    }
  };

  const isPending = application.status === 'Pending';

  return (
    <>
      <div className="space-y-6">
        {/* Action Buttons for Pending Applications */}
        {isPending && canEdit && (
          <div className="flex gap-2 justify-end">
            <Button
              variant="outline"
              className="text-green-600 border-green-600 hover:bg-green-50"
              onClick={handleApproveClick}
            >
              <CheckCircle className="h-4 w-4 mr-2" />
              Approve
            </Button>
            <Button
              variant="outline"
              className="text-red-600 border-red-600 hover:bg-red-50"
              onClick={() => setShowRejectionDialog(true)}
            >
              <XCircle className="h-4 w-4 mr-2" />
              Reject
            </Button>
          </div>
        )}

        {/* Basic Information */}
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
            <CardDescription>Contact and registrant details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <FileText className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">SME</p>
                  <p className="font-medium text-foreground">{application.sme}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <User className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">Registrant Name</p>
                  <p className="font-medium text-foreground">{application.registrant_name}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">Email</p>
                  <p className="font-medium text-foreground">{application.email}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">Phone</p>
                  <p className="font-medium text-foreground">{application.phone}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Registration Details */}
        <Card>
          <CardHeader>
            <CardTitle>Registration Details</CardTitle>
            <CardDescription>Official registration and tax information</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <FileText className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">SME Registration Number</p>
                  <p className="font-medium text-foreground">{application.sme_registration_number}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <FileText className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">SME Tax Identification Number</p>
                  <p className="font-medium text-foreground">{application.sme_tax_identification_number}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Status */}
        <Card>
          <CardHeader>
            <CardTitle>Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Tag className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Current Status</p>
                <Badge variant={
                  application.status === 'Approved' ? 'default' :
                    application.status === 'Rejected' ? 'destructive' : 'secondary'
                }>
                  {application.status}
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Metadata Section */}
        <div className="px-1">
          <h3 className="text-sm font-medium mb-2 text-muted-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">Created</p>
              <p className="font-medium text-foreground">{formatDate(application.createdAt || application.created_at || null)}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Last Updated</p>
              <p className="font-medium text-foreground">{formatDate(application.updatedAt || application.updated_at || null)}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Approval Dialog */}
      <Dialog open={showApprovalDialog} onOpenChange={setShowApprovalDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Approve Application</DialogTitle>
            <DialogDescription>
              Select an SME to associate with this application. A user account will be created for the applicant.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="sme">Select SME *</Label>
              <Select value={selectedSmeId} onValueChange={setSelectedSmeId} disabled={isLoadingSmes}>
                <SelectTrigger>
                  <SelectValue placeholder={isLoadingSmes ? "Loading SMEs..." : "Select an SME"} />
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
              Are you sure you want to reject this application? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowRejectionDialog(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleReject} disabled={isLoading}>
              {isLoading ? 'Rejecting...' : 'Reject Application'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
