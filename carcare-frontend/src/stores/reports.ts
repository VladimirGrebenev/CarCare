// src/stores/reports.ts
import { $state } from 'svelte/store';

export const reportsStore = $state({
  items: [],
  loading: false,
  error: null,
  filters: {},
  async load(filters = {}) {
    this.loading = true;
    this.error = null;
    this.filters = filters;
    try {
      // TODO: fetchReports API integration
      const res = await fetch('/api/reports?' + new URLSearchParams(filters), { credentials: 'include' });
      if (!res.ok) throw new Error('Ошибка загрузки отчётов');
      this.items = await res.json();
    } catch (e) {
      this.error = e instanceof Error ? e.message : 'Ошибка загрузки';
    } finally {
      this.loading = false;
    }
  },
  async export(format = 'csv') {
    try {
      const params = new URLSearchParams({ ...this.filters, format });
      const res = await fetch('/api/reports/export?' + params, { credentials: 'include' });
      if (!res.ok) throw new Error('Ошибка экспорта');
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `reports.${format}`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (e) {
      this.error = e instanceof Error ? e.message : 'Ошибка экспорта';
    }
  }
});
