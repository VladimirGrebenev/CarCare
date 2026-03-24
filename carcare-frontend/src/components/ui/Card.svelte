<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLAttributes } from 'svelte/elements';

  type Props = HTMLAttributes<HTMLDivElement> & {
    children?: Snippet;
    header?: Snippet | null;
    footer?: Snippet | null;
    className?: string;
  };

  let {
    children = null,
    header = null,
    footer = null,
    className = '',
    class: classAttr = ''
  }: Props = $props();
</script>
<div class={`card glassmorphism ${className} ${classAttr}`}>
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
  border-radius: 1rem;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  padding: 1.5rem;
  color: var(--text);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.card-header, .card-footer {
  font-weight: 600;
  color: var(--accent);
}
:global(.dark) .card {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --text: #fff;
  --accent: #7de2fc;
}
</style>