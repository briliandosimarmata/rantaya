<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import EventCard from '#lib/components/EventCard.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import { dateRangeLabel } from '#lib/format.js';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  function url(changes: any) {
    return (
      '/agenda?' +
      new URLSearchParams({
        period: data.period,
        category: data.category,
        from: data.from,
        to: data.to,
        ...changes,
      })
    );
  }
</script>

<SEO title={'Agenda pertunjukan ' + ctx.city} />
<PageHeading title="Agenda panggung" description="Pilih tanggal, temukan cerita berikutnya." />
<div class="content-grid">
  <section>
    <nav class="tabs" aria-label="Periode agenda">
      <a class:active={data.period === 'upcoming'} href={url({ period: 'upcoming' })}>Mendatang</a
      ><a class:active={data.period === 'past'} href={url({ period: 'past' })}>Selesai</a><a
        class:active={data.period === 'all'}
        href={url({ period: 'all' })}>Semua agenda</a
      >
    </nav>
    <div class="filters">
      <div class="filter-chips">
        {#each [['', 'Semua'], ['Teater', 'Teater'], ['Musik', 'Musik']] as [value, label]}<a
            class:active={data.category === value}
            href={url({ category: value })}>{label}</a
          >{/each}
      </div>
      <a
        class="filter-control agenda-date-sort"
        aria-label={'Filter tanggal: ' + dateRangeLabel(data.from, data.to, 'Semua tanggal')}
        href={'/tanggal?' +
          new URLSearchParams({
            from: data.from,
            to: data.to,
            category: data.category,
            period: data.period,
          })}
        ><Icon name="calendar" size={18} /><span>{dateRangeLabel(data.from, data.to)}</span><Icon
          name="down"
          size={14}
        /></a
      >
    </div>
    {#if data.events.length}<div class="events-list">
        {#each data.events as event}<EventCard {event} />{/each}
      </div>
      {#if data.events.length === 20}<a class="button" href={url({ page: data.page + 1 })}
          >Agenda berikutnya</a
        >{/if}{:else}<Empty
        title="Belum ada acara pada pilihan ini"
        body="Coba kategori lain atau lihat semua agenda."
      />{/if}
  </section>
  <DiscoveryAside events={data.nearby} organizers={data.organizers} />
</div>
