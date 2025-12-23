import React from 'react';
import { Calendar, BookOpen, FileText, Hash, Users, MapPin } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { CrudDetailViewProps } from '@/types/crud';
import { Event } from '@/types/event';
export function EventDetailView({
  item: event,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Event>) {
  const formatDate = (date: string | Date | null) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: 'numeric'
    });
  };

  const attendeeCount = event.attending_sme_details?.length || 0;

  return (
    <Tabs defaultValue="details" className="w-full">
      <TabsList className="grid w-full grid-cols-2">
        <TabsTrigger value="details">Details</TabsTrigger>
        <TabsTrigger value="attendees">
          Attendees ({attendeeCount})
        </TabsTrigger>
      </TabsList>

      <TabsContent value="details" className="space-y-6 mt-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Event Information</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <BookOpen className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Title</p>
                <p className="font-medium text-foreground">{event.title}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Description</p>
                <p className="font-medium text-foreground">{event.description}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Date</p>
                <p className="font-medium text-foreground">{formatDate(event.date)}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <MapPin className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Venue</p>
                <p className="font-medium text-foreground">{event.venue}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <MapPin className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">District</p>
                <p className="font-medium text-foreground">{event.district}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Users className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Partners</p>
                <div className="flex flex-wrap gap-2 mt-1">
                  {event.partners && event.partners.length > 0 ? (
                    event.partners.map((partner, index) => (
                      <Badge key={index} variant="secondary">
                        {partner}
                      </Badge>
                    ))
                  ) : (
                    <span className="text-muted-foreground italic">No partners selected</span>
                  )}
                </div>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Notes</p>
                <p className="font-medium text-foreground whitespace-pre-wrap">{event.notes || '-'}</p>
              </div>
            </div>
          </div>
        </div>

        <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">Created</p>
              <p className="font-medium text-sm text-foreground">{formatDate(event.createdAt || event.created_at || null)}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Last Updated</p>
              <p className="font-medium text-sm text-foreground">{formatDate(event.updatedAt || event.updated_at || null)}</p>
            </div>
          </div>
        </div>
      </TabsContent>

      <TabsContent value="attendees" className="space-y-6 mt-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Attending SMEs</h3>
          {event.attending_sme_details && event.attending_sme_details.length > 0 ? (
            <div className="space-y-2">
              <div className="rounded-md border">
                <div className="grid grid-cols-2 gap-4 p-3 bg-muted/50 border-b font-medium text-sm">
                  <div>SME Name</div>
                  <div>ID</div>
                </div>
                {event.attending_sme_details.map((sme) => (
                  <div key={sme.id} className="grid grid-cols-2 gap-4 p-3 border-b last:border-b-0 hover:bg-muted/50 transition-colors">
                    <div className="font-medium text-foreground">{sme.name}</div>
                    <div className="text-muted-foreground">#{sme.id}</div>
                  </div>
                ))}
              </div>
              <p className="text-sm text-muted-foreground mt-2">
                Total: {event.attending_sme_details.length} {event.attending_sme_details.length === 1 ? 'attendee' : 'attendees'}
              </p>
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center border rounded-md bg-muted/20">
              <Users className="h-12 w-12 text-muted-foreground/50 mb-3" />
              <p className="text-muted-foreground font-medium">No attendees yet</p>
              <p className="text-sm text-muted-foreground mt-1">SMEs can register for this event through their portal</p>
            </div>
          )}
        </div>
      </TabsContent>
    </Tabs>
  );
}
