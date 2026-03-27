<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLInputAttributes } from 'svelte/elements';

  type Props = Omit<HTMLInputAttributes, 'value'> & {
    value?: string;
    label?: string;
    error?: string;
    hint?: string;
    icon?: Snippet | null;
    className?: string;
    inputProps?: Record<string, unknown>;
  };

  let {
    value = $bindable(''),
    label = '',
    type = 'text',
    placeholder = '',
    error = '',
    hint = '',
    disabled = false,
    required = false,
    icon = null,
    id = '',
    className = '',
    class: classAttr = '',
    inputProps = {}
  }: Props = $props();

  const inputId = $derived(
    id || `input-${String(label || 'field').toLowerCase().replace(/[^a-z0-9а-яё]+/gi, '-')}`
  );
</script>

<div class="input-wrapper {className} {classAttr}">
  {#if label}
    <label class="input-label" for={inputId}>
      {label}{#if required}<span class="required">*</span>{/if}
    </label>
  {/if}
  <div class="input-inner" class:input-error-state={!!error} class:input-disabled={disabled}>
    {#if icon}
      <span class="input-icon">{@render icon()}</span>
    {/if}
    <input
      id={inputId}
      {type}
      {placeholder}
      bind:value
      aria-label={label || placeholder}
      aria-invalid={!!error}
      aria-describedby={error ? `${inputId}-error` : hint ? `${inputId}-hint` : undefined}
      {disabled}
      {required}
      class="input-field"
      {...inputProps}
    />
  </div>
  {#if error}
    <div id="{inputId}-error" class="input-message input-message-error" role="alert">{error}</div>
  {:else if hint}
    <div id="{inputId}-hint" class="input-message">{hint}</div>
  {/if}
</div>

<style>
.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.input-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
}
.required { color: var(--danger); margin-left: 2px; }

.input-inner {
  display: flex;
  align-items: center;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 0 0.875rem;
  transition: border-color var(--transition), box-shadow var(--transition);
}
.input-inner:focus-within {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
}
.input-inner.input-error-state {
  border-color: var(--danger);
}
.input-inner.input-error-state:focus-within {
  box-shadow: 0 0 0 3px var(--danger-light);
}
.input-inner.input-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.input-field {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-family: var(--font);
  font-size: 0.9375rem;
  outline: none;
  padding: 0.5625rem 0;
  min-width: 0;
}
.input-field::placeholder { color: var(--text-disabled); }
.input-field:disabled { cursor: not-allowed; }

.input-icon {
  display: flex;
  align-items: center;
  color: var(--text-secondary);
  margin-right: 0.5rem;
  flex-shrink: 0;
}

.input-message { font-size: 0.8125rem; color: var(--text-secondary); }
.input-message-error { color: var(--danger); }
</style>
