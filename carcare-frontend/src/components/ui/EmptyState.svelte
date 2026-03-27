<script lang="ts">
  import type { Snippet } from 'svelte';

  type Props = {
    children?: Snippet;
    message?: string;
    icon?: Snippet | null;
    class?: string;
  };

  let { children = null, message = 'Нет данных', icon = null, class: className = '' }: Props = $props();
</script>

<div class="empty-state {className}">
  {#if icon}
    <div class="empty-icon">{@render icon()}</div>
  {/if}
  <p class="empty-message">{message}</p>
  {@render children?.()}
</div>

<style>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 1.5rem;
  gap: 0.75rem;
  color: var(--text-secondary);
  text-align: center;
}
.empty-icon { font-size: 2.5rem; opacity: 0.6; }
.empty-message { font-size: 0.9375rem; font-weight: 500; }
</style>
