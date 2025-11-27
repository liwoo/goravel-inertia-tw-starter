import React from 'react';
import { Calendar as BigCalendar, CalendarProps, dateFnsLocalizer } from 'react-big-calendar';
import { format, parse, startOfWeek, getDay } from 'date-fns';
import { enUS } from 'date-fns/locale';
import "react-big-calendar/lib/css/react-big-calendar.css";
import "./calendar.css";

const locales = {
    'en-US': enUS,
};

const localizer = dateFnsLocalizer({
    format,
    parse,
    startOfWeek,
    getDay,
    locales,
});

export default function Calendar<TEvent extends object = object, TResource extends object = object>(
    props: Omit<CalendarProps<TEvent, TResource>, 'localizer'>
) {
    return (
        <div className="h-[600px] w-full bg-background border rounded-md p-4">
            <BigCalendar
                localizer={localizer}
                {...props}
            />
        </div>
    );
}
