<script lang="ts">
  import { getContext, onMount, tick } from 'svelte';
  import { goto } from '$app/navigation';
  import { request, mutate, upload } from '#lib/api.js';
  import type { AppContext, Mention } from '#lib/types.js';
  import Avatar from '#lib/components/Avatar.svelte';
  import Icon from '#lib/components/Icon.svelte';
  import SEO from '#lib/components/SEO.svelte';
  import { mentionPlacement } from '#lib/mention-placement.js';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let body = $state(''),
    title = $state(''),
    image = $state(''),
    link = $state(''),
    label = $state(''),
    mentions = $state<Mention[]>([]),
    titleOpen = $state(false),
    linkOpen = $state(false),
    busy = $state(false),
    uploading = $state(false),
    errorMessage = $state(''),
    ready = $state(false),
    searchTerm = $state(''),
    mentionOpen = $state(false),
    queryVersion = 0,
    results = $state<Mention[]>([]),
    selected = $state(0),
    panelTop = $state(0),
    panelLeft = $state(0),
    panelWidth = $state(300),
    panelHeight = $state(200),
    queryStart = 0;
  let textarea: HTMLTextAreaElement;
  let fileInput: HTMLInputElement;
  let editor: HTMLDivElement;
  let timer: ReturnType<typeof setTimeout>;
  let abort: AbortController | undefined;
  let draftKey = $derived('ruang-draft:' + ctx.user?.id);
  function resize() {
    if (textarea) {
      textarea.style.height = 'auto';
      textarea.style.height = Math.max(220, textarea.scrollHeight) + 'px';
    }
  }
  function position() {
    if (!textarea || !editor) return;
    const style = getComputedStyle(textarea);
    const mirror = document.createElement('div');
    for (const k of [
      'font',
      'lineHeight',
      'letterSpacing',
      'padding',
      'boxSizing',
      'wordBreak',
      'overflowWrap',
    ] as const)
      mirror.style[k] = style[k];
    mirror.style.cssText +=
      ';position:absolute;visibility:hidden;white-space:pre-wrap;width:' +
      textarea.clientWidth +
      'px';
    mirror.textContent = textarea.value.slice(0, textarea.selectionStart);
    const marker = document.createElement('span');
    marker.textContent = '.';
    mirror.append(marker);
    document.body.append(mirror);
    const caret = marker.offsetTop;
    mirror.remove();
    const rect = textarea.getBoundingClientRect(),
      vv = window.visualViewport;
    const placement = mentionPlacement(
      rect,
      {
        top: rect.top + caret - textarea.scrollTop,
        height: parseFloat(style.lineHeight),
      },
      {
        left: vv?.offsetLeft || 0,
        top: vv?.offsetTop || 0,
        width: vv?.width || innerWidth,
        height: vv?.height || innerHeight,
      },
      document.querySelector('.write-header')?.getBoundingClientRect().bottom || 0,
      Math.max(64, results.length * 64),
    );
    panelTop = placement.top;
    panelLeft = placement.left;
    panelWidth = placement.width;
    panelHeight = placement.height;
  }
  function closeMention() {
    queryVersion++;
    clearTimeout(timer);
    abort?.abort();
    results = [];
    searchTerm = '';
    mentionOpen = false;
  }

  async function inspect() {
    resize();
    if (!textarea) return;
    const before = body.slice(0, textarea.selectionStart);
    const match = /(?:^|\s)@([^\s@]*)$/.exec(before);
    clearTimeout(timer);
    abort?.abort();
    const version = ++queryVersion;
    if (!match) {
      closeMention();
      return;
    }
    queryStart = before.lastIndexOf('@');
    searchTerm = match[1];
    selected = 0;
    position();
    const term = searchTerm;
    timer = setTimeout(async () => {
      abort = new AbortController();
      try {
        const found = await request<Mention[]>('/mentions?q=' + encodeURIComponent(term), {
          signal: abort.signal,
        });
        if (queryVersion === version && searchTerm === term) {
          results = found;
          mentionOpen = true;
          await tick();
          position();
        }
      } catch (e: any) {
        if (e.name !== 'AbortError') errorMessage = e.message;
      }
    }, 140);
  }
  async function pick(m: Mention) {
    const end = textarea.selectionStart;
    const insertion = '@' + m.name + ' ';
    body = body.slice(0, queryStart) + insertion + body.slice(end);
    if (!mentions.some((x) => x.kind === m.kind && x.id === m.id)) mentions = [...mentions, m];
    closeMention();
    await tick();
    textarea.focus();
    textarea.setSelectionRange(queryStart + insertion.length, queryStart + insertion.length);
    resize();
  }
  function keys(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      closeMention();
      return;
    }
    if (!results.length) return;
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
      e.preventDefault();
      selected = (selected + (e.key === 'ArrowDown' ? 1 : -1) + results.length) % results.length;
      document.getElementById('mention-' + selected)?.scrollIntoView({ block: 'nearest' });
    } else if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      pick(results[selected]);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      results = [];
      searchTerm = '';
    }
  }
  async function insertAt() {
    textarea.focus();
    const start = textarea.selectionStart,
      end = textarea.selectionEnd;
    body = body.slice(0, start) + (start && body[start - 1] !== ' ' ? ' @' : '@') + body.slice(end);
    await tick();
    textarea.setSelectionRange(
      start + (start && body[start - 1] !== ' ' ? 2 : 1),
      start + (start && body[start - 1] !== ' ' ? 2 : 1),
    );
    await inspect();
  }
  async function chooseFile(e: Event) {
    const file = (e.currentTarget as HTMLInputElement).files?.[0];
    if (!file) return;
    uploading = true;
    errorMessage = '';
    try {
      image = (await upload(file)).url;
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      uploading = false;
    }
  }
  async function submit(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    errorMessage = '';
    try {
      const post = await mutate('/posts', {
        title: titleOpen ? title : '',
        body,
        image_url: image,
        city: ctx.city,
        link_url: ctx.user?.role === 'organizer' ? link : '',
        link_label: label,
        mentions: mentions.map((m) => ({ kind: m.kind, id: m.id })),
      });
      ready = false;
      localStorage.removeItem(draftKey);
      ctx.toast('Percakapanmu sudah dibagikan.');
      await goto('/post/' + post.id, { invalidateAll: true });
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
  $effect(() => {
    if (ready) {
      const draft = { body, title, image, link, label, mentions, titleOpen, linkOpen };
      try {
        localStorage.setItem(draftKey, JSON.stringify(draft));
      } catch {}
    }
  });
  onMount(() => {
    try {
      const d = JSON.parse(localStorage.getItem(draftKey) || 'null');
      if (d) {
        body = d.body || '';
        title = d.title || '';
        image = d.image || '';
        link = d.link || '';
        label = d.label || '';
        mentions = Array.isArray(d.mentions) ? d.mentions : [];
        titleOpen = !!d.titleOpen;
        linkOpen = !!d.linkOpen;
      }
    } catch {}
    ready = true;
    textarea.focus();
    resize();
    (async () => {
      if (data.event) {
        const e = await request('/events/' + data.event);
        if (!mentions.some((m) => m.id === e.id && m.kind === 'event'))
          mentions = [...mentions, { kind: 'event', id: e.id, slug: e.slug, name: e.title }];
      }
      if (data.organizer) {
        const o = await request('/organizers/' + data.organizer);
        if (!mentions.some((m) => m.id === o.id && m.kind === 'organizer'))
          mentions = [...mentions, { kind: 'organizer', id: o.id, slug: o.slug, name: o.name }];
      }
    })().catch((e) => (errorMessage = e.message));
    const handler = () => {
      if (mentionOpen) position();
    };
    window.addEventListener('resize', handler);
    window.addEventListener('scroll', handler, true);
    window.visualViewport?.addEventListener('resize', handler);
    window.visualViewport?.addEventListener('scroll', handler);
    return () => {
      clearTimeout(timer);
      abort?.abort();
      window.removeEventListener('resize', handler);
      window.removeEventListener('scroll', handler, true);
      window.visualViewport?.removeEventListener('resize', handler);
      window.visualViewport?.removeEventListener('scroll', handler);
    };
  });
</script>

<SEO title="Buat postingan" privatePage />
<form id="compose-form" class="write-surface" onsubmit={submit}>
  <div class="write-header">
    <span class="small muted">Buat postingan</span><button
      type="submit"
      class="button primary"
      disabled={busy || uploading || !body.trim()}
      >{busy ? 'Mengirim…' : uploading ? 'Mengunggah…' : 'Posting'}</button
    >
  </div>
  <h1 class="sr-only">Buat postingan</h1>
  <div class="write-main">
    <Avatar name={ctx.user?.name || ''} src={ctx.user?.avatar_url || ''} size={'lime ' + (ctx.user?.role === 'organizer' ? 'square' : '')} />
    <div class="write-content" bind:this={editor}>
      <strong>{ctx.user?.name}</strong>
      <div class="meta">
        {ctx.user?.role === 'organizer' ? 'Pengelola' : 'Akun pengguna'} · {ctx.city}
      </div>
      <label class="sr-only" for="post-body">Isi postingan</label><textarea
        id="post-body"
        bind:this={textarea}
        bind:value={body}
        oninput={inspect}
        onclick={inspect}
        onkeydown={keys}
        maxlength="5000"
        required
        placeholder="Mau cerita apa?"
        aria-autocomplete="list"
        aria-haspopup="listbox"
        aria-controls={mentionOpen ? 'mention-results' : undefined}
        aria-activedescendant={results.length && mentionOpen ? 'mention-' + selected : undefined}
      ></textarea>
      {#if titleOpen}<div class="optional-title">
          <label class="sr-only" for="post-title">Judul opsional</label><input
            id="post-title"
            bind:value={title}
            maxlength="120"
            placeholder="Judul (opsional)"
          />
        </div>{/if}
      {#if mentions.length}<div class="chips">
          {#each mentions as m}<span class="mention"
              >@{m.name}<button
                type="button"
                class="action"
                aria-label={'Hapus mention ' + m.name}
                onclick={() => (mentions = mentions.filter((x) => x !== m))}
                ><Icon name="close" size={14} /></button
              ></span
            >{/each}
        </div>{/if}
      {#if image}<div class="compose-photo">
          <img class="attachment" src={image} alt="Gambar yang akan diposting" /><button
            class="icon-button"
            type="button"
            aria-label="Hapus gambar"
            onclick={() => (image = '')}><Icon name="close" /></button
          >
        </div>{/if}
      {#if linkOpen && ctx.user?.role === 'organizer'}<div class="purchase-fields form-row">
          <div>
            <label for="post-link">Tautan tiket / merchandise</label><input
              id="post-link"
              type="url"
              bind:value={link}
              placeholder="https://..."
            />
          </div>
          <div>
            <label for="post-link-label">Teks tombol (opsional)</label><input
              id="post-link-label"
              bind:value={label}
              maxlength="80"
              placeholder="Beli tiket"
            />
          </div>
        </div>{/if}
    </div>
  </div>
  <input
    bind:this={fileInput}
    type="file"
    accept="image/jpeg,image/png,image/webp"
    hidden
    onchange={chooseFile}
  />
  <div class="write-toolbar">
    <button
      class="icon-button"
      type="button"
      onclick={() => fileInput.click()}
      aria-label="Unggah gambar"
      disabled={uploading}><Icon name="image" /></button
    >
    <button
      class="icon-button mention-tool"
      type="button"
      onclick={insertAt}
      aria-label="Mention event atau pengelola">@</button
    >
    {#if ctx.user?.role === 'organizer'}<button
        class="icon-button"
        type="button"
        onclick={() => (linkOpen = !linkOpen)}
        aria-label="Tambahkan tautan pembelian"><Icon name="link" /></button
      >{/if}
    <button
      class="action title-toggle"
      type="button"
      onclick={() => (titleOpen = !titleOpen)}
      aria-label="Tambahkan judul opsional">{titleOpen ? 'Hapus judul' : 'Tambah judul'}</button
    >
    <span class="meta">Draft tersimpan</span>
  </div>
  {#if mentionOpen}<div
      class="mention-panel"
      id="mention-results"
      role="listbox"
      aria-label="Saran mention"
      style={'top:' +
        panelTop +
        'px;left:' +
        panelLeft +
        'px;width:' +
        panelWidth +
        'px;max-height:' +
        panelHeight +
        'px'}
    >
      {#each results as m, i}<button
          type="button"
          id={'mention-' + i}
          class:selected={i === selected}
          role="option"
          aria-selected={i === selected}
          onpointerdown={(e) => e.preventDefault()}
          onclick={() => pick(m)}
        >
          <span class="avatar square"
            ><Icon name={m.kind === 'event' ? 'calendar' : 'people'} /></span
          >
          <span
            ><strong>{m.name}</strong><span class="meta"
              >{m.label || (m.kind === 'event' ? 'Event' : 'Pengelola')}{m.city
                ? ' · ' + m.city
                : ''}</span
            ></span
          >
        </button>{:else}<p class="mention-empty" role="status">
          Tidak ada hasil. Coba nama lain.
        </p>{/each}
    </div>{/if}
  {#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}
</form>
