<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import Avatar from '#lib/components/Avatar.svelte';
  import Follow from '#lib/components/Follow.svelte';
  import PostCard from '#lib/components/PostCard.svelte';
  import EventCard from '#lib/components/EventCard.svelte';
  import ProductCard from '#lib/components/ProductCard.svelte';
  import ReviewCard from '#lib/components/ReviewCard.svelte';
  import ComposerEntry from '#lib/components/ComposerEntry.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let o = $derived(data.organizer);
  const tabs = [['linimasa','Linimasa'],['mention','Mention'],['agenda','Agenda'],['merch','Merchandise'],['ulasan','Ulasan'],['tentang','Tentang']];
</script>
<SEO title={o.name} description={o.description} image={o.cover_url} />
<PageHeading title="Ruang komunitas" secondary />
<div class="content-grid"><section class="main-column">
  {#if o.cover_url}<div class="cover"><img src={o.cover_url} alt="Foto komunitas contoh" width="1100" height="300" /><span class="cover-label">{o.category.toUpperCase()} / {o.city.toUpperCase()}</span></div>{/if}
  <header class="organizer-header">
    <Avatar name={o.name} src={o.avatar_url} size="large big lime" />
    <div class="community-identity"><h1>{o.name}</h1><span class="meta">{o.category} · {o.city}</span></div>
    <p class="description">{o.description}</p>
    <div class="stats"><span><b>{o.followers}</b> pengikut</span><span><b>{data.events.length}{data.events.length >= 20 ? '+' : ''}</b> event</span><span><b>{data.reviews.length}{data.reviews.length >= 20 ? '+' : ''}</b> ulasan</span></div>
    <div class="follow"><Follow id={o.id} /></div>
  </header>
  <nav class="tabs compact" aria-label="Isi komunitas">{#each tabs as [value,label]}<a class:active={data.tab === value} href={'?tab=' + value}>{label}{data.tab === 'mention' && value === 'mention' && data.posts.length ? ' (' + data.posts.length + ')' : ''}</a>{/each}</nav>
  {#if data.tab === 'linimasa' || data.tab === 'mention'}
    {#if data.tab === 'mention' && ctx.user?.role !== 'organizer' || data.tab === 'linimasa' && ctx.user?.organizer_id === o.id}<ComposerEntry href={'/tulis?organizer=' + o.id} />{/if}
    {#each data.posts as post}<PostCard {post} events={data.events} />{:else}<Empty title={data.tab === 'mention' ? 'Belum ada mention' : 'Belum ada obrolan di sini'} body={data.tab === 'mention' ? 'Mulai cerita dan mention komunitas atau eventnya.' : 'Ikuti komunitas atau mulai cerita pertamamu.'} />{/each}
  {:else if data.tab === 'agenda'}<div class="events-list">{#each data.events as event}<EventCard {event} />{/each}</div>
  {:else if data.tab === 'merch'}<div class="product-grid">{#each data.products as p}<ProductCard product={p} />{:else}<Empty title="Belum ada merchandise" body="Kabar produk baru akan dibagikan di linimasa." />{/each}</div>
  {:else if data.tab === 'ulasan'}<div class="note"><b>{data.reviews.length} ulasan dari {new Set(data.reviews.map((r: any) => r.event_id)).size} event.</b><br />Setiap ulasan tetap terhubung ke acara yang ditonton.</div>{#each data.reviews as review}<ReviewCard {review} />{:else}<Empty title="Belum ada ulasan" body="Ulasan dari acara yang selesai akan tampil di sini." />{/each}
  {:else}<div class="card pad body-copy"><h2>Tentang {o.name}</h2><p>{o.about || o.description}</p><div class="divider"></div><div class="info-line"><Icon name="pin" />{o.city}</div><p class="small">Aktif sejak {new Date(o.created_at).getFullYear()} · {o.category}</p><a class="button" href={'/tulis?organizer=' + o.id}>Hubungi pengelola</a></div>{/if}
</section><DiscoveryAside events={data.nearby} organizers={data.organizers} /></div>
