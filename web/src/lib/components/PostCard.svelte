<script lang="ts">
  import { getContext, onDestroy } from 'svelte';
  import type { AppContext, Post, ArtsEvent } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { date, ended, money, share } from '#lib/format.js';
  import Avatar from './Avatar.svelte';
  import Icon from './Icon.svelte';
  let {
    post: p,
    detail = false,
    events = [],
  } = $props<{ post: Post; detail?: boolean; events?: ArtsEvent[] }>();
  const ctx = getContext<AppContext>('app');
  let busy = $state(false),
    pendingVote = $state<boolean | null>(null),
    votePulse = $state(false);
  let confirmedVote = $derived(ctx.account.votes.includes(p.id));
  let voted = $derived(pendingVote ?? confirmedVote);
  let voteCount = $derived(
    p.votes + (pendingVote === null ? 0 : Number(pendingVote) - Number(confirmedVote)),
  );
  let saved = $derived(ctx.account.bookmarks.some((b) => b.kind === 'post' && b.id === p.id));
  let menuOpen = $state(false);
  let pulseTimer: ReturnType<typeof setTimeout>;
  onDestroy(() => clearTimeout(pulseTimer));
  const relatedEvent = $derived(
    events.find((event: ArtsEvent) =>
      p.mentions.some(
        (mention: Post['mentions'][number]) => mention.kind === 'event' && mention.id === event.id,
      ),
    ),
  );
  async function act(kind: 'vote' | 'save') {
    if (busy || !ctx.requireLogin()) return;
    busy = true;
    if (kind === 'vote') {
      pendingVote = !confirmedVote;
      votePulse = true;
      clearTimeout(pulseTimer);
      pulseTimer = setTimeout(() => (votePulse = false), 300);
    }
    try {
      if (kind === 'vote') await mutate('/posts/' + p.id + '/vote', { active: pendingVote });
      else await mutate('/bookmarks', { kind: 'post', id: p.id, active: !saved });
      await ctx.refresh();
    } catch (e: any) {
      votePulse = false;
      ctx.toast(e.message);
    } finally {
      pendingVote = null;
      busy = false;
    }
  }
  async function sharing() {
    menuOpen = false;
    try {
      await share(p.title || p.author, '/post/' + p.id);
      ctx.toast('Tautan siap dibagikan.');
    } catch (e: any) {
      if (e.name !== 'AbortError') ctx.toast('Tautan belum bisa disalin.');
    }
  }
</script>

<svelte:window
  onkeydown={(e) => {
    if (e.key === 'Escape') menuOpen = false;
  }}
  onclick={(e) => {
    if (!(e.target instanceof Element) || !e.target.closest('[data-post-menu="' + p.id + '"]'))
      menuOpen = false;
  }}
/>
<article class="post-card card" data-post-id={p.id} aria-label={'Postingan ' + p.author}>
  <header class="post-head">
    <Avatar name={p.author} src={p.avatar_url} size={p.official ? 'square' : ''} />
    <div class="post-author">
      {#if p.organizer_id}<a class="author" href={'/ruang/' + p.organizer_id}>{p.author}</a
        >{:else if p.account_id === ctx.user?.id}<a class="author" href="/profil">{p.author}</a
        >{:else}<strong class="small">{p.author}</strong>{/if}{#if p.official}<span
          class="meta official-label"
        >
          · Pengelola</span
        >{/if}
      <div class="meta">
        {p.city} · {date(p.created_at, { day: 'numeric', month: 'short' })} · {date(p.created_at, {
          hour: '2-digit',
          minute: '2-digit',
        })}
      </div>
    </div>
    <div class="post-menu-wrap" data-post-menu={p.id}>
      <button
        class="icon-button post-menu-trigger"
        onclick={() => (menuOpen = !menuOpen)}
        aria-label="Pilihan postingan"
        aria-expanded={menuOpen}
        aria-controls={'post-menu-' + p.id}><Icon name="more" /></button
      >
      {#if menuOpen}<div class="post-menu" id={'post-menu-' + p.id}>
          <a href={'/lapor?kind=post&id=' + p.id} onclick={() => (menuOpen = false)}
            ><Icon name="flag" />Laporkan konten</a
          >
        </div>{/if}
    </div>
  </header>
  <div class="post-content">
    {#if p.pinned}<span class="pill highlight">Disematkan</span>{/if}
    {#if detail}{#if p.title}<h2>{p.title}</h2>{/if}
      <p class="post-body">{p.body}</p>{:else}<a class="post-open" href={'/post/' + p.id}
        >{#if p.title}<h2>{p.title}</h2>{/if}
        <p class="post-body">{p.body}</p></a
      >{/if}
    {#if p.mentions.length}<div class="mention-chips">
        {#each p.mentions as m}<a
            class="mention-chip"
            href={(m.kind === 'event' ? '/event/' : '/ruang/') + m.slug}>@{m.name}</a
          >{/each}
      </div>{/if}
    {#if relatedEvent}<a class="event-inline" href={'/event/' + relatedEvent.slug}>
        <span class="date-tile"
          ><span>{date(relatedEvent.starts_at, { month: 'short' }).toUpperCase()}</span><b
            >{date(relatedEvent.starts_at, { day: 'numeric' })}</b
          ></span
        >
        <span class="grow"
          ><b class="small">{relatedEvent.title}</b><span class="meta event-inline-meta"
            >{date(relatedEvent.starts_at)} · {ended(relatedEvent)
              ? 'Selesai'
              : money(relatedEvent.price)}</span
          ></span
        ><Icon name="next" />
      </a>{/if}
    {#if p.link_url}<a
        class="button purchase-link full"
        href={p.link_url}
        target="_blank"
        rel="noopener noreferrer"><Icon name="link" />{p.link_label || 'Lihat tautan pembelian'}</a
      >{/if}
  </div>
  {#if p.image_url}<a class="post-photo-link" href={'/post/' + p.id}
      ><img
        class="post-image"
        src={p.image_url}
        alt={'Foto pada postingan ' + p.author}
        loading="lazy"
        width="720"
        height="480"
      /></a
    >{/if}
  <footer class="post-actions">
    <button
      class:active={voted}
      class:vote-pulse={votePulse}
      class="action vote-action"
      onclick={() => act('vote')}
      onanimationend={() => (votePulse = false)}
      disabled={busy}
      aria-busy={busy}
      aria-pressed={voted}
      aria-label="Dukung postingan"><Icon name="up" /><span>{voteCount}</span></button
    >
    <a class="action" href={'/post/' + p.id + '#komentar'} aria-label="Buka komentar"
      ><Icon name="chat" />{p.comment_count} komentar</a
    >
    <span class="spacer"></span>
    <button
      class:active={saved}
      class="action"
      onclick={() => act('save')}
      disabled={busy}
      aria-pressed={saved}
      aria-label="Simpan postingan"><Icon name="bookmark" /></button
    >
    <button class="action" onclick={sharing} aria-label="Bagikan postingan"
      ><Icon name="share" /></button
    >
    {#if detail}<a class="action small" href={'/lapor?kind=post&id=' + p.id}>Laporkan</a>{/if}
  </footer>
</article>
