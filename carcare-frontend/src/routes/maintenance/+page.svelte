<script lang="ts">
  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Table from '../../components/ui/Table.svelte';
  import Modal from '../../components/ui/Modal.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Input from '../../components/ui/Input.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import FAB from '../../components/ui/FAB.svelte';
  import { writable, get } from 'svelte/store';
  import {
    fetchMaintenanceHistory,
    addMaintenance,
    updateMaintenance,
    deleteMaintenance,
    fetchCars,
  } from '../../lib/api';
  import type { Car } from '../../lib/types';

  const MAINTENANCE_TYPES = [
    { value: 'oil', label: 'Замена масла' },
    { value: 'tires', label: 'Шиномонтаж / смена резины' },
    { value: 'brakes', label: 'Тормозная система' },
    { value: 'filters', label: 'Замена фильтров' },
    { value: 'battery', label: 'Аккумулятор' },
    { value: 'inspection', label: 'Техосмотр' },
    { value: 'other', label: 'Другое' },
  ];

  const COLUMNS = [
    { label: 'Дата', key: '_date' },
    { label: 'Автомобиль', key: '_car' },
    { label: 'Тип обслуживания', key: '_typeLabel' },
    { label: 'Стоимость (₽)', key: '_cost' },
    { label: 'Описание', key: 'description' },
  ];

  type MaintenanceRecord = {
    id: string;
    date: string;
    carId: string;
    type: string;
    cost: number;
    description?: string;
    [key: string]: unknown;
  };

  let items = $state<MaintenanceRecord[]>([]);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let cars = $state<Car[]>([]);

  let filterType = $state('');
  let filterCarId = $state('');

  function getCarLabel(carId: string): string {
    if (!carId) return '—';
    const car = cars.find(c => c.id === carId);
    if (!car) return carId;
    return `${car.brand} ${car.model} (${car.plate})`;
  }

  function getTypeLabel(type: string): string {
    return MAINTENANCE_TYPES.find(t => t.value === type)?.label ?? type ?? '—';
  }

  let rows = $derived(
    items
      .filter(item => {
        const matchType = !filterType || item.type === filterType;
        const matchCar = !filterCarId || item.carId === filterCarId;
        return matchType && matchCar;
      })
      .map(item => ({
        ...item,
        _date: item.date ?? '—',
        _car: getCarLabel(String(item.carId ?? '')),
        _typeLabel: getTypeLabel(String(item.type ?? '')),
        _cost: Number(item.cost).toLocaleString('ru-RU', { minimumFractionDigits: 2 }) + ' ₽',
      }))
  );

  let showModal = $state(false);
  let editingId = $state<string | null>(null);
  let confirmDeleteId = $state<string | null>(null);
  let saving = $state(false);
  let toast = $state({ open: false, message: '', type: 'info' as 'info' | 'success' | 'error' });

  type MaintenanceForm = {
    date: string;
    carId: string;
    type: string;
    cost: string;
    description: string;
  };

  let form = $state<MaintenanceForm>({
    date: '',
    carId: '',
    type: 'oil',
    cost: '',
    description: '',
  });

  async function loadItems() {
    loading = true;
    error = null;
    try {
      const data = await fetchMaintenanceHistory({});
      items = Array.isArray(data) ? data : [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Ошибка загрузки';
      items = [];
    } finally {
      loading = false;
    }
  }

  function openAdd() {
    editingId = null;
    form = {
      date: new Date().toISOString().split('T')[0],
      carId: cars[0]?.id ?? '',
      type: 'oil',
      cost: '',
      description: '',
    };
    showModal = true;
  }

  function openEdit(row: Record<string, unknown>) {
    editingId = String(row.id ?? '');
    form = {
      date: String(row.date ?? ''),
      carId: String(row.carId ?? ''),
      type: String(row.type ?? 'other'),
      cost: String(row.cost ?? ''),
      description: String(row.description ?? ''),
    };
    showModal = true;
  }

  async function handleSave() {
    if (!form.carId) {
      toast = { open: true, message: 'Выберите автомобиль', type: 'error' };
      return;
    }
    if (!form.cost || Number(form.cost) <= 0) {
      toast = { open: true, message: 'Укажите стоимость', type: 'error' };
      return;
    }
    saving = true;
    try {
      const payload = { ...form, cost: Number(form.cost) };
      if (editingId) {
        const updated = await updateMaintenance(editingId, payload);
        items = items.map(item => item.id === editingId ? { ...item, ...updated } : item);
        toast = { open: true, message: 'Запись обновлена', type: 'success' };
      } else {
        const created = await addMaintenance(payload);
        items = [created, ...items];
        toast = { open: true, message: 'Запись добавлена', type: 'success' };
      }
      showModal = false;
    } catch (e: unknown) {
      toast = {
        open: true,
        message: e instanceof Error ? e.message : 'Ошибка сохранения',
        type: 'error',
      };
    } finally {
      saving = false;
    }
  }

  async function handleDelete(id: string) {
    try {
      await deleteMaintenance(id);
      items = items.filter(item => item.id !== id);
      toast = { open: true, message: 'Запись удалена', type: 'success' };
    } catch (e: unknown) {
      toast = {
        open: true,
        message: e instanceof Error ? e.message : 'Ошибка удаления',
        type: 'error',
      };
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
    await loadItems();
  });
</script>

<PageLayout title="Техобслуживание">
  <div class="page-toolbar">
    <div class="filters">
      <div class="filter-field">
        <label class="filter-label" for="maint-filter-type">Тип обслуживания</label>
        <select id="maint-filter-type" class="field-select" bind:value={filterType}>
          <option value="">Все типы</option>
          {#each MAINTENANCE_TYPES as mt}
            <option value={mt.value}>{mt.label}</option>
          {/each}
        </select>
      </div>
      <div class="filter-field">
        <label class="filter-label" for="maint-filter-car">Автомобиль</label>
        <select id="maint-filter-car" class="field-select" bind:value={filterCarId}>
          <option value="">Все автомобили</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
    </div>
    <Button variant="primary" onclick={openAdd}>+ Добавить ТО</Button>
  </div>

  <Table
    columns={COLUMNS}
    rows={rows}
    loading={loading}
    error={error ?? ''}
    emptyText="Нет записей о техобслуживании"
    onRowClick={openEdit}
  >
    {#snippet actions(row)}
      <div class="row-actions">
        <Button variant="ghost" onclick={(e) => { e.stopPropagation(); openEdit(row); }}>✏️</Button>
        <Button variant="danger" onclick={(e) => { e.stopPropagation(); confirmDeleteId = row.id as string; }}>🗑️</Button>
      </div>
    {/snippet}
  </Table>

  <Modal
    open={showModal}
    title={editingId ? 'Редактировать запись ТО' : 'Добавить запись ТО'}
    onClose={() => showModal = false}
    width="520px"
  >
    <div class="form-grid">
      <div class="field">
        <label class="field-label" for="maint-date">Дата *</label>
        <input id="maint-date" class="field-input" type="date" bind:value={form.date} required />
      </div>
      <div class="field">
        <label class="field-label" for="maint-car">Автомобиль *</label>
        <select id="maint-car" class="field-select" bind:value={form.carId} required>
          <option value="">Выберите автомобиль</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model} ({car.plate})</option>
          {/each}
        </select>
      </div>
      <div class="field">
        <label class="field-label" for="maint-type">Тип обслуживания *</label>
        <select id="maint-type" class="field-select" bind:value={form.type}>
          {#each MAINTENANCE_TYPES as mt}
            <option value={mt.value}>{mt.label}</option>
          {/each}
        </select>
      </div>
      <Input
        label="Стоимость (₽) *"
        type="number"
        min="0"
        step="1"
        placeholder="5000"
        bind:value={form.cost}
        required
      />
      <div class="col-span-2">
        <Input
          label="Описание"
          placeholder="Дополнительные детали, пробег, мастерская..."
          bind:value={form.description}
        />
      </div>
    </div>
    {#snippet footer()}
      <Button variant="secondary" onclick={() => showModal = false}>Отмена</Button>
      <Button variant="primary" loading={saving} onclick={handleSave}>Сохранить</Button>
    {/snippet}
  </Modal>

  <Modal open={!!confirmDeleteId} title="Подтвердите удаление" onClose={() => confirmDeleteId = null}>
    <p class="confirm-text">Удалить запись о техобслуживании? Это действие нельзя отменить.</p>
    {#snippet footer()}
      <Button variant="secondary" onclick={() => confirmDeleteId = null}>Отмена</Button>
      <Button variant="danger" onclick={() => confirmDeleteId && handleDelete(confirmDeleteId)}>Удалить</Button>
    {/snippet}
  </Modal>

  <Toast open={toast.open} message={toast.message} type={toast.type} />
  <FAB label="Добавить ТО" onClick={openAdd} />
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

.row-actions { display: flex; gap: 0.25rem; }
.confirm-text { color: var(--text-secondary); line-height: 1.6; }
</style>
