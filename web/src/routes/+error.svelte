<script lang="ts">
  import { page } from '$app/state';
  import SEO from '#lib/components/SEO.svelte';
</script>

<SEO title="Halaman belum dapat dibuka" privatePage />
<section class="card pad narrow">
  <p class="eyebrow">{page.status}</p>
  <h1>Halaman belum dapat dibuka.</h1>
  <p>{page.error?.message || 'Coba lagi sebentar.'}</p>
  <a class="button primary" href="/">Kembali ke beranda</a>
</section>
