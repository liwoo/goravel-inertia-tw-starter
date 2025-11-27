import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
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

export default function EventIndex({
  data,
  filters,
  permissions,
  meta
}: EventIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

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
          />
        </div>
      </div>


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
