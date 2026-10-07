<script lang="ts">
  import type { ArtsEvent, Organizer } from '#lib/types.js';
  import { date, time, money } from '#lib/format.js';
  import Follow from './Follow.svelte';
  import Avatar from './Avatar.svelte';
  import Icon from './Icon.svelte';
  let { events, organizers } = $props<{ events: ArtsEvent[]; organizers: Organizer[] }>();
</script>
{#snippet tile(s: string)}<span class="date-tile"><span>{date(s, { month: 'short' }).toUpperCase()}</span><b>{date(s, { day: 'numeric' })}</b></span>{/snippet}
<aside class="right-column discovery-aside">
  {#if events[0]}{@const e = events[0]}<div class="card next-event">
    <div class="between"><span class="eyebrow">Pentas berikutnya</span>{@render tile(e.starts_at)}</div>
    <h3>{e.title}</h3><div class="meta">{e.organizer.name} · {time(e.starts_at)}</div>
    <div class="info-line"><Icon name="pin" /><span>{e.venue}</span></div>
    <a class="button lime full" href={'/event/' + e.slug}>Lihat pertunjukan</a>
  </div>{/if}
  <section><div class="aside-title"><h3>Catat tanggalnya</h3><a class="link" href="/agenda">Semua</a></div>
    {#each events.slice(1, 3) as e}<a class="side-event" href={'/event/' + e.slug}>{@render tile(e.starts_at)}<div><h3>{e.title}</h3><div class="meta">{e.category} · {money(e.price)}</div><div class="meta">{e.venue}</div></div></a>{:else}<p class="small muted">Belum ada agenda tambahan.</p>{/each}
  </section>
  <section><div class="aside-title"><h3>Temukan ruangmu</h3><a href="/ruang" class="link">Jelajah</a></div>
    {#each organizers as o}<div class="community-mini"><a href={'/ruang/' + o.slug}><Avatar name={o.name} src={o.avatar_url} size="square" /></a><div class="community-text"><a href={'/ruang/' + o.slug}><b>{o.name}</b></a><div class="meta">{o.category} · {o.followers} pengikut</div></div><Follow id={o.id} small /></div>{/each}
  </section>
  <div class="aside-note">Obrolan tetap terbuka, bahkan setelah lampu panggung padam.<br /><a href="/profil?tab=settings" class="link">Atur kabar yang kamu terima</a></div>
</aside>
