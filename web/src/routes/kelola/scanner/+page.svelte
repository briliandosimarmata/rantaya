<script lang="ts">
  import { getContext, onDestroy } from 'svelte';
  import { untrack } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  const ctx = getContext<AppContext>('app');
  let { data } = $props();
  let eventID = $state(untrack(() => data.events[0]?.id || '')),
    code = $state(''),
    busy = $state(false),
    running = $state(false),
    message = $state(''),
    valid = $state(false);
  let camera: HTMLVideoElement;
  let stream: MediaStream | undefined;
  let frame = 0;
  let last = 0;
  function stop() {
    running = false;
    if (typeof window !== 'undefined') cancelAnimationFrame(frame);
    stream?.getTracks().forEach((t) => t.stop());
    stream = undefined;
    if (camera) camera.srcObject = null;
  }
  async function validate(e?: SubmitEvent) {
    e?.preventDefault();
    stop();
    busy = true;
    message = '';
    try {
      const result = await mutate('/check-in', { event_id: eventID, code });
      valid = true;
      message = result.name + ' — ' + result.message;
      code = '';
      await ctx.refresh();
    } catch (err: any) {
      valid = false;
      message = err.message;
    } finally {
      busy = false;
    }
  }
  async function start() {
    busy = true;
    message = '';
    try {
      if (!navigator.mediaDevices?.getUserMedia)
        throw new Error('Kamera tidak tersedia. Gunakan kode manual.');
      stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: { ideal: 'environment' }, width: { ideal: 1280 } },
        audio: false,
      });
      camera.srcObject = stream;
      await camera.play();
      const jsQR = (await import('jsqr')).default;
      const canvas = document.createElement('canvas');
      const draw = canvas.getContext('2d', { willReadFrequently: true })!;
      running = true;
      const scan = (now: number) => {
        if (!running) return;
        if (now - last > 200 && camera.readyState >= 2) {
          last = now;
          canvas.width = 640;
          canvas.height = Math.round((640 * camera.videoHeight) / camera.videoWidth);
          draw.drawImage(camera, 0, 0, canvas.width, canvas.height);
          const image = draw.getImageData(0, 0, canvas.width, canvas.height);
          const qr = jsQR(image.data, image.width, image.height, {
            inversionAttempts: 'dontInvert',
          });
          if (qr) {
            code = qr.data;
            validate();
            return;
          }
        }
        frame = requestAnimationFrame(scan);
      };
      frame = requestAnimationFrame(scan);
    } catch (err: any) {
      stop();
      valid = false;
      message =
        err.name === 'NotAllowedError' ? 'Izin kamera ditolak. Gunakan kode manual.' : err.message;
    } finally {
      busy = false;
    }
  }
  onDestroy(stop);
</script>

<div class="page-heading">
  <h1>Selamat datang di panggung.</h1>
  <p>Scan QR untuk mencatat satu kali masuk. Check-in memerlukan koneksi ke server.</p>
</div>
<section class="narrow card pad stack">
  <div>
    <label for="scan-event">Event yang sedang dibuka</label><select
      id="scan-event"
      bind:value={eventID}
      onchange={stop}
      >{#each data.events as e}<option value={e.id}>{e.title}</option>{/each}</select
    >
    <p class="hint">Gate tersedia dua jam sebelum sesi sampai dua jam setelah selesai.</p>
  </div>
  <video
    class="scanner-video"
    bind:this={camera}
    muted
    playsinline
    aria-label="Pratinjau kamera scanner"><track kind="captions" /></video
  >{#if running}<button class="button" onclick={stop}>Hentikan kamera</button>{:else}<button
      class="button primary"
      onclick={start}
      disabled={busy || !eventID}>Buka kamera & scan QR</button
    >{/if}
  <form class="form" onsubmit={validate}>
    <div>
      <label for="ticket-code">Kode manual tiket</label><input
        id="ticket-code"
        bind:value={code}
        placeholder="Tempel kode dari detail tiket"
        required
        autocomplete="off"
      />
      <p class="hint">Gunakan kode manual jika kamera sulit membaca QR.</p>
    </div>
    <button class="button" disabled={busy || !eventID}
      >{busy ? 'Memeriksa…' : 'Validasi & catat masuk'}</button
    >
  </form>
  {#if message}<div class="note" role="status" aria-live="polite">
      <strong>{valid ? 'Berhasil masuk' : 'Tiket belum dapat digunakan'}</strong>
      <p>{message}</p>
    </div>{/if}
</section>
