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
    finesLoading, finesError,
    filteredFinesList, loadFines, createFine, editFine, removeFine
  } from '../../stores/fines';
  import type { Fine } from '../../lib/types';
  import { fetchCars } from '../../lib/api';
  import type { Car } from '../../lib/types';

  const COLUMNS = [
    { label: 'Дата', key: '_date' },
    { label: 'Автомобиль', key: '_car' },
    { label: 'Описание', key: 'description' },
    { label: 'Сумма (₽)', key: '_amount' },
  ];

  let cars = $state<Car[]>([]);

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
    if (!carId) return '—';
    const car = cars.find(c => c.id === carId);
    if (!car) return carId;
    return `${car.brand} ${car.model} (${car.plate})`;
  }

  let rows = $derived(
    $filteredFinesList.map(f => ({
      ...f,
      _date: formatDate(f.date),
      _car: getCarLabel(String(f.carId ?? '')),
      _amount: Number(f.amount).toLocaleString('ru-RU', { minimumFractionDigits: 2 }) + ' ₽',
      _statusBadge: f.status === 'paid' ? 'paid' : 'unpaid',
      _dateRaw: f.date ?? '',
      _amountNum: Number(f.amount),
    }))
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

  function applySorting<T extends Record<string, unknown>>(arr: T[]): T[] {
    if (sortState.length === 0) return arr;
    return [...arr].sort((a, b) => {
      for (const { key, dir } of sortState) {
        const av = a[key] ?? '';
        const bv = b[key] ?? '';
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

  // Visible sort state for Table
  let tableSortState = $derived(
    sortState.map(s => ({ ...s, key: s.key === '_dateRaw' ? '_date' : s.key === '_amountNum' ? '_amount' : s.key }))
  );

  function handleTableSort(key: string, event: MouseEvent) {
    const internalKey = key === '_date' ? '_dateRaw' : key === '_amount' ? '_amountNum' : key;
    handleSort(internalKey, event);
  }

  const SORT_KEYS = ['_date', '_car', 'description', '_amount'];

  let showModal = $state(false);
  let editingId = $state<string | null>(null);
  let confirmDeleteId = $state<string | null>(null);
  let toast = $state({ open: false, message: '', type: 'info' as 'info' | 'success' | 'error' });

  type FineForm = {
    date: string;
    description: string;
    amount: string;
    status: 'paid' | 'unpaid';
    carId: string;
  };

  let form = $state<FineForm>({
    date: '',
    description: '',
    amount: '',
    status: 'unpaid',
    carId: '',
  });

  let filterStatus = $state('');
  let filterSearch = $state('');
  let filterCarId = $state('');

  // Pagination state
  let page = $state(1);
  let perPage = $state(5);
  const PER_PAGE_OPTIONS = [5, 10, 25];

  let displayRows = $derived(
    rows.filter(r => {
      const matchStatus = !filterStatus || r._statusBadge === filterStatus;
      const matchCar = !filterCarId || String(r.carId ?? '') === filterCarId;
      const search = filterSearch.toLowerCase();
      const matchSearch = !search ||
        String(r.description ?? '').toLowerCase().includes(search) ||
        String(r._car ?? '').toLowerCase().includes(search);
      return matchStatus && matchCar && matchSearch;
    })
  );

  // Pagination derived values
  let sortedRows = $derived(applySorting(displayRows));
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

  function setPage(p: number) {
    page = Math.max(1, Math.min(p, totalPages));
  }

  function setPerPage(value: number) {
    perPage = value;
    page = 1;
  }

  // Reset page when filters change
  $effect(() => {
    void filterStatus;
    void filterSearch;
    void filterCarId;
    page = 1;
  });

  function openAdd() {
    editingId = null;
    form = {
      date: new Date().toISOString().split('T')[0],
      description: '',
      amount: '',
      status: 'unpaid',
      carId: cars[0]?.id ?? '',
    };
    showModal = true;
  }

  function openEdit(row: Record<string, unknown>) {
    editingId = String(row.id ?? '');
    form = {
      date: String(row.date ?? ''),
      description: String(row.description ?? ''),
      amount: String(row.amount ?? ''),
      status: row.status === 'paid' ? 'paid' : 'unpaid',
      carId: String(row.carId ?? ''),
    };
    showModal = true;
  }

  async function handleSave() {
    if (!form.description.trim()) {
      toast = { open: true, message: 'Укажите описание нарушения', type: 'error' };
      return;
    }
    if (!form.amount || Number(form.amount) <= 0) {
      toast = { open: true, message: 'Укажите сумму штрафа', type: 'error' };
      return;
    }
    try {
      const payload: Omit<Fine, 'id'> = {
        ...form,
        amount: Number(form.amount),
      };
      if (editingId) {
        await editFine(editingId, payload);
        toast = { open: true, message: 'Штраф обновлён', type: 'success' };
      } else {
        await createFine(payload);
        toast = { open: true, message: 'Штраф добавлен', type: 'success' };
      }
      showModal = false;
    } catch {
      toast = { open: true, message: $finesError ?? 'Ошибка сохранения', type: 'error' };
    }
  }

  async function handleDelete(id: string) {
    try {
      await removeFine(id);
      toast = { open: true, message: 'Штраф удалён', type: 'success' };
    } catch {
      toast = { open: true, message: $finesError ?? 'Ошибка удаления', type: 'error' };
    } finally {
      confirmDeleteId = null;
    }
  }

  onMount(async () => {
    await ensureAuthenticated();
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch {
      cars = [];
    }
    loadFines();
  });
</script>

<PageLayout title="Штрафы">
  {#snippet toolbar()}
    <div class="filters">
      <div class="filter-field">
        <label class="filter-label" for="fines-search">Поиск</label>
        <Input
          id="fines-search"
          placeholder="Поиск по описанию или авто..."
          bind:value={filterSearch}
        />
      </div>
      <div class="filter-field">
        <label class="filter-label" for="fines-filter-car">Автомобиль</label>
        <select
          id="fines-filter-car"
          class="field-select"
          bind:value={filterCarId}
        >
          <option value="">Все автомобили</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
      <div class="filter-field">
        <label class="filter-label" for="fines-filter-status">Статус</label>
        <select
          id="fines-filter-status"
          class="field-select"
          bind:value={filterStatus}
        >
          <option value="">Все статусы</option>
          <option value="unpaid">Не оплачен</option>
          <option value="paid">Оплачен</option>
        </select>
      </div>
    </div>
    <Button variant="primary" onclick={openAdd}>+ Добавить</Button>
  {/snippet}

  <Table
    columns={COLUMNS}
    rows={pagedRows}
    loading={$finesLoading}
    error={$finesError ?? ''}
    emptyText="Нет записей о штрафах"
    onRowClick={openEdit}
    sortKeys={SORT_KEYS}
    sort={tableSortState}
    onSort={handleTableSort}
  >
    {#snippet actions(row)}
      <div class="row-actions">
        <button
          class="badge-btn"
          class:badge-paid={row._statusBadge === 'paid'}
          class:badge-unpaid={row._statusBadge !== 'paid'}
          onclick={(e) => {
            e.stopPropagation();
            const id = String(row.id ?? '');
            if (id) {
              const newStatus = row._statusBadge === 'paid' ? 'unpaid' : 'paid';
              editFine(id, {
                carId: String(row.carId ?? ''),
                amount: Number(row.amount),
                type: String(row.type ?? ''),
                date: String(row.date ?? ''),
                description: String(row.description ?? ''),
                status: newStatus,
              })
                .then(() => toast = { open: true, message: 'Статус обновлён', type: 'success' })
                .catch(() => toast = { open: true, message: 'Ошибка', type: 'error' });
            }
          }}
          title="Переключить статус"
        >
          {row._statusBadge === 'paid' ? '✓ Оплачен' : '✗ Не оплачен'}
        </button>
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
    title={editingId ? 'Редактировать штраф' : 'Добавить штраф'}
    onClose={() => showModal = false}
    width="520px"
  >
    <div class="form-grid">
      <div class="field">
        <label class="field-label" for="fine-date">Дата *</label>
        <input id="fine-date" class="field-input" type="date" bind:value={form.date} required />
      </div>
      <div class="field">
        <label class="field-label" for="fine-car">Автомобиль</label>
        <select id="fine-car" class="field-select" bind:value={form.carId}>
          <option value="">Без привязки к авто</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
      <div class="col-span-2">
        <Input
          label="Описание нарушения *"
          placeholder="Превышение скорости, парковка..."
          bind:value={form.description}
          required
        />
      </div>
      <Input
        label="Сумма (₽) *"
        type="number"
        min="0"
        step="1"
        placeholder="500"
        bind:value={form.amount}
        required
      />
      <div class="field">
        <label class="field-label" for="fine-status">Статус *</label>
        <select id="fine-status" class="field-select" bind:value={form.status}>
          <option value="unpaid">Не оплачен</option>
          <option value="paid">Оплачен</option>
        </select>
      </div>
    </div>
    {#snippet footer()}
      <Button variant="secondary" onclick={() => showModal = false}>Отмена</Button>
      <Button variant="primary" loading={$finesLoading} onclick={handleSave}>Сохранить</Button>
    {/snippet}
  </Modal>

  <Modal open={!!confirmDeleteId} title="Подтвердите удаление" onClose={() => confirmDeleteId = null}>
    <p class="confirm-text">Удалить запись о штрафе? Это действие нельзя отменить.</p>
    {#snippet footer()}
      <Button variant="secondary" onclick={() => confirmDeleteId = null}>Отмена</Button>
      <Button variant="danger" onclick={() => confirmDeleteId && handleDelete(confirmDeleteId)}>Удалить</Button>
    {/snippet}
  </Modal>

  <Toast open={toast.open} message={toast.message} type={toast.type} />
</PageLayout>

<style>
.filters { display: flex; gap: 0.75rem; flex: 1; flex-wrap: wrap; align-items: flex-end; }

.filter-field { display: flex; flex-direction: column; gap: 0.375rem; }
.filter-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  min-width: 380px;
}
.col-span-2 { grid-column: span 2; }

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

/* Цветные badge-кнопки статуса */
.badge-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.625rem;
  border-radius: 100px;
  font-size: 0.8125rem;
  font-weight: 600;
  font-family: var(--font);
  cursor: pointer;
  border: 1px solid transparent;
  transition: opacity var(--transition), transform var(--transition);
  white-space: nowrap;
}
.badge-btn:hover { opacity: 0.8; transform: scale(0.97); }

.badge-paid {
  background: var(--success-light);
  color: var(--success);
  border-color: rgba(63, 185, 80, 0.3);
}
.badge-unpaid {
  background: var(--danger-light);
  color: var(--danger);
  border-color: rgba(248, 81, 73, 0.3);
}

.row-actions { display: flex; gap: 0.375rem; align-items: center; }
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
