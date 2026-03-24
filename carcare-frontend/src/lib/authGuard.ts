// src/lib/authGuard.ts
// Auth guard utility for protected routes

import { goto } from '$app/navigation';
import { bootstrapAuth, authToken } from '../stores/auth';

export async function ensureAuthenticated(): Promise<boolean> {
  const unsubscribe = authToken.subscribe(() => {});
  unsubscribe();

  const isAuthenticated = await bootstrapAuth();
  if (!isAuthenticated) {
    await goto('/welcome', { replaceState: true });
    return false;
  }
  return true;
}
