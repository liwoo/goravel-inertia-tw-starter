import React, {ReactNode} from 'react';
import {AppSidebar} from "@/components/app-sidebar";
import {SiteHeader} from "@/components/site-header";
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {usePage} from "@inertiajs/react";
import {SharedData} from "@/types/app";
import {MessageProvider} from "@/contexts/MessageContext";
import {NotificationProvider} from "@/contexts/NotificationContext";
import {UIProvider} from "@/contexts/UIContext";
import {PresenceProvider} from "@/contexts/PresenceContext";

interface AdminLayoutProps {
    title?: string;
    children: ReactNode;
}


export default function AdminLayout({title, children}: AdminLayoutProps) {
    //get user from inertia shared data
    const { props } = usePage<SharedData>();
    const user = props.auth?.user;
    const appVersion = props.appVersion;

    return (
        <UIProvider>
            <PresenceProvider>
                <MessageProvider>
                    <NotificationProvider>
                        <SidebarProvider>
                        <AppSidebar variant="inset" user={user} />
                        <SidebarInset>
                            <SiteHeader title={title || "Dashboard"}/>
                            <div className="flex flex-1 flex-col min-w-0">
                                <div className="@container/main flex flex-1 flex-col gap-2 min-w-0 overflow-hidden">
                                    {/* Page-specific content will be rendered here */}
                                    {children}
                                </div>
                                {/* Footer */}
                                <footer className="border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
                                    <div className="flex h-10 items-center justify-between px-4 text-xs text-muted-foreground">
                                        <span>
                                            &copy; {new Date().getFullYear()} Ministry of Trade and Industry &amp; SMEDI. All rights reserved.
                                        </span>
                                        <span>
                                            v{appVersion || '0.0.0'}
                                        </span>
                                    </div>
                                </footer>
                            </div>
                        </SidebarInset>
                        </SidebarProvider>
                    </NotificationProvider>
                </MessageProvider>
            </PresenceProvider>
        </UIProvider>
    );
}
