<script lang="ts">
  import type { Snippet } from 'svelte';

  export let icon: Snippet | null = null;
  export let label: string = '';
  export let onClick: (() => void) | null = null;
  export let className: string = '';
</script>
<button class="fab glassmorphism {className}" aria-label={label} onclick={onClick}>
  {#if icon}
    <span class="fab-icon">{@render icon()}</span>
  {/if}
  {#if label}
    <span class="fab-label">{label}</span>
  {/if}
</button>
<style>
.fab {
  position: fixed;
  right: 2rem;
  bottom: 2rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  border-radius: 2rem;
  background: var(--glass-bg);
  color: var(--accent);
  box-shadow: var(--glass-shadow);
  font-size: 1.2rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  z-index: 1200;
  transition: background 0.2s;
}
.fab-icon {
  font-size: 1.5rem;
}
:global(.dark) .fab {
  --glass-bg: rgba(30, 30, 40, 0.8);
  --glass-shadow: 0 2px 32px 0 rgba(0,0,0,0.22);
  --accent: #7de2fc;
}
</style>
