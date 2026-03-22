<script lang="ts">
  let { type = 'button', variant = 'primary', disabled = false, loading = false, icon = null, children } = $props();
</script>
<button
  type={type}
  class="btn {variant} glassmorphism"
  disabled={disabled || loading}
  aria-busy={loading}
>
  {#if loading}
    <span class="loader"></span>
  {/if}
  {#if icon}
    <span class="btn-icon">{@render icon()}</span>
  {/if}
  {@render children()}
</button>
<style>
.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 0.75rem;
  border: none;
  background: var(--glass-bg);
  color: var(--text);
  font-size: 1rem;
  font-weight: 600;
  box-shadow: var(--glass-shadow);
  cursor: pointer;
  transition: background 0.2s;
}
.btn.primary {
  background: linear-gradient(90deg, #7de2fc 0%, #b9b6e5 100%);
  color: #222;
}
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn-icon {
  display: flex;
  align-items: center;
}
.loader {
  width: 1em;
  height: 1em;
  border: 2px solid #fff;
  border-top: 2px solid #7de2fc;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
:global(.dark) .btn {
  --glass-bg: rgba(30, 30, 40, 0.6);
  --glass-shadow: 0 2px 16px 0 rgba(0,0,0,0.15);
  --text: #fff;
}
</style>