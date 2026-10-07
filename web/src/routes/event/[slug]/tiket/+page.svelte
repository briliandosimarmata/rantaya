<script lang="ts">
  import { getContext } from 'svelte';
  import { untrack } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import { createIdempotencyKey } from '#lib/idempotency.js';
  import { date, time, money } from '#lib/format.js';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let sessionID = $state(
      untrack(
        () =>
          data.event.sessions.find(
            (s: any) =>
              s.id === data.selectedSession &&
              new Date(s.starts_at) > new Date() &&
              s.available > 0,
          )?.id ||
          data.event.sessions.find(
            (s: any) => new Date(s.starts_at) > new Date() && s.available > 0,
          )?.id ||
          '',
      ),
    ),
    quantity = $state(1),
    busy = $state(false),
    errorMessage = $state(''),
    key = '';
  let chosen = $derived(data.event.sessions.find((s: any) => s.id === sessionID));
  async function buy(e: SubmitEvent) {
    e.preventDefault();
    if (busy) return;
    busy = true;
    errorMessage = '';
    try {
      key = key || createIdempotencyKey();
      const order = await mutate('/orders', {
        session_id: sessionID,
        quantity,
        idempotency_key: key,
      });
      await goto('/transaksi/' + order.id, { invalidateAll: true });
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<SEO title={'Tiket ' + data.event.title} privatePage />
<div class="page-heading">
  <p class="eyebrow">Pilih pertunjukanmu</p>
  <h1>{data.event.title}</h1>
  <p>{data.event.venue} · {data.event.city}</p>
</div>
<form class="narrow form card pad" onsubmit={buy}>
  <fieldset style="border:0;padding:0">
    <legend><strong>Sesi pertunjukan</strong></legend>
    <div class="stack" style="margin-top:14px">
      {#each data.event.sessions as s}<label
          class:selected={sessionID === s.id}
          class="session-option"
          ><input
            type="radio"
            name="session"
            bind:group={sessionID}
            value={s.id}
            disabled={s.available <= 0 || new Date(s.starts_at) <= new Date()}
            onchange={() => (key = '')}
          />
          <div class="grow">
            <strong>{s.label}</strong>
            <div class="meta">{date(s.starts_at)} · {time(s.starts_at)}</div>
            <span class="small">{s.available} tiket tersedia</span>
          </div>
          <strong>{money(s.price)}</strong></label
        >{/each}
    </div>
  </fieldset>
  <div>
    <label for="quantity">Jumlah tiket</label><input
      id="quantity"
      type="number"
      bind:value={quantity}
      min="1"
      max={Math.min(6, chosen?.available || 1)}
      required
      onchange={() => (key = '')}
    />
    <p class="hint">Satu QR akan diterbitkan untuk setiap tiket.</p>
  </div>
  <div class="divider"></div>
  <div class="between">
    <strong>Total</strong><strong>{money((chosen?.price || 0) * quantity)}</strong>
  </div>
  <p class="hint">
    {chosen?.price === 0
      ? 'Tiket gratis akan langsung diterbitkan.'
      : 'Pembayaran lewat transfer langsung ke pengelola. Tiket diterbitkan setelah dana dikonfirmasi.'}
  </p>
  {#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}<button
    class="button primary"
    disabled={busy || !chosen}>{busy ? 'Membuat pesanan…' : 'Beli & lanjutkan pembayaran'}</button
  >
</form>
