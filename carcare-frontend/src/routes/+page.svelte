<!-- src/routes/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { loadAuth, user } from '../stores/auth';

  let isAuthenticated = false;

  onMount(() => {
    loadAuth();
    const unsub = user.subscribe((u) => {
      isAuthenticated = !!u;
      if (!isAuthenticated) {
        goto('/welcome');
      }
    });
    return unsub;
  });
</script>
<main><h1>CarCare</h1></main>
