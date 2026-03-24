const CACHE_NAME = 'carcare-pwa-v3';
const LEGACY_CACHE_PREFIXES = ['carcare-', 'carcare-pwa-'];
const OFFLINE_URL = '/offline.html';
const ASSETS = [
  '/',
  '/index.html',
  '/manifest.webmanifest',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png',
  '/icons/maskable-icon-192x192.png',
  '/icons/maskable-icon-512x512.png',
  '/favicon.ico',
  OFFLINE_URL
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys
          .filter((key) => key !== CACHE_NAME)
          .filter((key) => LEGACY_CACHE_PREFIXES.some((prefix) => key.startsWith(prefix)))
          .map((key) => caches.delete(key))
      )
    )
  );
  event.waitUntil(self.clients.claim());
});

// Protected routes that contain user-specific data and should not be cached
const PROTECTED_ROUTES = ['/profile', '/cars', '/fuel', '/maintenance', '/fines', '/reports', '/users'];

function isProtectedRoute(pathname) {
  return PROTECTED_ROUTES.some(route => pathname === route || pathname.startsWith(route + '/'));
}

self.addEventListener('fetch', (event) => {
  const request = event.request;
  const requestUrl = new URL(request.url);
  const isApiRequest = requestUrl.pathname.startsWith('/api/');
  const isGet = request.method === 'GET';
  const isNavigation = request.mode === 'navigate';
  const isProtected = isProtectedRoute(requestUrl.pathname);

  if (isApiRequest || !isGet) {
    event.respondWith(fetch(request));
    return;
  }

  // Protected routes should never be cached - always fetch fresh
  if (isProtected) {
    event.respondWith(
      fetch(request)
        .catch(async () => {
          const cached = await caches.match(request);
          return cached || caches.match(OFFLINE_URL);
        })
    );
    return;
  }

  if (isNavigation) {
    event.respondWith(
      fetch(request)
        .then((networkResponse) => {
          if (networkResponse && networkResponse.status === 200) {
            const responseToCache = networkResponse.clone();
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(request, responseToCache);
            });
          }
          return networkResponse;
        })
        .catch(async () => {
          const cached = await caches.match(request);
          return cached || caches.match(OFFLINE_URL);
        })
    );
    return;
  }

  event.respondWith(
    caches.match(request).then((cachedResponse) => {
      if (cachedResponse) {
        return cachedResponse;
      }

      return fetch(request)
        .then((networkResponse) => {
          if (!networkResponse || networkResponse.status !== 200 || networkResponse.type !== 'basic') {
            return networkResponse;
          }

          const responseToCache = networkResponse.clone();
          caches.open(CACHE_NAME).then((cache) => {
            cache.put(request, responseToCache);
          });

          return networkResponse;
        })
        .catch(() => {
          if (request.mode === 'navigate') {
            return caches.match(OFFLINE_URL);
          }
          return new Response(null, { status: 504, statusText: 'Gateway Timeout' });
        });
    })
  );
});

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});
