<!-- src/routes/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { bootstrapAuth } from '../stores/auth';

  onMount(() => {
    let disposed = false;

    (async () => {
      const isAuthenticated = await bootstrapAuth();
      if (disposed) {
        return;
      }

      goto(isAuthenticated ? '/profile' : '/welcome', { replaceState: true });
    })();

    return () => {
      disposed = true;
    };
  });
</script>
<main><h1>CarCare</h1></main>
