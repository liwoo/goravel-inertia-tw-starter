import React from 'react';
import { Member } from '@/types/member';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { User, Mail } from 'lucide-react';

/**
 * Member table columns configuration
 */
export const memberColumns: CrudColumn<Member>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    className: 'min-w-[150px]',
    render: (member) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <User className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{member.first_name} {member.last_name}</p>
          {member.other_names && <p className="text-xs text-muted-foreground">{member.other_names}</p>}
        </div>
      </div>
    ),
  },
  {
    key: 'email',
    label: 'Contact',
    sortable: true,
    className: 'min-w-[150px]',
    render: (member) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <Mail className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{member.email || '-'}</p>
          <p className="text-xs text-muted-foreground">{member.phone_number}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'national_id_number',
    label: 'National ID',
    sortable: true,
    className: 'min-w-[120px]',
    render: (member) => (
      <div className="text-sm font-medium">{member.national_id_number}</div>
    ),
  },
  {
    key: 'status',
    label: 'Status',
    sortable: false,
    className: 'w-32',
    render: (member) => (
      <div className="flex flex-col gap-1">
        {member.is_intern && <Badge variant="secondary" className="w-fit">Intern</Badge>}
        {member.is_part_time && <Badge variant="outline" className="w-fit">Part Time</Badge>}
        {!member.is_intern && !member.is_part_time && <Badge variant="default" className="w-fit">Full Time</Badge>}
      </div>
    ),
  },
];

/**
 * Compact columns for mobile/smaller screens
 */
export const memberColumnsMobile: CrudColumn<Member>[] = [
  {
    key: 'combined',
    label: 'Member',
    sortable: false,
    render: (member) => (
      <div className="space-y-2">
        <div className="font-medium text-foreground">{member.name}</div>
        <div className="text-sm text-muted-foreground">
          {new Date(member.createdAt || member.created_at || '').toLocaleDateString()}
        </div>
      </div>
    ),
  },
];

/**
 * Member filters configuration
 */
export const memberFilters: CrudFilter[] = [
  {
    key: 'name',
    label: 'Name',
    type: 'text',
    placeholder: 'Enter name',
  },
  {
    key: 'email',
    label: 'Email',
    type: 'text',
    placeholder: 'Enter email',
  }
];
