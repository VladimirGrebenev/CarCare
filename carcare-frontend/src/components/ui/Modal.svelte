<script lang="ts">
  import type { Snippet } from 'svelte';

  type Props = {
    open?: boolean;
    title?: string;
    onClose?: (() => void) | null;
    children?: Snippet;
    footer?: Snippet | null;
    width?: string;
  };

  let { open = false, title = '', onClose = null, children = null, footer = null, width = '480px' }: Props = $props();

  function handleClose() {
    if (onClose) onClose();
  }

  function handleBackdropKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') { e.preventDefault(); handleClose(); }
  }
</script>

{#if open}
  <div
    class="modal-backdrop"
    role="dialog"
    aria-modal="true"
    aria-label={title}
    tabindex="-1"
    onclick={(e) => { if (e.target === e.currentTarget) handleClose(); }}
    onkeydown={handleBackdropKeydown}
  >
    <div class="modal-panel" style="max-width: {width};" role="document">
      {#if title}
        <header class="modal-header">
          <h2 class="modal-title">{title}</h2>
          <button class="modal-close" aria-label="Закрыть" onclick={handleClose}>✕</button>
        </header>
      {/if}
      <div class="modal-body">{@render children?.()}</div>
      {#if footer}
        <footer class="modal-footer">{@render footer()}</footer>
      {/if}
    </div>
  </div>
{/if}

<style>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
  animation: backdrop-in 0.15s ease;
}

.modal-panel {
  background: var(--bg-layer);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  width: 100%;
  display: flex;
  flex-direction: column;
  animation: modal-in 0.15s ease;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem 1rem;
  border-bottom: 1px solid var(--border);
}
.modal-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
}
.modal-close {
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background var(--transition), color var(--transition);
  flex-shrink: 0;
}
.modal-close:hover { background: var(--danger-light); color: var(--danger); border-color: var(--danger); }

.modal-body { padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
.modal-footer {
  padding: 1rem 1.5rem 1.25rem;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

@keyframes backdrop-in { from { opacity: 0; } to { opacity: 1; } }
@keyframes modal-in { from { opacity: 0; transform: scale(0.96) translateY(-8px); } to { opacity: 1; transform: none; } }
</style>
