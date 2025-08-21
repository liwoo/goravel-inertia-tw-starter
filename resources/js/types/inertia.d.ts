// Global namespace for app-specific types
declare namespace App {
  interface PageProps {
    errors: Record<string, string>;
    [key: string]: any;
  }
}
