<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext, PaymentMethod } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let id = $state(''),
    kind = $state('bank'),
    provider = $state(''),
    number = $state(''),
    holder = $state(''),
    note = $state(''),
    enabled = $state(true),
    busy = $state(false),
    errorMessage = $state('');
  function edit(m: PaymentMethod) {
    id = m.id;
    kind = m.kind;
    provider = m.provider;
    number = m.number;
    holder = m.holder;
    note = m.note;
    enabled = m.enabled;
    document.getElementById('payment-master')?.scrollIntoView({ behavior: 'smooth' });
  }
  function reset() {
    id = '';
    provider = '';
    number = '';
    holder = '';
    note = '';
    enabled = true;
  }
  async function submit(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    errorMessage = '';
    try {
      await mutate(
        '/payment-methods' + (id ? '/' + id : ''),
        { kind, provider, number, holder, note, enabled },
        id ? 'PUT' : 'POST',
      );
      reset();
      await ctx.refresh();
      ctx.toast('Metode pembayaran tersimpan.');
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<div class="page-heading">
  <h1>Rekening & wallet.</h1>
  <p>
    Pembayaran langsung ke pengelola. Detail pada pesanan yang sudah memilih rekening tetap
    menggunakan snapshot sebelumnya.
  </p>
</div>
<section class="narrow stack">
  {#each data.methods as m}<article class="card pad">
      <div class="between">
        <strong>{m.provider}</strong><span class="pill">{m.enabled ? 'Aktif' : 'Nonaktif'}</span>
      </div>
      <p>{m.number}<br />{m.holder}</p>
      <button class="button" onclick={() => edit(m)}>Edit metode</button>
    </article>{/each}
  <form id="payment-master" class="card pad form" onsubmit={submit}>
    <h2>{id ? 'Edit metode pembayaran' : 'Tambah metode pembayaran'}</h2>
    <div>
      <label for="payment-kind">Jenis</label><select id="payment-kind" bind:value={kind}
        ><option value="bank">Rekening bank</option><option value="wallet">Wallet</option></select
      >
    </div>
    <div>
      <label for="provider">Bank / penyedia wallet</label><input
        id="provider"
        bind:value={provider}
        minlength="2"
        maxlength="60"
        required
      />
    </div>
    <div>
      <label for="number">Nomor rekening / wallet</label><input
        id="number"
        bind:value={number}
        minlength="3"
        maxlength="80"
        required
        inputmode="numeric"
      />
    </div>
    <div>
      <label for="holder">Nama pemilik</label><input
        id="holder"
        bind:value={holder}
        minlength="2"
        maxlength="100"
        required
      />
    </div>
    <div>
      <label for="payment-note">Instruksi tambahan</label><textarea
        id="payment-note"
        bind:value={note}
        maxlength="500"></textarea>
    </div>
    <label class="inline-label"
      ><input type="checkbox" bind:checked={enabled} />Aktif untuk pesanan baru</label
    >{#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}<button
      class="button primary"
      disabled={busy}>{busy ? 'Menyimpan…' : 'Simpan metode'}</button
    >{#if id}<button class="action" type="button" onclick={reset}>Batal edit</button>{/if}
  </form>
</section>
