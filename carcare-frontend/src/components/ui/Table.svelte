<script lang="ts">
  import Loader from './Loader.svelte';

  export let columns: { label: string; key: string }[] = [];
  export let rows: Array<Record<string, unknown>> = [];
  export let emptyText: string = 'Нет данных';
  export let loading: boolean = false;
  export let error: string = '';
  export let className: string = '';
</script>
<div class="table-wrapper {className}">
  {#if loading}
    <div class="table-loader"><Loader /></div>
  {:else if error}
    <div class="table-error" role="alert">{error}</div>
  {:else if rows.length === 0}
    <div class="table-empty">{emptyText}</div>
  {:else}
    <table>
      <thead>
        <tr>
          {#each columns as col}
            <th>{col.label}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each rows as row}
          <tr>
            {#each columns as col}
              <td>{row[col.key]}</td>
            {/each}
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
  background: var(--glass-bg);
  border-radius: 1rem;
  box-shadow: var(--glass-shadow);
  padding: 1rem;
}
table {
  width: 100%;
  border-collapse: collapse;
  color: var(--text);
}
th, td {
  padding: 0.75rem 1rem;
  text-align: left;
}
th {
  color: var(--accent);
  font-weight: 600;
  background: rgba(125,226,252,0.08);
}
tr:nth-child(even) {
  background: rgba(30,30,40,0.04);
}
.table-loader, .table-error, .table-empty {
  padding: 2rem;
  text-align: center;
  color: var(--accent);
}
:global(.dark) .table-wrapper {
  --glass-bg: rgba(30, 30, 40, 0.7);
  --glass-shadow: 0 2px 24px 0 rgba(0,0,0,0.18);
  --text: #fff;
  --accent: #7de2fc;
}
</style>
