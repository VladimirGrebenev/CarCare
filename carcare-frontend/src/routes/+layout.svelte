<script lang="ts">
  import NavBar from '../components/navigation/NavBar.svelte';
  import Sidebar from '../components/navigation/Sidebar.svelte';
  let isMobile = window.innerWidth < 768;
  $effect(() => {
    const handler = () => isMobile = window.innerWidth < 768;
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  });
</script>
<div class="app-layout">
  {#if isMobile}
    <slot />
    <NavBar />
  {:else}
    <Sidebar />
    <div class="desktop-content"><slot /></div>
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
