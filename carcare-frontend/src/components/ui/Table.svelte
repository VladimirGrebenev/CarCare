<script lang="ts">
  import type { Snippet } from 'svelte';
  import Loader from './Loader.svelte';

  type Column = { label: string; key: string };
  type SortDir = 'asc' | 'desc';
  type SortState = { key: string; dir: SortDir };

  type Props = {
    columns?: Column[];
    rows?: Array<Record<string, unknown>>;
    emptyText?: string;
    loading?: boolean;
    error?: string;
    className?: string;
    onRowClick?: ((row: Record<string, unknown>) => void) | null;
    actions?: Snippet<[Record<string, unknown>]> | null;
    sortKeys?: string[];
    sort?: SortState[];
    onSort?: ((key: string, event: MouseEvent) => void) | null;
  };

  let {
    columns = [],
    rows = [],
    emptyText = 'Нет данных',
    loading = false,
    error = '',
    className = '',
    onRowClick = null,
    actions = null,
    sortKeys = [],
    sort = [],
    onSort = null
  }: Props = $props();

  function getSortState(key: string): SortState | undefined {
    return sort.find(s => s.key === key);
  }

  function getSortRank(key: string): number | null {
    const idx = sort.findIndex(s => s.key === key);
    return idx === -1 ? null : idx + 1;
  }

  function isSortable(key: string): boolean {
    return sortKeys.includes(key);
  }
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
            {#if isSortable(col.key) && onSort}
              {@const state = getSortState(col.key)}
              {@const rank = getSortRank(col.key)}
              {@const isPrimary = rank === 1}
              {@const isSecondary = rank === 2}
              <th
                scope="col"
                class="sortable"
                class:sort-primary={isPrimary}
                class:sort-secondary={isSecondary}
                onclick={(e) => onSort?.(col.key, e)}
                title="Сортировать по «{col.label}» (Shift+клик — вторичная сортировка)"
              >
                <span class="th-content">
                  {col.label}
                  <span class="sort-icon" class:sort-icon-primary={isPrimary} class:sort-icon-secondary={isSecondary}>
                    {#if !state}
                      ↕
                    {:else if state.dir === 'desc'}
                      ↓
                    {:else}
                      ↑
                    {/if}
                    {#if isSecondary}
                      <sup class="sort-rank">2</sup>
                    {/if}
                  </span>
                </span>
              </th>
            {:else}
              <th scope="col">{col.label}</th>
            {/if}
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
  padding: 0.625rem 0.75rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  white-space: nowrap;
}

th.sortable {
  cursor: pointer;
  user-select: none;
  transition: background var(--transition, 0.15s ease);
}

th.sortable:hover {
  background: rgba(255, 255, 255, 0.06);
}

th.sort-primary {
  color: var(--accent-text, var(--accent));
}

th.sort-secondary {
  color: var(--text-secondary);
}

.th-content {
  display: inline-flex;
  align-items: center;
  gap: 0.3em;
}

.sort-icon {
  font-size: 0.85em;
  color: var(--text-secondary);
  opacity: 0.5;
  display: inline-flex;
  align-items: center;
  gap: 0.1em;
}

.sort-icon-primary {
  color: var(--accent-text, var(--accent));
  opacity: 1;
}

.sort-icon-secondary {
  color: var(--text-secondary);
  opacity: 0.8;
}

.sort-rank {
  font-size: 0.7em;
  line-height: 1;
  vertical-align: super;
}

td {
  padding: 0.625rem 0.75rem;
  border-top: 1px solid var(--border);
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
}

tr.clickable { cursor: pointer; }
tr.clickable:hover td { background: var(--accent-light); }
tr:focus-visible td { background: var(--accent-light); }

.col-actions { width: 1%; white-space: nowrap; padding: 0.375rem 0.5rem; }

.table-state {
  padding: 3rem 1rem;
  text-align: center;
  color: var(--text-secondary);
}
.table-error { color: var(--danger); }
</style>
