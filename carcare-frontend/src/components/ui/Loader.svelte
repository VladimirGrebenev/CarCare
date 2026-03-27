<script lang="ts">
  type Props = { size?: number; className?: string; };
  let { size = 32, className = '' }: Props = $props();
</script>
<div class="loader {className}" style="width:{size}px;height:{size}px;" aria-label="Загрузка..." role="status"></div>
<style>
.loader {
  border: 3px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  margin: 0 auto;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
