<!-- src/routes/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import WidgetTile from '../components/widgets/WidgetTile.svelte';
  import Loader from '../components/ui/Loader.svelte';
  import ErrorState from '../components/ui/ErrorState.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';
  import { fetchProfile, fetchCars, fetchFuelHistory, fetchMaintenanceHistory, fetchFines } from '../lib/api';
  import { $state } from 'svelte/store';

  // Dashboard widgets state (drag&drop order)
  let widgetOrder = $state([
    'profile',
    'cars',
    'fuel',
    'maintenance',
    'fines'
  ]);

  // Widget data and states
  let profile = $state({ data: null, loading: true, error: null });
  let cars = $state({ data: null, loading: true, error: null });
  let fuel = $state({ data: null, loading: true, error: null });
  let maintenance = $state({ data: null, loading: true, error: null });
  let fines = $state({ data: null, loading: true, error: null });

  // Fetch all widgets data
  async function loadAll() {
    profile.loading = cars.loading = fuel.loading = maintenance.loading = fines.loading = true;
    profile.error = cars.error = fuel.error = maintenance.error = fines.error = null;
    try {
      profile.data = await fetchProfile();
    } catch (e) { profile.error = e instanceof Error ? e.message : 'Ошибка профиля'; }
    profile.loading = false;
    try {
      cars.data = await fetchCars();
    } catch (e) { cars.error = e instanceof Error ? e.message : 'Ошибка авто'; }
    cars.loading = false;
    try {
      fuel.data = await fetchFuelHistory();
    } catch (e) { fuel.error = e instanceof Error ? e.message : 'Ошибка заправок'; }
    fuel.loading = false;
    try {
      maintenance.data = await fetchMaintenanceHistory();
    } catch (e) { maintenance.error = e instanceof Error ? e.message : 'Ошибка ТО'; }
    maintenance.loading = false;
    try {
      fines.data = await fetchFines();
    } catch (e) { fines.error = e instanceof Error ? e.message : 'Ошибка штрафов'; }
    fines.loading = false;
  }

  onMount(loadAll);

  // Drag&Drop reorder logic (local only)
  let dragIndex: number | null = null;
  function handleDragStart(idx: number) { dragIndex = idx; }
  function handleDrop(idx: number) {
    if (dragIndex === null || dragIndex === idx) return;
    const updated = [...widgetOrder];
    const [moved] = updated.splice(dragIndex, 1);
    updated.splice(idx, 0, moved);
    widgetOrder = updated;
    dragIndex = null;
    // TODO: persist order (e.g. localStorage)
  }
</script>

<main>
  <h1>CarCare</h1>
  <p>Добро пожаловать в CarCare — ваш помощник по учёту расходов, заправок, ТО и штрафов.</p>

  <div class="dashboard-widgets">
    {#each widgetOrder as key, idx}
      <div
        class="dashboard-widget"
        draggable="true"
        ondragstart={() => handleDragStart(idx)}
        ondragover|preventDefault
        ondrop={() => handleDrop(idx)}
      >
        {#if key === 'profile'}
          <WidgetTile title="Профиль" icon={() => '👤'}>
            {#if profile.loading}
              <Loader />
            {:else if profile.error}
              <ErrorState message={profile.error} onRetry={loadAll} />
            {:else if profile.data}
              <div>Здравствуйте, {profile.data.name || profile.data.email}</div>
            {/if}
          </WidgetTile>
        {:else if key === 'cars'}
          <WidgetTile title="Авто" icon={() => '🚗'} value={cars.data?.length ?? ''}>
            {#if cars.loading}
              <Loader />
            {:else if cars.error}
              <ErrorState message={cars.error} onRetry={loadAll} />
            {:else if cars.data?.length === 0}
              <EmptyState message="Нет добавленных авто" />
            {:else if cars.data}
              <div>{cars.data.length} авто</div>
            {/if}
          </WidgetTile>
        {:else if key === 'fuel'}
          <WidgetTile title="Заправки" icon={() => '⛽'} value={fuel.data?.length ?? ''}>
            {#if fuel.loading}
              <Loader />
            {:else if fuel.error}
              <ErrorState message={fuel.error} onRetry={loadAll} />
            {:else if fuel.data?.length === 0}
              <EmptyState message="Нет заправок" />
            {:else if fuel.data}
              <div>Последняя: {fuel.data[0]?.date} — {fuel.data[0]?.liters} л</div>
            {/if}
          </WidgetTile>
        {:else if key === 'maintenance'}
          <WidgetTile title="ТО" icon={() => '🛠️'} value={maintenance.data?.length ?? ''}>
            {#if maintenance.loading}
              <Loader />
            {:else if maintenance.error}
              <ErrorState message={maintenance.error} onRetry={loadAll} />
            {:else if maintenance.data?.length === 0}
              <EmptyState message="Нет записей ТО" />
            {:else if maintenance.data}
              <div>Последнее: {maintenance.data[0]?.date} — {maintenance.data[0]?.type || ''}</div>
            {/if}
          </WidgetTile>
        {:else if key === 'fines'}
          <WidgetTile title="Штрафы" icon={() => '💸'} value={fines.data?.length ?? ''}>
            {#if fines.loading}
              <Loader />
            {:else if fines.error}
              <ErrorState message={fines.error} onRetry={loadAll} />
            {:else if fines.data?.length === 0}
              <EmptyState message="Штрафов нет!" />
            {:else if fines.data}
              <div>Сумма: {fines.data.reduce((sum, f) => sum + (f.amount || 0), 0)} ₽</div>
            {/if}
          </WidgetTile>
        {/if}
      </div>
    {/each}
  </div>
</main>

<style>
.dashboard-widgets {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}
.dashboard-widget {
  cursor: grab;
}
.dashboard-widget:active {
  cursor: grabbing;
}
</style>
