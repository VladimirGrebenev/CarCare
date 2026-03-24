<script lang="ts">
  export let title: string = '';
</script>
<main class="page-layout">
  {#if title}
    <header class="page-header"><h1>{title}</h1></header>
  {/if}
  <section class="page-content">
    <slot />
  </section>
</main>
<style>
.page-layout {
  min-height: 100vh;
  padding: 2rem 1rem 5rem 1rem;
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: var(--accent);
  margin-bottom: 1rem;
}
.page-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
:global(.dark) .page-layout {
  --accent: #7de2fc;
  color: #fff;
}
</style>
