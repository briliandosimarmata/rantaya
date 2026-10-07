<script lang="ts">
  import { getContext } from 'svelte';
  import { untrack } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { ended, date } from '#lib/format.js';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let body = $state(untrack(() => data.review?.body || '')),
    attended = $state(untrack(() => !!data.review)),
    busy = $state(false),
    errorMessage = $state('');
  async function submit(event: SubmitEvent) {
    event.preventDefault();
    busy = true;
    errorMessage = '';
    try {
      await mutate('/reviews', { event_id: data.event.id, body, attended });
      ctx.toast('Ulasan tersimpan.');
      await goto('/event/' + data.event.slug + '?tab=ulasan', { invalidateAll: true });
    } catch (e: any) {
      errorMessage = e.message;
    } finally {
      busy = false;
    }
  }
</script>

<SEO title="Tulis ulasan" privatePage />
{#if ended(data.event)}<form class="editor-surface" onsubmit={submit}>
  <div class="editor-top"><h1>{data.review ? 'Edit ulasan' : 'Bagikan pengalamanmu'}</h1><button class="button primary" disabled={busy}>{busy ? 'Menyimpan…' : 'Kirim ulasan'}</button></div>
  <div class="review-context">Setelah menonton <b>{data.event.title}</b><br /><span class="meta">{date(data.event.starts_at)} · {data.event.organizer.name}</span></div>
  <div class="field"><label for="review-body">Ulasanmu</label><textarea id="review-body" bind:value={body} minlength="15" maxlength="3000" required rows="7" placeholder="Apa yang berkesan? Apa yang bisa diperbaiki?"></textarea></div>
  <label class="small"><input type="checkbox" bind:checked={attended} required /> Saya hadir di pertunjukan ini.</label>
  <p class="meta" style="margin-top:17px">Ulasan yang sama tampil pada event dan penyelenggara, dengan nama event asal. Pernyataan hadir tidak memberi badge verifikasi.</p>
  {#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}
</form>{:else}<div class="page-heading"><h1>Ulasan belum dibuka</h1><p>Gunakan tab Obrolan untuk bertanya sebelum acara.</p></div>{/if}
