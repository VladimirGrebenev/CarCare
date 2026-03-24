// src/stores/fines.ts
import { writable } from 'svelte/store';
import type { Fine } from '../lib/types';

export const fines = writable<Fine[]>([]);
export const finesLoading = writable(false);
export const finesError = writable<string | null>(null);
export const finesSuccess = writable<string | null>(null);

export const finesFilters = writable<Record<string, string | number | boolean>>({});
