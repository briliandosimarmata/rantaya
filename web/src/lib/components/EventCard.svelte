<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext, ArtsEvent } from '#lib/types.js';
  import { date, time, money, ended } from '#lib/format.js';
  import { mutate } from '#lib/api.js';
  import Icon from './Icon.svelte';
  let { event: e, compact = false } = $props<{ event: ArtsEvent; compact?: boolean }>();
  const ctx = getContext<AppContext>('app');
  let saved = $derived(ctx.account.bookmarks.some((b) => b.kind === 'event' && b.id === e.id));
  async function bookmark() {
    if (!ctx.requireLogin()) return;
    try {
      await mutate('/bookmarks', { kind: 'event', id: e.id, active: !saved });
      await ctx.refresh();
    } catch (err: any) {
      ctx.toast(err.message);
    }
  }
</script>

{#if compact}<a class="next-event" href={'/event/' + e.slug}>
    <span class="date-tile"
      ><strong>{date(e.starts_at, { day: 'numeric' })}</strong><span
        >{date(e.starts_at, { month: 'short' })}</span
      ></span
    >
    <span><strong>{e.title}</strong><span class="meta">{time(e.starts_at)} · {e.venue}</span></span>
  </a>{:else}<article class="event-card card">
    <a class="event-poster" href={'/event/' + e.slug} aria-label={'Detail ' + e.title}>
      {#if e.flyer_url}<img src={e.flyer_url} alt="" loading="lazy" width="260" height="340" />{/if}
      <span class="eyebrow">{e.category} / {e.city}</span><span class="poster-title">{e.title}</span
      >
      <span class="small">{date(e.starts_at, { day: 'numeric', month: 'short' })}</span>
    </a>
    <div class="event-card-info">
      <span class="pill" class:lime={!ended(e)}>{ended(e) ? 'Selesai' : e.category}</span>
      <h2><a href={'/event/' + e.slug}>{e.title}</a></h2>
      <a class="meta" href={'/ruang/' + e.organizer.slug}>{e.organizer.name}</a>
      <div class="event-info-line">
        <Icon name="calendar" size={16} /><span>{date(e.starts_at)} · {time(e.starts_at)}</span>
      </div>
      <div class="event-info-line"><Icon name="pin" size={16} /><span>{e.venue}</span></div>
      <div class="event-card-footer">
        <strong>{ended(e) ? 'Pertunjukan selesai' : money(e.price)}</strong>
        <div class="row">
          <button
            class="icon-button"
            class:active={saved}
            onclick={bookmark}
            aria-pressed={saved}
            aria-label={saved ? 'Batalkan simpan event' : 'Simpan event'}
            ><Icon name={saved ? 'check' : 'bookmark'} /></button
          >
          <a class="button small" href={'/event/' + e.slug + (ended(e) ? '?tab=ulasan' : '')}
            >{ended(e) ? 'Baca ulasan' : 'Lihat event'}</a
          >
        </div>
      </div>
    </div>
  </article>{/if}
