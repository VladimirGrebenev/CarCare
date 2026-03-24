<script lang="ts">
  import PageLayout from '../../components/layout/PageLayout.svelte';
  import EmptyState from '../../components/ui/EmptyState.svelte';
  import ErrorState from '../../components/ui/ErrorState.svelte';
  import Loader from '../../components/ui/Loader.svelte';
  import Table from '../../components/ui/Table.svelte';
  import { reportsStore } from '../../stores/reports';

  let filters = $state({ period: '', type: '' });

  $effect(() => {
    reportsStore.load(filters);
  });

  function onExport(format: string) {
    reportsStore.export(format);
  }
</script>

<PageLayout title="Отчёты">
  <form class="filters" onsubmit={(event) => event.preventDefault()}>
    <input type="text" placeholder="Период" bind:value={filters.period} />
    <input type="text" placeholder="Тип" bind:value={filters.type} />
    <button type="button" onclick={() => reportsStore.load(filters)}>Фильтровать</button>
    <button type="button" onclick={() => onExport('csv')}>Экспорт CSV</button>
    <button type="button" onclick={() => onExport('xlsx')}>Экспорт XLSX</button>
  </form>

  {#if reportsStore.loading}
    <Loader />
  {:else if reportsStore.error}
    <ErrorState message={reportsStore.error} />
  {:else if !reportsStore.items.length}
    <EmptyState message="Нет отчётов" />
  {:else}
    <Table columns={[{label:'Период',key:'period'},{label:'Тип',key:'type'},{label:'Сумма',key:'amount'}]} rows={reportsStore.items} />
  {/if}
</PageLayout>
