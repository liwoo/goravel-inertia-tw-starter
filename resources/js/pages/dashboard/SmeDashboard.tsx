import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { router } from '@inertiajs/react';

interface SmeDashboardProps {
    user: any;
}

export const SmeDashboard: React.FC<SmeDashboardProps> = ({ user }) => {
    router.visit('/portal', { replace: true });
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

            {/* Main Content Area */}
            <Card className="min-h-[300px]">
                <CardHeader>
                    <CardTitle>My Business Profile</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="flex items-center justify-center h-40 border-2 border-dashed rounded-lg">
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
