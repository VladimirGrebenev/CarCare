<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLButtonAttributes } from 'svelte/elements';

  type Props = HTMLButtonAttributes & {
    children?: Snippet;
    variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
    loading?: boolean;
    icon?: Snippet | null;
  };

  let {
    children = null,
    type = 'button',
    variant = 'primary',
    disabled = false,
    loading = false,
    icon = null,
    class: className = '',
    onclick
  }: Props = $props();
</script>

<button
  {type}
  class="btn btn-{variant} {className}"
  disabled={disabled || loading}
  aria-busy={loading}
  {onclick}
>
  {#if loading}
    <span class="btn-spinner" aria-hidden="true"></span>
  {/if}
  {#if icon && !loading}
    <span class="btn-icon">{@render icon()}</span>
  {/if}
  {@render children?.()}
</button>

<style>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.5625rem 1.25rem;
  border-radius: var(--radius-md);
  border: 1px solid transparent;
  font-family: var(--font);
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition), border-color var(--transition), opacity var(--transition), box-shadow var(--transition);
  white-space: nowrap;
}

.btn-primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}
.btn-primary:hover:not(:disabled) {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
  box-shadow: var(--shadow-sm);
}

.btn-secondary {
  background: var(--bg-input);
  color: var(--text-primary);
  border-color: var(--border);
}
.btn-secondary:hover:not(:disabled) {
  background: var(--accent-light);
  border-color: var(--accent);
  color: var(--accent-text);
}

.btn-danger {
  background: var(--danger-light);
  color: var(--danger);
  border-color: var(--danger);
}
.btn-danger:hover:not(:disabled) {
  background: var(--danger);
  color: #fff;
}

.btn-ghost {
  background: transparent;
  color: var(--text-secondary);
  border-color: transparent;
}
.btn-ghost:hover:not(:disabled) {
  background: var(--accent-light);
  color: var(--text-primary);
}

.btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.btn-icon {
  display: flex;
  align-items: center;
  font-size: 1em;
}

.btn-spinner {
  width: 1em;
  height: 1em;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: btn-spin 0.7s linear infinite;
  flex-shrink: 0;
}
@keyframes btn-spin {
  to { transform: rotate(360deg); }
}
</style>
