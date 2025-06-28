import * as React from "react"
import { ShieldIcon } from "lucide-react"

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
} from "@/components/ui/sidebar"
import { usePermissions } from "@/contexts/PermissionsContext"
import { navigationConfig } from "@/config/navigation"

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
    user?: any;
}

export function AppSidebar({user, ...props}: AppSidebarProps) {
    const { canPerformAction, isSuperAdmin: checkIsSuperAdmin, isAdmin } = usePermissions();
    
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
    
    return (
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
            </SidebarContent>
            <SidebarFooter>
                {user && (
                    <NavUser user={user}/>
                )}
            </SidebarFooter>
        </Sidebar>
    )
}
