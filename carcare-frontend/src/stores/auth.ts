// src/stores/auth.ts
import { writable } from 'svelte/store';

type AuthUser = Record<string, unknown>;

export const authToken = writable<string | null>(null);
export const user = writable<AuthUser | null>(null);
export const authReady = writable(false);

export function setAuth(token: string, userData?: AuthUser | null) {
  authToken.set(token);
  user.set(userData ?? null);
  if (typeof window !== 'undefined') {
    localStorage.setItem('authToken', token);
    if (userData) {
      localStorage.setItem('user', JSON.stringify(userData));
    } else {
      localStorage.removeItem('user');
    }
  }
  authReady.set(true);
}

export function clearAuth() {
  authToken.set(null);
  user.set(null);
  if (typeof window !== 'undefined') {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
  }
  authReady.set(true);
}

export function loadAuth() {
  if (typeof window === 'undefined') {
    return false;
  }

  const token = localStorage.getItem('authToken');
  const userData = localStorage.getItem('user');
  if (token) {
    authToken.set(token);
    if (userData) {
      try {
        user.set(JSON.parse(userData) as AuthUser);
      } catch {
        localStorage.removeItem('user');
        user.set(null);
      }
    }
    return true;
  }

  return false;
}

export async function bootstrapAuth(): Promise<boolean> {
  authReady.set(false);

  const hasLocalAuth = loadAuth();
  if (!hasLocalAuth) {
    clearAuth();
    return false;
  }

  try {
    const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;

    const response = await fetch('/api/profile', {
      method: 'GET',
      credentials: 'include',
      cache: 'no-store',
      headers: token ? { Authorization: `Bearer ${token}` } : undefined
    });

    // 401/403: Auth failed - clear and return false
    if (response.status === 401 || response.status === 403) {
      clearAuth();
      authReady.set(true);
      return false;
    }

    // Other status codes: try to validate token even if status is not perfect
    // (handles temporary server issues gracefully)
    if (!response.ok) {
      // If backend is temporarily unavailable (5xx), keep local auth as fallback
      // Only clear auth if we're sure it's invalid
      if (response.status >= 500) {
        authReady.set(true);
        return true; // Keep local auth as fallback
      }
      clearAuth();
      authReady.set(true);
      return false;
    }

    const profile = await response.json().catch(() => null);
    if (profile) {
      user.set(profile as AuthUser);
      if (typeof window !== 'undefined') {
        localStorage.setItem('user', JSON.stringify(profile));
      }
    }

    authReady.set(true);
    return true;
  } catch (_error) {
    // Network error: keep local auth as fallback
    authReady.set(true);
    return true;
  }
}
