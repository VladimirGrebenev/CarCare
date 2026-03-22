<script lang="ts">
  export let error;
</script>
<ErrorState message={error?.message || 'Неизвестная ошибка'} />
