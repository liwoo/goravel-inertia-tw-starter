import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}


// snakefiy json object keys (convert camelCase to snake_case)
export function snakefiyKeys(obj: any): any {
  if (Array.isArray(obj)) {
    return obj.map(snakefiyKeys);
  } else if (obj !== null && typeof obj === 'object') {
    return Object.fromEntries(
      Object.entries(obj).map(([key, value]) => [
        key.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`),
        snakefiyKeys(value)
      ])
    );
  }
  return obj;
}

// camelify json object keys (convert snake_case to camelCase)
export function camelifyKeys(obj: any): any {
  if (Array.isArray(obj)) {
    return obj.map(camelifyKeys);
  } else if (obj !== null && typeof obj === 'object') {
    return Object.fromEntries(
      Object.entries(obj).map(([key, value]) => [
        key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase()),
        camelifyKeys(value)
      ])
    );
  }
  return obj;
}