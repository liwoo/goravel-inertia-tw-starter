import React from 'react';
import { Calendar, User, MapPin, Tag, FileText, Building2, Briefcase, Clock, DollarSign } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Bdsp } from '@/types/bdsp';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { CopyableText } from '@/components/ui/copyable-text';

export function BdspDetailView({
  item: bdsp,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Bdsp>) {
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

  return (
    <div className="space-y-6">

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            BDSP Information
          </CardTitle>
          <CardDescription>Core details and location</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">UBDSP Number</p>
                <CopyableText
                  value={bdsp.ubdsp_number || (bdsp as any).ubdspNumber}
                  className="font-mono font-medium text-foreground"
                  iconSize="md"
                />
              </div>
            </div>

            <DetailRow icon={User} label="Name" value={bdsp.name} />
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Registration Status</p>
                <Badge
                  variant={
                    (bdsp.registration_status || '').toLowerCase() === 'active' ? 'default' :
                      (bdsp.registration_status || '').toLowerCase() === 'rejected' ? 'destructive' :
                        (bdsp.registration_status || '').toLowerCase() === 'inactive' ? 'outline' : 'secondary'
                  }
                  className="capitalize"
                >
                  {bdsp.registration_status || 'Pending'}
                </Badge>
              </div>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <h4 className="font-semibold">Location</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <DetailRow icon={MapPin} label="Postal Address" value={bdsp.postal_address} />
              <DetailRow icon={MapPin} label="Physical Address" value={bdsp.physical_address} />
            </div>
          </div>
        </CardContent>
      </Card>

      <Separator />

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Briefcase className="h-5 w-5" />
            Services Offered
          </CardTitle>
          <CardDescription>List of services provided by this BDSP</CardDescription>
        </CardHeader>
        <CardContent>
          {bdsp.service_list && bdsp.service_list.length > 0 ? (
            <div className="border rounded-lg overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Service Name</TableHead>
                    <TableHead>Cost</TableHead>
                    <TableHead>Duration</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {bdsp.service_list.map((service, index) => (
                    <TableRow key={index}>
                      <TableCell className="font-medium pl-4">{service.name}</TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <DollarSign className="h-3 w-3 text-muted-foreground" />
                          {service.cost.toLocaleString()}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <Clock className="h-3 w-3 text-muted-foreground" />
                          {service.duration}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              No services listed
            </div>
          )}
        </CardContent>
      </Card>

      <Separator />

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Tag className="h-5 w-5" />
            Partners & Products
          </CardTitle>
          <CardDescription>Collaborations and product offerings</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div>
            <h4 className="text-sm font-medium mb-3 text-muted-foreground">Partners</h4>
            <div className="flex flex-wrap gap-2">
              {bdsp.partners && bdsp.partners.length > 0 ? (
                bdsp.partners.map((partner, idx) => (
                  <Badge key={idx} variant="secondary" className="px-3 py-1">
                    {partner}
                  </Badge>
                ))
              ) : (
                <span className="text-sm text-muted-foreground">No partners listed</span>
              )}
            </div>
          </div>

          <Separator />

          <div>
            <h4 className="text-sm font-medium mb-3 text-muted-foreground">Product Types</h4>
            <div className="flex flex-wrap gap-2">
              {bdsp.product_types && bdsp.product_types.length > 0 ? (
                bdsp.product_types.map((product, idx) => (
                  <Badge key={idx} variant="outline" className="px-3 py-1">
                    {product}
                  </Badge>
                ))
              ) : (
                <span className="text-sm text-muted-foreground">No product types listed</span>
              )}
            </div>
          </div>
        </CardContent>
      </Card>


      {/* Metadata Section */}
      <div>
        <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Created</p>
            <p className="font-medium text-sm text-foreground">{formatDate(bdsp.createdAt || bdsp.created_at || null)}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Last Updated</p>
            <p className="font-medium text-sm text-foreground">{formatDate(bdsp.updatedAt || bdsp.updated_at || null)}</p>
          </div>
        </div>
      </div>
    </div>
  );
}

