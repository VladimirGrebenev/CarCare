<script lang="ts">
  let { value = $bindable(''), label = '', type = 'text', placeholder = '', error = '', disabled = false, required = false, icon = null } = $props();
</script>
<div class="input-wrapper glassmorphism">
  {#if label}
    <label>{label}{required ? '*' : ''}</label>
  {/if}
  <div class="input-inner">
    {#if icon}
      <span class="input-icon">{@render icon()}</span>
    {/if}
    <input
      type={type}
      placeholder={placeholder}
      bind:value
      aria-label={label}
      aria-invalid={!!error}
      disabled={disabled}
      required={required}
      class:error={!!error}
    />
  </div>
  {#if error}
    <div class="input-error" role="alert">{error}</div>
  {/if}
</div>
<style>
.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
.input-inner {
  display: flex;
  align-items: center;
  background: var(--glass-bg);
  border-radius: 0.75rem;
  padding: 0.5rem 1rem;
  box-shadow: var(--glass-shadow);
}
input {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--text);
  font-size: 1rem;
  outline: none;
  padding: 0.5rem 0;
}
.input-icon {
  margin-right: 0.5rem;
  color: var(--accent);
}
.input-error {
  color: var(--danger);
  font-size: 0.85rem;
}
:global(.dark) .input-wrapper {
  --glass-bg: rgba(30, 30, 40, 0.6);
  --glass-shadow: 0 2px 16px 0 rgba(0,0,0,0.15);
  --text: #fff;
  --accent: #7de2fc;
  --danger: #ff6b6b;
}
</style>