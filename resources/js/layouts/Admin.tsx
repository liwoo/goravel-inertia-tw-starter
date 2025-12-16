import React, {ReactNode, useEffect, useState} from 'react';
import {AppSidebar} from "@/components/app-sidebar";
import {SiteHeader} from "@/components/site-header";
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {usePage} from "@inertiajs/react";
import {SharedData} from "@/types/app";
import {MessageProvider} from "@/contexts/MessageContext";
import {NotificationProvider} from "@/contexts/NotificationContext";
import {UIProvider} from "@/contexts/UIContext";
import {PresenceProvider} from "@/contexts/PresenceContext";
import {MySmeModal} from "@/components/MySmeModal";
import axios from "axios";

interface AdminLayoutProps {
    title?: string;
    children: ReactNode;
}


export default function AdminLayout({title, children}: AdminLayoutProps) {
    //get user from inertia shared data
    const { props } = usePage<SharedData>();
    const user = props.auth?.user;
    const appVersion = props.appVersion;

    // MySME Modal state
    const [mySmeModalOpen, setMySmeModalOpen] = useState(false);
    const [userSmeId, setUserSmeId] = useState<number | null>(null);

    // Listen for nav-action events to open MySME modal
    useEffect(() => {
        const handleNavAction = async (event: Event) => {
            const customEvent = event as CustomEvent<{ action: string }>;
            console.log('Admin: nav-action received:', customEvent.detail);
            if (customEvent.detail?.action === 'openMySmeModal') {
                // Fetch user's linked SME ID if not already loaded
                if (!userSmeId && user?.email) {
                    console.log('Admin: Fetching SME for user email:', user.email);
                    try {
                        const response = await axios.get('/api/smes/by-email');
                        console.log('Admin: SME by-email response:', response.data);
                        const smeId = response.data?.data?.id || response.data?.data?.ID;
                        if (smeId) {
                            console.log('Admin: Setting userSmeId to:', smeId);
                            setUserSmeId(smeId);
                            setMySmeModalOpen(true);
                        } else {
                            console.error('Admin: No SME ID found in response');
                        }
                    } catch (error) {
                        console.error('Admin: Failed to fetch user SME:', error);
                    }
                } else if (userSmeId) {
                    console.log('Admin: Using cached userSmeId:', userSmeId);
                    setMySmeModalOpen(true);
                }
            }
        };

        window.addEventListener('nav-action', handleNavAction);
        return () => window.removeEventListener('nav-action', handleNavAction);
    }, [userSmeId, user?.email]);

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

                        {/* MySME Modal - accessible from nav */}
                        {userSmeId && (
                            <MySmeModal
                                open={mySmeModalOpen}
                                onOpenChange={setMySmeModalOpen}
                                smeId={userSmeId}
                            />
                        )}
                    </NotificationProvider>
                </MessageProvider>
            </PresenceProvider>
        </UIProvider>
    );
}
