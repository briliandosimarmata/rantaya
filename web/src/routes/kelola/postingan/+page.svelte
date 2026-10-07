<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import Icon from '#lib/components/Icon.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let confirming = $state('');
  async function pin(id: string, pinned: boolean) {
    try {
      await mutate('/posts/' + id + '/pin', { pinned }, 'PATCH');
      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    }
  }
  async function remove(id: string) {
    try {
      await mutate('/posts/' + id, {}, 'DELETE');
      confirming = '';
      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    }
  }
</script>

<div class="between" style="margin-bottom:20px"><h2>Postingan ruang</h2><a class="button primary" href="/tulis">Buat postingan</a></div>
<div class="community-grid">
{#each data.posts as post}<article class="card pad"><h3 style="margin:13px 0">{post.title || post.body.slice(0,70)}</h3><div class="meta">{post.votes} dukungan · {post.comment_count} komentar</div><div class="row wrap" style="margin-top:17px"><a class="button slim" href={'/post/' + post.id}>Buka</a><button class="button slim" onclick={() => pin(post.id,!post.pinned)}>{post.pinned ? 'Lepas sematan' : 'Sematkan'}</button><button class="icon-button" onclick={() => confirming = confirming === post.id ? '' : post.id} aria-label="Kelola postingan"><Icon name="more" /></button></div>
{#if confirming===post.id}<div class="row wrap" style="margin-top:12px"><span class="small">Hapus postingan ini?</span><button class="button" onclick={() => remove(post.id)}>Ya, hapus</button><button class="action" onclick={() => confirming = ''}>Batal</button></div>{/if}</article>{/each}
</div>
