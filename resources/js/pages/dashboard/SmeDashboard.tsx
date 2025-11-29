import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";

interface SmeDashboardProps {
    user: any;
}

export const SmeDashboard: React.FC<SmeDashboardProps> = ({ user }) => {
    return (
        <div className="flex flex-col gap-6 p-4 md:p-6">
            {/* Welcome Header */}
            <div className="flex flex-col gap-1">
                <h1 className="text-2xl font-semibold tracking-tight">
                    Welcome back, {user?.name}!
                </h1>
                <p className="text-muted-foreground">
                    Here is your MSME dashboard.
                </p>
            </div>

            {/* Formalisation Progress */}
            <Card className="w-full">
                <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-medium">
                        Formalisation Progress
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="flex items-center gap-4">
                        <Progress value={78} className="h-4 flex-1" />
                        <span className="text-sm font-bold">78% Formalisation</span>
                    </div>
                </CardContent>
            </Card>

            {/* Widgets Row */}
            <div className="grid gap-6 md:grid-cols-2">
                {/* Events Widget */}
                <Card>
                    <CardHeader>
                        <CardTitle>Events</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex flex-col gap-4">
                            <p className="text-sm text-muted-foreground">
                                Upcoming events relevant to your business.
                            </p>
                            {/* Placeholder for events list */}
                            <div className="space-y-2">
                                <div className="rounded-md border p-3">
                                    <h4 className="font-semibold">SME Workshop</h4>
                                    <p className="text-xs text-muted-foreground">Nov 30, 2025</p>
                                </div>
                                <div className="rounded-md border p-3">
                                    <h4 className="font-semibold">Networking Mixer</h4>
                                    <p className="text-xs text-muted-foreground">Dec 05, 2025</p>
                                </div>
                            </div>
                            <Button variant="outline" className="w-full">View All Events</Button>
                        </div>
                    </CardContent>
                </Card>

                {/* Opportunities Widget */}
                <Card>
                    <CardHeader>
                        <CardTitle>Opportunities</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex flex-col gap-4">
                            <p className="text-sm text-muted-foreground">
                                Business opportunities and procurement notices.
                            </p>
                            {/* Placeholder for opportunities list */}
                            <div className="space-y-2">
                                <div className="rounded-md border p-3">
                                    <h4 className="font-semibold">Supply of Office Stationery</h4>
                                    <p className="text-xs text-muted-foreground">Closing: Dec 10, 2025</p>
                                </div>
                                <div className="rounded-md border p-3">
                                    <h4 className="font-semibold">Catering Services Tender</h4>
                                    <p className="text-xs text-muted-foreground">Closing: Dec 15, 2025</p>
                                </div>
                            </div>
                            <Button variant="outline" className="w-full">View All Opportunities</Button>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Main Content Area */}
            <Card className="min-h-[300px]">
                <CardHeader>
                    <CardTitle>My Business Profile</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="flex items-center justify-center h-40 border-2 border-dashed rounded-lg">
                        <p className="text-muted-foreground">Main content area for business details, stats, etc.</p>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
