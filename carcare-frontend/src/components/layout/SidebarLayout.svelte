<script lang="ts">
  export let title: string = '';
  import Sidebar from '../navigation/Sidebar.svelte';
</script>
<div class="sidebar-layout">
  <Sidebar />
  <main class="sidebar-main">
    {#if title}
      <header class="sidebar-header"><h1>{title}</h1></header>
    {/if}
    <section class="sidebar-content">
      <slot />
    </section>
  </main>
</div>
<style>
.sidebar-layout {
  display: flex;
  min-height: 100vh;
}
.sidebar-main {
  flex: 1;
  padding: 2rem 1rem 2rem 100px;
  max-width: 900px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
.sidebar-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: var(--accent);
  margin-bottom: 1rem;
}
.sidebar-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
:global(.dark) .sidebar-main {
  --accent: #7de2fc;
  color: #fff;
}
</style>
