declare namespace App {
  interface PageProps {
    errors: Record<string, string>;
    [key: string]: any;
  }
}

declare module '@inertiajs/react' {
  import { ComponentType } from 'react';

  export interface InertiaPageProps extends App.PageProps {
    [key: string]: any;
  }

  export function usePage<T = {}>(): {
    props: InertiaPageProps & T;
    url: string;
    component: string;
    version: string;
  };
}
