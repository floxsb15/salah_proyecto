const unsafeMethods = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

export default defineNuxtPlugin((nuxtApp) => {
  const originalFetch = globalThis.$fetch;
  const secureFetch = originalFetch.create({
    credentials: 'same-origin',
    onRequest({ request, options }) {
      const method = String(options.method || 'GET').toUpperCase();
      const path = typeof request === 'string' ? request : request.toString();
      if (!path.startsWith('/api/') || !unsafeMethods.has(method) || path.endsWith('/login')) {
        return;
      }
      const csrfToken = readCookie('csrf_token');
      if (!csrfToken) {
        return;
      }
      const headers = new Headers(options.headers);
      headers.set('X-CSRF-Token', csrfToken);
      options.headers = headers;
    },
    async onResponseError({ request, response }) {
      const path = typeof request === 'string' ? request : request.toString();
      if (response.status === 401 && !path.endsWith('/login')) {
        localStorage.removeItem('user');
        await nuxtApp.runWithContext(() => navigateTo('/login'));
      }
    },
  });

  globalThis.$fetch = secureFetch as typeof globalThis.$fetch;
});

function readCookie(name: string): string {
  const prefix = `${encodeURIComponent(name)}=`;
  const entry = document.cookie.split('; ').find((value) => value.startsWith(prefix));
  return entry ? decodeURIComponent(entry.slice(prefix.length)) : '';
}
