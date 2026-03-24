// src/stores/auth.ts
import { writable } from 'svelte/store';

type AuthUser = Record<string, unknown>;

export const authToken = writable<string | null>(null);
export const user = writable<AuthUser | null>(null);

export function setAuth(token: string, userData: AuthUser) {
  authToken.set(token);
  user.set(userData);
  localStorage.setItem('authToken', token);
  localStorage.setItem('user', JSON.stringify(userData));
}

export function clearAuth() {
  authToken.set(null);
  user.set(null);
  localStorage.removeItem('authToken');
  localStorage.removeItem('user');
}

export function loadAuth() {
  const token = localStorage.getItem('authToken');
  const userData = localStorage.getItem('user');
  if (token && userData) {
    authToken.set(token);
    user.set(JSON.parse(userData) as AuthUser);
  }
}
