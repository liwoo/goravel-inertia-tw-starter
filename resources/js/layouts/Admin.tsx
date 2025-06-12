import React, {ReactNode} from 'react';
import {AppSidebar} from "@/components/app-sidebar";
import {SiteHeader} from "@/components/site-header";
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {usePage} from "@inertiajs/react";
import {SharedData} from "@/types/app";

interface AdminLayoutProps {
    title?: string;
    children: ReactNode;
}


export default function AdminLayout({title, children}: AdminLayoutProps) {
    //get user from inertia shared data
    const { props } = usePage<SharedData>();
    const user = props.auth?.user;
    return (
        <SidebarProvider>
            <AppSidebar variant="inset" user={user} />
            <SidebarInset>
                <SiteHeader title={title || "Dashboard"}/>
                <div className="flex flex-1 flex-col">
                    <div className="@container/main flex flex-1 flex-col gap-2">
                        {/* Page-specific content will be rendered here */}
                        {children}
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    );
}
