<script lang="ts">
  import { date, money, status } from '#lib/format.js';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
</script>

<SEO title="History transaksi" privatePage />
<div class="page-heading">
  <h1>Perjalanan tiketmu.</h1>
  <p>Lanjutkan pembayaran, cek konfirmasi, atau buka tiket yang sudah terbit.</p>
</div>
<div class="between narrow">
  <h2>History transaksi</h2>
  <a class="button" href="/tiket">Tiket saya</a>
</div>
<section class="narrow card order-list">
  {#each data.orders as o}<a href={'/transaksi/' + o.id}
      ><div class="between">
        <strong>{o.event.title}</strong><span class="pill">{status[o.status]}</span>
      </div>
      <p class="meta">{o.session.label} · {date(o.session.starts_at)} · {o.quantity} tiket</p>
      <strong>{money(o.total)}</strong></a
    >{:else}<Empty
      title="Belum ada transaksi"
      body="Temukan pertunjukan berikutnya dan pilih sesi yang ingin ditonton."
      href="/agenda"
      label="Buka agenda"
    />{/each}
</section>
