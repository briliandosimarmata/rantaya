<script lang="ts">
  import { onMount } from 'svelte';
  import { untrack } from 'svelte';
  import { goto } from '$app/navigation';
  import EventCard from '#lib/components/EventCard.svelte';
  import PostCard from '#lib/components/PostCard.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import { date } from '#lib/format.js';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  let input: HTMLInputElement;
  let query = $state(untrack(() => data.q));
  let timer: ReturnType<typeof setTimeout>;
  function search(e: SubmitEvent) {
    e.preventDefault();
    clearTimeout(timer);
    goto('/cari?q=' + encodeURIComponent(query), { reset: false });
  }
  function changed() {
    clearTimeout(timer);
    timer = setTimeout(() => goto('/cari?q=' + encodeURIComponent(query), { reset: false }), 350);
  }
  onMount(() => {
    input?.focus();
    return () => clearTimeout(timer);
  });
</script>

<SEO title="Cari event dan komunitas" privatePage />
<PageHeading title="Cari di Rantaya" />
<section class="main-column"><form class="page-search" onsubmit={search}><label class="sr-only" for="search-q">Cari komunitas, event, atau obrolan</label><Icon name="search" /><input bind:this={input} id="search-q" bind:value={query} oninput={changed} placeholder="Komunitas, event, atau obrolan..." autocomplete="off" maxlength="100" /><button class="button primary">Cari</button></form>
  {#if !data.q}<h3 style="margin:24px 0 12px">Ruang di Karawang</h3>{#each data.organizers as o}<a class="search-result" href={'/ruang/' + o.slug}><h3>{o.name}</h3><p class="meta">{o.category} · {o.city}</p></a>{/each}
  {:else}{#if data.organizers.length}<h3 style="margin-top:25px">Komunitas</h3>{#each data.organizers as o}<a class="search-result" href={'/ruang/' + o.slug}><h3>{o.name}</h3><p>{o.category} · {o.city}</p></a>{/each}{/if}
    {#if data.events.length}<h3 style="margin-top:25px">Event</h3>{#each data.events as e}<a class="search-result" href={'/event/' + e.slug}><h3>{e.title}</h3><p>{date(e.starts_at)} · {e.city}</p></a>{/each}{/if}
    {#if data.posts.length}<h3 style="margin-top:25px">Obrolan</h3>{#each data.posts as p}<a class="search-result" href={'/post/' + p.id}><h3>{p.title || p.body.slice(0,70)}</h3><p>{p.author} · {p.city}</p></a>{/each}{/if}
    {#if !data.organizers.length && !data.events.length && !data.posts.length}<Empty title="Belum ketemu" body="Coba judul event, nama komunitas, atau topik lain." />{/if}
  {/if}
</section>
