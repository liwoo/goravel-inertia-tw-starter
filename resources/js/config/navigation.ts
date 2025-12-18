import {
    CogIcon,
    HelpCircleIcon,
    LayoutDashboardIcon,
    ShieldIcon,
    UsersIcon,
    Book,
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
            icon: Book,
            requiredService: "books",
            requiredAction: "read" as const,
        },
    ],
    navClouds: [

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
    ],
}