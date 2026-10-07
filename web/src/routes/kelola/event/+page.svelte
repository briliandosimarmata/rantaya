<script lang="ts">
  import { date, ended } from '#lib/format.js';
  import Empty from '#lib/components/Empty.svelte';
  import Icon from '#lib/components/Icon.svelte';
  let { data } = $props();
</script>
<div class="between" style="margin-bottom:20px"><h2>Event kamu</h2><a class="button primary" href="/kelola/event/baru"><Icon name="plus" />Buat event</a></div>
{#if data.events.length}<div class="card table-wrap"><table class="table"><thead><tr><th>Event</th><th>Tanggal</th><th>Status</th><th>Tindakan</th></tr></thead><tbody>
{#each data.events as e}<tr><td><a href={'/event/' + e.slug}><b>{e.title}</b></a></td><td>{date(e.starts_at)}</td><td><span class="pill">{!e.published ? 'Draft' : ended(e) ? 'Selesai' : 'Mendatang'}</span></td><td><a class="button slim" href={'/kelola/event/' + e.id}><Icon name="edit" />Edit</a></td></tr>{/each}
</tbody></table></div>{:else}<Empty title="Belum ada event" body="Mulai dengan informasi utama, jadwal, dan kapasitas." />{/if}
