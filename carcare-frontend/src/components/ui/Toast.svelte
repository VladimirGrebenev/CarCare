<script lang="ts">
  let { message = '', type = 'info', open = $bindable(false), duration = 3000 } = $props();
  $effect(() => {
    if (open && duration > 0) {
      const timer = setTimeout(() => open = false, duration);
      return () => clearTimeout(timer);
    }
  });
</script>
{#if open}
  <div class="toast {type}" role="status" aria-live="polite">{message}</div>
{/if}
<style>
.toast {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  min-width: 220px;
  max-width: 90vw;
  padding: 1rem 2rem;
  border-radius: 0.75rem;
  background: var(--glass-bg);
  color: var(--text);
  box-shadow: var(--glass-shadow);
  z-index: 2000;
  font-size: 1rem;
  font-weight: 500;
  opacity: 0.98;
}
.toast.info { border-left: 4px solid #7de2fc; }
.toast.success { border-left: 4px solid #4ade80; }
.toast.error { border-left: 4px solid #ff6b6b; }
:global(.dark) .toast {
  --glass-bg: rgba(30, 30, 40, 0.9);
  --glass-shadow: 0 2px 32px 0 rgba(0,0,0,0.22);
  --text: #fff;
}
</style>
