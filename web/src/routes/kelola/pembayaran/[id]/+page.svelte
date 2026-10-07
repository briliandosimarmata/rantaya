<script lang="ts">
  import OrderDetail from '#lib/components/OrderDetail.svelte';
  let { data } = $props();
</script>

<OrderDetail order={data.order} reviewer />
