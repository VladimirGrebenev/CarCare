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
  {#if title}
    <h1 class="page-title">{title}</h1>
  {/if}
  {#if toolbar}
    <div class="page-toolbar">
      {@render toolbar()}
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
  padding: 1.25rem 1rem 5rem 1rem;
  max-width: 860px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.page-title {
  font-size: 1.625rem;
  font-weight: 700;
  color: #4da6d9;
  margin: 0;
}
.page-toolbar {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  flex-wrap: wrap;
}
.page-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
</style>
