<script lang="ts">
  import { date, time } from '#lib/format.js';
  import QR from '#lib/components/QR.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
</script>

<SEO title="Tiket saya" privatePage />
<div class="page-heading">
  <h1>Sampai jumpa di panggung.</h1>
  <p>Tiketmu ada di sini. Tunjukkan QR saat masuk.</p>
</div>
<section class="narrow stack">
  {#each data.tickets as t}<article class="card pad ticket">
      <span class="pill">{t.checked_at ? 'Sudah digunakan' : 'Siap digunakan'}</span>
      <h2>{t.event.title}</h2>
      <p>{t.session.label}<br />{date(t.session.starts_at)} · {time(t.session.starts_at)}</p>
      <QR code={t.qr} /><a class="button" href={'/transaksi/' + t.order_id}>Detail pesanan</a>
    </article>{:else}<Empty
      title="Belum ada tiket"
      body="Tiket diterbitkan setelah pembayaran dikonfirmasi."
      href="/transaksi"
      label="Lihat transaksi"
    />{/each}
</section>
