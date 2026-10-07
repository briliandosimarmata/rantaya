<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { date } from '#lib/format.js';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  async function read(id = '') {
    try {
      await mutate('/notifications/read', { id });
      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    }
  }
</script>

<SEO title="Notifikasi" privatePage />
<PageHeading title="Notifikasi" description="Kabar yang ingin kamu ikuti.">{#snippet action()}<button class="button" onclick={() => read()}>Tandai dibaca</button>{/snippet}</PageHeading>
<div class="content-grid"><section class="main-column"><div class="card">
  {#each data.items as n}<a class="notification" class:unread={!n.read_at} href={n.url} onclick={() => read(n.id)}><Avatar name={n.title} size="square" /><div><h3>{n.title}</h3><p>{n.body}</p><span class="meta">{date(n.created_at, { day:'numeric',month:'short',hour:'2-digit',minute:'2-digit' })}</span></div></a>{:else}<Empty title="Belum ada kabar baru" body="Kabar dari ruang yang kamu ikuti akan muncul di sini." />{/each}
</div><a class="button" href="/profil?tab=settings" style="margin-top:20px"><Icon name="settings" />Atur notifikasi</a></section><DiscoveryAside events={data.nearby} organizers={data.organizers} /></div>
