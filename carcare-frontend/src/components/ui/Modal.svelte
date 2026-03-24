<script lang="ts">
  import type { Snippet } from 'svelte';

  type Props = {
    open?: boolean;
    title?: string;
    onClose?: (() => void) | null;
    children?: Snippet;
  };

  let { open = false, title = '', onClose = null, children = null }: Props = $props();

  function handleClose() {
    if (onClose) onClose();
    open = false;
  }

  function handleBackdropKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape' || event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleClose();
    }
  }

</script>
{#if open}
  <div
    class="modal-backdrop"
    tabindex="0"
    aria-modal="true"
    role="dialog"
    onclick={(event) => {
      if (event.target === event.currentTarget) {
        handleClose();
      }
    }}
    onkeydown={handleBackdropKeydown}
  >
    <div
      class="modal glassmorphism"
      role="document"
      tabindex="-1"
    >
      <header class="modal-header">
        <h2>{title}</h2>
        <button class="modal-close" aria-label="Close" onclick={handleClose}>&times;</button>
      </header>
      <div class="modal-content">{@render children?.()}</div>
    </div>
  </div>
{/if}
<style>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(20, 20, 30, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  border-radius: 1rem;
  padding: 2rem;
  min-width: 320px;
  max-width: 90vw;
  color: var(--text);
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}
.modal-close {
  background: none;
  border: none;
  font-size: 2rem;
  color: var(--accent);
  cursor: pointer;
}
:global(.dark) .modal {
  --glass-bg: rgba(30, 30, 40, 0.85);
  --glass-shadow: 0 2px 32px 0 rgba(0,0,0,0.22);
  --text: #fff;
  --accent: #7de2fc;
}
</style>