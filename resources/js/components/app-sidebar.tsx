import * as React from "react"
import {
    BarChartIcon,
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
    SunIcon,
    UsersIcon,
} from "lucide-react"

import {NavDocuments} from "@/components/nav-documents"
import {NavMain} from "@/components/nav-main"
import {NavSecondary} from "@/components/nav-secondary"
import {NavUser} from "@/components/nav-user"
import {useTheme} from "@/context/ThemeContext"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuItem,
} from "@/components/ui/sidebar"

const data = {
    user: {
        name: "shadcn",
        email: "m@example.com",
        avatar: "/avatars/shadcn.jpg",
    },
    navMain: [
        {
            title: "Dashboard",
            url: "/dashboard",
            icon: LayoutDashboardIcon,
        },
        {
            title: "Lifecycle",
            url: "#",
            icon: ListIcon,
        },
        {
            title: "Analytics",
            url: "#",
            icon: BarChartIcon,
        },
        {
            title: "Projects",
            url: "#",
            icon: FolderIcon,
        },
        {
            title: "Team",
            url: "#",
            icon: UsersIcon,
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
            title: "Settings",
            url: "/settings",
            icon: SettingsIcon,
        },
        {
            title: "Get Help",
            url: "#",
            icon: HelpCircleIcon,
        },
        {
            title: "Search",
            url: "#",
            icon: SearchIcon,
        },
    ],
    documents: [
        {
            name: "Data Library",
            url: "#",
            icon: DatabaseIcon,
        },
        {
            name: "Reports",
            url: "#",
            icon: ClipboardListIcon,
        },
        {
            name: "Word Assistant",
            url: "#",
            icon: FileIcon,
        },
    ],
}

export function AppSidebar({user, ...props}: React.ComponentProps<typeof Sidebar>) {
    const { theme, toggleTheme } = useTheme();
    
    return (
        <Sidebar collapsible="offcanvas" {...props}>
            <SidebarHeader>
                <div className="flex items-center justify-between px-3 py-2">
                    <a href="/dashboard" className="flex items-center gap-2">
                        <img src="/placeholder.svg" alt="Logo" className="h-8 w-auto" />
                    </a>
                    <button
                        onClick={toggleTheme}
                        className="rounded-md p-2 hover:bg-sidebar-accent transition-colors"
                        aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                    >
                        {theme === 'dark' ? (
                            <SunIcon className="h-5 w-5" />
                        ) : (
                            <MoonIcon className="h-5 w-5" />
                        )}
                    </button>
                </div>
            </SidebarHeader>
            <SidebarContent>
                <NavMain items={data.navMain}/>
                <NavDocuments items={data.documents}/>
                <NavSecondary items={data.navSecondary} className="mt-auto"/>
            </SidebarContent>
            <SidebarFooter>
                {user && (
                    <NavUser user={user}/>
                )}
            </SidebarFooter>
        </Sidebar>
    )
}
