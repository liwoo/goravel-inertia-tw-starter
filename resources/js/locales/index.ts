import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import common from './en/common.json';
import crud from './en/crud.json';
import nav from './en/nav.json';
import auth from './en/auth.json';
import exportNs from './en/export.json';
import dashboard from './en/dashboard.json';
import users from './en/users.json';
import books from './en/books.json';
import authors from './en/authors.json';
import settings from './en/settings.json';

i18n.use(initReactI18next).init({
  lng: 'en',
  fallbackLng: 'en',
  ns: ['common', 'crud', 'nav', 'auth', 'export', 'dashboard', 'users', 'books', 'authors', 'settings'],
  defaultNS: 'common',
  resources: {
    en: {
      common,
      crud,
      nav,
      auth,
      export: exportNs,
      dashboard,
      users,
      books,
      authors,
      settings,
    },
  },
  interpolation: {
    escapeValue: false, // React already escapes
  },
  react: {
    useSuspense: false,
  },
});

export default i18n;
