<script lang="ts">
  import { ensureAuthenticated } from '../../lib/authGuard';
  import { onMount } from 'svelte';
  // TODO: подключить загрузку и отображение штрафов
  
  onMount(async () => {
    await ensureAuthenticated();
  });
</script>

<main>
  <h2>Штрафы</h2>
  <p>Здесь будет история штрафов.</p>
</main>
