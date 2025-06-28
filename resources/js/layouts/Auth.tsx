import React, { ReactNode } from 'react';
import { ThemeToggleIcon } from '@/components/ThemeToggleIcon';

interface AuthLayoutProps {
  children: ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col p-6 md:p-10">
        {/* Header section for the logo */}
        <div className="flex justify-between items-center mb-4 md:mb-8"> 
          <a href="/" className="flex items-center gap-2 font-medium">
            <img src="/placeholder.svg" alt="Logo" className="h-8 w-auto" /> 
            {/* Optionally, add text next to logo if desired */}
            {/* <span className="text-lg font-semibold">Acme Inc.</span> */}
          </a>
          <ThemeToggleIcon />
        </div>
        {/* Main content area */}
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            {children} {/* Page-specific form will be rendered here */}
          </div>
        </div>
      </div>
      <div className="relative hidden bg-muted lg:block">
        <img
          src="/placeholder.svg" 
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
      </div>
    </div>
  );
}
