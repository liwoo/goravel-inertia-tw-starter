import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { router } from "@inertiajs/react";

interface Event {
    id: number;
    title: string;
    date: string;
    venue: string;
    district?: string;
}

interface Procurement {
    id: number;
    organization: string;
    refNo: string;
    closeDate: string;
    procurementType?: string;
}

interface Formalisation {
    formalisation_score: number;
    // Add other fields if needed
}

interface PortalCardsProps {
    events?: Event[];
    procurements?: Procurement[];
    formalisation?: Formalisation;
}

export function PortalCards({ events = [], procurements = [], formalisation }: PortalCardsProps) {
    const formalisationScore = formalisation?.formalisation_score || 0;

    return (
        <>
            <Card className="w-full">
                <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-medium">
                        Formalisation Progress
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="flex items-center gap-4">
                        <Progress value={formalisationScore} className="h-4 flex-1" />
                        <span className="text-sm font-bold">{formalisationScore}% Formalisation</span>
                    </div>
                </CardContent>
            </Card>

            <div className="grid gap-6 md:grid-cols-2">
                <Card>
                    <CardHeader>
                        <CardTitle>Events</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex flex-col gap-4">
                            <p className="text-sm text-muted-foreground">
                                Upcoming events relevant to your business.
                            </p>
                            {events.length > 0 ? (
                                <div className="space-y-2">
                                    {events.map((event) => (
                                        <div key={event.id} className="rounded-md border p-3">
                                            <h4 className="font-semibold">{event.title}</h4>
                                            <p className="text-xs text-muted-foreground">{new Date(event.date).toLocaleDateString()}</p>
                                            <p className="text-xs text-muted-foreground">{event.venue}</p>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-sm text-muted-foreground italic">No upcoming events found.</p>
                            )}
                            <Button variant="outline" className="w-full" onClick={() => router.visit('/events')}>View All Events</Button>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>Opportunities</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex flex-col gap-4">
                            <p className="text-sm text-muted-foreground">
                                Business opportunities and procurement notices.
                            </p>
                            {procurements.length > 0 ? (
                                <div className="space-y-2">
                                    {procurements.map((proc) => (
                                        <div key={proc.id} className="rounded-md border p-3">
                                            <h4 className="font-semibold">{proc.organization}</h4>
                                            <p className="text-xs text-muted-foreground">{proc.refNo}</p>
                                            <p className="text-xs text-muted-foreground">Closing: {new Date(proc.closeDate).toLocaleDateString()}</p>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-sm text-muted-foreground italic">No active opportunities found.</p>
                            )}
                            <Button variant="outline" className="w-full" onClick={() => router.visit('/procurement-notices')}>View All Opportunities</Button>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </>
    )
}