<script lang="ts">
  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Table from '../../components/ui/Table.svelte';
  import Modal from '../../components/ui/Modal.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import FAB from '../../components/ui/FAB.svelte';
  import {
    fetchMaintenanceHistory,
    addMaintenance,
    updateMaintenance,
    deleteMaintenance,
    fetchCars,
  } from '../../lib/api';
  import type { Car } from '../../lib/types';

  // Встроенные типы услуг
  const BUILTIN_TYPES = [
    'Замена масла',
    'Шиномонтаж / смена резины',
    'Тормозная система',
    'Замена фильтров',
    'Аккумулятор',
    'Техосмотр',
    'Замена свечей зажигания',
    'Замена ремня ГРМ',
    'Антифриз / охлаждающая жидкость',
    'Развал-схождение',
    'Промывка форсунок',
    'Замена тормозных дисков и колодок',
  ];

  const CUSTOM_TYPES_KEY = 'carcare_maintenance_custom_types';

  function loadCustomTypes(): string[] {
    try {
      return JSON.parse(localStorage.getItem(CUSTOM_TYPES_KEY) ?? '[]');
    } catch { return []; }
  }

  function saveCustomType(typeName: string) {
    const custom = loadCustomTypes();
    const trimmed = typeName.trim();
    if (!trimmed || BUILTIN_TYPES.includes(trimmed) || custom.includes(trimmed)) return;
    custom.unshift(trimmed);
    localStorage.setItem(CUSTOM_TYPES_KEY, JSON.stringify(custom.slice(0, 50)));
  }

  const COLUMNS = [
    { label: 'Дата', key: '_date' },
    { label: 'Автомобиль', key: '_car' },
    { label: 'Услуга', key: 'type' },
    { label: 'Стоимость (₽)', key: '_cost' },
  ];

  type MaintenanceRecord = {
    id: string;
    date: string;
    carId: string;
    type: string;
    cost: number;
    [key: string]: unknown;
  };

  let items = $state<MaintenanceRecord[]>([]);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let cars = $state<Car[]>([]);
  let customTypes = $state<string[]>([]);

  let filterCarId = $state('');

  // Все доступные типы = встроенные + пользовательские
  let allTypes = $derived([...BUILTIN_TYPES, ...customTypes.filter(t => !BUILTIN_TYPES.includes(t))]);

  function getCarLabel(carId: string): string {
    if (!carId) return '—';
    const car = cars.find(c => c.id === carId);
    if (!car) return carId;
    return `${car.brand} ${car.model}${car.plate ? ` (${car.plate})` : ''}`;
  }

  let rows = $derived(
    items
      .filter(item => !filterCarId || item.carId === filterCarId)
      .map(item => ({
        ...item,
        _date: item.date ?? '—',
        _car: getCarLabel(String(item.carId ?? '')),
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
  };

  let form = $state<MaintenanceForm>({ date: '', carId: '', type: '', cost: '' });

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
      type: '',
      cost: '',
    };
    showModal = true;
  }

  function openEdit(row: Record<string, unknown>) {
    editingId = String(row.id ?? '');
    form = {
      date: String(row.date ?? ''),
      carId: String(row.carId ?? ''),
      type: String(row.type ?? ''),
      cost: String(row.cost ?? ''),
    };
    showModal = true;
  }

  async function handleSave() {
    if (!form.carId) {
      toast = { open: true, message: 'Выберите автомобиль', type: 'error' };
      return;
    }
    if (!form.type.trim()) {
      toast = { open: true, message: 'Укажите тип услуги', type: 'error' };
      return;
    }
    const cost = Number(form.cost);
    if (!form.cost || isNaN(cost) || cost < 0) {
      toast = { open: true, message: 'Укажите корректную стоимость (≥ 0)', type: 'error' };
      return;
    }

    // Сохраняем пользовательский тип если он не из встроенных
    saveCustomType(form.type);
    customTypes = loadCustomTypes();

    saving = true;
    try {
      const payload = { ...form, type: form.type.trim(), cost };
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
      toast = { open: true, message: e instanceof Error ? e.message : 'Ошибка сохранения', type: 'error' };
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
      toast = { open: true, message: e instanceof Error ? e.message : 'Ошибка удаления', type: 'error' };
    } finally {
      confirmDeleteId = null;
    }
  }

  onMount(async () => {
    await ensureAuthenticated();
    customTypes = loadCustomTypes();
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch { cars = []; }
    await loadItems();
  });
</script>

<PageLayout title="Техобслуживание">
  <div class="page-toolbar">
    <div class="filters">
      <div class="filter-field">
        <label class="filter-label" for="maint-filter-car">Автомобиль</label>
        <select id="maint-filter-car" class="field-select" bind:value={filterCarId}>
          <option value="">Все автомобили</option>
          {#each cars as car}
            <option value={car.id}>{car.brand} {car.model}{car.plate ? ` (${car.plate})` : ''}</option>
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
    width="480px"
  >
    {#snippet children()}
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
            <option value={car.id}>{car.brand} {car.model}{car.plate ? ` (${car.plate})` : ''}</option>
          {/each}
        </select>
      </div>

      <div class="field col-span-2">
        <label class="field-label" for="maint-type">Услуга *</label>
        <input
          id="maint-type"
          class="field-input"
          type="text"
          list="maint-type-list"
          placeholder="Начните вводить или выберите из списка..."
          bind:value={form.type}
          autocomplete="off"
        />
        <datalist id="maint-type-list">
          {#if customTypes.length > 0}
            <optgroup label="Мои услуги">
              {#each customTypes as t}
                <option value={t}>{t}</option>
              {/each}
            </optgroup>
          {/if}
          {#each allTypes.filter(t => !customTypes.includes(t)) as t}
            <option value={t}>{t}</option>
          {/each}
        </datalist>
        <span class="field-hint">Выберите из списка или введите свою услугу — она сохранится для следующих записей</span>
      </div>

      <div class="field col-span-2">
        <label class="field-label" for="maint-cost">Стоимость (₽) *</label>
        <input
          id="maint-cost"
          class="field-input"
          type="number"
          min="0"
          step="1"
          placeholder="5000"
          bind:value={form.cost}
          required
        />
      </div>
    </div>
    {/snippet}
    {#snippet footer()}
      <Button variant="secondary" onclick={() => showModal = false}>Отмена</Button>
      <Button variant="primary" loading={saving} onclick={handleSave}>Сохранить</Button>
    {/snippet}
  </Modal>

  <Modal open={!!confirmDeleteId} title="Подтвердите удаление" onClose={() => confirmDeleteId = null}>
    {#snippet children()}
    <p class="confirm-text">Удалить запись о техобслуживании? Это действие нельзя отменить.</p>
    {/snippet}
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
}
.col-span-2 { grid-column: span 2; }

.field { display: flex; flex-direction: column; gap: 0.375rem; }
.field-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }
.field-hint { font-size: 0.75rem; color: var(--text-secondary); opacity: 0.7; }

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
.confirm-text { color: var(--text-secondary); line-height: 1.6; margin: 0; }
</style>
