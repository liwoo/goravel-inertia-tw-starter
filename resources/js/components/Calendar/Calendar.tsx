import React from 'react';
import { Calendar as BigCalendar, CalendarProps, dateFnsLocalizer } from 'react-big-calendar';
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
    return (
        <Card>
            <CardHeader>
                <CardTitle>Event Calendar</CardTitle>
            </CardHeader>
            <CardContent>
                <BigCalendar
                    className="bg-card"
                    localizer={localizer}
                    {...props}
                />
            </CardContent>

        </Card>
    );
}
