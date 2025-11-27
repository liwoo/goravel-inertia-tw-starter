import React, { useState } from 'react';
import { Head, router, usePage } from '@inertiajs/react';
import { SharedData } from '@/types/app.d';
import { Download } from 'lucide-react';
import {
  Event,
  EventListResponse,
  EventListRequest
} from '@/types/event';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  EventCreateForm,
  EventEditForm,
  EventDetailView,
  eventColumns,
  eventColumnsMobile,
  eventFilters
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import Calendar from '@/components/Calendar/Calendar';
import { CrudDrawer } from '@/components/Crud/CrudDrawer';
import { toast } from 'sonner';
import { format } from 'date-fns';
import { ExportDialog } from '@/components/ExportDialog';
import { ExportField, ExportOptions, ExportColumn } from '@/types/export';
import { exportData, formatDateForExport } from '@/utils/exportUtils';
import { PageAction } from '@/types/crud';

// Props interface for the Event Index page
interface EventIndexProps {
  data: EventListResponse;
  filters: EventListRequest;
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

// Event export field definitions
const eventExportFields: ExportField[] = [
  { id: 'id', label: 'ID' },
  { id: 'title', label: 'Title' },
  { id: 'description', label: 'Description' },
  { id: 'date', label: 'Date' },
  { id: 'venue', label: 'Venue' },
  { id: 'district', label: 'District' },
  { id: 'partners', label: 'Partners' },
  { id: 'notes', label: 'Notes' },
  { id: 'createdAt', label: 'Date Added' },
];

// Default fields to export
const defaultEventExportFields = [
  'title', 'date', 'venue', 'district', 'partners'
];

export default function EventIndex({
  data,
  filters,
  permissions,
  meta
}: EventIndexProps) {
  const isMobile = useIsMobile();
  const [showExportDialog, setShowExportDialog] = useState(false);
  const { props } = usePage<SharedData>();
  const currentUser = props.auth?.user;

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  // Handle export
  const handleExport = async (options: ExportOptions) => {
    // Convert ExportFields to ExportColumns
    const columns: ExportColumn[] = eventExportFields
      .filter(f => options.fields.includes(f.id))
      .map(f => ({
        id: f.id,
        label: f.label,
        formatter: f.id === 'date' || f.id === 'createdAt'
          ? (value: any) => formatDateForExport(value)
          : f.id === 'partners'
            ? (value: any) => Array.isArray(value) ? value.join(', ') : value
            : undefined,
        width: f.id === 'description' ? 40 : 20,
      }));

    // Prepare data for export
    const exportRows = data.data.map(event => ({
      ...event,
      createdAt: event.createdAt || (event as any).created_at,
    }));

    await exportData(
      {
        rows: exportRows,
        columns,
        title: 'Events Export',
      },
      {
        ...options,
        filename: `events-export-${new Date().toISOString().split('T')[0]}`,
      }
    );
  };

  // Page actions including export
  const pageActions: PageAction[] = [
    {
      key: 'export',
      label: 'Export Data',
      icon: <Download className="h-4 w-4" />,
      handler: () => setShowExportDialog(true),
    },
  ];

  const [calendarDrawerState, setCalendarDrawerState] = useState<{
    isOpen: boolean;
    type: 'create' | 'view' | undefined;
    selectedEvent?: Event;
    selectedSlot?: { start: Date; end: Date };
  }>({ isOpen: false, type: undefined });

  const calendarCreateFormRef = React.useRef<any>(null);

  const calendarEvents = React.useMemo(() => {
    return data.data.map(event => ({
      ...event,
      start: new Date(event.date),
      end: new Date(new Date(event.date).getTime() + 60 * 60 * 1000), // Default 1 hour duration
      title: event.title,
    }));
  }, [data.data]);

  const handleSelectEvent = (event: any) => {
    setCalendarDrawerState({
      isOpen: true,
      type: 'view',
      selectedEvent: event as Event
    });
  };

  const handleSelectSlot = (slotInfo: { start: Date; end: Date }) => {
    if (!permissions.canCreate) return;
    setCalendarDrawerState({
      isOpen: true,
      type: 'create',
      selectedSlot: slotInfo
    });
  };

  const closeCalendarDrawer = () => {
    setCalendarDrawerState({ isOpen: false, type: undefined });
    calendarCreateFormRef.current = null;
  };

  const handleCalendarSave = () => {
    calendarCreateFormRef.current?.handleSubmit();
  };

  return (
    <Admin title={"Event"}>
      <Head title="Event - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Calendar Widget */}
        <div className="px-0">
          <Calendar
            events={calendarEvents}
            startAccessor="start"
            endAccessor="end"
            style={{ height: 600 }}
            onSelectEvent={handleSelectEvent}
            onSelectSlot={handleSelectSlot}
            selectable
          />
        </div>

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Event>
            data={data}
            filters={filters}
            title="Events"
            resourceName="events"
            columns={isMobile ? eventColumnsMobile : eventColumns}
            customFilters={eventFilters}
            paginationConfig={meta?.pagination}
            createForm={EventCreateForm}
            editForm={EventEditForm}
            detailView={EventDetailView}
            onRefresh={handleRefresh}
            canCreate={permissions.canCreate}
            canEdit={permissions.canEdit}
            canDelete={permissions.canDelete}
            canView={true}
            pageActions={pageActions}
          />
        </div>
      </div>

      {/* Export Dialog */}
      <ExportDialog
        open={showExportDialog}
        onClose={() => setShowExportDialog(false)}
        onExport={handleExport}
        totalItems={data.data.length}
        availableFields={eventExportFields}
        defaultFields={defaultEventExportFields}
        title="Export Events"
        description={`Export ${data.data.length.toLocaleString()} event records from the current page.`}
        showStatsOption={false}
        defaultFilename="events-export"
        preparedBy={currentUser?.name}
      />


      <CrudDrawer
        isOpen={calendarDrawerState.isOpen}
        onClose={closeCalendarDrawer}
        title={calendarDrawerState.type === 'create' ? 'Create Event' : 'Event Details'}
        type={calendarDrawerState.type === 'create' ? 'create' : 'view'}
        onSave={handleCalendarSave}
        canSave={calendarDrawerState.type === 'create'}
        canEdit={permissions.canEdit && calendarDrawerState.type === 'view'}
        onEdit={() => {
          // For now, we don't support direct edit from calendar view, 
          // user can use the table below. Or we could implement it.
          toast.info("Please use the table below to edit events.");
        }}
      >
        {calendarDrawerState.type === 'create' && (
          <EventCreateForm
            ref={calendarCreateFormRef}
            onSuccess={() => {
              toast.success('Event created successfully');
              closeCalendarDrawer();
              handleRefresh();
            }}
            onError={(error) => {
              toast.error('Failed to create event');
              console.error(error);
            }}
            onCancel={closeCalendarDrawer}
            initialData={{
              date: calendarDrawerState.selectedSlot ? format(calendarDrawerState.selectedSlot.start, 'yyyy-MM-dd') : ''
            }}
          />
        )}
        {calendarDrawerState.type === 'view' && calendarDrawerState.selectedEvent && (
          <EventDetailView
            item={calendarDrawerState.selectedEvent}
            onClose={closeCalendarDrawer}
            canEdit={permissions.canEdit}
          />
        )}
      </CrudDrawer>
    </Admin >
  );
}
