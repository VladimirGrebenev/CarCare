<!-- src/routes/fines/+page.svelte -->
<script lang="ts">
// Fines page logic

<script lang="ts">
import { onMount } from 'svelte';
import { fines, finesLoading, finesError, finesSuccess, finesFilters } from '../../stores/fines';
import { fetchFines, addFine, updateFine, deleteFine } from '../../lib/api';
import type { Fine } from '../../lib/types';

// Загрузка штрафов
const loadFines = async () => {
  finesLoading.set(true);
  finesError.set(null);
  try {
    const data = await fetchFines(finesFilters.get());
    fines.set(data);
  } catch (e) {
    finesError.set(e.message || 'Ошибка загрузки штрафов');
  } finally {
    finesLoading.set(false);
  }
};

onMount(loadFines);

// CRUD операции
const handleAdd = async (fine: Partial<Fine>) => {
  finesLoading.set(true);
  finesError.set(null);
  try {
    await addFine(fine);
    finesSuccess.set('Штраф добавлен');
    await loadFines();
  } catch (e) {
    finesError.set(e.message || 'Ошибка добавления штрафа');
  } finally {
    finesLoading.set(false);
  }
};

const handleUpdate = async (id: string, fine: Partial<Fine>) => {
  finesLoading.set(true);
  finesError.set(null);
  try {
    await updateFine(id, fine);
    finesSuccess.set('Штраф обновлён');
    await loadFines();
  } catch (e) {
    finesError.set(e.message || 'Ошибка обновления штрафа');
  } finally {
    finesLoading.set(false);
  }
};

const handleDelete = async (id: string) => {
  finesLoading.set(true);
  finesError.set(null);
  try {
    await deleteFine(id);
    finesSuccess.set('Штраф удалён');
    await loadFines();
  } catch (e) {
    finesError.set(e.message || 'Ошибка удаления штрафа');
  } finally {
    finesLoading.set(false);
  }
};
</script>
  <h2>Штрафы</h2>
  <p>Здесь будет история штрафов.</p>
</main>
