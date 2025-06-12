import React from 'react';
import { MoonIcon, SunIcon } from 'lucide-react';
import { useTheme } from '../context/ThemeContext';
import {
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

export function SidebarHeaderWithTheme() {
  const { theme, toggleTheme } = useTheme();
  
  return (
    <SidebarMenu>
      <div className="flex items-center justify-between w-full px-3 py-2">
        <SidebarMenuItem className="!p-0">
          <a href="/dashboard" className="flex items-center gap-2">
            <img src="/placeholder.svg" alt="Logo" className="h-8 w-auto" />
          </a>
        </SidebarMenuItem>
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
    </SidebarMenu>
  );
}
