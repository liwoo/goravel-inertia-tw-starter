import {
    BookIcon,
    CameraIcon,
    CogIcon,
    FileCodeIcon,
    FileTextIcon,
    HelpCircleIcon,
    LandmarkIcon,
    LayoutDashboardIcon,
    NotebookTabsIcon,
    PenToolIcon,
    ShieldIcon,
    UsersIcon,
} from "lucide-react"

// Navigation item types
export interface BaseNavItem {
    title: string; // i18n key within 'nav' namespace
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
        title: string; // i18n key within 'nav' namespace
        url: string;
    }[];
}

interface DocumentItem {
    name: string; // i18n key within 'nav' namespace
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
// title/name values are i18n keys resolved in nav components via useTranslation('nav')
export const navigationConfig: NavigationConfig = {
    navMain: [
        {
            title: "main.dashboard",
            url: "/dashboard",
            icon: LayoutDashboardIcon,
            // Dashboard is always accessible to authenticated users
        },
        {
            title: "main.books",
            url: "/admin/books",
            icon: BookIcon,
            requiredService: "books",
            requiredAction: "read" as const,
        },
        {
            title: "main.authors",
            url: "/admin/authors",
            icon: PenToolIcon,
            requiredService: "authors",
            requiredAction: "read" as const,
        },
        {
            title: "main.lenders",
            url: "/admin/lenders",
            icon: LandmarkIcon,
            requiredService: "lenders",
            requiredAction: "read" as const,
        },
    ],
    navClouds: [
        {
            title: "clouds.capture",
            icon: CameraIcon,
            isActive: true,
            url: "#",
            items: [
                {
                    title: "clouds.activeProposals",
                    url: "#",
                },
                {
                    title: "clouds.archived",
                    url: "#",
                },
            ],
        },
        {
            title: "clouds.proposal",
            icon: FileTextIcon,
            url: "#",
            items: [
                {
                    title: "clouds.activeProposals",
                    url: "#",
                },
                {
                    title: "clouds.archived",
                    url: "#",
                },
            ],
        },
        {
            title: "clouds.prompts",
            icon: FileCodeIcon,
            url: "#",
            items: [
                {
                    title: "clouds.activeProposals",
                    url: "#",
                },
                {
                    title: "clouds.archived",
                    url: "#",
                },
            ],
        },
    ],
    navSecondary: [
        {
            title: "secondary.rolesPermissions",
            url: "/admin/permissions",
            icon: ShieldIcon,
            requireSuperAdmin: true,
        },
        {
            title: "secondary.users",
            url: "/admin/users",
            icon: UsersIcon,
            requireSuperAdmin: true,
        },
        {
            title: "secondary.getHelp",
            url: "#",
            icon: HelpCircleIcon,
            // Help is always accessible
        },
    ],
    documents: [
        {
            name: "documents.applications",
            url: "/admin/applications",
            icon: NotebookTabsIcon,
            requiredService: "applications",
            requiredAction: "read" as const,
        },
        {
            name: "documents.configurations",
            url: "/admin/configs",
            icon: CogIcon,
            requiredService: "config",
            requiredAction: "read" as const,
        },
    ],

}
