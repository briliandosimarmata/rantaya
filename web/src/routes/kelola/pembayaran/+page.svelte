<script lang="ts">
  import { money, date, status } from '#lib/format.js';
  import Empty from '#lib/components/Empty.svelte';
  let { data } = $props();
  let tab = $state('pending');
  let items = $derived(
    tab === 'pending'
      ? data.orders.filter((o: any) => o.status === 'awaiting_review')
      : data.orders,
  );
</script>

<div class="page-heading">
  <h1>Konfirmasi pembayaran.</h1>
  <p>Periksa dana masuk sebelum menyetujui bukti transfer.</p>
</div>
<nav class="tabs" aria-label="Status pembayaran">
  <button class:active={tab === 'pending'} onclick={() => (tab = 'pending')}>Perlu diperiksa</button
  ><button class:active={tab === 'all'} onclick={() => (tab = 'all')}>Semua transaksi</button>
</nav>
<section class="card order-list">
  {#each items as o}<a href={'/kelola/pembayaran/' + o.id}
      ><div class="between">
        <strong>{o.event.title}</strong><span class="pill">{status[o.status]}</span>
      </div>
      <p class="meta">{o.customer.name} · {o.quantity} tiket · {date(o.created_at)}</p>
      <strong>{money(o.total)}</strong></a
    >{:else}<Empty
      title="Tidak ada pembayaran yang menunggu"
      body="Pengelola mendapat notifikasi ketika penonton mengirim bukti transfer."
    />{/each}
</section>
