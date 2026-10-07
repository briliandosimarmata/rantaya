<script lang="ts">
  import { getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import PostCard from '#lib/components/PostCard.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  function url(changes: any) {
    const q = new URLSearchParams({
      tab: data.tab,
      sort: data.sort,
      category: data.category,
      ...changes,
    });
    return '/?' + q.toString();
  }
</script>

<SEO title="Obrolan, proses, dan panggung terdekat" />
<div class="page-heading home-heading">
  <div class="meta"><Icon name="pin" size={14} /><span>{ctx.city} · Ruang seni lokal</span></div>
  <div class="between">
    <div>
      <h1>Ruang obrolan</h1>
      <p>Cerita dan percakapan dari komunitasmu.</p>
    </div>
  </div>
</div>
<div class="content-grid home-grid">
  <section class="home-feed" aria-label="Linimasa">
    <nav class="tabs" aria-label="Pilihan feed">
      <a
        class:active={data.tab === 'explore'}
        aria-current={data.tab === 'explore' ? 'page' : undefined}
        href={url({ tab: 'explore' })}>Jelajah</a
      ><a
        class:active={data.tab === 'following'}
        aria-current={data.tab === 'following' ? 'page' : undefined}
        href={url({ tab: 'following' })}>Diikuti</a
      >
    </nav>
    <div class="filters">
      <div class="filter-chips">
        {#each [['', 'Semua'], ['Teater', 'Teater'], ['Musik', 'Musik']] as [value, label]}<a
            class:active={data.category === value}
            aria-current={data.category === value ? 'true' : undefined}
            href={url({ category: value })}>{label}</a
          >{/each}
      </div>
      <label class="filter-control"
        ><span class="sr-only">Urutan postingan</span><select
          value={data.sort}
          onchange={(e) => goto(url({ sort: e.currentTarget.value }))}
          ><option value="latest">Terbaru</option><option value="popular">Paling didukung</option
          ></select
        ></label
      >
    </div>
    <div class="card compose-entry">
      <Avatar name={ctx.user?.name || '?'} src={ctx.user?.avatar_url || ''} size="lime" /><a
        class="composer-start"
        href="/tulis">Mau cerita apa hari ini?</a
      ><a class="icon-button" href="/tulis" aria-label="Buat postingan"><Icon name="edit" /></a>
    </div>
    {#if data.posts.length}{#each data.posts as post (post.id)}<PostCard
          {post}
          events={data.relatedEvents}
        />{/each}
      <div class="between">
        {#if data.page > 1}<a class="button" href={url({ page: data.page - 1 })}>Sebelumnya</a
          >{/if}{#if data.posts.length === 20}<a class="button" href={url({ page: data.page + 1 })}
            >Percakapan berikutnya</a
          >{/if}
      </div>{:else}<Empty
        title="Belum ada obrolan di sini"
        body="Ikuti komunitas atau mulai cerita pertamamu."
        href="/ruang"
        label="Jelajah komunitas"
      />{/if}
  </section>
  <DiscoveryAside events={data.events} organizers={data.organizers} />
</div>
