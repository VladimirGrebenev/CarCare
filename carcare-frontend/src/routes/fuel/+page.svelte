<!-- src/routes/fuel/+page.svelte -->
<script lang="ts">
// Fuel page logic
</script>
<script lang="ts">
import { onMount } from 'svelte';
import {
  fuelList, filteredFuelList, fuelLoading, fuelError, fuelFilters,
  loadFuel, createFuel, editFuel, removeFuel
} from '../../stores/fuel';

let newFuel = { date: '', liters: '', price: '', carId: '' };

function handleFilterChange(e) {
  fuelFilters.value = { ...fuelFilters.value, [e.target.name]: e.target.value };
  loadFuel();
}

function handleCreate(e) {
  e.preventDefault();
  createFuel(newFuel).then(() => {
    newFuel = { date: '', liters: '', price: '', carId: '' };
  });
}

onMount(() => {
  loadFuel();
});
</script>
  <p>Здесь будет история заправок.</p>
</main>
