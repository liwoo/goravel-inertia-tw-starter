import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}


// snakefiy json object keys
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