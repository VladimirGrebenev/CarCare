<!-- src/routes/cars/+page.svelte -->

<script lang="ts">
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import Card from '../../components/ui/Card.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';

  import { onMount } from 'svelte';

  import type { Car } from '../../lib/types';
  import { fetchCars } from '../../lib/api';





  let loading: boolean = true;
  let error: string = '';

  let cars: Car[] = [];

  async function loadCars() {
    loading = true;
    error = '';
    try {
      const result = await fetchCars();
      cars = Array.isArray(result) ? result : [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Ошибка загрузки авто';
    } finally {
      loading = false;
    }
  }

  onMount(loadCars);
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
          <div class="car-main">{car.brand} {car.model} ({car.year})</div>
          <div class="car-plate">Гос. номер: {car.plate}</div>
        </Card>
      {/each}
    </div>
  {/if}
</PageLayout>
