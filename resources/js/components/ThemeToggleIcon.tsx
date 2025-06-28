import React from 'react';
import { MoonIcon, SunIcon } from 'lucide-react';
import { useTheme } from '@/context/ThemeContext';

export function ThemeToggleIcon() {
  const { theme, toggleTheme } = useTheme();
  
  return (
    <button
      onClick={toggleTheme}
      className="rounded-md p-2 hover:bg-sidebar-accent transition-colors"
      aria-label="Toggle theme"
    >
      {theme === 'dark' ? (
        <SunIcon className="h-5 w-5" />
      ) : (
        <MoonIcon className="h-5 w-5" />
      )}
    </button>
  );
}
