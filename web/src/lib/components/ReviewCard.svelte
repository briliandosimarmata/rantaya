<script lang="ts">
  import { getContext } from 'svelte';
  import type { Review, AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { date } from '#lib/format.js';
  import Icon from './Icon.svelte';
  import Avatar from './Avatar.svelte';
  let { review: r } = $props<{ review: Review }>();
  const ctx = getContext<AppContext>('app');
  let reply = $state(''),
    open = $state(false),
    busy = $state(false);
  let helped = $derived(ctx.account.helpful_reviews?.includes(r.id) || false);
  async function respond(event: SubmitEvent) {
    event.preventDefault();
    busy = true;
    try {
      await mutate('/reviews/' + r.id + '/reply', { body: reply });
      reply = '';
      open = false;
      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    } finally {
      busy = false;
    }
  }
  async function helpful() {
    if (!ctx.requireLogin()) return;
    try {
      await mutate('/reviews/' + r.id + '/helpful', { active: !helped });

      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    }
  }
</script>

<article class="card review-card">
  <header class="row"><Avatar name={r.author} src={r.avatar_url} /><div class="grow"><b class="small">{r.author}</b><div class="meta">{date(r.created_at)}</div></div>{#if r.verified}<span class="pill lime"><Icon name="check" size={13} />Kehadiran terverifikasi</span>{/if}<a class="icon-button" href={'/lapor?kind=review&id=' + r.id} aria-label="Laporkan ulasan"><Icon name="more" /></a></header>
  <div class="review-context">Setelah menonton <a href={'/event/' + r.event_slug + '?tab=ulasan'}>{r.event_title}</a><br /><span class="meta">{date(r.event_starts_at || r.created_at)} · {r.organizer_name || ''}</span></div>
  <p class="review-text">{r.body}</p>
  <div class="post-actions"><button class="action" class:active={helped} aria-pressed={helped} onclick={helpful}><Icon name="up" />Membantu · {r.helpful}</button>
    {#if ctx.user?.organizer_id === r.organizer_id}<button class="action" onclick={() => open = !open}><Icon name="chat" />Balas</button>{/if}
    {#if ctx.user?.id === r.account_id}<a class="action" href={'/event/' + r.event_slug + '/ulasan'}><Icon name="edit" />Edit</a>{/if}
  </div>
  {#each r.replies as rr}<div class="reply"><b>{rr.author}</b><p>{rr.body}</p></div>{/each}
  {#if open}<form class="reply-form" onsubmit={respond}><label for={'reply-' + r.id}>Balasan pengelola</label><textarea id={'reply-' + r.id} bind:value={reply} maxlength="2000" required></textarea><div class="row"><button class="button primary" disabled={busy}>Kirim balasan</button><button class="button" type="button" onclick={() => open = false}>Batal</button></div></form>{/if}
</article>
