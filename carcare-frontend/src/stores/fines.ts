// src/stores/fines.ts
import { $state } from 'svelte/store';
import type { Fine } from '../lib/types';

export const fines = $state<Fine[]>([]);
export const finesLoading = $state(false);
export const finesError = $state<string | null>(null);
export const finesSuccess = $state<string | null>(null);

export const finesFilters = $state<Record<string, any>>({});
