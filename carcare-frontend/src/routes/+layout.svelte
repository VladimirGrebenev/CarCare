<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import NavBar from '../components/navigation/NavBar.svelte';
  import Sidebar from '../components/navigation/Sidebar.svelte';
  import { theme } from '../lib/theme';
  import './+layout.css';

  let { children } = $props();
  let isMobile = $state(false);
  let currentPath = $derived($page.url.pathname);

  onMount(() => {
    theme.init();
    const update = () => (isMobile = window.innerWidth < 768);
    update();
    window.addEventListener('resize', update);
    return () => window.removeEventListener('resize', update);
  });

  function navigate(path: string) {
    goto(path);
  }
</script>

<div class="app-layout">
  {#if isMobile}
    <main class="mobile-content">{@render children()}</main>
    <NavBar active={currentPath} onNavigate={navigate} onToggleTheme={() => theme.toggle()} />
  {:else}
    <Sidebar active={currentPath} onNavigate={navigate} onToggleTheme={() => theme.toggle()} />
    <main class="desktop-content">{@render children()}</main>
  {/if}
</div>

<style>
.app-layout {
  min-height: 100vh;
  width: 100vw;
  display: flex;
  background: var(--bg-base);
  transition: background var(--transition-slow);
}
.desktop-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  min-height: 100vh;
  padding: 2rem;
  overflow-y: auto;
}
.mobile-content {
  flex: 1;
  min-height: 100vh;
  padding: 1rem;
  padding-bottom: 80px;
}
</style>
