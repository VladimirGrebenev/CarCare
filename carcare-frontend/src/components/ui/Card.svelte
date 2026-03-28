<script lang="ts">
  import type { Snippet } from 'svelte';

  type Props = {
    children?: Snippet;
    header?: Snippet | null;
    footer?: Snippet | null;
    class?: string;
    className?: string;
  };

  let {
    children = null,
    header = null,
    footer = null,
    class: classAttr = '',
    className = ''
  }: Props = $props();
</script>

<div class="card spotlight {className} {classAttr}">
  {#if header}
    <div class="card-header">{@render header()}</div>
  {/if}
  <div class="card-content">{@render children?.()}</div>
  {#if footer}
    <div class="card-footer">{@render footer()}</div>
  {/if}
</div>

<style>
.card {
  border-radius: var(--radius-lg);
  background: var(--bg-layer);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-sm);
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  transition: background var(--transition-slow), border-color var(--transition-slow);
}
.card-content, .card-header, .card-footer { position: relative; z-index: 1; }

.card-header {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border);
}
.card-content { color: var(--text-primary); }
.card-footer {
  padding-top: 0.75rem;
  border-top: 1px solid var(--border);
  color: var(--text-secondary);
  font-size: 0.875rem;
}
</style>
