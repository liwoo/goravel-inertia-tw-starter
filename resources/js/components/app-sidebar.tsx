import * as React from "react"
import {
    BarChartIcon,
    BookIcon,
    CameraIcon,
    ClipboardListIcon,
    DatabaseIcon,
    FileCodeIcon,
    FileIcon,
    FileTextIcon,
    FolderIcon,
    HelpCircleIcon,
    LayoutDashboardIcon,
    ListIcon,
    MoonIcon,
    SearchIcon,
    SettingsIcon,
    ShieldIcon,
    SunIcon,
    UsersIcon,
} from "lucide-react"

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

// Navigation items with permission requirements
const navigationConfig = {
    navMain: [
        {
            title: "Dashboard",
            url: "/dashboard",
            icon: LayoutDashboardIcon,
            // Dashboard is always accessible to authenticated users
        },
        {
            title: "Books",
            url: "/admin/books",
            icon: BookIcon,
            requiredService: "books",
            requiredAction: "read" as const,
        },
        {
            title: "Users",
            url: "/admin/users",
            icon: UsersIcon,
            requiredService: "users",
            requiredAction: "read" as const,
        },
        {
            title: "Analytics",
            url: "#",
            icon: BarChartIcon,
            requiredService: "reports",
            requiredAction: "read" as const,
        },
    ],
    navClouds: [
        {
            title: "Capture",
            icon: CameraIcon,
            isActive: true,
            url: "#",
            items: [
                {
                    title: "Active Proposals",
                    url: "#",
                },
                {
                    title: "Archived",
                    url: "#",
                },
            ],
        },
        {
            title: "Proposal",
            icon: FileTextIcon,
            url: "#",
            items: [
                {
                    title: "Active Proposals",
                    url: "#",
                },
                {
                    title: "Archived",
                    url: "#",
                },
            ],
        },
        {
            title: "Prompts",
            icon: FileCodeIcon,
            url: "#",
            items: [
                {
                    title: "Active Proposals",
                    url: "#",
                },
                {
                    title: "Archived",
                    url: "#",
                },
            ],
        },
    ],
    navSecondary: [
        {
            title: "Permissions",
            url: "/admin/permissions",
            icon: ShieldIcon,
            requireSuperAdmin: true,
        },
        {
            title: "Settings",
            url: "/settings",
            icon: SettingsIcon,
            // Settings is always accessible to authenticated users
        },
        {
            title: "Get Help",
            url: "#",
            icon: HelpCircleIcon,
            // Help is always accessible
        },
        {
            title: "Search",
            url: "#",
            icon: SearchIcon,
            // Search is always accessible
        },
    ],
    documents: [
        {
            name: "Data Library",
            url: "#",
            icon: DatabaseIcon,
            requiredService: "reports",
            requiredAction: "read" as const,
        },
        {
            name: "Reports",
            url: "#",
            icon: ClipboardListIcon,
            requiredService: "reports", 
            requiredAction: "read" as const,
        },
        {
            name: "Word Assistant",
            url: "#",
            icon: FileIcon,
            // Always accessible
        },
    ],
}

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
    user?: any;
}

export function AppSidebar({user, ...props}: AppSidebarProps) {
    const { canPerformAction, isSuperAdmin: checkIsSuperAdmin, isAdmin } = usePermissions();
    
    // Filter navigation items based on permissions
    const navigationItems = React.useMemo(() => {
        // Filter main navigation
        const filteredNavMain = navigationConfig.navMain.filter(item => {
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
