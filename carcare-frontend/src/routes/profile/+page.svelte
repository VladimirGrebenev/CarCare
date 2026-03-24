<script lang="ts">
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Card from '../../components/ui/Card.svelte';
  import Button from '../../components/ui/Button.svelte';
  import Modal from '../../components/ui/Modal.svelte';
  import Toast from '../../components/ui/Toast.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import { onMount } from 'svelte';

  let loading = $state(true);
  let error = $state('');
  let user = $state(null); // { name, avatarUrl, cars: Car[] }
  let showAddCar = $state(false);
  let toast = $state({ open: false, message: '', type: 'info' });


  import { fetchProfile as apiFetchProfile, addCar as apiAddCar } from '../../lib/api';

  async function fetchProfile() {
    loading = true;
    error = '';
    try {
      user = await apiFetchProfile();
    } catch (e) {
      error = e?.message || 'Не удалось загрузить профиль';
    } finally {
      loading = false;
    }
  }

  async function handleAddCar(car) {
    try {
      loading = true;
      await apiAddCar(car);
      await fetchProfile();
      toast = { open: true, message: 'Авто добавлено', type: 'success' };
    } catch (e) {
      toast = { open: true, message: e?.message || 'Ошибка добавления авто', type: 'error' };
    } finally {
      loading = false;
      showAddCar = false;
    }
  }

  onMount(fetchProfile);
</script>

<PageLayout title="Профиль">
  {#if loading}
    <div class="profile-loader"><Loader size={48} /></div>
  {:else if error}
    <ErrorState message={error} />
  {:else if !user}
    <EmptyState message="Нет данных профиля" />
  {:else}
    <section class="profile-main">
      <div class="profile-header glassmorphism">
        <img
          src={user.avatarUrl}
          alt="Аватар пользователя"
          class="profile-avatar"
          width="96" height="96"
        />
        <div class="profile-info">
          <h2 class="profile-name">{user.name}</h2>
          <Button variant="secondary" aria-label="Настройки профиля">
            ⚙️ Настройки
          </Button>
        </div>
      </div>

      <div class="profile-cars-section">
        <div class="profile-cars-header">
          <h3>Мои авто</h3>
          <Button
            variant="primary"
            aria-label="Добавить авто"
            onclick={() => showAddCar = true}
          >
            ➕ Добавить авто
          </Button>
        </div>
        {#if user.cars.length === 0}
          <EmptyState message="У вас пока нет добавленных авто" />
        {:else}
          <div class="profile-cars-list">
            {#each user.cars as car}
              <Card className="profile-car-card">
                <div class="car-main">
                  <div class="car-brand-model">{car.brand} {car.model}</div>
                  <div class="car-year">{car.year}</div>
                </div>
                <div class="car-plate">Гос. номер: {car.plate}</div>
              </Card>
            {/each}
          </div>
        {/if}
      </div>
    </section>

    <Modal open={showAddCar} title="Добавить авто" onClose={() => showAddCar = false}>
      <!-- TODO: Форма добавления авто (с валидацией) -->
      <div style="min-width:220px">Форма добавления авто (TODO)</div>
      <Button variant="primary" onclick={() => handleAddCar({ brand: 'Lada', model: 'Vesta', year: 2022, plate: 'C789FG' })} aria-label="Сохранить авто">Сохранить</Button>
    </Modal>
    <Toast open={toast.open} message={toast.message} type={toast.type} />
  {/if}
</PageLayout>

<style>
.profile-main {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 600px;
  margin: 0 auto;
  padding: 1.5rem 0 3rem 0;
}
.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.5rem 2rem;
  border-radius: 1.5rem;
  box-shadow: var(--glass-shadow);
  background: var(--glass-bg);
}
.profile-avatar {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #7de2fc;
  background: #fff;
}
.profile-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.profile-name {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
}
.profile-cars-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
.profile-cars-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}
.profile-cars-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
:global(.profile-car-card) {
  min-width: 0;
  padding: 1rem 1.2rem;
  border-radius: 1rem;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.car-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 1.1rem;
  font-weight: 600;
}
.car-plate {
  font-size: 0.95rem;
  color: var(--accent);
}
.profile-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}
@media (max-width: 700px) {
  .profile-main {
    max-width: 100vw;
    padding: 0.5rem 0 2rem 0;
  }
  .profile-header {
    flex-direction: column;
    align-items: flex-start;
    padding: 1rem 1rem;
    gap: 1rem;
  }
  .profile-cars-list {
    grid-template-columns: 1fr;
  }
}
</style>
