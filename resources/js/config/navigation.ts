import {
    BarChartIcon,
    BookIcon, BuildingIcon, Calendar1Icon,
    CameraIcon,
    ClipboardListIcon,
    CogIcon,
    DatabaseIcon,
    FileCodeIcon,
    FileIcon,
    FileTextIcon,
    FolderIcon,
    HelpCircleIcon,
    LayoutDashboardIcon, NotebookTabsIcon, NotepadText, PercentSquareIcon, PersonStandingIcon,
    SettingsIcon,
    ShieldIcon, User2Icon,
    UsersIcon,
} from "lucide-react"

// Navigation item types
export interface BaseNavItem {
    title: string;
    url: string;
    icon: any;
    requiredService?: string;
    requiredAction?: "read" | "write" | "delete" | "manage" | "create" | "update" | "export" | "bulk_update" | "bulk_delete";
    requiredRole?: string;
    requireSuperAdmin?: boolean;
    variant?: "default" | "primary";
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
    requiredRole?: string;
}

export interface NavigationConfig {
    navMain: NavItem[];
    navClouds: NavItemWithChildren[];
    navSecondary: NavItem[];
    documents: DocumentItem[];
    navSme: NavItem[];
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
            title: "SMEs",
            url: "/admin/smes",
            icon: User2Icon,
            requiredService: "smes",
            requiredAction: "read" as const,
        },
        {
            title: "BDSPs",
            url: "/admin/bdsps",
            icon: BuildingIcon,
            requiredService: "bdsps",
            requiredAction: "read" as const,
        },

        {
            title: "Events",
            url: "/admin/events",
            icon: Calendar1Icon,
            requiredService: "events",
            requiredAction: "read" as const,
        },

        {
            title: "Procurement",
            url: "/admin/procurement-notices",
            icon: ClipboardListIcon,
            requiredService: "procurement_notices",
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
            name: "Configurations",
            url: "/admin/configs",
            icon: CogIcon,
            requiredService: "config",
            requiredAction: "read" as const,
        },
        {
            name: "Applications",
            url: "/admin/applications",
            icon: NotebookTabsIcon,
            requiredService: "applications",
            requiredAction: "read" as const,
        },
    ],

    navSme: [
        {
            title: "My MSME",
            url: "/dashboard",
            icon: LayoutDashboardIcon,
            requiredRole: "sme-user",
            variant: "primary",
        },
        {
            title: "Portal",
            url: "/portal",
            icon: LayoutDashboardIcon,
            requiredRole: "sme-user",
        },
        {
            title: "Directory",
            url: "/directory",
            icon: FolderIcon,
            requiredRole: "sme-user",
        },
    ],
}