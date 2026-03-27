<script lang="ts">
  import type { Snippet } from 'svelte';
  import Loader from './Loader.svelte';

  type Column = { label: string; key: string };

  type Props = {
    columns?: Column[];
    rows?: Array<Record<string, unknown>>;
    emptyText?: string;
    loading?: boolean;
    error?: string;
    className?: string;
    onRowClick?: ((row: Record<string, unknown>) => void) | null;
    actions?: Snippet<[Record<string, unknown>]> | null;
  };

  let {
    columns = [],
    rows = [],
    emptyText = 'Нет данных',
    loading = false,
    error = '',
    className = '',
    onRowClick = null,
    actions = null
  }: Props = $props();
</script>

<div class="table-wrapper {className}">
  {#if loading}
    <div class="table-state"><Loader /></div>
  {:else if error}
    <div class="table-state table-error" role="alert">{error}</div>
  {:else if rows.length === 0}
    <div class="table-state table-empty">{emptyText}</div>
  {:else}
    <table>
      <thead>
        <tr>
          {#each columns as col}
            <th scope="col">{col.label}</th>
          {/each}
          {#if actions}<th scope="col" class="col-actions">Действия</th>{/if}
        </tr>
      </thead>
      <tbody>
        {#each rows as row}
          <tr
            class:clickable={!!onRowClick}
            onclick={() => onRowClick?.(row)}
            onkeydown={(e) => { if (e.key === 'Enter') onRowClick?.(row); }}
            tabindex={onRowClick ? 0 : undefined}
            role={onRowClick ? 'button' : undefined}
          >
            {#each columns as col}
              <td>{row[col.key] ?? '—'}</td>
            {/each}
            {#if actions}
              <td class="col-actions" onclick={(e) => e.stopPropagation()}>
                {@render actions(row)}
              </td>
            {/if}
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
.table-wrapper {
  width: 100%;
  overflow-x: auto;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  background: var(--bg-layer);
}

table {
  width: 100%;
  border-collapse: collapse;
  color: var(--text-primary);
  font-size: 0.9375rem;
}

thead { background: var(--bg-input); }
th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  white-space: nowrap;
}

td {
  padding: 0.75rem 1rem;
  border-top: 1px solid var(--border);
  color: var(--text-primary);
}

tr.clickable { cursor: pointer; }
tr.clickable:hover td { background: var(--accent-light); }
tr:focus-visible td { background: var(--accent-light); }

.col-actions { width: 1%; white-space: nowrap; }

.table-state {
  padding: 3rem 1rem;
  text-align: center;
  color: var(--text-secondary);
}
.table-error { color: var(--danger); }
</style>
