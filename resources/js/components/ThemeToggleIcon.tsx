import React from 'react';
import { MoonIcon } from 'lucide-react';

export function ThemeToggleIcon() {
  return (
    <button
      className="rounded-md p-2 hover:bg-sidebar-accent transition-colors"
      aria-label="Theme toggle"
    >
      <MoonIcon className="h-5 w-5" />
    </button>
  );
}
