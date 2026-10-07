<script lang="ts">
  import { money } from '#lib/format.js';
  import Icon from '#lib/components/Icon.svelte';
  let { data } = $props();
  let activity = $derived((data.stats.activity || []).map((d: any) => ({...d, label: new Date(d.date + 'T12:00:00+07:00').toLocaleDateString('id-ID',{weekday:'short',timeZone:'Asia/Jakarta'})})));
  let peak = $derived(Math.max(1,...activity.map((d: any) => d.count)));
</script>
<div class="metrics">
  {#each [['Pengikut', data.stats.followers, 'Anggota yang mengikuti ruang'], ['Interaksi', data.stats.interactions, 'Dukungan dan percakapan'], ['Klik tiket', data.stats.ticket_clicks, 'Klik belum berarti pembelian'], ['Klik merchandise', data.stats.merch_clicks, 'Menuju kanal pemesanan']] as [label,n,description], i}
    <article class="card metric" class:highlight={i===0}><span class="small">{label}</span><b>{n}</b><span class="meta">{description}</span></article>
  {/each}
</div>
<div class="content-grid">
  <section class="card pad"><div class="between"><h2>Aktivitas ruang</h2><span class="meta">7 hari terakhir</span></div>
    <div class="chart" role="img" aria-label={'Aktivitas per hari: ' + activity.map((d: any) => d.label + ' ' + d.count).join(', ')}>
      {#each activity as day}<div class="chart-col"><div class="chart-bar" style:height={(day.count / peak * 100) + '%'} title={day.count + ' aktivitas'}></div>{day.label}</div>{/each}
    </div>
  </section>
  <aside class="right-column"><div class="card pad"><h3>Kelola ruangmu</h3><div class="divider"></div>
    <a class="button full" href="/kelola/event/baru"><Icon name="plus" />Buat event</a>
    <a class="button full" style="margin-top:10px" href="/kelola/merch/baru"><Icon name="bag" />Tambah merchandise</a>
    <a class="button full" style="margin-top:10px" href={'/ruang/' + data.organizer.slug}>Lihat ruang komunitas</a>
  </div></aside>
</div>
<div class="card pad" style="margin-top:24px"><div class="between"><div><h2>Penjualan yang disetujui</h2><strong>{money(data.stats.revenue)}</strong></div><a class="button primary" href="/kelola/pembayaran">Periksa pembayaran</a></div><p class="meta" style="margin-top:12px">{data.stats.pending_payments} perlu konfirmasi · {data.stats.tickets} tiket terbit · {data.stats.checked_in} sudah masuk</p></div>
