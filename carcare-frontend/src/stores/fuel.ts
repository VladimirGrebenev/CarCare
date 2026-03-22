// src/stores/fuel.ts
// Svelte 5 rune store for fuel data (CRUD, filters, loading, error states)

import { writable, derived, get } from 'svelte/store';
import { fetchFuelHistory, addFuel, updateFuel, deleteFuel } from '../lib/api';

export type Fuel = {
  id?: string;
  date: string;
  liters: number | string;
  price: number | string;
  carId: string;
};

export const fuelList = writable<Fuel[]>([]);
export const fuelLoading = writable(false);
export const fuelError = writable<string | null>(null);
export const fuelFilters = writable<{ [key: string]: string }>({});

export const filteredFuelList = derived([
  fuelList, fuelFilters
], ([$fuelList, $fuelFilters]) => {
  if (!$fuelFilters || Object.keys($fuelFilters).length === 0) return $fuelList;
  return $fuelList.filter((fuel: Record<string, any>) => {
    // Fallback for Object.entries and includes
    const entries = Object.entries ? Object.entries($fuelFilters) : Object.keys($fuelFilters).map(key => [key, $fuelFilters[key]]);
    return entries.every(([key, val]) => {
      if (!val) return true;
      const field = fuel[key];
      return field && (typeof field === 'string' ? field.includes(String(val)) : String(field).includes(String(val)));
    });
  });
});

export async function loadFuel() {
  fuelLoading.set(true);
  fuelError.set(null);
  try {
    const filters = get(fuelFilters);
    const data = await fetchFuelHistory(filters);
    fuelList.set(data);
  } catch (e: any) {
    fuelError.set(e.message || 'Ошибка загрузки заправок');
  } finally {
    fuelLoading.set(false);
  }
}

export async function createFuel(fuel: Fuel) {
  fuelLoading.set(true);
  fuelError.set(null);
  try {
    const newFuel = await addFuel(fuel);
    fuelList.update(list => [newFuel, ...list]);
  } catch (e: any) {
    fuelError.set(e.message || 'Ошибка добавления заправки');
  } finally {
    fuelLoading.set(false);
  }
}

export async function editFuel(id: string, fuel: Fuel) {
  fuelLoading.set(true);
  fuelError.set(null);
  try {
    const updated = await updateFuel(id, fuel);
    fuelList.update(list => list.map(f => f.id === id ? updated : f));
  } catch (e: any) {
    fuelError.set(e.message || 'Ошибка обновления заправки');
  } finally {
    fuelLoading.set(false);
  }
}

export async function removeFuel(id: string) {
  fuelLoading.set(true);
  fuelError.set(null);
  try {
    await deleteFuel(id);
    fuelList.update(list => list.filter(f => f.id !== id));
  } catch (e: any) {
    fuelError.set(e.message || 'Ошибка удаления заправки');
  } finally {
    fuelLoading.set(false);
  }
}
