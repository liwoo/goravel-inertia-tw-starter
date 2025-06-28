import React from 'react';
import {
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

export function SidebarHeaderWithTheme() {
  return (
    <SidebarMenu>
      <div className="flex items-center justify-between w-full px-3 py-2">
        <SidebarMenuItem className="!p-0">
          <a href="/dashboard" className="flex items-center gap-2">
            <img src="/placeholder.svg" alt="Logo" className="h-8 w-auto" />
          </a>
        </SidebarMenuItem>
      </div>
    </SidebarMenu>
  );
}
