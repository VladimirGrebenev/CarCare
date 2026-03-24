<script lang="ts">
  import ErrorState from '../components/ui/ErrorState.svelte';
  export let error;
</script>
<ErrorState message={error?.message || 'Неизвестная ошибка'} />
