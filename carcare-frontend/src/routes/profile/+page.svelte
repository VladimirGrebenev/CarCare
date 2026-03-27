<script lang="ts">
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Modal from '../../components/ui/Modal.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import Input from '../../components/ui/Input.svelte';
  import { onMount } from 'svelte';
  import { ensureAuthenticated } from '../../lib/authGuard';
  import { fetchProfile as apiFetchProfile, addCar as apiAddCar, fetchCars } from '../../lib/api';
  import type { Car } from '../../lib/types';

  const CURRENT_YEAR = new Date().getFullYear();

  type ProfileData = {
    name?: string;
    email?: string;
    avatarUrl?: string;
    cars?: Car[];
    [key: string]: unknown;
  };

  let loading = $state(true);
  let error = $state('');
  let saving = $state(false);
  let deletingCarId = $state<string | null>(null);

  let profile = $state<ProfileData | null>(null);
  let cars = $state<Car[]>([]);

  let showAddCar = $state(false);
  let confirmDeleteCarId = $state<string | null>(null);
  let toast = $state<{ open: boolean; message: string; type: 'info' | 'success' | 'error' | 'warning' }>({
    open: false, message: '', type: 'info',
  });

  type CarForm = {
    brand: string;
    model: string;
    year: string;
    plate: string;
    vin: string;
  };

  let carForm = $state<CarForm>({ brand: '', model: '', year: String(CURRENT_YEAR), plate: '', vin: '' });
  let carFormError = $state('');

  function normalizeProfile(data: unknown): ProfileData | null {
    if (!data || typeof data !== 'object') return null;
    return data as ProfileData;
  }

  async function loadProfile() {
    loading = true;
    error = '';
    try {
      const data = await apiFetchProfile();
      profile = normalizeProfile(data);
      // Also load cars separately in case they're not in profile
      const carsData = await fetchCars();
      cars = Array.isArray(carsData) ? carsData : (Array.isArray(profile?.cars) ? profile.cars : []);
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Не удалось загрузить профиль';
      error = message;
      if (message.includes('авторизац')) {
        setTimeout(() => ensureAuthenticated(), 1000);
      }
    } finally {
      loading = false;
    }
  }

  function openAddCar() {
    carForm = { brand: '', model: '', year: String(CURRENT_YEAR), plate: '', vin: '' };
    carFormError = '';
    showAddCar = true;
  }

  function validateCarForm(): string {
    if (!carForm.brand.trim()) return 'Укажите марку автомобиля';
    if (!carForm.model.trim()) return 'Укажите модель автомобиля';
    const year = Number(carForm.year);
    if (!year || year < 1990 || year > CURRENT_YEAR + 1) {
      return `Год должен быть от 1990 до ${CURRENT_YEAR + 1}`;
    }
    if (!carForm.plate.trim()) return 'Укажите государственный номер';
    if (!carForm.vin.trim()) return 'VIN является обязательным полем';
    if (carForm.vin.trim().length !== 17) return 'VIN должен содержать ровно 17 символов';
    return '';
  }

  async function submitCarForm() {
    carFormError = '';
    const validationError = validateCarForm();
    if (validationError) {
      carFormError = validationError;
      return;
    }
    saving = true;
    try {
      const newCar = await apiAddCar({
        brand: carForm.brand.trim(),
        model: carForm.model.trim(),
        year: Number(carForm.year),
        plate: carForm.plate.trim().toUpperCase(),
        vin: carForm.vin.trim().toUpperCase(),
      });
      // Add the new car to the list
      if (newCar && typeof newCar === 'object') {
        cars = [newCar as Car, ...cars];
      } else {
        // Reload from API if the response doesn't contain car data
        const refreshed = await fetchCars();
        cars = Array.isArray(refreshed) ? refreshed : cars;
      }
      showAddCar = false;
      toast = { open: true, message: 'Автомобиль успешно добавлен', type: 'success' };
    } catch (e: unknown) {
      toast = {
        open: true,
        message: e instanceof Error ? e.message : 'Ошибка добавления авто',
        type: 'error',
      };
    } finally {
      saving = false;
    }
  }

  async function handleDeleteCar(carId: string) {
    deletingCarId = carId;
    try {
      const res = await fetch(`/api/cars/${carId}`, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          Authorization: `Bearer ${typeof window !== 'undefined' ? localStorage.getItem('authToken') ?? '' : ''}`,
        },
      });
      if (!res.ok && res.status !== 404) {
        throw new Error(`Ошибка удаления (${res.status})`);
      }
      cars = cars.filter(c => c.id !== carId);
      toast = { open: true, message: 'Автомобиль удалён', type: 'success' };
    } catch (e: unknown) {
      toast = {
        open: true,
        message: e instanceof Error ? e.message : 'Ошибка удаления авто',
        type: 'error',
      };
    } finally {
      deletingCarId = null;
      confirmDeleteCarId = null;
    }
  }

  function getUserInitials(name: string | undefined): string {
    if (!name) return '?';
    return name
      .split(' ')
      .filter(Boolean)
      .slice(0, 2)
      .map(w => w[0].toUpperCase())
      .join('');
  }

  onMount(async () => {
    await ensureAuthenticated();
    loadProfile();
  });
</script>

<PageLayout title="Профиль">
  {#if loading}
    <div class="center-loader"><Loader size={48} /></div>
  {:else if error}
    <ErrorState message={error} />
  {:else}
    <div class="profile-layout">
      <!-- User card -->
      <section class="profile-header glassmorphism">
        <div class="avatar-wrap">
          {#if profile?.avatarUrl}
            <img
              src={profile.avatarUrl}
              alt="Аватар"
              class="avatar-img"
              width="80"
              height="80"
            />
          {:else}
            <div class="avatar-placeholder" aria-hidden="true">
              {getUserInitials(profile?.name as string | undefined)}
            </div>
          {/if}
        </div>
        <div class="profile-info">
          {#if profile?.name}
            <h2 class="profile-name">{profile.name}</h2>
          {/if}
          {#if profile?.email}
            <p class="profile-email">{profile.email}</p>
          {/if}
        </div>
      </section>

      <!-- Cars section -->
      <section class="cars-section">
        <div class="section-header">
          <h3 class="section-title">Мои автомобили <span class="cars-count">{cars.length}</span></h3>
          <Button variant="primary" onclick={openAddCar}>
            + Добавить авто
          </Button>
        </div>

        {#if cars.length === 0}
          <EmptyState message="У вас пока нет добавленных автомобилей" />
        {:else}
          <div class="cars-grid">
            {#each cars as car (car.id)}
              <div class="car-card glassmorphism">
                <div class="car-header">
                  <div class="car-icon" aria-hidden="true">🚗</div>
                  <div class="car-meta">
                    <div class="car-name">{car.brand} {car.model}</div>
                    <div class="car-year">{car.year} г.в.</div>
                  </div>
                  <button
                    class="car-delete-btn"
                    aria-label="Удалить {car.brand} {car.model}"
                    onclick={() => confirmDeleteCarId = car.id}
                    disabled={deletingCarId === car.id}
                    title="Удалить автомобиль"
                  >
                    {deletingCarId === car.id ? '⏳' : '🗑️'}
                  </button>
                </div>
                <div class="car-details">
                  <div class="car-detail-row">
                    <span class="detail-label">Гос. номер</span>
                    <span class="detail-value plate">{car.plate}</span>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </section>
    </div>

    <!-- Add car modal -->
    <Modal
      open={showAddCar}
      title="Добавить автомобиль"
      onClose={() => { showAddCar = false; carFormError = ''; }}
      width="500px"
    >
      <div class="car-form">
        <Input
          label="Марка *"
          placeholder="Toyota, BMW, Lada..."
          bind:value={carForm.brand}
          required
        />
        <Input
          label="Модель *"
          placeholder="Camry, X5, Vesta..."
          bind:value={carForm.model}
          required
        />
        <div class="field">
          <label class="field-label" for="car-year">Год выпуска *</label>
          <input
            id="car-year"
            class="field-input"
            type="number"
            min="1990"
            max={CURRENT_YEAR + 1}
            placeholder={String(CURRENT_YEAR)}
            bind:value={carForm.year}
          />
        </div>
        <Input
          label="Гос. номер *"
          placeholder="А123БВ77"
          value={carForm.plate}
          oninput={(e) => carForm.plate = (e.target as HTMLInputElement).value.toUpperCase()}
          required
        />
        <div class="col-span-2">
          <Input
            label="VIN * (17 символов)"
            placeholder="WVWZZZ1JZ3W386752"
            value={carForm.vin}
            oninput={(e) => {
              carForm.vin = (e.target as HTMLInputElement).value.toUpperCase();
            }}
            hint={carForm.vin ? `${carForm.vin.length}/17 символов` : ''}
            error={carForm.vin && carForm.vin.length > 0 && carForm.vin.length !== 17
              ? `Введено ${carForm.vin.length} из 17 символов`
              : ''}
            required
          />
        </div>
        {#if carFormError}
          <p class="form-error col-span-2" role="alert">{carFormError}</p>
        {/if}
      </div>
      {#snippet footer()}
        <Button variant="secondary" onclick={() => { showAddCar = false; carFormError = ''; }}>
          Отмена
        </Button>
        <Button variant="primary" loading={saving} onclick={submitCarForm}>
          Добавить
        </Button>
      {/snippet}
    </Modal>

    <!-- Confirm delete car modal -->
    <Modal
      open={!!confirmDeleteCarId}
      title="Удалить автомобиль?"
      onClose={() => confirmDeleteCarId = null}
    >
      {#if confirmDeleteCarId}
        {@const car = cars.find(c => c.id === confirmDeleteCarId)}
        <p class="confirm-text">
          {#if car}
            Вы уверены, что хотите удалить <strong>{car.brand} {car.model} ({car.plate})</strong>?
            Все связанные записи о заправках, ТО и штрафах будут потеряны.
          {:else}
            Удалить этот автомобиль? Это действие нельзя отменить.
          {/if}
        </p>
      {/if}
      {#snippet footer()}
        <Button variant="secondary" onclick={() => confirmDeleteCarId = null}>Отмена</Button>
        <Button
          variant="danger"
          loading={!!deletingCarId}
          onclick={() => confirmDeleteCarId && handleDeleteCar(confirmDeleteCarId)}
        >
          Удалить
        </Button>
      {/snippet}
    </Modal>

    <Toast open={toast.open} message={toast.message} type={toast.type} />
  {/if}
</PageLayout>

<style>
.center-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.profile-layout {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 640px;
  margin: 0 auto;
  padding-bottom: 3rem;
}

/* ── Profile header ── */
.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.75rem 2rem;
  border-radius: var(--radius-xl);
}

.avatar-wrap { flex-shrink: 0; }

.avatar-img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid var(--accent);
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.02em;
  user-select: none;
}

.profile-info { display: flex; flex-direction: column; gap: 0.375rem; }
.profile-name { font-size: 1.375rem; font-weight: 700; margin: 0; }
.profile-email { font-size: 0.9375rem; color: var(--text-secondary); margin: 0; }

/* ── Cars section ── */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.25rem;
}
.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.cars-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.5rem;
  height: 1.5rem;
  padding: 0 0.375rem;
  background: var(--accent-light);
  color: var(--accent-text);
  border-radius: 100px;
  font-size: 0.8125rem;
  font-weight: 700;
}

.cars-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

/* ── Car card ── */
.car-card {
  padding: 1.25rem 1.375rem;
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  gap: 0.875rem;
  transition: transform var(--transition), box-shadow var(--transition);
}
.car-card:hover { transform: translateY(-2px); box-shadow: var(--shadow-md); }

.car-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.car-icon { font-size: 1.5rem; flex-shrink: 0; }
.car-meta { flex: 1; min-width: 0; }
.car-name { font-size: 1rem; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.car-year { font-size: 0.8125rem; color: var(--text-secondary); }

.car-delete-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 0.375rem;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  transition: background var(--transition), color var(--transition);
  flex-shrink: 0;
  line-height: 1;
}
.car-delete-btn:hover:not(:disabled) {
  background: var(--danger-light);
  color: var(--danger);
}
.car-delete-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.car-details { display: flex; flex-direction: column; gap: 0.5rem; }
.car-detail-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}
.detail-label { font-size: 0.8125rem; color: var(--text-secondary); }
.detail-value { font-size: 0.9375rem; font-weight: 500; }
.plate {
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--accent-text);
  letter-spacing: 0.06em;
  background: var(--accent-light);
  padding: 0.125rem 0.5rem;
  border-radius: var(--radius-sm);
}

/* ── Add car form ── */
.car-form {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
.col-span-2 { grid-column: span 2; }

.field { display: flex; flex-direction: column; gap: 0.375rem; }
.field-label { font-size: 0.8125rem; font-weight: 600; color: var(--text-secondary); }
.field-input {
  padding: 0.5625rem 0.875rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text-primary);
  font-size: 0.9375rem;
  font-family: var(--font);
  outline: none;
  transition: border-color var(--transition), box-shadow var(--transition);
}
.field-input:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
}

.form-error {
  color: var(--danger);
  font-size: 0.875rem;
  margin: 0;
  padding: 0.5rem 0.75rem;
  background: var(--danger-light);
  border-radius: var(--radius-md);
  border: 1px solid rgba(248, 81, 73, 0.25);
}

.confirm-text {
  color: var(--text-secondary);
  line-height: 1.7;
  margin: 0;
}
.confirm-text strong { color: var(--text-primary); }

/* ── Responsive ── */
@media (max-width: 640px) {
  .profile-layout { max-width: 100%; padding: 0 0 2rem; }
  .profile-header { flex-direction: column; align-items: flex-start; padding: 1.25rem; }
  .cars-grid { grid-template-columns: 1fr; }
  .car-form { grid-template-columns: 1fr; }
  .col-span-2 { grid-column: 1; }
}
</style>
