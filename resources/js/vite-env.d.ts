/// <reference types="vite/client" />

interface ImportMeta {
  glob: (pattern: string) => Record<string, () => Promise<any>>;
}
