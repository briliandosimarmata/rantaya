<script lang="ts">
  import { getContext } from 'svelte';
  import { goto, invalidateAll } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate, request } from '#lib/api.js';
  import { returnPath } from '#lib/navigation.js';
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  let busy = $state(false),
    errorMessage = $state('');
  let label = $derived(
    data.role === 'organizer' ? 'pengelola' : data.role === 'admin' ? 'moderator' : 'penonton',
  );
  async function demo() {
    busy = true;
    try {
      await mutate('/auth/demo', { role: data.role });
      const session = await request('/auth/me');
      const next = returnPath(data.next, data.role === 'organizer' ? '/kelola' : data.role === 'admin' ? '/admin' : '/');
      await goto(session.user.onboarding_done ? next : '/onboarding?next=' + encodeURIComponent(next), { invalidateAll: true });
    } catch (e: any) {
      errorMessage = e.message;
    } finally {
      busy = false;
    }
  }
</script>

<SEO title={'Masuk sebagai ' + label} privatePage />
<section class="login-surface">
  <p class="eyebrow">{data.role === 'organizer' ? 'AKUN PENGELOLA' : data.role === 'admin' ? 'AKUN MODERATOR' : 'AKUN PENGGUNA'}</p>
  <h1>{data.role === 'customer' ? 'Masuk ke ruangmu' : 'Masuk sebagai ' + label}</h1>
  <p class="muted">
    {data.role === 'organizer'
      ? 'Kelola komunitas, event, dan merchandise.'
      : 'Ikuti komunitas, tulis cerita, dan temukan pentas.'}
  </p>
  <div class="stack">
    {#if data.config.google && data.role !== 'admin'}<a
        class="button dark full"
        href={'/api/auth/google/start?' + new URLSearchParams({ role: data.role, next: returnPath(data.next, data.role === 'organizer' ? '/kelola' : '/') })}>Lanjutkan dengan Google</a
      >{:else if data.role !== 'admin'}<p class="note">
        Login Google akan tersedia setelah konfigurasi akun selesai.
      </p>{/if}{#if data.config.demo}<button
        class="button primary full"
        onclick={demo}
        disabled={busy}>{busy ? 'Membuka akun…' : 'Coba akun demo ' + label}</button
      >
      <p class="hint">
        Akun demo menggunakan data fiktif untuk mencoba aplikasi lokal.
      </p>{/if}{#if !data.config.google && !data.config.demo}<p class="note">
        Pendaftaran belum tersedia. Hubungi pengelola platform.
      </p>{/if}{#if errorMessage || data.loginError}<p class="note error-note" role="alert">
        {errorMessage || 'Login Google dibatalkan. Silakan coba lagi.'}
      </p>{/if}
  </div>
  <div class="divider"></div>
  {#if data.role !== 'admin'}<p class="small">
      {data.role === 'customer' ? 'Mengelola komunitas?' : 'Ingin menonton?'}
      <a href={'/masuk/' + (data.role === 'customer' ? 'organizer' : 'customer')}
        >{data.role === 'customer' ? 'Masuk akun pengelola' : 'Masuk akun penonton'}</a
      >
    </p>
    <p class="hint">Akun penonton dan pengelola terpisah, termasuk profil dan aktivitasnya.</p>{/if}
</section>
