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
    { label: 'Объём (л)', key: 'liters' },
    { label: 'Цена/л (₽)', key: '_priceFormatted' },
    { label: 'Сумма (₽)', key: '_total' },
  ];

  let cars = $state<Car[]>([]);
  let carsLoading = $state(false);

  let showModal = $state(false);
  let editingId = $state<string | null>(null);
  let confirmDeleteId = $state<string | null>(null);
  let toast = $state({ open: false, message: '', type: 'info' as 'info' | 'success' | 'error' });

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
      };
    })
  );

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
    carsLoading = true;
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch {
      cars = [];
    } finally {
      carsLoading = false;
    }
  }

  onMount(async () => {
    await ensureAuthenticated();
    await loadCars();
    loadFuel();
  });
</script>

<PageLayout title="Заправки">
  <div class="page-toolbar">
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
  </div>

  <Table
    columns={COLUMNS}
    rows={rows}
    loading={$fuelLoading}
    error={$fuelError ?? ''}
    emptyText="Нет записей о заправках"
    onRowClick={openEdit}
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
.page-toolbar {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}
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
</style>
