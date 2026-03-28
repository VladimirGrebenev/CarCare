<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title = '',
    toolbar,
    children,
  }: {
    title?: string;
    toolbar?: Snippet;
    children?: Snippet;
  } = $props();
</script>
<main class="page-layout">
  {#if title || toolbar}
    <div class="page-header">
      {#if title}
        <h1 class="page-title">{title}</h1>
      {/if}
      {#if toolbar}
        <div class="page-toolbar-inline">
          {@render toolbar()}
        </div>
      {/if}
    </div>
  {/if}
  <section class="page-content">
    {#if children}
      {@render children()}
    {/if}
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
  gap: 1rem;
}
.page-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.25rem;
  flex-wrap: wrap;
}
.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--accent);
  white-space: nowrap;
  margin: 0;
}
.page-toolbar-inline {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
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
