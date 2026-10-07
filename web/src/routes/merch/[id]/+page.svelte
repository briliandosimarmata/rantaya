<script lang="ts">
  import { untrack } from 'svelte';
  import { recordClick } from '#lib/api.js';
  import { money } from '#lib/format.js';
  import SEO from '#lib/components/SEO.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import Icon from '#lib/components/Icon.svelte';
  let { data } = $props();
  let p = $derived(data.product);
  let variant = $state(untrack(() => data.product.variants[0] || ''));
  let link = $derived.by(() => {
    if (!p.purchase_url) return '';
    const u = new URL(p.purchase_url);
    if (u.hostname === 'wa.me') u.searchParams.set('text','Halo, saya ingin memesan ' + p.name + ' (' + (variant || p.variants[0]) + ').');
    return u.href;
  });
</script>
<SEO title={p.name} description={p.description} image={p.image_url} />
<PageHeading title="Merchandise" secondary />
<div class="crumb"><a href={'/ruang/' + p.organizer.slug + '?tab=merch'}>{p.organizer.name}</a><Icon name="next" /><span>Koleksi</span></div>
<article class="product-detail wide-page">
  {#if p.image_url}<img class="product-image" src={p.image_url} alt={'Foto contoh ' + p.name} width="720" height="720" />{/if}
  <section><span class="pill lime">{p.availability}</span><h1>{p.name}</h1>
    <a class="row" href={'/ruang/' + p.organizer.slug}><Avatar name={p.organizer.name} src={p.organizer.avatar_url} size="square" /><b class="small">{p.organizer.name}</b></a>
    <div class="price">{money(p.price)}</div><p class="small muted">{p.description}</p>
    <div class="field"><label for="variant">{p.variants.length > 1 ? 'Pilih ukuran' : 'Ukuran'}</label><select id="variant" bind:value={variant}>{#each p.variants as v}<option value={v}>{v}</option>{/each}</select></div>
    {#if p.availability === 'Habis'}<button class="button primary full" disabled><Icon name="bag" />Pesan ke penyelenggara</button>{:else}<a class="button primary full" href={link || '/ruang/' + p.organizer.slug} onclick={() => recordClick('merch',p.id)} target={link ? '_blank' : undefined} rel="noopener noreferrer"><Icon name="bag" />Pesan ke penyelenggara</a>{/if}
    <div class="note" style="margin-top:17px">Pemesanan melalui kanal penyelenggara. Produk tetap tersedia setelah acara selesai.</div>
    {#if p.event_title}<a class="event-inline" href={'/event/' + p.event_slug}><Icon name="calendar" /><span class="small">Dari pertunjukan<br /><b>{p.event_title}</b></span><Icon name="next" /></a>{/if}
  </section>
</article>
