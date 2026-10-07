<script lang="ts">
  import { page } from '$app/state';
  import SEO from '#lib/components/SEO.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  let { children, data } = $props();
  const links = [
    ['/kelola', 'Ringkasan'], ['/kelola/event', 'Event'], ['/kelola/postingan', 'Postingan'],
    ['/kelola/merch', 'Merchandise'], ['/kelola/ulasan', 'Ulasan'], ['/kelola/profil', 'Pengaturan'],
    ['/kelola/pembayaran', 'Pembayaran'], ['/kelola/scanner', 'Scan tiket'], ['/kelola/rekening', 'Rekening'],
  ];
  let editor = $derived(/^\/kelola\/(event|merch)\/[^/]+$/.test(page.url.pathname));
  let originalPage = $derived(links.slice(0,6).some(([path]) => path === page.url.pathname));
</script>
<SEO title="Kelola komunitas" privatePage />
{#if !editor}
<PageHeading title="Kelola komunitas" description={data.organizer.name} secondary={!originalPage}>
  {#snippet action()}<a class="button" href={'/ruang/' + data.organizer.slug}>Lihat ruang</a>{/snippet}
</PageHeading>
<div class="wide-page dashboard-page">
  <nav class="tabs manage-tabs" aria-label="Dashboard pengelola">
    {#each links as [path,label]}<a class:active={path === '/kelola' ? page.url.pathname === path : page.url.pathname.startsWith(path)} href={path}>{label}</a>{/each}
  </nav>
  {@render children()}
</div>
{:else}{@render children()}{/if}
