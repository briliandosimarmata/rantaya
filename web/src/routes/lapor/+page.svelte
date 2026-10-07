<script lang="ts">
  import { getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let reason = $state(''),
    busy = $state(false),
    errorMessage = $state('');
  async function submit(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    try {
      await mutate('/reports', { kind: data.kind, target_id: data.id, reason });
      ctx.toast('Laporan diterima untuk ditinjau moderator.');
      await goto('/');
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<SEO title="Laporkan konten" privatePage />
<div class="page-heading">
  <h1>Laporkan konten.</h1>
  <p>Jelaskan masalahnya agar moderator dapat meninjau dengan konteks.</p>
</div>
<form class="form narrow card pad" onsubmit={submit}>
  <div>
    <label for="reason">Alasan laporan</label><textarea
      id="reason"
      bind:value={reason}
      minlength="10"
      maxlength="1000"
      required
      rows="6"></textarea>
  </div>
  {#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}<button
    class="button primary"
    disabled={busy}>Kirim laporan</button
  >
</form>
