import {
    BarChartIcon,
    BookIcon,
    CameraIcon,
    ClipboardListIcon,
    DatabaseIcon,
    FileCodeIcon,
    FileIcon,
    FileTextIcon,
    HelpCircleIcon,
    LayoutDashboardIcon,
    SettingsIcon,
    ShieldIcon,
    UsersIcon,
} from "lucide-react"

// Navigation item types
export interface BaseNavItem {
    title: string;
    url: string;
    icon: any;
    requiredService?: string;
    requiredAction?: "read" | "write" | "delete" | "manage" | "create" | "update" | "export" | "bulk_update" | "bulk_delete";
    requireSuperAdmin?: boolean;
}

export type NavItem = BaseNavItem;

interface NavItemWithChildren extends BaseNavItem {
    isActive?: boolean;
    items?: {
        title: string;
        url: string;
    }[];
}

interface DocumentItem {
    name: string;
    url: string;
    icon: any;
    requiredService?: string;
    requiredAction?: "read" | "write" | "delete" | "manage" | "create" | "update" | "export" | "bulk_update" | "bulk_delete";
}

export interface NavigationConfig {
    navMain: NavItem[];
    navClouds: NavItemWithChildren[];
    navSecondary: NavItem[];
    documents: DocumentItem[];
}

// Navigation items with permission requirements
export const navigationConfig: NavigationConfig = {
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
            title: "Analysis",
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
            title: "Roles & Permissions",
            url: "/admin/permissions",
            icon: ShieldIcon,
            requireSuperAdmin: true,
        },
        {
            title: "Users",
            url: "/admin/users",
            icon: UsersIcon,
            requireSuperAdmin: true,
        },
        {
            title: "Get Help",
            url: "#",
            icon: HelpCircleIcon,
            // Help is always accessible
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