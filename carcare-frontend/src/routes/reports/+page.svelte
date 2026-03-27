<script lang="ts">
  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Table from '../../components/ui/Table.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import {
    reportsList,
    reportsLoading,
    reportsError,
    reportsSummary,
    loadReports,
  } from '../../stores/reports';
  import { fetchCars } from '../../lib/api';
  import type { Car } from '../../lib/types';

  const COLUMNS = [
    { label: 'Период', key: 'period' },
    { label: 'Тип', key: '_typeLabel' },
    { label: 'Кол-во', key: 'count' },
    { label: 'Сумма (₽)', key: '_amount' },
  ];

  const PERIOD_OPTIONS = [
    { value: '', label: 'Все периоды' },
    { value: 'month', label: 'Текущий месяц' },
    { value: 'quarter', label: 'Квартал' },
    { value: 'year', label: 'Текущий год' },
  ];

  const TYPE_OPTIONS = [
    { value: '', label: 'Все типы расходов' },
    { value: 'fuel', label: 'Заправки' },
    { value: 'maintenance', label: 'Техобслуживание' },
    { value: 'fine', label: 'Штрафы' },
  ];

  const TYPE_LABELS: Record<string, string> = {
    fuel: 'Заправки',
    maintenance: 'ТО',
    fine: 'Штрафы',
  };

  let cars = $state<Car[]>([]);
  let filterPeriod = $state('');
  let filterCarId = $state('');
  let filterType = $state('');
  let toast = $state({ open: false, message: '', type: 'error' as 'info' | 'success' | 'error' });

  let tableRows = $derived(
    $reportsList.map(r => ({
      ...r,
      _typeLabel: TYPE_LABELS[r.type] ?? r.type,
      _amount: Number(r.amount).toLocaleString('ru-RU', { minimumFractionDigits: 2 }) + ' ₽',
    }))
  );

  const maxBar = $derived(
    Math.max($reportsSummary.totalFuel, $reportsSummary.totalMaintenance, $reportsSummary.totalFines, 1)
  );

  const summaryCards = $derived([
    {
      label: 'Заправки',
      value: $reportsSummary.totalFuel,
      color: 'var(--accent)',
      bgColor: 'var(--accent-light)',
      icon: '⛽',
    },
    {
      label: 'Техобслуживание',
      value: $reportsSummary.totalMaintenance,
      color: 'var(--success)',
      bgColor: 'var(--success-light)',
      icon: '🔧',
    },
    {
      label: 'Штрафы',
      value: $reportsSummary.totalFines,
      color: 'var(--danger)',
      bgColor: 'var(--danger-light)',
      icon: '📋',
    },
    {
      label: 'Итого',
      value: $reportsSummary.total,
      color: 'var(--text-primary)',
      bgColor: 'rgba(255,255,255,0.04)',
      icon: '💰',
    },
  ]);

  function applyFilters() {
    loadReports({ period: filterPeriod, carId: filterCarId, type: filterType });
  }

  function resetFilters() {
    filterPeriod = '';
    filterCarId = '';
    filterType = '';
    loadReports({});
  }

  onMount(async () => {
    await ensureAuthenticated();
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch {
      cars = [];
    }
    loadReports({});
  });
</script>

<PageLayout title="Отчёты и статистика">
  <!-- Filters -->
  <div class="filters-panel">
    <div class="filters-row">
      <div class="filter-field">
        <label class="filter-label" for="report-period">Период</label>
        <select id="report-period" class="field-select" bind:value={filterPeriod}>
          {#each PERIOD_OPTIONS as opt}
            <option value={opt.value}>{opt.label}</option>
          {/each}
        </select>
      </div>
      <div class="filter-field">
        <label class="filter-label" for="report-car">Автомобиль</label>
        <select id="report-car" class="field-select" bind:value={filterCarId}>
          <option value="">Все автомобили</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
      <div class="filter-field">
        <label class="filter-label" for="report-type">Тип расходов</label>
        <select id="report-type" class="field-select" bind:value={filterType}>
          {#each TYPE_OPTIONS as opt}
            <option value={opt.value}>{opt.label}</option>
          {/each}
        </select>
      </div>
      <div class="filter-actions">
        <Button variant="primary" onclick={applyFilters}>Применить</Button>
        <Button variant="secondary" onclick={resetFilters}>Сброс</Button>
      </div>
    </div>
  </div>

  {#if $reportsLoading}
    <div class="center-loader">
      <Loader size={48} />
    </div>
  {:else if $reportsError}
    <div class="error-banner" role="alert">
      <span class="error-icon">⚠️</span>
      <span>{$reportsError}</span>
    </div>
  {:else}
    <!-- Summary cards -->
    <div class="summary-grid">
      {#each summaryCards as card}
        <div class="summary-card" style="--card-color:{card.color}; --card-bg:{card.bgColor};">
          <div class="summary-icon">{card.icon}</div>
          <div class="summary-body">
            <div class="summary-label">{card.label}</div>
            <div class="summary-value">{card.value.toLocaleString('ru-RU')} ₽</div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Bar chart -->
    {#if $reportsSummary.total > 0}
      <div class="chart-section">
        <h3 class="section-title">Распределение расходов</h3>
        <div class="bar-chart">
          {#each [
            { label: 'Заправки', value: $reportsSummary.totalFuel, color: 'var(--accent)' },
            { label: 'Техобслуживание', value: $reportsSummary.totalMaintenance, color: 'var(--success)' },
            { label: 'Штрафы', value: $reportsSummary.totalFines, color: 'var(--danger)' },
          ] as bar}
            <div class="bar-row">
              <span class="bar-label">{bar.label}</span>
              <div class="bar-track">
                <div
                  class="bar-fill"
                  style="width:{(bar.value / maxBar * 100).toFixed(1)}%; background:{bar.color};"
                  role="progressbar"
                  aria-valuenow={bar.value}
                  aria-valuemin={0}
                  aria-valuemax={maxBar}
                ></div>
              </div>
              <span class="bar-amount">
                {bar.value.toLocaleString('ru-RU')} ₽
                {#if $reportsSummary.total > 0}
                  <span class="bar-percent">({(bar.value / $reportsSummary.total * 100).toFixed(0)}%)</span>
                {/if}
              </span>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Detail table -->
    <div class="detail-section">
      <h3 class="section-title">Детализация</h3>
      {#if $reportsList.length === 0}
        <EmptyState message="Нет данных за выбранный период. Попробуйте изменить фильтры." />
      {:else}
        <Table
          columns={COLUMNS}
          rows={tableRows}
          emptyText="Нет данных"
        />
      {/if}
    </div>
  {/if}

  <Toast open={toast.open} message={toast.message} type={toast.type} />
</PageLayout>

<style>
/* ── Filters ── */
.filters-panel {
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 1.25rem 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: var(--shadow-sm);
}
.filters-row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-end;
  flex-wrap: wrap;
}
.filter-field { display: flex; flex-direction: column; gap: 0.375rem; min-width: 160px; }
.filter-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }
.filter-actions { display: flex; gap: 0.5rem; align-items: center; padding-top: 1.25rem; }

.field-select {
  padding: 0.5625rem 0.875rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text-primary);
  font-size: 0.9375rem;
  font-family: var(--font);
  outline: none;
  cursor: pointer;
  transition: border-color var(--transition), box-shadow var(--transition);
}
.field-select:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
}

/* ── Summary cards ── */
.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}
@media (max-width: 900px) { .summary-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 500px) { .summary-grid { grid-template-columns: 1fr; } }

.summary-card {
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-left: 3px solid var(--card-color);
  border-radius: var(--radius-lg);
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition), box-shadow var(--transition);
}
.summary-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.summary-icon { font-size: 1.75rem; flex-shrink: 0; }
.summary-body { display: flex; flex-direction: column; gap: 0.25rem; min-width: 0; }
.summary-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-secondary);
}
.summary-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--card-color);
  white-space: nowrap;
}

/* ── Bar chart ── */
.section-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 1rem;
}

.chart-section { margin-bottom: 1.5rem; }
.bar-chart { display: flex; flex-direction: column; gap: 0.875rem; }
.bar-row { display: flex; align-items: center; gap: 1rem; }
.bar-label { width: 120px; font-size: 0.875rem; color: var(--text-secondary); text-align: right; flex-shrink: 0; }
.bar-track {
  flex: 1;
  height: 14px;
  background: var(--border);
  border-radius: 100px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  border-radius: 100px;
  min-width: 2px;
  transition: width 0.7s cubic-bezier(.4,0,.2,1);
}
.bar-amount {
  min-width: 130px;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
}
.bar-percent {
  font-size: 0.75rem;
  font-weight: 400;
  color: var(--text-secondary);
  margin-left: 0.25rem;
}

/* ── Detail table ── */
.detail-section { margin-bottom: 2rem; }

/* ── States ── */
.center-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4rem 1rem;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  background: var(--danger-light);
  border: 1px solid rgba(248, 81, 73, 0.3);
  border-radius: var(--radius-lg);
  color: var(--danger);
  font-weight: 500;
  margin-bottom: 1.5rem;
}
.error-icon { font-size: 1.25rem; flex-shrink: 0; }
</style>
