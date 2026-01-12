import React from 'react';
import { Event } from '@/types/event';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { FileText, Calendar, Hash, BookOpen, Users, MapPin, NotepadText } from 'lucide-react';

/**
 * Event table columns configuration
 */
export const eventColumns: CrudColumn<Event>[] = [
  {
    key: 'title',
    label: 'Title',
    sortable: true,
    className: 'min-w-[150px]',
    render: (event) => (
      <div className="flex items-center gap-2">
        <p className="font-mono text-sm font-medium text-foreground">{event.title}</p>
      </div>
    ),
  },
  {
    key: 'date',
    label: 'Date',
    sortable: true,
    className: 'min-w-[150px]',
    render: (event) => (
      <div className="flex items-center gap-2">
        <p className="font-medium text-foreground">
          {event.date ? new Date(event.date).toLocaleDateString() : '-'} {event.end_date ? ' - ' + new Date(event.end_date).toLocaleDateString() : ''}
        </p>
      </div>
    ),
  },
  {
    key: 'venue',
    label: 'Venue',
    sortable: true,
    className: 'min-w-[150px]',
    render: (event) => (
      <div className="flex items-center gap-2">
        <p className="font-medium text-foreground">{event.venue} ({event.district})</p>
      </div>
    ),
  },
  {
    key: 'partners',
    label: 'Partners',
    sortable: false,
    className: 'min-w-[200px]',
    render: (event) => (
      <div className="flex items-start gap-3">
        <div className="flex-1 space-y-1">
          <div className="flex flex-wrap gap-1">
            {event.partners && event.partners.length > 0 ? (
              event.partners.slice(0, 3).map((partner, index) => (
                <Badge key={index} variant="secondary" className="text-xs">
                  {partner}
                </Badge>
              ))
            ) : (
              <span className="text-sm text-muted-foreground">-</span>
            )}
            {event.partners && event.partners.length > 3 && (
              <Badge variant="outline" className="text-xs">
                +{event.partners.length - 3}
              </Badge>
            )}
          </div>
        </div>
      </div>
    ),
  },
  {
    key: 'attending_smes',
    label: 'Attending MSMEs',
    sortable: false,
    className: 'min-w-[150px]',
    render: (event) => (
      <div className="flex items-start gap-3">
        <div className="flex-1 space-y-1">
          <div className="flex flex-wrap gap-1">
            <Badge variant="outline">{event.attending_sme_details?.length || 0} MSMEs</Badge>
          </div>
        </div>
      </div>
    ),
  },
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
    render: (event) => {
      const dateValue = event.createdAt || event.created_at;
      if (!dateValue) return <span className="text-sm text-muted-foreground">-</span>;

      const date = new Date(dateValue);
      return (
        <div className="text-sm text-muted-foreground">
          {date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
          })}
        </div>
      );
    },
  },
];

/**
 * Compact columns for mobile/smaller screens
 */
export const eventColumnsMobile: CrudColumn<Event>[] = [
  {
    key: 'combined',
    label: 'Event',
    sortable: false,
    render: (event) => (
      <div className="space-y-2">
        <div className="font-medium text-foreground">{event.title}</div>
        <div className="text-sm text-muted-foreground">
          {new Date(event.createdAt || event.created_at || '').toLocaleDateString()}
        </div>
      </div>
    ),
  },
];

/**
 * Event filters configuration
 */
export const eventFilters: CrudFilter[] = [
  {
    key: 'title',
    label: 'Title',
    type: 'text',
    placeholder: 'Enter title',
  },
  {
    key: 'description',
    label: 'Description',
    type: 'text',
    placeholder: 'Enter description',
  },
  {
    key: 'date',
    label: 'Date',
    type: 'text',
    placeholder: 'Enter date',
  },
  {
    key: 'venue',
    label: 'Venue',
    type: 'text',
    placeholder: 'Enter venue',
  },
  {
    key: 'partners',
    label: 'Partners',
    type: 'text',
    placeholder: 'Enter partners',
  },
  {
    key: 'district',
    label: 'District',
    type: 'text',
    placeholder: 'Enter district',
  },
  {
    key: 'attending_smes',
    label: 'Attending MSMEs',
    type: 'number',
    placeholder: 'Enter MSME ID',
  },
  {
    key: 'notes',
    label: 'Notes',
    type: 'text',
    placeholder: 'Enter notes',
  }
];
