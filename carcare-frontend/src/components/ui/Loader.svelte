<script lang="ts">
  let { size = 32, color = '#7de2fc', className = '' } = $props();
</script>
<div class="loader {className}" style="width:{size}px;height:{size}px;border-top-color:{color};"></div>
<style>
.loader {
  border: 4px solid rgba(125,226,252,0.2);
  border-top: 4px solid #7de2fc;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
