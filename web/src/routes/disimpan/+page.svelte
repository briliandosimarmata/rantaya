<script lang="ts">
  import PostCard from '#lib/components/PostCard.svelte'; import EventCard from '#lib/components/EventCard.svelte'; import Empty from '#lib/components/Empty.svelte'; import SEO from '#lib/components/SEO.svelte'; import PageHeading from '#lib/components/PageHeading.svelte'; import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  let { data } = $props();
  let events = $derived(data.items.filter((x: any) => x.item && x.kind === 'event'));
  let posts = $derived(data.items.filter((x: any) => x.item && x.kind === 'post'));
</script>
<SEO title="Disimpan" privatePage /><PageHeading title="Disimpan" description="Acara dan obrolan yang ingin kamu buka lagi." />
<div class="content-grid"><section class="main-column">
  {#if events.length}<h2 style="margin-bottom:17px">Event ({events.length})</h2><div class="events-list" style="margin-bottom:27px">{#each events as item}<EventCard event={item.item} />{/each}</div>{/if}
  {#if posts.length}<h2 style="margin-bottom:17px">Postingan ({posts.length})</h2>{#each posts as item}<PostCard post={item.item} events={data.relatedEvents} />{/each}{/if}
  {#if !events.length && !posts.length}<Empty title="Belum ada yang disimpan" body="Ketuk ikon simpan pada event atau postingan." href="/agenda" label="Jelajah agenda" />{/if}
</section><DiscoveryAside events={data.nearby} organizers={data.organizers} /></div>
