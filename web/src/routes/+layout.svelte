<script lang="ts">
  import '../app.css';
  import '../prototype.css';
  import { onMount, setContext, untrack } from 'svelte';
  import { page } from '$app/state';
  import { goto, invalidateAll } from '$app/navigation';
  import { request } from '#lib/api.js';
  import type { AppContext, SessionState } from '#lib/types.js';
  import Icon from '#lib/components/Icon.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  let { data, children } = $props();
  let account = $state<SessionState>(untrack(() => data.account));
  let cityOpen = $state(false),
    message = $state('');
  let timer: ReturnType<typeof setTimeout>;
  $effect(() => {
    account = data.account;
  });
  function toast(m: string) {
    message = m;
    clearTimeout(timer);
    timer = setTimeout(() => (message = ''), 4200);
  }
  const ctx: AppContext = {
    get account() {
      return account;
    },
    get user() {
      return account.user;
    },
    get city() {
      return data.city;
    },
    get demo() {
      return data.config.demo;
    },
    toast,
    requireLogin() {
      if (account.user) return true;
      goto('/masuk/customer?next=' + encodeURIComponent(page.url.pathname + page.url.search));
      return false;
    },
    async refresh() {
      await invalidateAll();
    },
  };
  setContext('app', ctx);
  const links = [
    ['/', 'Beranda', 'home'],
    ['/agenda', 'Agenda', 'calendar'],
    ['/ruang', 'Ruang komunitas', 'people'],
    ['/disimpan', 'Disimpan', 'bookmark'],
  ];
  const taskPage = $derived(
    /^\/(tulis|tanggal|masuk|onboarding)(\/|$)/.test(page.url.pathname) ||
      /^\/event\/[^/]+\/(tiket|ulasan|venue)(\/|$)/.test(page.url.pathname) ||
      /^\/transaksi\/[^/]+/.test(page.url.pathname) ||
      /^\/kelola\/(event|merch)\/[^/]+/.test(page.url.pathname),
  );
  const rootPage = $derived(['/', '/agenda', '/ruang', '/disimpan'].includes(page.url.pathname));
  function active(path: string) {
    return path === '/' ? page.url.pathname === '/' : page.url.pathname.startsWith(path);
  }
  function back() {
    if (history.length > 1) history.back();
    else goto('/');
  }
  async function changeCity(city: string) {
    document.cookie =
      'ruang_city=' + encodeURIComponent(city) + '; Path=/; SameSite=Lax; Max-Age=31536000';
    cityOpen = false;
    await invalidateAll();
  }
  onMount(() => {
    const poll = setInterval(async () => {
      if (!account.user) return;
      try {
        const fresh = await request<SessionState>('/auth/me');
        account = { ...account, unread: fresh.unread };
      } catch {}
    }, 30000);
    return () => {
      clearInterval(poll);
      clearTimeout(timer);
    };
  });
</script>

<svelte:window
  onclick={(e) => {
    if (!(e.target instanceof Element) || !e.target.closest('.location')) cityOpen = false;
  }}
  onkeydown={(e) => {
    if (e.key === 'Escape') cityOpen = false;
  }}
/>

<svelte:head
  ><link rel="icon" href="/favicon.svg" /><meta name="theme-color" content="#172322" /></svelte:head
>
<a class="skip" href="#content">Lewati navigasi</a>
<div
  class="app-shell"
  class:task-shell={taskPage}
  class:detail-event={/^\/event\/[^/]+$/.test(page.url.pathname)}
>
  <aside class="sidebar">
    <a class="brand" href="/" aria-label="Rantaya, beranda"><span class="brand-mark">r.</span>RANTAYA</a>
    <nav aria-label="Navigasi utama" style="margin-top:30px">
      {#each links as [path, label, icon]}<a class="navlink" class:active={active(path)} href={path} aria-current={active(path) ? 'page' : undefined}><Icon name={icon} />{label}</a>{/each}
    </nav>
    <div class="side-label">RUANG YANG DIIKUTI</div>
    {#each data.followedOrganizers as o}<a class="space-link" href={'/ruang/' + o.slug}><Avatar name={o.name} src={o.avatar_url} size="square" />{o.name}</a>{/each}
    <div class="side-bottom">
      <a class="manage-link" href={account.user?.role === 'organizer' ? '/kelola' : '/masuk/organizer?next=/kelola'}><Icon name="settings" />{account.user?.role === 'organizer' ? 'Kelola komunitas' : 'Masuk pengelola'}</a>
      {#if account.user}<a class="manage-link transaction-access" href="/transaksi"><Icon name="ticket" />Transaksi & tiket</a>{/if}
      <a href={account.user ? '/profil' : '/masuk/customer'} class="side-user"><Avatar name={account.user?.name || '?'} size={account.user?.role === 'organizer' ? 'square' : ''} /><span><b class="small">{account.user?.name || 'Pengunjung'}</b><span class="meta" style="display:block">{account.user?.role === 'organizer' ? 'Akun pengelola' : 'Akun pengguna'}</span></span></a>
      {#if data.config.demo}<span class="prototype">Demo · seluruh data contoh</span>{/if}
    </div>
  </aside>
  <div class="workspace">
<header class="topbar">
  <div class="topbar-inner">
    <div class="header-start">
      {#if rootPage}<a class="brand mobile-brand" href="/" aria-label="Rantaya, beranda"
          ><span class="brand-mark">r.</span><span>rantaya<span class="brand-dot">.</span></span></a
        >{:else}<button
          class="icon-button header-back"
          onclick={back}
          aria-label="Kembali ke halaman sebelumnya"><Icon name="back" /></button
        >{/if}
      <div class="location">
        <button
          class="location-trigger"
          onclick={() => (cityOpen = !cityOpen)}
          aria-expanded={cityOpen}
          aria-controls="city-menu"
          ><Icon name="pin" /><span>{data.city}</span><Icon name="down" /></button
        >{#if cityOpen}<div class="location-menu" id="city-menu">
            <p class="meta">Jelajahi kota</p>
            {#each data.config.cities as city}<button
                class:active={city === data.city}
                onclick={() => changeCity(city)}
                ><span>{city}</span>{#if city === data.city}<Icon
                    name="check"
                  />{:else if city === 'Jakarta'}<span class="meta">Segera hadir</span>{/if}</button
              >{/each}
          </div>{/if}
      </div>
    </div>
    {#if !taskPage}<a class="search-trigger" href="/cari" aria-label="Cari komunitas atau event"
        ><Icon name="search" /><span>Cari komunitas atau event...</span></a
      >{/if}
    <div class="top-actions">
      {#if !taskPage}{#if data.config.demo}{/if}<a
          class="icon-button notification-icon"
          href="/notifikasi"
          aria-label={'Notifikasi, ' + account.unread + ' belum dibaca'}
          ><Icon name="bell" />{#if account.unread}<span class="badge">{account.unread}</span
            >{/if}</a
        >{/if}<a
        class="profile-access"
        href={account.user ? '/profil' : '/masuk/customer'}
        aria-label={account.user ? 'Buka profil' : 'Masuk akun'}
        ><Avatar
          name={account.user?.name || '?'}
          src={account.user?.avatar_url || ''}
          size={'lime ' + (account.user?.role === 'organizer' ? 'square' : '')}
        /></a
      >
      {#if !taskPage}<a class="button slim desktop-compose" href="/tulis"><Icon name="plus" /><span>Buat post</span></a>{/if}
    </div>
  </div>
</header>
  <main id="content" tabindex="-1">
    {#key page.url.pathname}{@render children()}{/key}
  </main>
  </div>
</div>
{#if !taskPage}<nav class="bottom-nav" aria-label="Navigasi mobile">
    {#each links as [path, label, icon]}<a
        class:active={active(path)}
        href={path}
        aria-current={active(path) ? 'page' : undefined}
        ><Icon name={icon} /><span>{path === '/ruang' ? 'Ruang' : label}</span></a
      >{/each}
  </nav>
  {#if !page.url.pathname.startsWith('/event/') && !page.url.pathname.startsWith('/kelola') && ctx.user?.role !== 'admin'}<a
      class="mobile-compose-fab"
      href="/tulis"
      aria-label="Buat postingan"><Icon name="edit" size={24} /></a
    >{/if}{/if}
<div class:visible={!!message} class="toast" role="status" aria-live="polite">{message}</div>
