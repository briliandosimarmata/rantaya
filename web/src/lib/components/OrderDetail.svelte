<script lang="ts">
  import { getContext, untrack } from 'svelte';
  import type { AppContext, Order, PaymentMethod } from '#lib/types.js';
  import { mutate, upload } from '#lib/api.js';
  import { date, time, money, status } from '#lib/format.js';
  import QR from './QR.svelte';
  let {
    order,
    methods = [],
    reviewer = false,
  } = $props<{ order: Order; methods?: PaymentMethod[]; reviewer?: boolean }>();
  const ctx = getContext<AppContext>('app');
  let busy = $state(false),
    errorMessage = $state(''),
    chosen = $state(untrack(() => order.payment_snapshot?.id || '')),
    proof = $state(''),
    unpaid = $state(false),
    received = $state(false),
    note = $state('');
  async function action(path: string, body: any) {
    busy = true;
    errorMessage = '';
    try {
      await mutate('/orders/' + order.id + path, body, path === '/payment' ? 'PATCH' : 'POST');
      await ctx.refresh();
      ctx.toast('Status pesanan diperbarui.');
    } catch (e: any) {
      errorMessage = e.message;
    } finally {
      busy = false;
    }
  }
  async function method(id: string) {
    chosen = id;
    await action('/payment', { payment_method_id: id });
  }
  async function proofFile(e: Event) {
    const file = (e.currentTarget as HTMLInputElement).files?.[0];
    if (!file) return;
    busy = true;
    errorMessage = '';
    try {
      proof = (await upload(file, 'proof')).id;
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
  async function copy() {
    try {
      await navigator.clipboard.writeText(order.payment_snapshot?.number || '');
      ctx.toast('Nomor pembayaran disalin.');
    } catch {
      ctx.toast('Salin nomor pembayaran secara manual.');
    }
  }
</script>

<section class="narrow stack">
  <div class="card pad">
    <div class="between">
      <span class="pill highlight">{status[order.status]}</span><span class="meta"
        >#{order.id.slice(0, 8).toUpperCase()}</span
      >
    </div>
    <h1 style="font-size:1.7rem;margin-top:20px">{order.event.title}</h1>
    <div class="order-summary">
      <p>
        {order.session.label}<br />{date(order.session.starts_at)} · {time(order.session.starts_at)}
      </p>
      <div class="between">
        <span>{order.quantity} tiket · {reviewer ? order.customer.name : order.event.venue}</span
        ><strong>{money(order.total)}</strong>
      </div>
    </div>
  </div>
  {#if !reviewer && (order.status === 'awaiting_payment' || order.status === 'correction_requested')}<section
      class="card pad"
    >
      <h2>
        {order.status === 'correction_requested' ? 'Perbaiki bukti transfer' : 'Bayar ke pengelola'}
      </h2>
      {#if ctx.demo}<p class="note">
          Rekening dan wallet contoh. Jangan transfer uang sungguhan.
        </p>{/if}{#if order.expires_at}<p class="hint">
          Unggah dan kirim bukti sebelum {date(order.expires_at, {
            day: 'numeric',
            month: 'short',
          })} pukul {time(order.expires_at)}. Jika sudah transfer tetapi waktunya berakhir, hubungi
          pengelola.
        </p>{/if}{#if order.status === 'awaiting_payment'}<div class="stack">
          {#each methods as m}<label class:selected={chosen === m.id} class="payment-option"
              ><input
                type="radio"
                name="payment"
                value={m.id}
                checked={chosen === m.id}
                disabled={busy}
                onchange={() => method(m.id)}
              />
              <div>
                <strong>{m.provider}</strong>
                <div class="meta">
                  {m.kind === 'bank' ? 'Transfer bank' : 'Wallet'} · {m.holder}
                </div>
              </div></label
            >{/each}
        </div>{/if}{#if order.payment_snapshot}<div class="note" style="margin-top:18px">
          <strong>{order.payment_snapshot.provider}</strong>
          <div class="between">
            <span style="overflow-wrap:anywhere">{order.payment_snapshot.number}</span><button
              class="action"
              onclick={copy}>Salin</button
            >
          </div>
          <div>{order.payment_snapshot.holder}</div>
          <p class="hint">{order.payment_snapshot.note}</p>
        </div>{/if}{#if order.note}<p class="note error-note" style="margin-top:18px">
          {order.note}
        </p>{/if}
      <div class="divider"></div>
      <label for="proof-file">Bukti transfer</label><input
        id="proof-file"
        type="file"
        accept="image/jpeg,image/png,image/webp"
        onchange={proofFile}
        disabled={busy || !order.payment_snapshot}
      />
      <p class="hint">
        JPG, PNG, atau WebP, maksimal 8 MB. Bukti hanya dapat dilihat oleh kamu dan pengelola event.
      </p>
      {#if proof}<img
          class="proof-image"
          src={'/api/uploads/' + proof}
          alt="Bukti transfer yang dipilih"
        /><button
          class="button primary full"
          style="margin-top:16px"
          disabled={busy}
          onclick={() => action('/proof', { upload_id: proof })}>Kirim bukti transfer</button
        >{/if}
      <p class="hint">Kamu bisa menutup halaman dan melanjutkan lewat History transaksi.</p>
      {#if order.status === 'awaiting_payment'}<details>
          <summary>Batalkan pembelian</summary><label class="inline-label"
            ><input type="checkbox" bind:checked={unpaid} />Saya belum mentransfer pembayaran.</label
          ><button
            class="button"
            disabled={busy || !unpaid}
            onclick={() => action('/cancel', { confirm_unpaid: true })}>Batalkan pesanan</button
          >
        </details>{/if}
    </section>{/if}
  {#if order.status === 'awaiting_review'}<div class="card pad">
      <h2>{reviewer ? 'Periksa pembayaran' : 'Menunggu konfirmasi pengelola'}</h2>
      <p>
        {reviewer
          ? 'Cocokkan bukti dengan mutasi rekening. Persetujuan akan menerbitkan tiket.'
          : 'Bukti transfer sudah diterima. Kami akan mengabari setelah pengelola memeriksa dana.'}
      </p>
      {#if order.proof_id}<img
          class="proof-image"
          src={'/api/uploads/' + order.proof_id}
          alt="Bukti transfer pesanan"
        />{/if}{#if reviewer}<div class="divider"></div>
        {#if order.payment_snapshot}<p>
            {order.payment_snapshot.provider} · {order.payment_snapshot.number}<br />{order
              .payment_snapshot.holder}
          </p>{/if}<label class="inline-label"
          ><input type="checkbox" bind:checked={received} />Saya sudah memeriksa mutasi dan dana
          benar-benar masuk.</label
        ><button
          class="button primary full"
          disabled={busy || !received}
          onclick={() => action('/review', { action: 'approve', received_funds: true, note: '' })}
          >Setujui & terbitkan tiket</button
        >
        <details style="margin-top:22px">
          <summary>Minta perbaikan bukti</summary><label for="correction-note"
            >Alasan perbaikan</label
          ><textarea id="correction-note" bind:value={note} minlength="10" maxlength="500"
          ></textarea><button
            class="button"
            disabled={busy || note.trim().length < 10}
            onclick={() => action('/review', { action: 'correction', note, received_funds: false })}
            >Kirim permintaan perbaikan</button
          >
        </details>{/if}
    </div>{/if}
  {#if order.status === 'approved'}<section class="stack">
      {#each order.tickets as ticket}<article class="card pad ticket">
          <span class="pill"
            >Tiket {ticket.ordinal} dari {order.quantity} · {ticket.checked_at
              ? 'Sudah digunakan'
              : 'Siap digunakan'}</span
          >
          <h2>{order.event.title}</h2>
          <QR code={ticket.qr} />
          <details class="no-print" style="width:100%;overflow-wrap:anywhere">
            <summary>Kode manual tiket</summary><code>{ticket.qr.replace('ruang:ticket:', '')}</code
            >
          </details>
          <p class="hint">Tunjukkan QR saat masuk. Jangan bagikan kepada orang lain.</p>
        </article>{/each}<button class="button no-print" onclick={() => window.print()}
        >Cetak / simpan tiket sebagai PDF</button
      >
    </section>{/if}
  {#if order.status === 'cancelled' || order.status === 'expired'}<div class="note">
      <p>
        {order.status === 'expired'
          ? 'Batas pembayaran telah berakhir. Jika kamu sudah transfer, hubungi pengelola dan sebutkan nomor pesanan ini.'
          : 'Pesanan dibatalkan. Kamu dapat memilih sesi lagi.'}
      </p>
      <a class="button" href={'/ruang/' + order.event.organizer.slug}>Buka ruang pengelola</a>
    </div>{/if}{#if reviewer && order.status === 'correction_requested'}<p class="note">
      Menunggu penonton mengirim ulang bukti. Kuota tetap ditahan. {order.note}
    </p>{/if}{#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}
  <section class="card pad">
    <h2>Perjalanan pesanan</h2>
    <ol class="timeline">
      {#each order.history as h}<li>
          <strong>{status[h.status]}</strong><small
            >{date(h.created_at, {
              day: 'numeric',
              month: 'short',
              hour: '2-digit',
              minute: '2-digit',
            })}</small
          ><span>{h.note}</span>
        </li>{/each}
    </ol>
  </section>
</section>
