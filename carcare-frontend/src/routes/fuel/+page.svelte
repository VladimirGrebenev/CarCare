<script lang="ts">
  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Table from '../../components/ui/Table.svelte';
  import Modal from '../../components/ui/Modal.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Input from '../../components/ui/Input.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import {
    fuelLoading, fuelError, fuelFilters,
    filteredFuelList, loadFuel, createFuel, editFuel, removeFuel
  } from '../../stores/fuel';
  import type { Fuel } from '../../stores/fuel';
  import { fetchCars } from '../../lib/api';
  import type { Car } from '../../lib/types';

  const FUEL_TYPES = [
    { value: 'petrol', label: 'Бензин' },
    { value: 'diesel', label: 'Дизель' },
    { value: 'gas', label: 'Газ' },
    { value: 'electric', label: 'Электро' },
  ];

  const COLUMNS = [
    { label: 'Дата', key: '_date' },
    { label: 'Автомобиль', key: '_car' },
    { label: 'Тип топлива', key: '_fuelType' },
    { label: 'Объём', key: 'liters' },
    { label: 'Цена/л', key: '_priceFormatted' },
    { label: 'Сумма', key: '_total' },
  ];

  let cars = $state<Car[]>([]);
  let _carsLoading = $state(false);

  let showModal = $state(false);
  let editingId = $state<string | null>(null);
  let confirmDeleteId = $state<string | null>(null);
  let toast = $state({ open: false, message: '', type: 'info' as 'info' | 'success' | 'error' });

  // Pagination state
  let page = $state(1);
  let perPage = $state(5);
  const PER_PAGE_OPTIONS = [5, 10, 25];

  type FuelForm = {
    date: string;
    liters: string;
    price: string;
    carId: string;
    fuelType: string;
  };

  let form = $state<FuelForm>({
    date: '',
    liters: '',
    price: '',
    carId: '',
    fuelType: 'petrol',
  });

  function formatDate(raw: string | null | undefined): string {
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

  function getCarLabel(carId: string): string {
    const car = cars.find(c => c.id === carId);
    if (!car) return carId || '—';
    return `${car.brand} ${car.model} (${car.plate})`;
  }

  function getFuelTypeLabel(type: string): string {
    return FUEL_TYPES.find(t => t.value === type)?.label ?? type ?? '—';
  }

  let rows = $derived(
    $filteredFuelList.map(f => {
      const record = f as Record<string, unknown>;
      return {
        ...record,
        _date: formatDate(record.date as string | null | undefined),
        _car: getCarLabel(String(record.carId ?? '')),
        _fuelType: getFuelTypeLabel(String(record.fuelType ?? '')),
        _priceFormatted: Number(record.price).toLocaleString('ru-RU', { minimumFractionDigits: 2 }),
        _total: (Number(record.liters) * Number(record.price)).toLocaleString('ru-RU', { minimumFractionDigits: 2 }) + ' ₽',
        // numeric shadows for sorting
        _totalNum: Number(record.liters) * Number(record.price),
        _dateRaw: record.date as string ?? '',
      };
    })
  );

  // Sorting state
  let sortState = $state<Array<{ key: string; dir: 'asc' | 'desc' }>>([
    { key: '_dateRaw', dir: 'desc' }
  ]);

  function handleSort(key: string, _event: MouseEvent) {
    const existing = sortState.findIndex(s => s.key === key);
    if (existing === -1) {
      if (sortState.length < 2) {
        sortState = [...sortState, { key, dir: 'desc' as const }];
      } else {
        sortState = [sortState[0], { key, dir: 'desc' as const }];
      }
    } else if (sortState[existing].dir === 'desc') {
      sortState = sortState.map((s, i) => i === existing ? { ...s, dir: 'asc' as const } : s);
    } else {
      sortState = sortState.filter((_, i) => i !== existing);
    }
    page = 1;
  }

  function applySorting(arr: typeof rows): typeof rows {
    if (sortState.length === 0) return arr;
    return [...arr].sort((a, b) => {
      for (const { key, dir } of sortState) {
        const av = a[key as keyof typeof a] ?? '';
        const bv = b[key as keyof typeof b] ?? '';
        const an = Number(av), bn = Number(bv);
        let cmp = 0;
        if (!isNaN(an) && !isNaN(bn) && String(av) !== '' && String(bv) !== '') {
          cmp = an - bn;
        } else {
          cmp = String(av).localeCompare(String(bv), 'ru');
        }
        if (cmp !== 0) return dir === 'asc' ? cmp : -cmp;
      }
      return 0;
    });
  }

  // Pagination derived values
  let sortedRows = $derived(applySorting(rows));
  let totalPages = $derived(Math.ceil(sortedRows.length / perPage));
  let showPagination = $derived(sortedRows.length > 5);
  let pagedRows = $derived(sortedRows.slice((page - 1) * perPage, page * perPage));

  let pageNumbers = $derived((() => {
    if (totalPages <= 5) return Array.from({ length: totalPages }, (_, i) => i + 1);
    const nums: (number | '...')[] = [];
    if (page <= 3) {
      nums.push(1, 2, 3, 4, '...', totalPages);
    } else if (page >= totalPages - 2) {
      nums.push(1, '...', totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
    } else {
      nums.push(1, '...', page - 1, page, page + 1, '...', totalPages);
    }
    return nums;
  })());

  // Visible sort state for Table: map _dateRaw -> _date for display
  let tableSortState = $derived(
    sortState.map(s => ({ ...s, key: s.key === '_dateRaw' ? '_date' : s.key }))
  );

  function handleTableSort(key: string, event: MouseEvent) {
    // Map display key back to internal key
    const internalKey = key === '_date' ? '_dateRaw' : key;
    handleSort(internalKey, event);
  }

  function setPage(p: number) {
    page = Math.max(1, Math.min(p, totalPages));
  }

  function setPerPage(value: number) {
    perPage = value;
    page = 1;
  }

  function openAdd() {
    editingId = null;
    form = {
      date: new Date().toISOString().split('T')[0],
      liters: '',
      price: '',
      carId: cars[0]?.id ?? '',
      fuelType: 'petrol',
    };
    showModal = true;
  }

  function openEdit(row: Record<string, unknown>) {
    editingId = String(row.id ?? '');
    form = {
      date: String(row.date ?? ''),
      liters: String(row.liters ?? ''),
      price: String(row.price ?? ''),
      carId: String(row.carId ?? ''),
      fuelType: String(row.fuelType ?? 'petrol'),
    };
    showModal = true;
  }

  async function handleSave() {
    if (!form.carId) {
      toast = { open: true, message: 'Выберите автомобиль', type: 'error' };
      return;
    }
    if (!form.liters || Number(form.liters) <= 0) {
      toast = { open: true, message: 'Укажите объём топлива', type: 'error' };
      return;
    }
    if (!form.price || Number(form.price) <= 0) {
      toast = { open: true, message: 'Укажите цену за литр', type: 'error' };
      return;
    }
    try {
      const payload: Fuel & { fuelType: string } = {
        ...form,
        liters: Number(form.liters),
        price: Number(form.price),
      };
      if (editingId) {
        await editFuel(editingId, payload);
        toast = { open: true, message: 'Заправка обновлена', type: 'success' };
      } else {
        await createFuel(payload);
        toast = { open: true, message: 'Заправка добавлена', type: 'success' };
      }
      showModal = false;
    } catch {
      toast = { open: true, message: $fuelError ?? 'Ошибка сохранения', type: 'error' };
    }
  }

  async function handleDelete(id: string) {
    try {
      await removeFuel(id);
      toast = { open: true, message: 'Заправка удалена', type: 'success' };
    } catch {
      toast = { open: true, message: $fuelError ?? 'Ошибка удаления', type: 'error' };
    } finally {
      confirmDeleteId = null;
    }
  }

  async function loadCars() {
    _carsLoading = true;
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch {
      cars = [];
    } finally {
      _carsLoading = false;
    }
  }

  // Reset page when filters change
  $effect(() => {
    void $fuelFilters;
    page = 1;
  });

  const SORT_KEYS = ['_date', '_car', '_fuelType', 'liters', '_priceFormatted', '_total'];

  onMount(async () => {
    await ensureAuthenticated();
    await loadCars();
    loadFuel();
  });
</script>

<PageLayout title="Заправки">
  {#snippet toolbar()}
    <div class="filters">
      <div class="filter-field">
        <label class="filter-label" for="fuel-filter-car">Автомобиль</label>
        <select
          id="fuel-filter-car"
          class="field-select"
          value={$fuelFilters['carId'] ?? ''}
          onchange={(e) => fuelFilters.update(f => ({ ...f, carId: (e.target as HTMLSelectElement).value }))}
        >
          <option value="">Все автомобили</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
    </div>
    <Button variant="primary" onclick={openAdd}>+ Добавить</Button>
  {/snippet}

  <Table
    columns={COLUMNS}
    rows={pagedRows}
    loading={$fuelLoading}
    error={$fuelError ?? ''}
    emptyText="Нет записей о заправках"
    onRowClick={openEdit}
    sortKeys={SORT_KEYS}
    sort={tableSortState}
    onSort={handleTableSort}
    className="fuel-table"
  >
    {#snippet actions(row)}
      <div class="row-actions">
        <button class="icon-btn edit-btn" title="Редактировать" onclick={(e) => { e.stopPropagation(); openEdit(row); }}>
          <svg width="15" height="15" viewBox="0 0 16 16" fill="none"><path d="M11.013 1.427a1.75 1.75 0 012.474 2.474L4.81 12.578l-3.182.354.354-3.181 8.031-8.324z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
        <button class="icon-btn delete-btn" title="Удалить" onclick={(e) => { e.stopPropagation(); confirmDeleteId = row.id as string; }}>
          <svg width="15" height="15" viewBox="0 0 16 16" fill="none"><path d="M6 2h4M2 4h12M5 4l.5 8h5L11 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
      </div>
    {/snippet}
  </Table>

  {#if showPagination}
    <div class="pagination-bar">
      <span class="pagination-info">
        Показано {Math.min((page - 1) * perPage + 1, sortedRows.length)}–{Math.min(page * perPage, sortedRows.length)} из {sortedRows.length}
      </span>
      <div class="per-page-group">
        {#each PER_PAGE_OPTIONS as opt}
          <button
            class="page-btn"
            class:active={perPage === opt}
            onclick={() => setPerPage(opt)}
          >{opt}</button>
        {/each}
      </div>
      <div class="page-nav">
        <button class="page-btn nav-btn" disabled={page === 1} onclick={() => setPage(page - 1)}>←</button>
        {#each pageNumbers as p}
          {#if p === '...'}
            <span class="page-ellipsis">…</span>
          {:else}
            <button
              class="page-btn"
              class:active={page === p}
              onclick={() => setPage(p as number)}
            >{p}</button>
          {/if}
        {/each}
        <button class="page-btn nav-btn" disabled={page === totalPages} onclick={() => setPage(page + 1)}>→</button>
      </div>
    </div>
  {/if}

  <Modal
    open={showModal}
    title={editingId ? 'Редактировать заправку' : 'Добавить заправку'}
    onClose={() => showModal = false}
    width="520px"
  >
    <div class="form-grid">
      <div class="field">
        <label class="field-label" for="fuel-date">Дата *</label>
        <input id="fuel-date" class="field-input" type="date" bind:value={form.date} required />
      </div>
      <div class="field">
        <label class="field-label" for="fuel-car">Автомобиль *</label>
        <select id="fuel-car" class="field-select" bind:value={form.carId} required>
          <option value="">Выберите автомобиль</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
      <div class="field">
        <label class="field-label" for="fuel-type">Тип топлива *</label>
        <select id="fuel-type" class="field-select" bind:value={form.fuelType}>
          {#each FUEL_TYPES as ft}
            <option value={ft.value}>{ft.label}</option>
          {/each}
        </select>
      </div>
      <Input
        label="Объём (л) *"
        type="number"
        min="0"
        step="0.01"
        placeholder="45.5"
        bind:value={form.liters}
        required
      />
      <Input
        label="Цена за литр (₽) *"
        type="number"
        min="0"
        step="0.01"
        placeholder="58.40"
        bind:value={form.price}
        required
      />
      {#if form.liters && form.price}
        <div class="total-preview">
          <span class="total-label">Итого:</span>
          <span class="total-value">
            {(Number(form.liters) * Number(form.price)).toLocaleString('ru-RU', { minimumFractionDigits: 2 })} ₽
          </span>
        </div>
      {/if}
    </div>
    {#snippet footer()}
      <Button variant="secondary" onclick={() => showModal = false}>Отмена</Button>
      <Button variant="primary" loading={$fuelLoading} onclick={handleSave}>Сохранить</Button>
    {/snippet}
  </Modal>

  <Modal open={!!confirmDeleteId} title="Подтвердите удаление" onClose={() => confirmDeleteId = null}>
    <p class="confirm-text">Удалить запись о заправке? Это действие нельзя отменить.</p>
    {#snippet footer()}
      <Button variant="secondary" onclick={() => confirmDeleteId = null}>Отмена</Button>
      <Button variant="danger" onclick={() => confirmDeleteId && handleDelete(confirmDeleteId)}>Удалить</Button>
    {/snippet}
  </Modal>

  <Toast open={toast.open} message={toast.message} type={toast.type} />
</PageLayout>

<style>
/* Таблица заправок — компактные столбцы */
:global(.fuel-table table) { table-layout: fixed; }
:global(.fuel-table th:nth-child(1)) { width: 110px; } /* дата */
:global(.fuel-table th:nth-child(2)) { width: auto; }   /* авто */
:global(.fuel-table th:nth-child(3)) { width: 80px; }   /* тип */
:global(.fuel-table th:nth-child(4)) { width: 70px; }   /* объём */
:global(.fuel-table th:nth-child(5)) { width: 75px; }   /* цена */
:global(.fuel-table th:nth-child(6)) { width: 90px; }   /* сумма */
:global(.fuel-table th:nth-child(7)) { width: 80px; }   /* действия */

.filters { display: flex; gap: 0.75rem; flex: 1; flex-wrap: wrap; }

.filter-field { display: flex; flex-direction: column; gap: 0.375rem; }
.filter-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.field { display: flex; flex-direction: column; gap: 0.375rem; }
.field-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }

.field-select,
.field-input {
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
.field-select:focus,
.field-input:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
}

.total-preview {
  grid-column: span 2;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.5rem;
  padding: 0.625rem 0.875rem;
  background: var(--accent-light);
  border-radius: var(--radius-md);
  border: 1px solid rgba(0, 120, 212, 0.25);
}
.total-label { font-size: 0.875rem; color: var(--text-secondary); }
.total-value { font-size: 1.125rem; font-weight: 700; color: var(--accent-text); }

.row-actions { display: flex; gap: 0.25rem; align-items: center; }
.confirm-text { color: var(--text-secondary); line-height: 1.6; }

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg-input);
  cursor: pointer;
  color: var(--text-secondary);
  transition: background var(--transition), color var(--transition), border-color var(--transition);
}
.edit-btn:hover { background: var(--accent-light); color: var(--accent-text); border-color: rgba(0,120,212,0.4); }
.delete-btn:hover { background: var(--danger-light); color: var(--danger); border-color: var(--danger); }

/* Pagination */
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 1rem;
  flex-wrap: wrap;
}

.pagination-info {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  white-space: nowrap;
}

.per-page-group,
.page-nav {
  display: flex;
  gap: 0.25rem;
  align-items: center;
}

.page-btn {
  min-width: 2rem;
  height: 2rem;
  padding: 0 0.5rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-family: var(--font);
  cursor: pointer;
  transition: background var(--transition), color var(--transition), border-color var(--transition);
}
.page-btn:hover:not(:disabled):not(.active) {
  background: var(--bg-layer);
  color: var(--text-primary);
}
.page-btn.active {
  background: var(--accent-light);
  color: var(--accent-text);
  border-color: rgba(0, 120, 212, 0.4);
  font-weight: 600;
}
.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.nav-btn { font-size: 1rem; }
.page-ellipsis {
  padding: 0 0.25rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
  line-height: 2rem;
}
</style>
