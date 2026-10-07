<script lang="ts">
  import { getContext } from 'svelte';
  import { untrack } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate, upload } from '#lib/api.js';
  import { localDateValue } from '#lib/format.js';
  let { data } = $props();
  const ctx = getContext<AppContext>('app');
  const source = untrack(() => data.event);
  let title = $state(source?.title || ''),
    description = $state(source?.description || ''),
    category = $state(source?.category || 'Teater'),
    city = $state(source?.city || 'Karawang'),
    venue = $state(source?.venue || ''),
    address = $state(source?.address || ''),
    maps = $state(source?.maps_url || ''),
    duration = $state(source?.duration || '90 menit'),
    language = $state(source?.language || 'Bahasa Indonesia'),
    age = $state(source?.age || 'Semua usia'),
    flyer = $state(source?.flyer_url || ''),
    trailer = $state(source?.trailer_url || ''),
    layout = $state(source?.layout_url || ''),
    published = $state(source?.published || false),
    lineup = $state<any[]>(source?.lineup || []),
    sessions = $state<any[]>(
      source?.sessions.map((s: any) => ({
        ...s,
        starts_at: localDateValue(s.starts_at),
        ends_at: localDateValue(s.ends_at),
      })) || [
        {
          id: '',
          label: 'Pertunjukan utama',
          starts_at: '',
          ends_at: '',
          price: 45000,
          capacity: 60,
        },
      ],
    ),
    busy = $state(false),
    errorMessage = $state('');
  async function photo(e: Event, kind: string, index = -1) {
    const file = (e.currentTarget as HTMLInputElement).files?.[0];
    if (!file) return;
    busy = true;
    errorMessage = '';
    try {
      const image = (await upload(file)).url;
      if (kind === 'flyer') flyer = image;
      else if (kind === 'layout') layout = image;
      else lineup[index].photo = image;
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
  async function submit(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    errorMessage = '';
    try {
      const event = await mutate(
        '/events' + (source ? '/' + source.id : ''),
        {
          title,
          description,
          category,
          city,
          venue,
          address,
          maps_url: maps,
          duration,
          language,
          age,
          flyer_url: flyer,
          trailer_url: trailer,
          layout_url: layout,
          published,
          lineup,
          sessions: sessions.map((s) => ({
            id: s.id || '',
            label: s.label,
            starts_at: new Date(s.starts_at + ':00+07:00').toISOString(),
            ends_at: new Date(s.ends_at + ':00+07:00').toISOString(),
            price: Number(s.price),
            capacity: Number(s.capacity),
          })),
        },
        source ? 'PUT' : 'POST',
      );
      ctx.toast('Event tersimpan.');
      await goto('/kelola/event', { invalidateAll: true });
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<div class="page-heading">
  <h1>{source ? 'Edit event.' : 'Buat pertemuan baru.'}</h1>
  <p>Jadwal menggunakan WIB. Judul dan URL event tetap mudah dibagikan.</p>
</div>
<form class="form narrow card pad" onsubmit={submit}>
  <div>
    <label for="event-title">Nama event</label><input
      id="event-title"
      bind:value={title}
      minlength="3"
      maxlength="120"
      required
    />
  </div>
  <div>
    <label for="event-description">Deskripsi</label><textarea
      id="event-description"
      bind:value={description}
      minlength="10"
      maxlength="10000"
      rows="6"
      required></textarea>
  </div>
  <div class="form-row">
    <div>
      <label for="event-category">Kategori</label><select id="event-category" bind:value={category}
        >{#each ['Teater', 'Musik', 'Komedi', 'Lainnya'] as c}<option>{c}</option>{/each}</select
      >
    </div>
    <div>
      <label for="event-city">Kota</label><input
        id="event-city"
        bind:value={city}
        required
        maxlength="60"
      />
    </div>
  </div>
  <div>
    <label for="event-venue">Nama venue</label><input
      id="event-venue"
      bind:value={venue}
      required
      maxlength="150"
    />
  </div>
  <div>
    <label for="event-address">Alamat</label><input
      id="event-address"
      bind:value={address}
      maxlength="400"
    />
  </div>
  <div>
    <label for="event-maps">Tautan Google Maps (opsional)</label><input
      id="event-maps"
      type="url"
      bind:value={maps}
      placeholder="https://maps.app.goo.gl/..."
    />
  </div>
  <div class="form-row">
    <div>
      <label for="event-duration">Durasi</label><input
        id="event-duration"
        bind:value={duration}
        maxlength="80"
      />
    </div>
    <div>
      <label for="event-age">Batas usia</label><input
        id="event-age"
        bind:value={age}
        maxlength="60"
      />
    </div>
  </div>
  <div>
    <label for="event-language">Bahasa</label><input
      id="event-language"
      bind:value={language}
      maxlength="80"
    />
  </div>
  <div class="divider"></div>
  <section class="stack">
    <h2>Jadwal & kuota tiket</h2>
    {#each sessions as s, i}<fieldset class="card pad">
        <legend>Sesi {i + 1}</legend><label for={'session-label-' + i}>Nama sesi</label><input
          id={'session-label-' + i}
          bind:value={s.label}
          required
          maxlength="100"
        />
        <div class="form-row">
          <div>
            <label for={'session-start-' + i}>Mulai · WIB</label><input
              id={'session-start-' + i}
              type="datetime-local"
              bind:value={s.starts_at}
              required
            />
          </div>
          <div>
            <label for={'session-end-' + i}>Selesai · WIB</label><input
              id={'session-end-' + i}
              type="datetime-local"
              bind:value={s.ends_at}
              required
            />
          </div>
        </div>
        <div class="form-row">
          <div>
            <label for={'session-price-' + i}>Harga · Rp</label><input
              id={'session-price-' + i}
              type="number"
              min="0"
              max="100000000"
              bind:value={s.price}
              required
            />
          </div>
          <div>
            <label for={'session-capacity-' + i}>Kapasitas</label><input
              id={'session-capacity-' + i}
              type="number"
              min="1"
              max="100000"
              bind:value={s.capacity}
              required
            />
          </div>
        </div>
        {#if sessions.length > 1}<button
            type="button"
            class="action"
            onclick={() => (sessions = sessions.filter((_, n) => n !== i))}>Hapus sesi</button
          >{/if}
      </fieldset>{/each}<button
      class="button"
      type="button"
      onclick={() =>
        (sessions = [
          ...sessions,
          { id: '', label: '', starts_at: '', ends_at: '', price: 0, capacity: 60 },
        ])}>Tambah sesi</button
    >
    <p class="hint">
      Jadwal dan harga sesi yang sudah memiliki pesanan tidak bisa diubah. Kapasitas bisa ditambah.
    </p>
  </section>
  <div class="divider"></div>
  <section class="stack">
    <h2>Pengisi acara</h2>
    {#each lineup as p, i}<div class="card pad">
        <div class="form-row">
          <div>
            <label for={'performer-name-' + i}>Nama</label><input
              id={'performer-name-' + i}
              bind:value={p.name}
              maxlength="80"
              required
            />
          </div>
          <div>
            <label for={'performer-role-' + i}>Peran</label><input
              id={'performer-role-' + i}
              bind:value={p.role}
              maxlength="80"
              required
              placeholder="Aktor, sutradara, band…"
            />
          </div>
        </div>
        <label for={'performer-photo-' + i}>Foto</label><input
          id={'performer-photo-' + i}
          type="file"
          accept="image/jpeg,image/png,image/webp"
          onchange={(e) => photo(e, 'performer', i)}
        />{#if p.photo}<img class="avatar large" src={p.photo} alt={p.name} /><button
            type="button"
            class="action"
            onclick={() => (p.photo = '')}>Hapus foto</button
          >{/if}<button
          type="button"
          class="action"
          onclick={() => (lineup = lineup.filter((_, n) => n !== i))}>Hapus pengisi</button
        >
      </div>{/each}<button
      class="button"
      type="button"
      onclick={() => (lineup = [...lineup, { name: '', role: '', photo: '' }])}
      >Tambah pengisi</button
    >
  </section>
  <details>
    <summary>Trailer, flyer, dan denah venue (opsional)</summary>
    <div class="form" style="margin-top:18px">
      <div>
        <label for="trailer">Trailer YouTube / MP4 / WebM</label><input
          id="trailer"
          type="url"
          bind:value={trailer}
          placeholder="https://..."
        />
        <p class="hint">
          Jika kosong, flyer ditampilkan. Jika keduanya kosong, section media disembunyikan.
        </p>
      </div>
      <div>
        <label for="flyer">Flyer landscape</label><input
          id="flyer"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          onchange={(e) => photo(e, 'flyer')}
        />{#if flyer}<img class="attachment" src={flyer} alt="Flyer event" /><button
            class="action"
            type="button"
            onclick={() => (flyer = '')}>Hapus flyer</button
          >{/if}
      </div>
      <div>
        <label for="layout">Denah venue / tempat duduk</label><input
          id="layout"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          onchange={(e) => photo(e, 'layout')}
        />{#if layout}<img class="attachment" src={layout} alt="Denah venue" /><button
            class="action"
            type="button"
            onclick={() => (layout = '')}>Hapus denah</button
          >{/if}
      </div>
    </div>
  </details>
  <label class="inline-label"
    ><input type="checkbox" bind:checked={published} />Terbitkan event agar tampil di agenda publik.</label
  >{#if errorMessage}<p class="note error-note" role="alert">{errorMessage}</p>{/if}<button
    class="button primary"
    disabled={busy}>{busy ? 'Menyimpan…' : 'Simpan event'}</button
  >
</form>
