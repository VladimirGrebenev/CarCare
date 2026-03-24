<script lang="ts">
  import { ensureAuthenticated } from '../../lib/authGuard';
  import { onMount } from 'svelte';
  // TODO: подключить загрузку и отображение заправок
  
  onMount(async () => {
    await ensureAuthenticated();
  });
</script>

<main>
  <h2>Заправки</h2>
  <p>Здесь будет история заправок.</p>
</main>
