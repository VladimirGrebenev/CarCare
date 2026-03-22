self.addEventListener('install', (event) => {
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
});

self.addEventListener('fetch', (event) => {
  // Simple cache-first strategy
  event.respondWith(
    caches.open('carcare-v1').then((cache) =>
      cache.match(event.request).then((resp) =>
        resp || fetch(event.request).then((response) => {
          cache.put(event.request, response.clone());
          return response;
        })
      )
    )
  );
});
