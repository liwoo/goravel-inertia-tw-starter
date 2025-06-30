import * as React from "react"
import { ShieldIcon, Search, Command } from "lucide-react"

import {NavDocuments} from "@/components/nav-documents"
import {NavMain} from "@/components/nav-main"
import {NavSecondary} from "@/components/nav-secondary"
import {NavUser} from "@/components/nav-user"
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

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
    user?: any;
}

export function AppSidebar({user, ...props}: AppSidebarProps) {
    const { canPerformAction, isSuperAdmin: checkIsSuperAdmin, isAdmin } = usePermissions();
    const [searchOpen, setSearchOpen] = React.useState(false);
    
    // Filter navigation items based on permissions
    const navigationItems = React.useMemo(() => {
        // Filter main navigation
        const filteredNavMain = navigationConfig.navMain.filter(item => {
            // Check super admin requirement
            if (item.requireSuperAdmin) {
                return checkIsSuperAdmin();
            }
            
            // If no permission requirement, show the item
            if (!item.requiredService && !item.requiredAction) {
                return true;
            }
            
            // Check service permission
            if (item.requiredService && item.requiredAction) {
                return canPerformAction(item.requiredService, item.requiredAction);
            }
            
            return true;
        });
        
        // Filter secondary navigation
        const filteredNavSecondary = navigationConfig.navSecondary.filter(item => {
            // Check super admin requirement
            if (item.requireSuperAdmin) {
                return checkIsSuperAdmin();
            }
            
            // Check service permission
            if (item.requiredService && item.requiredAction) {
                return canPerformAction(item.requiredService, item.requiredAction);
            }
            
            return true;
        });
        
        // Filter documents
        const filteredDocuments = navigationConfig.documents.filter(item => {
            // Check service permission
            if (item.requiredService && item.requiredAction) {
                return canPerformAction(item.requiredService, item.requiredAction);
            }
            
            return true;
        });
        
        return {
            navMain: filteredNavMain,
            navSecondary: filteredNavSecondary,
            documents: filteredDocuments,
        };
    }, [canPerformAction, checkIsSuperAdmin]);

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
            <Sidebar collapsible="offcanvas" {...props}>
                <SidebarHeader>
                    <div className="flex items-center justify-between px-3 py-2">
                        <a href="/dashboard" className="flex items-center gap-2">
                            <img src="/placeholder.svg" alt="Logo" className="h-8 w-auto" />
                        </a>
                    </div>
                    {/* Super Admin Badge */}
                    {checkIsSuperAdmin() && (
                        <div className="px-3 pb-2">
                            <div className="flex items-center gap-2 rounded-md bg-red-100 dark:bg-red-900/20 px-2 py-1 text-xs">
                                <ShieldIcon className="h-3 w-3 text-red-600 dark:text-red-400" />
                                <span className="text-red-700 dark:text-red-300 font-medium">Super Admin</span>
                            </div>
                        </div>
                    )}
                </SidebarHeader>
                <SidebarContent>
                    <NavMain items={navigationItems.navMain}/>
                    <NavDocuments items={navigationItems.documents}/>
                    <NavSecondary items={navigationItems.navSecondary} className="mt-auto"/>
                    
                    {/* Hardcoded Search Option */}
                    <div className="mt-2 px-3 pb-3">
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton onClick={() => setSearchOpen(true)}>
                                    <Search className="h-4 w-4" />
                                    <span>Search</span>
                                    <kbd className="ml-auto pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100">
                                        <Command className="h-3 w-3" />K
                                    </kbd>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </div>
                </SidebarContent>
                <SidebarFooter>
                    {user && (
                        <NavUser user={user}/>
                    )}
                </SidebarFooter>
            </Sidebar>
            
            {/* Global Search Dialog */}
            <GlobalSearch isOpen={searchOpen} onClose={() => setSearchOpen(false)} />
        </>
    )
}
