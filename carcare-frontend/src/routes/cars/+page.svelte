<!-- src/routes/cars/+page.svelte -->

<script lang="ts">
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Card from '../../components/ui/Card.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';
  import FAB from '../../components/ui/FAB.svelte';
  import Modal from '../../components/ui/Modal.svelte';
  import Input from '../../components/ui/Input.svelte';
  import Button from '../../components/ui/Button.svelte';

  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import type { Car } from '../../lib/types';
  import { fetchCars, addCar, updateCar, deleteCar } from '../../lib/api';

  let loading = $state(true);
  let error = $state('');
  let cars = $state<Car[]>([]);

  // Форма добавления/редактирования
  let modalOpen = $state(false);
  let editingCar = $state<Car | null>(null);
  let saving = $state(false);
  let saveError = $state('');

  let formBrand = $state('');
  let formModel = $state('');
  let formYear = $state<string | number>('');
  let formVin = $state('');
  let formPlate = $state('');

  // Подтверждение удаления
  let deleteId = $state<string | null>(null);
  let deleting = $state(false);

  const currentYear = new Date().getFullYear();

  async function loadCars() {
    loading = true;
    error = '';
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : 'Ошибка загрузки авто';
      error = msg;
      if (msg.includes('авторизац')) setTimeout(() => ensureAuthenticated(), 1000);
    } finally {
      loading = false;
    }
  }

  function openAdd() {
    editingCar = null;
    formBrand = '';
    formModel = '';
    formYear = '';
    formVin = '';
    formPlate = '';
    saveError = '';
    modalOpen = true;
  }

  function openEdit(car: Car) {
    editingCar = car;
    formBrand = car.brand;
    formModel = car.model;
    formYear = car.year;
    formVin = car.vin;
    formPlate = car.plate ?? '';
    saveError = '';
    modalOpen = true;
  }

  function closeModal() {
    if (saving) return;
    modalOpen = false;
  }

  async function handleSave() {
    saveError = '';
    if (!String(formBrand).trim() || !String(formModel).trim() || !String(formYear).trim() || !String(formVin).trim()) {
      saveError = 'Заполните обязательные поля: Марка, Модель, Год, VIN';
      return;
    }
    const year = parseInt(String(formYear), 10);
    if (isNaN(year) || year < 1886 || year > currentYear) {
      saveError = `Укажите корректный год (1886 — ${currentYear})`;
      return;
    }

    saving = true;
    try {
      if (editingCar) {
        await updateCar({
          id: editingCar.id,
          brand: String(formBrand).trim(),
          model: String(formModel).trim(),
          year,
          vin: String(formVin).trim(),
          plate: String(formPlate).trim()
        });
      } else {
        await addCar({
          id: crypto.randomUUID(),
          brand: String(formBrand).trim(),
          model: String(formModel).trim(),
          year,
          vin: String(formVin).trim(),
          plate: String(formPlate).trim()
        });
      }
      modalOpen = false;
      await loadCars();
    } catch (e: unknown) {
      saveError = e instanceof Error ? e.message : 'Ошибка сохранения';
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!deleteId) return;
    deleting = true;
    try {
      await deleteCar(deleteId);
      cars = cars.filter(c => c.id !== deleteId);
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Ошибка удаления';
    } finally {
      deleting = false;
      deleteId = null;
    }
  }

  onMount(async () => {
    await ensureAuthenticated();
    loadCars();
  });
</script>

<PageLayout title="Мои автомобили">
  {#if loading}
    <Loader size={40} />
  {:else if error}
    <ErrorState message={error} />
  {:else if cars.length === 0}
    <EmptyState message="Нет добавленных авто" />
  {:else}
    <div class="cars-list">
      {#each cars as car (car.id)}
        <Card className="car-card">
          {#snippet children()}
          <div class="car-row">
            <div class="car-info">
              <div class="car-main">{car.brand} {car.model} <span class="car-year">({car.year})</span></div>
              {#if car.vin}
                <div class="car-vin">VIN: {car.vin}</div>
              {/if}
              {#if car.plate}
                <div class="car-plate">{car.plate}</div>
              {/if}
            </div>
            <div class="car-actions">
              <button class="icon-btn edit-btn" title="Редактировать" onclick={() => openEdit(car)}>
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                  <path d="M11.013 1.427a1.75 1.75 0 012.474 2.474L4.81 12.578l-3.182.354.354-3.181 8.031-8.324z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
              <button class="icon-btn delete-btn" title="Удалить" onclick={() => deleteId = car.id}>
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                  <path d="M6 2h4M2 4h12M5 4l.5 8h5L11 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
            </div>
          </div>
          {/snippet}
        </Card>
      {/each}
    </div>
  {/if}
</PageLayout>

<FAB label="Добавить авто" onClick={openAdd} position="fixed">
  {#snippet icon()}+{/snippet}
</FAB>

<!-- Модальное окно добавления/редактирования -->
<Modal
  open={modalOpen}
  title={editingCar ? 'Редактировать автомобиль' : 'Добавить автомобиль'}
  onClose={closeModal}
>
  {#snippet children()}
    <Input label="Марка" placeholder="Toyota" bind:value={formBrand} required />
    <Input label="Модель" placeholder="Camry" bind:value={formModel} required />
    <Input label="Год" type="number" placeholder="2022" bind:value={formYear} required inputProps={{min: 1886, max: currentYear}} />
    <Input label="VIN (17 символов)" placeholder="JTDBE30K673012345" bind:value={formVin} required />
    <Input label="Гос. номер" placeholder="А123БВ77" bind:value={formPlate} />
    {#if saveError}
      <div class="form-error" role="alert">{saveError}</div>
    {/if}
  {/snippet}
  {#snippet footer()}
    <Button variant="secondary" onclick={closeModal} disabled={saving}>Отмена</Button>
    <Button variant="primary" onclick={handleSave} loading={saving}>
      {editingCar ? 'Сохранить' : 'Добавить'}
    </Button>
  {/snippet}
</Modal>

<!-- Модальное окно подтверждения удаления -->
<Modal
  open={!!deleteId}
  title="Удалить автомобиль?"
  onClose={() => deleteId = null}
>
  {#snippet children()}
    {#if deleteId}
      {@const car = cars.find(c => c.id === deleteId)}
      <p class="confirm-text">
        Удалить <strong>{car?.brand} {car?.model}</strong>? Все связанные заправки, ТО и штрафы тоже будут удалены.
      </p>
    {/if}
  {/snippet}
  {#snippet footer()}
    <Button variant="secondary" onclick={() => deleteId = null} disabled={deleting}>Отмена</Button>
    <Button variant="danger" onclick={handleDelete} loading={deleting}>Удалить</Button>
  {/snippet}
</Modal>

<style>
.cars-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.car-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.car-info { flex: 1; min-width: 0; }

.car-main {
  font-weight: 600;
  font-size: 1rem;
  color: var(--text-primary);
}

.car-year {
  font-weight: 400;
  color: var(--text-secondary);
}

.car-vin {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
  font-family: monospace;
}

.car-plate {
  display: inline-block;
  margin-top: 0.375rem;
  font-size: 0.875rem;
  font-weight: 700;
  font-family: monospace;
  letter-spacing: 0.06em;
  background: var(--accent-light);
  color: var(--accent-text);
  padding: 0.125rem 0.5rem;
  border-radius: var(--radius-sm);
}

.car-actions {
  display: flex;
  gap: 0.375rem;
  flex-shrink: 0;
}

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

.edit-btn:hover {
  background: var(--accent-light);
  color: var(--accent);
  border-color: var(--accent);
}

.delete-btn:hover {
  background: var(--danger-light);
  color: var(--danger);
  border-color: var(--danger);
}

.confirm-text {
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
}
.confirm-text strong { color: var(--text-primary); }

.form-error {
  font-size: 0.875rem;
  color: var(--danger);
  padding: 0.5rem 0.75rem;
  background: var(--danger-light);
  border-radius: var(--radius-sm);
  border: 1px solid var(--danger);
}
</style>
