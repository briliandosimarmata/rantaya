<script lang="ts">
  import { getContext } from 'svelte';
  import { untrack } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate, upload } from '#lib/api.js';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  const o = untrack(() => data.organizer);
  let name = $state(o.name),
    description = $state(o.description),
    about = $state(o.about),
    city = $state(o.city),
    category = $state(o.category),
    avatar = $state(o.avatar_url),
    cover = $state(o.cover_url),
    busy = $state(false),
    errorMessage = $state('');
  async function photo(e: Event, kind: string) {
    const file = (e.currentTarget as HTMLInputElement).files?.[0];
    if (!file) return;
    busy = true;
    try {
      const url = (await upload(file)).url;
      if (kind === 'avatar') avatar = url;
      else cover = url;
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
        '/organizer',
        { name, description, about, city, category, avatar_url: avatar, cover_url: cover },
        'PATCH',
      );
      await ctx.refresh();
      ctx.toast('Profil ruang tersimpan.');
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<div class="page-heading">
  <h2>Identitas ruang</h2>
  <p>Informasi ini tampil pada profil publik pengelola.</p>
</div>
<form class="narrow card pad form" onsubmit={save}>
  <div>
    <label for="org-name">Nama komunitas</label><input
      id="org-name"
      bind:value={name}
      minlength="2"
      maxlength="80"
      required
    />
  </div>
  <div>
    <label for="org-description">Deskripsi singkat</label><textarea
      id="org-description"
      bind:value={description}
      minlength="10"
      maxlength="500"
      required></textarea>
  </div>
  <div>
    <label for="org-about">Tentang ruang</label><textarea
      id="org-about"
      bind:value={about}
      maxlength="5000"
      rows="6"></textarea>
  </div>
  <div class="form-row">
    <div>
      <label for="org-city">Kota</label><input
        id="org-city"
        bind:value={city}
        minlength="2"
        maxlength="60"
        required
      />
    </div>
    <div>
      <label for="org-category">Kategori</label><select id="org-category" bind:value={category}
        >{#each ['Teater', 'Musik', 'Komedi', 'Lainnya'] as c}<option>{c}</option>{/each}</select
      >
    </div>
  </div>
  <div>
    <label for="org-avatar">Foto profil ruang</label><input
      id="org-avatar"
      type="file"
      accept="image/jpeg,image/png,image/webp"
      onchange={(e) => photo(e, 'avatar')}
    />{#if avatar}<img class="avatar large" src={avatar} alt={name} /><button
        class="action"
        type="button"
        onclick={() => (avatar = '')}>Hapus foto</button
      >{/if}
  </div>
  <div>
    <label for="org-cover">Sampul ruang</label><input
      id="org-cover"
      type="file"
      accept="image/jpeg,image/png,image/webp"
      onchange={(e) => photo(e, 'cover')}
    />{#if cover}<img class="attachment" src={cover} alt={'Sampul ' + name} /><button
        class="action"
        type="button"
        onclick={() => (cover = '')}>Hapus sampul</button
      >{/if}
  </div>
  {#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}<button
    class="button primary"
    disabled={busy}>Simpan profil ruang</button
  >
</form>
