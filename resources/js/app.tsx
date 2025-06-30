import React from 'react';
import { createRoot } from 'react-dom/client';
// @ts-ignore
import { createInertiaApp, PageProps, AppType } from '@inertiajs/react';
import '../css/app.css';
import { Toaster } from 'sonner';
import { PermissionsProvider } from '@/contexts/PermissionsContext';
import { ThemeProvider } from '@/context/ThemeContext';
import '@/lib/axios'; // Configure axios defaults 

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

    // Wrap each page component with ThemeProvider and PermissionsProvider
    const PageComponent = pageModule.default;
    const WrappedPage = (props: any) => (
      <ThemeProvider>
        <PermissionsProvider>
          <PageComponent {...props} />
        </PermissionsProvider>
      </ThemeProvider>
    );

    return WrappedPage;
  },
  setup({ el, App, props }: { el: HTMLElement; App: AppType<PageProps>; props: PageProps }) {
    const root = createRoot(el);
    
    root.render(
      <>
        <App {...props} />
        <Toaster 
          position="top-right"
          expand={true}
          richColors
          closeButton
        />
      </>
    ); 
  },
  progress: {
    color: '#4B5563',
  },
});
