import axios from 'axios';

// Configure axios defaults
axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
axios.defaults.headers.common['Accept'] = 'application/json';

// Get CSRF token from meta tag if it exists
const token = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
if (token) {
  axios.defaults.headers.common['X-CSRF-TOKEN'] = token;
}

// Add interceptor to handle authentication errors
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Don't redirect if we're already on auth pages or making auth requests
      const currentPath = window.location.pathname;
      const isAuthPage = currentPath === '/login' || currentPath === '/' || currentPath === '/una';
      const isAuthRequest = error.config?.url?.includes('/login') || error.config?.url?.includes('/auth');

      // Only redirect to /una if we're not on auth pages and not making auth requests
      if (!isAuthPage && !isAuthRequest) {
        window.location.href = '/una';
      }
    } else if (error.response?.status === 409) {
      // Handle Inertia redirect conflicts
      const location = error.response.headers['x-inertia-location'];
      if (location) {
        window.location.href = location;
      }
    }
    return Promise.reject(error);
  }
);

export default axios;