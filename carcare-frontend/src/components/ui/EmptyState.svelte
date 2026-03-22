<script lang="ts">
  let { message = 'Нет данных', icon = null, children } = $props();
</script>
<div class="empty-state glassmorphism">
  {#if icon}
    <div class="empty-icon">{@render icon()}</div>
  {/if}
  <div class="empty-message">{message}</div>
  {@render children()}
</div>
<style>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  border-radius: 1rem;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  color: var(--text);
  gap: 1rem;
}
.empty-icon {
  font-size: 2.5rem;
  color: var(--accent);
}
.empty-message {
  font-size: 1.1rem;
  font-weight: 500;
}
:global(.dark) .empty-state {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --text: #fff;
  --accent: #7de2fc;
}
</style>
