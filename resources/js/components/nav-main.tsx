import { MailIcon, PlusCircleIcon, Command, type LucideIcon } from "lucide-react"
// @ts-ignore
import { Link, router, usePage } from '@inertiajs/react'
import { useEffect } from 'react'
import { usePermissions } from "@/contexts/PermissionsContext"

import { Button } from "@/components/ui/button"
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { cn } from "@/lib/utils"

export function NavMain({
  items,
}: {
  items: {
    title: string
    url: string
    icon?: LucideIcon
    variant?: "default" | "primary"
    action?: "openMySmeModal"
  }[]
}) {
  const { url } = usePage();
  const { canPerformAction } = usePermissions();

  // Check if a menu item is active based on current URL
  const isActive = (itemUrl: string) => {
    // Handle exact match for dashboard
    if (itemUrl === '/dashboard') {
      return url === '/dashboard';
    }
    // For other routes, check if current URL starts with the item URL
    // This handles /admin/smes, /admin/smes/1, /admin/smes/create, etc.
    return url.startsWith(itemUrl);
  };

  const primaryCta = items.find(item => item.variant === 'primary');

  // Handle primary CTA click - either dispatch action or navigate
  const handlePrimaryCtaClick = () => {
    if (primaryCta?.action) {
      // Dispatch custom event for actions
      window.dispatchEvent(new CustomEvent('nav-action', { detail: { action: primaryCta.action } }));
    } else if (primaryCta?.url) {
      router.visit(primaryCta.url);
    }
  };

  // Check if user can create any entity
  const canCreateAnything =
    canPerformAction('smes', 'create') ||
    canPerformAction('bdsps', 'create') ||
    canPerformAction('events', 'create') ||
    canPerformAction('procurement_notices', 'create');

  // Keyboard shortcuts for navigation items (Cmd/Ctrl + 1-9)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.metaKey || e.ctrlKey) {
        const key = parseInt(e.key);
        if (key >= 1 && key <= 9 && items[key - 1]) {
          e.preventDefault();
          router.visit(items[key - 1].url);
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [items]);

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        {primaryCta != null ? <SidebarMenu>
          <SidebarMenuItem className="flex items-center gap-2">
            <SidebarMenuButton
              tooltip={primaryCta.title}
              className="min-w-16 bg-primary text-primary-foreground duration-200 ease-linear hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
              onClick={handlePrimaryCtaClick}
            >
              <PlusCircleIcon />
              <span>{primaryCta.title}</span>
            </SidebarMenuButton>
            <Button
              size="icon"
              className="h-9 w-9 shrink-0 group-data-[collapsible=icon]:opacity-0"
              variant="outline"
            >
              <MailIcon />
              <span className="sr-only">Inbox</span>
            </Button>
          </SidebarMenuItem>
        </SidebarMenu> : canCreateAnything && (
          <SidebarMenu>
            <SidebarMenuItem className="flex items-center gap-2">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <SidebarMenuButton
                    tooltip="Quick Create"
                    className="min-w-8 bg-primary text-primary-foreground duration-200 ease-linear hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
                  >
                    <PlusCircleIcon />
                    <span>Quick Create</span>
                  </SidebarMenuButton>
                </DropdownMenuTrigger>
                <DropdownMenuContent side="right" align="start" className="w-48">
                  <DropdownMenuLabel>Create New</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  {canPerformAction('smes', 'create') && (
                    <DropdownMenuItem onClick={() => router.visit('/smes/create')}>
                      MSME
                    </DropdownMenuItem>
                  )}
                  {canPerformAction('bdsps', 'create') && (
                    <DropdownMenuItem onClick={() => router.visit('/bdsps/create')}>
                      BDSP
                    </DropdownMenuItem>
                  )}
                  {canPerformAction('events', 'create') && (
                    <DropdownMenuItem onClick={() => router.visit('/events/create')}>
                      Event
                    </DropdownMenuItem>
                  )}
                  {canPerformAction('procurement_notices', 'create') && (
                    <DropdownMenuItem onClick={() => router.visit('/procurements/create')}>
                      Procurement
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
              <Button
                size="icon"
                className="h-9 w-9 shrink-0 group-data-[collapsible=icon]:opacity-0"
                variant="outline"
              >
                <MailIcon />
                <span className="sr-only">Inbox</span>
              </Button>
            </SidebarMenuItem>
          </SidebarMenu>
        )}
        <SidebarMenu>
          {items.filter(item => item.variant !== 'primary').map((item, index) => {
            const active = isActive(item.url);
            return (
              <SidebarMenuItem key={item.title}>
                <Link href={item.url}>
                  <SidebarMenuButton
                    tooltip={item.title}
                    isActive={active}
                  >
                    {item.icon && <item.icon />}
                    <span>{item.title}</span>
                    {index < 9 && (
                      <kbd className="ml-auto pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 group-data-[collapsible=icon]:hidden">
                        <Command className="h-3 w-3" />
                        {index + 1}
                      </kbd>
                    )}
                  </SidebarMenuButton>
                </Link>
              </SidebarMenuItem>
            );
          })}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
