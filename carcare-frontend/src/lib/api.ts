// src/lib/api.ts
// API stubs for CarCare frontend


// --- Profile ---
export async function fetchProfile() {
  const res = await fetch('/api/profile', {
    credentials: 'include'
  });
  if (!res.ok) throw new Error('Ошибка загрузки профиля');
  return res.json();
}

export async function addCar(car) {
  const res = await fetch('/api/cars', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(car)
  });
  if (!res.ok) throw new Error('Ошибка добавления авто');
  return res.json();
}

// --- Cars ---
export async function fetchCars() {
  const res = await fetch('/api/cars', { credentials: 'include' });
  if (!res.ok) throw new Error('Ошибка загрузки авто');
  return res.json();
}


// --- Fuel CRUD ---
export async function fetchFuelHistory(filters = {}) {
  // Преобразуем фильтры в query string
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') params.append(key, val);
  });
  const res = await fetch(`/api/fuel?${params.toString()}`, { credentials: 'include' });
  if (!res.ok) throw new Error('Ошибка загрузки заправок');
  return res.json();
}

export async function addFuel(fuel) {
  const res = await fetch('/api/fuel', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(fuel)
  });
  if (!res.ok) throw new Error('Ошибка добавления заправки');
  return res.json();
}

export async function updateFuel(id, fuel) {
  const res = await fetch(`/api/fuel/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(fuel)
  });
  if (!res.ok) throw new Error('Ошибка обновления заправки');
  return res.json();
}

export async function deleteFuel(id) {
  const res = await fetch(`/api/fuel/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) throw new Error('Ошибка удаления заправки');
  return { success: true };
}


// --- Maintenance CRUD ---
export async function fetchMaintenanceHistory(filters = {}) {
  // Преобразуем фильтры в query string
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') params.append(key, val);
  });
  const res = await fetch(`/api/maintenance?${params.toString()}`, { credentials: 'include' });
  if (!res.ok) throw new Error('Ошибка загрузки ТО');
  return res.json();
}

export async function addMaintenance(maintenance) {
  const res = await fetch('/api/maintenance', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(maintenance)
  });
  if (!res.ok) throw new Error('Ошибка добавления ТО');
  return res.json();
}

export async function updateMaintenance(id, maintenance) {
  const res = await fetch(`/api/maintenance/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(maintenance)
  });
  if (!res.ok) throw new Error('Ошибка обновления ТО');
  return res.json();
}

export async function deleteMaintenance(id) {
  const res = await fetch(`/api/maintenance/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) throw new Error('Ошибка удаления ТО');
  return { success: true };
}


// --- Fines CRUD ---
export async function fetchFines(filters = {}) {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') params.append(key, val);
  });
  const res = await fetch(`/api/fines?${params.toString()}`, { credentials: 'include' });
  if (!res.ok) throw new Error('Ошибка загрузки штрафов');
  return res.json();
}

export async function addFine(fine) {
  const res = await fetch('/api/fines', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(fine)
  });
  if (!res.ok) throw new Error('Ошибка добавления штрафа');
  return res.json();
}

export async function updateFine(id, fine) {
  const res = await fetch(`/api/fines/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(fine)
  });
  if (!res.ok) throw new Error('Ошибка обновления штрафа');
  return res.json();
}

export async function deleteFine(id) {
  const res = await fetch(`/api/fines/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  });
  if (!res.ok) throw new Error('Ошибка удаления штрафа');
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
    body: JSON.stringify({ email, password })
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Ошибка авторизации');
  }
  return res.json(); // { token, user }
}

export async function oauthLogin(provider: string) {
  // Redirect to backend OAuth endpoint
  window.location.href = `/api/auth/oauth/${provider}`;
}
