<script lang="ts">
  import { onMount } from 'svelte';
  import NavBar from '../components/navigation/NavBar.svelte';
  import Sidebar from '../components/navigation/Sidebar.svelte';
  let { children } = $props();
  let isMobile = $state(false);

  onMount(() => {
    const update = () => (isMobile = window.innerWidth < 768);
    update();
    window.addEventListener('resize', update);
    return () => window.removeEventListener('resize', update);
  });
</script>
<div class="app-layout">
  {#if isMobile}
    {@render children()}
    <NavBar />
  {:else}
    <Sidebar />
    <div class="desktop-content">{@render children()}</div>
  {/if}
</div>
<style>
.app-layout {
  min-height: 100vh;
  width: 100vw;
  background: linear-gradient(135deg, #232526 0%, #414345 100%);
  display: flex;
  flex-direction: column;
}
.desktop-content {
  flex: 1;
  margin-left: 80px;
  min-height: 100vh;
}
:global(.dark) .app-layout {
  background: linear-gradient(135deg, #181a20 0%, #232526 100%);
}
</style>
