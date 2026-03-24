// src/lib/api.ts
// API stubs for CarCare frontend

function getAuthToken(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }
  return localStorage.getItem('authToken');
}

function withAuthHeaders(headers: Record<string, string> = {}) {
  const token = getAuthToken();
  return token
    ? { ...headers, Authorization: `Bearer ${token}` }
    : headers;
}

function normalizeAuthPayload(payload: unknown, fallbackEmail?: string) {
  if (!payload || typeof payload !== 'object') {
    throw new Error('Некорректный ответ авторизации');
  }

  const data = payload as Record<string, unknown>;
  const rawToken = data.token ?? data.access_token;
  if (typeof rawToken !== 'string' || rawToken.length === 0) {
    throw new Error('Токен авторизации не получен');
  }

  let authUser = data.user;
  if (!authUser || typeof authUser !== 'object') {
    authUser = fallbackEmail ? { email: fallbackEmail } : { id: 'current-user' };
  }

  return {
    token: rawToken,
    user: authUser as Record<string, unknown>
  };
}

/**
 * Handle API errors gracefully, distinguishing between auth and temporary errors
 */
function createErrorMessage(status: number, action: string): string {
  if (status === 401 || status === 403) {
    return 'Требуется авторизация';
  }
  if (status >= 500) {
    return `Служба недоступна (${action})`;
  }
  if (status === 404) {
    return `Не найдено (${action})`;
  }
  return `Ошибка ${status} (${action})`;
}

// --- Profile ---
export async function fetchProfile() {
  const res = await fetch('/api/profile', {
    credentials: 'include',
    headers: withAuthHeaders(),
    cache: 'no-store'
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'загрузка профиля'));
  }
  return res.json();
}

export async function addCar(car) {
  const res = await fetch('/api/cars', {
    method: 'POST',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(car)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'добавление авто'));
  }
  return res.json();
}

// --- Cars ---
export async function fetchCars() {
  const res = await fetch('/api/cars', {
    credentials: 'include',
    headers: withAuthHeaders()
  });
  
  // If not found or server error, try fetching from profile
  if (res.status === 404 || res.status === 500) {
    try {
      const profile = await fetchProfile();
      return Array.isArray(profile?.cars) ? profile.cars : [];
    } catch (_e) {
      return [];
    }
  }
  
  // 401/403: auth error
  if (res.status === 401 || res.status === 403) {
    throw new Error(createErrorMessage(res.status, 'авто'));
  }
  
  // Other errors: return empty gracefully
  if (!res.ok) return [];
  
  try {
    return await res.json();
  } catch (_e) {
    return [];
  }
}


// --- Fuel CRUD ---
export async function fetchFuelHistory(filters = {}) {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') params.append(key, String(val));
  });
  const res = await fetch(`/api/fuel?${params.toString()}`, {
    credentials: 'include',
    headers: withAuthHeaders()
  });
  
  // 401/403: throw error to trigger auth check
  if (res.status === 401 || res.status === 403) {
    throw new Error(createErrorMessage(res.status, 'заправки'));
  }
  
  // Other errors: return empty gracefully
  if (!res.ok) return [];
  
  try {
    return await res.json();
  } catch (_e) {
    return [];
  }
}

export async function addFuel(fuel) {
  const res = await fetch('/api/fuel', {
    method: 'POST',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(fuel)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'добавление заправки'));
  }
  return res.json();
}

export async function updateFuel(id, fuel) {
  const res = await fetch(`/api/fuel/${id}`, {
    method: 'PUT',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(fuel)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'обновление заправки'));
  }
  return res.json();
}

export async function deleteFuel(id) {
  const res = await fetch(`/api/fuel/${id}`, {
    method: 'DELETE',
    credentials: 'include',
    headers: withAuthHeaders()
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'удаление заправки'));
  }
  return { success: true };
}


// --- Maintenance CRUD ---
export async function fetchMaintenanceHistory(filters = {}) {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') params.append(key, String(val));
  });
  const res = await fetch(`/api/maintenance?${params.toString()}`, {
    credentials: 'include',
    headers: withAuthHeaders()
  });
  
  // 401/403: throw error to trigger auth check
  if (res.status === 401 || res.status === 403) {
    throw new Error(createErrorMessage(res.status, 'ТО'));
  }
  
  // Other errors: return empty gracefully
  if (!res.ok) return [];
  
  try {
    return await res.json();
  } catch (_e) {
    return [];
  }
}

export async function addMaintenance(maintenance) {
  const res = await fetch('/api/maintenance', {
    method: 'POST',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(maintenance)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'добавление ТО'));
  }
  return res.json();
}

export async function updateMaintenance(id, maintenance) {
  const res = await fetch(`/api/maintenance/${id}`, {
    method: 'PUT',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(maintenance)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'обновление ТО'));
  }
  return res.json();
}

export async function deleteMaintenance(id) {
  const res = await fetch(`/api/maintenance/${id}`, {
    method: 'DELETE',
    credentials: 'include',
    headers: withAuthHeaders()
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'удаление ТО'));
  }
  return { success: true };
}


// --- Fines CRUD ---
export async function fetchFines(filters = {}) {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') params.append(key, String(val));
  });
  const res = await fetch(`/api/fines?${params.toString()}`, {
    credentials: 'include',
    headers: withAuthHeaders()
  });
  
  // 401/403: throw error to trigger auth check
  if (res.status === 401 || res.status === 403) {
    throw new Error(createErrorMessage(res.status, 'штрафы'));
  }
  
  // Other errors: return empty gracefully
  if (!res.ok) return [];
  
  try {
    return await res.json();
  } catch (_e) {
    return [];
  }
}

export async function addFine(fine) {
  const res = await fetch('/api/fines', {
    method: 'POST',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(fine)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'добавление штрафа'));
  }
  return res.json();
}

export async function updateFine(id, fine) {
  const res = await fetch(`/api/fines/${id}`, {
    method: 'PUT',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(fine)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'обновление штрафа'));
  }
  return res.json();
}

export async function deleteFine(id) {
  const res = await fetch(`/api/fines/${id}`, {
    method: 'DELETE',
    credentials: 'include',
    headers: withAuthHeaders()
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'удаление штрафа'));
  }
  return { success: true };
}

export async function fetchUsers() {
  // TODO: Implement API call to backend
  return [];
}

// --- Auth API ---
export async function login(email: string, password: string) {
  const res = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    cache: 'no-store'
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Ошибка авторизации');
  }
  return normalizeAuthPayload(await res.json(), email);
}

export async function register(email: string, password: string) {
  const res = await fetch('/api/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    cache: 'no-store'
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Ошибка регистрации');
  }
  return normalizeAuthPayload(await res.json(), email);
}

export async function oauthLogin(provider: string) {
  // Redirect to backend OAuth endpoint
  window.location.href = `/api/auth/oauth/${provider}`;
}
