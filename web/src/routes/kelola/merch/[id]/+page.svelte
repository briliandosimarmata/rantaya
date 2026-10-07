<script lang="ts">
  import { getContext } from 'svelte';
  import { untrack } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate, upload } from '#lib/api.js';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  const p = untrack(() => data.product);
  let name = $state(p?.name || ''),
    description = $state(p?.description || ''),
    price = $state(p?.price || 0),
    image = $state(p?.image_url || ''),
    variants = $state(p?.variants.join(', ') || 'Satu ukuran'),
    availability = $state(p?.availability || 'Tersedia'),
    url = $state(p?.purchase_url || ''),
    eventID = $state(p?.event_id || ''),
    busy = $state(false),
    errorMessage = $state('');
  async function photo(e: Event) {
    const file = (e.currentTarget as HTMLInputElement).files?.[0];
    if (!file) return;
    busy = true;
    try {
      image = (await upload(file)).url;
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
  async function save(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    try {
      await mutate(
        '/products' + (p ? '/' + p.id : ''),
        {
          name,
          description,
          price: Number(price),
          image_url: image,
          variants: variants
            .split(',')
            .map((x: string) => x.trim())
            .filter(Boolean),
          availability,
          purchase_url: url,
          event_id: eventID,
        },
        p ? 'PUT' : 'POST',
      );
      ctx.toast('Merchandise tersimpan.');
      await goto('/kelola/merch', { invalidateAll: true });
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<div class="page-heading"><h1>{p ? 'Edit merchandise.' : 'Cerita untuk dibawa pulang.'}</h1></div>
<form class="narrow card pad form" onsubmit={save}>
  <div>
    <label for="product-name">Nama produk</label><input
      id="product-name"
      bind:value={name}
      minlength="2"
      maxlength="100"
      required
    />
  </div>
  <div>
    <label for="product-desc">Deskripsi</label><textarea
      id="product-desc"
      bind:value={description}
      minlength="10"
      maxlength="5000"
      required></textarea>
  </div>
  <div>
    <label for="product-price">Harga · Rp</label><input
      id="product-price"
      type="number"
      bind:value={price}
      min="0"
      max="100000000"
      required
    />
  </div>
  <div>
    <label for="product-photo">Gambar</label><input
      id="product-photo"
      type="file"
      accept="image/jpeg,image/png,image/webp"
      onchange={photo}
    />{#if image}<img class="attachment" src={image} alt={name} /><button
        class="action"
        type="button"
        onclick={() => (image = '')}>Hapus gambar</button
      >{/if}
  </div>
  <div>
    <label for="product-variants">Varian, pisahkan dengan koma</label><input
      id="product-variants"
      bind:value={variants}
      required
      placeholder="S, M, L, XL"
    />
  </div>
  <div>
    <label for="product-availability">Ketersediaan</label><input
      id="product-availability"
      bind:value={availability}
      maxlength="80"
      required
    />
  </div>
  <div>
    <label for="product-event">Event terkait (opsional)</label><select
      id="product-event"
      bind:value={eventID}
      ><option value="">Tidak terkait event</option>{#each data.events as e}<option value={e.id}
          >{e.title}</option
        >{/each}</select
    >
  </div>
  <div>
    <label for="product-url">Tautan pemesanan / WhatsApp (opsional)</label><input
      id="product-url"
      type="url"
      bind:value={url}
      placeholder="https://wa.me/..."
    />
  </div>
  <p class="hint">
    Produk tetap tampil ketika event selesai. Pemesanan merchandise dilanjutkan ke kanal pengelola.
  </p>
  {#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}<button
    class="button primary"
    disabled={busy}>Simpan merchandise</button
  >
</form>
