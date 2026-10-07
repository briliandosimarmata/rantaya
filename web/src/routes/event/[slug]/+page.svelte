<script lang="ts">
  import { getContext } from 'svelte';
  import { page } from '$app/state';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { date, time, money, ended, share, video } from '#lib/format.js';
  import Icon from '#lib/components/Icon.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import PostCard from '#lib/components/PostCard.svelte';
  import ReviewCard from '#lib/components/ReviewCard.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import Follow from '#lib/components/Follow.svelte';
  import SEO from '#lib/components/SEO.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import ComposerEntry from '#lib/components/ComposerEntry.svelte';
  import EventActions from '#lib/components/EventActions.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let e = $derived(data.event);
  let playing = $state(false);
  let media = $derived(video(e.trailer_url));
  let saved = $derived(ctx.account.bookmarks.some((b) => b.kind === 'event' && b.id === e.id));
  let schema = $derived(
    JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'Event',
      name: e.title,
      description: e.description,
      startDate: e.starts_at,
      endDate: e.ends_at,
      eventAttendanceMode: 'https://schema.org/OfflineEventAttendanceMode',
      eventStatus: 'https://schema.org/EventScheduled',
      location: {
        '@type': 'Place',
        name: e.venue,
        address: {
          '@type': 'PostalAddress',
          streetAddress: e.address,
          addressLocality: e.city,
          addressCountry: 'ID',
        },
      },
      organizer: {
        '@type': 'Organization',
        name: e.organizer.name,
        url: page.url.origin + '/ruang/' + e.organizer.slug,
      },
      ...(e.flyer_url ? { image: new URL(e.flyer_url, page.url.origin).href } : {}),
      offers: e.sessions.map((s: any) => ({
        '@type': 'Offer',
        name: s.label,
        price: s.price,
        priceCurrency: 'IDR',
        availability: s.available > 0 ? 'https://schema.org/InStock' : 'https://schema.org/SoldOut',
        url: page.url.origin + '/event/' + e.slug + '/tiket',
      })),
    }).replace(/</g, '\\u003c'),
  );
  async function bookmark() {
    if (!ctx.requireLogin()) return;
    try {
      await mutate('/bookmarks', { kind: 'event', id: e.id, active: !saved });
      await ctx.refresh();
    } catch (err: any) {
      ctx.toast(err.message);
    }
  }
  async function sharing() {
    try {
      await share(e.title, '/event/' + e.slug);
      ctx.toast('Tautan siap dibagikan.');
    } catch (err: any) {
      if (err.name !== 'AbortError') ctx.toast('Tautan belum bisa disalin.');
    }
  }
</script>

<SEO
  privatePage={!e.published}
  title={e.title + ' — ' + e.city}
  description={e.description}
  image={e.flyer_url}
/><svelte:head>{@html '<script type="application/ld+json">' + schema + '</script>'}</svelte:head>
<article class="event-page">
  <PageHeading title="Detail event" secondary />
  <div class="crumb"><a href="/agenda">Agenda</a><Icon name="next" size={14} /><span>{e.title}</span></div>
  <div class="content-grid"><section class="main-column">
    <header class="event-summary"><span class="pill" class:lime={!ended(e)}>{ended(e) ? 'Selesai' : e.category}</span><h1>{e.title}</h1>
      <a class="row event-organizer" href={'/ruang/' + e.organizer.slug}><Avatar name={e.organizer.name} src={e.organizer.avatar_url} size="square" /><b class="small">{e.organizer.name}</b></a>
      <div class="info-line"><Icon name="calendar" /><span>{date(e.starts_at)} · {time(e.starts_at)}</span></div>
      <div class="info-line"><Icon name="pin" /><span>{e.venue}, {e.city}</span></div>
      <a class="button event-map-link" href={e.maps_url || 'https://www.google.com/maps/search/?api=1&query=' + encodeURIComponent(e.venue + ', ' + e.city)} target="_blank" rel="noopener noreferrer"><Icon name="link" />{e.maps_url ? 'Buka Google Maps' : 'Cari venue di Google Maps'}</a>
    </header>
    <nav class="tabs" aria-label="Detail event">{#each [['tentang','Tentang'],['obrolan','Obrolan'],['ulasan','Ulasan (' + data.reviews.length + ')']] as [tab,label]}<a class:active={data.tab === tab} href={'?tab=' + tab}>{label}</a>{/each}</nav>
    {#if data.tab === 'tentang'}
      {#if media}<section class="event-media" aria-label={'Trailer ' + e.title}>
        {#if media.type === 'youtube'}{#if playing}<iframe title={'Trailer ' + e.title} src={media.url} allow="autoplay; encrypted-media; picture-in-picture" allowfullscreen></iframe>{:else}<button class="button trailer-play" onclick={() => playing = true}>Putar trailer</button>{#if e.flyer_url}<img src={e.flyer_url} alt={'Flyer ' + e.title} />{/if}{/if}
        {:else}<video src={media.url} controls playsinline preload="metadata" poster={e.flyer_url}><track kind="captions" /></video>{/if}
      </section>{:else if e.flyer_url}<section class="event-media" aria-label={'Flyer ' + e.title}><img src={e.flyer_url} alt={'Flyer ' + e.title} width="960" height="540" fetchpriority="high" /></section>{/if}
      <div class="body-copy"><h2>Tentang pertunjukan</h2><p class="event-about">{e.description}</p>
        <div class="row wrap"><span class="pill">{e.duration}</span><span class="pill">{e.language}</span><span class="pill">{e.age}</span></div>
        {#if e.lineup.length}<div class="divider"></div><h3>Pengisi acara & tim</h3>
          <!-- svelte-ignore a11y_no_noninteractive_tabindex, a11y_no_noninteractive_element_interactions (Keyboard-accessible horizontal scroll region.) -->
          <div class="lineup" role="region" aria-label="Pengisi acara, geser ke samping" tabindex="0" onkeydown={(ev) => { if (ev.key === 'ArrowLeft' || ev.key === 'ArrowRight') { ev.preventDefault(); ev.currentTarget.scrollBy({ left: ev.key === 'ArrowRight' ? 130 : -130 }); } }}>
            {#each e.lineup as p}<div class="performer">{#if p.photo}<img class="lineup-photo" src={p.photo} alt={'Foto ' + p.name} width="88" height="88" loading="lazy" />{:else}<span class="lineup-photo lineup-initials">{p.name.split(/\s+/).slice(0,2).map((s: string) => s[0]).join('')}</span>{/if}<b>{p.name}</b><span class="meta">{p.role}</span></div>{/each}
          </div>{/if}
        {#if e.layout_url}<a class="button venue-button" href={'/event/' + e.slug + '/venue'}><Icon name="pin" />Lihat layout venue & tempat duduk<Icon name="next" /></a>{/if}
        <div class="event-sessions"><div class="divider"></div><h3>{e.sessions.length > 1 ? 'Pertunjukan & sesi' : 'Jadwal'}</h3>
          {#if ended(e)}<p>Pertunjukan telah selesai. Cerita dan ulasannya tetap terbuka.</p>{:else}{#each e.sessions as s}<div class="session"><div><b class="small">{s.label}</b><div class="meta">{time(s.starts_at)} · {money(s.price)}</div></div><a class="button slim" href={'/event/' + e.slug + '/tiket?session=' + s.id}>Pilih sesi</a></div>{/each}{/if}
          <div class="divider"></div><h3>Kenali penyelenggaranya</h3>
          <a class="event-inline event-organizer-link" href={'/ruang/' + e.organizer.slug}><Avatar name={e.organizer.name} src={e.organizer.avatar_url} size="square" /><div><b class="small">{e.organizer.name}</b><div class="meta">{e.organizer.category} · {e.organizer.city}</div></div><Icon name="next" /></a>
        </div>
      </div>
    {:else if data.tab === 'obrolan'}<ComposerEntry href={'/tulis?event=' + e.id} />{#each data.posts as post}<PostCard {post} events={[e]} />{:else}<Empty title="Obrolan acara ini masih terbuka" body="Mulai pertanyaan atau bagikan cerita tentang pertunjukan." />{/each}
    {:else}<div class="between" style="margin-bottom:20px"><h2>{data.reviews.length} ulasan penonton</h2>{#if ended(e) && ctx.user?.role !== 'organizer'}<a class="button primary" href={'/event/' + e.slug + '/ulasan'}>{data.reviews.some((r: any) => r.account_id === ctx.user?.id) ? 'Edit ulasan' : 'Tulis ulasan'}</a>{/if}</div>
      {#if !ended(e)}<div class="note">Ulasan dibuka setelah acara selesai. Untuk pertanyaan sebelum acara, gunakan tab Obrolan.</div>{/if}
      {#each data.reviews as review}<ReviewCard {review} />{:else}<Empty title={ended(e) ? 'Belum ada ulasan' : 'Pertunjukannya belum berlangsung'} body={ended(e) ? 'Bagikan pengalamanmu setelah menonton.' : 'Simpan acara dan ikuti penyelenggaranya.'} />{/each}
    {/if}
  </section><aside class="right-column"><div class="card ticket-panel"><span class="eyebrow">{ended(e) ? 'CERITA MASIH BERLANJUT' : 'TIKET PERTUNJUKAN'}</span><h2>{ended(e) ? 'Setelah pentas' : money(e.price)}</h2><p class="small muted">{ended(e) ? 'Ikut percakapan dan baca cerita penonton.' : 'Pembelian melalui kanal resmi penyelenggara.'}</p>
    <EventActions event={e} {saved} {bookmark} {sharing} canReview={ctx.user?.role !== 'organizer' && ctx.user?.role !== 'admin'} /><div class="divider"></div><div class="community-mini"><Avatar name={e.organizer.name} src={e.organizer.avatar_url} size="square" /><div class="community-text"><b>{e.organizer.name}</b><div class="meta">{e.organizer.followers} pengikut</div></div></div><Follow id={e.organizer.id} />
  </div></aside></div>
  <div class="mobile-ticket"><div class="mobile-event-actions">{#if !ended(e)}<span class="meta ticket-price">Mulai <b>{money(e.price)}</b></span>{/if}<EventActions event={e} {saved} {bookmark} {sharing} canReview={ctx.user?.role !== 'organizer' && ctx.user?.role !== 'admin'} /></div></div>
</article>
