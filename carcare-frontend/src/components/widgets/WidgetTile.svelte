<script lang="ts">
  import type { Snippet } from 'svelte';

  export let title: string = '';
  export let icon: Snippet | null = null;
  export let value: string = '';
  export let onClick: (() => void) | null = null;
  export let className: string = '';
  export let children: Snippet | null = null;

  function handleKeydown(event: KeyboardEvent) {
    if (!onClick) return;
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      onClick();
    }
  }
</script>
<div
  class="widget-tile glassmorphism {className}"
  tabindex="0"
  role="button"
  onclick={onClick}
  onkeydown={handleKeydown}
>
  {#if icon}
    <div class="widget-icon">{@render icon()}</div>
  {/if}
  <div class="widget-content">
    <div class="widget-title">{title}</div>
    <div class="widget-value">{value}</div>
    {@render children?.()}
  </div>
</div>
<style>
.widget-tile {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem 1rem;
  border-radius: 1rem;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  color: var(--text);
  cursor: pointer;
  transition: background 0.2s;
  outline: none;
}
.widget-tile:focus {
  box-shadow: 0 0 0 2px #7de2fc;
}
.widget-icon {
  font-size: 2rem;
  color: var(--accent);
}
.widget-title {
  font-size: 1.1rem;
  font-weight: 600;
}
.widget-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--accent);
}
:global(.dark) .widget-tile {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --text: #fff;
  --accent: #7de2fc;
}
</style>
