// src/stores/maintenance.ts
// Svelte 5 rune store for maintenance data
import { fetchMaintenanceHistory, addMaintenance, updateMaintenance, deleteMaintenance } from '../lib/api';

export const maintenanceStore = {
  items: [],
  loading: false,
  error: null,
  async load(filters = {}) {
    this.loading = true;
    this.error = null;
    try {
      this.items = await fetchMaintenanceHistory(filters);
    } catch (e) {
      this.error = e instanceof Error ? e.message : 'Ошибка загрузки';
    } finally {
      this.loading = false;
    }
  },
  async add(data) {
    this.loading = true;
    this.error = null;
    try {
      const created = await addMaintenance(data);
      this.items = [created, ...this.items];
      return created;
    } catch (e) {
      this.error = e instanceof Error ? e.message : 'Ошибка добавления';
      throw e;
    } finally {
      this.loading = false;
    }
  },
  async update(id, data) {
    this.loading = true;
    this.error = null;
    try {
      const updated = await updateMaintenance(id, data);
      this.items = this.items.map(item => item.id === id ? updated : item);
      return updated;
    } catch (e) {
      this.error = e instanceof Error ? e.message : 'Ошибка обновления';
      throw e;
    } finally {
      this.loading = false;
    }
  },
  async remove(id) {
    this.loading = true;
    this.error = null;
    try {
      await deleteMaintenance(id);
      this.items = this.items.filter(item => item.id !== id);
    } catch (e) {
      this.error = e instanceof Error ? e.message : 'Ошибка удаления';
      throw e;
    } finally {
      this.loading = false;
    }
  }
};
