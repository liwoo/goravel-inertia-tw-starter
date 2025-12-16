import React from 'react';
import { Calendar as BigCalendar, CalendarProps, dateFnsLocalizer, View } from 'react-big-calendar';
import { format, parse, startOfWeek, getDay } from 'date-fns';
import { enUS } from 'date-fns/locale';
import "react-big-calendar/lib/css/react-big-calendar.css";
import "./calendar.css";
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';

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
    // Only show month view - remove week/day/agenda views
    const allowedViews: View[] = ['month'];

    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-4">
                <CardTitle className="text-2xl font-semibold tracking-tight">Event Calendar</CardTitle>
            </CardHeader>
            <CardContent className="pt-0">
                <BigCalendar
                    className="bg-card"
                    localizer={localizer}
                    views={allowedViews}
                    defaultView="month"
                    {...props}
                />
            </CardContent>
        </Card>
    );
}
