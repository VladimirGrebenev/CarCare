<script lang="ts">
  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';
  import {
    fetchFuelHistory,
    fetchMaintenanceHistory,
    fetchFines,
    fetchCars,
  } from '../../lib/api';
  import type { Car } from '../../lib/types';

  // ── Types ────────────────────────────────────────────────────────────────────

  type FuelEvent = {
    id: string;
    date: string;
    liters: number;
    price: number;
    fuelType?: string;
    carId?: string;
  };

  type MaintenanceEvent = {
    id: string;
    date: string;
    type: string;
    cost: number;
    carId?: string;
  };

  type Fine = {
    id: string;
    date: string;
    type?: string;
    amount: number;
    status?: string;
    description?: string;
    carId?: string;
  };

  // ── Constants ─────────────────────────────────────────────────────────────────

  const FUEL_TYPE_LABELS: Record<string, string> = {
    petrol: 'Бензин',
    diesel: 'Дизель',
    gas: 'Газ',
    electric: 'Электро',
  };

  type ReportSummary = {
    total_fuel_cost: number;
    total_maintenance_cost: number;
    total_fines_amount: number;
    fuel_count: number;
    maintenance_count: number;
    fines_count: number;
  };

  type ViewMode = 'table' | 'chart';
  type Period = 'week' | 'month' | 'custom';

  // ── State ────────────────────────────────────────────────────────────────────

  type Tab = 'overview' | 'fuel' | 'maintenance' | 'fines';
  let activeTab = $state<Tab>('overview');

  let loading = $state(true);
  let error = $state<string | null>(null);

  let fuelList = $state<FuelEvent[]>([]);
  let maintenanceList = $state<MaintenanceEvent[]>([]);
  let finesList = $state<Fine[]>([]);
  let carsList = $state<Car[]>([]);
  let summary = $state<ReportSummary>({
    total_fuel_cost: 0,
    total_maintenance_cost: 0,
    total_fines_amount: 0,
    fuel_count: 0,
    maintenance_count: 0,
    fines_count: 0,
  });

  // Car filter
  let filterCarId = $state('');

  // Tooltip state
  let tooltip = $state({ visible: false, x: 0, y: 0, text: '' });

  function showTooltip(e: MouseEvent, text: string) {
    tooltip = { visible: true, x: e.clientX + 12, y: e.clientY - 8, text };
  }
  function moveTooltip(e: MouseEvent) {
    if (tooltip.visible) {
      tooltip = { ...tooltip, x: e.clientX + 12, y: e.clientY - 8 };
    }
  }
  function hideTooltip() {
    tooltip = { ...tooltip, visible: false };
  }

  // View/Period state
  let fuelView = $state<ViewMode>('table');
  let fuelPeriod = $state<Period>('month');
  let fuelDateFrom = $state('');
  let fuelDateTo = $state('');

  let maintenanceView = $state<ViewMode>('table');
  let maintenancePeriod = $state<Period>('month');
  let maintenanceDateFrom = $state('');
  let maintenanceDateTo = $state('');

  let finesView = $state<ViewMode>('table');
  let finesPeriod = $state<Period>('month');
  let finesDateFrom = $state('');
  let finesDateTo = $state('');

  // ── Derived ──────────────────────────────────────────────────────────────────

  function fuelCost(e: FuelEvent): number {
    return Number(e.liters) * Number(e.price);
  }

  // Filtered by car
  const filteredFuelList = $derived(
    filterCarId ? fuelList.filter(e => e.carId === filterCarId) : fuelList
  );
  const filteredMaintenanceList = $derived(
    filterCarId ? maintenanceList.filter(e => e.carId === filterCarId) : maintenanceList
  );
  const filteredFinesList = $derived(
    filterCarId ? finesList.filter(e => e.carId === filterCarId) : finesList
  );

  const fuelTotal = $derived(
    filteredFuelList.reduce((sum, e) => sum + fuelCost(e), 0)
  );

  const maintenanceTotal = $derived(
    filteredMaintenanceList.reduce((sum, e) => sum + (Number(e.cost) || 0), 0)
  );

  const finesTotal = $derived(
    filteredFinesList.reduce((sum, e) => sum + (Number(e.amount) || 0), 0)
  );

  const grandTotal = $derived(fuelTotal + maintenanceTotal + finesTotal);

  // ── Chart helpers ─────────────────────────────────────────────────────────────

  function getDateRange(period: Period, from: string, to: string): [Date, Date] {
    const now = new Date();
    if (period === 'week') {
      const start = new Date(now); start.setDate(now.getDate() - 6);
      return [start, now];
    }
    if (period === 'month') {
      const start = new Date(now); start.setDate(now.getDate() - 29);
      return [start, now];
    }
    return [from ? new Date(from) : new Date(0), to ? new Date(to) : now];
  }

  type DayPoint = { day: string; value: number };

  function groupByDay<T extends { date: string }>(
    items: T[],
    getValue: (item: T) => number,
    period: Period, from: string, to: string
  ): DayPoint[] {
    const [start, end] = getDateRange(period, from, to);
    const map = new Map<string, number>();
    const cur = new Date(start);
    while (cur <= end) {
      map.set(cur.toISOString().split('T')[0], 0);
      cur.setDate(cur.getDate() + 1);
    }
    for (const item of items) {
      const d = item.date?.split('T')[0] ?? '';
      if (map.has(d)) map.set(d, (map.get(d) ?? 0) + getValue(item));
    }
    return Array.from(map.entries()).map(([day, value]) => ({ day, value }));
  }

  const fuelChartData = $derived(
    groupByDay(filteredFuelList, e => fuelCost(e), fuelPeriod, fuelDateFrom, fuelDateTo)
  );
  const maintenanceChartData = $derived(
    groupByDay(filteredMaintenanceList, e => Number(e.cost) || 0, maintenancePeriod, maintenanceDateFrom, maintenanceDateTo)
  );
  const finesChartData = $derived(
    groupByDay(filteredFinesList, e => Number(e.amount) || 0, finesPeriod, finesDateFrom, finesDateTo)
  );

  // ── Helpers ──────────────────────────────────────────────────────────────────

  function fmt(value: number): string {
    return value.toLocaleString('ru-RU') + ' ₽';
  }

  function fmtDate(raw: string | null | undefined): string {
    if (!raw) return '—';
    const d = new Date(raw);
    if (isNaN(d.getTime())) return String(raw);
    const day = String(d.getDate()).padStart(2, '0');
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const year = d.getFullYear();
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');
    return `${day}.${month}.${year} ${hours}:${minutes}`;
  }

  // ── Car helpers ───────────────────────────────────────────────────────────────

  function getCarLabel(carId: string): string {
    if (!carId) return '—';
    const car = carsList.find(c => c.id === carId);
    if (!car) return '—';
    return [car.brand, car.model, car.plate ? `(${car.plate})` : ''].filter(Boolean).join(' ');
  }

  // ── Data loading ─────────────────────────────────────────────────────────────

  async function loadAll() {
    loading = true;
    error = null;
    try {
      const [fuelRaw, maintenanceRaw, finesRaw, reportRaw, carsRaw] = await Promise.all([
        fetchFuelHistory().catch(() => []),
        fetchMaintenanceHistory().catch(() => []),
        fetchFines().catch(() => []),
        fetch('/api/reports', {
          credentials: 'include',
          headers: {
            Authorization: `Bearer ${typeof window !== 'undefined' ? localStorage.getItem('authToken') ?? '' : ''}`
          }
        }).then(r => r.ok ? r.json() : null).catch(() => null),
        fetchCars().catch(() => []),
      ]);

      fuelList = Array.isArray(fuelRaw) ? fuelRaw : [];
      maintenanceList = Array.isArray(maintenanceRaw) ? maintenanceRaw : [];
      finesList = Array.isArray(finesRaw) ? finesRaw : [];
      carsList = Array.isArray(carsRaw) ? carsRaw : [];

      if (reportRaw && typeof reportRaw === 'object' && !Array.isArray(reportRaw)) {
        summary = {
          total_fuel_cost: Number(reportRaw.total_fuel_cost) || 0,
          total_maintenance_cost: Number(reportRaw.total_maintenance_cost) || 0,
          total_fines_amount: Number(reportRaw.total_fines_amount) || 0,
          fuel_count: Number(reportRaw.fuel_count) || 0,
          maintenance_count: Number(reportRaw.maintenance_count) || 0,
          fines_count: Number(reportRaw.fines_count) || 0,
        };
      }
    } catch (e) {
      error = e instanceof Error ? e.message : 'Ошибка загрузки данных';
    } finally {
      loading = false;
    }
  }

  onMount(async () => {
    await ensureAuthenticated();
    await loadAll();
  });
</script>

{#snippet barChart(data: DayPoint[], color: string)}
  {@const maxVal = Math.max(...data.map(d => d.value), 1)}
  {@const w = 600}
  {@const h = 200}
  {@const pad = 40}
  {@const barW = Math.max(2, (w - pad * 2) / data.length - 2)}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="bar-chart-wrap" onmousemove={moveTooltip}>
    <svg viewBox="0 0 {w} {h + pad}" class="bar-svg">
      {#each [0, 0.25, 0.5, 0.75, 1] as frac}
        {@const y = pad + (h - h * frac)}
        <line x1={pad} y1={y} x2={w} y2={y} stroke="var(--border)" stroke-width="1"/>
        <text x={pad - 4} y={y + 4} text-anchor="end" font-size="10" fill="var(--text-secondary)">
          {Math.round(maxVal * frac).toLocaleString('ru-RU')}
        </text>
      {/each}
      {#each data as point, i}
        {@const x = pad + i * ((w - pad * 2) / data.length)}
        {@const barH = point.value > 0 ? Math.max(2, (point.value / maxVal) * h) : 0}
        {@const y = pad + h - barH}
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <rect
          x={x} y={y} width={barW} height={barH} fill={color} rx="2" opacity="0.85"
          style="cursor: pointer"
          onmouseenter={(e) => showTooltip(e, `${point.day} — ${point.value.toLocaleString('ru-RU')} ₽`)}
          onmouseleave={hideTooltip}
        />
        {#if data.length <= 14}
          <text x={x + barW / 2} y={h + pad + 12} text-anchor="middle" font-size="9" fill="var(--text-secondary)">
            {point.day.slice(5)}
          </text>
        {/if}
      {/each}
    </svg>
  </div>
{/snippet}

<PageLayout title="Отчёты и статистика">
  <!-- Tooltip -->
  {#if tooltip.visible}
    <div
      class="chart-tooltip"
      style="left:{tooltip.x}px; top:{tooltip.y}px"
    >{tooltip.text}</div>
  {/if}

  <!-- Filter bar -->
  <div class="report-filters">
    <div class="filter-field">
      <label class="filter-label" for="report-filter-car">Автомобиль</label>
      <select
        id="report-filter-car"
        class="field-select-sm"
        bind:value={filterCarId}
      >
        <option value="">Все автомобили</option>
        {#each carsList as car}
          <option value={car.id}>{car.brand} {car.model}{car.plate ? ` (${car.plate})` : ''}</option>
        {/each}
      </select>
    </div>
  </div>

  <!-- Tabs nav -->
  <div class="tabs-nav" role="tablist">
    {#each ([
      { id: 'overview',     label: 'Общий' },
      { id: 'fuel',         label: 'Топливо' },
      { id: 'maintenance',  label: 'ТО' },
      { id: 'fines',        label: 'Штрафы' },
    ] as const) as tab}
      <button
        role="tab"
        class="tab-btn spotlight"
        class:active={activeTab === tab.id}
        aria-selected={activeTab === tab.id}
        onclick={() => activeTab = tab.id as Tab}
      >
        {tab.label}
      </button>
    {/each}
  </div>

  {#if loading}
    <div class="center-loader">
      <Loader size={48} />
    </div>
  {:else if error}
    <div class="error-banner" role="alert">
      <span class="error-icon">⚠</span>
      <span>{error}</span>
    </div>
  {:else}

    <!-- ── Tab: Overview ─────────────────────────────────────────────────────── -->
    {#if activeTab === 'overview'}
      <div class="overview-section">
        <!-- Summary cards -->
        <div class="summary-grid">
          <div class="summary-card" style="--card-color: var(--accent);">
            <div class="summary-label">Топливо</div>
            <div class="summary-value">{fmt(fuelTotal)}</div>
            <div class="summary-sub">{filteredFuelList.length} заправок</div>
          </div>
          <div class="summary-card" style="--card-color: var(--success);">
            <div class="summary-label">ТО</div>
            <div class="summary-value">{fmt(maintenanceTotal)}</div>
            <div class="summary-sub">{filteredMaintenanceList.length} записей</div>
          </div>
          <div class="summary-card" style="--card-color: var(--danger);">
            <div class="summary-label">Штрафы</div>
            <div class="summary-value">{fmt(finesTotal)}</div>
            <div class="summary-sub">{filteredFinesList.length} штрафов</div>
          </div>
        </div>

        <!-- Grand total -->
        <div class="total-row">
          <span class="total-label">Итого расходов:</span>
          <span class="total-value">{fmt(grandTotal)}</span>
        </div>

        <!-- SVG Pie chart -->
        {#if grandTotal > 0}
          {@const pieItems = [
            { label: 'Топливо',  value: fuelTotal,        color: '#0078d4' },
            { label: 'ТО',       value: maintenanceTotal, color: '#3fb950' },
            { label: 'Штрафы',   value: finesTotal,       color: '#f85149' },
          ]}
          {@const slices = pieItems.filter(s => s.value > 0)}
          {@const pieTotal = slices.reduce((s, sl) => s + sl.value, 0)}
          <div class="chart-section">
            <h3 class="section-title">Распределение расходов</h3>
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="pie-wrap" onmousemove={moveTooltip}>
              <svg viewBox="-1 -1 2 2" class="pie-svg" style="transform: rotate(-90deg)">
                {#each slices as slice, i}
                  {@const startAngle = slices.slice(0, i).reduce((s, sl) => s + sl.value / pieTotal * Math.PI * 2, 0)}
                  {@const endAngle = startAngle + slice.value / pieTotal * Math.PI * 2}
                  {@const x1 = Math.cos(startAngle)}
                  {@const y1 = Math.sin(startAngle)}
                  {@const x2 = Math.cos(endAngle)}
                  {@const y2 = Math.sin(endAngle)}
                  {@const largeArc = slice.value / pieTotal > 0.5 ? 1 : 0}
                  <!-- svelte-ignore a11y_no_static_element_interactions -->
                  <path
                    d="M 0 0 L {x1} {y1} A 1 1 0 {largeArc} 1 {x2} {y2} Z"
                    fill={slice.color}
                    opacity="0.9"
                    class="pie-slice"
                    onmouseenter={(e) => showTooltip(e, `${slice.label} — ${slice.value.toLocaleString('ru-RU')} ₽ (${(slice.value / pieTotal * 100).toFixed(0)}%)`)}
                    onmouseleave={hideTooltip}
                  />
                {/each}
              </svg>
              <div class="pie-legend">
                {#each pieItems as item}
                  {#if item.value > 0}
                    <div class="legend-item">
                      <span class="legend-dot" style="background:{item.color}"></span>
                      <span class="legend-label">{item.label}</span>
                      <span class="legend-value">{fmt(item.value)}</span>
                      <span class="legend-pct">({(item.value / grandTotal * 100).toFixed(0)}%)</span>
                    </div>
                  {/if}
                {/each}
              </div>
            </div>
          </div>
        {:else}
          <EmptyState message="Нет данных для отображения. Добавьте заправки, ТО или штрафы." />
        {/if}
      </div>
    {/if}

    <!-- ── Tab: Fuel ──────────────────────────────────────────────────────────── -->
    {#if activeTab === 'fuel'}
      <div class="tab-section">
        <div class="section-header">
          <h3 class="section-title">История заправок</h3>
          <span class="section-total">Итого: {fmt(fuelTotal)}</span>
        </div>
        <div class="view-controls">
          <div class="view-toggle">
            <button class="toggle-btn spotlight" class:active={fuelView === 'table'} onclick={() => fuelView = 'table'}>Таблица</button>
            <button class="toggle-btn spotlight" class:active={fuelView === 'chart'} onclick={() => fuelView = 'chart'}>График</button>
          </div>
          {#if fuelView === 'chart'}
            <div class="period-controls">
              <button class="toggle-btn spotlight" class:active={fuelPeriod === 'week'} onclick={() => fuelPeriod = 'week'}>Неделя</button>
              <button class="toggle-btn spotlight" class:active={fuelPeriod === 'month'} onclick={() => fuelPeriod = 'month'}>Месяц</button>
              <button class="toggle-btn spotlight" class:active={fuelPeriod === 'custom'} onclick={() => fuelPeriod = 'custom'}>Период</button>
              {#if fuelPeriod === 'custom'}
                <input type="date" class="field-input-sm" bind:value={fuelDateFrom} />
                <span>—</span>
                <input type="date" class="field-input-sm" bind:value={fuelDateTo} />
              {/if}
            </div>
          {/if}
        </div>
        {#if filteredFuelList.length === 0}
          <EmptyState message="Нет данных о заправках." />
        {:else if fuelView === 'table'}
          <div class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Дата</th>
                  <th>Тип топлива</th>
                  <th>Объём, л</th>
                  <th>Цена за л</th>
                  <th>Сумма</th>
                </tr>
              </thead>
              <tbody>
                {#each filteredFuelList as e}
                  <tr>
                    <td>{fmtDate(e.date)}</td>
                    <td>{FUEL_TYPE_LABELS[e.fuelType ?? ''] ?? e.fuelType ?? '—'}</td>
                    <td>{Number(e.liters).toLocaleString('ru-RU')}</td>
                    <td>{fmt(Number(e.price))}</td>
                    <td class="amount-cell">{fmt(fuelCost(e))}</td>
                  </tr>
                {/each}
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="4" class="foot-label">Итого</td>
                  <td class="amount-cell foot-total">{fmt(fuelTotal)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        {:else}
          {@render barChart(fuelChartData, '#0078d4')}
        {/if}
      </div>
    {/if}

    <!-- ── Tab: Maintenance ───────────────────────────────────────────────────── -->
    {#if activeTab === 'maintenance'}
      <div class="tab-section">
        <div class="section-header">
          <h3 class="section-title">История техобслуживания</h3>
          <span class="section-total">Итого: {fmt(maintenanceTotal)}</span>
        </div>
        <div class="view-controls">
          <div class="view-toggle">
            <button class="toggle-btn spotlight" class:active={maintenanceView === 'table'} onclick={() => maintenanceView = 'table'}>Таблица</button>
            <button class="toggle-btn spotlight" class:active={maintenanceView === 'chart'} onclick={() => maintenanceView = 'chart'}>График</button>
          </div>
          {#if maintenanceView === 'chart'}
            <div class="period-controls">
              <button class="toggle-btn spotlight" class:active={maintenancePeriod === 'week'} onclick={() => maintenancePeriod = 'week'}>Неделя</button>
              <button class="toggle-btn spotlight" class:active={maintenancePeriod === 'month'} onclick={() => maintenancePeriod = 'month'}>Месяц</button>
              <button class="toggle-btn spotlight" class:active={maintenancePeriod === 'custom'} onclick={() => maintenancePeriod = 'custom'}>Период</button>
              {#if maintenancePeriod === 'custom'}
                <input type="date" class="field-input-sm" bind:value={maintenanceDateFrom} />
                <span>—</span>
                <input type="date" class="field-input-sm" bind:value={maintenanceDateTo} />
              {/if}
            </div>
          {/if}
        </div>
        {#if filteredMaintenanceList.length === 0}
          <EmptyState message="Нет данных о техобслуживании." />
        {:else if maintenanceView === 'table'}
          <div class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Дата</th>
                  <th>Автомобиль</th>
                  <th>Тип работ</th>
                  <th>Стоимость</th>
                </tr>
              </thead>
              <tbody>
                {#each filteredMaintenanceList as e}
                  <tr>
                    <td>{fmtDate(e.date)}</td>
                    <td>{getCarLabel(String(e.carId ?? ''))}</td>
                    <td>{e.type ?? '—'}</td>
                    <td class="amount-cell">{fmt(Number(e.cost))}</td>
                  </tr>
                {/each}
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="3" class="foot-label">Итого</td>
                  <td class="amount-cell foot-total">{fmt(maintenanceTotal)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        {:else}
          {@render barChart(maintenanceChartData, '#3fb950')}
        {/if}
      </div>
    {/if}

    <!-- ── Tab: Fines ─────────────────────────────────────────────────────────── -->
    {#if activeTab === 'fines'}
      <div class="tab-section">
        <div class="section-header">
          <h3 class="section-title">Штрафы</h3>
          <span class="section-total">Итого: {fmt(finesTotal)}</span>
        </div>
        <div class="view-controls">
          <div class="view-toggle">
            <button class="toggle-btn spotlight" class:active={finesView === 'table'} onclick={() => finesView = 'table'}>Таблица</button>
            <button class="toggle-btn spotlight" class:active={finesView === 'chart'} onclick={() => finesView = 'chart'}>График</button>
          </div>
          {#if finesView === 'chart'}
            <div class="period-controls">
              <button class="toggle-btn spotlight" class:active={finesPeriod === 'week'} onclick={() => finesPeriod = 'week'}>Неделя</button>
              <button class="toggle-btn spotlight" class:active={finesPeriod === 'month'} onclick={() => finesPeriod = 'month'}>Месяц</button>
              <button class="toggle-btn spotlight" class:active={finesPeriod === 'custom'} onclick={() => finesPeriod = 'custom'}>Период</button>
              {#if finesPeriod === 'custom'}
                <input type="date" class="field-input-sm" bind:value={finesDateFrom} />
                <span>—</span>
                <input type="date" class="field-input-sm" bind:value={finesDateTo} />
              {/if}
            </div>
          {/if}
        </div>
        {#if filteredFinesList.length === 0}
          <EmptyState message="Нет данных о штрафах." />
        {:else if finesView === 'table'}
          <div class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Дата</th>
                  <th>Автомобиль</th>
                  <th>Нарушение</th>
                  <th>Статус</th>
                  <th>Сумма</th>
                </tr>
              </thead>
              <tbody>
                {#each filteredFinesList as e}
                  <tr>
                    <td>{fmtDate(e.date)}</td>
                    <td>{getCarLabel(String(e.carId ?? ''))}</td>
                    <td>{e.description || e.type || '—'}</td>
                    <td>
                      <span class="badge" class:badge-paid={e.status === 'paid'} class:badge-unpaid={e.status !== 'paid'}>
                        {e.status === 'paid' ? 'Оплачен' : 'Не оплачен'}
                      </span>
                    </td>
                    <td class="amount-cell">{fmt(Number(e.amount))}</td>
                  </tr>
                {/each}
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="4" class="foot-label">Итого</td>
                  <td class="amount-cell foot-total">{fmt(finesTotal)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        {:else}
          {@render barChart(finesChartData, '#f85149')}
        {/if}
      </div>
    {/if}

  {/if}
</PageLayout>

<style>
/* ── Tooltip ─────────────────────────────────────────────────────────────── */
.chart-tooltip {
  position: fixed;
  background: var(--bg-overlay, var(--bg-layer));
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 6px 10px;
  font-size: 0.8125rem;
  color: var(--text-primary);
  pointer-events: none;
  z-index: 1000;
  white-space: nowrap;
  box-shadow: var(--shadow-md, 0 4px 16px rgba(0,0,0,0.3));
}

/* ── Report filters ──────────────────────────────────────────────────────── */
.report-filters {
  display: flex;
  gap: 0.75rem;
  align-items: flex-end;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}
.filter-field { display: flex; flex-direction: column; gap: 0.375rem; }
.filter-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }
.field-select-sm {
  padding: 0.4375rem 0.75rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text-primary);
  font-size: 0.875rem;
  font-family: var(--font);
  outline: none;
  cursor: pointer;
  transition: border-color var(--transition), box-shadow var(--transition);
}
.field-select-sm:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
}

/* ── Tabs nav ────────────────────────────────────────────────────────────── */
.tabs-nav {
  display: flex;
  gap: 0.25rem;
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 0.375rem;
  box-shadow: var(--shadow-sm);
}

.tab-btn {
  flex: 1;
  padding: 0.625rem 1rem;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--text-secondary);
  font-size: 0.9375rem;
  font-family: var(--font);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition), color var(--transition);
  white-space: nowrap;
}

.tab-btn:hover:not(.active) {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-primary);
}

.tab-btn.active {
  background: rgba(0, 120, 212, 0.15);
  border: 1px solid rgba(0, 120, 212, 0.4);
  color: var(--accent);
  font-weight: 600;
}

/* ── Summary cards ───────────────────────────────────────────────────────── */
.overview-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

@media (max-width: 700px) {
  .summary-grid { grid-template-columns: 1fr; }
}

.summary-card {
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-left: 3px solid var(--card-color);
  border-radius: var(--radius-lg);
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition), box-shadow var(--transition);
}

.summary-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.summary-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-secondary);
}

.summary-value {
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--card-color);
}

.summary-sub {
  font-size: 0.8125rem;
  color: var(--text-secondary);
}

/* ── Grand total ─────────────────────────────────────────────────────────── */
.total-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 1rem 1.5rem;
  box-shadow: var(--shadow-sm);
}

.total-label {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.total-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
}

/* ── Chart section ───────────────────────────────────────────────────────── */
.chart-section {
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  box-shadow: var(--shadow-sm);
}

.section-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 1.25rem;
}

/* ── Pie chart ───────────────────────────────────────────────────────────── */
.pie-wrap {
  display: flex;
  align-items: center;
  gap: 2rem;
  flex-wrap: wrap;
}
.pie-svg {
  width: 200px;
  height: 200px;
  flex-shrink: 0;
  border-radius: 50%;
}
.pie-slice {
  transition: opacity 0.2s;
  cursor: pointer;
}
.pie-slice:hover { opacity: 1 !important; filter: brightness(1.15); }
.pie-legend {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}
.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}
.legend-label { color: var(--text-secondary); min-width: 60px; }
.legend-value { font-weight: 600; color: var(--text-primary); }
.legend-pct { font-size: 0.8rem; color: var(--text-secondary); }

/* ── Bar chart ───────────────────────────────────────────────────────────── */
.bar-chart-wrap { overflow-x: auto; }
.bar-svg { width: 100%; min-width: 300px; max-width: 700px; display: block; }

/* ── View controls ───────────────────────────────────────────────────────── */
.view-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}
.view-toggle, .period-controls {
  display: flex;
  gap: 0.25rem;
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 0.25rem;
}
.toggle-btn {
  padding: 0.375rem 0.875rem;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-family: var(--font);
  cursor: pointer;
  transition: background var(--transition), color var(--transition);
}
.toggle-btn.active {
  background: rgba(0, 120, 212, 0.15);
  border: 1px solid rgba(0, 120, 212, 0.4);
  color: var(--accent);
  font-weight: 600;
}
.field-input-sm {
  padding: 0.375rem 0.625rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text-primary);
  font-size: 0.875rem;
  font-family: var(--font);
  outline: none;
}

/* ── Tab sections ────────────────────────────────────────────────────────── */
.tab-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.section-header .section-title {
  margin-bottom: 0;
}

.section-total {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-primary);
}

/* ── Data table ──────────────────────────────────────────────────────────── */
.table-wrap {
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9375rem;
}

.data-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-secondary);
  background: var(--bg-input);
  border-bottom: 1px solid var(--border);
}

.data-table td {
  padding: 0.75rem 1rem;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border);
}

.data-table tbody tr:last-child td {
  border-bottom: none;
}

.data-table tbody tr:hover td {
  background: var(--bg-input);
}

.amount-cell {
  text-align: right;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.data-table tfoot td {
  border-top: 1px solid var(--border);
  border-bottom: none;
  background: var(--bg-input);
  padding: 0.75rem 1rem;
}

.foot-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.foot-total {
  font-size: 1rem;
  color: var(--text-primary);
}

/* ── Badges ──────────────────────────────────────────────────────────────── */
.badge {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  border-radius: 100px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-paid {
  background: var(--success-light);
  color: var(--success);
}

.badge-unpaid {
  background: var(--danger-light);
  color: var(--danger);
}

/* ── States ──────────────────────────────────────────────────────────────── */
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
}

.error-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}
</style>
