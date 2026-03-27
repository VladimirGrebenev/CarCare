// src/stores/reports.ts
import { writable, derived } from 'svelte/store';

export type ReportItem = {
  period: string;
  type: string;
  amount: number;
  count: number;
};

export type ReportSummary = {
  totalFuel: number;
  totalMaintenance: number;
  totalFines: number;
  total: number;
};

export const reportsList = writable<ReportItem[]>([]);
export const reportsLoading = writable(false);
export const reportsError = writable<string | null>(null);
export const reportsFilters = writable<{ period?: string; carId?: string; type?: string }>({});

export const reportsSummary = derived(reportsList, ($list) => {
  return $list.reduce<ReportSummary>(
    (acc, item) => {
      const amount = Number(item.amount) || 0;
      if (item.type === 'fuel') acc.totalFuel += amount;
      else if (item.type === 'maintenance') acc.totalMaintenance += amount;
      else if (item.type === 'fine') acc.totalFines += amount;
      acc.total += amount;
      return acc;
    },
    { totalFuel: 0, totalMaintenance: 0, totalFines: 0, total: 0 }
  );
});

export async function loadReports(filters: { period?: string; carId?: string; type?: string } = {}) {
  reportsLoading.set(true);
  reportsError.set(null);
  reportsFilters.set(filters);
  try {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([k, v]) => { if (v) params.append(k, v); });
    const res = await fetch(`/api/reports?${params.toString()}`, {
      credentials: 'include',
      headers: {
        Authorization: `Bearer ${typeof window !== 'undefined' ? localStorage.getItem('authToken') ?? '' : ''}`
      }
    });
    if (!res.ok) throw new Error('Ошибка загрузки отчётов');
    const data = await res.json();
    reportsList.set(Array.isArray(data) ? data : []);
  } catch (e: unknown) {
    reportsError.set(e instanceof Error ? e.message : 'Ошибка загрузки');
    reportsList.set([]);
  } finally {
    reportsLoading.set(false);
  }
}
