import React from 'react';
import { createRoot } from 'react-dom/client';
// @ts-ignore
import { createInertiaApp, PageProps, AppType } from '@inertiajs/react';
import '../css/app.css';
import { ThemeProvider } from './context/ThemeContext';
import { Toaster } from './components/ui/sonner';

const appName = import.meta.env.VITE_APP_NAME || 'Blog';

createInertiaApp({
  title: (title: string) => `${title} - ${appName}`,
  resolve: (name: string) => {
    const pages = import.meta.glob('./pages/**/*.tsx', { eager: true });
    const pageModule = pages[`./pages/${name}.tsx`] as { default: React.ComponentType<any> } | undefined;

    if (!pageModule) {
      console.error(`Page component not found: ${name}. Available pages:`, Object.keys(pages));
      // You might want to throw an error or return a specific 404 component here
      // For now, throwing an error to make it clear if a page is missing.
      throw new Error(`Page component "${name}" not found.`);
    }

    return pageModule.default;
  },
  setup({ el, App, props }: { el: HTMLElement; App: AppType<PageProps>; props: PageProps }) {
    const root = createRoot(el);
    root.render(
      <ThemeProvider>
        <App {...props} />
        <Toaster richColors position="top-right" />
      </ThemeProvider>
    );
  },
  progress: {
    color: '#4B5563',
  },
});
