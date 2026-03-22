<script lang="ts">
  let { message = 'Произошла ошибка', icon = null, children, onRetry = null } = $props();
</script>
<div class="error-state glassmorphism" role="alert">
  {#if icon}
    <div class="error-icon">{@render icon()}</div>
  {/if}
  <div class="error-message">{message}</div>
  {@render children()}
  {#if onRetry}
    <button class="retry-btn" onclick={onRetry}>Повторить</button>
  {/if}
</div>
<style>
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  border-radius: 1rem;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  color: var(--danger);
  gap: 1rem;
}
.error-icon {
  font-size: 2.5rem;
  color: var(--danger);
}
.error-message {
  font-size: 1.1rem;
  font-weight: 500;
}
.retry-btn {
  margin-top: 1rem;
  padding: 0.5rem 1.5rem;
  border-radius: 0.75rem;
  border: none;
  background: #7de2fc;
  color: #222;
  font-weight: 600;
  cursor: pointer;
}
:global(.dark) .error-state {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --danger: #ff6b6b;
}
</style>
