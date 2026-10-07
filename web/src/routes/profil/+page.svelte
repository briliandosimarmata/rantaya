<script lang="ts">
  import { getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import ProfileForm from '#lib/components/ProfileForm.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  async function logout() {
    try {
      await mutate('/auth/logout');
      await goto('/', { invalidateAll: true });
    } catch (e: any) {
      ctx.toast(e.message);
    }
  }
</script>

<SEO title="Profil saya" privatePage />
<PageHeading title={ctx.user?.role === 'organizer' ? 'Akun pengelola' : 'Profil kamu'} />
<div class="profile-layout wide-page">
  <aside class="card profile-card"><Avatar name={ctx.user?.name || '?'} src={ctx.user?.avatar_url || ''} size={'big large lime ' + (ctx.user?.role === 'organizer' ? 'square' : '')} /><h2>{ctx.user?.name}</h2><span class="pill account-label">{ctx.user?.role === 'organizer' ? 'Pengelola' : ctx.user?.role === 'admin' ? 'Moderator' : 'Pengguna'}</span>
    <span class="meta" style="display:block;margin-top:14px">{ctx.user?.email}</span><p>{ctx.user?.bio}</p><div class="divider"></div><div class="between small"><span>Mengikuti</span><b>{ctx.account.follows.length} ruang</b></div>
    {#if ctx.user?.role === 'organizer'}<a href="/kelola" class="button full" style="margin-top:17px">Kelola komunitas</a>{:else}<a class="button full" href="/transaksi" style="margin-top:17px"><Icon name="ticket" />Transaksi & tiket</a><a class="link" href="/tiket" style="display:block;margin-top:12px">Tiket saya</a>{/if}
    <button class="button full" style="margin-top:25px" onclick={logout}><Icon name="logout" />Keluar akun</button>
    <a href={'/masuk/' + (ctx.user?.role === 'organizer' ? 'customer' : 'organizer')} class="link" style="display:block;margin-top:16px">Buka login {ctx.user?.role === 'organizer' ? 'pengguna' : 'pengelola'}</a>
  </aside>
  <section><nav class="tabs" aria-label="Pengaturan akun"><a class:active={data.tab !== 'settings'} href="/profil">Profil</a><a class:active={data.tab === 'settings'} href="/profil?tab=settings">Notifikasi</a></nav>
    {#if data.tab === 'settings'}<div class="card pad"><h2>Kabar yang kamu terima</h2><ProfileForm notifications /></div>
    {:else if ctx.user?.role === 'organizer' && data.organizer}<div class="card pad"><h2>{data.organizer.name}</h2><p class="small muted" style="margin:15px 0">{data.organizer.description}</p><a class="button primary" href="/kelola/profil">Edit identitas komunitas</a></div>
    {:else}<div class="card pad"><h2>Edit profil</h2><ProfileForm /></div>{/if}
  </section>
</div>
