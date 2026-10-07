<script lang="ts">
  let { code } = $props<{ code: string }>();
  let canvas: HTMLCanvasElement;
  let errorMessage = $state('');
  $effect(() => {
    const value = code;
    if (canvas)
      import('qrcode')
        .then((q) =>
          q.default.toCanvas(canvas, value, {
            width: 280,
            margin: 3,
            errorCorrectionLevel: 'M',
            color: { dark: '#172322', light: '#F4F5F7' },
          }),
        )
        .catch(() => (errorMessage = 'QR belum bisa ditampilkan. Gunakan kode tiket.'));
  });
  function download() {
    const a = document.createElement('a');
    a.href = canvas.toDataURL('image/png');
    a.download = 'tiket-ruang.png';
    a.click();
  }
</script>

<canvas bind:this={canvas} aria-label="QR tiket masuk"></canvas>{#if errorMessage}<p role="alert">
    {errorMessage}
  </p>{/if}<button class="button no-print" onclick={download}>Unduh QR</button>
