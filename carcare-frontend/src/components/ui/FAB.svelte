<script lang="ts">
  import type { Snippet } from 'svelte';

  type Props = {
    icon?: Snippet | null;
    label?: string;
    onClick?: (() => void) | null;
    class?: string;
    className?: string;
    position?: 'fixed' | 'relative';
  };

  let {
    icon = null,
    label = '',
    onClick = null,
    class: classAttr = '',
    className = '',
    position = 'fixed'
  }: Props = $props();
</script>

<button
  class="fab {className} {classAttr}"
  class:fab-fixed={position === 'fixed'}
  aria-label={label}
  onclick={onClick}
>
  {#if icon}
    <span class="fab-icon">{@render icon()}</span>
  {/if}
  {#if label}
    <span class="fab-label">{label}</span>
  {/if}
</button>

<style>
.fab {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.875rem 1.375rem;
  border-radius: 100px;
  background: var(--accent);
  color: #fff;
  border: none;
  box-shadow: var(--shadow-md);
  font-family: var(--font);
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  z-index: 200;
  transition: background var(--transition), box-shadow var(--transition), transform var(--transition);
}
.fab:hover {
  background: var(--accent-hover);
  box-shadow: var(--shadow-lg);
  transform: translateY(-1px);
}
.fab:active { transform: translateY(0); }

.fab-fixed {
  position: fixed;
  right: 2rem;
  bottom: 2rem;
}

.fab-icon { font-size: 1.25rem; }
</style>
