// i18n TypeScript augmentation
// Relaxed typing to support cross-namespace keys and dynamic key variables.
// Strict key autocomplete can be re-enabled after all migrations are complete.
import 'i18next';

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'common';
  }
}
