<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import Empty from '#lib/components/Empty.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let busy = $state(false);
  async function resolve(id: string, status: string) {
    busy = true;
    try {
      await mutate('/admin/reports/' + id, { status }, 'PATCH');
      await ctx.refresh();
      ctx.toast('Laporan diperbarui.');
    } catch (e: any) {
      ctx.toast(e.message);
    } finally {
      busy = false;
    }
  }
</script>

<SEO title="Moderasi platform" privatePage />
<div class="page-heading">
  <h1>Jaga percakapan tetap terbuka.</h1>
  <p>Tinjau laporan konten dengan konteks sebelum mengambil tindakan.</p>
</div>
<section class="narrow stack">
  {#each data.reports as r}<article class="card pad">
      <div class="between">
        <strong>{r.kind === 'post' ? 'Postingan' : 'Ulasan'}</strong><span class="pill"
          >{r.status}</span
        >
      </div>
      <p>{r.reason}</p>
      <blockquote class="post-body">{r.content?.body}</blockquote>
      <div class="row wrap">
        <button class="button" disabled={busy} onclick={() => resolve(r.id, 'reviewed')}
          >Tandai ditinjau</button
        ><button class="button" disabled={busy} onclick={() => resolve(r.id, 'hidden')}
          >Sembunyikan konten</button
        >
      </div>
    </article>{:else}<Empty
      title="Belum ada laporan"
      body="Laporan pengguna akan muncul di sini."
    />{/each}
</section>
