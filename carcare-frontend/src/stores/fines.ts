// src/stores/fines.ts
import { writable, derived, get } from 'svelte/store';
import type { Fine } from '../lib/types';
import { fetchFines, addFine, updateFine, deleteFine } from '../lib/api';

export const finesList = writable<Fine[]>([]);
export const finesLoading = writable(false);
export const finesError = writable<string | null>(null);
export const finesFilters = writable<{ [key: string]: string }>({});

export const filteredFinesList = derived([finesList, finesFilters], ([$finesList, $finesFilters]) => {
  if (!$finesFilters || Object.keys($finesFilters).length === 0) return $finesList;
  return $finesList.filter((fine) => {
    return Object.entries($finesFilters).every(([key, val]) => {
      if (!val) return true;
      const field = (fine as Record<string, unknown>)[key];
      return field && String(field).toLowerCase().includes(String(val).toLowerCase());
    });
  });
});

export async function loadFines() {
  finesLoading.set(true);
  finesError.set(null);
  try {
    const filters = get(finesFilters);
    const data = await fetchFines(filters);
    finesList.set(Array.isArray(data) ? data : []);
  } catch (e: unknown) {
    finesError.set(e instanceof Error ? e.message : 'Ошибка загрузки штрафов');
  } finally {
    finesLoading.set(false);
  }
}

export async function createFine(fine: Omit<Fine, 'id'>) {
  finesLoading.set(true);
  finesError.set(null);
  try {
    const newFine = await addFine(fine);
    finesList.update(list => [newFine, ...list]);
    return newFine;
  } catch (e: unknown) {
    finesError.set(e instanceof Error ? e.message : 'Ошибка добавления штрафа');
    throw e;
  } finally {
    finesLoading.set(false);
  }
}

export async function editFine(id: string, fine: Partial<Fine>) {
  finesLoading.set(true);
  finesError.set(null);
  try {
    const updated = await updateFine(id, fine);
    finesList.update(list => list.map(f => f.id === id ? updated : f));
    return updated;
  } catch (e: unknown) {
    finesError.set(e instanceof Error ? e.message : 'Ошибка обновления штрафа');
    throw e;
  } finally {
    finesLoading.set(false);
  }
}

export async function removeFine(id: string) {
  finesLoading.set(true);
  finesError.set(null);
  try {
    await deleteFine(id);
    finesList.update(list => list.filter(f => f.id !== id));
  } catch (e: unknown) {
    finesError.set(e instanceof Error ? e.message : 'Ошибка удаления штрафа');
    throw e;
  } finally {
    finesLoading.set(false);
  }
}
