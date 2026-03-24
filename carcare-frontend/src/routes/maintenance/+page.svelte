<!-- src/routes/maintenance/+page.svelte -->
<script lang="ts">
import { maintenanceStore } from '../../stores/maintenance';
import Loader from '../../components/ui/Loader.svelte';
import ErrorState from '../../components/ui/ErrorState.svelte';
import EmptyState from '../../components/ui/EmptyState.svelte';
import Card from '../../components/ui/Card.svelte';
import Button from '../../components/ui/Button.svelte';
import { ensureAuthenticated } from '../../lib/authGuard';

import { onMount } from 'svelte';
let filters = { type: '', carId: '' };
onMount(async () => {
  await ensureAuthenticated();
  maintenanceStore.load(filters);
});

function handleRetry() {
  maintenanceStore.load(filters);
}

function handleAdd() {
  // TODO: реализовать модалку/форму добавления
  alert('Добавление ТО пока не реализовано');
}
</script>

<main>
  <h2>Техобслуживание</h2>
  <div style="margin-bottom:1rem;display:flex;gap:1rem;align-items:center;">
    <input
      placeholder="Тип ТО"
      value={filters.type}
      oninput={(e) => (filters.type = (e.currentTarget as HTMLInputElement).value)}
      style="padding:0.5rem;border-radius:0.5rem;border:1px solid #ccc;"
    />
    <input
      placeholder="ID авто"
      value={filters.carId}
      oninput={(e) => (filters.carId = (e.currentTarget as HTMLInputElement).value)}
      style="padding:0.5rem;border-radius:0.5rem;border:1px solid #ccc;"
    />
    <Button onclick={handleAdd}>Добавить ТО</Button>
  </div>

  {#if maintenanceStore.loading}
    <Loader size={40} />
  {:else if maintenanceStore.error}
    <ErrorState message={maintenanceStore.error} onRetry={handleRetry} />
  {:else if maintenanceStore.items.length === 0}
    <EmptyState message="Нет записей о ТО" />
  {:else}
    <div style="display:grid;gap:1rem;">
      {#each maintenanceStore.items as item (item.id)}
        <Card>
          <div style="font-weight:600;">{item.type} ({item.date})</div>
          <div>Машина: {item.carId}</div>
          <div>Стоимость: {item.cost} ₽</div>
        </Card>
      {/each}
    </div>
  {/if}
</main>
