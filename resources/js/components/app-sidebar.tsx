import * as React from "react"
import { ShieldIcon, Search, Command } from "lucide-react"

import { NavDocuments } from "@/components/nav-documents"
import { NavMain } from "@/components/nav-main"
import { NavSecondary } from "@/components/nav-secondary"
import { NavUser } from "@/components/nav-user"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
} from "@/components/ui/sidebar"
import { usePermissions } from "@/contexts/PermissionsContext"
import { navigationConfig } from "@/config/navigation"
import { GlobalSearch } from "@/components/GlobalSearch"
import { useUI } from "@/contexts/UIContext"

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
    user?: any;
}

export function AppSidebar({ user, ...props }: AppSidebarProps) {
    const { canPerformAction, isSuperAdmin: checkIsSuperAdmin, isAdmin } = usePermissions();
    const [searchOpen, setSearchOpen] = React.useState(false);
    const { openMessages, openNotifications } = useUI();

    // Filter navigation items based on permissions
    const navigationItems = React.useMemo(() => {
        const userRoles = user?.roles || [];

        // Helper to check role requirement
        const hasRequiredRole = (item: any) => {
            if (!item.requiredRole) return true;
            return userRoles.some((role: any) => role.slug === item.requiredRole);
        };

        // Helper to check all requirements for an item
        const checkItemRequirements = (item: any) => {
            // Check super admin requirement
            if (item.requireSuperAdmin) {
                return checkIsSuperAdmin();
            }

            // Check role requirement
            if (!hasRequiredRole(item)) return false;

            // If no permission requirement, show the item
            if (!item.requiredService && !item.requiredAction) {
                return true;
            }

            // Check service permission
            if (item.requiredService && item.requiredAction) {
                return canPerformAction(item.requiredService, item.requiredAction);
            }

            return true;
        };

        // Determine which main navigation to use
        let activeNavMain = navigationConfig.navMain;

        // Filter main navigation
        const filteredNavMain = activeNavMain.filter(checkItemRequirements);

        // Filter secondary navigation
        const filteredNavSecondary = navigationConfig.navSecondary.filter(checkItemRequirements);

        // Filter documents
        const filteredDocuments = navigationConfig.documents.filter(checkItemRequirements);

        return {
            navMain: filteredNavMain,
            navSecondary: filteredNavSecondary,
            documents: filteredDocuments,
        };
    }, [canPerformAction, checkIsSuperAdmin, user]);

    // Global keyboard shortcut for search
    React.useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
                e.preventDefault();
                setSearchOpen(true);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

    return (
        <>
            <Sidebar collapsible="icon" {...props}>
                <SidebarHeader>
                    <div className="flex items-center justify-between px-1 py-2">
                        <a href="/dashboard" className="flex items-center gap-2 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:w-full">
                            <img src="/images/logo.svg" alt="Logo" className="h-8 w-auto shrink-0" />
                            <span className="text-lg font-semibold group-data-[collapsible=icon]:hidden">Starter Project</span>
                        </a>
                    </div>
                </SidebarHeader>
                <SidebarContent>
                    <NavMain items={navigationItems.navMain} />
                    <NavDocuments items={navigationItems.documents} />
                    <NavSecondary items={navigationItems.navSecondary} className="mt-auto" />

                    {/* Hardcoded Search Option */}
                    <div className="mt-2 px-3 pb-3 group-data-[collapsible=icon]:px-2">
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton tooltip="Search" onClick={() => setSearchOpen(true)}>
                                    <Search className="h-4 w-4" />
                                    <span>Search</span>
                                    <kbd className="ml-auto pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 group-data-[collapsible=icon]:hidden">
                                        <Command className="h-3 w-3" />K
                                    </kbd>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </div>
                </SidebarContent>
                <SidebarFooter>
                    {user && (
                        <NavUser
                            user={user}
                            isSuperAdmin={checkIsSuperAdmin()}
                            onMessagesClick={openMessages}
                            onNotificationsClick={openNotifications}
                        />
                    )}
                </SidebarFooter>
            </Sidebar>

            {/* Global Search Dialog */}
            <GlobalSearch isOpen={searchOpen} onClose={() => setSearchOpen(false)} />
        </>
    )
}
