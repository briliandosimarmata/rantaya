<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { date } from '#lib/format.js';
  import PostCard from '#lib/components/PostCard.svelte';
  import Avatar from '#lib/components/Avatar.svelte';
  import PageHeading from '#lib/components/PageHeading.svelte';
  import DiscoveryAside from '#lib/components/DiscoveryAside.svelte';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let body = $state(''),
    reply = $state(''),
    replyID = $state(''),
    busy = $state(false);
  async function comment(event: SubmitEvent, parent = '') {
    event.preventDefault();
    if (!ctx.requireLogin()) return;
    busy = true;
    try {
      await mutate('/posts/' + data.post.id + '/comments', {
        body: parent ? reply : body,
        parent_id: parent,
      });
      body = '';
      reply = '';
      replyID = '';
      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    } finally {
      busy = false;
    }
  }
</script>

<SEO
  title={data.post.title || 'Percakapan ' + data.post.author}
  description={data.post.body}
  image={data.post.image_url}
/>
<PageHeading title="Percakapan" />
<div class="content-grid"><section class="main-column">
  <PostCard post={data.post} events={data.relatedEvents} detail />
  <section id="komentar">
    <h2 style="margin-bottom:17px">{data.comments.filter((c: any) => !c.parent_id).length} komentar</h2>
    {#if ctx.user}<form class="form card pad" style="margin-bottom:20px" onsubmit={(e) => comment(e)}>
        <div>
          <label for="comment-body">Ikut percakapan</label><textarea
            id="comment-body"
            bind:value={body}
            maxlength="1500"
            required
            placeholder="Bagikan pendapat atau pengalamanmu..."></textarea>
        </div>
        <button class="button primary" disabled={busy}>Kirim komentar</button>
      </form>{:else}<a class="button" href={'/masuk/customer?next=/post/' + data.post.id}
        >Masuk untuk ikut mengobrol</a
      >{/if}
    <div class="divider"></div>
    {#each data.comments.filter((c: any) => !c.parent_id) as c}<article class="card pad" style="margin-bottom:14px">
        <div class="row">
          <Avatar name={c.author} src={c.avatar_url} />
          <div>
            <strong>{c.author}</strong><span class="meta"> {c.official ? '· Pengelola' : ''}</span>
            <div class="meta">{date(c.created_at)}</div>
          </div>
        </div>
        <p class="post-body" style="margin-top:12px">{c.body}</p>
        <button
          class="action"
          onclick={() => {
            if (ctx.requireLogin()) replyID = replyID === c.id ? '' : c.id;
          }}>Balas</button
        >{#each data.comments.filter((r: any) => r.parent_id === c.id) as r}<div class="reply">
            <strong>{r.author}{r.official ? ' · Pengelola' : ''}</strong>
            <p>{r.body}</p>
          </div>{/each}{#if replyID === c.id}<form
            class="reply-form"
            onsubmit={(e) => comment(e, c.id)}
          >
            <label for={'reply-' + c.id}>Balasan untuk {c.author}</label><textarea
              id={'reply-' + c.id}
              bind:value={reply}
              maxlength="1500"
              required></textarea><button class="button primary" disabled={busy}
              >Kirim balasan</button
            >
          </form>{/if}
        <div class="divider"></div>
      </article>{/each}
  </section>
</section><DiscoveryAside events={data.nearby} organizers={data.organizers} /></div>
