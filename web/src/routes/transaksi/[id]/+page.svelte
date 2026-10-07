<script lang="ts">
  import OrderDetail from '#lib/components/OrderDetail.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
</script>

<SEO title="Detail transaksi" privatePage /><OrderDetail
  order={data.order}
  methods={data.methods}
/>
