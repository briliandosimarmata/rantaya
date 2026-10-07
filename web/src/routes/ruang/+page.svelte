<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import Follow from '#lib/components/Follow.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
</script>
<SEO title={'Komunitas seni ' + ctx.city} />
<PageHeading title="Temukan ruangmu" description="Ikuti komunitas dan kenali cerita di balik panggung." />
<div class="content-grid"><section class="main-column">
  <div class="filters"><div class="filter-chips">{#each [['', 'Semua'], ['Teater', 'Teater'], ['Musik', 'Musik']] as [value,label]}<a class:active={data.category === value} href={'?category=' + value}>{label}</a>{/each}</div></div>
  <div class="community-grid">{#each data.organizers as o}<article class="card community-card">
    <div class="between"><Avatar name={o.name} src={o.avatar_url} size="square lime" /><span class="pill">{o.category}</span></div>
    <a href={'/ruang/' + o.slug}><h3>{o.name}</h3></a><p>{o.description}</p>
    <div class="meta"><Icon name="pin" size={14} /> {o.city} · {o.followers} pengikut</div>
    <div class="between"><a href={'/ruang/' + o.slug} class="link">Masuk ke ruang</a><Follow id={o.id} /></div>
  </article>{/each}</div>
  {#if !data.organizers.length}<Empty title="Ruang di kota ini belum tersedia" body="Peluncuran awal berfokus di Karawang." href="/ruang?city=Karawang" label="Lihat Karawang" />{/if}
</section><DiscoveryAside events={data.nearby} organizers={data.allOrganizers} /></div>
