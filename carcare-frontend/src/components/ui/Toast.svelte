<script lang="ts">
  type Props = {
    message?: string;
    type?: 'info' | 'success' | 'error' | 'warning';
    show?: boolean;
    open?: boolean;
    duration?: number;
  };

  let { message = '', type = 'info', show = false, open = false, duration = 3000 }: Props = $props();
  const visible = $derived(show || open);

  let timer: ReturnType<typeof setTimeout> | undefined;
  $effect(() => {
    if (visible && duration > 0) {
      clearTimeout(timer);
      timer = setTimeout(() => { show = false; open = false; }, duration);
    }
  });

  const icons: Record<string, string> = {
    info: 'ℹ️', success: '✅', error: '❌', warning: '⚠️'
  };
</script>

{#if visible}
  <div class="toast toast-{type}" role="status" aria-live="polite">
    <span class="toast-icon">{icons[type] ?? 'ℹ️'}</span>
    <span class="toast-message">{message}</span>
  </div>
{/if}

<style>
.toast {
  position: fixed;
  top: 1.5rem;
  right: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 240px;
  max-width: 420px;
  padding: 0.875rem 1.25rem;
  border-radius: var(--radius-lg);
  background: var(--bg-layer);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-lg);
  z-index: 2000;
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--text-primary);
  animation: toast-in 0.2s ease;
}

.toast-info    { border-left: 3px solid var(--accent); }
.toast-success { border-left: 3px solid var(--success); }
.toast-error   { border-left: 3px solid var(--danger); }
.toast-warning { border-left: 3px solid var(--warning); }

.toast-icon { font-size: 1.125rem; flex-shrink: 0; }
.toast-message { flex: 1; }

@keyframes toast-in {
  from { opacity: 0; transform: translateX(16px); }
  to   { opacity: 1; transform: none; }
}
</style>
