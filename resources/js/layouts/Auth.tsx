import React, { ReactNode } from 'react';
import { ThemeToggleIcon } from '@/components/ThemeToggleIcon';

interface AuthLayoutProps {
  children: ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="relative flex flex-col p-6 md:p-10 overflow-hidden">
        {/* Spider web texture overlay - fades in from bottom right */}
        <div className="absolute inset-0 pointer-events-none opacity-[0.06] dark:opacity-[0.08]">
          <svg
            className="absolute -bottom-1/4 -right-1/4 w-[150%] h-[150%]"
            viewBox="0 0 800 800"
            xmlns="http://www.w3.org/2000/svg"
          >
            <defs>
              <radialGradient id="webFade" cx="100%" cy="100%" r="100%">
                <stop offset="0%" stopColor="currentColor" stopOpacity="1" />
                <stop offset="60%" stopColor="currentColor" stopOpacity="0.4" />
                <stop offset="100%" stopColor="currentColor" stopOpacity="0" />
              </radialGradient>
            </defs>

            {/* Concentric web circles */}
            <circle cx="800" cy="800" r="100" fill="none" stroke="url(#webFade)" strokeWidth="2.5" />
            <circle cx="800" cy="800" r="200" fill="none" stroke="url(#webFade)" strokeWidth="2.5" />
            <circle cx="800" cy="800" r="300" fill="none" stroke="url(#webFade)" strokeWidth="2" />
            <circle cx="800" cy="800" r="400" fill="none" stroke="url(#webFade)" strokeWidth="1.8" />
            <circle cx="800" cy="800" r="500" fill="none" stroke="url(#webFade)" strokeWidth="1.5" />
            <circle cx="800" cy="800" r="600" fill="none" stroke="url(#webFade)" strokeWidth="1.3" />
            <circle cx="800" cy="800" r="700" fill="none" stroke="url(#webFade)" strokeWidth="1" />

            {/* Radial web lines */}
            <line x1="800" y1="800" x2="800" y2="100" stroke="url(#webFade)" strokeWidth="2.5" />
            <line x1="800" y1="800" x2="100" y2="800" stroke="url(#webFade)" strokeWidth="2.5" />
            <line x1="800" y1="800" x2="600" y2="200" stroke="url(#webFade)" strokeWidth="2" />
            <line x1="800" y1="800" x2="200" y2="600" stroke="url(#webFade)" strokeWidth="2" />
            <line x1="800" y1="800" x2="400" y2="300" stroke="url(#webFade)" strokeWidth="1.8" />
            <line x1="800" y1="800" x2="300" y2="400" stroke="url(#webFade)" strokeWidth="1.8" />
            <line x1="800" y1="800" x2="500" y2="200" stroke="url(#webFade)" strokeWidth="1.5" />
            <line x1="800" y1="800" x2="200" y2="500" stroke="url(#webFade)" strokeWidth="1.5" />
            <line x1="800" y1="800" x2="650" y2="300" stroke="url(#webFade)" strokeWidth="1.5" />
            <line x1="800" y1="800" x2="300" y2="650" stroke="url(#webFade)" strokeWidth="1.5" />
            <line x1="800" y1="800" x2="450" y2="200" stroke="url(#webFade)" strokeWidth="1.3" />
            <line x1="800" y1="800" x2="200" y2="450" stroke="url(#webFade)" strokeWidth="1.3" />
          </svg>
        </div>

        {/* Header section for the logo */}
        <div className="relative z-10 flex justify-between items-center mb-4 md:mb-8">
          <a href="/" className="flex items-center gap-2 font-medium">
            <img src="/images/mw-coat.svg" alt="Logo" className="h-8 w-auto" />
            {/* Optionally, add text next to logo if desired */}
            {/* <span className="text-lg font-semibold">Acme Inc.</span> */}
          </a>
          <ThemeToggleIcon />
        </div>
        {/* Main content area */}
        <div className="relative z-10 flex flex-1 items-center justify-center">
          <div className="w-full max-w-sm">
            {children} {/* Page-specific form will be rendered here */}
          </div>
        </div>
      </div>
      <div className="relative hidden bg-muted lg:block">
        <img
          src="https://msmedb.trade.gov.mw/dist/media/images/backgrounds/bg-3.jpg"
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
      </div>
    </div>
  );
}
